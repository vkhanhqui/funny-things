package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		panic(err)
	}

	updates := map[string]string{
		"github":   "new vkhanhqui",
		"website":  "new khanhqui.com",
		"linkedin": "new khanhqui",
		"name":     "new Qui Vo",
	}

	// Run in parallel
	var wg sync.WaitGroup
	for k, v := range updates {
		wg.Add(1)
		go func() {
			query := fmt.Sprintf(`
			UPDATE example_tbl 
			SET data = jsonb_set(
				COALESCE(data, '{}'), 
				'{%s}', '"%s"'::jsonb
			)
			WHERE id = $1`, k, v)
			_, err = conn.Exec(ctx, query, 1)
			if err != nil {
				panic(err)
			}

			defer wg.Done()
		}()
	}
	wg.Wait()
}
