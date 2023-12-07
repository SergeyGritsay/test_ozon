package postgres

import (
	"database/sql"
	"fmt"
	"github.com/tetovske/advertisement-service/pkg/models"
)

const photosTableName = "PHOTOS"

type PhotoPSQL struct {
	conn *sql.DB
}

func NewPhotoPSQL(conn *sql.DB) *PhotoPSQL {
	return &PhotoPSQL{conn: conn}
}

func (r *PhotoPSQL) CreatePhoto(pic models.Photo, adId int) (int, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (link, tag, advertisementid) values ($1, $2, $3) RETURNING id", photosTableName)

	row := tx.QueryRow(createItemQuery, pic.Link, pic.Tag, adId)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *PhotoPSQL) GetPhoto(id int) (models.Photo, error) {
	var photo models.Photo

	query := fmt.Sprintf(`SELECT id, link, tag FROM %s WHERE id = $1 ORDER BY tag ASC`, photosTableName)
	row := r.conn.QueryRow(query, id)
	err := row.Scan(&photo.Id, &photo.Link, &photo.Tag)
	if err != nil {
		return models.Photo{}, err
	}

	return photo, nil
}

func (r *PhotoPSQL) GetPhotoList(id int) ([]models.Photo, error) {
	var pics []models.Photo

	query := fmt.Sprintf(`SELECT id, link, tag FROM %s WHERE advertisementid = $1 ORDER BY tag ASC`, photosTableName)
	rows, err := r.conn.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		photo := models.Photo{}
		if err = rows.Scan(&photo.Id, &photo.Link, &photo.Tag); err != nil {
			return nil, err
		}

		pics = append(pics, photo)
	}

	return pics, nil
}
