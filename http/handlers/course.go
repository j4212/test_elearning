package handlers

import (
	"fmt"
	"strconv"

	"github.com/cvzamannow/E-Learning-API/helper"
	"github.com/cvzamannow/E-Learning-API/http"
	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
)

/*
* E Learning API (C) 2024
* This file contains related endpoint routes
* Route should starts with prefix version 'v1' and following with plural nouns.
* Example: /api/v1/courses
 */

/*
* Route setup should naming with pascal case
* Tips: If use word Route before your route name
* Ex: RouteCourses -> Add word 'Route' before `Courses`
 */

func (h *Handlers) RouteCourses(app *fiber.App) {
	v1 := app.Group("/api/v1")
	// Course Routes
	v1.Post("/courses", h.Middleware.Protected(), h.AddCourse)
	v1.Get("/courses", h.Middleware.Protected(), h.GetCourses)
	v1.Get("/courses/detail/:id", h.Middleware.Protected(), h.GetDetailCourseByID)
	v1.Put("/courses/:id", h.Middleware.Protected(), h.EditCourse)
	v1.Delete("/courses/:id", h.Middleware.Protected(), h.DeleteCourse)

	// Chapter Routes
	v1.Post("/chapters", h.Middleware.Protected(), h.CreateChapter)
	v1.Get("/chapters/:id", h.FindChapter)
	v1.Put("/chapters/:id", h.Middleware.Protected(), h.UpdateChapter)
	v1.Delete("/chapters/:id", h.Middleware.Protected(), h.DeleteChapter)

	// Theories
	v1.Post("/theories", h.Middleware.Protected(), h.CreateTheory)
	v1.Get("/theories/:id", h.Middleware.Protected(), h.FindTheory)
	v1.Delete("/theories/:id", h.Middleware.Protected(), h.DeleteTheory)
	v1.Put("/theories/:id", h.Middleware.Protected(), h.EditTheory)

	// Submission
	v1.Post("/submission", h.Middleware.Protected(), h.CreateSubmission)
	v1.Get("/submission/:id", h.Middleware.Protected(), h.FindSubmission)
	v1.Put("/submission/:id", h.Middleware.Protected(), h.EditSubmission)
	v1.Delete("/submission/:id", h.Middleware.Protected(), h.DeleteSubmission)
	v1.Get("/submission/detail/:id", h.Middleware.Protected(), h.GetDetailSubmissionTeacher)

	// Submission Student
	v1.Post("/submission-student", h.Middleware.Protected(), h.CreateSubmissionStudent)
	v1.Get("/submission-student/:student_id", h.Middleware.Protected(), h.FindStudentSubmission)
	v1.Get("/submission-student/detail/:id", h.Middleware.Protected(), h.GetDetailSubmission)
	v1.Get("/submission-student", h.Middleware.Protected(), h.ListingSubmission)
	v1.Put("/reset-submission", h.Middleware.Protected(), h.ResetSubmission)
	v1.Get("/placeholder/submission-student", h.Middleware.Protected(), h.GetSubmissionPlaceholder)

	// Submission Student Approve and Reject
	v1.Put("/submission-student/approve/:student_id", h.Middleware.Protected(), h.Approvesubmisson)
	v1.Put("/submission-student/reject/:student_id", h.Middleware.Protected(), h.Rejectsubmisson)

	// Mencari kelas student yang login entah itu sudah kelas 11 atau kelas 10
	v1.Get("/class/student", h.Middleware.Protected(), h.FindClassStudent)

	// Update Progress
	v1.Put("/progress", h.Middleware.Protected(), h.UpdateProgress)
	// v1.Get("/enroll-courses/:slug", h.Middleware.Protected(), h.GetEnrollCourse)

	// Complete Course
	v1.Put("/complete-course", h.Middleware.Protected(), h.SaveCompleteCourse)
}

/*
* Handlers name should be explicit and not using buzz word.
* Tips: Use verb for naming function
* Ex: AddCourses(c *fiber.Ctx),
 */

func (h *Handlers) AddCourse(c *fiber.Ctx) error {
	var request http.CourseHTTP
	c.BodyParser(&request)

	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "TEACHER" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Couldn't access resource because unauthorized request!",
			Data:    nil,
		})
	}

	id, _ := gonanoid.New(20)

	var class []model.CourseClass

	for _, el := range request.Classes {
		idClass, _ := gonanoid.New(20)
		class = append(class, model.CourseClass{
			ID:    idClass,
			Class: el,
			Slug:  slug.Make(el),
		})
	}

	isDraft := true

	if !request.IsDraft {
		isDraft = false
	}

	if request.ThumbnailImg == "" {
		request.ThumbnailImg = "https://pub-c883652fcf2a4a4f9b0d7321a986a773.r2.dev/Thumbnail.png"
	}

	// Find Teacher to return the response.
	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"user_id": userID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "fail",
			Message: "Couldn't create course because internal error",
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.CreateCourse(model.Course{
		ID:               id,
		Title:            request.Title,
		Description:      request.Description,
		ThumbnailImg:     request.ThumbnailImg,
		Detail:           request.Detail,
		EstimationHour:   request.EstimationHour,
		IsDraft:          isDraft,
		EstimationMinute: request.EstimationMinute,
		Slug:             slug.Make(request.Title),
		TeacherID:        findTeacher.ID,
		CourseClasses:    class,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "fail",
			Message: "Couldn't create course because internal error",
			Data:    nil,
		})
	}

	var classResp []string

	for _, el := range res.CourseClasses {
		classResp = append(classResp, el.Class)
	}

	dataRes := http.CourseHTTP{
		ID:               res.ID,
		Title:            res.Title,
		Description:      res.Description,
		ThumbnailImg:     res.ThumbnailImg,
		EstimationHour:   res.EstimationHour,
		EstimationMinute: res.EstimationMinute,
		IsDraft:          res.IsDraft,
		Detail:           res.Detail,
		Slug:             res.Slug,
		Classes:          classResp,
		Teacher:          &findTeacher.Name,
	}

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Course has been created!",
		Data:    dataRes,
	})
}

