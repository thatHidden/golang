package validation

import (
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"
	"regexp"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var NameRX = regexp.MustCompile("^[a-zA-Z ,.'-]+$")
var UserNameRX = regexp.MustCompile("^[a-zA-Z0-9_]*$")
var PhoneRX = regexp.MustCompile("^\\+\\d{1,2} \\(\\d{3}\\) \\d{3}-\\d{2}-\\d{2}$\n")

type UserValidator struct {
	base *validator.BaseValidator
}

func NewUserValidator(b *validator.BaseValidator) validator.InterfaceValidator[models.UserInputDTO] {
	return &UserValidator{
		base: b,
	}
}

func (uv *UserValidator) validatePhone(phone string) {
	uv.base.Check(phone != "", "phone", "must be provided")
	uv.base.Check(uv.base.Matches(phone, PhoneRX), "phone", "must be a valid phone")
}

func (uv *UserValidator) validateUsername(username string) {
	uv.base.Check(username != "", "username", "must be provided")
	uv.base.Check(len(username) <= 12, "username", "must not be more than 10 bytes long")
	uv.base.Check(len(username) >= 2, "username", "must be at least 2 bytes long")
	uv.base.Check(uv.base.Matches(username, UserNameRX), "username", "must be a valid username")
}

func (uv *UserValidator) validateName(name string) {
	uv.base.Check(name != "", "name", "must be provided")
	uv.base.Check(len(name) <= 10, "name", "must not be more than 10 bytes long")
	uv.base.Check(len(name) >= 2, "name", "must be at least 2 bytes long")
	uv.base.Check(uv.base.Matches(name, NameRX), "name", "must be a valid name")
}

func (uv *UserValidator) validateEmail(email string) {
	uv.base.Check(email != "", "email", "must be provided")
	uv.base.Check(uv.base.Matches(email, EmailRX), "email", "must be a valid email address")
}

func (uv *UserValidator) validatePasswordPlaintext(password string) {
	uv.base.Check(password != "", "password", "must be provided")
	uv.base.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	uv.base.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func (uv *UserValidator) Validate(u *models.UserInputDTO) {
	uv.validateUsername(u.Username)
	uv.validateEmail(u.Email)
	uv.validatePasswordPlaintext(u.Password)
}
