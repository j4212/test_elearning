package handlers

import (
	"net/http"

	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Route setup for certificates
func (h *Handlers) RouteCertificates(app *fiber.App) {
	v1 := app.Group("/test")
	v1.Post("/certificates", h.CreateCertificateHandler)
	// Add other certificate routes if needed
}

// Handler for creating a certificate
func (h *Handlers) CreateCertificateHandler(c *fiber.Ctx) error {
	var request model.Certificate
	if err := c.BodyParser(&request); err != nil {
		logrus.Error("[HANDLER] Failed to parse request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	// Call the CreateCertificate function from the repository
	certificate, err := h.CertificateRepo.CreateCertificate(request)
	if err != nil {
		logrus.Error("[HANDLER] Failed to create certificate:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Certificate created successfully",
		"data":    certificate,
	})
}
