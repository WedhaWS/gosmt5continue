package web

import (
	"time"
)

type PekerjaanAlumniResponse struct {
	Id				int64	`json:"id" `
    AlumniId        int64 	`json:"id_alumni" `
    NamaPerusahaan  string	`json:"nama_perusahaan" `
    PosisiJabatan   string	`json:"posisi_jabatan" `
    BidangIndustri  string	`json:"bidang_industri" `
    LokasiKerja     string	`json:"lokasi_kerja" `
    GajiRange       string	`json:"gaji_range" `
    TanggalMulai    time.Time `json:"tanggal_mulai" `
    TanggalSelesai  *time.Time `json:"tanggal_selesai" `
    StatusPekerjaan string	`json:"status_pekerjaan" `
    Deskripsi       string	`json:"deskripsi" `
}