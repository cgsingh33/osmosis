package concentrated_liquidity_test

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/osmosis-labs/osmosis/v13/x/concentrated-liquidity/types"
)

type lpTest struct {
	poolId          uint64
	owner           sdk.AccAddress
	currentTick     sdk.Int
	lowerTick       int64
	upperTick       int64
	currentSqrtP    sdk.Dec
	amount0Desired  sdk.Int
	amount0Minimum  sdk.Int
	amount0Expected sdk.Int
	amount1Desired  sdk.Int
	amount1Minimum  sdk.Int
	amount1Expected sdk.Int
	liquidityAmount sdk.Dec
	tickSpacing     uint64
	expectedError   error
}

var (
	baseCase = &lpTest{
		poolId:          1,
		currentTick:     DefaultCurrTick,
		lowerTick:       DefaultLowerTick,
		upperTick:       DefaultUpperTick,
		currentSqrtP:    DefaultCurrSqrtPrice,
		amount0Desired:  DefaultAmt0,
		amount0Minimum:  sdk.ZeroInt(),
		amount0Expected: DefaultAmt0Expected,
		amount1Desired:  DefaultAmt1,
		amount1Minimum:  sdk.ZeroInt(),
		amount1Expected: DefaultAmt1Expected,
		liquidityAmount: DefaultLiquidityAmt,
		tickSpacing:     DefaultTickSpacing,
	}
)

