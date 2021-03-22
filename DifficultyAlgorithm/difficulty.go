package DifficultyAlgorithm

import (
	"fmt"
	"time"

	"github.com/paulbellamy/ratecounter"
)

type DifficultyCalculator struct {
	limiter          *ratecounter.RateCounter
	BlockTargetHours int64
	MinDifficulty    int
	MaxDifficulty    int
}

func (dc *DifficultyCalculator) Initialize(blockTargetHours int64, minDifficulty int, maxdifficulty int, reportratestats bool) {
	dc.BlockTargetHours = blockTargetHours
	dc.MinDifficulty = minDifficulty
	dc.MaxDifficulty = maxdifficulty

	dc.limiter = ratecounter.NewRateCounter(time.Hour)

	if reportratestats {
		go dc.ReportRateStatistics()
	}
}

func (dc *DifficultyCalculator) IncreaseCounter() {
	dc.limiter.Incr(1)
}

func (dc *DifficultyCalculator) CalculateNextDifficulty(difficulty int) int {
	if dc.limiter.Rate() > int64(dc.BlockTargetHours) && difficulty < dc.MaxDifficulty {
		return difficulty + 1
	} else if dc.limiter.Rate() < int64(dc.BlockTargetHours) && dc.MinDifficulty < difficulty {
		return difficulty - 1
	} else {
		return difficulty
	}
}

func (dc *DifficultyCalculator) ReportRateStatistics() {
	for {
		time.Sleep(time.Second)
		fmt.Println("Total rate:", dc.limiter.Rate(), "/ hour. Target:", dc.BlockTargetHours, "/ hour")
	}
}
