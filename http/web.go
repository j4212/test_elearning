package http

import (
	"time"

	"github.com/cvzamannow/E-Learning-API/helper"
	"github.com/cvzamannow/E-Learning-API/model"
)

type WebResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type CertificateHTTP struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CourseHTTP struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Teacher          *string  `json:"teacher,omitempty"`
	Description      string   `json:"description"`
	Detail           string   `json:"detail"`
	IsDraft          bool     `json:"is_draft"`
	IsComplete       *bool    `json:"is_complete,omitempty"`
	Classes          []string `json:"classes,omitempty"`
	TotalChapter     *int     `json:"total_chapter,omitempty"`
	TotalStudent     *int     `json:"total_student"`
	EstimationHour   string   `json:"estimation_hour"`
	EstimationMinute string   `json:"estimation_minute"`
	Slug             string   `json:"slug"`
	ThumbnailImg     string   `json:"thumbnail_img"`
}

type CourseDetailHTTP struct {
	Course   CourseHTTP    `json:"course"`
	Chapters []ChapterHTTP `json:"chapters"`
}

type CoursePaginationResponse struct {
	Entries    []CourseHTTP       `json:"entries"`
	Pagination *helper.Pagination `json:"pagination"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register struct {
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Status   string     `json:"status"`
	Role     model.ROLE `json:"role"`
}

// User Management Request
type UserManagement struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type EditUserManagementRequest struct {
	Id string `json:"id"`
}

type ChapterHTTP struct {
	ID        string         `json:"id"`
	CourseID  string         `json:"course_id"`
	Title     string         `json:"title"`
	Slug      string         `json:"slug"`
	Materials []MaterialHTTP `json:"materials,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type MaterialHTTP struct {
	ID         string    `json:"id"`
	ChapterID  string    `json:"chapter_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	IsLock     *bool     `json:"is_lock,omitempty"`
	IsComplete *bool     `json:"is_complete,omitempty"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NextMaterialHTTP struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type TheoryHTTP struct {
	ID        string            `json:"id"`
	ChapterID string            `json:"chapter_id"`
	Title     string            `json:"title"`
	Slug      string            `json:"slug"`
	Type      *string           `json:"type,omitempty"`
	Content   string            `json:"content"`
	Next      *NextMaterialHTTP `json:"next"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Student struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IDNumber  int    `json:"id_number"`
	SchoolsID string `json:"school_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Teacher struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	IDNumber int     `json:"id_number"`
	Username string  `json:"username"`
	SchoolID *string `json:"school_id"`
	Password string  `json:"password"`
}

type TeacherHTTP struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	IDNumber int     `json:"id_number"`
	Username string  `json:"username"`
	SchoolID *string `json:"school_id"`
}

type ActiveStudent struct {
	ID         string `json:"id"`
	StudentID  string `json:"student_id"`
	SchoolYear string `json:"school_year"`
	Class      string `json:"class"`
}

type Schools struct {
	ID         string `json:"id"`
	SchoolYear string `json:"school_year"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Logo       string `json:"logo"`
}
type Certificate struct {
	RecipientName     string `json:"RecipientName"`
	CourseName        string`json:"CourseName"`
	CertificateNo      string `json:"CertificateNo "`
	Score             string `json:"Score"`
	CompletionEndDate time.Time `json:"CompletionEndDate"`
	CertificateUrl    string `json:"CertificateUrl"`
}


// quizz
type Quiz struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	CourseID    string     `json:"course_id"`
	ChapterID   string     `json:"chapter_id"`
	Description string     `json:"description"`
	Type        *string    `json:"type,omitempty"`
	Quizes      []QuizData `json:"quizzes"`
}

type QuizData struct {
	ID      string       `json:"id"`
	Quiz    string       `json:"quiz"`
	QuizID  string       `json:"quiz_id"`
	ImgURL  *string      `json:"img_url,omitempty"`
	Delete  *bool        `json:"delete,omitempty"`
	Answers []QuizAnswer `json:"answers"`
}

type QuizAnswer struct {
	ID        string `json:"id"`
	Answer    string `json:"answer"`
	QuizesID  string `json:"quizes_id"`
	Delete    *bool  `json:"delete,omitempty"`
	IsCorrect bool   `json:"is_correct"`
}

// End Quizz

type Admin struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SchoolID string `json:"school_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClassHTTP struct {
	Class []string `json:"classes"`
}

