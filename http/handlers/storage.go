package handlers

import (
	"fmt"
	"os"

	"github.com/cvzamannow/E-Learning-API/http"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *Handlers) RouteStorage(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/storage", h.Middleware.Protected(), h.UploadViaHTTP)
}

func (h *Handlers) UploadViaHTTP(c *fiber.Ctx) error {
	os.Mkdir("temp", os.ModePerm)

	file, err := c.FormFile("file")

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	filename := fmt.Sprintf("./temp/%s", file.Filename)
	err = c.SaveFile(file, filename)

	if err != nil {
		logrus.Warnln("[storage-handlers] An error occured", err)
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("storage", "saving file"),
			Data:    nil,
		})
	}

	object, err := h.R2Cloudflare.Upload(filename)

	if err != nil {
		logrus.Warnln("[storage-handlers] An error occured", err)
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("storage", "saving file"),
			Data:    nil,
		})
	}

	return c.Status(201).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("file", "uploaded"),
		Data: map[string]interface{}{
			"url": object,
		},
	})
}
