package handlers

import (
	"github.com/cvzamannow/E-Learning-API/helper"
	"github.com/cvzamannow/E-Learning-API/http"
	"github.com/cvzamannow/E-Learning-API/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
)

func (h *Handlers) RouterQuiz(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Post("/quizz/:id", h.Middleware.Protected(), h.CreateQuizHandler)
	v1.Put("/quizz/:id", h.Middleware.Protected(), h.UpdateQuizHandler)
	v1.Delete("/quizz/:id", h.Middleware.Protected(), h.DeleteQuizHandler)
	v1.Get("/quizz/:id", h.Middleware.Protected(), h.GetDetailQuizByIdHandler)

	// quiz answer
	v1.Post("/quizz/answer-student/:id", h.Middleware.Protected(), h.CreateQuizAnswerHandler)
}

// CreateQuizHandler handles HTTP request to create a quiz.
func (h *Handlers) CreateQuizHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	chapter, err := h.CourseRepository.FindChapter(map[string]interface{}{
		"id": id,
	})

	var quizzDataResponse []http.QuizData

	role := c.Locals("role").(string)

	if role != "TEACHER" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	var request http.Quiz
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	materialID, err := helper.GenerateNanoId()
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error generating nano ID",
			Data:    nil,
		})
	}
	result, err := h.QuizRepository.CreateQuiz(
		model.Quiz{
			ID:          string(materialID),
			Title:       request.Title,
			Description: request.Description,
			ChapterID:   string(chapter.ID),
			Material: model.Material{
				ID:        string(materialID),
				ChapterID: chapter.ID,
				Title:     request.Title,
				Type:      "QUIZ",
				Slug:      slug.Make(request.Title),
			},
		},
	)
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Quiz", "create"),
			Data:    nil,
		})
	}

	for _, quizes := range request.Quizes {
		quizesID, err := helper.GenerateNanoId()
		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Error generating nano ID",
				Data:    nil,
			})
		}

		QuizAnswerResponse := []http.QuizAnswer{}

		quizzData, err := h.QuizRepository.CreateQuizes(model.Quizes{
			ID:     string(quizesID),
			Quiz:   quizes.Quiz,
			QuizID: string(result.ID),
			ImgURL: *quizes.ImgURL,
		})
		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("Quizzes", "create"),
				Data:    nil,
			})
		}

		for _, answer := range quizes.Answers {
			answerID, err := helper.GenerateNanoId()
			if err != nil {
				return c.Status(500).JSON(&http.WebResponse{
					Status:  "error",
					Message: "Error generating nano ID",
					Data:    nil,
				})
			}

			quizzAnswer, err := h.QuizRepository.CreateQuizAnswer(model.QuizAnswer{
				ID:        string(answerID),
				Answer:    answer.Answer,
				QuizesID:  string(quizesID),
				IsCorrect: answer.IsCorrect,
			})

			if err != nil {
				return c.Status(500).JSON(&http.WebResponse{
					Status:  "error",
					Message: h.errorInternal("Quiz Answers", "create"),
					Data:    nil,
				})
			}

			var found bool
			for _, itemAnswer := range QuizAnswerResponse {
				if itemAnswer.ID == quizzAnswer.QuizesID {
					found = true
					break
				}
			}

			if !found {
				QuizAnswerResponse = append(QuizAnswerResponse, http.QuizAnswer{
					ID:        string(quizzAnswer.ID),
					Answer:    quizzAnswer.Answer,
					QuizesID:  string(quizesID),
					IsCorrect: quizzAnswer.IsCorrect,
				})
			}
		}

		quizzDataResponse = append(quizzDataResponse, http.QuizData{
			ID:      quizzData.ID,
			Quiz:    quizzData.Quiz,
			QuizID:  string(result.ID),
			ImgURL:  &quizzData.ImgURL,
			Answers: QuizAnswerResponse,
		})
	}

	typeOfMaterial := "QUIZ"

	response := http.Quiz{
		ID:          string(materialID),
		Title:       request.Title,
		Description: request.Description,
		ChapterID:   id,
		Type:        &typeOfMaterial,
		Quizes:      quizzDataResponse,
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Quiz", "created"),
		Data:    response,
	})
}

