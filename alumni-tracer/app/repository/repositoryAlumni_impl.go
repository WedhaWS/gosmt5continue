package repository

import (
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type AlumniRepository interface{
	Save(Ctx context.Context,tx *sql.Tx,Alumni domain.Alumni) domain.Alumni
	Update(Ctx context.Context, tx *sql.Tx, Alumni domain.Alumni) domain.Alumni
	Delete(Ctx context.Context, tx *sql.Tx, AlumniId int64)
	FindById(Ctx context.Context, tx *sql.Tx, AlumniId int64) (domain.Alumni,error)
	FindAll(Ctx context.Context, tx *sql.Tx,search,sortBy,order string,limit,offset int) []domain.Alumni
	FindPengangguran(Ctx context.Context,tx *sql.Tx) ([]domain.Alumni,int64)
	
}

type AlumniRepositoryImpl struct{}

func NewAlumniRepository() AlumniRepository {
	return &AlumniRepositoryImpl{}
}

func (repository *AlumniRepositoryImpl) Save(ctx context.Context,tx *sql.Tx,Alumni domain.Alumni) (domain.Alumni) {
	SQL := "insert into alumni(nim,nama,jurusan,angkatan,tahun_lulus,email,no_telepon,alamat) values($1,$2,$3,$4,$5,$6,$7,$8) returning id"
	var LastInsertId int64
	err := tx.QueryRowContext(ctx,SQL,Alumni.Nim,Alumni.Nama,Alumni.Jurusan,Alumni.Angkatan,Alumni.Tahun_lulus,Alumni.Email,Alumni.No_telepon,Alumni.Alamat).Scan(&LastInsertId)
	if err != nil {
		panic(err)
	}
	Alumni.Id = LastInsertId
	return Alumni
}

func (repository *AlumniRepositoryImpl) Update(ctx context.Context, tx *sql.Tx,Alumni domain.Alumni) (domain.Alumni){
	SQL := `
        UPDATE alumni
        SET nim=$1, nama=$2, jurusan=$3, angkatan=$4, tahun_lulus=$5,
            email=$6, no_telepon=$7, alamat=$8, updated_at=NOW()
        WHERE id=$9
    `
	_,err := tx.ExecContext(ctx,SQL,Alumni.Nim,Alumni.Nama,Alumni.Jurusan,Alumni.Angkatan,Alumni.Tahun_lulus,Alumni.Email,Alumni.No_telepon,Alumni.Alamat,Alumni.Id)
	if err != nil{
		panic(err)
	}

	return Alumni
}

func (repository *AlumniRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, AlumniId int64){
	SQL := `DELETE from alumni where id = $1`
	_,err := tx.ExecContext(ctx,SQL,AlumniId)
	if err != nil {
		panic(err)
	}
}

func (repository *AlumniRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, AlumniId int64) (domain.Alumni, error){
	SQL := `SELECT id,nim,nama,jurusan,angkatan,tahun_lulus,email,no_telepon,alamat from alumni where id = $1`
	rows,err := tx.QueryContext(ctx,SQL,AlumniId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	Alumni := domain.Alumni{}
	if rows.Next() {
		err := rows.Scan(&Alumni.Id,&Alumni.Nim,&Alumni.Nama,&Alumni.Jurusan,&Alumni.Angkatan,&Alumni.Tahun_lulus,&Alumni.Email,&Alumni.No_telepon,&Alumni.Alamat)
		if err != nil {
			panic(err)
		}
		return Alumni,nil
	}else{
		return Alumni,errors.New("alumni not found")
	}
}

func (repository *AlumniRepositoryImpl) FindAll(ctx context.Context,tx *sql.Tx,search,sortBy,order string,limit,offset int) []domain.Alumni{

	allowedSort := map[string]bool{
		"id": true, "nama": true, "jurusan": true, "tahun_lulus": true,
	}

	if !allowedSort[sortBy] {
		sortBy = "id"
	}

	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	SQL := fmt.Sprintf(
    `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat
     FROM alumni
     WHERE nama ILIKE $1 OR jurusan ILIKE $1 OR tahun_lulus::text ILIKE $1
     ORDER BY %s %s
     LIMIT $2 OFFSET $3`, sortBy, order)


	rows, err := tx.QueryContext(ctx, SQL, "%"+search+"%", limit, offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var Alumnis []domain.Alumni
	for rows.Next() {
		Alumni := domain.Alumni{}
		err := rows.Scan(&Alumni.Id,&Alumni.Nim,&Alumni.Nama,&Alumni.Jurusan,&Alumni.Angkatan,&Alumni.Tahun_lulus,&Alumni.Email,&Alumni.No_telepon,&Alumni.Alamat)
		if err != nil {
			panic(err)
		}
		Alumnis = append(Alumnis,Alumni)
	}

	return Alumnis
}

func (respostory *AlumniRepositoryImpl) FindPengangguran(ctx context.Context,tx *sql.Tx) ([]domain.Alumni,int64){
	SQL := `SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.no_telepon, a.alamat FROM alumni a 
			LEFT JOIN pekerjaan_alumni p
			ON a.id = p.alumni_id
			WHERE p.id IS NULL`

	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var Alumnis []domain.Alumni
	for rows.Next() {
		Alumni := domain.Alumni{}
		err := rows.Scan(&Alumni.Id,&Alumni.Nim,&Alumni.Nama,&Alumni.Jurusan,&Alumni.Angkatan,&Alumni.Tahun_lulus,&Alumni.Email,&Alumni.No_telepon,&Alumni.Alamat)
		if err != nil {
			panic(err)
		}
		Alumnis = append(Alumnis,Alumni)
	}

	return Alumnis,int64(len(Alumnis))
}