package http

import (
	"homework-dontpanicw/app/api/http/types"
	"homework-dontpanicw/app/repository"
	"homework-dontpanicw/app/usecases"
	"homework-dontpanicw/app/usecases/auth"
	"net/http"
)

type User struct {
	userService usecases.User
}

func NewUserHandler(service usecases.User) *User {
	return &User{
		userService: service,
	}
}

// @Summary Registration
// @Description Создает нового пользователя в системе.
// @Tags users
// @Accept json
// @Produce json
// @Param request body types.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} types.RegisterResponse "Пользователь успешно зарегистрирован"
// @Failure 400 {object} types.ErrorResponse "Некорректные данные"
// @Failure 500 {object} types.ErrorResponse "Ошибка сервера"
// @Router /user/register [post]
func (u *User) registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := types.CreateUserRequest(r)
	if err != nil {
		types.UserProcessError(w, repository.BadRequest, nil, 0)
		return
	}

	err = u.userService.RegisterNewUser(ctx, user.Username, user.Password)
	if err != nil {
		types.UserProcessError(w, err, nil, 0)
		return
	}
	types.UserProcessError(w, nil, types.RegisterResponse{Result: "user registered"}, http.StatusCreated)
}

// @Summary Authorization
// @Description Login user and create token
// @Tags users
// @Accept json
// @Produce json
// @Param request body types.RegisterRequest true "Данные для входа"
// @Success 200 {object} types.LoginResponse "Успешная авторизация"
// @Failure 400 {object} types.ErrorResponse "Некорректные данные"
// @Failure 401 {object} types.ErrorResponse "Неверный логин или пароль"
// @Failure 500 {object} types.ErrorResponse "Ошибка сервера"
// @Router /user/login [post]
func (u *User) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := types.CreateUserRequest(r)
	if err != nil {
		types.UserProcessError(w, repository.BadRequest, nil, 0)
		return
	}

	var storedPassword string
	storedUser, err := u.userService.GetUserByUsername(ctx, user.Username)
	if err != nil {
		types.UserProcessError(w, repository.UserNotFound, nil, 0)
		return
	}

	storedId := storedUser.Id
	storedPassword = storedUser.Password

	if !auth.CheckPassword(user.Password, storedPassword) {
		types.UserProcessError(w, repository.Unauthorized, nil, 0)
		return
	}

	sessionId, err := u.userService.CreateNewSession(storedId)
	if err != nil {
		types.UserProcessError(w, err, nil, 0)
		return
	}
	types.UserProcessError(w, nil, types.LoginResponse{SessionId: sessionId}, http.StatusCreated)
}
