package helper

import (
	"time"

	"github.com/cvzamannow/E-Learning-API/entity"
	"github.com/cvzamannow/E-Learning-API/model"
)

type NextMaterialArg struct {
	ChapterID         string
	CurrentMaterialID string
	CreatedAt         time.Time
}

func CountStatusSubmission(data []model.SubmissionStudent) *entity.TotalSubmissionStatus {
	var pending int
	var reject int 
	var approve int 

	if len(data) == 0 {
		return &entity.TotalSubmissionStatus{
			Pending: 0,
			Rejected: 0,
			Approved: 0,
		}
	}

	for _, el := range data {
		if el.Status == "PENDING" {
			pending = pending + 1
		} else if el.Status == "REJECTED" || el.Status == "REV_REJECT" {
			reject = reject + 1 
		} else if el.Status == "APPROVE" {
			approve = approve + 1
		} else if el.Status == "APPROVED" {
			approve = approve + 1
		} 
	}

	return &entity.TotalSubmissionStatus{
		Pending: pending,
		Rejected: reject,
		Approved: approve,
	}
}