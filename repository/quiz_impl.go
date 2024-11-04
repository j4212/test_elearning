package repository

// Import necessary packages and libraries
import (
	"errors"

	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// quizImpl implements the QuizRepository interface.
type quizImpl struct {
	DB *gorm.DB
}

// NewQuizRepository creates a new instance of QuizRepository.
func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizImpl{
		DB: db,
	}
}

// CreateQuiz creates a new quiz.
func (repos *quizImpl) CreateQuiz(request model.Quiz) (*model.Quiz, error) {
	// Check if Chapter exists
	var subChapter model.Chapter
	if err := repos.DB.Model(&subChapter).Where("id = ?", request.ChapterID).First(&subChapter).Error; err != nil {
		logrus.Warningln("[DATABASE] SubMaterial not found", err.Error)
		return nil, errors.New("[DATABASE] SubMaterial not found")
	}

	// Create Quiz
	if err := repos.DB.Create(&request).Error; err != nil {
		logrus.Warningln("[DATABASE] Error creating Quiz")
		return nil, errors.New("[DATABASE] Error creating Quiz")
	}

	return &request, nil

}

func (repos *quizImpl) UpdateQuiz(codd map[string]interface{}, request model.Quiz) (*model.Quiz, error) {
	// Update Quiz
	if err := repos.DB.Model(&request).Where(codd).Updates(&request).Error; err != nil {
		logrus.Warningln("[DATABASE] Error updating Quiz")
		return nil, errors.New("[DATABASE] Error updating Quiz")
	}

	return &request, nil
}

func (repos *quizImpl) GetQuizByMaterialId(id string) (*model.Quiz, error) {
	var quiz model.Quiz
	if err := repos.DB.Where("material_id = ?", id).First(&quiz).Error; err != nil {
		logrus.Warningln("[DATABASE] Quizes not found", err.Error)
		return nil, errors.New("[DATABASE] Quizes not found")
	}

	return &quiz, nil
}

// FindQuiz finds a quiz.
func (repos *quizImpl) FindQuiz(codd map[string]interface{}) (*model.Quiz, error) {
	var quiz model.Quiz
	if err := repos.DB.Where(codd).First(&quiz).Error; err != nil {
		logrus.Warningln("[DATABASE] Quiz not found")
		return nil, errors.New("[DATABASE] Quiz not found")
	}

	return &quiz, nil
}

// DeleteQuiz deletes a quiz.
func (repos *quizImpl) DeleteQuiz(codd map[string]interface{}) error {
	// checking if quiz exists
	var quiz model.Quiz
	if err := repos.DB.Where(codd).First(&quiz).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warningln("[DATABASE] Quiz not found")
			return errors.New("[DATABASE] Quiz not found")
		}
		return err
	}

	var quizes []model.Quizes
	if err := repos.DB.Where("quiz_id = ?", quiz.ID).Find(&quizes).Error; err != nil {
		return err
	}

	for _, q := range quizes {
		if err := repos.DB.Where("quizes_id = ?", q.ID).Delete(&model.QuizAnswer{}).Error; err != nil {
			return err
		}
	}

	if err := repos.DB.Where("quiz_id = ?", quiz.ID).Delete(&model.Quizes{}).Error; err != nil {
		return err
	}

	if err := repos.DB.Delete(&quiz).Error; err != nil {
		return err
	}

	return nil
}

// CRUD Quizes
// CreateQuizes creates multiple quizes.
func (repos *quizImpl) CreateQuizes(request model.Quizes) (*model.Quizes, error) {

	var quiz model.Quiz
	if err := repos.DB.Model(&model.Quiz{}).Where("id = ?", request.QuizID).First(&quiz).Error; err != nil {
		logrus.Warningln("[DATABASE] Quiz not found")
		return nil, errors.New("[DATABASE] Quiz not found")
	}

	if err := repos.DB.Create(&request).Error; err != nil {
		logrus.Warningln("[DATABASE] Error creating Quizes")
		return nil, errors.New("[DATABASE] Error creating Quizes")
	}

	return &request, nil
}

// Find Quizes
// FindQuizes finds multiple quizes based on conditions.
func (repos *quizImpl) FindQuizes(cond map[string]interface{}) (*model.Quizes, error) {
	var quizes model.Quizes
	if err := repos.DB.Where(cond).First(&quizes).Error; err != nil {
		logrus.Warningln("[DATABASE] Quizes not found")
		return nil, errors.New("[DATABASE] Quizes not found")
	}

	return &quizes, nil
}

func (repos *quizImpl) FindQuizesNoError(cond map[string]interface{}) *model.Quizes {
	var quizes model.Quizes
	repos.DB.Where(cond).First(&quizes)

	return &quizes
}

// Get Quizes By Id Quiz
// GetQuizesByIdQuiz finds multiple quizes by Quiz ID.
func (repos *quizImpl) GetQuizesByIdQuiz(id string) (*[]model.Quizes, error) {
	var quizes []model.Quizes
	if err := repos.DB.Preload("QuizAnswers").Order("created_at ASC").Where("quiz_id = ?", id).Find(&quizes).Error; err != nil {
		logrus.Warningln("[DATABASE] Quizes not found")
		return nil, errors.New("[DATABASE] Quizes not found")
	}

	return &quizes, nil
}

// Update Quizes
// UpdateQuizes updates multiple quizes based on conditions.
func (repos *quizImpl) UpdateQuizes(codd map[string]interface{}, request model.Quizes) (*model.Quizes, error) {
	if err := repos.DB.Model(&model.Quizes{}).Where(codd).Updates(request).Error; err != nil {
		logrus.Warningln("[DATABASE] Error updating Quizes")
		return nil, errors.New("[DATABASE] Error updating Quizes")
	}

	return &request, nil
}