type SubmittedSubmissionHTTP struct {
	ID          string `json:"id"`
	FileUrl     string `json:"attachment_file"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Grade       int    `json:"grade"`
	Comment     string `json:"comment"`
	Date        string `json:"date"`
}

type HistorySubmissionHTTP struct {
	ID          string `json:"id"`
	FileUrl     string `json:"attachment_file"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Grade       int    `json:"grade"`
	Comment     string `json:"comment"`
	Date        string `json:"date"`
}

type SubmissionHTTP struct {
	ID                string                   `json:"id"`
	ChapterID         *string                  `json:"chapter_id"`
	Title             string                   `json:"title"`
	Slug              string                   `json:"slug"`
	Type              *string                  `json:"type,omitempty"`
	Content           string                   `json:"content"`
	Next              *NextMaterialHTTP        `json:"next"`
	Submitted         *SubmittedSubmissionHTTP `json:"submitted"`
	HistorySubmission []HistorySubmissionHTTP  `json:"history"`
	Date              time.Time                `json:"date"`
}
type SubmissionStudent struct {
	ID          string `json:"id"`
	MaterialID  string `json:"material_id"`
	CourseID    string `json:"course_id"`
	ChapterID   string `json:"chapter_id"`
	FileURL     string `json:"file_url"`
	Description string `json:"description"`
}

type SubmissionStudentAprrove struct {
	SubmissionID string `json:"submission_id"`
	Grade        int    `json:"grade"`
}

type SubmissionStudentReject struct {
	SubmissionID string `json:"submission_id"`
	Comment      string `json:"comment"`
}

type SubmissionStudentDetailStudent struct {
	ID              string  `json:"id"`
	Material        string  `json:"material"`
	Chapter         string  `json:"chapter"`
	Grade           int     `json:"grade"`
	ActiveStudentID string  `json:"active_student_id"`
	FileURL         string  `json:"file_url"`
	Status          string  `json:"status"`
	Date            string  `json:"date"`
	Comment         *string `json:"comment"`
	Description     string  `json:"description"`
}

type SubmissionStudentDetailTeacher struct {
	ID              string `json:"id"`
	Material        string `json:"material"`
	Chapter         string `json:"chapter"`
	Grade           int    `json:"grade"`
	ActiveStudentID string `json:"active_student_id"`
	FileURL         string `json:"file_url"`
	Status          string `json:"status"`
	Date            string `json:"date"`
	Student         string `json:"student"`
	Description     string `json:"description"`
}

// Grades

type GradeStudents struct {
	ID         string `json:"id"`
	Course     string `json:"course"`
	Grade      int    `json:"grade"`
	StudentID  string `json:"student_id"`
	Material   string `json:"material"`
	Type       string `json:"type"`
	MaterialID string `json:"material_id"`
	Date       string `json:"date"`
}

type GradeStudentResponse struct {
	Data       []GradeStudents `json:"grade"`
	Count      int             `json:"count"`
	Class      []string        `json:"class"`
	SchoolYear []string        `json:"school_year"`
}

type GradesResponse struct {
	Data    []StudentData `json:"data"`
	Student StudentInfo   `json:"student"`
}

type StudentData struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Grades     []Submission `json:"grades"`
	Class      string       `json:"class"`
	SchoolYear string       `json:"school_year"`
	Average    float64      `json:"average"`
}

type Submission struct {
	ID      string  `json:"id"`
	Course  string  `json:"course"`
	Grades  float64 `json:"grades"`
	Chapter string  `json:"chapter"`
}

type StudentInfo struct {
	SchoolYears   []map[string]string `json:"school_years"`
	Classes       []map[string]string `json:"classes"`
	TotalStudents int                 `json:"total_students"`
	AverageCourse map[string]float64  `json:"average_course"`
}

// Request quiz answers student
type Answer struct {
	QuizAnswerID string `json:"quiz_answer_id"`
	QuizesID     string `json:"quizes_id"`
}
type QuizAnswerStudent struct {
	ID     string   `json:"id"`
	Answer []Answer `json:"answer"`
}

// response

type QuizAnswerStudentResponseHTTP struct {
	ID            string  `json:"id"`
	QuizID        string  `json:"quiz_id"`
	QuizesID      string  `json:"quizes_id"`
	Quiz          string  `json:"quiz"`
	Answer        string  `json:"answer"`
	AnswerCorrect *string `json:"answer_correct"`
}

