package repository

import (
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/web"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PekerjaanAlumniRepository interface{
	Save(ctx context.Context,tx *sql.Tx,PekerjaanAlumni domain.PekerjaanAlumni) (domain.PekerjaanAlumni)
	Update(ctx context.Context,tx *sql.Tx, PekerjaanAlumni domain.PekerjaanAlumni) (domain.PekerjaanAlumni,error)
	Delete(ctx context.Context,tx *sql.Tx, request web.PekerjaanAlumniDelete) error
	DeleteAdmin(ctx context.Context,tx *sql.Tx, id int64)
	FindById(ctx context.Context,tx *sql.Tx, PekerjaanAlumniId int64) (domain.PekerjaanAlumni,error)
	FindAll(ctx context.Context,tx *sql.Tx,search,sortBy,order string,limit,offset int) ([]domain.PekerjaanAlumni)
	DeleteHardAdmin(Ctx context.Context,tx *sql.Tx,id int64) error
	DeleteHardUser(Ctx context.Context,tx *sql.Tx,id int64) error
	Trash(Ctx context.Context,tx *sql.Tx) ([]domain.PekerjaanAlumni,error)
	RecoverAdmin(Ctx context.Context,tx *sql.Tx,id int64) error
	RecoverUser(Ctx context.Context,tx *sql.Tx,id int64) error
}

type PekerjaanAlumniRepository_impl struct{}

func NewPekerjaanAlumniRepository() PekerjaanAlumniRepository {
	return &PekerjaanAlumniRepository_impl{}
}

func (repository *PekerjaanAlumniRepository_impl) Save(ctx context.Context,tx *sql.Tx, PekerjaanAlumni domain.PekerjaanAlumni) domain.PekerjaanAlumni{
	SQL := `INSERT into pekerjaan_alumni(
	alumni_id,
	nama_perusahaan,
	posisi_jabatan,
	bidang_industri,
	lokasi_kerja,
	gaji_range,
	tanggal_mulai_kerja,
	tanggal_selesai_kerja,
	status_pekerjaan,
	deskripsi_pekerjaan
	) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`

	var LastInsertId int64

	err := tx.QueryRowContext(ctx,SQL,PekerjaanAlumni.AlumniId,PekerjaanAlumni.NamaPerusahaan,
		PekerjaanAlumni.PosisiJabatan,PekerjaanAlumni.BidangIndustri,PekerjaanAlumni.LokasiKerja,
		PekerjaanAlumni.GajiRange,PekerjaanAlumni.TanggalMulai,PekerjaanAlumni.TanggalSelesai,PekerjaanAlumni.StatusPekerjaan,PekerjaanAlumni.Deskripsi).Scan(&LastInsertId)
	
	if err != nil {
		panic(err)
	}

	PekerjaanAlumni.Id = LastInsertId

	return PekerjaanAlumni
}

func (repository *PekerjaanAlumniRepository_impl) DeleteHardAdmin(ctx context.Context, tx *sql.Tx,id int64) error {
	SQL:= `DELETE FROM pekerjaan_alumni WHERE id = $1`
	res,err := tx.ExecContext(ctx,SQL,id)
	if err != nil {
		return errors.New("something wrong")
	}
	val,err := res.RowsAffected()
	if err != nil {
		return errors.New("something wrong")
	}
	if val <= 0 {
		return errors.New("no data affected")
	}
	return nil
}

func (repository *PekerjaanAlumniRepository_impl) DeleteHardUser(ctx context.Context, tx *sql.Tx,id int64) error {
	SQL:= `DELETE FROM pekerjaan_alumni pa
USING alumni a
INNER JOIN users u ON a.user_id = u.id
WHERE pa.alumni_id = a.id
  AND pa.id = $1;
`
	res,err := tx.ExecContext(ctx,SQL,id)
	if err != nil {
		return errors.New("something wrong")
	}
	val,err := res.RowsAffected()
	if err != nil {
		return errors.New("something wrong")
	}
	if val <= 0 {
		return errors.New("no data affected")
	}
	return nil
}

