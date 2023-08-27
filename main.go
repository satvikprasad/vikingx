package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/satvikprasad/vikingx/api"
	"github.com/satvikprasad/vikingx/db"
	"github.com/satvikprasad/vikingx/okx"
)

func main() {
	godotenv.Load(".env")

	db, err := db.NewDB()
	if err != nil {
		fmt.Printf("Error creating database: %s", err)
	}

	a := okx.NewOkApi(true, ".env")

	api.ListenAndServe(db, a, os.Getenv("PORT"))
}
