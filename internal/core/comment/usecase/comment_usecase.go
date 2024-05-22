package usecase

import (
	"cleanstandarts/internal/core/comment/usecase/validation"
	"cleanstandarts/internal/domain"
	"cleanstandarts/pkg/validator"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
}

func NewCommentUsecase(cr domain.CommentRepository) domain.CommentUsecase {
	return &commentUsecase{
		commentRepo: cr,
	}
}

func (cu *commentUsecase) validateInputFormat(base *validator.BaseValidator, c *domain.Comment) {
	v := validation.NewCommentFormatValidator(base)
	v.Validate(c)
}

func (cu *commentUsecase) validateComment(c *domain.Comment) (errors map[string]string, ok bool) {
	base := validator.NewBaseValidator()
	cu.validateInputFormat(base, c)
	if !base.Valid() {
		return base.Errors, false
	}
	return errors, true
}

func (cu *commentUsecase) Fetch(userID uint64, auctionID uint64) (result []domain.Comment, err error) {
	result, err = cu.commentRepo.Fetch(userID, auctionID)
	return result, err
}

func (cu *commentUsecase) GetByID(id uint64) (result domain.Comment, err error) {
	result, err = cu.commentRepo.GetByID(id)
	return result, err
}

func (cu *commentUsecase) Delete(id uint64) (err error) {
	err = cu.commentRepo.Delete(id)
	return err
}

func (cu *commentUsecase) Create(c *domain.Comment, u *domain.User) (errors map[string]string, ok bool) {
	errors, ok = cu.validateComment(c)
	if !ok {
		return errors, false
	}
	errors = map[string]string{}
	c.UserID = uint64(u.ID)
	err := cu.commentRepo.Create(c)
	if err != nil {
		errors["internal"] = err.Error()
	}
	return errors, true
}
