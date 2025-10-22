package fetcher

import (
	"fmt"
	"log"
	"time"

	"davisbacon/internal/api"
	"davisbacon/internal/db"
)

type WDDocResponse struct {
	Document string `json:"document"`
}

// FetchDocuments retrieves detailed WD text for each record in the DB.
func FetchDocuments(client *api.APIClient) error {
	log.Println("Fetching WD documents...")

	rows, err := db.DB.Query(`SELECT id, revision_number, published_date FROM wage_determinations ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var wdID string
		var rev int
		var pubDate *time.Time
		if err := rows.Scan(&wdID, &rev, &pubDate); err != nil {
			log.Println("Scan error:", err)
			continue
		}

		var exists int
		err = db.DB.QueryRow(`SELECT COUNT(1) FROM wd_documents WHERE wd_id = $1 AND revision_number = $2`, wdID, rev).Scan(&exists)
		if err != nil || exists > 0 {
			continue
		}

		url := fmt.Sprintf("%s/wdol/v1/wd/%s/%d", api.BaseURL, wdID, rev)
		var resp WDDocResponse
		if err := client.GetJSON(url, &resp); err != nil {
			log.Printf("Failed fetching doc for %s rev %d: %v", wdID, rev, err)
			continue
		}

		doc := db.WDDoc{
			WDID:           wdID,
			RevisionNumber: rev,
			PublishDate:    pubDate,
			Document:       resp.Document,
		}
		_ = db.InsertDocument(doc)
		time.Sleep(1 * time.Second)
	}
	log.Println("All WD documents fetched.")
	return nil
}
