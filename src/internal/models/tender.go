package models

type TenderStatus string

const (
	CreatedT   TenderStatus = "CREATED"
	PublishedT TenderStatus = "PUBLISHED"
	ClosedT    TenderStatus = "CLOSED"
)

type Tender struct {
	ID              int          `json:"id" gorm:"primaryKey"`
	Name            string       `json:"name" gorm:"not null" validate:"required"`
	Description     string       `json:"description" validate:"required"`
	ServiceType     string       `json:"serviceType" validate:"required"`
	Status          TenderStatus `json:"status"`
	OrganizationID  int          `json:"organizationId" gorm:"not null" validate:"required"`
	CreatorUsername string       `json:"creatorUsername" validate:"required"`
	Version         int          `json:"version"`
}

/*
{
    "name": "FIRST TENDER",
    "description": "NEW",
    "serviceType": "it",
    "organizationId": 1,
    "creatorUsername": "user1"
}
*/
