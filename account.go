package brackets

import (
	"github.com/gocraft/dbr/v2"
	"log"
)

type Account struct {
	Id          int64  `db:"id"`
	AdminUser   int64  `db:"admin_user_id"`
	AccountName string `db:"account_name"`
}

func LoadAccountByName(tx dbr.SessionRunner, name string) (*Account, error) {
	var account Account
	err := tx.Select("*").From("accounts").Where("account_name = ?", name).
		LoadOne(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil

}

func CreateAccount(tx dbr.SessionRunner, name string) (*Account, error) {

	var id int64
	err := tx.InsertInto("accounts").
		Pair("admin_user_id", 0).
		Pair("account_name", name).
		Returning("id").Load(&id)

	if err != nil {
		log.Fatalf("Create Account failed: %v", err)
		return nil, err
	}

	var account Account
	err = tx.Select("*").From("accounts").Where("id = ?", id).LoadOne(&account)

	if err != nil {
		log.Fatalf("Select User failed: %v", err)
		return nil, err
	}

	return &account, err
}
