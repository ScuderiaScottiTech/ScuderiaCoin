package main

import (
	"errors"
	"flag"
	"math/rand"
	"time"

	coindb "github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/CoinDB"
	"github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/DifficultyAlgorithm"
	"github.com/gin-gonic/gin"
	ginlimiter "github.com/julianshen/gin-limiter"
)

var (
	host             = flag.String("host", "0.0.0.0", "Listening address")
	port             = flag.String("port", "8989", "Listening port")
	databaseAddr     = flag.String("dbconn", "", "Address and port of the postgres database")
	databasePassword = flag.String("dbpass", "", "Database password for user postgres")
	difficulty       = flag.Int("miningdifficulty", 4, "Mining difficulty")
	mindifficulty    = flag.Int("mindifficulty", 4, "Minimum difficulty")
	reward           = flag.Int("minereward", 10, "Reward for a mined challenge")
	blocktargethours = flag.Int("targethours", 10, "Block target time")
	reportratestats  = flag.Bool("reportrate", false, "Report rate statistics")
)

var difficultyCalculator = &DifficultyAlgorithm.DifficultyCalculator{}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	coindb.InitDatabaseConnection(*databaseAddr, *databasePassword)
	difficultyCalculator.Initialize(int64(*blocktargethours), *mindifficulty, *reportratestats)
	RefreshChallenge()

	api := gin.Default()

	resultLimiter := ginlimiter.NewRateLimiter(time.Second, 1, func(c *gin.Context) (string, error) {
		walletid := c.Query("walletid")
		magic := c.Query("magic")

		if len(magic) >= 150 || len(walletid) >= 38 {
			c.String(400, "A query parameter is too long")
			return "", errors.New("query parameter length exceeded")
		}

		if walletid != "" && magic != "" {
			return "", nil
		}

		c.String(400, "A query parameter is missing")
		return "", errors.New("a query parameter is missing")
	})

	api.GET("/mine/getChallenge", getChallengeRoute)
	api.GET("/mine/resultChallenge", resultLimiter.Middleware(), resultChallengeRoute)
	
	api.GET("/transactions/list")

	api.Run(*host + ":" + *port)
}
