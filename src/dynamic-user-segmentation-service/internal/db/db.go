package db

import (
	"database/sql"
	"dynamic-user-segmentation-service/internal/utils"
	"fmt"
	_ "github.com/lib/pq"
)

func CreateConnection(cnf utils.DatabaseConfiguration) *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cnf.Host, cnf.User, cnf.Password, cnf.Name, cnf.Port)
	db, err := sql.Open(cnf.Type, dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
