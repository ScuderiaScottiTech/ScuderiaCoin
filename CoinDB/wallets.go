package CoinDB

import "github.com/go-pg/pg/v10"

type Wallet struct {
	Id      int64
	Balance uint64
}

func (w *Wallet) GetWallet() error {
	err := pdb.Model(w).WherePK().Select()

	if err == pg.ErrNoRows {
		newWallet := &Wallet{Id: w.Id}
		_, err = pdb.Model(newWallet).Insert()
	}

	return err
}

func (w *Wallet) Exists() (bool, error) {
	return pdb.Model(w).WherePK().Exists()
}

func (w *Wallet) IncrementBalance(amount uint64) error {
	err := w.GetWallet()
	if err != nil {
		return err
	}

	w.Balance += amount

	_, err = pdb.Model(w).WherePK().Update()
	return err
}

func (w *Wallet) DecrementBalance(amount uint64) error {
	err := w.GetWallet()
	if err != nil {
		return err
	}

	w.Balance -= amount

	_, err = pdb.Model(w).WherePK().Update()
	return err
}

func GetTopWallets() ([]Wallet, error) {
	var wallets []Wallet

	err := pdb.Model(&wallets).OrderExpr("balance DESC").Where("wallet.balance != 0").Limit(10).Select()

	return wallets, err
}
