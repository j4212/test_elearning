package repository

import "github.com/cvzamannow/E-Learning-API/model"

type GradesRepository interface {

	// Student
	GetGradesStudent(codd map[string]interface{}, class, schoolYear, schoolID string) (*[]model.SubmissionStudent, error)
	GetQuizGradesStudent(codd map[string]interface{}) (*[]model.QuizAnswerStudent, error)

	// Teacher
	GetGradesTeacher(schoolYear, class, schoolID string) (*[]model.SubmissionStudent, error)
	FindGradesTeachers(codd map[string]interface{}) (*[]model.SubmissionStudent, error)

	// get all the school year data for which there is only data
	GetAvailableSchoolYears(codd map[string]interface{}) (*[]string, error)
	GetAvailableClasses(codd map[string]interface{}) (*[]string, error)
}
