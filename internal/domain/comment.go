package domain

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID    uint64 `gorm:"not null; foreignkey:UserID" json:"user_id"`
	AuctionID uint64 `gorm:"not null; foreignkey:AuctionID" json:"auction_id"`
	Text      string `gorm:"not null" json:"text"`
	UpVotes   uint64 `gorm:"not null" json:"up_votes"`
	ReplyID   uint64 `json:"reply_id; foreignkey:CommentID"`
}

type CommentUsecase interface {
	Fetch(userID uint64, auctionID uint64) (result []Comment, err error)
	GetByID(id uint64) (result Comment, err error)
	Create(c *Comment, u *User) (errors map[string]string, ok bool)
	Delete(id uint64) (err error)
}

type CommentRepository interface {
	Fetch(userID uint64, auctionID uint64) (result []Comment, err error)
	GetByID(id uint64) (result Comment, err error)
	Create(c *Comment) (err error)
	Delete(id uint64) (err error)
}
