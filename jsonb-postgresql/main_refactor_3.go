package main

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main_refactor_3() {
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
			// Retrieve
			eTbl, err := retrieve(ctx, conn, 1)
			if err != nil {
				panic(err)
			}

			// Update
			eTbl.Data[k] = v
			err = update(ctx, conn, 1, eTbl.Data)
			if err != nil {
				panic(err)
			}

			defer wg.Done()
		}()
	}
	wg.Wait()
}
