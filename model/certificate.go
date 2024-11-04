package model

import (
	// "time"
)

type Certificate struct {
	ID                string `gorm:"primaryKey"`
	RecipientName     string
	CourseName        string
	CertificateNo      string
	Score             string
	CompletionEndDate string
	CertificateUrl    string
}
