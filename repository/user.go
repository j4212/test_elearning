package repository

import (
	"github.com/cvzamannow/E-Learning-API/model"
)

type UserRepository interface {
	// User
	CreateUser(request model.User) (*model.User, error)
	FindUser(cond map[string]interface{}) (*model.User, error)

	// Teacher
	CreateTeacher(model.Teacher) (*model.Teacher, error)
	FindTeacher(cond map[string]interface{}) (*model.Teacher, error)
	GetAllTeacher() (*[]model.Teacher, error)
	UpdateTeacher(codd map[string]interface{}, request model.Teacher) (*model.Teacher, error)
	DeleteTeacher(codd map[string]interface{}) (*model.Teacher, error)
	FindTeachers(codd map[string]interface{}) (*[]model.Teacher, error)

	// CRUD Student
	CreateStudent(request model.Student) (*model.Student, error)
	FindStudent(codd map[string]interface{}) (*model.Student, error)
	GetAllStudent() (*[]model.Student, error)
	UpdateStudent(codd map[string]interface{}, request model.Student) (*model.Student, error)
	DeleteStudent(codd map[string]interface{}) (*model.Student, error)

	// CRUD Active Student
	CreateActiveStudent(request model.ActiveStudent) (*model.ActiveStudent, error)
	// * Note (RageNeko26)
	// I'm adding 1 method for querying Active Student with map condition.
	// Also the method name have 's' at  the end of the word, it means the student is more than one.
	FindActiveStudents(cond map[string]interface{}) ([]model.ActiveStudent, error)
	// * End of note
	FindActiveStudent(codd map[string]interface{}) (*model.ActiveStudent, error)
	FindActiveStudentNoError(codd map[string]interface{}) *model.ActiveStudent
	GetAllActiveStudent() (*[]model.ActiveStudent, error)
	UpdateActiveStudent(codd map[string]interface{}, request model.ActiveStudent) (*model.ActiveStudent, error)
	DeleteActiveStudent(codd map[string]interface{}) (*model.ActiveStudent, error)

	// CRUD Admin
	CreateAdmin(request model.AdminSchool) (*model.AdminSchool, error)
	FindAdminSchool(codd map[string]interface{}) (*model.AdminSchool, error)
	GetAllAdminSchool() (*[]model.AdminSchool, error)
	UpdateAdminSchool(codd map[string]interface{}, request model.AdminSchool) (*model.AdminSchool, error)
	DeleteAdminSchool(codd map[string]interface{}) (*model.AdminSchool, error)
}
