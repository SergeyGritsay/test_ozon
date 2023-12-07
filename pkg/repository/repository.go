package repository

import (
	"advertisement-service/pkg/models"
	"advertisement-service/pkg/repository/postgres"
	"database/sql"
)

type AdvertisementRepository interface {
	CreateAdvertisement(ad models.Advertisement) (int, error)
	GetAdvertisement(id int) (models.Advertisement, error)
	GetAdvertisementList(sort string) ([]models.Advertisement, error)
}

type PhotoRepository interface {
	CreatePhoto(pic models.Photo, adId int) (int, error)
	GetPhoto(id int) (models.Photo, error)
	GetPhotoList(id int) ([]models.Photo, error)
}

type Repository struct {
	AdvertisementRepository
	PhotoRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AdvertisementRepository: postgres.NewAdvertisementPSQL(db),
		PhotoRepository:         postgres.NewPhotoPSQL(db),
	}
}
