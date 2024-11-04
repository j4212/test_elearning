package model

import (
	"time"

	"gorm.io/gorm"
)

/*
* Model is simply database enitty for query data.
* Course is a where teacher can add their own study material.
 */

type STATUS string

const (
	APPROVE STATUS = "APPROVE"
	REJECT  STATUS = "REJECTED"
	PENDING STATUS = "PENDNG"
)

type Course struct {
	ID               string `gorm:"primaryKey"`
	TeacherID        string
	Title            string
	Description      string
	Detail           string
	EstimationHour   string
	EstimationMinute string
	Slug             string
	IsDraft          bool
	ThumbnailImg     string
	CompleteCourses  []CompleteCourse `gorm:"foreignKey:CourseID"`
	CourseClasses    []CourseClass    `gorm:"foreignKey:CourseID"`
	Chapters         []Chapter        `gorm:"foreignKey:CourseID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type Chapter struct {
	ID        string `gorm:"primaryKey"`
	CourseID  string
	Title     string
	Slug      string
	IsDraft   bool
	Materials []Material `gorm:"foreignKey:ChapterID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Material struct {
	ID         string `gorm:"primaryKey"            json:"id"`
	ChapterID  string `                             json:"chapter_id"`
	CourseID   string
	Title      string                `                             json:"title"`
	Type       string                `                             json:"type"`
	Slug       string                `                             json:"slug"`
	Theory     Theory                `gorm:"foreignKey:MaterialID" json:"theory"`
	Submission Submission            `gorm:"foreignKey:MaterialID" json:"submission"`
	Progress   []ActiveStudentCourse `gorm:"foreignKey:MaterialID"`
	CreatedAt  time.Time             `                             json:"created_at"`
	UpdatedAt  time.Time             `                             json:"updated_at"`
	DeletedAt  gorm.DeletedAt        `gorm:"index"                 json:"deleted_at"`
}

type Theory struct {
	ID         string `gorm:"primaryKey"`
	MaterialID string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Submission struct {
	ID         string `gorm:"primaryKey"`
	MaterialID string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type CourseClass struct {
	ID        string `gorm:"primaryKey"`
	CourseID  string
	Class     string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SubmissionStudent struct {
	ID              string `gorm:"primaryKey"`
	MaterialID      string
	CourseID        string
	Material        Material `gorm:"foreignKey:MaterialID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SchoolID        string
	Schools         Schools `gorm:"foreignKey:SchoolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TeacherID       string
	ActiveStudentID string
	ActiveStudent   ActiveStudent `gorm:"foreignKey:ActiveStudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status          STATUS
	FileUrl         string
	Grade           int
	Course          string
	Comment         *string
	Description     string
	Class           string
	SchoolYear      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type ActiveStudentCourse struct {
	ID              string `gorm:"primaryKey"`
	ActiveStudentID string
	ActiveStudent   ActiveStudent `gorm:"foreignKey:ActiveStudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CourseID        string
	Course          Course `gorm:"foreignKey:CourseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MaterialID      string
	Material        Material `gorm:"foreignKey:MaterialID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type CompleteCourse struct {
	ID              string `gorm:"primaryKey"`
	CourseID        string
	ActiveStudentID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
