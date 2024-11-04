package repository

import (
	"errors"

	"github.com/cvzamannow/E-Learning-API/helper"
	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type courseImpl struct {
	DB *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseImpl{
		DB: db,
	}
}

func (repos *courseImpl) CreateCourse(data model.Course) (*model.Course, error) {
	result := repos.DB.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}

func (repos *courseImpl) FindCourses(
	page int,
	limit int,
	queryClass string,
	queryMajor string,
	teacherID string,
	isComplete bool,
	activeStudentID string,
	teacherPOV string,
	isActive bool,
) (*helper.Pagination, []model.Course) {

	var courses []model.Course
	tx := repos.DB.Model(&courses).
		Distinct("course_classes.course_id").
		Select("courses.id, courses.title, courses.teacher_id, courses.estimation_hour, courses.estimation_minute, courses.description, courses.detail, courses.slug, courses.thumbnail_img, courses.created_at, courses.updated_at").
		Joins("inner join course_classes on course_classes.course_id = courses.id")

	if queryClass != "" && queryMajor == "" {
		q := "%" + queryClass + "-" + "%"
		tx.Where("course_classes.slug ILIKE  ?", q)
	}

	if queryClass != "" && queryMajor != "" {
		q := queryClass + "-" + queryMajor
		tx.Where("course_classes.slug = ?", q)
	}

	if teacherID != "" {
		tx.Where("courses.teacher_id = ?", teacherID)
	}

	if activeStudentID != "" && isComplete {
		tx.Joins("inner join complete_courses on complete_courses.course_id = courses.id").
			Where("complete_courses.active_student_id = ?", activeStudentID)
	}

	if teacherPOV != "" {
		tx.Where("teacher_id = ?", teacherPOV)
	}

	if activeStudentID != "" {
		tx.Preload("CompleteCourses", "active_student_id = ?", activeStudentID) 
	}

	if isActive {
		logrus.Infoln("[repository] [Is Active Scope] Trigerred")
		var completePlaceholder []model.CompleteCourse
		txActive := repos.DB.Model(&model.CompleteCourse{}).Where("active_student_id = ?", activeStudentID).Find(&completePlaceholder)

		if txActive.RowsAffected > 0 {
			var ids []string
			for _, el := range completePlaceholder {
				ids = append(ids, el.CourseID)
			}

			tx.Not("courses.id IN ?", ids)
		}

	} 

	tx.Where("courses.is_draft = ?", false)

	logrus.Infoln("[repository] Course  Classes:", queryClass)

	// Counting rows
	var c int64
	tx.Find(&courses)
	c = tx.RowsAffected

	pagination, txPaginator := helper.Paginator(page, limit, tx, c)

	txPaginator.Preload("CourseClasses").Preload("Chapters").Order("created_at ASC").Find(&courses)
	if txPaginator.Error != nil {
		logrus.Warnln("[database] There is error:", txPaginator.Error.Error())
		return nil, nil
	}

	logrus.Infoln("[database] Total Rows:", tx.RowsAffected)

	return pagination, courses

}

func (repos *courseImpl) FindCourse(cond map[string]interface{}, isStudentPOV bool, activeStudentID string) (*model.Course, error) {
	var course model.Course

	tx := repos.DB.Model(&model.Course{})

	if isStudentPOV && activeStudentID != "" {
		logrus.Infoln("[repository] Current POV: Student")
		tx.Preload("Chapters.Materials", func(d *gorm.DB) *gorm.DB {
			return d.Order("created_at ASC").Preload("Progress", "active_student_id = ?", activeStudentID)
		}).Preload("CompleteCourses")
	} else {
		tx.Preload("Chapters.Materials", func(d *gorm.DB) *gorm.DB {
			return d.Order("created_at ASC")
		})
	}

	tx.Preload("CourseClasses").
		Where(cond).
		First(&course)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &course, nil
}

func (repos *courseImpl) DeleteCourse(cond map[string]interface{}) (*model.Course, error) {
	var course model.Course

	err := repos.DB.Where(cond).First(&course).Error

	if err != nil {
		return nil, err
	}

	err = repos.DB.Delete(&course).Error

	if err != nil {
		return nil, err
	}

	var submissionStudent []model.SubmissionStudent

	tx := repos.DB.Where("course_id = ?", course.ID).Find(&submissionStudent)
	if tx.RowsAffected > 0 {
		repos.DB.Delete(&submissionStudent)
	}

	_, err = repos.DeleteChapter(map[string]interface{}{
		"course_id": course.ID,
	})

	if err != nil {
		logrus.Warnln("[repository] Something is wrong when deleting course...")
	}

	// Do delete complete course
	var completedCourse []model.CompleteCourse
	tx = repos.DB.Where("course_id = ?", course.ID).Find(&completedCourse)

	if tx.RowsAffected > 0 {
		repos.DB.Delete(&completedCourse)
	}

	return &course, nil
}

func (repos *courseImpl) EditCourse(
	cond map[string]interface{},
	data model.Course,
) (*model.Course, error) {
	var course model.Course

	err := repos.DB.Where(cond).First(&course).Error
	if err != nil {
		return nil, err
	}

	data.ID = course.ID
	data.TeacherID = course.TeacherID
	data.CreatedAt = course.CreatedAt

	err = repos.DB.Model(&course).Save(&data).Error
	if err != nil {
		logrus.Warnln("[database] Couldn't update Course data because error:", err)
		return nil, err
	}

	var class []model.CourseClass

	tx := repos.DB.Find(&class, "course_id = ?", data.ID)

	if tx.RowsAffected > 0 {
		for _, el := range class {
			err = repos.DB.Delete(&el).Error
			if err != nil {
				return nil, err
			}
		}
	}

	for _, el := range data.CourseClasses {
		err = repos.DB.Save(&el).Error

		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}


func (repos *courseImpl) CreateChapter(data model.Chapter) (*model.Chapter, error) {
	err := repos.DB.Create(&data).Error

	if err != nil {
		return nil, err
	}

	// Undo completed progress
	var completedCourse []model.CompleteCourse 
	tx := repos.DB.Where("course_id", data.CourseID).Find(&completedCourse)

	if tx.RowsAffected > 0 {
		repos.DB.Delete(&completedCourse)
	}

	return &data, nil
}

func (repos *courseImpl) FindChapter(cond map[string]interface{}) (*model.Chapter, error) {
	var chapter model.Chapter

	err := repos.DB.Where(cond).First(&chapter).Error

	if err != nil {
		return nil, err
	}

	return &chapter, nil
}

func (repos *courseImpl) FindChapters(course_id string) ([]model.Chapter, error) {
	var chapters []model.Chapter

	err := repos.DB.Where("course_id = ?", course_id).Find(&chapters).Error

	if err != nil {
		return nil, err
	}

	return chapters, nil
}

func (repos *courseImpl) UpdateChapter(id string, data model.Chapter) (*model.Chapter, error) {
	var chapter model.Chapter

	err := repos.DB.First(&chapter, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	data.ID = chapter.ID
	data.CourseID = chapter.CourseID
	data.CreatedAt = chapter.CreatedAt

	err = repos.DB.Save(&data).Error

	if err != nil {
		return nil, err
	}

	return &data, nil

}

func (repos *courseImpl) DeleteChapter(
	cond map[string]interface{},
) (*model.Chapter, error) {
	var chapter model.Chapter
	err := repos.DB.Where(cond).First(&chapter).Error

	if err != nil {
		return nil, err
	}

	err = repos.DB.Delete(&chapter).Error

	if err != nil {
		return nil, err
	}

	var materials []model.Material
	tx := repos.DB.Where(map[string]interface{}{
		"chapter_id": chapter.ID,
	}).Find(&materials)

	if tx.RowsAffected > 0 {
		repos.DB.Delete(&materials)

		for _, el := range materials {
			if el.Type == "SUBMISSION" {
				var s []model.SubmissionStudent
				sTx := repos.DB.Where("material_id = ?", el.ID).Find(&s)

				if sTx.RowsAffected > 0 {
					repos.DB.Delete(&s)
				}
			}
		}
	}

	return &chapter, nil
}

func (repos *courseImpl) CreateMaterial(data model.Material) (*model.Material, error) {
	var chapter model.Chapter 

	tx := repos.DB.Where("id = ?", data.ChapterID).First(&chapter)

	if tx.RowsAffected == 0 || tx.Error != nil {
		logrus.Warnln("[db] Something is wrong:", tx.Error)
		return nil, errors.New("internal error")
	}

	data.CourseID = chapter.CourseID

	err := repos.DB.Create(&data).Error

	if err != nil {
		return nil, err
	}

	// Undo completed course student
	var completedCourse []model.CompleteCourse
	tx = repos.DB.Where("course_id = ?", data.CourseID).Find(&completedCourse)

	if tx.RowsAffected > 0 {
		repos.DB.Delete(&completedCourse)
	}


	return &data, nil
}

func (repos *courseImpl) FindMaterial(cond map[string]interface{}) (*model.Material, error) {
	var material model.Material

	err := repos.DB.Where(cond).First(&material).Error

	if err != nil {
		return nil, err
	}

	return &material, nil
}

func (repos *courseImpl) UpdateMaterial(
	cond map[string]interface{},
	data model.Material,
) (*model.Material, error) {
	var material model.Material

	err := repos.DB.Where(cond).First(&material).Error

	if err != nil {
		return nil, err
	}

	data.ID = material.ID
	data.ChapterID = material.ChapterID
	data.CourseID = material.CourseID
	data.Type = material.Type
	data.CreatedAt = material.CreatedAt

	err = repos.DB.Save(&data).Error

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repos *courseImpl) DeleteMaterial(cond map[string]interface{}) (*model.Material, error) {
	var material model.Material

	err := repos.DB.Where(cond).First(&material).Error

	if err != nil {
		return nil, err
	}

	err = repos.DB.Delete(&material).Error

	if err != nil {
		return nil, err
	}

	if material.Type == "SUBMISSION" {
		// Delete Submission Student
		var submissionStudent []model.SubmissionStudent
		tx := repos.DB.Where("material_id = ?", material.ID).Find(&submissionStudent)

		if tx.RowsAffected > 0 {
			repos.DB.Delete(&submissionStudent)
		}

		
	}

	return &material, nil
}

func (repos *courseImpl) FindMaterials(cond map[string]interface{}) ([]model.Material, error) {
	var materials []model.Material

	err := repos.DB.Where(cond).
		Order("created_at ASC").
		Find(&materials).
		Error

	if err != nil {
		logrus.Warnln("[database] Failed to retrieve Materials because error:", err)
		return nil, err
	}

	return materials, nil
}

func (repos *courseImpl) CreateTheory(
	data model.Theory,
) (*model.Theory, error) {
	err := repos.DB.Create(&data).Error

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repos *courseImpl) FindTheory(
	cond map[string]interface{},
) (*model.Theory, error) {
	var theory model.Theory

	err := repos.DB.Where(cond).First(&theory).Error

	if err != nil {
		return nil, err
	}

	return &theory, nil
}

func (repos *courseImpl) CreateCourseClass(data model.CourseClass) (*model.CourseClass, error) {
	tx := repos.DB.Create(&data)

	if tx.Error != nil {
		tx.Debug()
		return nil, tx.Error
	}
	return &data, nil
}

func (repos *courseImpl) DeleteCourseClass(codd map[string]interface{}) (*model.CourseClass, error) {

	var courseClass model.CourseClass

	tx := repos.DB.Where(codd).First(&courseClass)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = repos.DB.Delete(&courseClass)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &courseClass, nil

}

func (repos *courseImpl) CreateSubmission(data model.Submission) (*model.Submission, error) {
	err := repos.DB.Create(&data).Error

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repos *courseImpl) FindSubmission(cond map[string]interface{}) (*model.Submission, error) {
	var submission model.Submission

	err := repos.DB.Where(cond).First(&submission).Error
	if err != nil {
		return nil, err
	}

	return &submission, nil
}

func (repos *courseImpl) EditSubmission(
	id string,
	data model.Submission,
) (*model.Submission, error) {
	var submission model.Submission

	// Tx Begin
	tx := repos.DB.Where("material_id = ?", id).First(&submission)

	if tx.Error != nil {
		return nil, tx.Error
	}

	data.ID = submission.ID
	data.MaterialID = submission.MaterialID
	data.CreatedAt = submission.CreatedAt
	tx = repos.DB.Save(&data)

	if tx.Error != nil {
		// Rollback
		return nil, tx.Error

	}

	return &data, nil
}

// TODO: Create Submissions Student

func (repos *courseImpl) CreateSubmissionStudent(
	data model.SubmissionStudent,
) (*model.SubmissionStudent, error) {
	if err := repos.DB.Preload("ActiveStudent").Create(&data).Error; err != nil {
		logrus.Warningln("[Database] Failed to create submission student: ", err.Error())
		return nil, err

	}

	return &data, nil
}

func (repos *courseImpl) FindSubmissionStudent(condition map[string]interface{}, isSubmitted bool) (*model.SubmissionStudent, error) {
	var submissionStudent model.SubmissionStudent

	tx := repos.DB.Preload("ActiveStudent").Preload("ActiveStudent.Student").Preload("Material").
		Where(condition)

	if isSubmitted {
		tx.Not("status = ?", "REJECTED")
	}

	tx.First(&submissionStudent)

	if tx.RowsAffected == 0 || tx.Error != nil {
		logrus.Infoln("[repository] No record")
		return nil, errors.New("no record")
	}

	return &submissionStudent, nil
}

func (repos *courseImpl) FindSubmissionStudents(condition map[string]interface{}) ([]model.SubmissionStudent, error) {
	var submissions []model.SubmissionStudent

	tx := repos.DB.Model(&model.SubmissionStudent{}).Preload("ActiveStudent").Preload("Material").
		Where(condition).Find(&submissions)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, errors.New("no rows")
	}

	return submissions, nil
}

func (repos *courseImpl) DeleteSubmission(material_id string) (*model.Submission, error) {
	var submission model.Submission

	tx := repos.DB.Where("material_id = ?", material_id).First(&submission)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = repos.DB.Delete(&submission)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &submission, nil
}

func (repos *courseImpl) ApproveSubmission(
	codd map[string]interface{},
	grade int,
	teacherID string,
) (*model.SubmissionStudent, error) {
	var submission model.SubmissionStudent

	whereClause := repos.DB.Where(codd)

	tx := whereClause.First(&submission)
	if tx.Error != nil {
		return nil, tx.Error
	}

	submission.Status = "APPROVED"
	submission.Grade = grade
	submission.TeacherID = teacherID

	tx = repos.DB.Save(&submission)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &submission, nil
}

func (repos *courseImpl) RejectionsSubmission(
	codd map[string]interface{},
	comment *string,
	teacherID string,
) (*model.SubmissionStudent, error) {
	var submission model.SubmissionStudent

	tx := repos.DB.Where(codd).First(&submission)

	if tx.Error != nil {
		return nil, tx.Error
	}

	submission.Status = "REV_REJECT"
	submission.Comment = comment
	submission.TeacherID = teacherID

	tx = repos.DB.Save(&submission)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &submission, nil
}

func (repos *courseImpl) EnrollCourse(
	data model.ActiveStudentCourse,
) (*model.ActiveStudentCourse, error) {

	if err := repos.DB.Create(&data).Error; err != nil {
		logrus.Println("[Databases] Failed to enroll course")
		return nil, err
	}

	return &data, nil

}

func (repos *courseImpl) FindEnrollCourse(
	codd map[string]interface{},
) (*[]model.ActiveStudentCourse, error) {

	var enrollCourse []model.ActiveStudentCourse
	tx := repos.DB.Model(&model.ActiveStudentCourse{}).Preload("ActiveStudent").Preload("Material").Preload("Course").Where(codd).Find(&enrollCourse)

	if tx.RowsAffected == 0 || tx.Error != nil {
		logrus.Warnln("[error] Something is wrong:", tx.Error)
		return nil, errors.New("no rows")
	}

	logrus.Infoln("[repository] Enroll Rows:", len(enrollCourse))
	return &enrollCourse, nil

}

func (repos *courseImpl) SaveCompleteCourse(
	data model.CompleteCourse,
) (*model.CompleteCourse, error) {
	// First find the existing data
	var complete model.CompleteCourse

	tx := repos.DB.Where(map[string]interface{}{
		"active_student_id": data.ActiveStudentID,
		"course_id": data.CourseID,
	}).Find(&complete)

	if tx.RowsAffected > 0 {
		return &complete, nil
	}

	tx = repos.DB.Create(&data)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &data, nil
}

func (repos *courseImpl) EditTheory(material_id string, data model.Theory) (*model.Theory, error) {
	var theory model.Theory

	tx := repos.DB.First(&theory, "material_id = ?", material_id)

	if tx.Error != nil {
		logrus.Warnln("[database] Error in edit theory", tx.Error)
		return nil, tx.Error
	}

	data.ID = theory.ID
	data.MaterialID = theory.MaterialID

	tx = repos.DB.Save(&data)

	if tx.Error != nil {
		logrus.Warnln("[database] Error in update theory", tx.Error)
		return nil, tx.Error
	}

	return &data, nil

}

func (repos *courseImpl) DeleteTheory(cond map[string]interface{}) (*model.Theory, error) {
	var theory model.Theory
	tx := repos.DB.Where(cond).First(&theory)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = repos.DB.Delete(&theory)
	return &theory, nil
}

func (repos *courseImpl) NextMaterial(v *helper.NextMaterialArg, isOtherChapter bool) (*model.Material, error) {
	// Placeholder data
	var m model.Material

	tx := repos.DB.Order("created_at ASC").Where("id = ?", v.CurrentMaterialID).First(&m)

	if tx.RowsAffected == 0 {
		return nil, errors.New("no record")
	}

	logrus.Infoln("[db][func: NextMaterial] Rows:", tx.RowsAffected)

	var result model.Material

	if isOtherChapter {
		tx = repos.DB.Order("created_at ASC").Where("chapter_id = ?", v.ChapterID).First(&result)
	} else {
		tx = repos.DB.Order("created_at ASC").Where("chapter_id = ?", v.ChapterID).Where("created_at > ?", m.CreatedAt).First(&result)
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &result, nil
}

func (repos *courseImpl) FindSubmissionPreload(teacher_id string, active_student_id string, filterTeacherID string, filterCourseID string, filterMaterialID string, filterClass string) ([]model.SubmissionStudent, error) {
	if teacher_id != "" && active_student_id == "" {
		var submissions []model.SubmissionStudent
		tx := repos.DB.Model(&model.SubmissionStudent{}).Preload("Material").
			Preload("ActiveStudent").
			Where("teacher_id = ?", teacher_id)

		if filterMaterialID != "" {
			tx.Where("material_id = ?", filterMaterialID)
		}

		if filterClass != "" {
			tx.Where("class = ?", filterClass)
		}

		tx.Find(&submissions)

		if tx.Error != nil {
			logrus.Warnln("[repository][func: FindSubmissionPreload] Something went wrong in executing query:", tx.Error)
			return nil, tx.Error
		}

		return submissions, nil
	}

	var submissions []model.SubmissionStudent

	tx := repos.DB.Model(&model.SubmissionStudent{}).Preload("Material")
		
	if filterTeacherID != "" {
		tx.Where("teacher_id = ?", filterTeacherID)
	}

	if filterCourseID != "" {
		tx.Where("course_id = ?", filterCourseID)
	}

	tx.Where("active_student_id = ?", active_student_id)
	tx.Find(&submissions)

	if tx.Error != nil {
		logrus.Warnln("[repository][func: FindSubmissionPreload] Something went wrong in executing query:", tx.Error)
		return nil, tx.Error
	}

	return submissions, nil
}

func (repos *courseImpl) NextChapter(chapter_id string, course_id string) *model.Chapter {
	var previousChapter model.Chapter

	tx := repos.DB.Model(&model.Chapter{}).Order("created_at ASC").Where("id = ?", chapter_id).First(&previousChapter)

	if tx.RowsAffected == 0 {
		logrus.Warnln("[repository] Couldn't find previous Chapter")
		return nil
	}

	var nextChapter model.Chapter
	tx = repos.DB.Model(&model.Chapter{}).Order("created_at ASC").Where("course_id = ?", course_id).Where("created_at > ?", previousChapter.CreatedAt).First(&nextChapter)

	if tx.RowsAffected == 0 {
		logrus.Warnln("[repository] Couldn't find next chapter")
	}

	logrus.Infoln("[repository][func: NextChapter] Rows found:", nextChapter.Title)
	return &nextChapter
}

func (repos *courseImpl) ResetSubmittedSubmission(cond map[string]interface{}) error {
	tx := repos.DB.Model(&model.SubmissionStudent{}).Where(cond).Update("status", "REJECTED")

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}


func (repos *courseImpl) PlaceholderFilterSubmission() ([]model.SubmissionStudent, error) {
	var submissionStudent []model.SubmissionStudent

	tx := repos.DB.Model(&model.SubmissionStudent{}).Distinct("material_id").Find(&submissionStudent)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, errors.New("no rows")
	}

	return submissionStudent, nil

}

