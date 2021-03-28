package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/gookit/color"

	coindb "github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/CoinDB"
	model "github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/DataModels"
)

func getChallengeRoute(c *gin.Context) {
	c.JSON(200, &ChallengeApiResponse{
		Reward:           *reward,
		CurrentChallenge: GetChallenge(),
		Difficulty:       *difficulty,
	})
}

func mineStatistics(c *gin.Context) {
	recentMiners, err := (&coindb.Transaction{Senderid: 1}).LastBySender(20)
	if err != nil && err != pg.ErrNoRows {
		color.Error.Printf("Error in mineStatistics(): %v\n", err.Error())
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.JSON(200, model.MineStatisticsResponse{
		Difficulty: model.DifficultyStatistics{
			Reward:        *reward,
			RewardTarget:  *rewardtarget,
			TargetMinutes: *targetminutes,

			Minimum: *mindifficulty,
			Maximum: *maxdifficulty,
			Current: *difficulty,
		},
		Mining: model.MiningStatistics{
			MineRate:     difficultyCalculator.Rate(),
			RecentMiners: recentMiners,
		},
	})
}

func resultChallengeRoute(c *gin.Context) {
	resultChallengeQuery := &model.ResultChallengeQuery{}

	if c.ShouldBind(resultChallengeQuery) != nil {
		c.String(http.StatusBadRequest, "Invalid input data!")
		return
	}

	walletid, err := strconv.ParseInt(resultChallengeQuery.WalletID, 10, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "wallet id couldn't be converted to int from string")
		return
	}

	challengeCorrect := CheckChallenge(resultChallengeQuery.Magic, *difficulty)
	if challengeCorrect {
		// Get tokens for a specific user
		// incrementerr := coindb.IncrementBalance(walletid, uint64(*reward))
		incrementerr := (&coindb.Wallet{
			Id: walletid,
		}).IncrementBalance(uint64(*reward))

		transactionerr := (&coindb.Transaction{
			Senderid:   1,
			Receiverid: walletid,
			Amount:     uint64(*reward),
		}).Register()

		if incrementerr != nil || transactionerr != nil {
			color.Error.Println(incrementerr, transactionerr)
			c.String(http.StatusInternalServerError, "Internal server error, please report to @stack_smash")
			return
		}

		difficultyCalculator.IncreaseCounter()
		// *difficulty = difficultyCalculator.CalculateNextDifficulty(*difficulty)

		c.String(202, "Correct")
		color.Info.Tips(resultChallengeQuery.WalletID + " succesfully mined a block!")

		RefreshChallenge()
		difficultyCalculator.ReportRateStats()

		return
	} else {
		c.String(406, "Incorrect")
		return
	}
}

func walletBalanceRoute(c *gin.Context) {
	walletBalanceQuery := &model.WalletBalanceQuery{}

	if c.ShouldBindQuery(walletBalanceQuery) != nil {
		c.String(http.StatusBadRequest, "Invalid input data!")
		return
	}

	walletid, err := strconv.ParseInt(walletBalanceQuery.WalletID, 10, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid input data!")
		return
	}

	wallet := &coindb.Wallet{Id: walletid}

	err = wallet.GetWallet()
	if err != nil {
		color.Error.Printf("Database error in walletBalanceRoute(): %v", err)
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.JSON(200, wallet)
}
