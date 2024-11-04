package model

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	ID          string `gorm:"primaryKey"`
	ChapterID   string
	Chapter     Chapter `gorm:"foreignKey:ChapterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MaterialID  string
	Material    Material `gorm:"foreignKey:MaterialID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title       string
	Description string
	Quizes      []Quizes `gorm:"foreignKey:QuizID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Quizes struct {
	ID          string `gorm:"primaryKey"`
	Quiz        string
	QuizID      string
	ImgURL      string
	QuizAnswers []QuizAnswer `gorm:"foreignKey:QuizesID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type QuizAnswer struct {
	ID                 string `gorm:"primaryKey"`
	Answer             string
	QuizesID           string
	QuizAnswerStudents []QuizAnswerStudent `gorm:"foreignKey:QuizAnswerID"`
	IsCorrect          bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type QuizAnswerStudent struct {
	ID              string `gorm:"primaryKey"`
	QuizID          string // quiz nya
	Quiz            Quiz   `gorm:"foreignKey:QuizID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuizAnswerID    string
	QuizAnswer      QuizAnswer `gorm:"foreignKey:QuizAnswerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuizesID        string     `gorm:"foreignKey:QuizesID"`
	ActiveStudentID string
	ActiveStudent   ActiveStudent `gorm:"foreignKey:ActiveStudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Grades          int
	Score           int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
