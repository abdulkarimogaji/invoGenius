// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
	"time"
)

type Invoice struct {
	ID        int32
	UserID    int32
	Amount    float64
	Vat       float64
	Type      string
	IssuedAt  time.Time
	FromDate  time.Time
	UntilDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Receipt struct {
	ID            int32
	TransactionID int32
	UploadedBy    sql.NullInt32
	Filename      sql.NullString
	File          sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Setting struct {
	ID           int32
	SettingKey   string
	SettingValue sql.NullString
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}

type Transaction struct {
	ID            int32
	InvoiceID     int32
	PaymentMethod sql.NullString
	PaidAt        sql.NullTime
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type User struct {
	ID        int32
	FirstName sql.NullString
	LastName  sql.NullString
	Email     string
	Role      sql.NullString
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