func (s *KeeperTestSuite) TestCreatePosition() {
	tests := map[string]lpTest{
		"base case": {},
		"error: non-existent pool": {
			poolId:        2,
			expectedError: types.PoolNotFoundError{PoolId: 2},
		},
		"error: lower tick out of bounds": {
			lowerTick:     types.MinTick - 1,
			expectedError: types.InvalidTickError{Tick: types.MinTick - 1, IsLower: true},
		},
		"error: upper tick out of bounds": {
			upperTick:     types.MaxTick + 1,
			expectedError: types.InvalidTickError{Tick: types.MaxTick + 1, IsLower: false},
		},
		"error: upper tick is below the lower tick, but both are in bounds": {
			lowerTick:     50,
			upperTick:     40,
			expectedError: types.InvalidLowerUpperTickError{LowerTick: 50, UpperTick: 40},
		},
		"error: amount of token 0 is smaller than minimum; should fail and not update state": {
			amount0Minimum: baseCase.amount0Expected.Mul(sdk.NewInt(2)),
			expectedError:  types.InsufficientLiquidityCreatedError{Actual: baseCase.amount0Expected, Minimum: baseCase.amount0Expected.Mul(sdk.NewInt(2)), IsTokenZero: true},
		},
		"error: amount of token 1 is smaller than minimum; should fail and not update state": {
			amount1Minimum: baseCase.amount1Expected.Mul(sdk.NewInt(2)),
			expectedError:  types.InsufficientLiquidityCreatedError{Actual: baseCase.amount1Expected, Minimum: baseCase.amount1Expected.Mul(sdk.NewInt(2))},
		},
		"error: amount of token 1 is smaller than minimum; should fail and not update state test": {
			amount0Desired: sdk.ZeroInt(),
			amount1Desired: sdk.ZeroInt(),
			expectedError:  errors.New("liquidityDelta calculated equals zero"),
		},
		"create a position with non default tick spacing (10) with ticks that fall into tick spacing requirements": {
			lowerTick:       int64(84220),
			upperTick:       int64(86130),
			amount0Expected: sdk.NewInt(997568),
			liquidityAmount: sdk.MustNewDecFromStr("1514719247.706987476085061790"),
			tickSpacing:     10,
		},
		"error: attempt to use and upper and lower tick that are not divisible by tick spacing": {
			tickSpacing:   10,
			expectedError: types.TickSpacingError{TickSpacing: 10, LowerTick: DefaultLowerTick, UpperTick: DefaultUpperTick},
		},
		// TODO: add more tests
		// - custom hand-picked values
		// - think of overflows
		// - think of truncations
	}

	for name, tc := range tests {
		s.Run(name, func() {
			tc := tc
			s.SetupTest()

			// Merge tc with baseCase and update tc to the merged result. This is done to reduce the amount of boilerplate in test cases.
			baseConfigCopy := *baseCase
			mergeConfigs(&baseConfigCopy, &tc)
			tc = baseConfigCopy

			// Create a CL pool with custom tickSpacing
			pool, err := s.App.ConcentratedLiquidityKeeper.CreateNewConcentratedLiquidityPool(s.Ctx, 1, ETH, USDC, tc.tickSpacing)
			s.Require().NoError(err)

			// Fund test account and create the desired position
			s.FundAcc(s.TestAccs[0], sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))

			// Note user and pool account balances before create position is called
			userBalancePrePositionCreation := s.App.BankKeeper.GetAllBalances(s.Ctx, s.TestAccs[0])
			poolBalancePrePositionCreation := s.App.BankKeeper.GetAllBalances(s.Ctx, pool.GetAddress())

			asset0, asset1, liquidityCreated, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(s.Ctx, tc.poolId, s.TestAccs[0], tc.amount0Desired, tc.amount1Desired, tc.amount0Minimum, tc.amount1Minimum, tc.lowerTick, tc.upperTick)

			// Note user and pool account balances to compare after create position is called
			userBalancePostPositionCreation := s.App.BankKeeper.GetAllBalances(s.Ctx, s.TestAccs[0])
			poolBalancePostPositionCreation := s.App.BankKeeper.GetAllBalances(s.Ctx, pool.GetAddress())

			// If we expect an error, make sure:
			// - some error was emitted
			// - asset0 and asset1 that were calculated from create position is nil
			// - the error emitted was the expected error
			// - account balances are untouched
			if tc.expectedError != nil {
				s.Require().Error(err)
				s.Require().Equal(asset0, sdk.Int{})
				s.Require().Equal(asset1, sdk.Int{})
				s.Require().ErrorAs(err, &tc.expectedError)

				// Check account balances
				s.Require().Equal(userBalancePrePositionCreation.String(), userBalancePostPositionCreation.String())
				s.Require().Equal(poolBalancePrePositionCreation.String(), poolBalancePostPositionCreation.String())

				// Redundantly ensure that position was not created
				position, err := s.App.ConcentratedLiquidityKeeper.GetPosition(s.Ctx, tc.poolId, s.TestAccs[0], tc.lowerTick, tc.upperTick)
				s.Require().Error(err)
				s.Require().ErrorAs(err, &types.PositionNotFoundError{PoolId: tc.poolId, LowerTick: tc.lowerTick, UpperTick: tc.upperTick})
				s.Require().Nil(position)
				return
			}

			// Else, check that we had no error from creating the position, and that the liquidity and assets that were returned are expected
			s.Require().NoError(err)
			s.Require().Equal(tc.amount0Expected.String(), asset0.String())
			s.Require().Equal(tc.amount1Expected.String(), asset1.String())
			s.Require().Equal(tc.liquidityAmount.String(), liquidityCreated.String())

			// Check account balances
			s.Require().Equal(userBalancePrePositionCreation.Sub(sdk.NewCoins(sdk.NewCoin(ETH, asset0), (sdk.NewCoin(USDC, asset1)))).String(), userBalancePostPositionCreation.String())
			s.Require().Equal(poolBalancePrePositionCreation.Add(sdk.NewCoin(ETH, asset0), (sdk.NewCoin(USDC, asset1))).String(), poolBalancePostPositionCreation.String())

			// Check position state
			s.validatePositionUpdate(s.Ctx, tc.poolId, s.TestAccs[0], tc.lowerTick, tc.upperTick, tc.liquidityAmount)

			// Check tick state
			s.validateTickUpdates(s.Ctx, tc.poolId, s.TestAccs[0], tc.lowerTick, tc.upperTick, tc.liquidityAmount)
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawPosition() {
	tests := map[string]struct {
		setupConfig *lpTest
		// when this is set, it ovewrites the setupConfig
		// and gives the overwritten configuration to
		// the system under test.
		sutConfigOverwrite *lpTest
	}{
		"base case: withdraw full liquidity amount": {
			// setup parameters for creating a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				amount0Expected: baseCase.amount0Expected, // 0.998587 eth
				amount1Expected: baseCase.amount1Expected, // 5000 usdc
			},
		},
		"withdraw partial liquidity amount": {
			// setup parameters for creating a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				liquidityAmount: baseCase.liquidityAmount.QuoInt64(2),

				amount0Expected: baseCase.amount0Expected.QuoRaw(2),                   // 0.4992935 / 2 eth
				amount1Expected: baseCase.amount1Expected.QuoRaw(2).Sub(sdk.OneInt()), // 2499 usdc, one is lost due to truncation
			},
		},
		"error: no position created": {
			// setup parameters for creation a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				lowerTick:     -1, // valid tick at which no position exists
				expectedError: types.PositionNotFoundError{PoolId: 1, LowerTick: -1, UpperTick: 86129},
			},
		},
		"error: pool id for pool that does not exist": {
			// setup parameters for creating a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				poolId:        2, // does not exist
				expectedError: types.PoolNotFoundError{PoolId: 2},
			},
		},
		"error: upper tick out of bounds": {
			// setup parameters for creating a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				upperTick:     types.MaxTick + 1, // invalid tick
				expectedError: types.InvalidTickError{Tick: types.MaxTick + 1, IsLower: false},
			},
		},
		"error: lower tick out of bounds": {
			// setup parameters for creating a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				lowerTick:     types.MinTick - 1, // invalid tick
				expectedError: types.InvalidTickError{Tick: types.MinTick - 1, IsLower: true},
			},
		},
		"error: insufficient liquidity": {
			// setup parameters for creating a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				liquidityAmount: baseCase.liquidityAmount.Add(sdk.OneDec()), // 1 more than available
				expectedError:   types.InsufficientLiquidityError{Actual: baseCase.liquidityAmount.Add(sdk.OneDec()), Available: baseCase.liquidityAmount},
			},
		},
		"error: upper tick is below the lower tick, but both are in bounds": {
			// setup parameters for creating a pool and position.
			setupConfig: baseCase,

			// system under test parameters
			// for withdrawing a position.
			sutConfigOverwrite: &lpTest{
				lowerTick:     50,
				upperTick:     40,
				expectedError: types.InvalidLowerUpperTickError{LowerTick: 50, UpperTick: 40},
			},
		},
		// TODO: test with custom amounts that potentially lead to truncations.
	}

	for name, tc := range tests {
		s.Run(name, func() {
			// Setup.
			s.SetupTest()

			var (
				ctx                         = s.Ctx
				concentratedLiquidityKeeper = s.App.ConcentratedLiquidityKeeper
				liquidityCreated            = sdk.ZeroDec()
				owner                       = s.TestAccs[0]
				tc                          = tc
				config                      = *tc.setupConfig
				sutConfigOverwrite          = *tc.sutConfigOverwrite
			)

			// If a setupConfig is provided, use it to create a pool and position.
			if tc.setupConfig != nil {
				s.PrepareDefaultPool(ctx)
				var err error
				s.FundAcc(s.TestAccs[0], sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))
				_, _, liquidityCreated, err = concentratedLiquidityKeeper.CreatePosition(ctx, config.poolId, owner, config.amount0Desired, config.amount1Desired, sdk.ZeroInt(), sdk.ZeroInt(), config.lowerTick, config.upperTick)
				s.Require().NoError(err)
			}

			// If specific configs are provided in the test case, overwrite the config with those values.
			mergeConfigs(&config, &sutConfigOverwrite)

			// System under test.
			amtDenom0, amtDenom1, err := concentratedLiquidityKeeper.WithdrawPosition(ctx, config.poolId, owner, config.lowerTick, config.upperTick, config.liquidityAmount)

			if config.expectedError != nil {
				s.Require().Error(err)
				s.Require().Equal(amtDenom0, sdk.Int{})
				s.Require().Equal(amtDenom1, sdk.Int{})
				s.Require().ErrorAs(err, &config.expectedError)
				return
			}

			s.Require().NoError(err)
			s.Require().Equal(config.amount0Expected.String(), amtDenom0.String())
			s.Require().Equal(config.amount1Expected.String(), amtDenom1.String())

			// Determine the liquidity expected to remain after the withdraw.
			expectedRemainingLiquidity := liquidityCreated.Sub(config.liquidityAmount)

			// Check that the position was updated.
			s.validatePositionUpdate(ctx, config.poolId, owner, config.lowerTick, config.upperTick, expectedRemainingLiquidity)

			// check tick state
			s.validateTickUpdates(ctx, config.poolId, owner, config.lowerTick, config.upperTick, expectedRemainingLiquidity)
		})
	}
}

