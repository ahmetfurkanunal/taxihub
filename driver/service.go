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
	return &Service{
		Repo: repo,
	}
}

func (s *Service) CreateDriver(ctx context.Context, req CreateDriverRequest) (primitive.ObjectID, error) {
	if req.FirstName == "" || req.LastName == "" {
		return primitive.NilObjectID, errors.New("first name and last name are required")
	}

	if req.Plate == "" {
		return primitive.NilObjectID, errors.New("plate number is required")
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
	if id == "" {
		return errors.New("driver id is required")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid driver id format")
	}

	updatedDriver := Driver{
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

	return s.Repo.Update(ctx, objID, updatedDriver)
}

func (s *Service) DeleteDriver(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("driver id is required")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid driver id format")
	}

	return s.Repo.DeleteByID(ctx, objID)
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
	allDrivers, err := s.Repo.List(ctx)
	if err != nil {
		return nil, err
	}

	const maxDistance = 6.0
	var nearbyList []NearbyDriver

	for _, driver := range allDrivers {
		if taxiType != "" && driver.TaxiType != taxiType {
			continue
		}

		distance := util.DistanceKm(
			lat,
			lon,
			driver.Location.Latitude,
			driver.Location.Longitude,
		)

		if distance <= maxDistance {
			nearbyList = append(nearbyList, NearbyDriver{
				FirstName:  driver.FirstName,
				LastName:   driver.LastName,
				Plate:      driver.Plate,
				DistanceKm: distance,
			})
		}
	}

	sort.Slice(nearbyList, func(i, j int) bool {
		return nearbyList[i].DistanceKm < nearbyList[j].DistanceKm
	})

	return nearbyList, nil
}
