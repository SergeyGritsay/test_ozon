package services

import (
	"advertisement-service/pkg/models"
	"advertisement-service/pkg/repository"
	"errors"
	"log"
)

type AdvertisementService struct {
	repo *repository.Repository
}

func NewAdvertisementService(repo *repository.Repository) *AdvertisementService {
	return &AdvertisementService{repo: repo}
}

func (r *AdvertisementService) CreateAdvertisement(ad models.Advertisement) (int, error) {
	if !ad.ValidateFields() {
		return 0, errors.New("invalid fields")
	}

	adId, err := r.repo.AdvertisementRepository.CreateAdvertisement(ad)

	for i, pic := range ad.Photos {
		var tag uint

		if i == 0 {
			tag = 0
		} else {
			tag = 1
		}

		if _, err = r.repo.PhotoRepository.CreatePhoto(models.Photo{
			Link: pic,
			Tag:  tag,
		}, adId); err != nil {
			log.Fatal(err)

			return 0, err
		}
	}

	return adId, err
}

func (r *AdvertisementService) GetAdvertisement(id int, fields []string) (map[string]interface{}, error) {
	descPresence, photosPresence := false, false
	for _, field := range fields {
		if field == "photos" {
			photosPresence = true
		} else if field == "description" {
			descPresence = true
		}
	}

	ad, err := r.repo.AdvertisementRepository.GetAdvertisement(id)
	if err != nil {
		return map[string]interface{}{}, err
	}

	photos, err := r.repo.PhotoRepository.GetPhotoList(id)
	if err != nil {
		return map[string]interface{}{}, err
	}

	resp := map[string]interface{}{
		"title": ad.Title,
		"price": ad.Price,
	}

	if photosPresence {
		var temp []string

		for _, photo := range photos {
			temp = append(temp, photo.Link)
		}

		resp["photos"] = temp
	} else {
		resp["photos"] = photos[0].Link
	}

	if descPresence {
		resp["description"] = ad.Description
	}

	return resp, err
}

func (r *AdvertisementService) GetAdvertisements(sort string) ([]map[string]interface{}, error) {
	ads, err := r.repo.AdvertisementRepository.GetAdvertisementList(sort)
	if err != nil {
		return nil, err
	}

	var resp []map[string]interface{}

	for _, ad := range ads {
		photo, err := r.repo.PhotoRepository.GetPhoto(int(ad.Id))
		if err != nil {
			return []map[string]interface{}{}, err
		}

		resp = append(resp, map[string]interface{}{
			"title":  ad.Title,
			"price":  ad.Price,
			"photos": photo.Link,
		})
	}

	return resp, err
}