// mergeConfigs merges every desired non-zero field from overwrite
// into dst. dst is mutated due to being a pointer.
func mergeConfigs(dst *lpTest, overwrite *lpTest) {
	if overwrite != nil {
		if overwrite.poolId != 0 {
			dst.poolId = overwrite.poolId
		}
		if overwrite.lowerTick != 0 {
			dst.lowerTick = overwrite.lowerTick
		}
		if overwrite.upperTick != 0 {
			dst.upperTick = overwrite.upperTick
		}
		if !overwrite.liquidityAmount.IsNil() {
			dst.liquidityAmount = overwrite.liquidityAmount
		}
		if !overwrite.amount0Minimum.IsNil() {
			dst.amount0Minimum = overwrite.amount0Minimum
		}
		if !overwrite.amount0Desired.IsNil() {
			dst.amount0Desired = overwrite.amount0Desired
		}
		if !overwrite.amount0Expected.IsNil() {
			dst.amount0Expected = overwrite.amount0Expected
		}
		if !overwrite.amount1Minimum.IsNil() {
			dst.amount1Minimum = overwrite.amount1Minimum
		}
		if !overwrite.amount1Desired.IsNil() {
			dst.amount1Desired = overwrite.amount1Desired
		}
		if !overwrite.amount1Expected.IsNil() {
			dst.amount1Expected = overwrite.amount1Expected
		}
		if overwrite.expectedError != nil {
			dst.expectedError = overwrite.expectedError
		}
		if overwrite.tickSpacing != 0 {
			dst.tickSpacing = overwrite.tickSpacing
		}
	}
}

