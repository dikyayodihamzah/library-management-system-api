package usersvc

import (
	"context"
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/user"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
)

func (s *userService) Login(c context.Context, api *user.LoginRequest) (*user.LoginResponse, error) {
	userRes, err := s.UserRepository.FindByColumn(c, "email", api.Email)
	if userRes == nil || err != nil {
		return nil, exception.ErrorNotFound("We couldn't find an account with that NIK, please check again.")
	}

	if !lib.PasswordCompare(userRes.Password, api.Password) {
		return nil, exception.ErrorBadRequest("That password isn’t right, please check again.")
	}

	token, err := lib.GenerateJwt(&lib.Claims{
		IsAdmin: userRes.Role == "ADMIN",
		StandardClaims: jwt.StandardClaims{
			Issuer: userRes.ID.String(),
		},
	}, *lib.GenID())
	if err != nil {
		return nil, exception.ErrorInternal("Failed When create token")
	}

	refreshJwt, err := lib.GenerateRefreshJwt(userRes.ID.String())
	if err != nil {
		return nil, exception.ErrorInternal("Failed When create token")
	}

	userRes.Password = ""
	resp := user.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshJwt,
		UserData:     userRes,
	}

	return &resp, nil
}

func (s *userService) SignUp(c context.Context, api *user.SignUpRequest) (*user.User, error) {
	// validate request
	if err := s.Validate.Struct(api); err != nil {
		return nil, exception.ErrorBadRequest(err.Error())
	}

	userRes, err := s.UserRepository.FindByColumn(c, "email", api.Email)
	if userRes != nil || err == nil {
		return nil, exception.ErrorBadRequest("Email already exists")
	}

	now := time.Now()
	userRes = &user.User{
		SignUpRequest:    *api,
		Role:             "USER",
		LastActivityDate: now,
	}

	userRes.GenerateID()
	userRes.CreatedAt = &now
	if err := userRes.HashPassword(); err != nil {
		return nil, exception.ErrorInternal("Failed When hash password")
	}

	if err := s.TxManager.WithTx(c, func(tx pgx.Tx) error {
		return s.UserRepository.Create(c, tx, userRes)
	}); err != nil {
		return nil, exception.ErrorInternal("Failed When create user")
	}

	userRes.Password = ""
	return userRes, nil
}
