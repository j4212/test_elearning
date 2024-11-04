package repository

import (
	"github.com/cvzamannow/E-Learning-API/helper"
	"github.com/cvzamannow/E-Learning-API/model"
)

/*
* This folder handler logic in persistance layer.
* Simply, this is where our query database is doing their job.

* Tips: Param in the method SHOULD NOT BE POINTER.
* Because it has potential error 'POINTER DEFERENCE' when its value is nil.
 */
type CourseRepository interface {
	// Course
	CreateCourse(data model.Course) (*model.Course, error)
	CreateCourseClass(data model.CourseClass) (*model.CourseClass, error)
	DeleteCourseClass(codd map[string]interface{}) (*model.CourseClass, error)
	FindCourses(
		page int,
		limit int,
		queryClass string,
		queryMajor string,
		teacherID string,
		isComplete bool,
		activeStudentID string,
		teacherPOV string,
		isActive bool,
	) (*helper.Pagination, []model.Course)
	FindCourse(cond map[string]interface{}, isStudentPOV bool, activeStudentID string) (*model.Course, error)
	DeleteCourse(cond map[string]interface{}) (*model.Course, error)
	EditCourse(cond map[string]interface{}, data model.Course) (*model.Course, error)

	// Chapter
	CreateChapter(data model.Chapter) (*model.Chapter, error)
	FindChapter(cond map[string]interface{}) (*model.Chapter, error)
	FindChapters(course_id string) ([]model.Chapter, error)
	UpdateChapter(id string, data model.Chapter) (*model.Chapter, error)
	DeleteChapter(cond map[string]interface{}) (*model.Chapter, error)

	// Material
	CreateMaterial(data model.Material) (*model.Material, error)
	FindMaterial(cond map[string]interface{}) (*model.Material, error)
	FindMaterials(cond map[string]interface{}) ([]model.Material, error)
	UpdateMaterial(cond map[string]interface{}, data model.Material) (*model.Material, error)
	DeleteMaterial(cond map[string]interface{}) (*model.Material, error)

	// Study Material
	CreateTheory(
		data model.Theory,
	) (*model.Theory, error)
	FindTheory(cond map[string]interface{}) (*model.Theory, error)
	EditTheory(material_id string, data model.Theory) (*model.Theory, error)
	DeleteTheory(cond map[string]interface{}) (*model.Theory, error)

	// Submission
	CreateSubmission(data model.Submission) (*model.Submission, error)
	FindSubmission(cond map[string]interface{}) (*model.Submission, error)
	DeleteSubmission(id string) (*model.Submission, error)
	EditSubmission(id string, data model.Submission) (*model.Submission, error)

	CreateSubmissionStudent(data model.SubmissionStudent) (*model.SubmissionStudent, error)
	FindSubmissionStudent(condition map[string]interface{}, isSubmitted bool) (*model.SubmissionStudent, error)
	FindSubmissionStudents(condition map[string]interface{}) ([]model.SubmissionStudent, error)

	// Aprrove Rejections for Submissions
	ApproveSubmission(codd map[string]interface{}, grade int, teacherID string) (*model.SubmissionStudent, error)
	RejectionsSubmission(
		cond map[string]interface{},
		comment *string,
		teacherID string,
	) (*model.SubmissionStudent, error)

	// EnrollCourse
	EnrollCourse(data model.ActiveStudentCourse) (*model.ActiveStudentCourse, error)
	FindEnrollCourse(codd map[string]interface{}) (*[]model.ActiveStudentCourse, error)

	// Complete Course
	SaveCompleteCourse(data model.CompleteCourse) (*model.CompleteCourse, error)

	// Next Action Material
	NextMaterial(v *helper.NextMaterialArg, isOtherChapter bool) (*model.Material, error)
	NextChapter(chapter_id string, course_id string) *model.Chapter

	FindSubmissionPreload(teacher_id string, student_id string, filterTeacherID string, courseID string, materialID string, filterClass string) ([]model.SubmissionStudent, error)

	ResetSubmittedSubmission(cond map[string]interface{}) error

	PlaceholderFilterSubmission() ([]model.SubmissionStudent, error)
}