func (repository *PekerjaanAlumniRepository_impl) Update(ctx context.Context, tx *sql.Tx, PekerjaanAlumni domain.PekerjaanAlumni) (domain.PekerjaanAlumni,error){
	SQL := `UPDATE pekerjaan_alumni SET alumni_id=$1,
	nama_perusahaan = $2,
	posisi_jabatan=$3,
	bidang_industri=$4,
	lokasi_kerja=$5,
	gaji_range=$6,
	tanggal_mulai_kerja=$7,
	tanggal_selesai_kerja=$8,
	status_pekerjaan=$9,
	deskripsi_pekerjaan = $10,
	updated_at=NOW() 
	WHERE id = $11 
	returning id`

	var id int64
	err := tx.QueryRowContext(ctx,SQL,PekerjaanAlumni.AlumniId,PekerjaanAlumni.NamaPerusahaan,PekerjaanAlumni.PosisiJabatan,PekerjaanAlumni.BidangIndustri,PekerjaanAlumni.LokasiKerja,
		PekerjaanAlumni.GajiRange,PekerjaanAlumni.TanggalMulai,PekerjaanAlumni.TanggalSelesai,PekerjaanAlumni.StatusPekerjaan,PekerjaanAlumni.Deskripsi,PekerjaanAlumni.Id).Scan(&id)

	if(err == sql.ErrNoRows){
		return PekerjaanAlumni,errors.New("pekerjaan alumni tidak ditemukan")
	}
	if err != nil {
		panic(err)
	}

	return PekerjaanAlumni,nil
}

func (repository *PekerjaanAlumniRepository_impl) Delete(ctx context.Context,tx *sql.Tx, request web.PekerjaanAlumniDelete) error {
	SQL := `
    UPDATE pekerjaan_alumni pa
    SET is_deleted = TRUE
    FROM alumni a
    WHERE pa.id = $2
      AND pa.alumni_id = a.id
      AND a.user_id = $1
      AND pa.is_deleted = FALSE
	`	
	res,err := tx.ExecContext(ctx,SQL,request.UserId,request.PekerjaanId)
	if err != nil {
		return errors.New("terjadi kesalahan")
	}

	ref,_ := res.RowsAffected()
	if ref == 0 {
		return errors.New("not authorized")
	}

	return nil
}

func (repository *PekerjaanAlumniRepository_impl) DeleteAdmin(ctx context.Context,tx *sql.Tx, id int64){
	SQL := `update pekerjaan_alumni pa SET is_deleted = TRUE where pa.id = $1`
	_,err := tx.ExecContext(ctx,SQL,id)
	if err != nil {
		panic(err)
	}
	
}

func (repository *PekerjaanAlumniRepository_impl) RecoverAdmin(ctx context.Context,tx *sql.Tx, id int64) error{
	SQL := `update pekerjaan_alumni pa SET pa.is_deleted = FALSE where pa.id = $1`
	_,err := tx.ExecContext(ctx,SQL,id)
	if err != nil {
		return errors.New("terjadi Kesalahan")
	}
	return nil
}

func (repository *PekerjaanAlumniRepository_impl) RecoverUser(ctx context.Context,tx *sql.Tx, id int64) error{
	SQL := `UPDATE pekerjaan_alumni pa
SET is_deleted = FALSE
FROM alumni a
INNER JOIN users u ON a.user_id = u.id
WHERE pa.alumni_id = a.id
  AND pa.id = $1;
`
	_,err := tx.ExecContext(ctx,SQL,id)
	if err != nil {
		return errors.New("terjadi Kesalahan")
	}
	return nil
}


