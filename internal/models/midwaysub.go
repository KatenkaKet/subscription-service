package models

import (
	"fmt"
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type MidwaySub struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date,omitempty"`
}

func (mid MidwaySub) FromMidwaySub() (Subscription, error) {
	// Парсим StartDate в формате "MM-YYYY"
	start, err := time.Parse("01-2006", mid.StartDate)
	if err != nil {
		return Subscription{}, err
	}

	// Парсим EndDate, если оно задано
	var endPtr *time.Time
	if mid.EndDate != "" {
		end, err := time.Parse("01-2006", mid.EndDate)
		if err != nil {
			return Subscription{}, err
		}
		endPtr = &end
	}

	sub := Subscription{
		ID:          mid.ID,
		ServiceName: mid.ServiceName,
		Price:       mid.Price,
		UserID:      mid.UserID,
		StartDate:   start,
		EndDate:     endPtr,
	}

	return sub, nil
}

func (m MidwaySub) PrintFields() {
	fmt.Printf("ID: %d\n", m.ID)
	fmt.Printf("ServiceName: %s\n", m.ServiceName)
	fmt.Printf("Price: %d\n", m.Price)
	fmt.Printf("UserID: %s\n", m.UserID)
	fmt.Printf("StartDate: %s\n", m.StartDate)
	fmt.Printf("EndDate: %s\n", m.EndDate)
}
