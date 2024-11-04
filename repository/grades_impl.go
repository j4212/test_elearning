package repository

import (
	"github.com/cvzamannow/E-Learning-API/model"
	"gorm.io/gorm"
)

type gradesImpl struct {
	DB *gorm.DB
}

func NewGradesRepository(db *gorm.DB) GradesRepository {
	return &gradesImpl{
		DB: db,
	}
}

// GetGradesStudent implements GradesRepository.
func (repos *gradesImpl) GetGradesStudent(codd map[string]interface{}, class, schoolYear, schoolID string) (*[]model.SubmissionStudent, error) {

	var result []model.SubmissionStudent
	var schoolModel model.Schools

	query := repos.DB.Preload("ActiveStudent").Preload("Material").Preload("ActiveStudent.Student").Where(codd)
	repos.DB.Model(&model.Schools{}).Where("id = ?", schoolID).First(&schoolModel)

	if schoolYear != "" && schoolYear != " " {

		if class != "" && class != " " {
			query.Where("class = ?", class).Where("school_year = ?", schoolYear)
		} else {
			query.Where("school_year = ?", schoolYear)
		}
	} else if class != "" && class != " " {
		if schoolYear != "" && schoolYear != " " {
			query.Where("class = ?", class).Where("school_year = ?", schoolYear)
		} else {
			query.Where("class = ?", class)
		}
	} else {
		query.Where("school_year = ?", schoolModel.SchoolYear)
	}

	query.Where("status != ?", "REJECTED").Find(&result)

	return &result, nil

}

func (repos *gradesImpl) GetQuizGradesStudent(codd map[string]interface{}) (*[]model.QuizAnswerStudent, error) {
	var quizAnswerStudent []model.QuizAnswerStudent

	// Subquery untuk mencari entri terbaru
	subquery := repos.DB.Model(&model.QuizAnswerStudent{}).
		Select("MAX(created_at)").
		Where(codd)

	// Query utama untuk mendapatkan data berdasarkan subquery
	if err := repos.DB.Where("created_at = (?)", subquery).
		Preload("QuizAnswer").
		Preload("Quiz").
		Preload("Quiz.Material").
		Preload("ActiveStudent").
		Find(&quizAnswerStudent).
		Error; err != nil {
		return nil, err
	}

	return &quizAnswerStudent, nil
}

// GetGradesTeacher implements GradesRepository.
func (repos *gradesImpl) GetGradesTeacher(schoolYear, class, schoolID string) (*[]model.SubmissionStudent, error) {
	var result []model.SubmissionStudent
	var schoolModel model.Schools

	query := repos.DB.Preload("ActiveStudent").Preload("Material").Preload("ActiveStudent.Student").Where("school_id = ?", schoolID)
	repos.DB.Model(&model.Schools{}).Where("id = ?", schoolID).First(&schoolModel)

	if schoolYear != "" && schoolYear != " " {

		if class != "" && class != " " {
			query.Where("class = ?", class).Where("school_year = ?", schoolYear)
		} else {
			query.Where("school_year = ?", schoolYear)
		}
	} else if class != "" && class != " " {
		if schoolYear != "" && schoolYear != " " {
			query.Where("class = ?", class).Where("school_year = ?", schoolYear)
		} else {
			query.Where("class = ?", class)
		}
	} else {
		query.Where("school_year = ?", schoolModel.SchoolYear)
	}

	query.Find(&result)

	return &result, nil
}

func (repos *gradesImpl) FindGradesTeachers(codd map[string]interface{}) (*[]model.SubmissionStudent, error) {
	var result []model.SubmissionStudent

	if err := repos.DB.Where(codd).Preload("Material").Find(&result).Error; err != nil {

		return nil, err
	}

	return &result, nil

}

func (repos *gradesImpl) GetAvailableSchoolYears(codd map[string]interface{}) (*[]string, error) {
	var result []string

	if err := repos.DB.Model(&model.SubmissionStudent{}).Where(codd).Distinct().Pluck("school_year", &result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repos *gradesImpl) GetAvailableClasses(codd map[string]interface{}) (*[]string, error) {
	var result []string

	if err := repos.DB.Model(&model.SubmissionStudent{}).Where(codd).Distinct().Pluck("class", &result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
