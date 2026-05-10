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
	query := `INSERT INTO  parcels (number, client_id, status, address) VALUES (:number, :client_id, :status, :address)`
	res, err := p.db.Exec(query,
		sql.Named("number", parcel.Number),
		sql.Named("client_id", parcel.ClientId),
		sql.Named("status", parcel.Status),
		sql.Named("address", parcel.Address))
	if err != nil {
		return parcel, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return parcel, err
	}
	parcel.ID = int(id)

	return parcel, nil
}

func (p *ParcelRepository) AllParcels() ([]model.Parcel, error) {
	query := `SELECT id, number, client_id, status, address, created_at FROM parcels`

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
			&parcel.ClientId,
			&parcel.Status,
			&parcel.Address,
			&parcel.CreatedAt,
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

func (p *ParcelRepository) AddClient(client model.Client) (model.Client, error) {
	query := `INSERT INTO clients (name, phone, address, email) VALUES (:name, :phone, :address, :email)`
	res, err := p.db.Exec(query,
		sql.Named("name", client.Name),
		sql.Named("phone", client.Phone),
		sql.Named("address", client.Address),
		sql.Named("email", client.Email))
	if err != nil {
		return client, err
	}

	id, _ := res.LastInsertId()
	client.ID = int(id)

	return client, nil
}

func (p *ParcelRepository) GetClientByID(id int) (model.Client, error) {

	var client model.Client
	query := `SELECT id, name, phone, address, email, created_at FROM clients WHERE id = :id`
	row := p.db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(
		&client.ID,
		&client.Name,
		&client.Phone,
		&client.Address,
		&client.Email,
		&client.CreatedAt,
	)
	if err != nil {
		return client, err
	}

	return client, nil
}

func (p *ParcelRepository) GetParcelByID(id int) (model.Parcel, error) {
	var parcel model.Parcel
	query := `SELECT id, client_id, number, status, address, created_at FROM parcels WHERE id = ?`
	row := p.db.QueryRow(query, id)
	err := row.Scan(
		&parcel.ID,
		&parcel.ClientId,
		&parcel.Number,
		&parcel.Status,
		&parcel.Address,
		&parcel.CreatedAt,
	)
	if err != nil {
		return parcel, err
	}

	return parcel, nil
}

func (p *ParcelRepository) AllClients() ([]model.Client, error) {
	query := `SELECT id, name, phone, address, email, created_at FROM clients`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []model.Client

	for rows.Next() {
		var client model.Client

		err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Phone,
			&client.Address,
			&client.Email,
			&client.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (p *ParcelRepository) DeleteParcelByID(id int) (int, error) {
	query := `DELETE FROM parcels WHERE id = ?`
	res, err := p.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows), nil
}

func (p *ParcelRepository) ModifyParcel(parcel model.Parcel) (int, error) {
	query := `UPDATE parcels SET status = :status, address = :address WHERE id = :id`

	res, err := p.db.Exec(query,
		sql.Named("status", parcel.Status),
		sql.Named("address", parcel.Address),
		sql.Named("id", parcel.ID),
	)

	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows), nil
}