func (h *Handlers) GetCourses(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	class := c.Query("class")
	major := c.Query("major")
	teacher_id := c.Query("teacher_id")
	isCompleteQuery := c.Query("is_complete")

	if c.Query("page") == "" || c.Query("limit") == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Request must specify with query params 'page', 'limit' and 'class'",
			Data:    nil,
		})
	}
	// Default value
	isComplete := false
	activeStudent := ""

	isComplete, err := strconv.ParseBool(isCompleteQuery)

	if err != nil {
		isComplete = false
	}

	if c.Locals("role").(string) == "STUDENT" {
		// Find Active StudentID
		student, err := h.UserRepository.FindStudent(map[string]interface{}{
			"user_id": c.Locals("user_id"),
		})

		// If not found then Default Value is not changing
		if err != nil {
			activeStudent = ""
		}

		activeStudentData, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
			"student_id": student.ID,
		})

		if err != nil {
			activeStudent = ""
		}

		activeStudent = activeStudentData.ID

		queryIsActive := c.Query("is_active")

		isActiveVal := false

		isActive, err := strconv.ParseBool(queryIsActive)

		if err == nil {
			isActiveVal = isActive
		}

		logrus.Infoln("[handler] Is Active:", isActiveVal)
		pagination, data := h.CourseRepository.FindCourses(
			page,
			limit,
			class,
			major,
			teacher_id,
			isComplete,
			activeStudent,
			"",
			isActiveVal,
		)

		if pagination == nil || data == nil {
			return c.Status(404).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Couldn't find courses data because it's not found!",
				Data:    nil,
			})
		}

		var entries []http.CourseHTTP
		logrus.Infoln("[handlers] Pagination has been triggered")

		for _, el := range data {
			logrus.Infoln("[handler] Teacher ID:", el.TeacherID)
			findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
				"id": el.TeacherID,
			})

			if err != nil {
				return c.Status(500).JSON(&http.WebResponse{
					Status:  "error",
					Message: "Failed to retrieve Courses because internal error",
					Data:    nil,
				})
			}

			var classResp []string

			totalStudent := 0
			for _, j := range el.CourseClasses {
				logrus.Info("[handlers] Class Data:", j.Class)
				classResp = append(classResp, j.Class)

				student, err := h.UserRepository.FindActiveStudents(map[string]interface{}{
					"class": j.Class,
				})

				if err != nil && len(student) == 0 {
					totalStudent = 0
				}

				totalStudent += len(student)

			}

			var totalChapter = 0

			if len(el.Chapters) > 0 {
				totalChapter += len(el.Chapters)
			}

			thumbnail := "https://pub-c883652fcf2a4a4f9b0d7321a986a773.r2.dev/Thumbnail.png"

			if el.ThumbnailImg != "" {
				thumbnail = el.ThumbnailImg
			}

			var complete *bool

			if len(el.CompleteCourses) > 0 {
				status := true
				complete = &status
			}

			entries = append(entries, http.CourseHTTP{
				ID:               el.ID,
				Title:            el.Title,
				Slug:             el.Slug,
				Description:      el.Description,
				Teacher:          &findTeacher.Name,
				Detail:           el.Detail,
				IsDraft:          el.IsDraft,
				TotalStudent:     &totalStudent,
				IsComplete:       complete,
				EstimationHour:   el.EstimationHour,
				Classes:          classResp,
				TotalChapter:     &totalChapter,
				EstimationMinute: el.EstimationMinute,
				ThumbnailImg:     thumbnail,
			})
		}

		resp := http.CoursePaginationResponse{
			Pagination: pagination,
			Entries:    entries,
		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Successfully retrieve Courses!",
			Data:    resp,
		})
	}

	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"user_id": c.Locals("user_id").(string),
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("course", "retrieve"),
			Data:    []interface{}{},
		})
	}

	pagination, data := h.CourseRepository.FindCourses(
		page,
		limit,
		"",
		"",
		"",
		false,
		"",
		findTeacher.ID,
		false,
	)

	if pagination == nil || data == nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Couldn't find courses data because it's not found!",
			Data:    nil,
		})
	}

	var entries []http.CourseHTTP
	logrus.Infoln("[handlers] Pagination has been triggered")

	for _, el := range data {
		logrus.Infoln("[handler] Teacher ID:", el.TeacherID)
		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": el.TeacherID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Failed to retrieve Courses because internal error",
				Data:    nil,
			})
		}

		var classResp []string

		totalStudent := 0
		for _, j := range el.CourseClasses {
			logrus.Info("[handlers] Class Data:", j.Class)
			classResp = append(classResp, j.Class)

			student, err := h.UserRepository.FindActiveStudents(map[string]interface{}{
				"class": j.Class,
			})

			if err != nil && len(student) == 0 {
				totalStudent = 0
			}

			totalStudent += len(student)

		}

		var totalChapter = 0

		if len(el.Chapters) > 0 {
			totalChapter += len(el.Chapters)
		}

		thumbnail := "https://pub-c883652fcf2a4a4f9b0d7321a986a773.r2.dev/Thumbnail.png"

		if el.ThumbnailImg != "" {
			thumbnail = el.ThumbnailImg
		}

		entries = append(entries, http.CourseHTTP{
			ID:               el.ID,
			Title:            el.Title,
			Slug:             el.Slug,
			Description:      el.Description,
			Teacher:          &findTeacher.Name,
			Detail:           el.Detail,
			IsDraft:          el.IsDraft,
			TotalStudent:     &totalStudent,
			IsComplete:       &isComplete,
			EstimationHour:   el.EstimationHour,
			Classes:          classResp,
			TotalChapter:     &totalChapter,
			EstimationMinute: el.EstimationMinute,
			ThumbnailImg:     thumbnail,
		})
	}

	resp := http.CoursePaginationResponse{
		Pagination: pagination,
		Entries:    entries,
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Successfully retrieve Courses!",
		Data:    resp,
	})
}

