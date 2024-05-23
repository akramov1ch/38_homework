package models

type Album struct {
	Id     string  `json:"-"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
	Genre  string  `json:"genre"`
}

type Genre struct {
	Id   int    `json:"-"`
	Name string `json:"name"`
}
