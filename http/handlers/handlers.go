package handlers

import (
	"fmt"

	"github.com/cvzamannow/E-Learning-API/middleware"
	"github.com/cvzamannow/E-Learning-API/repository"
	"github.com/cvzamannow/E-Learning-API/service"
)

type Handlers struct {
	CourseRepository repository.CourseRepository
	UserRepository   repository.UserRepository
	JWT_SECRET       []byte
	Middleware       middleware.Middleware
	SchoolRepository repository.SchoolRepository
	QuizRepository repository.QuizRepository
	R2Cloudflare     *service.R2Stub
	GradesRepository repository.GradesRepository
	CertificateRepo repository.CertificateRepository
}

// This is a reusable response message
// Method should be in private visibility (Not Pascal Case).

// Error when not specify path params
func (h *Handlers) errorSpecifyResource(params string) string {
	// Example: /api/v1/courses/:id
	// If path params /:id is not included, this method should be invoked
	// This method will produce 'Must specify resource with id params' because 'id' is the params that we passed in the method.
	return fmt.Sprintf("Must specify resource with '%s' params!", params)
}

// Error when there is internal error such as database query error.
// 'Resource' parameter refer to the path params in the route.
// 'action' refer to the action of the resource. Ex: Create, Edit, Retrieve, Delete
func (h *Handlers) errorInternal(resource, action string) string {
	return fmt.Sprintf("Failed to %s in %s because there is internal error!", action, resource)
}

// Error when failed parsing body request
func (h *Handlers) errorParseBodyRequest() string {
	return "Failed to parse body request"
}

// Success Response message
func (h *Handlers) successResponse(resource, action string) string {
	return fmt.Sprintf("%s has been %s", resource, action)
}
