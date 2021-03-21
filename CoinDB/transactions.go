package CoinDB

type Transaction struct {
	Id uint64

	Senderid   int64
	Receiverid int64

	Amount uint64
}

func (t *Transaction) Register() error {
	_, err := pdb.Model(t).Insert()
	return err
}

func (t *Transaction) FindById() error {
	return pdb.Model(t).WherePK().Select(t)
}

func (t *Transaction) LastBySender(amount int) ([]Transaction, error) {
	var transactions []Transaction

	err := pdb.Model(&transactions).Where("transaction.id = ?", t.Senderid).OrderExpr("id DESC").Limit(amount).Select()

	return transactions, err
}