func (s *KeeperTestSuite) TestSendCoinsBetweenPoolAndUser() {
	type sendTest struct {
		coin0       sdk.Coin
		coin1       sdk.Coin
		poolToUser  bool
		expectError bool
	}
	tests := map[string]sendTest{
		"asset0 and asset1 are positive, position creation (user to pool)": {
			coin0: sdk.NewCoin("eth", sdk.NewInt(1000000)),
			coin1: sdk.NewCoin("usdc", sdk.NewInt(1000000)),
		},
		"only asset0 is positive, position creation (user to pool)": {
			coin0: sdk.NewCoin("eth", sdk.NewInt(1000000)),
			coin1: sdk.NewCoin("usdc", sdk.NewInt(0)),
		},
		"only asset1 is positive, position creation (user to pool)": {
			coin0: sdk.NewCoin("eth", sdk.NewInt(0)),
			coin1: sdk.NewCoin("usdc", sdk.NewInt(1000000)),
		},
		"only asset0 is greater than sender has, position creation (user to pool)": {
			coin0:       sdk.NewCoin("eth", sdk.NewInt(100000000000000)),
			coin1:       sdk.NewCoin("usdc", sdk.NewInt(1000000)),
			expectError: true,
		},
		"only asset1 is greater than sender has, position creation (user to pool)": {
			coin0:       sdk.NewCoin("eth", sdk.NewInt(1000000)),
			coin1:       sdk.NewCoin("usdc", sdk.NewInt(100000000000000)),
			expectError: true,
		},
		"asset0 and asset1 are positive, withdraw (pool to user)": {
			coin0:      sdk.NewCoin("eth", sdk.NewInt(1000000)),
			coin1:      sdk.NewCoin("usdc", sdk.NewInt(1000000)),
			poolToUser: true,
		},
		"only asset0 is positive, withdraw (pool to user)": {
			coin0:      sdk.NewCoin("eth", sdk.NewInt(1000000)),
			coin1:      sdk.NewCoin("usdc", sdk.NewInt(0)),
			poolToUser: true,
		},
		"only asset1 is positive, withdraw (pool to user)": {
			coin0:      sdk.NewCoin("eth", sdk.NewInt(0)),
			coin1:      sdk.NewCoin("usdc", sdk.NewInt(1000000)),
			poolToUser: true,
		},
		"only asset0 is greater than sender has, withdraw (pool to user)": {
			coin0:       sdk.NewCoin("eth", sdk.NewInt(100000000000000)),
			coin1:       sdk.NewCoin("usdc", sdk.NewInt(1000000)),
			poolToUser:  true,
			expectError: true,
		},
		"only asset1 is greater than sender has, withdraw (pool to user)": {
			coin0:       sdk.NewCoin("eth", sdk.NewInt(1000000)),
			coin1:       sdk.NewCoin("usdc", sdk.NewInt(100000000000000)),
			poolToUser:  true,
			expectError: true,
		},
	}

	for name, tc := range tests {
		s.Run(name, func() {
			tc := tc
			s.SetupTest()

			// create a CL pool
			s.PrepareDefaultPool(s.Ctx)

			// store pool interface
			poolI, err := s.App.ConcentratedLiquidityKeeper.GetPoolById(s.Ctx, 1)
			s.Require().NoError(err)
			concentratedPool := poolI.(types.ConcentratedPoolExtension)

			// fund pool address and user address
			s.FundAcc(poolI.GetAddress(), sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))
			s.FundAcc(s.TestAccs[0], sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))

			// set sender and receiver based on test case
			sender := s.TestAccs[0]
			receiver := concentratedPool.GetAddress()
			if tc.poolToUser {
				sender = concentratedPool.GetAddress()
				receiver = s.TestAccs[0]
			}

			// store pre send balance of sender and receiver
			preSendBalanceSender := s.App.BankKeeper.GetAllBalances(s.Ctx, sender)
			preSendBalanceReceiver := s.App.BankKeeper.GetAllBalances(s.Ctx, receiver)

			// system under test
			err = s.App.ConcentratedLiquidityKeeper.SendCoinsBetweenPoolAndUser(s.Ctx, concentratedPool.GetToken0(), concentratedPool.GetToken1(), tc.coin0.Amount, tc.coin1.Amount, sender, receiver)

			// store post send balance of sender and receiver
			postSendBalanceSender := s.App.BankKeeper.GetAllBalances(s.Ctx, sender)
			postSendBalanceReceiver := s.App.BankKeeper.GetAllBalances(s.Ctx, receiver)

			// check error if expected
			if tc.expectError {
				s.Require().Error(err)
				return
			}

			// otherwise, ensure balances are added/deducted appropriately
			expectedPostSendBalanceSender := preSendBalanceSender.Sub(sdk.NewCoins(tc.coin0, tc.coin1))
			expectedPostSendBalanceReceiver := preSendBalanceReceiver.Add(tc.coin0, tc.coin1)

			s.Require().NoError(err)
			s.Require().Equal(expectedPostSendBalanceSender.String(), postSendBalanceSender.String())
			s.Require().Equal(expectedPostSendBalanceReceiver.String(), postSendBalanceReceiver.String())
		})
	}
}