func (h *Handlers) GetDetailCourseByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Request must specify 'id' params!",
			Data:    []interface{}{},
		})
	}

	role := c.Locals("role").(string)

	if role == "TEACHER" {
		res, err := h.CourseRepository.FindCourse(map[string]interface{}{
			"id": id,
		}, false, "")

		if err != nil {
			return c.Status(404).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Couldn't find course because it's not found!",
				Data:    []interface{}{},
			})
		}

		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": res.TeacherID,
		})

		if err != nil {
			return c.Status(404).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Couldn't find course because it's not found!",
				Data:    []interface{}{},
			})
		}

		var chaptersResponse []http.ChapterHTTP

		totalStudent := 0

		for _, el := range res.CourseClasses {
			students, err := h.UserRepository.FindActiveStudents(map[string]interface{}{
				"class": el.Class,
			})

			if err == nil && len(students) > 0 {
				totalStudent += len(students)
			}
		}

		if len(res.Chapters) == 0 {
			isComplete := false

			var classes []string

			if len(res.CourseClasses) > 0 {
				for _, cl := range res.CourseClasses {
					classes = append(classes, cl.Class)
				}
			}

			totalChapter := 0

			return c.Status(200).JSON(&http.WebResponse{
				Status:  "success",
				Message: h.successResponse("course", "retrieve"),
				Data: http.CourseHTTP{
					ID:               res.ID,
					Title:            res.Title,
					Teacher:          &findTeacher.Name,
					Description:      res.Description,
					Detail:           res.Detail,
					IsDraft:          res.IsDraft,
					IsComplete:       &isComplete,
					Classes:          classes,
					Slug:             res.Slug,
					ThumbnailImg:     res.ThumbnailImg,
					EstimationHour:   res.EstimationHour,
					EstimationMinute: res.EstimationMinute,
					TotalChapter:     &totalChapter,
					TotalStudent:     &totalStudent,
				},
			})
		}

		for _, chap := range res.Chapters {
			var materials []http.MaterialHTTP

			if len(chap.Materials) > 0 {
				for _, el := range chap.Materials {

					materials = append(materials, http.MaterialHTTP{
						ID:        el.ID,
						ChapterID: el.ChapterID,
						Title:     el.Title,
						Slug:      el.Slug,
						Type:      el.Type,
						CreatedAt: el.CreatedAt,
						UpdatedAt: el.UpdatedAt,
					})

				}
			}

			chaptersResponse = append(chaptersResponse, http.ChapterHTTP{
				ID:        chap.ID,
				Title:     chap.Title,
				Slug:      chap.Slug,
				CourseID:  chap.CourseID,
				Materials: materials,
				CreatedAt: res.CreatedAt,
				UpdatedAt: res.UpdatedAt,
			})
		}

		var courseClasses []string

		for _, el := range res.CourseClasses {
			courseClasses = append(courseClasses, el.Class)
		}

		totalChapter := 0

		if len(res.CourseClasses) > 0 {
			totalChapter += len(res.Chapters)
		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Successfully getting detail Course",
			Data: http.CourseDetailHTTP{
				Course: http.CourseHTTP{
					ID:               res.ID,
					Title:            res.Title,
					Teacher:          &findTeacher.Name,
					Description:      res.Description,
					Slug:             res.Slug,
					TotalChapter:     &totalChapter,
					TotalStudent:     &totalStudent,
					Detail:           res.Detail,
					EstimationHour:   res.EstimationHour,
					EstimationMinute: res.EstimationMinute,
					ThumbnailImg:     res.ThumbnailImg,
					Classes:          courseClasses,
				},
				Chapters: chaptersResponse,
			},
		})

	}

	userID := c.Locals("user_id")

	student, err := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": userID,
	})

	activeStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": student.ID,
	})

	res, err := h.CourseRepository.FindCourse(map[string]interface{}{
		"id": id,
	}, true, activeStudent.ID)

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Couldn't find course because it's not found!",
			Data:    []interface{}{},
		})
	}

	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"id": res.TeacherID,
	})

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Couldn't find course because it's not found!",
			Data:    []interface{}{},
		})
	}

	var chaptersResponse []http.ChapterHTTP

	totalStudent := 0

	for _, el := range res.CourseClasses {
		students, err := h.UserRepository.FindActiveStudents(map[string]interface{}{
			"class": el.Class,
		})

		if err == nil && len(students) > 0 {
			totalStudent += len(students)
		}
	}

	if len(res.Chapters) == 0 {
		isComplete := false

		var classes []string

		if len(res.CourseClasses) > 0 {
			for _, cl := range res.CourseClasses {
				classes = append(classes, cl.Class)
			}
		}

		totalChapter := 0

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: h.successResponse("course", "retrieve"),
			Data: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Detail:           res.Detail,
				IsDraft:          res.IsDraft,
				IsComplete:       &isComplete,
				Classes:          classes,
				Slug:             res.Slug,
				ThumbnailImg:     res.ThumbnailImg,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				TotalChapter:     &totalChapter,
				TotalStudent:     &totalStudent,
			},
		})
	}

	for _, chap := range res.Chapters {
		var materials []http.MaterialHTTP

		if len(chap.Materials) > 0 {
			for _, el := range chap.Materials {

				// Find Progress Course
				// TRUE = complete
				// FALSE = not complete
				var statusProgress bool

				defaultLock := true

				if len(el.Progress) > 0 {
					statusProgress = true
				} else {
					statusProgress = false
				}

				materials = append(materials, http.MaterialHTTP{
					ID:         el.ID,
					ChapterID:  el.ChapterID,
					Title:      el.Title,
					Slug:       el.Slug,
					IsComplete: &statusProgress,
					IsLock:     &defaultLock,
					Type:       el.Type,
					CreatedAt:  el.CreatedAt,
					UpdatedAt:  el.UpdatedAt,
				})

			}
		}

		chaptersResponse = append(chaptersResponse, http.ChapterHTTP{
			ID:        chap.ID,
			Title:     chap.Title,
			Slug:      chap.Slug,
			CourseID:  chap.CourseID,
			Materials: materials,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		})
	}

	var courseClasses []string

	for _, el := range res.CourseClasses {
		courseClasses = append(courseClasses, el.Class)
	}

	totalChapter := 0

	if len(res.CourseClasses) > 0 {
		totalChapter += len(res.Chapters)
	}

	for indexChap := 0; indexChap < len(chaptersResponse); indexChap++ {
		for indexMaterial := 0; indexMaterial < len(chaptersResponse[indexChap].Materials); indexMaterial++ {
			next := indexMaterial + 1
			current := indexMaterial

			if *chaptersResponse[indexChap].Materials[current].IsComplete {
				isLock := false
				chaptersResponse[indexChap].Materials[current].IsLock = &isLock
			}

			if next < len(chaptersResponse[indexChap].Materials) {
				if !*chaptersResponse[indexChap].Materials[next].IsComplete && *chaptersResponse[indexChap].Materials[current].IsComplete {
					isLock := false
					chaptersResponse[indexChap].Materials[next].IsLock = &isLock
				} else {
					isLock := true
					chaptersResponse[indexChap].Materials[next].IsLock = &isLock
				}
			}
		}

		currentIndexChapter := indexChap
		nextChapterIndex := indexChap + 1

		if nextChapterIndex < len(chaptersResponse) {
			m := chaptersResponse[currentIndexChapter].Materials
			lastMaterialIndex := len(m) - 1
			lastMaterial := chaptersResponse[currentIndexChapter].Materials[lastMaterialIndex]

			if *lastMaterial.IsComplete {
				nextMaterials := chaptersResponse[nextChapterIndex].Materials
				if len(nextMaterials) > 0 {
					isLock := false
					nextMaterials[0].IsLock = &isLock
				}
			}
		}

	}

	var complete *bool

	if len(res.CompleteCourses) > 0 {
		val := true
		complete = &val
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Successfully getting detail Course",
		Data: http.CourseDetailHTTP{
			Course: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Slug:             res.Slug,
				TotalChapter:     &totalChapter,
				TotalStudent:     &totalStudent,
				Detail:           res.Detail,
				IsComplete:       complete,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				ThumbnailImg:     res.ThumbnailImg,
				Classes:          courseClasses,
			},
			Chapters: chaptersResponse,
		},
	})

}

func (h *Handlers) EditCourse(c *fiber.Ctx) error {
	idParams := c.Params("id")

	if idParams == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "fail",
			Message: "Request must specify 'slug' params",
			Data:    []interface{}{},
		})
	}

	var request http.CourseHTTP
	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	logrus.Infoln("[handler] Updated Title:", request.Title)

	var courseClasses []model.CourseClass

	for _, el := range request.Classes {
		id, _ := gonanoid.New(20)

		courseClasses = append(courseClasses, model.CourseClass{
			ID:       id,
			Class:    el,
			Slug:     slug.Make(el),
			CourseID: idParams,
		})
	}

	isDraft := true

	if !request.IsDraft {
		isDraft = false
	}

	data := model.Course{
		Title:            request.Title,
		Description:      request.Description,
		ThumbnailImg:     request.ThumbnailImg,
		EstimationHour:   request.EstimationHour,
		EstimationMinute: request.EstimationMinute,
		Detail:           request.Detail,
		IsDraft:          isDraft,
		CourseClasses:    courseClasses,
		Slug:             slug.Make(request.Title),
	}

	res, err := h.CourseRepository.EditCourse(map[string]interface{}{
		"id": idParams,
	}, data)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "fail",
			Message: "Couldn't edit course because there's internal error",
			Data:    []interface{}{},
		})
	}

	logrus.Infoln("[controller] Teacher ID:", res.TeacherID)
	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"id": res.TeacherID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to retrieve Courses because internal error",
			Data:    nil,
		})

	}

	var classResp []string

	for _, el := range res.CourseClasses {
		classResp = append(classResp, el.Class)
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Course has been edited!",
		Data: http.CourseHTTP{
			ID:               res.ID,
			Title:            res.Title,
			Teacher:          &findTeacher.Name,
			Description:      res.Description,
			Slug:             res.Slug,
			Detail:           res.Detail,
			Classes:          classResp,
			EstimationHour:   res.EstimationHour,
			EstimationMinute: res.EstimationMinute,
			ThumbnailImg:     res.ThumbnailImg,
		},
	})
}

