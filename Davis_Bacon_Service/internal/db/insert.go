package db

import (
	"log"
)

func InsertWageDetermination(wd WageDetermination) error {
	_, err := DB.Exec(`
		INSERT INTO wage_determinations (id, wd_number, state, revision_number, published_date, modified_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			revision_number = EXCLUDED.revision_number,
			published_date = EXCLUDED.published_date,
			modified_date = EXCLUDED.modified_date
	`, wd.ID, wd.WDNumber, wd.State, wd.RevisionNumber, wd.PublishedDate, wd.ModifiedDate)

	if err != nil {
		log.Printf("InsertWageDetermination failed: %v", err)
	}
	return err
}

func InsertCounty(wdID string, county string) error {
	_, err := DB.Exec(`INSERT INTO wd_counties (wd_id, county) VALUES ($1, $2) ON CONFLIC DO NOTHING`, wdID, county)
	return err
}

func InsertConstructionType(wdID string, ctype string) error {
	_, err := DB.Exec(`
		INSERT INTO wd_construction_types (wd_id, construction_type)
		VALUES ($1, $2)
		ON CONFLIC (wd_id, construction_type) DO NOTHING
	`, wdID, ctype)
	if err != nil {
		log.Printf("InsertConstructionType failed for %s: %v", ctype, err)
	}
	return err
}

func InsertDocument(doc WDDoc) error {
	_, err := DB.Exec(`
		INSERT INTO wd_documents (wd_id, revision_number, publish_date, document)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (wd_id, revision_number)
		DO UPDATE SET document = EXCLUDED.document
	`, doc.WDID, doc.RevisionNumber, doc.PublishDate, doc.Document)
	return err
}

func InsertDetail(d WDDetail) error {
	_, err := DB.Exec(`
		INSERT INTO wd_details (wd_id, classification, group_number, base_rate, fringe_rate, effective_date, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, d.WDID, d.Classification, d.GroupNumber, d.BaseRate, d.FringeRate, d.EffectiveDate, d.Notes)
	return err
}
