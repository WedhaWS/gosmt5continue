package helper

import (
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/web"
)



func ToUserResponse(user domain.Users) web.UserResponse {
	return web.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func ToAlumniResponse(alumni domain.Alumni) web.AlumniResponse {
	return web.AlumniResponse{
		Id:          alumni.Id,
		Nim:         alumni.Nim,
		Nama:        alumni.Nama,
		Jurusan:     alumni.Jurusan,
		Angkatan:    alumni.Angkatan,
		Tahun_lulus: alumni.Tahun_lulus,
		Email:       alumni.Email,
		No_telepon:  alumni.No_telepon,
		Alamat:      alumni.Alamat,
	}
}

func ToPekerjaanAlumniResponse(PekerjaanAlumni domain.PekerjaanAlumni) web.PekerjaanAlumniResponse {
	return web.PekerjaanAlumniResponse{
		Id:             PekerjaanAlumni.Id,
		AlumniId:       PekerjaanAlumni.AlumniId,
		NamaPerusahaan: PekerjaanAlumni.NamaPerusahaan,
		PosisiJabatan:  PekerjaanAlumni.PosisiJabatan,
		BidangIndustri: PekerjaanAlumni.BidangIndustri,
		LokasiKerja:	PekerjaanAlumni.LokasiKerja,
		GajiRange:      PekerjaanAlumni.GajiRange,
		TanggalMulai: 	PekerjaanAlumni.TanggalMulai,
		TanggalSelesai: PekerjaanAlumni.TanggalSelesai,
		StatusPekerjaan: PekerjaanAlumni.StatusPekerjaan,
		Deskripsi: 		PekerjaanAlumni.Deskripsi,
	}
}
