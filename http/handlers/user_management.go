package handlers

import (
	"github.com/cvzamannow/E-Learning-API/helper"
	"github.com/cvzamannow/E-Learning-API/http"
	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// RouterUserManagemet defines routes for user management endpoints.
func (h *Handlers) RouterUserManagemet(app *fiber.App) {
	v1 := app.Group("/api/v1/admin")

	v1Base := app.Group("/api/v1")

	// CRUD operations for students
	v1.Post("/students", h.Middleware.Protected(), h.CreateStudent)
	v1.Get("/students", h.Middleware.Protected(), h.GetAllStudents)
	v1.Get("/students/:id", h.Middleware.Protected(), h.GetStudentById)
	v1.Put("/students/:id", h.Middleware.Protected(), h.UpdateStudentById)
	v1.Delete("/students/:id", h.Middleware.Protected(), h.DeleteStudentById)

	// CRUD operations for teachers
	v1.Post("/teachers", h.Middleware.Protected(), h.CreateTeacher)
	v1.Get("/teachers", h.Middleware.Protected(), h.GetAllTeachers)
	v1.Get("/teachers/:id", h.Middleware.Protected(), h.GetTeacherById)
	v1.Put("/teachers/:id", h.Middleware.Protected(), h.UpdateTeacherById)
	v1.Delete("/teachers/:id", h.Middleware.Protected(), h.DeleteTeacherById)

	v1Base.Get("/placeholder/teachers", h.Middleware.Protected(), h.GetAllTeachersHandler)

	// CRUD operations for active students
	v1.Post("/active-students", h.Middleware.Protected(), h.CreateActiveStudent)
	v1.Get("/active-students", h.Middleware.Protected(), h.GetAllActiveStudents)
	v1.Get("/active-students/:id", h.Middleware.Protected(), h.GetActiveStudentById)
	v1.Put("/active-students/:id", h.Middleware.Protected(), h.UpdateActiveStudentsById)
	v1.Delete("/active-students/:id", h.Middleware.Protected(), h.DeleteActiveStudentById)

	// CRUD operations for Admin School
	v1.Post("/super-admin/admin-schools", h.Middleware.Protected(), h.CreateAdminSchool)
	v1.Get("/super-admin/admin-schools", h.Middleware.Protected(), h.GetAllAdminSchool)
	v1.Get("/super-admin/admin-schools/:id", h.Middleware.Protected(), h.GetAdminSchoolById)
	v1.Put("/super-admin/admin-schools/:id", h.Middleware.Protected(), h.UpdateAdminSchool)
	v1.Delete("/super-admin/admin-schools/:id", h.Middleware.Protected(), h.DeleteAdminSchool)
}

// CreateStudent handles the creation of a new student.
func (h *Handlers) CreateStudent(c *fiber.Ctx) error {
	var createStudentRequest http.Student

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	// Parse request body into createStudentRequest struct
	if err := c.BodyParser(&createStudentRequest); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request to create student",
			Data:    nil,
		})
	}

	// Generate a nano ID string for student and user ID
	id, err := helper.GenerateNanoId()
	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating nano ID",
			Data:    nil,
		})
	}

	// Create a hash for the password
	hashString, err := bcrypt.GenerateFromPassword(
		[]byte(createStudentRequest.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating password hash",
			Data:    nil,
		})
	}

	// Create student record in the database
	result, err := h.UserRepository.CreateStudent(model.Student{
		ID:        id,
		Name:      createStudentRequest.Name,
		IdNumber:  createStudentRequest.IDNumber,
		SchoolsID: createStudentRequest.SchoolsID,
		User: model.User{
			ID:       id,
			Username: createStudentRequest.Username,
			Password: string(hashString),
			Status:   "ACTIVE",
			Role:     "STUDENT",
		},
	})

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Student already exists",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Student created successfully",
		Data:    result,
	})
}

// GetAllStudents handles the retrieval of all students.
func (h *Handlers) GetAllStudents(c *fiber.Ctx) error {
	result, err := h.UserRepository.GetAllStudent()

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error getting students",
			Data:    nil,
		})
	}

	var response []map[string]interface{}

	for _, item := range *result {
		itemMap := map[string]interface{}{
			"id":        item.ID,
			"name":      item.Name,
			"id_number": item.IdNumber,
			"user": map[string]interface{}{
				"id":       item.User.ID,
				"username": item.User.Username,
				"status":   item.User.Status,
				"role":     item.User.Role,
			},
		}

		response = append(response, itemMap)
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Students retrieved successfully",
		Data:    response,
	})
}

