package driver

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/ahmetfurkanunal/bitaksi-taxihub/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) CreateDriver(ctx context.Context, req CreateDriverRequest) (primitive.ObjectID, error) {
	if req.FirstName == "" || req.LastName == "" || req.Plate == "" {
		return primitive.NilObjectID, errors.New("firstName, lastName ve plate zorunludur")
	}

	now := time.Now()

	driver := Driver{
		ID:        primitive.NewObjectID(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plate:     req.Plate,
		TaxiType:  req.TaxiType,
		CarBrand:  req.CarBrand,
		CarModel:  req.CarModel,
		Location: Location{
			Latitude:  req.Lat,
			Longitude: req.Lon,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.Repo.Create(ctx, driver)
}

func (s *Service) UpdateDriver(ctx context.Context, id string, req UpdateDriverRequest) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("geçersiz id formatı")
	}

	updated := Driver{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plate:     req.Plate,
		TaxiType:  req.TaxiType,
		CarBrand:  req.CarBrand,
		CarModel:  req.CarModel,
		Location: Location{
			Latitude:  req.Lat,
			Longitude: req.Lon,
		},
		UpdatedAt: time.Now(),
	}

	return s.Repo.Update(ctx, objID, updated)
}

func (s *Service) ListDrivers(ctx context.Context) ([]Driver, error) {
	return s.Repo.List(ctx)
}

type NearbyDriver struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Plate      string  `json:"plate"`
	DistanceKm float64 `json:"distanceKm"`
}

func (s *Service) FindNearbyDrivers(ctx context.Context, lat, lon float64, taxiType string) ([]NearbyDriver, error) {
	all, err := s.Repo.List(ctx)
	if err != nil {
		return nil, err
	}

	const radiusKm = 6.0

	var list []NearbyDriver

	for _, d := range all {
		if taxiType != "" && d.TaxiType != taxiType {
			continue
		}

		dist := util.DistanceKm(lat, lon, d.Location.Latitude, d.Location.Longitude)
		if dist <= radiusKm {
			list = append(list, NearbyDriver{
				FirstName:  d.FirstName,
				LastName:   d.LastName,
				Plate:      d.Plate,
				DistanceKm: dist,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].DistanceKm < list[j].DistanceKm
	})

	return list, nil
}
