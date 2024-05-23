CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

ALTER TABLE albums
ADD COLUMN genre_id INT,
ADD CONSTRAINT fk_genre
FOREIGN KEY (genre_id)
REFERENCES genres(id);