// GetStudentById handles the retrieval of a student by ID.
func (h *Handlers) GetStudentById(c *fiber.Ctx) error {
	result, err := h.UserRepository.FindStudent(map[string]interface{}{
		"id": c.Params("id"),
	})

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Student not found",
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id":        result.ID,
		"name":      result.Name,
		"id_number": result.IdNumber,
		"school_id": result.Schools.ID,
		"user": map[string]interface{}{
			"id":       result.User.ID,
			"username": result.User.Username,
			"status":   result.User.Status,
			"role":     result.User.Role,
		},
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Student retrieved successfully",
		Data:    response,
	})
}

// UpdateStudentById handles the update of a student by ID.
func (h *Handlers) UpdateStudentById(c *fiber.Ctx) error {
	var requestUpdate http.Student

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	// Parse request body into requestUpdate struct
	if err := c.BodyParser(&requestUpdate); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request to update student",
			Data:    nil,
		})
	}

	hashString, err := bcrypt.GenerateFromPassword([]byte(requestUpdate.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "generate password hash"),
			Data:    nil,
		})
	}

	// Update student record in the database
	_, err = h.UserRepository.UpdateStudent(
		map[string]interface{}{
			"id": c.Params("id"),
		},
		model.Student{
			Name:      requestUpdate.Name,
			IdNumber:  requestUpdate.IDNumber,
			SchoolsID: requestUpdate.SchoolsID,
			User: model.User{
				Username: requestUpdate.Username,
				Password: string(hashString),
			},
		},
	)

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error updating student",
			Data:    err.Error(),
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Student updated successfully",
		Data:    nil,
	})
}

// DeleteStudentById handles the deletion of a student by ID.
func (h *Handlers) DeleteStudentById(c *fiber.Ctx) error {
	// Delete student record from the database
	_, err := h.UserRepository.DeleteStudent(map[string]interface{}{
		"id": c.Params("id"),
	})

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error deleting student",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Student deleted successfully",
		Data:    nil,
	})
}

// CreateTeacher handles the creation of a new teacher.
func (h *Handlers) CreateTeacher(c *fiber.Ctx) error {
	var requestBody http.Teacher

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	// Parse request body into requestBody struct
	if err := c.BodyParser(&requestBody); err != nil {
		logrus.Warnln("[handlers] Error parsing payload:", err)
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request to create teacher",
			Data:    nil,
		})
	}

	// Generate a nano ID string for teacher and user ID
	id, err := helper.GenerateNanoId()
	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating nano ID",
			Data:    nil,
		})
	}

	// Generate a nano ID string for teacher ID
	teacherID, err := helper.GenerateNanoId()
	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating teacher ID",
			Data:    nil,
		})
	}

	// Create a hash for the password
	hashString, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating password hash",
			Data:    nil,
		})
	}

	// Create teacher record in the database
	result, err := h.UserRepository.CreateTeacher(model.Teacher{
		ID:        teacherID,
		Name:      requestBody.Name,
		IdNumber:  requestBody.IDNumber,
		SchoolsID: requestBody.SchoolID,
		User: model.User{
			ID:       id,
			Username: requestBody.Username,
			Password: string(hashString),
			Status:   "ACTIVE",
			Role:     "TEACHER",
		},
	})

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Teacher already exists",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Teacher created successfully",
		Data:    result,
	})
}

// GetAllTeachers handles the retrieval of all teachers.
func (h *Handlers) GetAllTeachers(c *fiber.Ctx) error {
	result, err := h.UserRepository.GetAllTeacher()

	role := c.Locals("role").(string)

	if role == "STUDENT" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error getting teachers",
			Data:    nil,
		})
	}

	var response []map[string]interface{}

	for _, item := range *result {
		itemMap := map[string]interface{}{
			"id":        item.ID,
			"id_number": item.IdNumber,
			"name":      item.Name,
			"user": map[string]interface{}{
				"username": item.User.Username,
				"status":   item.User.Status,
				"role":     item.User.Role,
			},
		}
		response = append(response, itemMap)

	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Teachers retrieved successfully",
		Data:    response,
	})
}

// GetTeacherById handles the retrieval of a teacher by ID.
func (h *Handlers) GetTeacherById(c *fiber.Ctx) error {
	result, err := h.UserRepository.FindTeacher(map[string]interface{}{
		"id": c.Params("id"),
	})

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Teacher not found",
			Data:    err.Error(),
		})
	}

	response := map[string]interface{}{
		"id":        result.ID,
		"id_number": result.IdNumber,
		"name":      result.Name,
		"user": map[string]interface{}{
			"username": result.User.Username,
			"status":   result.User.Status,
			"role":     result.User.Role,
		},
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Teacher retrieved successfully",
		Data:    response,
	})
}

