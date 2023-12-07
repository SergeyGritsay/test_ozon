package postgres

import (
	"database/sql"
	"fmt"
	"github.com/tetovske/advertisement-service/pkg/models"
)

const advertisementsTableName = "ADVERTISEMENTS"

type AdvertisementPSQL struct {
	conn *sql.DB
}

func NewAdvertisementPSQL(conn *sql.DB) *AdvertisementPSQL {
	return &AdvertisementPSQL{conn: conn}
}

func (r *AdvertisementPSQL) CreateAdvertisement(ad models.Advertisement) (int, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description, price) values ($1, $2, $3) RETURNING id", advertisementsTableName)

	row := tx.QueryRow(createItemQuery, ad.Title, ad.Description, ad.Price)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *AdvertisementPSQL) GetAdvertisement(id int) (models.Advertisement, error) {
	var ad models.Advertisement

	query := fmt.Sprintf(`SELECT id, title, description, price FROM %s WHERE id = $1`, advertisementsTableName)
	row := r.conn.QueryRow(query, id)
	err := row.Scan(&ad.Id, &ad.Title, &ad.Description, &ad.Price)
	if err != nil {
		return models.Advertisement{}, err
	}

	return ad, nil
}

func (r *AdvertisementPSQL) GetAdvertisementList(sort string) ([]models.Advertisement, error) {
	var ads []models.Advertisement

	query := fmt.Sprintf(`SELECT id, title, description, price FROM %s`, advertisementsTableName)

	sortQuery := " "

	if len(sort) > 0 {
		if string(sort[0]) == "-" {
			sortQuery += "ORDER BY " + sort[1:] + " DESC"
		} else if string(sort[0]) == "+" {
			sortQuery += "ORDER BY " + sort[1:] + " ASC"
		}
	}

	rows, err := r.conn.Query(query + sortQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		ad := models.Advertisement{}
		if err = rows.Scan(&ad.Id, &ad.Title, &ad.Description, &ad.Price); err != nil {
			return nil, err
		}

		ads = append(ads, ad)
	}

	return ads, nil
}
