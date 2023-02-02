package concentrated_liquidity

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/osmoutils/accum"
	"github.com/osmosis-labs/osmosis/v14/x/concentrated-liquidity/internal/math"
	cltypes "github.com/osmosis-labs/osmosis/v14/x/concentrated-liquidity/types"
)

const (
	feeAccumPrefix = "fee"
	keySeparator   = "/"
	uintBase       = 10
)

var (
	emptyCoins = sdk.DecCoins(nil)
)

// createFeeAccumulator creates an accumulator object in the store using the given poolId.
// The accumulator is initialized with the default(zero) values.
func (k Keeper) createFeeAccumulator(ctx sdk.Context, poolId uint64) error {
	err := accum.MakeAccumulator(ctx.KVStore(k.storeKey), getFeeAccumulatorName(poolId))
	if err != nil {
		return err
	}
	return nil
}

// nolint: unused
// getFeeAccumulator gets the fee accumulator object using the given poolOd
// returns error if accumulator for the given poolId does not exist.
func (k Keeper) getFeeAccumulator(ctx sdk.Context, poolId uint64) (accum.AccumulatorObject, error) {
	acc, err := accum.GetAccumulator(ctx.KVStore(k.storeKey), getFeeAccumulatorName(poolId))
	if err != nil {
		return accum.AccumulatorObject{}, err
	}

	return acc, nil
}

// chargeFee charges the given fee on the pool with the given id by updating
// the internal per-pool accumulator that tracks fee growth per one unit of
// liquidity. Returns error if fails to get accumulator.
// nolint: unused
func (k Keeper) chargeFee(ctx sdk.Context, poolId uint64, feeUpdate sdk.DecCoin) error {
	feeAccumulator, err := k.getFeeAccumulator(ctx, poolId)
	if err != nil {
		return err
	}

	feeAccumulator.AddToAccumulator(sdk.NewDecCoins(feeUpdate))

	return nil
}

// initializeFeeAccumulatorPosition initializes the pool fee accumulator with zero liquidity delta
// and zero value for the accumulator.
// Returns nil on success. Returns error if:
// - fails to get an accumulator for a given poold id
// - attempts to re-initialize an existing fee accumulator liqudity position
// - fails to create a position
func (k Keeper) initializeFeeAccumulatorPosition(ctx sdk.Context, poolId uint64, owner sdk.AccAddress, lowerTick, upperTick int64) error {
	// get fee accumulator for the pool
	feeAccumulator, err := k.getFeeAccumulator(ctx, poolId)
	if err != nil {
		return err
	}

	positionKey := formatPositionAccumulatorKey(poolId, owner, lowerTick, upperTick)

	hasPosition, err := feeAccumulator.HasPosition(positionKey)
	if err != nil {
		return err
	}

	// assure that existing position has zero liquidity
	if hasPosition {
		return fmt.Errorf("attempted to re-initialize fee accumulator position (%s) with non-zero liquidity", positionKey)
	}

	// initialize the owner's position with liquidity Delta and zero accumulator value
	if err := feeAccumulator.NewPosition(positionKey, sdk.ZeroDec(), nil); err != nil {
		return err
	}

	return nil
}

// updateFeeAccumulatorPosition updates the fee accumulator position for a given pool, owner, and tick range.
// It retrieves the current fee growth outside of the given tick range and updates the position's accumulator
// with the provided liquidity delta and the retrieved fee growth outside.
func (k Keeper) updateFeeAccumulatorPosition(ctx sdk.Context, poolId uint64, owner sdk.AccAddress, liquidityDelta sdk.Dec, lowerTick int64, upperTick int64) error {
	feeGrowthOutside, err := k.getFeeGrowthOutside(ctx, poolId, lowerTick, upperTick)
	if err != nil {
		return err
	}

	feeAccumulator, err := k.getFeeAccumulator(ctx, poolId)
	if err != nil {
		return err
	}

	// replace position's accumulator with the updated liquidity and the feeGrowthOutside
	err = feeAccumulator.UpdatePositionCustomAcc(
		formatPositionAccumulatorKey(poolId, owner, lowerTick, upperTick),
		liquidityDelta,
		feeGrowthOutside)
	if err != nil {
		return err
	}

	return nil
}

