package main

import (
	"flag"
	"math/rand"
	"time"

	coindb "github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/CoinDB"
	"github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/DifficultyAlgorithm"
	"github.com/gin-gonic/gin"
)

var (
	host             = flag.String("host", "0.0.0.0", "Listening address")
	port             = flag.String("port", "8989", "Listening port")
	databaseAddr     = flag.String("dbconn", "", "Address and port of the postgres database")
	databasePassword = flag.String("dbpass", "", "Database password for user postgres")

	difficulty    = flag.Int("miningdifficulty", 4, "Mining difficulty")
	mindifficulty = flag.Int("mindifficulty", 4, "Minimum difficulty")
	maxdifficulty = flag.Int("maxdifficulty", 10, "Maximum difficulty")

	reward        = flag.Int("minereward", 10, "Reward for a mined challenge")
	rewardtarget  = flag.Int("targethours", 30, "Reward target time")
	targetminutes = flag.Int("targetminutes", 15, "Reward minutes")

	reportratestats = flag.Bool("reportrate", false, "Report rate statistics")
	apiratelimiter  = flag.Int("apiratelimiter", 6, "Rate limiter allowance")
)

var difficultyCalculator = &DifficultyAlgorithm.DifficultyCalculator{}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	coindb.InitDatabaseConnection(*databaseAddr, *databasePassword)
	difficultyCalculator.Initialize(*rewardtarget, *targetminutes, *reward, *mindifficulty, *maxdifficulty, *reportratestats)
	RefreshChallenge()

	go difficultyCalculator.PeriodicCalculator(difficulty)

	api := gin.Default()

	// apiLimiter := ginlimiter.NewRateLimiter(time.Second, int64(*apiratelimiter), func(c *gin.Context) (string, error) {
	// 	return "", nil
	// })

	mine := api.Group("/mine")
	{
		mine.GET("/getChallenge", getChallengeRoute)
		mine.GET("/resultChallenge", resultChallengeRoute)
		mine.GET("/statistics", mineStatistics)
	}

	wallet := api.Group("/wallet")
	{
		wallet.GET("/balance", walletBalanceRoute)
	}

	api.Run(*host + ":" + *port)
}
