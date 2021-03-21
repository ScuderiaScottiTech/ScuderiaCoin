package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"

	coindb "github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/CoinDB"
)

func getChallengeRoute(c *gin.Context) {
	c.JSON(200, &ChallengeApiResponse{
		Reward:           *reward,
		CurrentChallenge: GetChallenge(),
		Difficulty:       *difficulty,
	})
}

func resultChallengeRoute(c *gin.Context) {
	swalletid := c.Query("walletid")
	magic := c.Query("magic")

	walletid, err := strconv.ParseInt(swalletid, 10, 0)
	if err != nil {
		c.String(400, "wallet id couldn't be converted to int from string")
		return
	}

	challengeCorrect := CheckChallenge(magic, *difficulty)
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
			c.String(500, "Internal server error, please report to @stack_smash")
			return
		}

		difficultyCalculator.IncreaseCounter()
		*difficulty = difficultyCalculator.CalculateNextDifficulty(*difficulty)

		c.String(202, "Correct")
		color.Info.Tips(swalletid + " succesfully mined a block!")

		RefreshChallenge()
		return
	} else {
		c.String(406, "Incorrect")
		return
	}
}
