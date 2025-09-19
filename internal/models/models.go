package models

import "database/sql"

type Models struct {
	Subscriptions SubscriptionDB
}

func NewModels(db *sql.DB) Models {
	return Models{
		Subscriptions: SubscriptionDB{DB: db},
	}
}
