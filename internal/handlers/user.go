package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/users"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type UserHandler struct {
	UserService users.UserService
}

func ProvideUserHandler(UserService users.UserService) UserHandler {
	return UserHandler{UserService: UserService}
}

func (h *UserHandler) Router(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/", h.CreateUser)
	})
}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat users.UserRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, _ := uuid.NewV4() // TODO: read from context

	foo, err := h.UserService.Create(requestFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, foo)
}
