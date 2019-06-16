package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/bigquery"
)

func main() {
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	rows, err := query(proj)
	if err != nil {
		log.Fatal(err)
	}
	if err := printResults(os.Stdout, rows); err != nil {
		log.Fatal(err)
	}
}

// query returns a slice of the results of a query.
func query(proj string) (*bigquery.RowIterator, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		return nil, err
	}

	query := client.Query(
		`SELECT * FROM base.members`)

	return query.Read(ctx)
}

type StackOverflowRow struct {
	FistName string `bigquery:"first_name"`
	LastName string `bigquery:"last_name"`
	nickname string `bigquery:"nickname"`
}

// printResults prints results from a query to the Stack Overflow public dataset.
func printResults(w io.Writer, iter *bigquery.RowIterator) error {
	for {
		var row StackOverflowRow
		err := iter.Next(&row)
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "url: %s views: %d\n", row.URL, row.ViewCount)
	}
}
