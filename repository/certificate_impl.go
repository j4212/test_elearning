// repository/certificate_impl.go

package repository

import (
    "errors"

    "github.com/cvzamannow/E-Learning-API/model"
    "gorm.io/gorm"
)

// certificateImpl mengimplementasikan CertificateRepository
type certificateImpl struct {
    DB *gorm.DB
}

// NewCertificateRepository membuat instance CertificateRepository dengan gorm DB
func NewCertificateRepository(db *gorm.DB) CertificateRepository {
    return &certificateImpl{DB: db}
}

// CreateCertificate mengimplementasikan metode untuk membuat sertifikat baru
func (repos *certificateImpl) CreateCertificate(request model.Certificate) (*model.Certificate, error) {
    if err := repos.DB.Create(&request).Error; err != nil {
        return nil, errors.New("error creating certificate")
    }
    return &request, nil
}
