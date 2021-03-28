package DataModels

type ResultChallengeQuery struct {
	WalletID string `form:"walletid" binding:"required"`
	Magic    string `form:"magic" binding:"required"`
}

type WalletBalanceQuery struct {
	WalletID string `form:"walletid" binding:"required"`
}
