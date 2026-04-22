package service

import (
	"parcel_tracker/internal/model"
)

type ParcelRepository interface {
	Create(parcel model.Parcel) (model.Parcel, error)
	GetAll() ([]model.Parcel, error)
}

type ParcelService struct {
	repo ParcelRepository
}

func NewParcelService(repo ParcelRepository) *ParcelService {
	return &ParcelService{repo: repo}
}

func (s *ParcelService) CreateParcel(parcel model.Parcel) (model.Parcel, error) {
	return s.repo.Create(parcel)
}

func (s *ParcelService) GetAllParcels() ([]model.Parcel, error) {
	return s.repo.GetAll()
}
