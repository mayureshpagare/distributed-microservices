package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)


func initDB() *sql.DB{
	dsn:= "ecom_user:localhost@tcp(127.0.0.1:3306)/ecommerce"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	
	if err := db.Ping(); err != nil {
		log.Fatal("Error Connecting to database : %v",err)
	}
	log.Println("Database Connected")
	return db
}