// response get quiz detail

type QuizResponseHTTP struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Grades      int                  `json:"grades"`
	Score       int                  `json:"score"`
	Quiz        []QuizesResponseHTTP `json:"quizes"`
	Next        *NextMaterialHTTP    `json:"next"`
	CreatedAt   *time.Time           `json:"created_at"`
	UpdatedAt   *time.Time           `json:"updated_at"`
}

type QuizesResponseHTTP struct {
	ID            string                      `json:"id"`
	Quiz          string                      `json:"quiz"`
	ImgURL        *string                     `json:"img_url,omitempty"`
	Answer        []QuizAnswerHTTP            `json:"answer"`
	AnswerStudent []AnswerStuedntResponseHTTP `json:"answer_student"`
	CreatedAt     *time.Time                  `json:"created_at"`
	UpdatedAt     *time.Time                  `json:"updated_at"`
}
type QuizAnswerHTTP struct {
	ID        string     `json:"id"`
	Answer    string     `json:"answer"`
	IsCorrect bool       `json:"is_correct"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type AnswerStuedntResponseHTTP struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Answer    string     `json:"answer"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// response student class
type StudentClassHTTP struct {
	Classes string `json:"classes"`
}

// Student Course
type EnrollCourseHTTP struct {
	MaterialID string `json:"material_id"`
	CourseID   string `json:"course_id"`
}

type EnrollCourseResponseHTTP struct {
	ID         string `json:"id"`
	Student    string `json:"student"`
	Class      string `json:"class"`
	SchoolYear string `json:"school_year"`
	Material   string `json:"material"`
	Course     string `json:"course"`
}

// grades teacher

type GraadeResponseHTTP struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Material string `json:"material"`
	Grade    int    `json:"grade"`
}
type GraadesStudentHTTP struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Class      string               `json:"class"`
	SchoolYear string               `json:"school_year"`
	Avarage    string               `json:"avarage"`
	Grades     []GraadeResponseHTTP `json:"grades"`
}

type AvarageResponseHTTP struct {
	Material string `json:"material"`
	Avarage  string `json:"avarage"`
}

type StatusSubmissionHTTP struct {
	Pending  int `json:"pending"`
	Rejected int `json:"rejected"`
	Approved int `json:"approved"`
}

type SubmissionStudentHTTP struct {
	ID              string    `json:"id"`
	SubmissionID    string    `json:"submission_id"`
	StudentID       string    `json:"student_id"`
	StudentName     string    `json:"student_name"`
	Class           string    `json:"class"`
	CourseTitle     string    `json:"course_title"`
	ChapterTitle    string    `json:"chapter_title"`
	SubmissionTitle string    `json:"submission_title"`
	File            string    `json:"submission_file"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	Comment         string    `json:"comment"`
	Date            time.Time `json:"date"`
}

type SubmissionStudentPOVUserHTTP struct {
	ID              string    `json:"id"`
	SubmissionID    string    `json:"submission_id"`
	CourseID        string    `json:"course_id"`
	ChapterID       string    `json:"chapter_id"`
	CourseTitle     string    `json:"course_title"`
	Status          string    `json:"status"`
	SubmissionTitle string    `json:"submission_title"`
	ChapterTitle    string    `json:"chapter_title"`
	Description     string    `json:"description"`
	File            string    `json:"submission_file"`
	Comment         string    `json:"comment"`
	Date            time.Time `json:"date"`
}

type ListSubmissionResponseHTTP struct {
	Status      StatusSubmissionHTTP    `json:"status_count"`
	Submissions []SubmissionStudentHTTP `json:"submissions"`
}

type ListSubmissionStudentPOVUserHTTP struct {
	Status      StatusSubmissionHTTP           `json:"status_count"`
	Submissions []SubmissionStudentPOVUserHTTP `json:"submissions"`
}

type ResetSubmittedRequestHTTP struct {
	SubmittedID string `json:"submitted_id"`
}

type TCoursePlaceholder struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TMaterialPlaceholder struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TSubmissionPlaceholder struct {
	Course       TCoursePlaceholder     `json:"course"`
	ListMaterial []TMaterialPlaceholder `json:"list_material"`
}
