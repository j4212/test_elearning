package repository

import "github.com/cvzamannow/E-Learning-API/model"

type CertificateRepository interface {
	CreateCertificate(request model.Certificate) (*model.Certificate, error)
}
