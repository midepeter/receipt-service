package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dchest/uniuri"
	"github.com/jackc/pgx"
)

var (
	TransactionReceiptTemplateFile string = "receiptnewtwo.html"
	TransactionSummaryTemplateFile string = "transaction-history-new.html"
	TransactionReceipt             string = "transaction_receipt.pdf"
	TransactionSummary             string = "transaction_summary_receipt.pdf"
)

func (s *Server) getTransactionReceipt(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("id")
	if transactionID == "" {
		http.Error(w, "Transaction ID is required", http.StatusBadRequest)
		return
	}
	var transaction Transaction
	err := s.TransactionDB.DBPool.QueryRow(context.Background(), `
		SELECT amount, from_acct, to_acct, transaction_id, timestamp
		FROM transactions
		WHERE transaction_id = $1`, transactionID).Scan(&transaction.Amount, &transaction.FromAccountNumber, &transaction.ToAccountNumber, &transaction.TransactionID, &transaction.TransactionDate)
	if err != nil {
		http.Error(w, "Failed to fetch transaction details", http.StatusInternalServerError)
		return
	}

	var user User
	err = s.UsersDB.DBPool.QueryRow(context.Background(), `
		SELECT firstname, lastname
		FROM users
		WHERE accountid = $1`, transaction.FromAccountNumber).Scan(&user.FirstName, &user.LastName)
	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	data := ReceiptTemplateData{
		Amount:          fmt.Sprintf("%.2f", transaction.Amount), // Format to 2 decimal places, full amount
		FromAccount:     lastFourDigits(transaction.FromAccountNumber),
		Recipient:       user.FirstName + " " + user.LastName,
		ToAccount:       lastFourDigits(transaction.ToAccountNumber),
		TransactionID:   transaction.TransactionID,
		Reference:       uniuri.NewLen(10),
		Fee:             0.0,
		Timestamp:       time.Now().Format(time.RFC3339),
		TransactionDate: transaction.TransactionDate.String(),
	}

	// Compile the template
	html, err := CompileTemplate(TransactionReceiptTemplateFile, data)
	if err != nil {
		log.Fatalf("Failed to compile template: %v", err)
	}

	if err := GeneratePDF(html, TransactionReceipt); err != nil {
		log.Fatalf("Failed to generate PDF: %v", err)
	}

	http.ServeFile(w, r, TransactionReceipt)
}

func (s *Server) getTransactionSummary(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get "accountid", "from" and "to"
	type SummaryRequest struct {
		AccountID string `json:"accountid"`
		From      string `json:"from"`
		To        string `json:"to"`
	}
	var req SummaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.AccountID == "" || req.From == "" || req.To == "" {
		http.Error(w, "'accountid', 'from', and 'to' fields are required", http.StatusBadRequest)
		return
	}

	toDate, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		http.Error(w, "Invalid 'to' date format", http.StatusBadRequest)
		return
	}

	fromDate, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		http.Error(w, "Invalid 'from' date format", http.StatusBadRequest)
		return
	}

	if toDate.Before(fromDate) {
		http.Error(w, "'to' date must be after 'from' date", http.StatusBadRequest)
		return
	}

	summary, err := s.TransactionDB.DBPool.Query(context.Background(), `
		SELECT timestamp, amount, from_acct, to_acct
		FROM transactions
		WHERE (from_acct = $1 OR to_acct = $1) AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC`, req.AccountID, fromDate, toDate)
	if err != nil {
		http.Error(w, "Failed to fetch transaction summary", http.StatusInternalServerError)
		return
	}

	var transactions []Transaction
	for summary.Next() {
		var transaction Transaction
		if err := summary.Scan(&transaction.TransactionDate, &transaction.Amount, &transaction.FromAccountNumber, &transaction.ToAccountNumber); err != nil {
			http.Error(w, "Failed to scan transaction", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}
	if err := summary.Err(); err != nil {
		http.Error(w, "Error iterating over transactions", http.StatusInternalServerError)
		return
	}

	var user User
	s.UsersDB.DBPool.QueryRow(context.Background(), `
		SELECT firstname, lastname
		FROM users
		WHERE accountid = $1`, req.AccountID).Scan(&user.FirstName, &user.LastName)

	data := SummaryTemplateData{
		AccountOwner:  user.FirstName + " " + user.LastName,
		AccountNumber: lastFourDigits(req.AccountID),
		FromDate:      fromDate.Format("2006-01-02"),
		ToDate:        toDate.Format("2006-01-02"),
		Transactions:  transactions,
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	for k, v := range transactions {
		transactions[k].TransactionDateString = v.TransactionDate.Format("2006-01-02 15:04:05")
		transactions[k].Reference = uniuri.NewLen(10)
		transactions[k].Amountstring = fmt.Sprintf("%.2f", v.Amount)
		if v.FromAccountNumber == req.AccountID {
			transactions[k].Type = "Debit"
			transactions[k].Balance = -v.Amount
			if k != 0 {
				transactions[k].Balance = transactions[k-1].Balance - v.Amount
			}
		} else {
			transactions[k].Type = "Credit"
			transactions[k].Balance = v.Amount
			if k != 0 {
				transactions[k].Balance = transactions[k-1].Balance + v.Amount
			}
		}

		transactions[k].BalanceString = fmt.Sprintf("%.2f", transactions[k].Balance)
		transactions[k].FromRoute = lastFourDigits(v.FromAccountNumber)
	}

	// Compile the template
	html, err := CompileTemplate(TransactionSummaryTemplateFile, data)
	if err != nil {
		log.Fatalf("Failed to compile template: %v", err)
	}

	if err := GeneratePDF(html, TransactionSummary); err != nil {
		log.Fatalf("Failed to generate PDF: %v", err)
	}

	http.ServeFile(w, r, TransactionSummary)
}

func lastFourDigits(accountNumber string) string {
	if len(accountNumber) < 4 {
		return accountNumber
	}
	return accountNumber[len(accountNumber)-4:]
}
