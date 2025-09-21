package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type SubscriptionDB struct {
	DB *sql.DB
}
type Subscription struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}
type SubscriptionWithCost struct {
	Subscription
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
	Cost     int       `json:"сost"`
}

//********************************************************************//
//  							 READ								  //
//********************************************************************//

func (m *SubscriptionDB) GetAll() ([]*Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM infosub.subscriptions  ORDER BY user_id, service_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subscriptions := []*Subscription{}
	for rows.Next() {
		var sub Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (m *SubscriptionDB) Get(id int) (*Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM infosub.subscriptions WHERE id = $1  ORDER BY user_id, service_name`

	var sub Subscription
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (m *SubscriptionDB) GetByUserID(uid uuid.UUID) ([]*Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM infosub.subscriptions WHERE user_id = $1  ORDER BY user_id, service_name`

	rows, err := m.DB.QueryContext(ctx, query, uid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subscriptions := []*Subscription{}
	for rows.Next() {
		var sub Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (m *SubscriptionDB) GetByUserSubscription(serviceName string) ([]*Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//log.Println(serviceName)

	query := `SELECT id, service_name, price, user_id, start_date, end_date
              FROM infosub.subscriptions
              WHERE service_name = $1 
              ORDER BY user_id, service_name`

	rows, err := m.DB.QueryContext(ctx, query, serviceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := []*Subscription{}
	for rows.Next() {
		var sub Subscription
		var endDate sql.NullTime

		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			sub.EndDate = &endDate.Time
		} else {
			sub.EndDate = nil
		}

		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

//********************************************************************//
//  							 CREATE								  //
//********************************************************************//

func (m *SubscriptionDB) Insert(sub *Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO infosub.subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	return m.DB.QueryRowContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
}

//********************************************************************//
//  							 DELETE								  //
//********************************************************************//

func (m *SubscriptionDB) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM infosub.subscriptions WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *SubscriptionDB) DeleteByUserID(uid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM infosub.subscriptions WHERE user_id = $1`
	_, err := m.DB.ExecContext(ctx, query, uid)
	return err
}

func (m *SubscriptionDB) DeleteByServiceName(serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM infosub.subscriptions WHERE service_name = $1`
	_, err := m.DB.ExecContext(ctx, query, serviceName)
	return err
}

//********************************************************************//
//  							 UPDATE								  //
//********************************************************************//

func (m *SubscriptionDB) Update(upd Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE infosub.subscriptions SET service_name = $1, price = $2,user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`
	_, err := m.DB.ExecContext(
		ctx, query,
		upd.ServiceName,
		upd.Price,
		upd.UserID,
		upd.StartDate,
		upd.EndDate,
		upd.ID,
	)
	return err
}

//********************************************************************//
//  							 FILTER								  //
//********************************************************************//

func (m *SubscriptionDB) GetSummary(from, to time.Time, userID, serviceName string) ([]*SubscriptionWithCost, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query, args, err := subQuery(userID, serviceName)
	if err != nil {
		return nil, 0, err
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	subCosts := []*SubscriptionWithCost{}
	totalcost := 0

	for rows.Next() {
		var sub Subscription
		var cost int

		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, 0, err
		}

		dateFrom := maxTime(sub.StartDate, from)
		dateTo := minTimePtr(sub.EndDate, to)
		cost = monthsDiff(dateFrom, dateTo) * sub.Price
		totalcost += cost

		subCost := SubscriptionWithCost{
			Subscription: sub,
			DateFrom:     dateFrom,
			DateTo:       dateTo,
			Cost:         cost,
		}

		subCosts = append(subCosts, &subCost)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return subCosts, totalcost, nil
}

//********************************************************************//
//  							 HELPERS							  //
//********************************************************************//

func subQuery(userID, serviceName string) (string, []interface{}, error) {
	var query string
	var args []interface{}

	switch {
	case userID != "" && serviceName != "":
		query = `SELECT * FROM infosub.subscriptions WHERE user_id = $1 AND service_name = $2 ORDER BY user_id, service_name`
		var parsedUUID uuid.UUID
		if err := parsedUUID.Scan(userID); err != nil {
			return "", nil, fmt.Errorf("invalid user_id: %w", err)
		}
		args = append(args, parsedUUID, serviceName)

	case userID == "" && serviceName != "":
		query = `SELECT * FROM infosub.subscriptions WHERE service_name = $1 ORDER BY user_id, service_name`
		args = append(args, serviceName)

	case userID != "" && serviceName == "":
		query = `SELECT * FROM infosub.subscriptions WHERE user_id = $1 ORDER BY user_id, service_name`
		var parsedUUID uuid.UUID
		if err := parsedUUID.Scan(userID); err != nil {
			return "", nil, fmt.Errorf("invalid user_id: %w", err)
		}
		args = append(args, parsedUUID)

	default:
		query = `SELECT * FROM infosub.subscriptions ORDER BY user_id, service_name`
	}

	return query, args, nil
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTimePtr(ptr *time.Time, t time.Time) time.Time {
	if ptr == nil {
		return t
	}
	if ptr.Before(t) {
		return *ptr
	}
	return t
}

func monthsDiff(from, to time.Time) int {
	years := to.Year() - from.Year()
	months := int(to.Month()) - int(from.Month())
	return years*12 + months + 1
}