func (h *Handlers) DeleteCourse(c *fiber.Ctx) error {
	idParams := c.Params("id")

	if idParams == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "fail",
			Message: "Request must specify resource using 'slug' params",
			Data:    []interface{}{},
		})
	}

	// find chapter by id course
	findChapter, err := h.CourseRepository.FindChapter(map[string]interface{}{
		"course_id": idParams,
	})

	if err != nil {
		res, err := h.CourseRepository.DeleteCourse(map[string]interface{}{
			"id": idParams,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "fail",
				Message: "Couldn't delete courses because internal error",
				Data:    []interface{}{},
			})
		}

		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": res.TeacherID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Failed to retrieve Courses because internal error",
				Data:    nil,
			})

		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Course has been deleted!",
			Data: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Detail:           res.Detail,
				Slug:             res.Slug,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				ThumbnailImg:     res.ThumbnailImg,
			},
		})
	}

	// find material by id chapter
	findMaterialByIDChapter, err := h.CourseRepository.FindMaterial(map[string]interface{}{
		"chapter_id": findChapter.ID,
	})

	if err != nil {

		res, err := h.CourseRepository.DeleteCourse(map[string]interface{}{
			"id": idParams,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "fail",
				Message: "Couldn't delete courses because internal error",
				Data:    []interface{}{},
			})
		}

		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": res.TeacherID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Failed to retrieve Courses because internal error",
				Data:    nil,
			})

		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Course has been deleted!",
			Data: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Detail:           res.Detail,
				Slug:             res.Slug,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				ThumbnailImg:     res.ThumbnailImg,
			},
		})
	}

	// delete theory by materialID
	_, err = h.CourseRepository.DeleteTheory(map[string]interface{}{
		"material_id": findMaterialByIDChapter.ID,
	})

	if err != nil {

		res, err := h.CourseRepository.DeleteCourse(map[string]interface{}{
			"id": idParams,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "fail",
				Message: "Couldn't delete courses because internal error",
				Data:    []interface{}{},
			})
		}

		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": res.TeacherID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Failed to retrieve Courses because internal error",
				Data:    nil,
			})

		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Course has been deleted!",
			Data: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Detail:           res.Detail,
				Slug:             res.Slug,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				ThumbnailImg:     res.ThumbnailImg,
			},
		})
	}

	// delete submission by materialID
	_, err = h.CourseRepository.DeleteSubmission(findMaterialByIDChapter.ID)

	if err != nil {
		res, err := h.CourseRepository.DeleteCourse(map[string]interface{}{
			"id": idParams,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "fail",
				Message: "Couldn't delete courses because internal error",
				Data:    []interface{}{},
			})
		}

		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": res.TeacherID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Failed to retrieve Courses because internal error",
				Data:    nil,
			})

		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Course has been deleted!",
			Data: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Detail:           res.Detail,
				Slug:             res.Slug,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				ThumbnailImg:     res.ThumbnailImg,
			},
		})
	}

	// delete Course class by CurseID
	_, err = h.CourseRepository.DeleteCourseClass(map[string]interface{}{
		"course_id": idParams,
	})

	if err != nil {

		res, err := h.CourseRepository.DeleteCourse(map[string]interface{}{
			"id": idParams,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "fail",
				Message: "Couldn't delete courses because internal error",
				Data:    []interface{}{},
			})
		}

		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": res.TeacherID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Failed to retrieve Courses because internal error",
				Data:    nil,
			})

		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Course has been deleted!",
			Data: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Detail:           res.Detail,
				Slug:             res.Slug,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				ThumbnailImg:     res.ThumbnailImg,
			},
		})

	}

	// delete material by id chapter
	_, err = h.CourseRepository.DeleteMaterial(map[string]interface{}{
		"chapter_id": findChapter.ID,
	})

	if err != nil {
		res, err := h.CourseRepository.DeleteCourse(map[string]interface{}{
			"id": idParams,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "fail",
				Message: "Couldn't delete courses because internal error",
				Data:    []interface{}{},
			})
		}

		findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"id": res.TeacherID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Failed to retrieve Courses because internal error",
				Data:    nil,
			})

		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: "Course has been deleted!",
			Data: http.CourseHTTP{
				ID:               res.ID,
				Title:            res.Title,
				Teacher:          &findTeacher.Name,
				Description:      res.Description,
				Detail:           res.Detail,
				Slug:             res.Slug,
				EstimationHour:   res.EstimationHour,
				EstimationMinute: res.EstimationMinute,
				ThumbnailImg:     res.ThumbnailImg,
			},
		})
	}

	// delete chapter by id course
	_, err = h.CourseRepository.DeleteChapter(map[string]interface{}{
		"id": findChapter.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error deleting chapter by id course",
			Data:    nil,
		})
	}

	// delete course

	res, err := h.CourseRepository.DeleteCourse(map[string]interface{}{
		"id": idParams,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "fail",
			Message: "Couldn't delete courses because internal error",
			Data:    []interface{}{},
		})
	}

	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"id": res.TeacherID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to retrieve Courses because internal error",
			Data:    nil,
		})

	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Course has been deleted!",
		Data: http.CourseHTTP{
			ID:               res.ID,
			Title:            res.Title,
			Teacher:          &findTeacher.Name,
			Description:      res.Description,
			Detail:           res.Detail,
			Slug:             res.Slug,
			EstimationHour:   res.EstimationHour,
			EstimationMinute: res.EstimationMinute,
			ThumbnailImg:     res.ThumbnailImg,
		},
	})
}

func (h *Handlers) CreateChapter(c *fiber.Ctx) error {
	var request http.ChapterHTTP

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to parse body request",
			Data:    nil,
		})
	}

	id, _ := gonanoid.New(20)
	res, err := h.CourseRepository.CreateChapter(model.Chapter{
		ID:       id,
		CourseID: request.CourseID,
		Title:    request.Title,
		Slug:     slug.Make(request.Title),
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to create chapter because internal error",
			Data:    nil,
		})
	}

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Chapter has been created!",
		Data: http.ChapterHTTP{
			ID:        res.ID,
			Title:     res.Title,
			CourseID:  res.CourseID,
			Slug:      res.Slug,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})
}

func (h *Handlers) FindChapter(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Must specify resource with 'id' params",
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.FindChapter(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Chapter is not found!",
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: fmt.Sprintf("Successfully retrieve chapter %s", id),
		Data: http.ChapterHTTP{
			ID:        res.ID,
			Title:     res.Title,
			Slug:      res.Slug,
			CourseID:  res.CourseID,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})
}

