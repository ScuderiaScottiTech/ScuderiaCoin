package DataModels

import "github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/CoinDB"

type DifficultyStatistics struct {
	Minimum int `json:"minimum"`
	Maximum int `json:"maximum"`
	Current int `json:"current"`

	Reward        int `json:"reward"`
	RewardTarget  int `json:"reward_target_rate"`
	TargetMinutes int `json:"reward_rate_minutes"`
}

type MiningStatistics struct {
	MineRate     int64                `json:"current_rate"`
	RecentMiners []CoinDB.Transaction `json:"recent_miners"`
}

type MineStatisticsResponse struct {
	Difficulty DifficultyStatistics `json:"difficulty"`
	Mining     MiningStatistics     `json:"mining"`
}
