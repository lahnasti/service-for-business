package models

type Review struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	BidID          int    `json:"bidId" validate:"bidId"`
	Username       string `json:"username" validate:"required"`
	OrganizationID int    `json:"organizationId" validate:"required"`
	Comment        string `json:"comment" validate:"required"`
}

/*
{
	"bidId": "1",
	"username": "user2",
	"organizationId": "1",
    "comment": "Good job!"
}
*/
