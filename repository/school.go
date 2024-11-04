package repository

import (
	"github.com/cvzamannow/E-Learning-API/entity"
	"github.com/cvzamannow/E-Learning-API/model"
)

type SchoolRepository interface {
	CreateSchool(request model.Schools) (*model.Schools, error)
	FindSchool(codd map[string]interface{}) (*model.Schools, error)
	GetAllSchool() (*[]model.Schools, error)
	UpdateSchool(id string, request model.Schools) (*model.Schools, error)
	DeleteSchool(id string) (*model.Schools, error)

	// Class
	FindClasses(q string) ([]entity.ClassEntity, error)
}
