package services

import (
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/web"
	"WedhaWS/utsgosmt5/alumni-tracer/app/repository"
	"WedhaWS/utsgosmt5/alumni-tracer/helper"
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2/log"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	Create(C *fiber.Ctx) error
	Update(C *fiber.Ctx) error
	Delete(C *fiber.Ctx) error
	FindById(C *fiber.Ctx) error
	FindAll(C *fiber.Ctx) error
	Login(C *fiber.Ctx) error
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	validate       validator.Validate
}

func NewUserService(UserRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: UserRepository,
		DB:             DB,
		validate:       *validate,
	}
}

func (service *UserServiceImpl) Login(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var request web.LoginRequest
	if err := C.BodyParser(&request); err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := service.validate.Struct(request)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	user := domain.Users{
		Email:         request.Email,
		Password_Hash: request.Password,
	}
	user, err = service.UserRepository.FindByEmailPassword(ctx, tx, user)
	if err != nil {
		return C.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	if !helper.CheckPassword(request.Password, user.Password_Hash) {
		return C.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Credentials"})
	}

	token, err := helper.GenerateToken(user)
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	response := web.LoginResponse{
		User:  user,
		Token: token,
	}
	log.Info(C.IP() + " " + response.User.Username + " Has Login")

	return C.Status(fiber.StatusOK).JSON(response)
}

func (service *UserServiceImpl) Create(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var request web.UserRequest
	if err := C.BodyParser(&request); err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := service.validate.Struct(request)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	password, err := helper.HashPassword(request.Password_Hash)
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user := domain.Users{
		Username:      request.Username,
		Email:         request.Email,
		Password_Hash: password,
		Role:          request.Role,
	}
	user = service.UserRepository.Save(ctx, tx, user)
	return C.Status(fiber.StatusCreated).JSON(user)
}

func (service *UserServiceImpl) Update(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var request web.UserUpdateRequest
	if err := C.BodyParser(&request); err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	request.Id = id
	err = service.validate.Struct(request)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	password, err := helper.HashPassword(request.Password_Hash)
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user := domain.Users{
		Id:            request.Id,
		Username:      request.Name,
		Email:         request.Email,
		Password_Hash: password,
		Role:          request.Role,
	}
	user = service.UserRepository.Update(ctx, tx, user)
	return C.Status(fiber.StatusCreated).JSON(user)
}

func (service *UserServiceImpl) Delete(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := strconv.ParseInt(C.Params("id"), 10, 64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	UserId := int64(id)
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	service.UserRepository.Delete(ctx, tx, UserId)
	return C.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "user deleted",
	})	
}

func (service *UserServiceImpl) FindById(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := strconv.ParseInt(C.Params("id"), 10, 64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	user, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return C.Status(fiber.StatusOK).JSON(user)
}

func (service *UserServiceImpl) FindAll(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	page, _ := strconv.Atoi(C.Query("page", "1"))
	limit, _ := strconv.Atoi(C.Query("limit", "10"))
	sortBy := C.Query("sortBy", "id")
	order := C.Query("order", "asc")
	search := C.Query("search", "")
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	Users := service.UserRepository.FindAll(ctx, tx, search, sortBy, order, limit, offset)

	var UserResponses []web.UserResponse
	for _, user := range Users {
		UserResponses = append(UserResponses, helper.ToUserResponse(user))
	}
	return C.Status(fiber.StatusOK).JSON(UserResponses)
}
