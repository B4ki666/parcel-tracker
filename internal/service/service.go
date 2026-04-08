package service

type Parcel struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	Status string `json:"status"`
}

type ParcelService struct {
	Parcels   map[int]Parcel
	idCounter int
}

func NewParcelService() *ParcelService {
	return &ParcelService{
		Parcels: make(map[int]Parcel),
	}
}

func (p *ParcelService) CreateParcel(parcel Parcel) Parcel {
	p.idCounter++
	parcel.ID = p.idCounter

	p.Parcels[parcel.ID] = parcel

	return parcel
}

func (p *ParcelService) GetAllParcels() []Parcel {
	parcels := make([]Parcel, 0, len(p.Parcels))

	for _, value := range p.Parcels {
		parcels = append(parcels, value)
	}

	return parcels
}
