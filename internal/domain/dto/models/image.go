package models

import "mime/multipart"

type ImagesInputDTO struct {
	Interior []*multipart.FileHeader
	Exterior []*multipart.FileHeader
	AucID    uint64
}