// UpdateTeacherById handles the update of a teacher by ID.
func (h *Handlers) UpdateTeacherById(c *fiber.Ctx) error {
	var requestBody http.Teacher
	if err := c.BodyParser(&requestBody); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request to update teacher",
			Data:    nil,
		})
	}

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	hashString, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "generate password hash"),
			Data:    nil,
		})
	}

	_, err = h.UserRepository.UpdateTeacher(map[string]interface{}{
		"id": c.Params("id"),
	}, model.Teacher{
		Name:     requestBody.Name,
		IdNumber: requestBody.IDNumber,
		User: model.User{
			Username: requestBody.Username,
			Password: string(hashString),
		},
	})

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error updating teacher",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Teacher updated successfully",
		Data:    nil,
	})

}

// DeleteTeacherById handles the deletion of a teacher by ID.
func (h *Handlers) DeleteTeacherById(c *fiber.Ctx) error {
	_, err := h.UserRepository.DeleteTeacher(map[string]interface{}{
		"id": c.Params("id"),
	})

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error deleting teacher",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Teacher deleted successfully",
		Data:    nil,
	})
}

// TODO: implement repository for active student

// CreateActiveStudent handles the creation of a new active student.
func (h *Handlers) CreateActiveStudent(c *fiber.Ctx) error {
	var requestActiveStudent http.ActiveStudent

	// Parse request body into requestActiveStudent struct
	if err := c.BodyParser(&requestActiveStudent); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request to create active student",
			Data:    nil,
		})
	}

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	// Generate a nano ID string for active student ID
	id, err := helper.GenerateNanoId()
	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating nano ID",
			Data:    nil,
		})
	}

	// Create active student record in the database
	_, err = h.UserRepository.CreateActiveStudent(model.ActiveStudent{
		ID:         id,
		StudentID:  requestActiveStudent.StudentID,
		SchoolYear: requestActiveStudent.SchoolYear,
		ClassSlug:  slug.Make(requestActiveStudent.Class),
		Class:      requestActiveStudent.Class,
	})

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error creating active student",
			Data:    err.Error(),
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Active student created successfully",
		Data:    nil,
	})
}

// GetActiveStudentById handles the retrieval of an active student by ID.
func (h *Handlers) GetActiveStudentById(c *fiber.Ctx) error {
	result, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"id": c.Params("id"),
	})

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Active student not found",
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id":          result.ID,
		"school_year": result.SchoolYear,
		"class":       result.Class,
		"student": map[string]interface{}{
			"id":        result.Student.ID,
			"name":      result.Student.Name,
			"id_number": result.Student.IdNumber,
		},
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Active student retrieved successfully",
		Data:    response,
	})
}

// GetAllActiveStudents handles the retrieval of all active students.
func (h *Handlers) GetAllActiveStudents(c *fiber.Ctx) error {
	result, err := h.UserRepository.GetAllActiveStudent()

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error getting active students",
			Data:    nil,
		})
	}

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	var response []map[string]interface{}

	for _, items := range *result {
		response = append(response, map[string]interface{}{
			"id":          items.ID,
			"school_year": items.SchoolYear,
			"class":       items.Class,
			"student": map[string]interface{}{
				"id":        items.Student.ID,
				"name":      items.Student.Name,
				"id_number": items.Student.IdNumber,
				"user_id":   items.Student.UserID,
			},
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Active students retrieved successfully",
		Data:    response,
	})
}

// DeleteActiveStudentById handles the deletion of an active student by ID.
func (h *Handlers) DeleteActiveStudentById(c *fiber.Ctx) error {
	// Delete active student record from the database
	_, err := h.UserRepository.DeleteActiveStudent(map[string]interface{}{
		"id": c.Params("id"),
	})

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error deleting active student",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Active student deleted successfully",
		Data:    nil,
	})
}

// UpdateActiveStudentsById handles the update of an active student by ID.
func (h *Handlers) UpdateActiveStudentsById(c *fiber.Ctx) error {
	var requestBody http.ActiveStudent
	if err := c.BodyParser(&requestBody); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request to update active student",
			Data:    nil,
		})
	}

	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	// Update active student record in the database
	_, err := h.UserRepository.UpdateActiveStudent(map[string]interface{}{
		"id": c.Params("id"),
	}, model.ActiveStudent{
		StudentID:  requestBody.StudentID,
		SchoolYear: requestBody.SchoolYear,
		Class:      requestBody.Class,
		ClassSlug:  slug.Make(requestBody.Class),
	})

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error updating active student",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Active student updated successfully",
		Data:    nil,
	})
}

