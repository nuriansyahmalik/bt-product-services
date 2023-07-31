package users

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/producer"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type UserService interface {
	Create(requestFormat UserRequestFormat, userId uuid.UUID) (user User, err error)
}

type UserSerivceImpl struct {
	UserRepository UserRepository
	Producer       producer.Producer
	Config         *configs.Config
}

func ProvideUserServiceImpl(userRepository UserRepository, producer producer.Producer, config *configs.Config) *UserSerivceImpl {
	return &UserSerivceImpl{UserRepository: userRepository, Producer: producer, Config: config}
}

func (u *UserSerivceImpl) Create(requestFormat UserRequestFormat, userId uuid.UUID) (user User, err error) {
	user, err = user.NewFromRequestFormat(requestFormat, userId)
	if err != nil {
		return
	}
	if err != nil {
		return user, failure.BadRequest(err)
	}
	err = u.UserRepository.Create(user)
	if err != nil {
		return
	}
	return
}
