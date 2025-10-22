package main

import (
	"log"
	"time"

	"davisbacon/internal/api"
	"davisbacon/internal/db"
	"davisbacon/internal/fetcher"
	"davisbacon/internal/parser"
	"davisbacon/internal/scraper"
	"davisbacon/internal/util"
)

func main() {
	start := time.Now()
	util.Info("Davis-Bacon Data Service starting up...")

	// 1. Load environment + connect to DB
	dsn := util.MustGetEnv("DATABASE_URL")
	db.InitDB(dsn)
	defer db.DB.Close()

	client := api.NewClient()

	// 2. Scrape metadata for all states
	states := []string{
		"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
		"HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD",
		"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
		"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
		"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY", "DC",
	}

	if err := scraper.ScrapeAllStates(client, states); err != nil {
		log.Fatalf("Scraper failed: %v", err)
	}

	// 3. Fetch wage determination documents
	if err := fetcher.FetchDocuments(client); err != nil {
		log.Fatalf("Fetcher failed: %v", err)
	}

	// 4. Parse text documents into details
	if err := parser.ParseDocuments(); err != nil {
		log.Fatalf("Parser failed: %v", err)
	}

	util.Success("All data processed successfully!")
	util.Info("⏱Total runtime:", time.Since(start))
}