// CRUD ADMIN
func (h *Handlers) CreateAdminSchool(c *fiber.Ctx) error {

	role := c.Locals("role").(string)

	if role != "SUPER_ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	var request http.Admin

	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	id, err := helper.GenerateNanoId()
	userId, _ := helper.GenerateNanoId()

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "generate nano ID"),
			Data:    nil,
		})
	}

	// Create a hash for the password
	hashString, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "generate password hash"),
			Data:    nil,
		})
	}

	result, err := h.UserRepository.CreateAdmin(model.AdminSchool{
		ID:       string(id),
		SchoolID: request.SchoolID,
		User: model.User{
			ID:       string(userId),
			Username: request.Username,
			Password: string(hashString),
			Status:   "ACTIVE",
			Role:     "ADMIN",
		},
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error creating admin / Username alerdy exist",
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id":        id,
		"user_id":   userId,
		"school_id": request.SchoolID,
		"user": map[string]interface{}{
			"id":       userId,
			"username": result.User.Username,
			"password": hashString,
			"status":   result.User.Status,
			"role":     result.User.Role,
		},
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Admin", "created"),
		Data:    response,
	})
}

func (h *Handlers) GetAllAdminSchool(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "SUPER_ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	result, err := h.UserRepository.GetAllAdminSchool()

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admins", "retrieve"),
			Data:    nil,
		})
	}

	var response []map[string]interface{}
	for _, item := range *result {
		response = append(response, map[string]interface{}{
			"id": item.ID,
			"school": map[string]interface{}{
				"school_year": item.Schools.SchoolYear,
				"name":        item.Schools.Name,
				"address":     item.Schools.Address,
				"logo":        item.Schools.Logo,
			},
			"user": map[string]interface{}{
				"id":       item.User.ID,
				"username": item.User.Username,
				"status":   item.User.Status,
				"role":     item.User.Role,
			},
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Admins", "retrieved"),
		Data:    response,
	})
}

func (h *Handlers) GetAdminSchoolById(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "SUPER_ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	result, err := h.UserRepository.FindAdminSchool(map[string]interface{}{
		"id": c.Params("id"),
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "retrieve"),
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id": result.ID,
		"school": map[string]interface{}{
			"SchoolYear": result.Schools.SchoolYear,
			"Name":       result.Schools.Name,
			"Address":    result.Schools.Address,
			"Logo":       result.Schools.Logo,
		},
		"user": map[string]interface{}{
			"id":       result.User.ID,
			"username": result.User.Username,
			"status":   result.User.Status,
			"role":     result.User.Role,
		},
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Admin", "retrieved"),
		Data:    response,
	})
}

func (h *Handlers) UpdateAdminSchool(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "SUPER_ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	var request http.Admin

	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	hashString, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "generate password hash"),
			Data:    nil,
		})
	}

	// Panggil fungsi UpdateAdminSchool dengan menggunakan transaksi
	_, err = h.UserRepository.UpdateAdminSchool(map[string]interface{}{
		"id": c.Params("id"),
	},
		model.AdminSchool{
			SchoolID: request.SchoolID,
			User: model.User{
				Username: request.Username,
				Password: string(hashString),
			},
		},
	)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "update"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Admin", "updated"),
		Data:    nil,
	})
}

func (h *Handlers) DeleteAdminSchool(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "SUPER_ADMIN" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	_, err := h.UserRepository.DeleteAdminSchool(map[string]interface{}{
		"id": c.Params("id"),
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Admin", "delete"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Admin", "deleted"),
		Data:    nil,
	})
}

// get all teacher

func (h *Handlers) GetAllTeachersHandler(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(string)

	// mencari siswa
	findStundet, err := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": userID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error getting student",
			Data:    nil,
		})
	}

	result, err := h.UserRepository.FindTeachers(map[string]interface{}{
		"schools_id": findStundet.SchoolsID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error getting teachers",
			Data:    nil,
		})
	}

	var response []http.TeacherHTTP

	for _, item := range *result {
		response = append(response, http.TeacherHTTP{
			ID:       item.ID,
			Name:     item.Name,
			IDNumber: item.IdNumber,
			Username: item.User.Username,
			SchoolID: &findStundet.SchoolsID,
		})

	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Teachers retrieved successfully",
		Data:    response,
	})
}
