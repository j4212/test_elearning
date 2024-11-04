package repository

import (
	"errors"

	"github.com/cvzamannow/E-Learning-API/entity"
	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type schoolImpl struct {
	DB *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) SchoolRepository {
	return &schoolImpl{
		DB: db,
	}
}

func (repos *schoolImpl) CreateSchool(request model.Schools) (*model.Schools, error) {
	var existingSchool model.Schools

	if err := repos.DB.Where("name = ?", request.Name).First(&existingSchool).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := repos.DB.Create(&request).Error; err != nil {
				return nil, err
			}
			return &request, nil
		} else {
			logrus.Warnln("[database] Gagal karena:", err)
			return nil, err
		}
	}

	logrus.Warnln("[database] School name is already registered")

	return nil, errors.New("School name is already registered")
}

// DeleteSchool implements SchoolRepository.
func (repos *schoolImpl) DeleteSchool(id string) (*model.Schools, error) {
	var school model.Schools

	if err := repos.DB.Where("id = ?", id).Delete(&school); err != nil {
		return nil, errors.New("Id for School does not exist")
	}

	return &school, nil
}

// FindSchool implements SchoolRepository.
func (repos *schoolImpl) FindSchool(codd map[string]interface{}) (*model.Schools, error) {
	var school model.Schools

	if err := repos.DB.Where(codd).First(&school).Error; err != nil {
		return nil, errors.New("Id for school not found")
	}

	return &school, nil
}

// GetAllSchool implements SchoolRepository.
func (repos *schoolImpl) GetAllSchool() (*[]model.Schools, error) {
	var school []model.Schools

	if err := repos.DB.Find(&school).Error; err != nil {
		return nil, err
	}

	return &school, nil
}

// UpdateSchool implements SchoolRepository.
func (repos *schoolImpl) UpdateSchool(id string, request model.Schools) (*model.Schools, error) {

	if err := repos.DB.Where("id = ?", id).Updates(&request).Error; err != nil {
		return nil, err
	}

	return &request, nil
}

func (repos *schoolImpl) FindClasses(q string) ([]entity.ClassEntity, error) {
	var classes []entity.ClassEntity

	tx := repos.DB.Model(&model.ActiveStudent{}).Distinct("class")

	if q != "" {
		q = "%" + q + "-%"
		tx.Where("class_slug ILIKE ?", q)
	}

	tx.Find(&classes)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return classes, nil
}
