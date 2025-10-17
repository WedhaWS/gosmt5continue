package services

import (
	"WedhaWS/utsgosmt5/alumni-tracer/helper"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/web"
	"WedhaWS/utsgosmt5/alumni-tracer/app/repository"
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator"
	"time"
	"strconv"
)

type PekerjaanAlumniService interface{
	Create(C *fiber.Ctx) error
	Update(C *fiber.Ctx) error
	Delete(C *fiber.Ctx) error
	FindById(C *fiber.Ctx) error
	FindAll(C *fiber.Ctx) error
	Trash(C *fiber.Ctx) error
	DeleteHard(C *fiber.Ctx) error
	Recover(C *fiber.Ctx) error
}

type PekerjaanAlumniServiceImpl struct{
	PekerjaanAlumniRespository repository.PekerjaanAlumniRepository
	DB *sql.DB
	validate *validator.Validate
}

func NewPekerjaanAlumniService(PekerjaanAlumniRepository repository.PekerjaanAlumniRepository,DB *sql.DB,validate *validator.Validate) PekerjaanAlumniService {
	return &PekerjaanAlumniServiceImpl{
		PekerjaanAlumniRespository: PekerjaanAlumniRepository,
		DB: DB,
		validate: validate,
	}
}

func (service *PekerjaanAlumniServiceImpl) Recover(C *fiber.Ctx) error{
	ctx := context.Background()

	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx,err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	if C.Locals("role").(string) == "user"{
		err = service.PekerjaanAlumniRespository.RecoverUser(ctx,tx,id)
		
	}else {
		service.PekerjaanAlumniRespository.RecoverAdmin(ctx,tx,id)
	}

	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	
	return C.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "pekerjaan alumni recover",
	})	
}

func (service *PekerjaanAlumniServiceImpl) Create(C *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var request web.PekerjaanAlumniRequest
	if err := C.BodyParser(&request); err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := service.validate.Struct(request)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tx,err := service.DB.Begin()
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer helper.CommitOrRollback(tx)
	PekerjaanAlumni := domain.PekerjaanAlumni{
		AlumniId: request.AlumniId,
		NamaPerusahaan: request.NamaPerusahaan,
		PosisiJabatan: request.PosisiJabatan,
		BidangIndustri: request.BidangIndustri,
		LokasiKerja: request.LokasiKerja,
		GajiRange: request.GajiRange,
		TanggalMulai: request.TanggalMulai,
		TanggalSelesai: request.TanggalSelesai,
		StatusPekerjaan: request.StatusPekerjaan,
		Deskripsi: request.Deskripsi,

	}
	PekerjaanAlumni = service.PekerjaanAlumniRespository.Save(ctx,tx,PekerjaanAlumni)

	return C.Status(fiber.StatusOK).JSON(helper.ToPekerjaanAlumniResponse(PekerjaanAlumni))
}

func (service *PekerjaanAlumniServiceImpl) Update(C *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var request web.PekerjaanAlumniUpdateRequest
	if err := C.BodyParser(&request); err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	request.AlumniId = id
	err = service.validate.Struct(request)
	if err != nil {
		panic(err)
	}
	tx,err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	PekerjaanAlumni := domain.PekerjaanAlumni{
		Id: request.Id,
		AlumniId: request.AlumniId,
		NamaPerusahaan: request.NamaPerusahaan,
		PosisiJabatan: request.PosisiJabatan,
		BidangIndustri: request.BidangIndustri,
		LokasiKerja: request.LokasiKerja,
		GajiRange: request.GajiRange,
		TanggalMulai: request.TanggalMulai,
		TanggalSelesai: request.TanggalSelesai,
		StatusPekerjaan: request.StatusPekerjaan,
		Deskripsi: request.Deskripsi,
	}
	PekerjaanAlumni,err = service.PekerjaanAlumniRespository.Update(ctx,tx,PekerjaanAlumni)
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return C.Status(fiber.StatusOK).JSON(helper.ToPekerjaanAlumniResponse(PekerjaanAlumni))
}

func (service *PekerjaanAlumniServiceImpl) Delete(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	UserId := C.Locals("user_id").(int64)
	tx,err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	if C.Locals("role").(string) == "user"{
		request := web.PekerjaanAlumniDelete{
			UserId: UserId,
			PekerjaanId: id,
		}

		err = service.PekerjaanAlumniRespository.Delete(ctx,tx,request)
		
	}else {
		service.PekerjaanAlumniRespository.DeleteAdmin(ctx,tx,id)
	}

	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	
	return C.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "pekerjaan alumni deleted",
	})	
}

func (service *PekerjaanAlumniServiceImpl) FindById(C *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx,err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	PekerjaanAlumni,err := service.PekerjaanAlumniRespository.FindById(ctx,tx,id)
	if err != nil {
		return C.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "pekerjaan alumni not found"})
	}
	return C.Status(fiber.StatusOK).JSON(helper.ToPekerjaanAlumniResponse(PekerjaanAlumni))
}

func (service *PekerjaanAlumniServiceImpl) FindAll(C *fiber.Ctx) error{
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
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	PekerjaanAlumnis := service.PekerjaanAlumniRespository.FindAll(ctx,tx,search,sortBy,order,limit,offset)
	var PekerjaanAlumniResponses []web.PekerjaanAlumniResponse

	for _,PekerjaanAlumni := range PekerjaanAlumnis{
		PekerjaanAlumniResponses = append(PekerjaanAlumniResponses, helper.ToPekerjaanAlumniResponse(PekerjaanAlumni))
	}
	return C.Status(fiber.StatusOK).JSON(PekerjaanAlumniResponses)
}

func (service *PekerjaanAlumniServiceImpl) Trash(C *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx,err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	PekerjaanAlumnis,err := service.PekerjaanAlumniRespository.Trash(ctx,tx)
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	var PekerjaanAlumniResponses []web.PekerjaanAlumniResponse

	for _,PekerjaanAlumni := range PekerjaanAlumnis{
		PekerjaanAlumniResponses = append(PekerjaanAlumniResponses, helper.ToPekerjaanAlumniResponse(PekerjaanAlumni))
	}
	return C.Status(fiber.StatusOK).JSON(PekerjaanAlumniResponses)
}

func (service *PekerjaanAlumniServiceImpl) DeleteHard(C *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id,err := strconv.ParseInt(C.Params("id"),10,64)
	if err != nil {
		return C.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tx,err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)
	if C.Locals("role").(string) == "user"{
		err = service.PekerjaanAlumniRespository.DeleteHardUser(ctx,tx,id)
		
	}else {
		err = service.PekerjaanAlumniRespository.DeleteHardAdmin(ctx,tx,id)
	}
	
	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return C.Status(fiber.StatusOK).JSON(fiber.Map{"Success" : "Berhasil Menghapus Data"})
}