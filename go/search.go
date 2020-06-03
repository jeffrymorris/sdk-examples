package main

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/couchbase/gocb/v2/search"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("172.23.111.3", opts)
	if err != nil {
		panic(err)
	}

	// For Server versions 6.5 or later you do not need to open a bucket here
	bucket := cluster.Bucket("travel-sample")

	// We wait until the bucket is definitely connected and setup.
	// For Server versions 6.5 or later if we hadn't opened a bucket then we could use cluster.WaitUntilReady here.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	matchResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		search.NewMatchQuery("swanky"),
		&gocb.SearchOptions{
			Limit:  10,
			Fields: []string{"description"},
		},
	)
	if err != nil {
		panic(err)
	}

	for matchResult.Next() {
		row := matchResult.Row()
		docID := row.ID
		score := row.Score

		var fields interface{}
		err := row.Fields(&fields)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Document ID: %s, search score: %f, fields included in result: %v\n", docID, score, fields)
	}

	// always check for errors after iterating
	err = matchResult.Err()
	if err != nil {
		panic(err)
	}

	dateRangeResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		search.NewDateRangeQuery().Start("2019-01-01", true).End("2019-02-01", false),
		&gocb.SearchOptions{
			Limit: 10,
		},
	)
	if err != nil {
		panic(err)
	}

	for dateRangeResult.Next() {
		row := dateRangeResult.Row()
		docID := row.ID
		score := row.Score

		fmt.Printf("Document ID: %s, search score: %f, fields included in result: %v\n", docID, score)
	}

	// always check for errors after iterating
	err = dateRangeResult.Err()
	if err != nil {
		panic(err)
	}

	conjunctionResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		search.NewConjunctionQuery(
			search.NewMatchQuery("swanky"),
			search.NewDateRangeQuery().Start("2019-01-01", true).End("2019-02-01", false),
		),
		&gocb.SearchOptions{
			Limit: 10,
		},
	)
	if err != nil {
		panic(err)
	}

	for conjunctionResult.Next() {
		row := conjunctionResult.Row()
		docID := row.ID
		score := row.Score

		var fields interface{}
		err := row.Fields(&fields)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Document ID: %s, search score: %f, fields included in result: %v\n", docID, score, fields)
	}

	// always check for errors after iterating
	err = conjunctionResult.Err()
	if err != nil {
		panic(err)
	}

	facetsResult, err := cluster.SearchQuery(
		"my-facet-index",
		search.NewMatchAllQuery(),
		&gocb.SearchOptions{
			Limit: 10,
			Facets: map[string]search.Facet{
				"type": search.NewTermFacet("type", 5),
			},
		},
	)
	if err != nil {
		panic(err)
	}

	for facetsResult.Next() {
	}

	facets, err := facetsResult.Facets()
	if err != nil {
		panic(err)
	}
	for _, facet := range facets {
		field := facet.Field
		total := facet.Total

		fmt.Printf("Facet field: %s, total: %d\n", field, total)
	}

	collection := bucket.DefaultCollection()

	hotel := struct {
		Description string `json:"description"`
	}{Description: "super swanky"}
	myWriteResult, err := collection.Upsert("a-new-hotel", hotel, nil)
	if err != nil {
		panic(err)
	}

	consistentWithResult, err := cluster.SearchQuery(
		"travel-sample-index-hotel-description",
		search.NewMatchQuery("swanky"),
		&gocb.SearchOptions{
			Limit:          10,
			ConsistentWith: gocb.NewMutationState(*myWriteResult.MutationToken()),
		},
	)
	if err != nil {
		panic(err)
	}

	for consistentWithResult.Next() {
		row := consistentWithResult.Row()
		docID := row.ID
		score := row.Score

		fmt.Printf("Document ID: %s, search score: %f\n", docID, score)
	}

	// always check for errors after iterating
	err = consistentWithResult.Err()
	if err != nil {
		panic(err)
	}
}
