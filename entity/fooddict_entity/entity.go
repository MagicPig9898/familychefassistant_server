package fooddict_entity

type TbFoodDict struct {
	ID          int64  `json:"id" db:"id"`
	ClassId     int64  `json:"class_id" db:"class_id"`
	Name        string `json:"name" db:"name"`
	Image       string `json:"image" db:"image"`
	Description string `json:"description" db:"description"`
	Good        string `json:"good" db:"good"`
	Bad         string `json:"bad" db:"bad"`
}

type TbFoodClass struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
