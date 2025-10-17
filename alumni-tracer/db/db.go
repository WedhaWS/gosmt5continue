package db

import (
	"database/sql"
	"os"
	"time"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func NewDB() *sql.DB{
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db,err := sql.Open("pgx",os.Getenv("DATABASE"))
	if err != nil{
		panic(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
