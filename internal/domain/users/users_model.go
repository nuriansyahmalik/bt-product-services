package users

import (
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type User struct {
	ID        uuid.UUID   `db:"userId"`
	Username  string      `db:"username"`
	Email     string      `db:"email"`
	UserType  string      `db:"userType"`
	CreatedAt time.Time   `db:"createdAt"`
	CreatedBy uuid.UUID   `db:"createdBy"`
	UpdatedAt null.Time   `db:"updatedAt"`
	UpdatedBy nuuid.NUUID `db:"updatedBy"`
	Deleted   null.Time   `db:"deletedAt"`
	DeletedBy nuuid.NUUID `db:"deletedBy"`
}

type UserRequestFormat struct {
	Username string `json:"username" `
	Email    string `json:"email"`
	UserType string `json:"userType"`
}

func (u User) NewFromRequestFormat(req UserRequestFormat, userId uuid.UUID) (newUser User, err error) {
	userId, _ = uuid.NewV4()
	newUser = User{
		ID:        userId,
		Username:  req.Username,
		Email:     req.Email,
		UserType:  req.UserType,
		CreatedAt: time.Now(),
		CreatedBy: userId,
	}
	users := make([]User, 0)
	users = append(users, newUser)
	return
}