func (s *KeeperTestSuite) TestIsInitialPosition() {
	type sendTest struct {
		initialSqrtPrice sdk.Dec
		initialTick      sdk.Int
		expectedResponse bool
	}
	tests := map[string]sendTest{
		"happy path: is initial position": {
			initialSqrtPrice: sdk.ZeroDec(),
			initialTick:      sdk.ZeroInt(),
			expectedResponse: true,
		},
		"happy path: is not initial position": {
			initialSqrtPrice: DefaultCurrSqrtPrice,
			initialTick:      DefaultCurrTick,
			expectedResponse: false,
		},
		"tick is zero but initialSqrtPrice is not, should not be detected as initial potion": {
			initialSqrtPrice: DefaultCurrSqrtPrice,
			initialTick:      sdk.ZeroInt(),
			expectedResponse: false,
		},
		"initialSqrtPrice is zero but tick is not, should not be detected as initial position (should not happen)": {
			initialSqrtPrice: sdk.ZeroDec(),
			initialTick:      DefaultCurrTick,
			expectedResponse: false,
		},
	}

	for name, tc := range tests {
		s.Run(name, func() {
			// System under test
			if s.App.ConcentratedLiquidityKeeper.IsInitialPosition(tc.initialSqrtPrice, tc.initialTick) {
				// If we expect the response to be true, then we should check that it is true
				s.Require().True(tc.expectedResponse)
			} else {
				// Else, we should check that it is false
				s.Require().False(tc.expectedResponse)
			}
		})
	}
}
