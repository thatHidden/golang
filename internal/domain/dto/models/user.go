package models

type UserInputDTO struct {
	Username string
	Email    string
	Password string
}

type UserOutputDTO struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
	//Email       string `json:"email"`
	//Phone       string `json:"phone"`
	Name string `json:"name"`
	//IsActivated bool   `json:"is_activated"`
	//UserComments
	//UserBids
}

//func NewUserOutputDTO(u *domain.User) UserOutputDTO {
//	return UserOutputDTO{
//		Username:    u.Username,
//		Photo:       u.Photo,
//		Email:       u.Email,
//		Phone:       u.Phone,
//		Name:        u.Name,
//		IsActivated: u.IsActivated,
//	}
//}
