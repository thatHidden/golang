package usecase

import (
	"cleanstandarts/internal/domain"
	"cleanstandarts/internal/domain/dto/models"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

type imageUsecase struct {
	imageRepo domain.ImageRepository
}

func NewCommentUsecase(ir domain.ImageRepository) domain.ImageUsecase {
	return &imageUsecase{
		imageRepo: ir,
	}
}

// toto
func (iu *imageUsecase) Fetch(userID uint64, auctionID uint64) (result []domain.Image, err error) {
	//result, err = iu.imageRepo.Fetch(userID, auctionID)
	return result, err
}

func (iu *imageUsecase) GetByID(id uint64) (result domain.Image, err error) {
	result, err = iu.imageRepo.GetByID(id)
	return result, err
}

func (iu *imageUsecase) Delete(id uint64) (err error) {
	err = iu.imageRepo.Delete(id)
	return err
}

func (iu *imageUsecase) Create(dto *models.ImagesInputDTO) (errors map[string]string, ok bool) {
	//errors, ok = iu.validateImage(dto)
	//if !ok {
	//	return errors, false
	//}
	errors = map[string]string{}

	savePath := "static/" + "auction/" + strconv.FormatUint(dto.AucID, 10) + "/"

	//savePath := filepath.Join("static", "auction", strconv.FormatUint(dto.AucID, 10)) //.env

	images := make([]domain.Image, 0, len(dto.Interior)+len(dto.Exterior))

	for _, image := range dto.Interior {
		path := filepath.Join(savePath, "interior")
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			errors["internal"] = err.Error()
			return errors, false
		}
		images = append(images, domain.Image{
			Mark:      "interior",
			AuctionID: dto.AucID,
			Path:      savePath + "interior/" + image.Filename,
		})
		err = iu.saveImage(filepath.Join(path, filepath.Base(image.Filename)), image)
		if err != nil {
			errors["internal"] = err.Error()
			return errors, false
		}
	}

	for _, image := range dto.Exterior {
		path := filepath.Join(savePath, "exterior")
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			errors["internal"] = err.Error()
			return errors, false
		}
		images = append(images, domain.Image{
			Mark:      "exterior",
			AuctionID: dto.AucID,
			Path:      savePath + "exterior/" + image.Filename,
		})
		err = iu.saveImage(filepath.Join(path, filepath.Base(image.Filename)), image)
		if err != nil {
			errors["internal"] = err.Error()
			return errors, false
		}
	}

	err := iu.imageRepo.MultipleCreate(images)
	if err != nil {
		errors["internal"] = err.Error()
	}
	return errors, true
}

func (iu *imageUsecase) saveImage(path string, image *multipart.FileHeader) error {
	file, err := image.Open()
	if err != nil {
		return err
	}

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	defer file.Close()

	return nil
}
