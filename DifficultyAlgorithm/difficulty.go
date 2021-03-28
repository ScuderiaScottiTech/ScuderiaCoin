package DifficultyAlgorithm

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/paulbellamy/ratecounter"
)

type DifficultyCalculator struct {
	limiter *ratecounter.RateCounter

	Reward              int
	RewardTarget        int
	RewardTargetMinutes int

	MinDifficulty int
	MaxDifficulty int
}

func (dc *DifficultyCalculator) Initialize(rewardtarget int, rewardTargetMinutes int, reward int, minDifficulty int, maxdifficulty int, reportratestats bool) {
	dc.RewardTarget = rewardtarget
	dc.RewardTargetMinutes = rewardTargetMinutes
	dc.Reward = reward

	dc.MinDifficulty = minDifficulty
	dc.MaxDifficulty = maxdifficulty

	dc.limiter = ratecounter.NewRateCounter(time.Duration(rewardTargetMinutes) * time.Minute)

	if reportratestats {
		dc.PeriodicRateReporter()
	}
}

func (dc *DifficultyCalculator) Rate() int64 {
	return dc.limiter.Rate()
}

func (dc *DifficultyCalculator) IncreaseCounter() {
	dc.limiter.Incr(int64(dc.Reward))
}

func (dc *DifficultyCalculator) CalculateNextDifficulty(difficulty int) int {
	if dc.limiter.Rate() > int64(dc.RewardTarget) && difficulty < dc.MaxDifficulty {
		return difficulty + 1
	} else if dc.limiter.Rate() < int64(dc.RewardTarget) && dc.MinDifficulty < difficulty {
		return difficulty - 1
	} else {
		return difficulty
	}
}

func (dc *DifficultyCalculator) PeriodicCalculator(difficulty *int) {
	for {
		time.Sleep((time.Minute * time.Duration(dc.RewardTargetMinutes)) / 2)

		*difficulty = dc.CalculateNextDifficulty(*difficulty)
		color.Info.Printf("New difficulty level computed: %d\n", *difficulty)
	}
}

func (dc *DifficultyCalculator) PeriodicRateReporter() {
	go func() {
		for {
			time.Sleep(time.Second)
			dc.ReportRateStats()
		}
	}()
}

func (dc *DifficultyCalculator) ReportRateStats() {
	fmt.Printf("Rate: %v / %v per %v minutes.\n", dc.Rate(), dc.RewardTarget, dc.RewardTargetMinutes)
}
