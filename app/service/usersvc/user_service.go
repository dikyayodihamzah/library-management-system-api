package usersvc

import (
	"context"

	"github.com/dikyayodihamzah/library-management-api/app/repository/userrepo"
	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/user"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
	"github.com/dikyayodihamzah/library-management-api/pkg/transaction"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type UserService interface {
	Login(c context.Context, api *user.LoginRequest) (*user.LoginResponse, error)
	SignUp(c context.Context, api *user.SignUpRequest) (*user.User, error)

	FindAll(c context.Context, filter *model.QueryParam) ([]user.User, int, error)
	AssignAdmin(c context.Context, id uuid.UUID) (*user.User, error)
}

type userService struct {
	Logger         *zap.SugaredLogger
	Validate       *validator.Validate
	TxManager      transaction.Manager
	UserRepository userrepo.UserRepository
}

func New(
	logger *zap.SugaredLogger,
	validate *validator.Validate,
	txManager transaction.Manager,
	userRepository userrepo.UserRepository,
) UserService {
	return &userService{
		Logger:         logger,
		Validate:       validate,
		TxManager:      txManager,
		UserRepository: userRepository,
	}
}

func (s *userService) FindAll(c context.Context, filter *model.QueryParam) ([]user.User, int, error) {
	// validate filter
	if filter.Sort == "" {
		filter.Sort = "-created_at"
	}

	if _, _, err := query.ValidateSort(filter.Sort, userrepo.SortUserMap); err != nil {
		return nil, 0, exception.ErrorBadRequest(err.Error())
	}

	// get all users data
	users, err := s.UserRepository.FindAll(c, filter)
	if err != nil {
		return nil, 0, exception.ErrorInternal("Failed to get users")
	}

	// get total users data
	total, err := s.UserRepository.Count(c, filter)
	if err != nil {
		return nil, 0, exception.ErrorInternal("Failed to get total users")
	}

	return users, total, nil
}

func (s *userService) AssignAdmin(c context.Context, id uuid.UUID) (*user.User, error) {
	// get user data
	userData, err := s.UserRepository.FindByColumn(c, "id", id)
	if err != nil {
		return nil, exception.ErrorNotFound("User not found")
	}

	if userData.Role == "ADMIN" {
		return nil, exception.ErrorBadRequest("User already has admin role")
	}

	// update user data
	userData.Role = "ADMIN"

	if err := s.TxManager.WithTx(c, func(tx pgx.Tx) error {
		return s.UserRepository.Update(c, tx, userData)
	}); err != nil {
		return nil, exception.ErrorInternal("Failed to assign admin")
	}

	userData.Password = ""
	return userData, nil
}
