package main

import (
	"fmt"
	"os"

	"github.com/ostheperson/go-auth-service/internal/database"
	"github.com/ostheperson/go-auth-service/internal/server"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		database.SeedDB()
		return
	}
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
