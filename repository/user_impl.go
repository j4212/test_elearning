package repository

import (
	"errors"

	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type userImpl struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) UserRepository {
	return &userImpl{
		DB: db,
	}
}

// FindUser implements AuthRepository.
// Use map[string]interface{} for flexbility query.
func (a *userImpl) FindUser(cond map[string]interface{}) (*model.User, error) {
	// check whether the username and password in the database are correct
	var user model.User
	err := a.DB.Where(cond).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil

}

// Create User implements AuthRepository.
func (a *userImpl) CreateUser(request model.User) (*model.User, error) {
	// check if the user is already registered
	var user model.User
	rows := a.DB.Where("username = ?", request.Username).First(&user).RowsAffected

	if rows != 0 {
		return nil, errors.New("already registered")
	}

	// create new user
	err := a.DB.Create(&request).Error

	if err != nil {
		return nil, err
	}

	// response
	return &request, nil

}

// TODO: Crud teacher implements UserRepository interface

func (repos *userImpl) CreateTeacher(data model.Teacher) (*model.Teacher, error) {

	var teacher model.Teacher

	if err := repos.DB.Where("id_number = ?", data.IdNumber).First(&teacher).Error; err == nil {
		logrus.Warnln("[database] Failed Teacher already exists:", err)
		return nil, errors.New("Teacher already exists")
	}

	err := repos.DB.Create(&data).Error

	if err != nil {
		logrus.Warnln("[database] Failed to insert data because:", err)
		return nil, err
	}

	return &data, nil
}

func (repos *userImpl) FindTeacher(codd map[string]interface{}) (*model.Teacher, error) {
	var teacher model.Teacher

	err := repos.DB.Preload("User").First(&teacher, codd).Error

	if err != nil {
		logrus.Warnln("[database] Failed to find records Teacher:", err.Error())
		return nil, err
	}

	return &teacher, nil
}

func (repos *userImpl) DeleteTeacher(codd map[string]interface{}) (*model.Teacher, error) {
	var teacher model.Teacher

	if err := repos.DB.Preload("User").Where(codd).Delete(&teacher).Error; err != nil {
		return nil, err
	}

	return &teacher, nil
}

// GetAllTeacher implements UserRepository.
func (repos *userImpl) GetAllTeacher() (*[]model.Teacher, error) {
	var teacher []model.Teacher

	if err := repos.DB.Preload("User").Find(&teacher).Error; err != nil {
		return nil, err
	}

	return &teacher, nil
}

// UpdateTeacher implements UserRepository.
func (repos *userImpl) UpdateTeacher(codd map[string]interface{}, request model.Teacher) (*model.Teacher, error) {

	tx := repos.DB.Begin()

	var teacher model.Teacher
	if err := repos.DB.Preload("User").Where(codd).First(&teacher).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&teacher).Updates(&request).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	user := &teacher.User

	user.Username = request.User.Username
	user.Password = request.User.Password

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// TODO: CRUD Student implements UserRepository interface

// CreateStudent implements UserRepository.
func (repos *userImpl) CreateStudent(request model.Student) (*model.Student, error) {
	var student model.Student

	if err := repos.DB.Where("id_number = ?", request.IdNumber).First(&student).Error; err == nil {
		logrus.Warnln("[database] Failed Student already exists:", err)
		return nil, errors.New("Student already exists")
	}

	err := repos.DB.Preload("Schools").Create(&request).Error

	if err != nil {
		logrus.Warnln("[database] Failed to insert data because:", err)
		return nil, err
	}

	return &request, nil
}

