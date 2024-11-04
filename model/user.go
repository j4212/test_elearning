package model

import (
	"time"

	"gorm.io/gorm"
)

type ROLE string

const (
	SUPER_ADMIN ROLE = "SUPER_ADMIN"
	ADMIN       ROLE = "ADMIN"
	TEACHER     ROLE = "TEACHER"
	STUDENT     ROLE = "STUDENT"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	Password  string
	Status    string `gorm:"type:varchar(20)"`
	Role      ROLE
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Student struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	IdNumber  int
	UserID    string
	SchoolsID string
	Schools   Schools `gorm:"foreignKey:SchoolsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Teacher struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	User      User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SchoolsID *string
	Schools   Schools `gorm:"foreignKey:SchoolsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name      string
	IdNumber  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ActiveStudent struct {
	ID         string `gorm:"primaryKey"`
	StudentID  string
	Student    Student `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SchoolYear string
	Class      string
	ClassSlug  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type AdminSchool struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	User      User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SchoolID  string
	Schools   Schools `gorm:"foreignKey:SchoolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