// Delete Quizes
// DeleteQuizes deletes multiple quizes based on conditions.
func (repos *quizImpl) DeleteQuizes(codd map[string]interface{}) error {
	var quizes model.Quizes
	if err := repos.DB.Model(&quizes).Where(codd).Delete(&quizes).Error; err != nil {
		logrus.Warningln("[DATABASE] Quizes not found")
		return errors.New("[DATABASE] Quizes not found")
	}

	return nil
}

// CRUD QuizAnswers

// Create QuizAnswers
// CreateQuizAnswer creates a new quiz answer.
func (repos *quizImpl) CreateQuizAnswer(request model.QuizAnswer) (*model.QuizAnswer, error) {
	var quizAnswer model.Quizes
	if err := repos.DB.Model(&model.Quizes{}).Where("id = ?", request.QuizesID).First(&quizAnswer).Error; err != nil {
		logrus.Warningln("[DATABASE] QuizAnswer not found")
		return nil, errors.New("[DATABASE] QuizAnswer not found")
	}

	if err := repos.DB.Create(&request).Error; err != nil {
		logrus.Warningln("[DATABASE] Error creating QuizAnswer")
		return nil, errors.New("[DATABASE] Error creating QuizAnswer")
	}

	return &request, nil
}

// FindQuizAnswer
// FindQuizAnswer finds a quiz answer based on conditions.
func (repos *quizImpl) FindQuizAnswer(cond map[string]interface{}) (*model.QuizAnswer, error) {
	var quizAnswer model.QuizAnswer
	if err := repos.DB.Where(cond).First(&quizAnswer).Error; err != nil {
		logrus.Warningln("[DATABASE] QuizAnswer not found")
		return nil, errors.New("[DATABASE] QuizAnswer not found")
	}

	return &quizAnswer, nil
}

// get Quiz Answer by id Quizes
// GetQuizAnswerByIdQuizes finds a quiz answer by Quizes ID.
func (repos *quizImpl) GetQuizAnswerByIdQuizes(id string) (*[]model.QuizAnswer, error) {
	var quizAnswer []model.QuizAnswer
	if err := repos.DB.Model(&model.QuizAnswer{}).Preload("QuizAnswerStudents").Order("created_at ASC").Where("quizes_id = ?", id).Find(&quizAnswer).Error; err != nil {
		logrus.Warningln("[DATABASE] QuizAnswer not found")
		return nil, errors.New("[DATABASE] QuizAnswer not found")
	}

	return &quizAnswer, nil
}

// Update QuizAnswers
// UpdateQuizAnswer updates a quiz answer based on conditions.
func (repos *quizImpl) UpdateQuizAnswer(codd map[string]interface{}, request model.QuizAnswer) (*model.QuizAnswer, error) {

	var quizAnswer model.QuizAnswer
	if err := repos.DB.Model(&quizAnswer).Where(codd).First(&quizAnswer).Error; err != nil {
		return nil, err
	}
	quizAnswer.Answer = request.Answer
	quizAnswer.IsCorrect = request.IsCorrect
	if err := repos.DB.Save(&quizAnswer).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// Delete QuizAnswers
// DeleteQuizAnswer deletes a quiz answer based on conditions.
func (repos *quizImpl) DeleteQuizAnswer(codd map[string]interface{}) error {
	var quizAnswer model.QuizAnswer
	if err := repos.DB.Model(&quizAnswer).Where(codd).Delete(&quizAnswer).Error; err != nil {
		logrus.Warningln("[DATABASE] QuizAnswer not found")
		return errors.New("[DATABASE] QuizAnswer not found")
	}

	return nil
}

// Create quiz answer student
func (repos *quizImpl) CreateQuizAnswerStudent(request model.QuizAnswerStudent) (*model.QuizAnswerStudent, error) {

	var quizCheck model.Quiz
	if err := repos.DB.Model(model.Quiz{}).Where("id = ?", request.QuizID).First(&quizCheck).Error; err != nil {
		return nil, err
	}

	if err := repos.DB.Create(&request).Error; err != nil {
		return nil, err
	}

	return &request, nil
}

// find QuizAnswer Students
func (repos *quizImpl) FindQuizAnswerStudent(codd map[string]interface{}) (*model.QuizAnswerStudent, error) {
	var quizAnswerStudent model.QuizAnswerStudent

	// Subquery untuk mencari entri terbaru
	subquery := repos.DB.Model(&model.QuizAnswerStudent{}).
		Select("MAX(created_at)").
		Where(codd)

	// Query utama untuk mendapatkan data berdasarkan subquery
	if err := repos.DB.Where("created_at = (?)", subquery).
		Preload("QuizAnswer").
		Find(&quizAnswerStudent).
		Error; err != nil {
		return nil, err
	}

	return &quizAnswerStudent, nil
}

func (repos *quizImpl) GetQuizAnswerStudentByIdQuiz(id string) (*[]model.QuizAnswerStudent, error) {
	var quizAnswerStudent []model.QuizAnswerStudent
	if err := repos.DB.Preload("QuizAnswer").Preload("ActiveStudent").Preload("ActiveStudent.Student").Order("created_at ASC").Where("quiz_id = ?", id).Find(&quizAnswerStudent).Error; err != nil {
		return nil, err
	}

	return &quizAnswerStudent, nil
}