func (repository *PekerjaanAlumniRepository_impl) FindById(ctx context.Context,tx *sql.Tx, PekerjaanAlumniId int64) (domain.PekerjaanAlumni,error) {
	SQL := `SELECT id,alumni_id,nama_perusahaan,posisi_jabatan,bidang_industri,lokasi_kerja,gaji_range,tanggal_mulai_kerja,tanggal_selesai_kerja,status_pekerjaan,deskripsi_pekerjaan from pekerjaan_alumni where id=$1`

	rows,err := tx.QueryContext(ctx,SQL,PekerjaanAlumniId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	PekerjaanAlumni := domain.PekerjaanAlumni{}
	if rows.Next(){
		err := rows.Scan(&PekerjaanAlumni.Id,&PekerjaanAlumni.AlumniId,&PekerjaanAlumni.NamaPerusahaan,&PekerjaanAlumni.PosisiJabatan,&PekerjaanAlumni.BidangIndustri,
		&PekerjaanAlumni.LokasiKerja,&PekerjaanAlumni.GajiRange,&PekerjaanAlumni.TanggalMulai,&PekerjaanAlumni.TanggalSelesai,&PekerjaanAlumni.StatusPekerjaan,&PekerjaanAlumni.Deskripsi)
		if err != nil {
			panic(err)
		}
		return PekerjaanAlumni,nil
	}else{
		return PekerjaanAlumni,errors.New("pekerjaan alumni not found")
	}
}

func (repository *PekerjaanAlumniRepository_impl) FindAll(ctx context.Context,tx *sql.Tx,search,sortBy,order string,limit,offset int) ([]domain.PekerjaanAlumni){

	allowedSort := map[string]bool{
		"id": true, "tanggal_mulai_kerja": true, "tanggal_selesai_kerja": true,
	}

	if !allowedSort[sortBy] {
		sortBy = "id"
	}

	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	SQL := fmt.Sprintf(`SELECT id,alumni_id,nama_perusahaan,posisi_jabatan,bidang_industri,lokasi_kerja,gaji_range,tanggal_mulai_kerja,tanggal_selesai_kerja,status_pekerjaan,deskripsi_pekerjaan from pekerjaan_alumni
						 WHERE nama_perusahaan ILIKE $1 AND is_deleted = FALSE
						 ORDER BY %s %s
						 LIMIT $2 OFFSET $3`,sortBy,order)

	rows,err := tx.QueryContext(ctx,SQL,"%"+search+"%", limit, offset)
	if err != nil {
		panic(err)
	}
	
	defer rows.Close()

	var PekerjaanAlumnis []domain.PekerjaanAlumni
	for rows.Next() {
		PekerjaanAlumni := domain.PekerjaanAlumni{}
		err := rows.Scan(&PekerjaanAlumni.Id,&PekerjaanAlumni.AlumniId,&PekerjaanAlumni.NamaPerusahaan,&PekerjaanAlumni.PosisiJabatan,&PekerjaanAlumni.BidangIndustri,
		&PekerjaanAlumni.LokasiKerja,&PekerjaanAlumni.GajiRange,&PekerjaanAlumni.TanggalMulai,&PekerjaanAlumni.TanggalSelesai,&PekerjaanAlumni.StatusPekerjaan,&PekerjaanAlumni.Deskripsi)
		if err != nil {
			panic(err)
		}
		PekerjaanAlumnis = append(PekerjaanAlumnis, PekerjaanAlumni)
	}
	return PekerjaanAlumnis
}

func (repository *PekerjaanAlumniRepository_impl) Trash(Ctx context.Context,tx *sql.Tx) ([]domain.PekerjaanAlumni,error){
	SQL := `SELECT id,alumni_id,nama_perusahaan,posisi_jabatan,bidang_industri,lokasi_kerja,gaji_range,tanggal_mulai_kerja,tanggal_selesai_kerja,status_pekerjaan,deskripsi_pekerjaan from pekerjaan_alumni
						 WHERE is_deleted IS TRUE
						`

	rows,err := tx.QueryContext(Ctx,SQL)
	if err != nil {
		return nil,err
	}
	
	defer rows.Close()

	var PekerjaanAlumnis []domain.PekerjaanAlumni
	for rows.Next() {
		PekerjaanAlumni := domain.PekerjaanAlumni{}
		err := rows.Scan(&PekerjaanAlumni.Id,&PekerjaanAlumni.AlumniId,&PekerjaanAlumni.NamaPerusahaan,&PekerjaanAlumni.PosisiJabatan,&PekerjaanAlumni.BidangIndustri,
		&PekerjaanAlumni.LokasiKerja,&PekerjaanAlumni.GajiRange,&PekerjaanAlumni.TanggalMulai,&PekerjaanAlumni.TanggalSelesai,&PekerjaanAlumni.StatusPekerjaan,&PekerjaanAlumni.Deskripsi)
		if err != nil {
			return nil,err
		}
		PekerjaanAlumnis = append(PekerjaanAlumnis, PekerjaanAlumni)
	}
	return PekerjaanAlumnis,nil
}