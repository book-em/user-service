package domain

type UserRole string

const (
	Guest UserRole = "guest"
	Host  UserRole = "host"
	Admin UserRole = "admin"
)

type User struct {
	ID       uint     `                gorm:"primaryKey"`
	Username string   `json:"username" gorm:"type:varchar(30);not null;uniqueIndex"`
	Password string   `json:"password" gorm:"type:varchar(100);not null;"`
	Email    string   `json:"email"    gorm:"type:varchar(60);not null;uniqueIndex"`
	Name     string   `json:"name"     gorm:"type:varchar(60);not null;"`
	Surname  string   `json:"surname"  gorm:"type:varchar(60);not null;"`
	Role     UserRole `json:"role"     gorm:"type:varchar(5);not null;check:role IN ('guest','host','admin')"`
	Address  string   `json:"address"  gorm:"type:varchar(150);not null"`
}
