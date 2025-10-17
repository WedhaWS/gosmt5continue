package services

type AppServices struct {
	AlumniService AlumniService
	PekerjaanAlumniService PekerjaanAlumniService
	UserService UserService
}

func NewAppService(
	alumni AlumniService,
	pekerjaan PekerjaanAlumniService,
	user UserService,
) *AppServices{
	return &AppServices{
		AlumniService: alumni,
		PekerjaanAlumniService: pekerjaan,
		UserService: user,
	}
}