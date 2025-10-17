package web

import (
	"time"
)

type PekerjaanAlumniUpdateRequest struct {
	Id				int64	`json:"-" validate:"required"`
    AlumniId        int64 	`json:"id_alumni" validate:"required"`
    NamaPerusahaan  string	`json:"nama_perusahaan" validate:"required"`
    PosisiJabatan   string	`json:"posisi_jabatan" validate:"required"`
    BidangIndustri  string	`json:"bidang_industri" validate:"required"`
    LokasiKerja     string	`json:"lokasi_kerja" validate:"required"`
    GajiRange       string	`json:"gaji_range" validate:"required"`
    TanggalMulai    time.Time `json:"tanggal_mulai" validate:"required"`
    TanggalSelesai  *time.Time `json:"tanggal_selesai" validate:"required"`
    StatusPekerjaan string	`json:"status_pekerjaan" validate:"required"`
    Deskripsi       string	`json:"deskripsi" validate:"required"`
}