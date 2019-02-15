package main

import (
	"github.com/hattaya92/finalexam/database"
	"github.com/hattaya92/finalexam/handler"
	_ "github.com/lib/pq"
)

func main() {
	database.ConnDB()
	handler.CreateTable()

	r := handler.NewRouter()
	r.Run(":2019")

}
