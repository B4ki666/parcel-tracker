package repository

import (
	"database/sql"
	"parcel_tracker/internal/model"
)

type ParcelRepository struct {
	db *sql.DB
}

func NewParcelRepository(db *sql.DB) *ParcelRepository {
	return &ParcelRepository{
		db: db,
	}
}

func (p *ParcelRepository) Create(parcel model.Parcel) (model.Parcel, error) {
	query := `INSERT INTO  parcels (number, status, address) VALUES (:number, :status, :address)`
	res, err := p.db.Exec(query,
		sql.Named("number", parcel.Number),
		sql.Named("status", parcel.Status),
		sql.Named("address", parcel.Address))
	if err != nil {
		return parcel, err
	}

	id, _ := res.LastInsertId()
	parcel.ID = int(id)

	return parcel, nil
}

func (p *ParcelRepository) GetAll() ([]model.Parcel, error) {
	query := `SELECT id, number, status, address, created_at FROM parcels`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []model.Parcel

	for rows.Next() {
		var parcel model.Parcel

		err := rows.Scan(
			&parcel.ID,
			&parcel.Number,
			&parcel.Status,
			&parcel.Address,
			&parcel.Created_at,
		)
		if err != nil {
			return nil, err
		}

		parcels = append(parcels, parcel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}
