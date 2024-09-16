package models

import "time"

type Decision string

const (
	SubmittedD Decision = "SUBMITTED"
	DeclinedD  Decision = "DECLINED"
)

type BidDecision struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	BidID          int       `json:"bidID" gorm:"primaryKey"`
	Username       string    `json:"username" validate:"required"`
	DecisionStatus Decision  `json:"decisionStatus"`
	DecisionTime   time.Time `json:"decisionDate"`
}
