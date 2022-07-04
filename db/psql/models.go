package psql

import "time"

type Transaction struct {
	Id    string `db:"id"`
	Name  string `db:"name"`
	State int    `db:"state"`
}

type Tag struct {
	TransactionId string `db:"transaction_id"`
	Name          string `db:"name"`
	Value         string `db:"value"`
}

type Action struct {
	TransactionId string    `db:"transaction_id"`
	ActionId      string    `db:"action_id"`
	Name          string    `db:"name"`
	Time          time.Time `db:"time"`
}

type HttpAction struct {
	Id      string `db:"id"`
	Method  string `db:"method"`
	URL     string `db:"url"`
	Headers string `db:"headers"`
	Body    string `db:"body"`
}
