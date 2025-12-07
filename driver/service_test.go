package driver

import (
	"context"
	"testing"
)

func TestCreateDriver_EmptyPlateShouldFail(t *testing.T) {
	service := &Service{Repo: nil}

	req := CreateDriverRequest{
		FirstName: "Ahmet",
		LastName:  "Yılmaz",
		Plate:     "",
	}

	_, err := service.CreateDriver(context.Background(), req)
	if err == nil {
		t.Errorf("boş plate için hata bekleniyordu, nil döndü")
	}
}

func TestCreateDriver_EmptyFirstNameShouldFail(t *testing.T) {
	service := &Service{Repo: nil}

	req := CreateDriverRequest{
		FirstName: "",
		LastName:  "Yılmaz",
		Plate:     "34 ABC 123",
	}

	_, err := service.CreateDriver(context.Background(), req)
	if err == nil {
		t.Errorf("boş firstName için hata bekleniyordu, nil döndü")
	}
}
