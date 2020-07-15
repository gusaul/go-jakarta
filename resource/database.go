package resource

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DatabaseConn *sqlx.DB

func InitDatabase() {
	DatabaseConn = GetDatabaseConn()
}

func GetDatabaseConn() *sqlx.DB {
	conn, err := sqlx.Connect("postgres", "host=localhost user=gojakarta password='password' dbname=product sslmode=disable")
	if err != nil {
		panic(err)
	}
	return conn
}
