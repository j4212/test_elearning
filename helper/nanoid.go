package helper

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateNanoId() (string, error) {
	id, err := gonanoid.New(20)

	if err != nil {
		return "", err
	}

	return id, nil
}
