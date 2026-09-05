package model

type MenuItem struct {
	ExternalID  string
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	IsActive    bool
	IsAvailable bool
}

type MenuSyncResult struct {
	CreatedCount     int
	UpdatedCount     int
	DeactivatedCount int
}
