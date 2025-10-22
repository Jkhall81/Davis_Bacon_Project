package db

import "time"

type WageDetermination struct {
	ID             string
	WDNumber       string
	State          string
	RevisionNumber int
	PublishedDate  *time.Time
	ModifiedDate   *time.Time
}

type WDDoc struct {
	WDID           string
	RevisionNumber int
	PublishDate    *time.Time
	Document       string
}

type WDDetail struct {
	WDID           string
	Classification string
	GroupNumber    *int
	BaseRate       float64
	FringeRate     float64
	EffectiveDate  *time.Time
	Notes          *string
}
