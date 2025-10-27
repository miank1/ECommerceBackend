package service

import (
	jwtutil "ecommerce-backend/pkg/jwt"
	"ecommerce-backend/services/userservice/internal/model"
	"ecommerce-backend/services/userservice/internal/repository"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo      *repository.UserRepository
	JWTTTL    int // minutes
	JWTSecret string
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{Repo: r}
}

// Login validates credentials and returns JWT + user (without password)
func (s *UserService) Login(email, password string) (string, *model.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// generate token
	token, err := jwtutil.GenerateToken(s.JWTSecret, user.ID, s.JWTTTL)
	if err != nil {
		return "", nil, err
	}

	user.PasswordHash = ""
	return token, user, nil
}

func (s *UserService) Register(name, email, password string) (*model.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email and password are required")
	}
	// check email uniqueness
	existing, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.Repo.Create(u); err != nil {
		return nil, err
	}

	// don't return password hash
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("🔐 UserService secret:", os.Getenv("JWT_SECRET"))

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *UserService) Authenticate(email, password string) (*model.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