func (h *Handlers) UpdateChapter(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Must specify resource with 'id' params ",
			Data:    nil,
		})
	}

	var request http.ChapterHTTP
	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.UpdateChapter(id, model.Chapter{
		Title: request.Title,
		Slug:  slug.Make(request.Title),
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("chapter", "update"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("chapter", "updated"),
		Data: &http.ChapterHTTP{
			ID:        res.ID,
			Title:     res.Title,
			Slug:      res.Slug,
			CourseID:  res.CourseID,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})
}

func (h *Handlers) DeleteChapter(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.DeleteChapter(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("chapter", "delete"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("chapter", "deleted"),
		Data: http.ChapterHTTP{
			ID:        res.ID,
			CourseID:  res.CourseID,
			Title:     res.Title,
			Slug:      res.Slug,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})
}

func (h *Handlers) CreateTheory(c *fiber.Ctx) error {
	var request http.TheoryHTTP

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	idMaterial, _ := gonanoid.New(20)

	m, err := h.CourseRepository.CreateMaterial(model.Material{
		ID:        idMaterial,
		ChapterID: request.ChapterID,
		Title:     request.Title,
		Slug:      slug.Make(request.Title),
		Type:      "THEORY",
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("material", "create"),
			Data:    nil,
		})
	}

	idTheory, _ := gonanoid.New(20)
	t, err := h.CourseRepository.CreateTheory(model.Theory{
		ID:         idTheory,
		MaterialID: m.ID,
		Content:    request.Content,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "create"),
			Data:    nil,
		})
	}

	typeOfMaterial := "THEORY"

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("theory", "created"),
		Data: http.TheoryHTTP{
			ID:        m.ID,
			Title:     m.Title,
			Slug:      m.Slug,
			Type:      &typeOfMaterial,
			Content:   t.Content,
			ChapterID: m.ChapterID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
	})
}

func (h *Handlers) FindTheory(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.FindMaterial(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "retrieve"),
			Data:    nil,
		})
	}

	t, err := h.CourseRepository.FindTheory(map[string]interface{}{
		"material_id": res.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "retrieve"),
			Data:    nil,
		})
	}

	// Next Material Action
	var nextMaterial *http.NextMaterialHTTP
	next, err := h.CourseRepository.NextMaterial(&helper.NextMaterialArg{
		ChapterID:         res.ChapterID,
		CurrentMaterialID: res.ID,
		CreatedAt:         res.CreatedAt,
	}, false)

	if err == nil && next.ChapterID == res.ChapterID {
		nextMaterial = &http.NextMaterialHTTP{
			ID:   next.ID,
			Type: next.Type,
		}
	}

	// Check if next chapter exist
	if next == nil {
		findChapterDetail, _ := h.CourseRepository.FindChapter(map[string]interface{}{
			"id": res.ChapterID,
		})
		nextChapter := h.CourseRepository.NextChapter(res.ChapterID, findChapterDetail.CourseID)

		if nextChapter != nil {
			doNext, err := h.CourseRepository.NextMaterial(&helper.NextMaterialArg{
				ChapterID:         nextChapter.ID,
				CurrentMaterialID: res.ID,
				CreatedAt:         res.CreatedAt,
			}, true)

			if err == nil {
				nextMaterial = &http.NextMaterialHTTP{
					ID:   doNext.ID,
					Type: doNext.Type,
				}
			}
		}
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("theory", "retrieve"),
		Data: http.TheoryHTTP{
			ID:        res.ID,
			ChapterID: res.ChapterID,
			Title:     res.Title,
			Slug:      res.Slug,
			Content:   t.Content,
			Next:      nextMaterial,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})

}

func (h *Handlers) DeleteMaterial(c *fiber.Ctx) error {
	id := c.Query("id")

	res, err := h.CourseRepository.DeleteMaterial(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("material", "delete"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("material", "deleted"),
		Data: http.MaterialHTTP{
			ID:        res.ID,
			ChapterID: res.ChapterID,
			Title:     res.Title,
			Type:      res.Type,
			Slug:      res.Slug,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})
}

func (h *Handlers) CreateSubmission(c *fiber.Ctx) error {
	var request http.SubmissionHTTP

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	idMaterial, err := gonanoid.New(20)

	if err != nil {
		logrus.Warnln("[handlers-submission] logtrace line 1019")
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "create"),
			Data:    nil,
		})
	}

	idSubmission, _ := gonanoid.New(20)

	m, err := h.CourseRepository.CreateMaterial(model.Material{
		ID:        idMaterial,
		Title:     request.Title,
		ChapterID: *request.ChapterID,
		Slug:      slug.Make(request.Title),
		Type:      "SUBMISSION",
		Submission: model.Submission{
			ID:         idSubmission,
			MaterialID: idMaterial,
			Content:    request.Content,
		},
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "create"),
			Data:    nil,
		})
	}

	typeOfMaterial := "SUBMISSION"

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "created"),
		Data: http.SubmissionHTTP{
			ID:        m.ID,
			Title:     m.Title,
			ChapterID: &m.ChapterID,
			Type:      &typeOfMaterial,
			Content:   request.Content,
			Slug:      m.Slug,
			Date:      m.CreatedAt,
		},
	})

}

// TODO: Crud submissions Student

func (h *Handlers) CreateSubmissionStudent(c *fiber.Ctx) error {

	studentID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "STUDENT" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	findActiveStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": studentID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("student", "retrieve"),
			Data:    nil,
		})
	}

	var request http.SubmissionStudent
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	id, err := helper.GenerateNanoId()

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "create"),
			Data:    nil,
		})
	}

	// Find Course to get teacher
	course, err := h.CourseRepository.FindCourse(map[string]interface{}{
		"id": request.CourseID,
	}, false, "")

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "create"),
			Data:    []interface{}{},
		})
	}

	result, err := h.CourseRepository.CreateSubmissionStudent(model.SubmissionStudent{
		ID:              string(id),
		MaterialID:      request.MaterialID,
		FileUrl:         request.FileURL,
		Description:     request.Description,
		ActiveStudentID: findActiveStudent.ID,
		CourseID:        course.ID,
		Status:          "PENDING",
		TeacherID:       course.TeacherID,
		Class:           findActiveStudent.Class,
		SchoolYear:      findActiveStudent.SchoolYear,
		SchoolID:        findActiveStudent.Student.SchoolsID,
		Course:          course.Title,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "create"),
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id":          string(id),
		"material_id": result.MaterialID,
		"file_url":    result.FileUrl,
		"description": result.Description,
		"student_id":  findActiveStudent.Student.ID,
		"status":      result.Status,
		"created_at":  result.CreatedAt,
		"updated_at":  result.UpdatedAt,
		"deleted_at":  result.DeletedAt,
	}

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "created"),
		Data:    response,
	})

}

