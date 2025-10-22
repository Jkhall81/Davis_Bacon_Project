package scraper

import (
	"fmt"
	"log"
	"time"

	"davisbacon/internal/api"
	"davisbacon/internal/db"
	"davisbacon/internal/util"
)

type WageDeterminationResponse struct {
	Embedded struct {
		Results []struct {
			ID                  string `json:"_id"`
			FullReferenceNumber string `json:"fullReferenceNumber"`
			RevisionNumber      int    `json:"revisionNumber"`
			ModifiedDate        int64  `json:"modifiedDate"`
			PublishDate         int64  `json:"publishDate"`
			Location            struct {
				State struct {
					Name     string `json:"name"`
					Counties []struct {
						Value string `json:"value"`
					} `json:"counties"`
				} `json:"state"`
			} `json:"location"`
			ConstructionTypes []string `json:"constructionTypes"`
		} `json:"results"`
	} `json:"_embedded"`
	Page struct {
		TotalPages int `json:"totalPages"`
	} `json:"page"`
}

// ScrapeAllStates scrapes all WD metadata for a list of state codes.
func ScrapeAllStates(client *api.APIClient, states []string) error {
	util.Info("📡 Scraping wage determinations from SAM.gov...")

	for _, state := range states {
		page := 0
		for {
			url := fmt.Sprintf("%s/sgs/v1/search?index=dbra&state=%s&page=%d&size=25&is_active=true", api.BaseURL, state, page)
			var resp WageDeterminationResponse
			if err := client.GetJSON(url, &resp); err != nil {
				log.Printf("Failed scraping %s page %d: %v", state, page, err)
				break
			}

			results := resp.Embedded.Results
			if len(results) == 0 {
				break
			}

			for _, item := range results {
				published := time.UnixMilli(item.PublishDate)
				modified := time.UnixMilli(item.ModifiedDate)

				wd := db.WageDetermination{
					ID:             item.ID,
					WDNumber:       item.FullReferenceNumber,
					State:          item.Location.State.Name,
					RevisionNumber: item.RevisionNumber,
					PublishedDate:  &published,
					ModifiedDate:   &modified,
				}
				_ = db.InsertWageDetermination(wd)

				for _, c := range item.Location.State.Counties {
					_ = db.InsertCounty(item.ID, c.Value)
				}
				for _, t := range item.ConstructionTypes {
					_ = db.InsertConstructionType(item.ID, t)
				}
			}

			page++
			if page >= resp.Page.TotalPages {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}
