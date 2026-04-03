package fooddict_entity

type TbFoodDict struct {
	ID          int64  `json:"id" db:"id,auto_increment"`
	ClassId     int64  `json:"class_id" db:"class_id"`
	Name        string `json:"name" db:"name"`
	Image       string `json:"image" db:"image"`
	Description string `json:"description" db:"description"`
}

type TbFoodClass struct {
	ID   int64  `json:"id" db:"id,auto_increment"`
	Name string `json:"name" db:"name"`
}