func (h *Handlers) FindSubmission(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	m, err := h.CourseRepository.FindMaterial(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.FindSubmission(map[string]interface{}{
		"material_id": m.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	var nextMaterial *http.NextMaterialHTTP
	next, err := h.CourseRepository.NextMaterial(&helper.NextMaterialArg{
		ChapterID:         m.ChapterID,
		CurrentMaterialID: m.ID,
		CreatedAt:         m.CreatedAt,
	}, false)

	if err == nil && next.ChapterID == m.ChapterID {
		nextMaterial = &http.NextMaterialHTTP{
			ID:   next.ID,
			Type: next.Type,
		}
	}

	// Check if next chapter exist
	if next == nil {
		findChapterDetail, _ := h.CourseRepository.FindChapter(map[string]interface{}{
			"id": m.ChapterID,
		})
		nextChapter := h.CourseRepository.NextChapter(m.ChapterID, findChapterDetail.CourseID)

		if nextChapter != nil {
			doNext, err := h.CourseRepository.NextMaterial(&helper.NextMaterialArg{
				ChapterID:         nextChapter.ID,
				CurrentMaterialID: m.ID,
				CreatedAt:         m.CreatedAt,
			}, true)

			if err == nil {
				nextMaterial = &http.NextMaterialHTTP{
					ID:   doNext.ID,
					Type: doNext.Type,
				}
			}
		}
	}

	role := c.Locals("role").(string)

	var submitted *http.SubmittedSubmissionHTTP
	var historySubmission []http.HistorySubmissionHTTP = nil

	if role == "STUDENT" {
		userID := c.Locals("user_id")
		findStudent, _ := h.UserRepository.FindStudent(map[string]interface{}{
			"user_id": userID,
		})

		activeStudent, _ := h.UserRepository.FindActiveStudent(map[string]interface{}{
			"student_id": findStudent.ID,
		})

		submit, _ := h.CourseRepository.FindSubmissionStudent(map[string]interface{}{
			"active_student_id": activeStudent.ID,
			"material_id":       m.ID,
		}, true)

		if submit != nil {
			status := string(submit.Status)

			if status == "REV_REJECT" {
				status = "REJECTED"
			}

			comment := ""

			if submit.Comment != nil {
				comment = *submit.Comment
			}

			submitted = &http.SubmittedSubmissionHTTP{
				ID:          submit.ID,
				FileUrl:     submit.FileUrl,
				Description: submit.Description,
				Status:      status,
				Grade:       submit.Grade,
				Comment:     comment,
				Date:        submit.UpdatedAt.String(),
			}
		}

		history, _ := h.CourseRepository.FindSubmissionStudents(map[string]interface{}{
			"active_student_id": activeStudent.ID,
			"material_id":       m.ID,
			"status":            "REJECTED",
		})

		if history != nil {
			for _, el := range history {
				comment := ""

				if el.Comment != nil {
					comment = *el.Comment
				}
				historySubmission = append(historySubmission, http.HistorySubmissionHTTP{
					ID:          el.ID,
					FileUrl:     el.FileUrl,
					Grade:       el.Grade,
					Comment:     comment,
					Description: el.Description,
					Status:      string(el.Status),
					Date:        el.UpdatedAt.String(),
				})
			}
		}
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "retrieve"),
		Data: http.SubmissionHTTP{
			ID:                m.ID,
			ChapterID:         &m.ChapterID,
			Title:             m.Title,
			Slug:              m.Slug,
			Content:           res.Content,
			Next:              nextMaterial,
			HistorySubmission: historySubmission,
			Submitted:         submitted,
			Date:              res.UpdatedAt,
		},
	})
}

func (h *Handlers) EditSubmission(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	var request http.SubmissionHTTP
	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	m, err := h.CourseRepository.FindMaterial(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("sumbmission", "edit"),
			Data:    nil,
		})
	}

	editM, err := h.CourseRepository.UpdateMaterial(map[string]interface{}{
		"id": id,
	}, model.Material{
		ID:    m.ID,
		Title: request.Title,
		Slug:  slug.Make(request.Title),
	})
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "edit"),
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.EditSubmission(editM.ID, model.Submission{
		Content: request.Content,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "edit"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "edited"),
		Data: http.SubmissionHTTP{
			ID:        res.ID,
			Title:     editM.Title,
			ChapterID: &editM.ChapterID,
			Slug:      editM.Slug,
			Content:   res.Content,
			Date:      res.UpdatedAt,
		},
	})
}

func (h *Handlers) DeleteSubmission(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	f, err := h.CourseRepository.FindMaterial(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "delete"),
			Data:    nil,
		})
	}

	res, err := h.CourseRepository.DeleteSubmission(f.ID)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "delete"),
			Data:    nil,
		})
	}

	_, err = h.CourseRepository.DeleteMaterial(map[string]interface{}{
		"id": f.ID,
	})

	if err != nil {
		logrus.Warnln("[handler] Failed to delete Submission")
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "delete"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "deleted"),
		Data: http.SubmissionHTTP{
			ID:        f.ID,
			Title:     f.Title,
			ChapterID: &f.ChapterID,
			Content:   res.Content,
			Slug:      f.Slug,
			Date:      f.UpdatedAt,
		},
	})

}

func (h *Handlers) GetDetailSubmission(c *fiber.Ctx) error {
	id := c.Params("id")

	submission, err := h.CourseRepository.FindSubmissionStudent(map[string]interface{}{
		"id": id,
	}, false)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	findChapter, err := h.CourseRepository.FindChapter(map[string]interface{}{
		"id": submission.Material.ChapterID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	submissionStudent := http.SubmissionStudentDetailStudent{
		ID:              submission.ID,
		ActiveStudentID: submission.ActiveStudentID,
		Material:        submission.Material.Title,
		Grade:           submission.Grade,
		FileURL:         submission.FileUrl,
		Status:          string(submission.Status),
		Date:            submission.CreatedAt.String(),
		Comment:         submission.Comment,
		Description:     submission.Description,
		Chapter:         findChapter.Title,
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "retrieve"),
		Data:    submissionStudent,
	})

}

// find Submissions detail pov teacher
func (h *Handlers) GetDetailSubmissionTeacher(c *fiber.Ctx) error {
	id := c.Params("id")

	submission, err := h.CourseRepository.FindSubmissionStudent(map[string]interface{}{
		"id": id,
	}, false)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	findChapter, err := h.CourseRepository.FindChapter(map[string]interface{}{
		"id": submission.Material.ChapterID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	submissionDetail := http.SubmissionStudentDetailTeacher{
		ID:              submission.ID,
		Material:        submission.Material.Title,
		Grade:           submission.Grade,
		FileURL:         submission.FileUrl,
		Chapter:         findChapter.Title,
		Status:          string(submission.Status),
		Date:            submission.CreatedAt.String(),
		Description:     submission.Description,
		ActiveStudentID: submission.ActiveStudentID,
		Student:         submission.ActiveStudent.Student.Name,
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "retrieve"),
		Data:    submissionDetail,
	})
}

func (h *Handlers) FindStudentSubmission(c *fiber.Ctx) error {
	student_id := c.Params("student_id")

	if student_id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("student_id"),
			Data:    nil,
		})
	}

	result, err := h.CourseRepository.FindSubmissionStudent(map[string]interface{}{
		"student_id": student_id,
	}, false)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission student", "retrieve"),
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id":          result.ID,
		"material_id": result.MaterialID,
		"file_url":    result.FileUrl,
		"status":      result.Status,
		"description": result.Description,
		"student_id":  result.ActiveStudent.StudentID,
		"created_at":  result.CreatedAt,
		"updated_at":  result.UpdatedAt,
		"deleted_at":  result.DeletedAt,
	}

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "retrieve"),
		Data:    response,
	})

}

func (h *Handlers) Approvesubmisson(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	teacherID := c.Locals("user_id").(string)
	student_id := c.Params("student_id")

	if role != "TEACHER" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	var request http.SubmissionStudentAprrove

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	// find active stuednt
	findActiveStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": student_id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("active student", "retrieve"),
			Data:    nil,
		})
	}

	// find teacher
	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"user_id": teacherID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("teacher", "retrieve"),
			Data:    nil,
		})
	}

	result, err := h.CourseRepository.ApproveSubmission(map[string]interface{}{
		"id":                request.SubmissionID,
		"active_student_id": findActiveStudent.ID,
	}, request.Grade, findTeacher.ID)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission student", "approve"),
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id":          string(result.ID),
		"material_id": result.MaterialID,
		"file_url":    result.FileUrl,
		"description": result.Description,
		"student_id":  result.ActiveStudent.StudentID,
		"grade":       result.Grade,
		"comment":     result.Comment,
		"status":      result.Status,
		"created_at":  result.CreatedAt,
		"updated_at":  result.UpdatedAt,
		"deleted_at":  result.DeletedAt,
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission student", "approved"),
		Data:    response,
	})
}

