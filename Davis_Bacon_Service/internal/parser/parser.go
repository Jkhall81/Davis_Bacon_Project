package parser

import (
	"fmt"
	"log"
	"regexp"

	"davisbacon/internal/db"
)

func ParseDocuments() error {
	log.Println("Parsing wage determination documents...")

	rows, err := db.DB.Query(`SELECT wd_id, revision_number, publish_date, document FROM wd_documents`)
	if err != nil {
		return err
	}
	defer rows.Close()

	ratePattern := regexp.MustCompile(`^(?P<label>[A-Za-z0-9\s/\-\(\):]+?)\.*\s*\$?(?P<base>\d+\.\d{2}|\*\*)\s+(?P<fringe>\d+\.\d{2}|\*\*)`)

	for rows.Next() {
		var wdID string
		var rev int
		var pubDate *string
		var text string
		if err := rows.Scan(&wdID, &rev, &pubDate, &text); err != nil {
			continue
		}

		for _, line := range regexp.MustCompile(`\r?\n`).Split(text, -1) {
			m := ratePattern.FindStringSubmatch(line)
			if m == nil {
				continue
			}

			base := 0.0
			fringe := 0.0
			var note *string

			if m[2] != "**" {
				fmt.Sscanf(m[2], "%f", &base)
			} else {
				n := "See WD notes (**)"
				note = &n
			}

			if m[3] != "**" {
				fmt.Sscanf(m[3], "%f", &fringe)
			} else {
				n := "See WD notes (**)"
				note = &n
			}

			detail := db.WDDetail{
				WDID:           wdID,
				Classification: m[1],
				BaseRate:       base,
				FringeRate:     fringe,
				Notes:          note,
			}
			_ = db.InsertDetail(detail)
		}
	}
	log.Println("All documents parsed successfully.")
	return nil
}
