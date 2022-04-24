package CoinDB

import (
	"time"

	"github.com/go-pg/pg/v10/orm"
)

type Transaction struct {
	Id uint64 `json:"transactionid"`

	Time time.Time `json:"time"`

	Senderid   int64 `json:"sender"`
	Receiverid int64 `json:"receiver"`

	Amount uint64 `json:"amount"`
}

func (t *Transaction) Register() error {
	t.Time = time.Now()
	_, err := pdb.Model(t).Insert()
	return err
}

func (t *Transaction) FindById() error {
	return pdb.Model(t).WherePK().Select(t)
}

func (t *Transaction) LastBySender(amount int) ([]Transaction, error) {
	var transactions []Transaction

	var query *orm.Query
	if t.Receiverid != 0 {
		query = pdb.Model(t).Where("transaction.senderid = ? AND transaction.receiverid = ?", t.Senderid, t.Receiverid)
	} else {
		query = pdb.Model(t).Where("transaction.senderid = ?", t.Senderid)
	}

	err := query.OrderExpr("id DESC").Limit(amount).Select(&transactions)

	return transactions, err
}
