package main

import "time"

// TemplateData holds the data to be rendered in templates
type ReceiptTemplateData struct {
	Amount          string
	FromAccount     string
	Recipient       string
	ToAccount       string
	TransactionID   string
	Reference       string
	Timestamp       string
	Fee             float64
	TransactionDate string
}

type SummaryTemplateData struct {
	AccountOwner  string
	AccountNumber string
	FromDate      string
	ToDate        string
	Transactions  []Transaction
	Timestamp     string
}
type Transaction struct {
	TransactionID         string    `json:"transaction_id" db:"transaction_id"` // corrected typo
	Reference             string    `json:"reference" db:"reference"`           //generate random ref
	FromRoute             string    `json:"from_route" db:"from_route"`
	ToRoute               string    `json:"to_route" db:"to_route"`
	FromAccountNumber     string    `json:"from_account" db:"from_acct"`
	ToAccountNumber       string    `json:"to_account" db:"to_acct"`
	Amount                float64   `json:"amount" db:"amount"`
	Amountstring          string    `json:"amount_string" db:"amount_string"`
	Type                  string    `json:"type" db:"type"`
	Balance               float64   `json:"balance" db:"balance"`
	BalanceString         string    `json:"balance_string" db:"balance_string"`
	TransactionDateString string    `json:"transaction_date" db:"transaction_date"`
	TransactionDate       time.Time `json:"timestamp" db:"timestamp"`
}

type User struct {
	AccountID string `json:"account_id" db:"accountid"`
	Username  string `json:"username" db:"username"`
	PassHash  string `json:"-" db:"passhash"`
	FirstName string `json:"first_name" db:"firstname"`
	LastName  string `json:"last_name" db:"lastname"`
	BirthDay  string `json:"birth_day" db:"birthday"`
	Timezone  string `json:"timezone" db:"timezone"`
	Address   string `json:"address" db:"address"`
	State     string `json:"state" db:"state"`
	Zip       string `json:"zip" db:"zip"`
	SSN       string `json:"ssn" db:"ssn"`
}
