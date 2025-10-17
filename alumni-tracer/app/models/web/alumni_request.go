package web

type AlumniRequest struct {
	Nim string `json:"nim" validate:"required"`
	Nama string `json:"nama" validate:"required"`
	Jurusan string `json:"jurusan" validate:"required"`
	Angkatan int64 `json:"angkatan" validate:"required"`
	Tahun_lulus int64 `json:"tahun_lulus" validate:"required"`
	Email string `json:"email" validate:"required"`
	No_telepon string `json:"no_telepon" validate:"required"`
	Alamat string `json:"alamat" validate:"required"`
}