func (h *Handlers) Rejectsubmisson(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	teacherID := c.Locals("user_id").(string)
	student_id := c.Params("student_id")

	if role != "TEACHER" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	var request http.SubmissionStudentReject

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	findActiveStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": student_id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("active student", "retrieve"),
			Data:    nil,
		})
	}

	// find teacher
	findTeacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"user_id": teacherID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("teacher", "retrieve"),
			Data:    nil,
		})
	}

	result, err := h.CourseRepository.RejectionsSubmission(map[string]interface{}{
		"id":                request.SubmissionID,
		"active_student_id": findActiveStudent.ID,
	}, &request.Comment, findTeacher.ID)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission student", "reject"),
			Data:    nil,
		})
	}
	response := map[string]interface{}{
		"id":          string(result.ID),
		"material_id": result.MaterialID,
		"file_url":    result.FileUrl,
		"description": result.Description,
		"student_id":  result.ActiveStudent.StudentID,
		"grade":       result.Grade,
		"comment":     result.Comment,
		"status":      result.Status,
		"created_at":  result.CreatedAt,
		"updated_at":  result.UpdatedAt,
		"deleted_at":  result.DeletedAt,
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission student", "reject"),
		Data:    response,
	})

}

// Find Class at active student by user id

func (h *Handlers) FindClassStudent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var response []http.StudentClassHTTP

	student, err := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": userID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("student", "retrieve"),
			Data:    nil,
		})
	}

	// find school by student id
	findSchool, err := h.SchoolRepository.FindSchool(map[string]interface{}{
		"id": student.SchoolsID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("school", "retrieve"),
			Data:    nil,
		})
	}

	activeStudent, err := h.UserRepository.FindActiveStudents(map[string]interface{}{
		"student_id":  student.ID,
		"school_year": findSchool.SchoolYear,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("active student", "retrieve"),
			Data:    nil,
		})
	}

	for _, item := range activeStudent {
		response = append(response, http.StudentClassHTTP{
			Classes: fmt.Sprintf("Kelas %s (%s)", item.Class, item.SchoolYear),
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("active student", "retrieve"),
		Data:    response,
	})

}

func (h *Handlers) UpdateProgress(c *fiber.Ctx) error {
	var request http.EnrollCourseHTTP
	user_id := c.Locals("user_id").(string)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	// find Active student
	findActiveStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": user_id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("active student", "retrieve"),
			Data:    nil,
		})
	}

	// find course
	findCourse, err := h.CourseRepository.FindCourse(map[string]interface{}{
		"id": request.CourseID,
	}, false, "")

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("course", "retrieve"),
			Data:    nil,
		})
	}

	id, _ := helper.GenerateNanoId()

	// Check if existing enrolled course
	_, err = h.CourseRepository.FindEnrollCourse(map[string]interface{}{
		"active_student_id": findActiveStudent.ID,
		"course_id":         findCourse.ID,
		"material_id":       request.MaterialID,
	})

	if err == nil {
		logrus.Infoln("[handler] Enrolls:")
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Student already finished this material",
			Data:    nil,
		})
	}

	_, err = h.CourseRepository.EnrollCourse(model.ActiveStudentCourse{
		ID:              id,
		ActiveStudentID: findActiveStudent.ID,
		CourseID:        findCourse.ID,
		MaterialID:      request.MaterialID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("active student course", "enroll"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("active student course", "enroll"),
		Data:    nil,
	})

}

func (h *Handlers) GetEnrollCourse(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	slugCourse := c.Params("slug")

	var response []http.EnrollCourseResponseHTTP

	// find active student by userId
	findAtiveStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": userID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to find active student",
			Data:    nil,
		})
	}

	// Find Student
	findStudent, err := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": userID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to find student",
			Data:    nil,
		})
	}

	// Find Course
	findCourse, err := h.CourseRepository.FindCourse(map[string]interface{}{
		"slug": slugCourse,
	}, false, "")

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("course", "retrieve"),
			Data:    nil,
		})
	}

	// find enroll course by active student id
	findEnrollCourse, err := h.CourseRepository.FindEnrollCourse(map[string]interface{}{
		"active_student_id": findAtiveStudent.ID,
		"course_id":         findCourse.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to find enroll course",
			Data:    nil,
		})
	}

	for _, item := range *findEnrollCourse {
		response = append(response, http.EnrollCourseResponseHTTP{
			ID:         item.ID,
			Student:    findStudent.Name,
			Class:      findAtiveStudent.Class,
			SchoolYear: findAtiveStudent.SchoolYear,
			Material:   item.Material.Title,
			Course:     item.Course.Title,
		})
	}

	return c.Status(500).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("enroll course", "retrieve"),
		Data:    response,
	})

}

func (h *Handlers) EditTheory(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	var request http.TheoryHTTP

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	editM, err := h.CourseRepository.UpdateMaterial(map[string]interface{}{
		"id": id,
	}, model.Material{
		Title: request.Title,
		Slug:  slug.Make(request.Title),
	})

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "edit"),
			Data:    nil,
		})
	}

	editTheory, err := h.CourseRepository.EditTheory(editM.ID, model.Theory{
		Content: request.Content,
	})

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "edit"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("theory", "edited"),
		Data: http.TheoryHTTP{
			ID:        editM.ID,
			Title:     editM.Title,
			ChapterID: editM.ChapterID,
			Slug:      editM.Slug,
			Content:   editTheory.Content,
			CreatedAt: editM.CreatedAt,
			UpdatedAt: editM.UpdatedAt,
		},
	})

}

func (h *Handlers) DeleteTheory(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	m, err := h.CourseRepository.FindMaterial(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "find"),
			Data:    nil,
		})
	}

	t, err := h.CourseRepository.DeleteTheory(map[string]interface{}{
		"material_id": m.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "delete"),
			Data:    nil,
		})
	}

	_, err = h.CourseRepository.DeleteMaterial(map[string]interface{}{
		"id": m.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("theory", "delete"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("theory", "delete"),
		Data: http.TheoryHTTP{
			ID:        m.ID,
			ChapterID: m.ChapterID,
			Title:     m.Title,
			Slug:      m.Slug,
			Content:   t.Content,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
	})
}

func (h *Handlers) SaveCompleteCourse(c *fiber.Ctx) error {
	var request map[string]string

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	role := c.Locals("role")
	userID := c.Locals("user_id")

	if role != "STUDENT" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Unauthorized",
			Data:    nil,
		})
	}

	// Find Student
	student, err := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": userID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("complete course", "save"),
			Data:    nil,
		})
	}

	// Find Active Student
	activeStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": student.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("complete course", "save"),
			Data:    nil,
		})
	}

	// Insert Complete Course
	id, _ := gonanoid.New(20)
	res, err := h.CourseRepository.SaveCompleteCourse(model.CompleteCourse{
		ID:              id,
		CourseID:        request["course_id"],
		ActiveStudentID: activeStudent.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("complete course", "save"),
			Data:    nil,
		})
	}

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("complete course", "saved"),
		Data: map[string]string{
			"course_id": res.CourseID,
		},
	})

}

