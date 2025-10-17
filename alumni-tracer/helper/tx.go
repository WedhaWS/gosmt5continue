package helper

import (
	"database/sql"
)

func CommitOrRollback(tx *sql.Tx){
	err := recover()
	if err != nil {
		error :=tx.Rollback()
		if error != nil {
			panic(error)
		}
		panic(err)
	}else{
		errorCommit :=tx.Commit()
		if errorCommit != nil {
			panic(errorCommit)
		}
	}
}