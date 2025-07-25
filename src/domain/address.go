package domain

type Address struct {
	ID      uint   `                gorm:"primaryKey"`
	Street  string `json:"street"   gorm:"type:varchar(80);not null"`
	City    string `json:"city"     gorm:"type:varchar(60);not null"`
	State   string `json:"state"    gorm:"type:varchar(60)"`
	Zip     string `json:"zip"      gorm:"type:varchar(20)"`
	Country string `json:"country"  gorm:"type:varchar(80);not null"`
}
