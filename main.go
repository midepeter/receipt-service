package main

import (
	"context"
	"log/slog"
)

func main() {
	// // Example template directory
	// templateDir := "./"

	// // Example template file
	// templateFile := filepath.Join(templateDir, "receiptnew.html")

	// // Example data for template
	// data := ReceiptTemplateData{
	// 	Amount:          12500.75,
	// 	FromAccount:     "1234567890",
	// 	Recipient:       "Jane Doe",
	// 	ToAccount:       "0987654321",
	// 	TransactionID:   "TXN123456789",
	// 	Reference:       "REF987654321",
	// 	Fee:             0.0,
	// 	Timestamp:       time.Now().Format(time.RFC3339),
	// 	TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
	// }

	// // Compile the template
	// html, err := CompileTemplate(templateFile, data)
	// if err != nil {
	// 	log.Fatalf("Failed to compile template: %v", err)
	// }

	// // Generate PDF
	// outputFile := "receipt.pdf"
	// if err := GeneratePDF(html, outputFile); err != nil {
	// 	log.Fatalf("Failed to generate PDF: %v", err)
	// }

	// fmt.Printf("PDF generated successfully: %s\n", outputFile)

	ctx := context.Background()
	transactionDB, err := NewDatabase(ctx, "postgresql://admin:password@35.223.9.138:5432/ledger-db", slog.Default())
	if err != nil {
		panic(err)
	}

	userDB, err := NewDatabase(ctx, "postgresql://accounts-admin:accounts-pwd@35.223.9.138:5432/accounts-db", slog.Default())
	if err != nil {
		panic(err)
	}

	server := NewServer(userDB, transactionDB, slog.Default())
	server.Start()
	//defer server.Shutdown()
}