// DeleteStudent implements UserRepository.
func (repos *userImpl) DeleteStudent(codd map[string]interface{}) (*model.Student, error) {
	var student model.Student

	if err := repos.DB.Preload("User").Where(codd).Delete(&student).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

// FindStudent implements UserRepository.
func (repos *userImpl) FindStudent(codd map[string]interface{}) (*model.Student, error) {
	var student model.Student

	if err := repos.DB.Preload("User").First(&student, codd).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

// GetAllStudent implements UserRepository.
func (repos *userImpl) GetAllStudent() (*[]model.Student, error) {
	var student []model.Student

	if err := repos.DB.Preload("User").Find(&student).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

// UpdateStudent implements UserRepository.
func (repos *userImpl) UpdateStudent(codd map[string]interface{}, request model.Student) (*model.Student, error) {
	tx := repos.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var student model.Student
	if err := tx.Preload("User").Where(codd).First(&student).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&student).Updates(&request).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	user := &student.User

	user.Username = request.User.Username
	user.Password = request.User.Password
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &request, nil
}

// CreateActiveStudent implements UserRepository.
func (repos *userImpl) CreateActiveStudent(request model.ActiveStudent) (*model.ActiveStudent, error) {

	// Cheking student_id is existing or not existing
	var student model.Student
	if err := repos.DB.Where("id = ? ", request.StudentID).First(&student).Error; err != nil {
		return nil, errors.New("Student not already exists")
	}

	err := repos.DB.Create(&request).Error

	if err != nil {
		return nil, err
	}

	return &request, nil

}

func (repos *userImpl) FindActiveStudents(cond map[string]interface{}) ([]model.ActiveStudent, error) {
	var activeStudent []model.ActiveStudent

	err := repos.DB.Where(cond).Find(&activeStudent).Error

	if err != nil {
		return nil, err
	}

	return activeStudent, nil
}

// DeleteActiveStudent implements UserRepository.
func (repos *userImpl) DeleteActiveStudent(codd map[string]interface{}) (*model.ActiveStudent, error) {

	var activeStudent model.ActiveStudent
	if err := repos.DB.Where(codd).Delete(&activeStudent).Error; err != nil {
		return nil, err
	}

	return &activeStudent, nil
}

// FindActiveStudent implements UserRepository.
func (repos *userImpl) FindActiveStudent(codd map[string]interface{}) (*model.ActiveStudent, error) {

	var activeStudent model.ActiveStudent
	if err := repos.DB.Preload("Student").Where(codd).Last(&activeStudent).Error; err != nil {
		return nil, err
	}

	return &activeStudent, nil
}

func (repos *userImpl) FindActiveStudentNoError(codd map[string]interface{}) *model.ActiveStudent {

	var activeStudent model.ActiveStudent
	if err := repos.DB.Preload("Student").Where(codd).Last(&activeStudent).Error; err != nil {
		return nil
	}

	return &activeStudent
}

// GetAllActiveStudent implements UserRepository.
func (repos *userImpl) GetAllActiveStudent() (*[]model.ActiveStudent, error) {

	var activeStudent []model.ActiveStudent
	if err := repos.DB.Preload("Student").Find(&activeStudent).Error; err != nil {
		return nil, err
	}

	return &activeStudent, nil

}

// UpdateActiveStudent implements UserRepository.
func (repos *userImpl) UpdateActiveStudent(codd map[string]interface{}, request model.ActiveStudent) (*model.ActiveStudent, error) {
	tx := repos.DB.Begin()

	if err := repos.DB.Preload("Student").Where(codd).Updates(&request).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &request, nil
}

// CRUD Admin

func (repos *userImpl) CreateAdmin(request model.AdminSchool) (*model.AdminSchool, error) {
	var admin model.User

	if err := repos.DB.Model(&admin).Where("username = ?", request.User.Username).First(&admin).Error; err == nil {
		logrus.Warnln("[database] Failed User already exist:", err)
		return nil, errors.New("User already exist")
	}

	err := repos.DB.Preload("User").Preload("Schools").Create(&request).Error

	if err != nil {
		logrus.Warnln("[database] Failed to insert data because:", err)
		return nil, err
	}

	return &request, nil
}

func (repos *userImpl) FindAdminSchool(codd map[string]interface{}) (*model.AdminSchool, error) {
	var admin model.AdminSchool

	if err := repos.DB.Preload("Schools").Preload("User").Model(&admin).Where(codd).First(&admin).Error; err != nil {
		logrus.Warningln("[database] Failed to find admin school")
		return nil, err
	}

	return &admin, nil
}

func (repos *userImpl) GetAllAdminSchool() (*[]model.AdminSchool, error) {

	var admins []model.AdminSchool

	if err := repos.DB.Preload("User").Preload("Schools").Model(&admins).Find(&admins).Error; err != nil {
		logrus.Warningln("[database] Failed get all admin school")
		return nil, err
	}

	return &admins, nil

}

func (repos *userImpl) UpdateAdminSchool(codd map[string]interface{}, request model.AdminSchool) (*model.AdminSchool, error) {
	tx := repos.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var admin model.AdminSchool
	if err := tx.Preload("User").Preload("Schools").Where(codd).First(&admin).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&admin).Updates(&request).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	user := &admin.User

	user.Username = request.User.Username
	user.Password = request.User.Password
	user.Status = request.User.Status
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &request, nil
}

func (repos *userImpl) DeleteAdminSchool(codd map[string]interface{}) (*model.AdminSchool, error) {
	var admin model.AdminSchool

	if err := repos.DB.Model(&admin).Where(codd).First(&admin).Error; err != nil {
		logrus.Warningln("[database] Failed find to admin school")
		return nil, err
	}

	if err := repos.DB.Delete(&admin).Error; err != nil {
		logrus.Warningln("[database] Failed to delete admin school")
		return nil, err
	}

	return &admin, nil
}

func (repos *userImpl) FindTeachers(codd map[string]interface{}) (*[]model.Teacher, error) {
	var teachers []model.Teacher

	if err := repos.DB.Preload("User").Model(&teachers).Where(codd).Find(&teachers).Error; err != nil {
		return nil, err
	}

	return &teachers, nil
}
