package handlers

import (
	"strconv"

	"github.com/cvzamannow/E-Learning-API/http"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *Handlers) RouteGrades(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Get("/grades-student", h.Middleware.Protected(), h.GetGradesStudent)
	v1.Get("/grades-teacher", h.Middleware.Protected(), h.GetGradesTeacher)
}

// grades student

func (h *Handlers) GetGradesStudent(c *fiber.Ctx) error {
	class := c.Query("class")

	schoolYear := c.Query("school_year")

	// role := c.Locals("role").(string)
	studentID := c.Locals("user_id").(string)

	findActiveStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": studentID,
	})
	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Active Student not found",
			Data:    nil,
		})
	}

	findStudent, err := h.UserRepository.FindStudent(map[string]interface{}{
		"id": findActiveStudent.StudentID,
	})

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Student not found",
			Data:    nil,
		})
	}

	result, err := h.GradesRepository.GetGradesStudent(map[string]interface{}{
		"active_student_id": findActiveStudent.ID,
	}, class, schoolYear, findStudent.SchoolsID)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to retrieve grades",
			Data:    nil,
		})
	}

	var response []http.GradeStudents

	for _, item := range *result {
		// find quiz answer student
		resultAnswerStudent, err := h.GradesRepository.GetQuizGradesStudent(map[string]interface{}{
			"active_student_id": findActiveStudent.ID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Could not find",
				Data:    nil,
			})
		}

		for _, quiz := range *resultAnswerStudent {
			response = append(response, http.GradeStudents{
				ID:         quiz.ID,
				Course:     item.Course,
				Grade:      quiz.Grades,
				StudentID:  findActiveStudent.StudentID,
				Date:       quiz.CreatedAt.Format("02 January 2006"),
				Material:   quiz.Quiz.Material.Title,
				Type:       quiz.Quiz.Material.Type,
				MaterialID: quiz.Quiz.Material.ID,
			})
		}

		response = append(response, http.GradeStudents{
			ID:         item.ID,
			Course:     item.Course,
			Grade:      item.Grade,
			StudentID:  findActiveStudent.StudentID,
			Date:       item.CreatedAt.Format("02 January 2006"),
			Material:   item.Material.Title,
			Type:       item.Material.Type,
			MaterialID: item.Material.ID,
		})

	}

	resultSchoolYear, err := h.GradesRepository.GetAvailableSchoolYears(map[string]interface{}{
		"school_id": findActiveStudent.Student.SchoolsID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error retrieving available school year: " + err.Error(),
			Data:    nil,
		})
	}

	resultClass, err := h.GradesRepository.GetAvailableClasses(map[string]interface{}{
		"school_id": findActiveStudent.Student.SchoolsID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error retrieving available class: " + err.Error(),
			Data:    nil,
		})
	}

	resultResponse := http.GradeStudentResponse{
		Data:       response,
		Count:      len(response),
		Class:      *resultClass,
		SchoolYear: *resultSchoolYear,
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Successfully getting grades for student",
		Data:    resultResponse,
	})
}

// GetGradesTeacher handles the request for retrieving grades by teacher
func (h *Handlers) GetGradesTeacher(c *fiber.Ctx) error {

	teacherID := c.Locals("user_id").(string)

	// // find teacher
	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"user_id": teacherID,
	})

	logrus.Println(findTeacher)
	logrus.Println(teacherID)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("teacher", "retrieve"),
			Data:    nil,
		})
	}

	schoolYear := c.Query("school_year")
	class := c.Query("class")

	result, err := h.GradesRepository.GetGradesTeacher(schoolYear, class, *findTeacher.SchoolsID)
	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error retrieving grades for teacher: " + err.Error(),
			Data:    nil,
		})
	}

	// untuk menyimpan rata rata nilai
	averageMap := make(map[string]int)
	countMap := make(map[string]int)

	studentMap := make(map[string]http.GraadesStudentHTTP)

	for _, grade := range *result {
		var grades []http.GraadeResponseHTTP

		student, err := h.GradesRepository.FindGradesTeachers(map[string]interface{}{
			"active_student_id": grade.ActiveStudentID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Error retrieving grades for student: " + err.Error(),
				Data:    nil,
			})
		}

		var gradeStudent int
		for _, item := range *student {
			if grade.ActiveStudentID == item.ActiveStudentID {
				grades = append(grades, http.GraadeResponseHTTP{
					ID:       item.ID,
					Type:     item.Material.Type,
					Material: item.Material.Title,
					Grade:    item.Grade,
				})
				gradeStudent += item.Grade
			}
		}

		// Menghitung rata-rata nilai siswa
		averageStudent := 0
		if len(*student) > 0 {
			averageStudent = gradeStudent / len(*student)
		}

		// menghitung nilai rata-rata per material
		averageMap[grade.Material.Title] += averageStudent
		countMap[grade.Material.Title]++

		studentResponse, ok := studentMap[grade.ActiveStudentID]
		if !ok {
			studentResponse = http.GraadesStudentHTTP{
				ID:         grade.ActiveStudentID,
				Name:       grade.ActiveStudent.Student.Name,
				Class:      grade.ActiveStudent.Class,
				SchoolYear: grade.ActiveStudent.SchoolYear,
				Avarage:    strconv.Itoa(averageStudent),
				Grades:     grades,
			}
			studentMap[grade.ActiveStudentID] = studentResponse
		} else {
			studentResponse.Grades = append(studentResponse.Grades, grades...)
			studentMap[grade.ActiveStudentID] = studentResponse
		}
	}

	// Mengonversi peta nilai rata-rata menjadi slice untuk response
	var avarageResponse []http.AvarageResponseHTTP
	for material, average := range averageMap {
		count := countMap[material]
		avarageResponse = append(avarageResponse, http.AvarageResponseHTTP{
			Material: material,
			Avarage:  strconv.Itoa(average / count),
		})
	}

	// Konversi peta siswa menjadi slice awokawokawok
	var studentSlice []http.GraadesStudentHTTP
	for _, student := range studentMap {
		studentSlice = append(studentSlice, student)
	}

	resultSchoolYear, err := h.GradesRepository.GetAvailableSchoolYears(map[string]interface{}{
		"school_id": findTeacher.SchoolsID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error retrieving available school year: " + err.Error(),
			Data:    nil,
		})
	}

	resultClass, err := h.GradesRepository.GetAvailableClasses(map[string]interface{}{
		"school_id": findTeacher.SchoolsID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error retrieving available class: " + err.Error(),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Successfuly get all grade information",
		Data: map[string]interface{}{
			"student":     studentSlice,
			"avarage":     avarageResponse,
			"school_year": resultSchoolYear,
			"class":       resultClass,
		},
	})
}
