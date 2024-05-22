package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"os"
	"time"
)

type Token struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignkey:UserID"`
	Plaintext string    `gorm:"unique; not null"`
	Exp       time.Time `gorm:"not null"`
}

func NewToken(u *User, purpose string, duration time.Duration) (*Token, error) { //в бизнес логику
	expire := time.Now().Add(duration)
	tokenBase := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,
		"exp": expire.Unix(),
		"pur": purpose,
	})
	tokenString, err := tokenBase.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return &Token{}, err
	}

	var token = Token{
		UserID:    u.ID,
		User:      *u,
		Plaintext: tokenString,
		Exp:       expire,
	}

	return &token, nil
}

type TokenRepository interface {
	GetByPlaintext(plaintext string) (result Token, err error)
	GetByUserID(id uint64) (result Token, err error)
	Create(t *Token) (err error)
	Delete(id uint64) (err error)
}
