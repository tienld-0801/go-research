package users_model

type User struct {
	UUID     string `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"uuid"`
	Name     string `gorm:"type:varchar(100)" json:"name"`
	Email    string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password string `gorm:"type:varchar(100)" json:"password"`
	Role     int    `json:"role"`
}

type UserResponse struct {
	UUID  string `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"uuid"`
	Name  string `gorm:"type:varchar(100)" json:"name"`
	Email string `gorm:"type:varchar(100);unique_index" json:"email"`
	Role  int    `json:"role"`
}
