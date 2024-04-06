package entity

type RadiusRange struct {
	ID        string  `gorm:"primary_key" json:"id"`
	UserID    string  `gorm:"type:varchar(36);index" json:"user_id"`
	Longitude float64 `gorm:"longitude" json:"longitude"`
	Latitude  float64 `gorm:"latitude" json:"latitude"`
}
