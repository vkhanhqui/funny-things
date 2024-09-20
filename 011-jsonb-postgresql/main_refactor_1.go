package main

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbURL = "postgresql://admin:123456@localhost:5432/example_db"

type ExampleTbl struct {
	ID   int            `json:"id"`
	Data map[string]any `json:"data"`
}

func main_refactor_1() {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		panic(err)
	}

	// Insert
	// data := map[string]any{
	// 	"name":    "Qui Vo",
	// 	"age":     24,
	// 	"country": "Vietnam",
	// }
	// query := "INSERT INTO example_tbl (data) VALUES ($1)"
	// _, err = conn.Exec(ctx, query, data)
	// if err != nil {
	// 	panic(err)
	// }

	// Retrieve
	query := "SELECT id, data FROM example_tbl WHERE id = $1"
	var eTbl ExampleTbl
	err = pgxscan.Get(ctx, conn, &eTbl, query, 1)
	if err != nil {
		panic(err)
	}

	// Update
	eTbl.Data["github"] = "vkhanhqui"
	query = "UPDATE example_tbl SET data = $1 WHERE id = $2"
	_, err = conn.Exec(ctx, query, eTbl.Data, eTbl.ID)
	if err != nil {
		panic(err)
	}
}
