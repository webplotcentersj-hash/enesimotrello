package service

import (
	"errors"
	"time"
	"task-board/internal/domain"
	"task-board/internal/repository"
	"task-board/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	Register(email, username, password, firstName, lastName string) (*domain.User, error)
	Login(email, password string) (string, *domain.User, error)
	GetProfile(userID uint) (*domain.User, error)
	UpdateProfile(userID uint, firstName, lastName string) (*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) SetConfig(config *config.Config) {
	s.config = config
}

func (s *userService) Register(email, username, password, firstName, lastName string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUser, _ = s.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Create new user
	user := &domain.User{
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	}

	// Set password
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(email, password string) (string, *domain.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(password) {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *userService) GetProfile(userID uint) (*domain.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *userService) UpdateProfile(userID uint, firstName, lastName string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) generateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(s.config.JWTExpiry)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}
