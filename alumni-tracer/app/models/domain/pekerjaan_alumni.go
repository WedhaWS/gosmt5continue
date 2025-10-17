package domain

import "time"

type PekerjaanAlumni struct {
    Id              int64
    AlumniId        int64
    NamaPerusahaan  string
    PosisiJabatan   string
    BidangIndustri  string
    LokasiKerja     string
    GajiRange       string
    TanggalMulai    time.Time
    TanggalSelesai  *time.Time
    StatusPekerjaan string
    Deskripsi       string
}
