package service

import (
	"errors"
	"testing"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/memory"
)

func TestCreate(t *testing.T) {
	repo := memory.NewSpyCatRepository()
	service := NewSpyCatService(repo)
	spyCat := &model.SpyCat{
		Name:              "Pickachu",
		YearsOfExperience: 2,
		Breed:             "pokemon",
		Salary:            100,
	}
	err := service.Create(t.Context(), spyCat)
	if err != nil {
		t.Fatal(err)
	}
	if spyCat.Id == 0 {
		t.Fatal("Spy cat ID was not set")
	}
}

func TestGetById(t *testing.T) {
	repo := memory.NewSpyCatRepository()
	service := NewSpyCatService(repo)
	spyCat := &model.SpyCat{
		Name:              "Pickachu",
		YearsOfExperience: 2,
		Breed:             "pokemon",
		Salary:            100,
	}
	err := service.Create(t.Context(), spyCat)
	if err != nil {
		t.Fatal(err)
	}

	gotSpyCat, err := service.GetById(t.Context(), spyCat.Id)
	if err != nil {
		t.Fatal(err)
	}
	if gotSpyCat == nil {
		t.Fatal("Spy cat was not found")
	}
	if gotSpyCat.Id != spyCat.Id {
		t.Fatal("Spy cat ID was not found")
	}
}

func TestRemove(t *testing.T) {
	repo := memory.NewSpyCatRepository()
	service := NewSpyCatService(repo)
	spyCat := &model.SpyCat{
		Name:              "Pickachu",
		YearsOfExperience: 2,
		Breed:             "pokemon",
		Salary:            100,
	}
	err := service.Create(t.Context(), spyCat)
	if err != nil {
		t.Fatal(err)
	}

	err = service.Remove(t.Context(), spyCat.Id)
	if err != nil {
		t.Fatal(err)
	}

	gotSpyCat, err := service.GetById(t.Context(), spyCat.Id)
	if !errors.Is(err, storage.ErrorModelNotFound) {
		t.Fatal(err)
	}
	if gotSpyCat != nil {
		t.Fatal("Spy cat was not removed")
	}
}

func TestUpdateSalary(t *testing.T) {
	repo := memory.NewSpyCatRepository()
	service := NewSpyCatService(repo)
	spyCat := &model.SpyCat{
		Name:              "Pickachu",
		YearsOfExperience: 2,
		Breed:             "pokemon",
		Salary:            100,
	}
	err := service.Create(t.Context(), spyCat)
	if err != nil {
		t.Fatal(err)
	}

	newSalary := 200.0
	err = service.UpdateSalary(t.Context(), spyCat.Id, newSalary)
	if err != nil {
		t.Fatal(err)
	}
	updatedSpyCat, err := service.GetById(t.Context(), spyCat.Id)
	if err != nil {
		t.Fatal(err)
	}
	if updatedSpyCat.Salary != newSalary {
		t.Fatal("salary was not updated")
	}
}

func TestGetAll(t *testing.T) {
	repo := memory.NewSpyCatRepository()
	service := NewSpyCatService(repo)
	spyCat1 := &model.SpyCat{
		Name:              "Pickachu",
		YearsOfExperience: 2,
		Breed:             "pokemon",
		Salary:            100,
	}
	err := service.Create(t.Context(), spyCat1)
	if err != nil {
		t.Fatal(err)
	}
	spyCat2 := &model.SpyCat{
		Name:              "Charizard",
		YearsOfExperience: 3,
		Breed:             "pokemon",
		Salary:            200,
	}
	err = service.Create(t.Context(), spyCat2)
	if err != nil {
		t.Fatal(err)
	}

	spyCats, err := service.GetAll(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(spyCats) != 2 {
		t.Fatal("Spy cats were not found")
	}
	if spyCats[0].Id != spyCat1.Id || spyCats[1].Id != spyCat2.Id {
		t.Fatal("Spy cats were not found")
	}
}
