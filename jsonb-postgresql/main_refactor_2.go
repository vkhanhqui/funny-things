package main

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

func retrieve(ctx context.Context, conn *pgxpool.Pool, id int) (ExampleTbl, error) {
	query := "SELECT id, data FROM example_tbl WHERE id = $1"
	var eTbl ExampleTbl
	err := pgxscan.Get(ctx, conn, &eTbl, query, id)
	return eTbl, err
}

func update(ctx context.Context, conn *pgxpool.Pool, id int, newData map[string]any) error {
	query := `
	UPDATE example_tbl 
	SET data = $2 
	WHERE id = $1`
	_, err := conn.Exec(ctx, query, id, newData)
	return err
}

func main_refactor_2() {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		panic(err)
	}

	eTbl, err := retrieve(ctx, conn, 1)
	if err != nil {
		panic(err)
	}

	newURLs := map[string]string{
		"github":   "vkhanhqui",
		"website":  "khanhqui.com",
		"linkedin": "khanhqui",
	}
	for k, v := range newURLs {
		eTbl.Data[k] = v
	}
	err = update(ctx, conn, 1, eTbl.Data)
	if err != nil {
		panic(err)
	}
}
