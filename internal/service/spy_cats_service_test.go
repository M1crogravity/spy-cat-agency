package service

import (
	"errors"
	"testing"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/memory"
)

func TestSpyCatsCreate(t *testing.T) {
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
	err = service.Create(t.Context(), spyCat)
	if !errors.Is(err, storage.ErrorUniqueConstraintViolation) {
		t.Fatal("unique name constaint violation")
	}
}

func TestSpyCatsGetById(t *testing.T) {
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

func TestSpyCatsRemove(t *testing.T) {
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

	err = service.Remove(t.Context(), spyCat.Id)
	if !errors.Is(err, storage.ErrorModelNotFound) {
		t.Fatal("Expected error to be storage.ErrorModelNotFound")
	}
}

func TestSpyCatsUpdateSalary(t *testing.T) {
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
	_, err = service.UpdateSalary(t.Context(), spyCat.Id, newSalary)
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

func TestSpyCatsGetAll(t *testing.T) {
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
	ids := map[int64]struct{}{
		spyCat1.Id: {},
		spyCat2.Id: {},
	}
	for _, spyCat := range spyCats {
		delete(ids, spyCat.Id)
	}
	if len(ids) != 0 {
		t.Fatal("Spy cats were not found")
	}
}
