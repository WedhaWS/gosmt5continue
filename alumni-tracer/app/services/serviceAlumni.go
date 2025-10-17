package services

import (
	"WedhaWS/utsgosmt5/alumni-tracer/helper"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/web"
	"WedhaWS/utsgosmt5/alumni-tracer/app/repository"
	"context"
	"database/sql"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"time"
	"strconv"
)

type AlumniService interface{
	Create(C *fiber.Ctx) error
	Update(C *fiber.Ctx) error
	Delete(C *fiber.Ctx) error
	FindById(C *fiber.Ctx) error
	FindAll(C *fiber.Ctx) error
	FindPengangguran(C *fiber.Ctx) error
}

type AlumniServiceImpl struct{
	AlumniRepository repository.AlumniRepository
	DB *sql.DB
	validate *validator.Validate
}

func NewAlumniService(AlumniRepository repository.AlumniRepository,DB *sql.DB,validate *validator.Validate) AlumniService{
	return &AlumniServiceImpl{
		AlumniRepository: AlumniRepository,
		DB: DB,
		validate: validate,
	}
}

func (service *AlumniServiceImpl) Create(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var request web.AlumniRequest
	if err := C.BodyParser(&request); err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err:= service.validate.Struct(request)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx,err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	Alumni := domain.Alumni{
		Nim: request.Nim,
		Nama: request.Nama,
		Jurusan: request.Jurusan,
		Angkatan: request.Angkatan,
		Tahun_lulus: request.Tahun_lulus,
		Email: request.Email,
		No_telepon: request.No_telepon,
		Alamat: request.Alamat,
	}
	Alumni = service.AlumniRepository.Save(ctx,tx,Alumni)

	return C.Status(fiber.StatusOK).JSON(helper.ToAlumniResponse(Alumni))
}

func (service *AlumniServiceImpl) Update(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var request web.AlumniUpdateRequest
	if err := C.BodyParser(&request); err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	request.Id = id
	err = service.validate.Struct(request)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx,err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)

	Alumni := domain.Alumni{
		Id: request.Id,
		Nim: request.Nim,
		Nama: request.Nama,
		Jurusan: request.Jurusan,
		Angkatan: request.Angkatan,
		Tahun_lulus: request.Angkatan,
		Email: request.Email,
		No_telepon: request.No_telepon,
		Alamat: request.No_telepon,
	}
	Alumni = service.AlumniRepository.Update(ctx,tx,Alumni)

	return C.Status(fiber.StatusOK).JSON(helper.ToAlumniResponse(Alumni))
}

func (service *AlumniServiceImpl) Delete(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx,err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	service.AlumniRepository.Delete(ctx,tx,id)
	return C.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "alumni deleted",
	})	
}

func (service *AlumniServiceImpl) FindById(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx,err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	Alumni,err := service.AlumniRepository.FindById(ctx,tx,id)
	if err != nil {
		return C.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "alumni not found"})
	}

	return C.Status(fiber.StatusOK).JSON(helper.ToAlumniResponse(Alumni))
}

func (service *AlumniServiceImpl) FindAll(C *fiber.Ctx) error {
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
	tx,err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)

	Alumnis := service.AlumniRepository.FindAll(ctx,tx,search,sortBy,order,limit,offset)
	var AlumniResponses []web.AlumniResponse
	for _,Alumni := range Alumnis{
		AlumniResponses = append(AlumniResponses,helper.ToAlumniResponse(Alumni))
	}
	return C.Status(fiber.StatusOK).JSON(AlumniResponses)
}

func (service *AlumniServiceImpl) FindPengangguran(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx,err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	Alumnis,count := service.AlumniRepository.FindPengangguran(ctx,tx)
	var AlumniResponses []web.AlumniResponse
	for _,Alumni := range Alumnis{
		AlumniResponses = append(AlumniResponses,helper.ToAlumniResponse(Alumni))
	}

	response := web.Pengangguran{
		Data: AlumniResponses,
		Jumlah: count,
	}

	return C.Status(fiber.StatusOK).JSON(response)
}