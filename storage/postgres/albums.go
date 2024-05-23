package postgres

import (
	"context"
	"fmt"
	m "38hw/models"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type AlbumRepo struct {
	DB *sqlx.DB
}

func NewAlbumsrepo(db *sqlx.DB) *AlbumRepo {
	return &AlbumRepo{
		DB: db,
	}
}

func (a *AlbumRepo) CreateAlbum(ctx context.Context, alb m.Album) error {
	genre := m.Genre{}
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	query := `
		SELECT id FROM genres WHERE name = $1	
	`
	row := tx.QueryRowContext(ctx, query, strings.ToLower(alb.Genre))

	err = row.Scan(&genre.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `
		INSERT INTO albums (
			title, 
			artist, 
			price,
			genre_id) 
		VALUES(	$1, $2, $3, $4)
		RETURNING
			id,         
			title,       
			artist,      
			price
	`

	_, err = tx.ExecContext(ctx, query, strings.ToLower(alb.Title), strings.ToLower(alb.Artist), alb.Price, genre.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (a *AlbumRepo) GetAlbumsById(ctx context.Context, id string) (m.Album, error) {
	query := `
		SELECT albums.title, albums.artist, albums.price, genres.name AS genre
		FROM albums
		INNER JOIN genres ON albums.genre_id = genres.id
		WHERE albums.id = $1;
	`
	row := a.DB.QueryRowContext(ctx, query, id)

	album := m.Album{}

	err := row.Scan(&album.Title, &album.Artist, &album.Price, &album.Genre)
	if err != nil {
		return album, err
	}

	return album, nil
}

func (a *AlbumRepo) GetAlbums(ctx context.Context) (albums []m.Album, err error) {
	query := `
		SELECT albums.title, albums.artist, albums.price, genres.name AS genre_name
		FROM albums
		INNER JOIN genres ON albums.genre_id = genres.id;
	`
	row, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return albums, err
	}

	for row.Next() {
		album := m.Album{}
		err := row.Scan(&album.Title, &album.Artist, &album.Price, &album.Genre)
		if err != nil {
			log.Println("Error scanning sql query: ", err)
			return albums, err
		}
		albums = append(albums, album)
	}
	return albums, nil
}

func (a *AlbumRepo) UpdateAlbumById(ctx context.Context, alb m.Album, id string) (m.Album, error) {
	album := m.Album{}
	genre := m.Genre{}
	tx, err := a.DB.Begin()
	if err != nil {
		return album, err
	}

	query := `
		SELECT id FROM genres WHERE name = $1	
	`
	row := tx.QueryRowContext(ctx, query, strings.ToLower(alb.Genre))

	err = row.Scan(&genre.Id)
	if err != nil {
		fmt.Println("Eror 1")
		tx.Rollback()
		return album, err
	}

	query = `
		UPDATE albums
		SET title = $1, artist = $2, price = $3, genre_id = $4                                     
		WHERE id = $5
		RETURNING
			id,         
			title,       
			artist,      
			price,
			(SELECT name FROM genres WHERE id = albums.genre_id) AS genre;
	`

	row = tx.QueryRowContext(ctx, query, strings.ToLower(alb.Title), strings.ToLower(alb.Artist), alb.Price, genre.Id, id)

	err = row.Scan(&album.Id, &album.Title, &album.Artist, &album.Price, &album.Genre)
	if err != nil {
		fmt.Println("Eror 2")
		tx.Rollback()
		return album, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Eror 3")
		return album, err
	}

	return album, nil
}

func (a *AlbumRepo) GetAlbumsByTitle(ctx context.Context, title string) (albums []m.Album, err error) {
	query := `
		SELECT albums.title, albums.artist, albums.price, genres.name AS genre
		FROM albums
		INNER JOIN genres ON albums.genre_id = genres.id
		WHERE albums.title = $1;	
    `
	rows, err := a.DB.QueryContext(ctx, query, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var album m.Album
		err := rows.Scan(&album.Title, &album.Artist, &album.Price, &album.Genre)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (a *AlbumRepo) GetAlbumsByArtist(ctx context.Context, artist string) (albums []m.Album, err error) {
	query := `
		SELECT albums.title, albums.artist, albums.price, genres.name AS genre
		FROM albums
		INNER JOIN genres ON albums.genre_id = genres.id
		WHERE albums.artist = $1;	
    `
	rows, err := a.DB.QueryContext(ctx, query, artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var album m.Album
		err := rows.Scan(&album.Title, &album.Artist, &album.Price, &album.Genre)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (a *AlbumRepo) GetAlbumsByPrice(ctx context.Context, price float64) (albums []m.Album, err error) {
	query := `
		SELECT albums.title, albums.artist, albums.price, genres.name AS genre
		FROM albums
		INNER JOIN genres ON albums.genre_id = genres.id
		WHERE albums.price = $1;
    `
	rows, err := a.DB.QueryContext(ctx, query, price)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var album m.Album
		err := rows.Scan(&album.Title, &album.Artist, &album.Price, &album.Genre)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (a *AlbumRepo) GetAlbumsByGenre(ctx context.Context, genre string) (albums []m.Album, err error) {
	query := `
		SELECT albums.title, albums.artist, albums.price, genres.name AS genre
		FROM albums
		INNER JOIN genres ON albums.genre_id = genres.id
		WHERE genres.name = $1;	
    `
	rows, err := a.DB.QueryContext(ctx, query, genre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var album m.Album
		err := rows.Scan(&album.Title, &album.Artist, &album.Price, &album.Genre)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (a *AlbumRepo) DeletAlbumsById(ctx context.Context, id string) error {
	query := `
	DELETE FROM albums WHERE id = $1;
	`
	_, err := a.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}