// getFeeGrowthOutside returns fee growth upper tick - fee growth lower tick
// nolint: unused
func (k Keeper) getFeeGrowthOutside(ctx sdk.Context, poolId uint64, lowerTick, upperTick int64) (sdk.DecCoins, error) {
	pool, err := k.getPoolById(ctx, poolId)
	if err != nil {
		return sdk.DecCoins{}, err
	}
	currentTick := pool.GetCurrentTick().Int64()

	// get lower, upper tick info
	lowerTickInfo, err := k.getTickInfo(ctx, poolId, lowerTick)
	if err != nil {
		return sdk.DecCoins{}, err
	}
	upperTickInfo, err := k.getTickInfo(ctx, poolId, upperTick)
	if err != nil {
		return sdk.DecCoins{}, err
	}

	poolFeeAccumulator, err := k.getFeeAccumulator(ctx, poolId)
	if err != nil {
		return sdk.DecCoins{}, err
	}
	poolFeeGrowth := poolFeeAccumulator.GetValue()

	// calculate fee growth for upper tick and lower tick
	feeGrowthAboveUpperTick := calculateFeeGrowth(upperTick, upperTickInfo.FeeGrowthOutside, currentTick, poolFeeGrowth, true)
	feeGrowthBelowLowerTick := calculateFeeGrowth(lowerTick, lowerTickInfo.FeeGrowthOutside, currentTick, poolFeeGrowth, false)

	return feeGrowthAboveUpperTick.Add(feeGrowthBelowLowerTick...), nil
}

// getInitialFeeGrowthOutsideForTick returns the initial value of fee growth outside for a given tick.
// This value depends on the tick's location relative to the current tick.
//
// feeGrowthOutside =
// { feeGrowthGlobal current tick >= tick }
// { 0               current tick <  tick }
//
// The value is chosen as if all of the fees earned to date had occurrd below the tick.
// Returns error if the pool with the given id does exist or if fails to get the fee accumulator.
func (k Keeper) getInitialFeeGrowthOutsideForTick(ctx sdk.Context, poolId uint64, tick int64) (sdk.DecCoins, error) {
	pool, err := k.getPoolById(ctx, poolId)
	if err != nil {
		return sdk.DecCoins{}, err
	}

	currentTick := pool.GetCurrentTick().Int64()
	if currentTick >= tick {
		feeAccumulator, err := k.getFeeAccumulator(ctx, poolId)
		if err != nil {
			return sdk.DecCoins{}, err
		}
		return feeAccumulator.GetValue(), nil
	}

	return emptyCoins, nil
}

// collectFees collects fees from the fee accumulator for the position given by pool id, owner, lower tick and upper tick.
// Upon successful collection, it bank sends the fees from the pool address to the owner and returns the collected coins.
// Returns error if:
// - pool with the given id does not exist
// - position given by pool id, owner, lower tick and upper tick does not exist
// - other internal database or math errors.
// nolint: unused
func (k Keeper) collectFees(ctx sdk.Context, poolId uint64, owner sdk.AccAddress, lowerTick int64, upperTick int64) (sdk.Coins, error) {
	feeAccumulator, err := k.getFeeAccumulator(ctx, poolId)
	if err != nil {
		return sdk.Coins{}, err
	}

	positionKey := formatPositionAccumulatorKey(poolId, owner, lowerTick, upperTick)

	hasPosition, err := feeAccumulator.HasPosition(positionKey)
	if err != nil {
		return sdk.Coins{}, err
	}

	if !hasPosition {
		return sdk.Coins{}, cltypes.PositionNotFoundError{PoolId: poolId, LowerTick: lowerTick, UpperTick: upperTick}
	}

	// compute fee growth outside of the range between lower tick and upper tick.
	feeGrowthOutside, err := k.getFeeGrowthOutside(ctx, poolId, lowerTick, upperTick)
	if err != nil {
		return sdk.Coins{}, err
	}

	// We need to update the position's accumulator to the current fee growth outside
	// before we claim rewards.
	if err := feeAccumulator.SetPositionCustomAcc(positionKey, feeGrowthOutside); err != nil {
		return sdk.Coins{}, err
	}

	// claim fees.
	feesClaimed, err := feeAccumulator.ClaimRewardsCustomAcc(positionKey, feeGrowthOutside)
	if err != nil {
		return sdk.Coins{}, err
	}

	// Once we have iterated through all the positions, we do a single bank send from the pool to the owner.
	pool, err := k.getPoolById(ctx, poolId)
	if err != nil {
		return sdk.Coins{}, err
	}
	if err := k.bankKeeper.SendCoins(ctx, pool.GetAddress(), owner, feesClaimed); err != nil {
		return sdk.Coins{}, err
	}
	return feesClaimed, nil
}

