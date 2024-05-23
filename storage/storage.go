package storage

import (
	"38hw/storage/postgres"
	"38hw/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Album() repo.AlbumsStorageI
}

type storagePg struct {
	db        *sqlx.DB
	albumRepo repo.AlbumsStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:        db,
		albumRepo: postgres.NewAlbumsrepo(db),
	}
}

func (s storagePg) Album() repo.AlbumsStorageI {
	return s.albumRepo
}
