package models

type Books struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}
