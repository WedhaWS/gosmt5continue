//go:build wireinject
// +build wireinject

package injector

import (
	"WedhaWS/alumni-tracer/repository"
	"WedhaWS/alumni-tracer/services"
	"WedhaWS/alumni-tracer/db"
	"WedhaWS/alumni-tracer/helper"
	"github.com/google/wire"
)

func InitializedService() *services.AppServices{
	wire.Build(
		db.NewDB,
		helper.NewValidator,
		repository.NewUsersRepository,
		repository.NewAlumniRepository,
		repository.NewPekerjaanAlumniRepository,
		services.NewUserService,
		services.NewAlumniService,
		services.NewPekerjaanAlumniService,
		services.NewAppService,
	)

	return nil
}