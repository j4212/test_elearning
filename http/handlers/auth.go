package handlers

import (
	"time"

	"github.com/cvzamannow/E-Learning-API/http"
	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) RouteAuth(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Post("/login", h.LoginHandler)
	v1.Get("/verify", h.Middleware.Protected(), h.Verify)
	v1.Post("/register", h.RegisterHandler)
}

func (h *Handlers) RegisterHandler(c *fiber.Ctx) error {
	var request http.Register

	err := c.BodyParser(&request)

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	id, _ := gonanoid.New(20)

	// Hash Password
	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	res, err := h.UserRepository.CreateUser(model.User{
		ID:       id,
		Username: request.Username,
		Password: string(hash),
		Role:     request.Role,
		Status:   request.Status,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("user", "register"),
			Data:    nil,
		})
	}

	if request.Role == "TEACHER" {
		id, _ := gonanoid.New(20)
		_, err := h.UserRepository.CreateTeacher(model.Teacher{
			ID:     id,
			UserID: res.ID,
			Name:   request.Name,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("teacher", "create"),
				Data:    nil,
			})
		}
	}

	return c.Status(201).JSON(http.WebResponse{
		Status:  "success",
		Message: h.successResponse("user", "registered"),
		Data: map[string]interface{}{
			"username": res.Username,
		},
	})
}

func (h *Handlers) LoginHandler(c *fiber.Ctx) error {

	var requestLogin http.Login

	if err := c.BodyParser(&requestLogin); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request login",
			Data:    nil,
		})
	}

	response, err := h.UserRepository.FindUser(map[string]interface{}{
		"username": requestLogin.Username,
	})

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "User is not registered!",
			Data:    nil,
		})
	}

	// search shcool_name
	resultstudent, _ := h.UserRepository.FindStudent(map[string]interface{}{
		"user_id": response.ID,
	})

	var claims *jwt.Token
	if resultstudent != nil {
		// mengecek apakah student id nya sudah terdaftar di active student atau belum
		activeStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
			"student_id": resultstudent.ID,
		})

		if err != nil {
			return c.Status(404).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Active student not found",
				Data:    nil,
			})
		}

		if activeStudent == nil {
			return c.Status(400).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Active student not found",
				Data:    nil,
			})
		}

		// mencari sekolah
		school, err := h.SchoolRepository.FindSchool(map[string]interface{}{
			"id": resultstudent.SchoolsID,
		})

		if err != nil {
			return c.Status(404).JSON(&http.WebResponse{
				Status:  "error",
				Message: "School is not found",
				Data:    nil,
			})
		}

		claims = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username":    response.Username,
			"name": resultstudent.Name,
			"user_id":     response.ID,
			"role":        response.Role,
			"school_name": school.Name,
			"exp":         time.Now().Add(time.Hour * 24).Unix(),
		})

	} else {
		// Find Teacher Name 
		teacher := ""

		teacherQuery, _ := h.UserRepository.FindTeacher(map[string]interface{}{
			"user_id": response.ID,
		})

		if teacherQuery != nil {
			teacher = teacherQuery.Name
		}

		claims = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": response.Username,
			"name": teacher,
			"user_id":  response.ID,
			"role":     response.Role,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(response.Password), []byte(requestLogin.Password))

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Password is incorrect!",
			Data:    nil,
		})
	}

	// create jwt token from the request username

	token, err := claims.SignedString(h.JWT_SECRET)

	if err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generate jwt token",
			Data:    nil,
		})
	}

	return c.JSON(&http.WebResponse{
		Status:  "success",
		Message: "Credentials is valid and Login successfuly!",
		Data: map[string]interface{}{
			"token": token,
			"role":  response.Role,
		},
	})

}

func (h *Handlers) Verify(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims == nil {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Unauthorized",
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "User verified",
		Data:    claims,
	})

}
