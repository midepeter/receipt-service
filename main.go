package main

import (
	"context"
	"log/slog"
)

func main() {

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