// UpdateQuizHandler menangani permintaan HTTP untuk memperbarui sebuah kuis.
func (h *Handlers) UpdateQuizHandler(c *fiber.Ctx) error {
	// Ambil peran pengguna dari konteks lokal
	role := c.Locals("role").(string)

	// Cek apakah pengguna adalah TEACHER, jika bukan kembalikan status 401 Unauthorized
	if role != "TEACHER" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	// Ambil ID kuis dari parameter URL
	id := c.Params("id")

	// Parse tubuh permintaan menjadi struktur http.Quiz
	var request http.Quiz
	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	// Validasi input: cek apakah judul atau deskripsi kosong
	if request.Title == "" || request.Title == " " || request.Description == "" || request.Description == " " {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Title atau description tidak boleh kosong",
			Data:    nil,
		})
	}

	// Perbarui kuis di repository
	result, err := h.QuizRepository.UpdateQuiz(
		map[string]interface{}{
			"id": id,
		},
		model.Quiz{
			Title:       request.Title,
			Description: request.Description,
			Material: model.Material{
				ID:    id,
				Title: request.Title,
				Slug:  slug.Make(request.Title),
			},
		},
	)
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Quiz", "update"),
			Data:    nil,
		})
	}

	// Jika kuis dalam permintaan kosong, hapus semua kuis terkait
	if request.Quizes == nil || len(request.Quizes) == 0 {
		findQuizes, err := h.QuizRepository.GetQuizesByIdQuiz(id)

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("Quizzes", "find"),
				Data:    nil,
			})
		}

		for _, quiz := range *findQuizes {
			// Hapus jawaban kuis
			findQuizAnswer, err := h.QuizRepository.GetQuizAnswerByIdQuizes(quiz.ID)

			if err != nil {
				return c.Status(500).JSON(&http.WebResponse{
					Status:  "error",
					Message: h.errorInternal("Quiz Answers", "find"),
					Data:    nil,
				})
			}

			for _, quizAnswer := range *findQuizAnswer {
				err := h.QuizRepository.DeleteQuizAnswer(map[string]interface{}{
					"id": quizAnswer.ID,
				})

				if err != nil {
					return c.Status(500).JSON(&http.WebResponse{
						Status:  "error",
						Message: h.errorInternal("Quiz Answers", "delete"),
						Data:    nil,
					})
				}

			}

			err = h.QuizRepository.DeleteQuizes(map[string]interface{}{
				"quiz_id": quiz.QuizID,
			})

			if err != nil {
				return c.Status(500).JSON(&http.WebResponse{
					Status:  "error",
					Message: h.errorInternal("Quizzes", "delete"),
					Data:    nil,
				})
			}

		}

	}

	// Perbarui material di repository
	_, err = h.CourseRepository.UpdateMaterial(map[string]interface{}{
		"id": id,
	},
		model.Material{
			Title: request.Title,
			Slug:  slug.Make(request.Title),
		},
	)

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Material", "update"),
			Data:    nil,
		})
	}

	findChapterByID, err := h.CourseRepository.FindChapter(map[string]interface{}{
		"id": request.ChapterID,
	})

	// Buat respons kuis
	var response http.Quiz
	typeQuiz := "QUIZ"
	response = http.Quiz{
		ID:          id,
		Title:       result.Title,
		ChapterID:   request.ChapterID,
		CourseID:    findChapterByID.CourseID,
		Description: result.Description,
		Type:        &typeQuiz,
	}

	response.Quizes = make([]http.QuizData, len(request.Quizes))

	// Iterasi melalui setiap kuis dalam permintaan
	for i, quizes := range request.Quizes {

		if quizes.Quiz == "" || quizes.Quiz == " " {
			return c.Status(400).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Quiz tidak boleh kosong",
				Data:    nil,
			})
		}

		_, err := h.QuizRepository.UpdateQuizes(
			map[string]interface{}{
				"id": quizes.ID,
			},
			model.Quizes{
				Quiz:   quizes.Quiz,
				ImgURL: *quizes.ImgURL,
			},
		)
		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("Quizzes", "update"),
				Data:    nil,
			})
		}

		// Jika kuis ditandai untuk dihapus, hapus jawaban kuis dan kuis itu sendiri
		if quizes.Delete != nil && *quizes.Delete == true {
			err = h.QuizRepository.DeleteQuizAnswer(map[string]interface{}{
				"quizes_id": quizes.ID,
			})

			if err != nil {
				return c.Status(500).JSON(&http.WebResponse{
					Status:  "error",
					Message: h.errorInternal("Quiz Answers", "delete"),
					Data:    nil,
				})
			}

			err = h.QuizRepository.DeleteQuizes(map[string]interface{}{
				"id": quizes.ID,
			})

			if err != nil {
				return c.Status(500).JSON(&http.WebResponse{
					Status:  "error",
					Message: h.errorInternal("Quizzes", "delete"),
					Data:    nil,
				})
			}

		}

		if quizes.Delete != nil && *quizes.Delete == true {
			response.Quizes[i] = http.QuizData{}
			continue
		} else {
			if quizes.ID != "" {
				response.Quizes[i].ID = quizes.ID
			}
			response.Quizes[i].Quiz = quizes.Quiz
			response.Quizes[i].QuizID = quizes.QuizID
			response.Quizes[i].ImgURL = quizes.ImgURL
		}
		response.Quizes[i].Answers = make([]http.QuizAnswer, len(quizes.Answers))

		// Iterasi melalui setiap jawaban dalam kuis
		for j, answer := range quizes.Answers {

			if answer.Answer == "" || answer.Answer == " " {
				return c.Status(400).JSON(&http.WebResponse{
					Status:  "error",
					Message: "Answer tidak boleh kosong",
					Data:    nil,
				})
			}

			// Jika kuis tidak ada, buat kuis baru
			if quizes.ID == "" || quizes.ID == " " {
				response.Quizes[i].Answers[j].QuizesID = response.Quizes[i].ID

				findQuizes := h.QuizRepository.FindQuizesNoError(map[string]interface{}{
					"id": quizes.ID,
				})

				if findQuizes == nil || findQuizes.ID == "" {

					var idQuizAnswerID string

					idQuizes, err := helper.GenerateNanoId()
					if err != nil {
						return c.Status(500).JSON(&http.WebResponse{
							Status:  "error",
							Message: "Error generating nano ID",
							Data:    nil,
						})
					}
					if response.Quizes[i].Answers[j].QuizesID == "" {

						findQuizes.ID = idQuizes

						// Buat kuis baru
						resultQuizID, err := h.QuizRepository.CreateQuizes(model.Quizes{
							ID:     string(idQuizes),
							Quiz:   quizes.Quiz,
							QuizID: string(id),
						})

						if err != nil {
							return c.Status(500).JSON(&http.WebResponse{
								Status:  "error",
								Message: h.errorInternal("Quizzes", "create"),
								Data:    nil,
							})
						}

						idQuizAnswer, err := helper.GenerateNanoId()

						if err != nil {
							return c.Status(500).JSON(&http.WebResponse{
								Status:  "error",
								Message: "Error generating nano ID",
								Data:    nil,
							})
						}

						// Buat jawaban kuis baru
						_, err = h.QuizRepository.CreateQuizAnswer(model.QuizAnswer{
							ID:        string(idQuizAnswer),
							Answer:    answer.Answer,
							QuizesID:  resultQuizID.ID,
							IsCorrect: answer.IsCorrect,
						})

						if err != nil {
							return c.Status(500).JSON(&http.WebResponse{
								Status:  "error",
								Message: h.errorInternal("Quiz Answers", "create"),
								Data:    nil,
							})
						}
						idQuizAnswerID = resultQuizID.ID
						response.Quizes[i].ID = resultQuizID.ID

						// response.Quizes[i].Answers[j].QuizesID = resultQuizID.ID
					}

					idQuizAnswer, err := helper.GenerateNanoId()

					if err != nil {
						return c.Status(500).JSON(&http.WebResponse{
							Status:  "error",
							Message: "Error generating nano ID",
							Data:    nil,
						})
					}

					if response.Quizes[i].Answers[j].QuizesID != "" {
						// Buat jawaban kuis baru
						result, err := h.QuizRepository.CreateQuizAnswer(model.QuizAnswer{
							ID:        string(idQuizAnswer),
							Answer:    answer.Answer,
							QuizesID:  response.Quizes[i].Answers[j].QuizesID,
							IsCorrect: answer.IsCorrect,
						})

						if err != nil {
							return c.Status(500).JSON(&http.WebResponse{
								Status:  "error",
								Message: h.errorInternal("Quiz Answers", "create"),
								Data:    nil,
							})
						}

						response.Quizes[i].Answers[j].QuizesID = result.QuizesID
						response.Quizes[i].ID = result.QuizesID

					}

					if idQuizAnswerID != "" {
						response.Quizes[i].Answers[j].QuizesID = idQuizAnswerID
					}

					if answer.ID != "" {
						response.Quizes[i].Answers[j].ID = answer.ID
					} else {
						response.Quizes[i].Answers[j].ID = idQuizes
					}

					if quizes.ID != "" {
						response.Quizes[i].Answers[j].QuizesID = quizes.ID
					}

					if answer.Delete != nil && *answer.Delete == true {
						response.Quizes[i].Answers[j] = http.QuizAnswer{}
					} else {
						// response.Quizes[i].Answers[j].QuizesID =
						response.Quizes[i].Answers[j].ID = idQuizAnswer
						response.Quizes[i].QuizID = id
					}

				} else {
					idQuizAnswer, err := helper.GenerateNanoId()

					if err != nil {
						return c.Status(500).JSON(&http.WebResponse{
							Status:  "error",
							Message: "Error generating nano ID",
							Data:    nil,
						})
					}

					// Buat jawaban kuis baru
					resultQuizAnswerCreate, err := h.QuizRepository.CreateQuizAnswer(model.QuizAnswer{
						ID:        string(idQuizAnswer),
						Answer:    answer.Answer,
						QuizesID:  findQuizes.ID,
						IsCorrect: answer.IsCorrect,
					})

					if err != nil {
						return c.Status(500).JSON(&http.WebResponse{
							Status:  "error",
							Message: h.errorInternal("Quiz Answers", "create"),
							Data:    nil,
						})
					}

					if answer.ID != "" {
						response.Quizes[i].Answers[j].ID = answer.ID
					} else {
						response.Quizes[i].Answers[j].ID = idQuizAnswer
					}

					if quizes.ID != "" {
						response.Quizes[i].Answers[j].QuizesID = quizes.ID
					} else {
						response.Quizes[i].Answers[j].QuizesID = findQuizes.ID
					}

					if answer.Delete != nil && *answer.Delete == true {
						response.Quizes[i].Answers[j] = http.QuizAnswer{}
					} else {
						response.Quizes[i].Answers[j].ID = resultQuizAnswerCreate.ID
						response.Quizes[i].Answers[j].QuizesID = findQuizes.ID
					}
				}

			} else if answer.Delete != nil && *answer.Delete == true {

				// Hapus jawaban kuis
				err = h.QuizRepository.DeleteQuizAnswer(map[string]interface{}{
					"id": answer.ID,
				})

				if err != nil {
					return c.Status(500).JSON(&http.WebResponse{
						Status:  "error",
						Message: h.errorInternal("Quiz Answers", "delete"),
						Data:    nil,
					})
				}

			} else {
				// Perbarui jawaban kuis
				_, err = h.QuizRepository.UpdateQuizAnswer(
					map[string]interface{}{
						"id": answer.ID,
					},
					model.QuizAnswer{
						Answer:    answer.Answer,
						IsCorrect: answer.IsCorrect,
					},
				)

				if err != nil {
					idQuizAnswer, err := helper.GenerateNanoId()

					if err != nil {
						return c.Status(500).JSON(&http.WebResponse{
							Status:  "error",
							Message: "Error generating nano ID",
							Data:    nil,
						})
					}
					result, err := h.QuizRepository.CreateQuizAnswer(model.QuizAnswer{
						ID:        string(idQuizAnswer),
						Answer:    answer.Answer,
						QuizesID:  quizes.ID,
						IsCorrect: answer.IsCorrect,
					})
					response.Quizes[i].Answers[j].ID = result.ID
				}
			}

			if answer.Delete != nil && *answer.Delete == true {
				response.Quizes[i].Answers[j] = http.QuizAnswer{}
				continue
			} else {
				if answer.ID != "" {
					response.Quizes[i].Answers[j].ID = answer.ID
				}
				if quizes.ID != "" {
					response.Quizes[i].Answers[j].QuizesID = quizes.ID
				}
				response.Quizes[i].Answers[j].Answer = answer.Answer
				response.Quizes[i].Answers[j].IsCorrect = answer.IsCorrect

			}
		}

		// end quiz create nya jika data quiz nya tidak ada

	}

	// Kembalikan respons sukses dengan data kuis yang diperbarui
	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Quiz", "updated"),
		Data:    response,
	})
}

