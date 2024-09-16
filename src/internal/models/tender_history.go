package models

type TenderHistory struct {
	ID              int          `json:"id" gorm:"primaryKey"`
	TenderID        int          `json:"tenderID"`
	Name            string       `json:"name" gorm:"not null"`
	Description     string       `json:"description"`
	ServiceType     string       `json:"serviceType"`
	Status          TenderStatus `json:"status"`
	OrganizationID  int          `json:"organizationId" gorm:"not null"`
	CreatorUsername string       `json:"creatorUsername"`
	Version         int          `json:"version"`
}