func (h *Handlers) ListingSubmission(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	user_id := c.Locals("user_id")

	logrus.Infoln("[handler] Listing Submission within POV:", role)
	if role == "TEACHER" {
		teacher, err := h.UserRepository.FindTeacher(map[string]interface{}{
			"user_id": user_id.(string),
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("submission", "retrieve"),
				Data:    nil,
			})
		}

		materialID := c.Query("material_id")
		class := c.Query("class")

		submission, err := h.CourseRepository.FindSubmissionPreload(teacher.ID, "", "", "", materialID, class)
		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("submission", "retrieve"),
				Data:    nil,
			})
		}

		countStatus := helper.CountStatusSubmission(submission)

		var resp []http.SubmissionStudentHTTP

		for _, el := range submission {
			student, err := h.UserRepository.FindStudent(map[string]interface{}{
				"id": el.ActiveStudent.StudentID,
			})

			if err != nil {
				logrus.Warnln("[handler] Couldn't find active student")
			}
			comment := ""

			if el.Comment != nil {
				comment = *el.Comment
			}

			status := string(el.Status)

			if status == "REV_REJECT" {
				status = "REJECTED"
			}

			material, err := h.CourseRepository.FindMaterial(map[string]interface{}{
				"id": el.MaterialID,
			})

			if material == nil {
				logrus.Warnln("[handler] Data material is nil!")
			}

			course, err := h.CourseRepository.FindCourse(
				map[string]interface{}{
					"id": material.CourseID,
				}, false, "")

			if err != nil || course == nil {
				logrus.Warnln("[handler] Data is nil!")
			}

			chapter, err := h.CourseRepository.FindChapter(map[string]interface{}{
				"id": material.ChapterID,
			})

			if err != nil || chapter == nil {
				logrus.Warnln("[handler] Data is nil")
			}

			resp = append(resp, http.SubmissionStudentHTTP{
				ID:              el.ID,
				SubmissionID:    el.MaterialID,
				StudentName:     student.Name,
				StudentID:       el.ActiveStudent.StudentID,
				Class:           el.ActiveStudent.Class,
				CourseTitle:     course.Title,
				ChapterTitle:    chapter.Title,
				Description:     el.Description,
				Comment:         comment,
				File:            el.FileUrl,
				SubmissionTitle: el.Material.Title,
				Status:          status,
				Date:            el.UpdatedAt,
			})
		}

		return c.Status(200).JSON(&http.WebResponse{
			Status:  "success",
			Message: h.successResponse("submission", "retrieve"),
			Data: http.ListSubmissionResponseHTTP{
				Status: http.StatusSubmissionHTTP{
					Pending:  countStatus.Pending,
					Rejected: countStatus.Rejected,
					Approved: countStatus.Approved,
				},
				Submissions: resp,
			},
		})

	}

	// Find Current StudentIDt
	student, err := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": user_id.(string),
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	activeStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": student.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	filterTeacher := c.Query("teacher_id")
	filterCourseID := c.Query("course_id")
	submission, err := h.CourseRepository.FindSubmissionPreload("", activeStudent.ID, filterTeacher, filterCourseID, "", "")

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	countStatus := helper.CountStatusSubmission(submission)

	var resp []http.SubmissionStudentPOVUserHTTP

	for _, el := range submission {
		comment := ""

		if el.Comment != nil {
			comment = *el.Comment
		}

		status := string(el.Status)

		if status == "REV_REJECT" {
			status = "REJECTED"
		}

		material, _ := h.CourseRepository.FindMaterial(map[string]interface{}{
			"id": el.MaterialID,
		})

		if material == nil {
			logrus.Warnln("[handler] Data material is nil!")
		}

		course, _ := h.CourseRepository.FindCourse(
			map[string]interface{}{
				"id": material.CourseID,
			}, false, "")

		if course == nil {
			logrus.Warnln("[handler] Data is nil!")
		}

		chapter, _ := h.CourseRepository.FindChapter(map[string]interface{}{
			"id": material.ChapterID,
		})

		if chapter == nil {
			logrus.Warnln("[handler] Data is nil")
		}

		resp = append(resp, http.SubmissionStudentPOVUserHTTP{
			ID:              el.ID,
			SubmissionID:    el.MaterialID,
			CourseID:        course.ID,
			ChapterID:       chapter.ID,
			CourseTitle:     course.Title,
			ChapterTitle:    chapter.Title,
			SubmissionTitle: el.Material.Title,
			Status:          status,
			Description:     el.Description,
			Comment:         comment,
			File:            el.FileUrl,
			Date:            el.UpdatedAt,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("submission", "retrieve"),
		Data: http.ListSubmissionStudentPOVUserHTTP{
			Status: http.StatusSubmissionHTTP{
				Pending:  countStatus.Pending,
				Rejected: countStatus.Rejected,
				Approved: countStatus.Approved,
			},
			Submissions: resp,
		},
	})

}

func (h *Handlers) ResetSubmission(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	if role != "STUDENT" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Unauthorized",
			Data:    nil,
		})
	}

	student, err := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": userID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	activeStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": student.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "retrieve"),
			Data:    nil,
		})
	}

	var request http.ResetSubmittedRequestHTTP
	err = c.BodyParser(&request)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	err = h.CourseRepository.ResetSubmittedSubmission(map[string]interface{}{
		"active_student_id": activeStudent.ID,
		"id":                request.SubmittedID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("submission", "reset"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Submisison has been revisioned!",
		Data: map[string]string{
			"id": request.SubmittedID,
		},
	})

}

func (h *Handlers) GetSubmissionPlaceholder(c *fiber.Ctx) error {
	res, err := h.CourseRepository.PlaceholderFilterSubmission()

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("placeholder", "retrieve"),
			Data:    nil,
		})
	}

	var response []http.TSubmissionPlaceholder

	// Find Material
	for _, el := range res {
		// Because above method will only 'material_id', we have to do our own custom query finder
		material, err := h.CourseRepository.FindMaterial(map[string]interface{}{
			"id": el.MaterialID,
		})

		if err != nil {
			logrus.Warnln("[handler] Material is not found!")
			logrus.Warnln("[handler] Error:", err)
		}

		course, err := h.CourseRepository.FindCourse(map[string]interface{}{
			"id": material.CourseID,
		}, false, "")

		if err != nil {
			logrus.Warnln("[handler] Course is not found!")
			logrus.Warnln("[handler] Error:", err)
		}

		// Do query again to find material within Course ID
		materials, err := h.CourseRepository.FindMaterials(map[string]interface{}{
			"course_id": course.ID,
			"type":      "SUBMISSION",
		})

		if err != nil {
			logrus.Warnln("[handler] Course is not found!")
			logrus.Warnln("[handler] Error:", err)
		}

		var materialResp []http.TMaterialPlaceholder
		if err == nil {
			for _, j := range materials {
				materialResp = append(materialResp, http.TMaterialPlaceholder{
					ID:    j.ID,
					Title: j.Title,
				})
			}
		}

		response = append(response, http.TSubmissionPlaceholder{
			Course: http.TCoursePlaceholder{
				ID:    course.ID,
				Title: course.Title,
			},
			ListMaterial: materialResp,
		})

	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("placeholder", "retrieve"),
		Data:    response,
	})

}
