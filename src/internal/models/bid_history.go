package models

type BidHistory struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	BidID           int       `json:"bidID" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null" validate:"required"`
	Description     string    `json:"description" validate:"required"`
	Status          BidStatus `json:"status"`
	TenderID        int       `json:"tenderId" gorm:"not null" validate:"required"`
	OrganizationID  *int      `json:"organizationId" gorm:"default:null"`
	CreatorUsername string    `json:"creatorUsername" gorm:"not null" validate:"required"`
	Version         int       `json:"version"`
}
