package repository

import (
	"WedhaWS/utsgosmt5/alumni-tracer/app/models/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository interface{
	Save(ctx context.Context,tx *sql.Tx,User domain.Users) domain.Users
	Update(ctx context.Context, tx *sql.Tx,User domain.Users) domain.Users
	Delete(ctx context.Context, tx *sql.Tx, UserId int64) 
	FindById(ctx context.Context, tx *sql.Tx, UserId int64) (domain.Users,error)
	FindAll(ctx context.Context, tx *sql.Tx,search,sortBy,order string,limit,offset int) []domain.Users
	FindByEmailPassword(ctx context.Context, tx *sql.Tx,User domain.Users) (domain.Users,error)
}

type UserRepositoryImpl struct{}

func NewUsersRepository() UserRepository{
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) FindByEmailPassword(ctx context.Context, tx *sql.Tx,User domain.Users) (domain.Users,error){
	SQL := "select id,username,email,password_hash,role,created_at from users where email = $1 OR password_hash = $2"
	rows,err := tx.QueryContext(ctx,SQL,User.Email,User.Password_Hash)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	user := domain.Users{}
	if(rows.Next()){
		err := rows.Scan(&user.Id,&user.Username,&user.Email,&user.Password_Hash,&user.Role,&user.CreatedAt)
		if err!=nil{
			panic(err)
		}
		return user,nil
	}else{
		return user,errors.New("user not found")
	}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context,tx *sql.Tx,User domain.Users) (domain.Users) {
	SQL := "insert into users(username,email,password_hash,role) values ($1,$2,$3,$4) returning id"
	var LastInsertId int64
	err := tx.QueryRowContext(ctx,SQL,User.Username,User.Email,User.Password_Hash,User.Role).Scan(&LastInsertId)
	if err != nil{
		panic(err)
	}
	Id := LastInsertId
	User.Id = Id
	return User
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.Users) (domain.Users) {
    SQL := "UPDATE users SET username = $1, email = $2, password_hash = $3, role = $4 WHERE id = $5"
    _, err := tx.ExecContext(ctx, SQL, user.Username, user.Email, user.Password_Hash, user.Role, user.Id)
    if err != nil {
        panic(err)
    }
    return user
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, UserId int64){
	SQL := "delete from users where id = $1"
	_,err := tx.ExecContext(ctx,SQL,UserId)
	if err != nil{
		panic(err)
	}
}

func (repository *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, UserId int64)(domain.Users,error){
	SQL := "select id,username,email,role,created_at from users where id = $1"
	rows,err := tx.QueryContext(ctx,SQL,UserId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	user := domain.Users{}
	if(rows.Next()){
		err := rows.Scan(&user.Id,&user.Username,&user.Email,&user.Role,&user.CreatedAt)
		if err!=nil{
			panic(err)
		}
		return user,nil
	}else{
		return user,errors.New("user not found")
	}
}

func (repostory *UserRepositoryImpl) FindAll(ctx context.Context,tx *sql.Tx,search,sortBy,order string,limit,offset int)([]domain.Users){
		
	allowedSort := map[string]bool{
		"id": true, "username": true, "email": true, "created_at": true, "role": true,
	}

	if !allowedSort[sortBy] {
		sortBy = "id" // default
	}


	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	SQL := fmt.Sprintf(`
		SELECT id, username, email, role, created_at
		FROM users
		WHERE username ILIKE $1 OR email ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := tx.QueryContext(ctx, SQL, "%"+search+"%", limit, offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var users []domain.Users
	for rows.Next() {
		user := domain.Users{}
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return users
}
