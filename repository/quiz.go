package repository

import "github.com/cvzamannow/E-Learning-API/model"

// QuizRepository represents the interface for interacting with quiz-related data in the repository.
type QuizRepository interface {
	// CRUD Quiz

	// CreateQuiz creates a new quiz.
	CreateQuiz(request model.Quiz) (*model.Quiz, error)
	// FindQuiz finds a quiz based on provided conditions.
	FindQuiz(cond map[string]interface{}) (*model.Quiz, error)
	// DeleteQuiz deletes a quiz based on provided conditions.
	DeleteQuiz(codd map[string]interface{}) error
	// Update Quiz updates a quiz based on provided conditions.
	UpdateQuiz(codd map[string]interface{}, request model.Quiz) (*model.Quiz, error)
	// GetQuiz by submaterial id
	GetQuizByMaterialId(id string) (*model.Quiz, error)

	// CRUD Quizes

	// CreateQuizes creates multiple quizes.
	CreateQuizes(request model.Quizes) (*model.Quizes, error)
	// FindQuizes finds multiple quizes based on provided conditions.
	FindQuizes(cond map[string]interface{}) (*model.Quizes, error)
	FindQuizesNoError(cond map[string]interface{}) *model.Quizes

	// GetQuizesByIdQuiz finds multiple quizes by Quiz ID.
	GetQuizesByIdQuiz(id string) (*[]model.Quizes, error)
	// UpdateQuizes updates multiple quizes based on provided conditions.
	UpdateQuizes(codd map[string]interface{}, request model.Quizes) (*model.Quizes, error)
	// DeleteQuizes deletes multiple quizes based on provided conditions.
	DeleteQuizes(codd map[string]interface{}) error

	// CRUD QuizAnswer

	// CreateQuizAnswer creates a new quiz answer.
	CreateQuizAnswer(request model.QuizAnswer) (*model.QuizAnswer, error)
	// FindQuizAnswer finds a quiz answer based on provided conditions.
	FindQuizAnswer(cond map[string]interface{}) (*model.QuizAnswer, error)
	// GetQuizAnswerByIdQuizes finds a quiz answer by Quizes ID.
	GetQuizAnswerByIdQuizes(id string) (*[]model.QuizAnswer, error)
	// UpdateQuizAnswer updates a quiz answer based on provided conditions.
	UpdateQuizAnswer(codd map[string]interface{}, request model.QuizAnswer) (*model.QuizAnswer, error)
	// DeleteQuizAnswer deletes a quiz answer based on provided conditions.
	DeleteQuizAnswer(codd map[string]interface{}) error

	// CRUD QuizAnswer for student
	CreateQuizAnswerStudent(request model.QuizAnswerStudent) (*model.QuizAnswerStudent, error)
	FindQuizAnswerStudent(codd map[string]interface{}) (*model.QuizAnswerStudent, error)
	GetQuizAnswerStudentByIdQuiz(id string) (*[]model.QuizAnswerStudent, error)
}
