package web

type AlumniResponse struct {
	Id int64 `json:"id"`
	Nim string `json:"nim"`
	Nama string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Angkatan int64 `json:"angkatan"`
	Tahun_lulus int64 `json:"tahun_lulus"`
	Email string `json:"email" `
	No_telepon string `json:"no_telepon"`
	Alamat string `json:"alamat"`
}