func getFeeAccumulatorName(poolId uint64) string {
	poolIdStr := strconv.FormatUint(poolId, uintBase)
	return strings.Join([]string{feeAccumPrefix, poolIdStr}, "/")
}

// calculateFeeGrowth for the given targetTicks.
// If calculating fee growth for an upper tick, we consider the following two cases
// 1. currentTick >= upperTick: If current Tick is GTE than the upper Tick, the fee growth would be pool fee growth - uppertick's fee growth outside
// 2. currentTick < upperTick: If current tick is smaller than upper tick, fee growth would be the upper tick's fee growth outside
// this goes vice versa for calculating fee growth for lower tick.
// nolint: unused
func calculateFeeGrowth(targetTick int64, feeGrowthOutside sdk.DecCoins, currentTick int64, feesGrowthGlobal sdk.DecCoins, isUpperTick bool) sdk.DecCoins {
	if (isUpperTick && currentTick >= targetTick) || (!isUpperTick && currentTick < targetTick) {
		return feesGrowthGlobal.Sub(feeGrowthOutside)
	}
	return feeGrowthOutside
}

// formatPositionAccumulatorKey formats the position accumulator key prefixed by pool id, owner, lower tick
// and upper tick with a key separator in-between.
// nolint: unused
func formatPositionAccumulatorKey(poolId uint64, owner sdk.AccAddress, lowerTick, upperTick int64) string {
	return strings.Join([]string{strconv.FormatUint(poolId, uintBase), owner.String(), strconv.FormatInt(lowerTick, uintBase), strconv.FormatInt(upperTick, uintBase)}, keySeparator)
}

// computeFeeChargePerSwapStepOutGivenIn returns the total fee charge per swap step given the parameters.
// Assumes swapping for token out given token in.
// - currentSqrtPrice the sqrt price at which the swap step begins.
// - nextTickSqrtPrice the next tick's sqrt price.
// - sqrtPriceLimit the sqrt price corresponding to the sqrt of the price representing price impact protection.
// - amountIn the amount of token in to be consumed during the swap step
// - amountSpecifiedRemaining is the total remaining amount of token in that needs to be consumed to complete the swap.
// - swapFee the swap fee to be charged.
//
// If swap fee is negative, it panics.
// If swap fee is 0, returns 0. Otherwise, computes and returns the fee charge per step.
// TODO: test this function.
func computeFeeChargePerSwapStepOutGivenIn(currentSqrtPrice, nextTickSqrtPrice, sqrtPriceLimit, amountIn, amountSpecifiedRemaining, swapFee sdk.Dec) sdk.Dec {
	feeChargeTotal := sdk.ZeroDec()

	if swapFee.IsNegative() {
		// This should never happen but is added as a defense-in-depth measure.
		panic(fmt.Errorf("swap fee must be non-negative, was (%s)", swapFee))
	}

	if swapFee.IsZero() {
		return feeChargeTotal
	}

	// 1. The current tick does not have enough liqudity to fulfill the swap.
	didReachNextSqrtPrice := currentSqrtPrice.Equal(nextTickSqrtPrice)
	// 2. The next sqrt price was not reached due to price impact protection.
	isPriceImpactProtection := currentSqrtPrice.Equal(sqrtPriceLimit)

	// In both cases, charge fee on the full amount that the tick
	// originally had.
	if didReachNextSqrtPrice || isPriceImpactProtection {
		// Multiply with rounding up to avoid under charging fees.
		feeChargeTotal = math.MulRoundUp(amountIn, swapFee)
	} else {
		// Otherwise, the current tick had enough liquidity to fulfill the swap
		// In that case, the fee is the difference between
		// the amount needed to fulfill and the actual amount we ended up charging.
		feeChargeTotal = amountSpecifiedRemaining.Sub(amountIn)
	}

	return feeChargeTotal
}
