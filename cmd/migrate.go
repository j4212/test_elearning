package cmd

import (
	"os"

	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DoMigrateUpCMD() cli.Command {

	return cli.Command{
		Name:  "migrate-up",
		Usage: "Run migration up with specific database source address",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "source",
				Value: "postgresql://simaku:RGY4JfM3CRA1RnnsfIIc@localhost:5432/simaku?sslmode=disable",
			},
		},
		Action: func(c *cli.Context) error {
			sourceDBArg := c.String("source")

			db, err := gorm.Open(postgres.Open(sourceDBArg), &gorm.Config{})

			if err != nil {
				logrus.Fatalf("[migrate-up] Failed to connect database source %s \n", err.Error())
				os.Exit(1)
			}

			err = db.AutoMigrate(
				&model.User{},
				&model.Course{},
				&model.CourseClass{},
				&model.Chapter{},
				&model.Material{},
				&model.Theory{},
				&model.Teacher{},
				&model.Submission{},
				&model.SubmissionStudent{},
				&model.ActiveStudent{},
				&model.Schools{},
				&model.AdminSchool{},
				&model.Quiz{},
				&model.Quizes{},
				&model.QuizAnswer{},
				&model.QuizAnswerStudent{},
				&model.ActiveStudentCourse{},
				&model.CompleteCourse{},
				&model.Certificate{},
			)

			if err != nil {
				logrus.Fatalf("[migrate-up] Failed to run migration because %s \n", err.Error())
				os.Exit(1)
			}

			logrus.Info("[migrate-up] Successfuly migrate up database...")

			return nil
		},
	}
}
