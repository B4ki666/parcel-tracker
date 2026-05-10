package service

import (
	"database/sql"
	"errors"
	"net/mail"
	"parcel_tracker/internal/model"
	"strings"

	"modernc.org/sqlite"
)

const (
	ParcelStatusCreated    = "created"
	ParcelStatusDelivered  = "delivered"
	ParcelStatusInTransit  = "in_transit"
	sqliteUniqueConstraint = 2067
)

type ParcelAndClientRepository interface {
	Create(parcel model.Parcel) (model.Parcel, error)
	AllParcels() ([]model.Parcel, error)
	GetClientByID(id int) (model.Client, error)
	GetParcelByID(id int) (model.Parcel, error)
	AddClient(client model.Client) (model.Client, error)
	AllClients() ([]model.Client, error)
	DeleteParcelByID(id int) (int, error)
	ModifyParcel(model.Parcel) (int, error)
}

var parcelStatuses = map[string]struct{}{
	ParcelStatusCreated:   {},
	ParcelStatusDelivered: {},
	ParcelStatusInTransit: {},
}

type ParcelService struct {
	repo ParcelAndClientRepository
}

func NewParcelService(repo ParcelAndClientRepository) *ParcelService {
	return &ParcelService{repo: repo}
}

func (s *ParcelService) CreateParcel(parcel model.Parcel) (model.Parcel, error) {

	validParcel, err := s.checkAndFixParcel(parcel)
	if err != nil {
		return parcel, err
	}

	_, err = s.repo.GetClientByID(*validParcel.ClientId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return parcel, model.NewError(404, "client not found")
		}
		return parcel, err
	}

	return s.repo.Create(validParcel)

}

func (s *ParcelService) GetAllParcels() ([]model.Parcel, error) {
	return s.repo.AllParcels()
}

func (s *ParcelService) GetAllClients() ([]model.Client, error) {
	return s.repo.AllClients()
}

func isEmailValid(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return addr.Address == email
}

func (s *ParcelService) CreateClient(client model.Client) (model.Client, error) {
	if strings.TrimSpace(client.Name) == "" {
		return client, model.NewError(400, "the name field is empty")
	}

	if strings.TrimSpace(client.Phone) == "" {
		return client, model.NewError(400, "the phone field is empty")
	}

	if strings.TrimSpace(client.Address) == "" {
		return client, model.NewError(400, "the address field is empty")
	}

	if emailIsValid := isEmailValid(client.Email); !emailIsValid {
		return client, model.NewError(400, "the Email field has an invalid format")
	}

	resClient, err := s.repo.AddClient(client)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqliteUniqueConstraint {
				return client, model.NewError(409, "email already registered")
			}

		}
		return client, err
	}

	return resClient, nil
}

func (s *ParcelService) checkAndFixParcel(parcel model.Parcel) (model.Parcel, error) {

	if strings.TrimSpace(parcel.Number) == "" {
		return parcel, model.NewError(400, "incorrect parcel number")
	}

	if parcel.ClientId == nil {
		return parcel, model.NewError(400, "the field client_id is empty")
	}

	if strings.TrimSpace(parcel.Status) == "" {
		parcel.Status = ParcelStatusCreated
		return parcel, nil
	}

	if _, ok := parcelStatuses[parcel.Status]; ok {
		return parcel, nil
	}

	return parcel, model.NewError(400, "incorrect parcel status")
}

func (s *ParcelService) GetParcel(parcelID int) (model.Parcel, error) {
	var parcel model.Parcel

	if parcelID < 1 {
		return model.Parcel{}, model.NewError(400, "incorrect id")
	}

	parcel, err := s.repo.GetParcelByID(parcelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return parcel, model.NewError(404, "parcel not found")
		}
		return parcel, err
	}

	return parcel, nil
}

func (s *ParcelService) DeleteParcel(parcelID int) error {
	if parcelID < 1 {
		return model.NewError(400, "incorrect id")
	}

	rowsAffected, err := s.repo.DeleteParcelByID(parcelID)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return model.NewError(404, "parcel not found")
	}

	return nil
}

func (s *ParcelService) PatсhParcel(parcelID int, patchParcel model.UpdateParcelRequest) (model.Parcel, error) {

	if parcelID < 1 {
		return model.Parcel{}, model.NewError(400, "incorrect id")
	}

	parcel, err := s.repo.GetParcelByID(parcelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Parcel{}, model.NewError(404, "parcel not found")
		}
		return model.Parcel{}, err
	}

	if patchParcel.Status != nil {
		if _, ok := parcelStatuses[*patchParcel.Status]; ok {
			parcel.Status = *patchParcel.Status
		} else {
			return model.Parcel{}, model.NewError(400, "incorrect status")
		}
	}

	if patchParcel.Address != nil {
		if strings.TrimSpace(*patchParcel.Address) == "" {
			return model.Parcel{}, model.NewError(400, "no address specified")
		}
		parcel.Address = *patchParcel.Address

	}

	rowsAffected, err := s.repo.ModifyParcel(parcel)

	if err != nil {
		return model.Parcel{}, err
	}

	if rowsAffected == 0 {
		return model.Parcel{}, model.NewError(404, "parcel not found")
	}

	return parcel, nil
}
