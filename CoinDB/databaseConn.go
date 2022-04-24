package CoinDB

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// var rdb *redis.Client
var pdb *pg.DB

func InitDatabaseConnection(databaseAddr string, databasePassword string) {
	pdb = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: databasePassword,
		Addr:     databaseAddr,
	})

	createSchema()
}

func createSchema() error {
	models := []interface{}{
		(*Wallet)(nil),
		(*Transaction)(nil),
	}

	for _, model := range models {
		err := pdb.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
