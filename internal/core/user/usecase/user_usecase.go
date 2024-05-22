package usecase

import (
	"cleanstandarts/internal/core/user/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto"
	"cleanstandarts/internal/domain/dto/models"
	"cleanstandarts/pkg/validator"

	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
}

func NewUserUsecase(ur domain.UserRepository, tr domain.TokenRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo:  ur,
		tokenRepo: tr,
	}
}

func (uu *userUsecase) Fetch() (result []models.UserOutputDTO, err error) {
	raw, err := uu.userRepo.Fetch()
	if err != nil {
		return nil, err
	}
	for _, user := range raw {
		DTO := dto.UserToOutputDto(&user)
		result = append(result, DTO)
	}
	return result, err
}

func (uu *userUsecase) validateInputFormat(base *validator.BaseValidator, u *models.UserInputDTO) {
	v := validation.NewUserValidator(base)
	v.Validate(u)
}

//func (uu *userUsecase) validateBusinessRules(b validator.BaseValidator,u *domain.User) {
//
//}

func (uu *userUsecase) ValidateUser(u *models.UserInputDTO) (errors map[string]string, ok bool) {
	ok = true
	base := validator.NewBaseValidator()
	uu.validateInputFormat(base, u)
	if !base.Valid() {
		errors = base.Errors
		ok = false
	}
	return errors, ok
}

func (uu *userUsecase) Activate(authUser domain.User, tokenRequest string) (err error) {
	requestTokenClaims, err := ParseJWT(tokenRequest)
	if err != nil {
		return err
	}

	findToken, err := uu.tokenRepo.GetByPlaintext(tokenRequest)
	if findToken.ID == 0 {
		return errors.New("err no token")
	}

	if requestTokenClaims["sub"] != float64(authUser.ID) {
		return errors.New("tokens not eq")
	}

	err = uu.tokenRepo.Delete(uint64(findToken.ID))
	if err != nil {
		return err
	}

	idFloat, ok := requestTokenClaims["sub"].(float64)
	if !ok {
		return errors.New("unable to convert 'sub' to float64")
	}

	id := uint64(idFloat)

	err = uu.userRepo.Activate(id)

	return nil
}

func (uu *userUsecase) Auth(email string, password string) (result string, err error) {
	var user domain.User
	user, err = uu.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("bad password or email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("server error")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

	if err != nil {
		return "", errors.New("error to create a token")
	}

	return tokenString, nil
}

func (uu *userUsecase) GetByID(id uint64) (result models.UserOutputDTO, err error) {
	raw, err := uu.userRepo.GetByID(id)
	if err != nil {
		return result, err
	}
	result = dto.UserToOutputDto(&raw)
	return result, err
}

func (uu *userUsecase) GetByIDRaw(id uint64) (result domain.User, err error) {
	result, err = uu.userRepo.GetByID(id)
	return result, err
}

func (uu *userUsecase) Create(u *models.UserInputDTO) (errors map[string]string, ok bool) {
	errors, ok = uu.ValidateUser(u)
	if !ok {
		return errors, false
	}

	errors = make(map[string]string)

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		errors["internal"] = err.Error()
		return errors, false
	}

	user := dto.UserDtoToDomain(u)
	user.Password = string(hash)

	err = uu.userRepo.Create(&user)
	if err != nil {
		errors["internal"] = err.Error()
		return errors, false
	}

	token, err := domain.NewToken(&user, "activation", time.Minute*5)

	err = uu.tokenRepo.Create(token)
	if err != nil {
		errors["internal"] = err.Error()
		return errors, false
	}

	return errors, true
}

func (uu *userUsecase) Delete(id uint64) (err error) {
	err = uu.userRepo.Delete(id)
	return err
}

//func (uu *userUsecase) Update(u *domain.User) (err error) {
//	err = uu.userRepo.Update(u)
//	return err
//}

func ParseJWT(tokenPlaintext string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenPlaintext, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return jwt.MapClaims{}, errors.New("unable to read token claims")
	}

	return claims, nil
}
