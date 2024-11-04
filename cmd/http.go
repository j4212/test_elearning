package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cvzamannow/E-Learning-API/config"
	"github.com/cvzamannow/E-Learning-API/http/handlers"
	"github.com/cvzamannow/E-Learning-API/middleware"
	"github.com/cvzamannow/E-Learning-API/repository"
	"github.com/cvzamannow/E-Learning-API/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/urfave/cli"
)

func HTTPGatewayServer(port int) {
	app := fiber.New()

	confDB := config.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	JWT_SECRET := os.Getenv("JWT_SECRET")

	// S3 Cloudflare Provider
	r2Cloudflare := service.NewR2Stub(&service.R2Cloudlfare{
		Bucket:       os.Getenv("R2_BUCKET"),
		AccountID:    os.Getenv("R2_ACCOUNT_ID"),
		Key:          os.Getenv("R2_KEY"),
		Secret:       os.Getenv("R2_SECRET"),
		PubBucketUrl: os.Getenv("R2_PUB"),
	})

	// Dependency Injection
	newDB := config.NewDBConfig(&confDB)
	courseRepos := repository.NewCourseRepository(newDB)

	// Dependency Injection auth / user
	userRepos := repository.NewAuthRepository(newDB)

	// Depedency injection User Management

	// Setting Schools
	schoolsRepos := repository.NewSchoolRepository(newDB)

	// Quiz
	quizRepos := repository.NewQuizRepository(newDB)
	//test
	certificateRepos := repository.NewCertificateRepository(newDB)
	// Grades
	gradesRepos := repository.NewGradesRepository(newDB)

	handlersDep := handlers.Handlers{
		R2Cloudflare:     r2Cloudflare,
		UserRepository:   userRepos,
		CourseRepository: courseRepos,
		JWT_SECRET:       []byte(JWT_SECRET),
		Middleware: middleware.Middleware{
			JwtSecret: []byte(JWT_SECRET),
		},
		SchoolRepository: schoolsRepos,
		QuizRepository:   quizRepos,
		GradesRepository: gradesRepos,
		CertificateRepo: certificateRepos,
	}

	// Setup global middleware
	app.Use(logger.New())

	// Setup Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	handlersDep.RouteCertificates(app)

	// Setup Course Route
	handlersDep.RouteCourses(app)

	// Setup Auth Route
	handlersDep.RouteAuth(app)

	// Setup User Management Route
	handlersDep.RouterUserManagemet(app)

	// Setup Schools Route
	handlersDep.RouterSchool(app)

	// Setup Quiz Route
	handlersDep.RouterQuiz(app)
	// R2 Storage Route
	handlersDep.RouteStorage(app)

	// grades router
	handlersDep.RouteGrades(app)

	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]string{
			"message": "E Learning API Version 1.0.0",
		})
	})

	app.Listen(fmt.Sprintf(":%d", port))
}

func HTTPGatewayServerCMD() cli.Command {
	return cli.Command{
		Name:  "http-gw-srv",
		Usage: "Run HTTP Gateway Server with specific port",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 8080,
			},
		},
		Action: func(c *cli.Context) error {
			port := c.Int("port")

			HTTPGatewayServer(port)
			return nil
		},
	}
}
