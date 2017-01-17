package main

import (
	"database/sql"
	"math/rand"
	"testing"

	_ "github.com/lib/pq"
)

func BenchmarkQuery(b *testing.B) {
	db, err := sql.Open("postgres", "host=localhost port=26257 dbname=bench user=root sslmode=disable")
	if err != nil {
		panic(err)
	}
	numberOfUsers := 100

	for ii := 0; ii < b.N; ii++ {
		uid := rand.Intn(numberOfUsers)
		rows, err := db.Query("SELECT permissions.rid from (permissions join users_permissions on users_permissions.permission_id = permissions.id) where users_permissions.user_id = $1 limit 1", uid)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				panic(err)
			}
		}
	}
}