// DeleteQuizHandler handles HTTP request to delete a quiz.
func (h *Handlers) DeleteQuizHandler(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "TEACHER" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	id := c.Params("id")

	err := h.QuizRepository.DeleteQuiz(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Quiz", "delete"),
			Data:    nil,
		})
	}

	_, err = h.CourseRepository.DeleteMaterial(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorInternal("Material", "delete"),
			Data:    nil,
		})
	}

	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Quiz", "deleted"),
		Data:    err,
	})

}

// GetDetailQuizBySlugHandler handles HTTP request to get detailed quiz information by slug.
func (h *Handlers) GetDetailQuizByIdHandler(c *fiber.Ctx) error {
	USERID := c.Locals("user_id").(string)

	findStudentActuveStudent := h.UserRepository.FindActiveStudentNoError(map[string]interface{}{
		"student_id": USERID,
	})

	// declare struct response
	var quizesResponse []http.QuizesResponseHTTP
	var quizResponse http.QuizResponseHTTP

	id := c.Params("id")

	// Get quiz by  ID
	resultQuiz, err := h.QuizRepository.FindQuiz(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error get detail quiz",
			Data:    nil,
		})
	}

	// find quiz answer student
	resultAnswerStudent, err := h.QuizRepository.FindQuizAnswerStudent(map[string]interface{}{
		"quiz_id": resultQuiz.ID,
	})

	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error get quiz answer student",
			Data:    nil,
		})
	}

	// find quizes by id quiz
	resultQuizes, err := h.QuizRepository.GetQuizesByIdQuiz(resultQuiz.ID)
	if err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Error get detail quiz",
			Data:    nil,
		})
	}

	for _, item := range *resultQuizes {

		// declare struct response for each loop
		var quizAnswerResponse []http.QuizAnswerHTTP
		var answerQuizStudent []http.AnswerStuedntResponseHTTP

		// find Answer
		resultAnswerStudentItem, err := h.QuizRepository.GetQuizAnswerStudentByIdQuiz(id)

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Error get quiz answer student",
				Data:    nil,
			})
		}

		// untuk response answer nya
		for _, answer := range item.QuizAnswers {
			if answer.QuizesID == item.ID {
				quizAnswerResponse = append(quizAnswerResponse, http.QuizAnswerHTTP{
					ID:        answer.ID,
					Answer:    answer.Answer,
					IsCorrect: answer.IsCorrect,
					CreatedAt: &answer.CreatedAt,
					UpdatedAt: &answer.UpdatedAt,
				})
			}

			if findStudentActuveStudent != nil {

				for _, answerStudent := range *resultAnswerStudentItem {
					if answerStudent.QuizesID == item.ID {
						if (len(answerQuizStudent) != 1) &&
							(answerStudent.QuizAnswer.Answer == answer.Answer) && (answerStudent.ActiveStudentID == findStudentActuveStudent.ID) {
							answerQuizStudent = append(
								answerQuizStudent,
								http.AnswerStuedntResponseHTTP{
									ID:        answerStudent.ID,
									Name:      answerStudent.ActiveStudent.Student.Name,
									Answer:    answerStudent.QuizAnswer.Answer,
									CreatedAt: &answerStudent.CreatedAt,
									UpdatedAt: &answerStudent.UpdatedAt,
								},
							)
						}
					}
				}

			}

		}

		if findStudentActuveStudent != nil {
			// untuk response quizes nya
			quizesResponse = append(quizesResponse, http.QuizesResponseHTTP{
				ID:            item.ID,
				Quiz:          item.Quiz,
				Answer:        quizAnswerResponse,
				AnswerStudent: answerQuizStudent,
				ImgURL:        &item.ImgURL,
				CreatedAt:     &item.CreatedAt,
				UpdatedAt:     &item.UpdatedAt,
			})
		} else {
			// untuk response quizes nya
			quizesResponse = append(quizesResponse, http.QuizesResponseHTTP{
				ID:            item.ID,
				Quiz:          item.Quiz,
				Answer:        quizAnswerResponse,
				ImgURL:        &item.ImgURL,
				AnswerStudent: []http.AnswerStuedntResponseHTTP{},
				CreatedAt:     &item.CreatedAt,
				UpdatedAt:     &item.UpdatedAt,
			})
		}

	}

	var nextMaterial *http.NextMaterialHTTP
	next, err := h.CourseRepository.NextMaterial(&helper.NextMaterialArg{
		ChapterID:         resultQuiz.ChapterID,
		CurrentMaterialID: resultQuiz.ID,
		CreatedAt:         resultQuiz.CreatedAt,
	}, false)

	if err == nil && next.ChapterID == resultQuiz.ChapterID {
		nextMaterial = &http.NextMaterialHTTP{
			ID:   next.ID,
			Type: next.Type,
		}
	}

	// Check if next chapter exist
	if next == nil {
		findChapterDetail, err := h.CourseRepository.FindChapter(map[string]interface{}{
			"id": resultQuiz.ChapterID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "Error",
				Message: "Quiz not found",
				Data:    nil,
			})
		}

		nextChapter := h.CourseRepository.NextChapter(resultQuiz.ChapterID, findChapterDetail.CourseID)

		if nextChapter != nil {
			doNext, err := h.CourseRepository.NextMaterial(&helper.NextMaterialArg{
				ChapterID:         nextChapter.ID,
				CurrentMaterialID: resultQuiz.ID,
				CreatedAt:         resultQuiz.CreatedAt,
			}, true)

			if err == nil {
				nextMaterial = &http.NextMaterialHTTP{
					ID:   doNext.ID,
					Type: doNext.Type,
				}
			}
		}
	}

	quizResponse = http.QuizResponseHTTP{
		ID:          resultQuiz.ID,
		Title:       resultQuiz.Title,
		Description: resultQuiz.Description,
		Grades:      resultAnswerStudent.Grades,
		Score:       resultAnswerStudent.Score,
		Quiz:        quizesResponse,
		Next:        nextMaterial,
		CreatedAt:   &resultQuiz.CreatedAt,
		UpdatedAt:   &resultQuiz.UpdatedAt,
	}

	// Return success response
	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: "Get detail quiz successfully",
		Data:    quizResponse,
	})
}

