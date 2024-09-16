package models

type BidStatus string

const (
	CreatedB   BidStatus = "CREATED"
	PublishedB BidStatus = "PUBLISHED"
	CanceledB  BidStatus = "CANCELED"
	SubmittedB BidStatus = "SUBMITTED"
	DeclinedB  BidStatus = "DECLINED"
)

type Bid struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null" validate:"required"`
	Description     string    `json:"description" validate:"required"`
	Status          BidStatus `json:"status"`
	TenderID        int       `json:"tenderId" gorm:"not null" validate:"required"`
	OrganizationID  *int      `json:"organizationId" gorm:"default:null"`
	CreatorUsername string    `json:"creatorUsername" gorm:"not null" validate:"required"`
	Version         int       `json:"version"`
}

/*
{
    "name": "NEW BID",
    "description": "new",
    "tenderId": 1,
	"organizationId": 1,
    "creatorUsername": "user1"
}
*/
