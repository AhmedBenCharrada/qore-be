package person

import (
	"context"
	"errors"
	"fmt"
	"qore-be/internal/domain/dto"
	"qore-be/internal/domain/entities"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = fmt.Errorf("not found")
)

// Repo represents the person repository interface.
type Repo struct {
	db *gorm.DB
}

// NewRepository create a new instance of the person repository.
func NewRepository(db *gorm.DB) *Repo {
	if db == nil {
		panic("nil db")
	}

	if err := db.AutoMigrate(&entities.Person{}, &entities.Phone{}, &entities.Address{}, &entities.PersonAddress{}); err != nil {
		panic(err)
	}

	return &Repo{
		db: db,
	}
}

// Add saves new user to the database.
func (r *Repo) Add(ctx context.Context, d dto.PersonDTO) (dto.PersonDTO, error) {
	pr := entities.Person{
		Age:  d.Age,
		Name: d.Name,
	}

	ph := &entities.Phone{
		Number: d.Number,
	}

	addr := entities.Address{
		City:    d.City,
		State:   d.State,
		Street1: d.Street1,
		Street2: d.Street2,
		Zip:     d.Zip,
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&pr); res != nil && res.Error != nil {
			return res.Error
		}

		ph.PersonID = pr.ID
		if res := tx.Create(&ph); res != nil && res.Error != nil {
			return res.Error
		}

		if res := tx.Create(&addr); res != nil && res.Error != nil {
			return res.Error
		}

		if res := tx.Create(&entities.PersonAddress{
			PersonID:  pr.ID,
			AddressID: addr.ID,
		}); res != nil && res.Error != nil {
			return res.Error
		}

		return nil
	})

	return dto.PersonDTO{
		Name:    pr.Name,
		Age:     pr.Age,
		Number:  ph.Number,
		State:   addr.State,
		City:    addr.City,
		Street1: addr.Street1,
		Street2: addr.Street2,
		Zip:     addr.Zip,
	}, err
}

// GetByID retrieves a person data by its ID.
func (r *Repo) GetByID(ctx context.Context, id int) (dto.PersonDTO, error) {
	pr := &entities.Person{}
	tx := r.db.WithContext(ctx).Table(pr.TableName()).First(&pr, "id= ?", id)
	if tx != nil && tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return dto.PersonDTO{}, ErrRecordNotFound
		}
		return dto.PersonDTO{}, fmt.Errorf("failed to get person data: %v", tx.Error)
	}

	ph := &entities.Phone{}
	tx = r.db.WithContext(ctx).Table(ph.TableName()).First(&ph, "person_id= ?", id)
	if tx != nil && tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return dto.PersonDTO{}, fmt.Errorf("failed to get phone data: %v", tx.Error)
	}

	pAddr := &entities.PersonAddress{}
	tx = r.db.WithContext(ctx).Table(pAddr.TableName()).First(&pAddr, "person_id= ?", id)
	if tx != nil && tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return dto.PersonDTO{}, fmt.Errorf("failed to get address_join data: %v", tx.Error)
	}

	addr := &entities.Address{}
	tx = r.db.WithContext(ctx).Table(addr.TableName()).First(&addr, "id= ?", pAddr.AddressID)
	if tx != nil && tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return dto.PersonDTO{}, fmt.Errorf("failed to get address: %v", tx.Error)
	}

	return dto.PersonDTO{
		Name:    pr.Name,
		Age:     pr.Age,
		Number:  ph.Number,
		State:   addr.State,
		City:    addr.City,
		Street1: addr.Street1,
		Street2: addr.Street2,
		Zip:     addr.Zip,
	}, nil
}

// GetAll retrieves person rows.
func (r *Repo) GetAll(ctx context.Context, offset int, limit int) ([]dto.PersonDTO, error) {
	persons := []entities.Person{}
	tx := r.db.WithContext(ctx).Table((entities.Person{}).TableName()).Offset(offset).Limit(limit).Find(&persons)
	if tx != nil && tx.Error != nil {
		return nil, fmt.Errorf("failed to get person data: %v", tx.Error)
	}

	dtos := make([]dto.PersonDTO, len(persons))
	fmt.Println(len(persons))
	for i, p := range persons {
		dtos[i] = dto.PersonDTO{
			ID:   p.ID,
			Name: p.Name,
			Age:  p.Age,
		}
	}
	return dtos, nil
}