// TODO: Create quiz answer
func (h *Handlers) CreateQuizAnswerHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if role != "STUDENT" {
		return c.Status(401).JSON(&http.WebResponse{
			Status:  "error",
			Message: "You are not authorized to perform this action",
			Data:    nil,
		})
	}

	var request http.QuizAnswerStudent
	requestID := c.Params("id")
	var score int
	var grade float64
	var totalQuestions int
	var response []http.QuizAnswerStudentResponseHTTP

	findStudentActuveStudent, err := h.UserRepository.FindActiveStudent(map[string]interface{}{
		"student_id": userID,
	})

	if err != nil {
		return c.Status(400).JSON(&http.WebResponse{
			Status:  "error",
			Message: "Active Student not found",
			Data:    nil,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(500).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorParseBodyRequest(),
			Data:    nil,
		})
	}

	// Cari materi kuis berdasarkan slug
	material, err := h.CourseRepository.FindMaterial(map[string]interface{}{
		"id": requestID,
	})

	if err != nil {
		return c.Status(404).JSON(&http.WebResponse{
			Status:  "error",
			Message: h.errorSpecifyResource("id"),
			Data:    nil,
		})
	}

	for _, item := range request.Answer {
		var AnswerCorrect string

		findAnswer, err := h.QuizRepository.FindQuizAnswer(map[string]interface{}{
			"id": item.QuizAnswerID,
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("Quiz Answers", "find"),
				Data:    nil,
			})
		}

		// Periksa apakah jawaban benar atau tidak
		if findAnswer.IsCorrect {
			score++
		}

		totalQuestions++

		// Hitung nilai dalam bentuk persentase
		grade = (float64(score) / float64(totalQuestions)) * 100

		resultQuiz, err := h.QuizRepository.GetQuizByMaterialId(material.ID)

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: "Error get detail quiz",
				Data:    nil,
			})
		}

		id, err := helper.GenerateNanoId()
		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("Quiz Answer", "create"),
				Data:    nil,
			})
		}

		result, err := h.QuizRepository.CreateQuizAnswerStudent(model.QuizAnswerStudent{
			ID:              string(id),
			QuizID:          resultQuiz.ID,
			ActiveStudentID: findStudentActuveStudent.ID,
			QuizesID:        item.QuizesID,
			QuizAnswerID:    item.QuizAnswerID,
			Score:           score,
			Grades:          int(grade),
		})

		if err != nil {
			return c.Status(500).JSON(&http.WebResponse{
				Status:  "error",
				Message: h.errorInternal("Quiz Answer", "create"),
				Data:    nil,
			})
		}

		findQuizes, _ := h.QuizRepository.FindQuizes(map[string]interface{}{
			"id": item.QuizesID,
		})

		// mencari jawaban yang benar
		if !findAnswer.IsCorrect {
			logrus.Println("ITEM ", item.QuizAnswerID)
			findAnswerCorrect, _ := h.QuizRepository.FindQuizAnswer(map[string]interface{}{
				"quizes_id": item.QuizesID,
			})
			logrus.Println("findAnswerCorrect ", findAnswerCorrect)
			if findAnswerCorrect.IsCorrect {
				AnswerCorrect = findAnswerCorrect.Answer
			}

		}

		response = append(response, http.QuizAnswerStudentResponseHTTP{
			ID:            result.ID,
			QuizID:        result.QuizID,
			Quiz:          findQuizes.Quiz,
			QuizesID:      result.QuizesID,
			Answer:        findAnswer.Answer,
			AnswerCorrect: &AnswerCorrect,
		})

	}
	return c.Status(200).JSON(&http.WebResponse{
		Status:  "success",
		Message: h.successResponse("Quiz Answer", "created"),
		Data:    response,
	})
}
