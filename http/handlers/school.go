package handlers

import (
	"github.com/cvzamannow/E-Learning-API/helper"
	"github.com/cvzamannow/E-Learning-API/http"
	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) RouterSchool(c *fiber.App) {
	v1 := c.Group("/api/v1/super-admin")
	v1.Post("/schools", h.Middleware.Protected(), h.CreateSchoolHandler)
	v1.Get("/schools", h.Middleware.Protected(), h.GetAllSchools)
	v1.Get("/schools/:id", h.Middleware.Protected(), h.GetSchoolById)
	v1.Put("/schools/:id", h.Middleware.Protected(), h.UpdateSchoolHandler)
	v1.Delete("/schools/:id", h.Middleware.Protected(), h.DeleteSchoolById)

	v1Base := c.Group("/api/v1")
	v1Base.Get("/classes", h.Middleware.Protected(), h.FindClasses)
}

func (h *Handlers) CreateSchoolHandler(c *fiber.Ctx) error {
	var request http.Schools
	if err := c.BodyParser(&request); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request create school",
			Data:    nil,
		})
	}

	uuid, err := helper.GenerateNanoId()

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating nano id",
			Data:    nil,
		})
	}

	result, err := h.SchoolRepository.CreateSchool(model.Schools{
		ID:         uuid,
		SchoolYear: request.SchoolYear,
		Name:       request.Name,
		Address:    request.Address,
		Logo:       request.Logo,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to register school because internal error",
			Data:    err.Error(),
		})
	}

	response := map[string]interface{}{
		"id":          result.ID,
		"school_year": result.SchoolYear,
		"name":        result.Name,
		"address":     result.Address,
		"logo":        result.Logo,
	}

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: "School has been created!",
		Data:    response,
	})

}

func (h *Handlers) GetSchoolById(c *fiber.Ctx) error {
	result, err := h.SchoolRepository.FindSchool(map[string]interface{}{
		"id": c.Params("id"),
	})

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "School id not found",
			Data:    nil,
		})
	}

	response := map[string]interface{}{
		"id":          result.ID,
		"school_year": result.SchoolYear,
		"name":        result.Name,
		"address":     result.Address,
		"logo":        result.Logo,
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "School found successfully",
		Data:    response,
	})

}

func (h *Handlers) GetAllSchools(c *fiber.Ctx) error {
	result, err := h.SchoolRepository.GetAllSchool()

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to get schools because internal error",
			Data:    nil,
		})
	}

	var response []map[string]interface{}

	for _, item := range *result {
		response = append(response, map[string]interface{}{
			"id":         item.ID,
			"schoolYear": item.SchoolYear,
			"name":       item.Name,
			"address":    item.Address,
			"logo":       item.Logo,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Schools found successfully",
		Data:    response,
	})
}

func (h *Handlers) UpdateSchoolHandler(c *fiber.Ctx) error {
	var request http.Schools
	if err := c.BodyParser(&request); err != nil {
		return c.JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error parsing request update school",
			Data:    nil,
		})
	}

	result, err := h.SchoolRepository.UpdateSchool(c.Params("id"), model.Schools{
		SchoolYear: request.SchoolYear,
		Name:       request.Name,
		Address:    request.SchoolYear,
		Logo:       request.Logo,
	})

	response := map[string]interface{}{
		"id":          result.ID,
		"school_year": result.SchoolYear,
		"name":        result.Name,
		"address":     result.Address,
		"logo":        result.Logo,
	}

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to update school because internal error",
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "School has been updated!",
		Data:    response,
	})

}

func (h *Handlers) DeleteSchoolById(c *fiber.Ctx) error {
	_, err := h.SchoolRepository.DeleteSchool(c.Params("id"))

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Failed to delete school because internal error",
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "School has been deleted!",
		Data:    nil,
	})
}

func (h *Handlers) FindClasses(c *fiber.Ctx) error {
	q := c.Query("q")

	res, err := h.SchoolRepository.FindClasses(q)

	if err != nil || len(res) == 0 || res == nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Classes is not found!",
			Data:    nil,
		})
	}

	var classResponse []map[string]interface{}

	for index, el := range res {
		classResponse = append(classResponse, map[string]interface{}{
			"id":    index + 1,
			"label": el.Class,
		})
	}

	return c.Status(200).JSON(map[string]interface{}{
		"status":  "success",
		"message": h.successResponse("class", "retrieve"),
		"classes": classResponse,
	})
}
