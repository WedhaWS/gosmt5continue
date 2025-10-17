package routes

import (
	"WedhaWS/utsgosmt5/alumni-tracer/app/services"
	"WedhaWS/utsgosmt5/alumni-tracer/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(S *services.AppServices) *fiber.App{
	app := fiber.New(fiber.Config{
        Prefork: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
        
        code := fiber.StatusInternalServerError

        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
        }

        return c.Status(code).JSON(fiber.Map{
            "error":   true,
            "message": err.Error(),
        })
    },
    })

	app.Post("/login",S.UserService.Login)
	app.Get("/trash",S.PekerjaanAlumniService.Trash)
	

	api := app.Group("/api",middleware.AuthRequired())
	api.Delete("/deleted/:id",S.PekerjaanAlumniService.DeleteHard)
	api.Get("/recover/:id",S.PekerjaanAlumniService.Recover)
	api.Delete("/pekerjaan/:id",S.PekerjaanAlumniService.Delete)
	api.Get("/users/:id",S.UserService.FindById)
	api.Get("/users",S.UserService.FindAll)
	api.Get("/alumni/:id",S.AlumniService.FindById)
	api.Get("/alumni",S.AlumniService.FindAll)
	api.Get("/pekerjaan/:id",S.PekerjaanAlumniService.FindById)
	api.Get("/pekerjaan",S.PekerjaanAlumniService.FindAll)
	app.Get("/pengangguran",S.AlumniService.FindPengangguran)

	protected := api.Group("/admin",middleware.AdminOnly())
	protected.Post("/users",S.UserService.Create)
	protected.Put("/users/:id",S.UserService.Update)
	protected.Delete("/users/:id",S.UserService.Delete)
	protected.Post("/alumni",S.AlumniService.Create)
	protected.Put("/alumni/:id",S.AlumniService.Update)
	protected.Delete("/alumni/:id",S.AlumniService.Delete)
	protected.Post("/pekerjaan",S.PekerjaanAlumniService.Create)
	protected.Put("/pekerjaan/:id",S.PekerjaanAlumniService.Update)
	
	
	return app
}