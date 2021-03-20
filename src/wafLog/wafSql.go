package wafLog

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type LogDB struct {
	Id        string `db:"id"`
	LogPath   string `db:"path"`
	Timestamp string `db:"timestamp"`
}

var Db *sqlx.DB

//caonimacaoca
func init() {
	//database, err := sqlx.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/waf")
	database, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/log")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
}

//func Insert(l LogDB){
//	r, err := Db.Exec("insert into log(id, logPath, logTime)values(?, ?, ?)", l.Id, l.LogPath, l.Timestamp)
//	if err != nil {
//		fmt.Println("exec failed, ", err)
//		return
//	}
//	id, err := r.LastInsertId()
//	if err != nil {
//		fmt.Println("exec failed, ", err)
//		return
//	}
//
//	fmt.Println("insert succ:", id)
//}

func Select(id string) ([]LogDB, error) {
	var log []LogDB
	fmt.Println(id)
	err := Db.Select(&log, "select id, path, timestamp from log where id=?", id)
	return log, err
}

func SelectAll() ([]LogDB, error) {
	var log []LogDB
	err := Db.Select(&log, "select id, path, timestamp from log")
	return log, err
}

//func Update(l LogDB){
//	res, err := Db.Exec("update log set id=?,logPath=?,logTime=? where id=?", l.Id,l.LogPath,l.Timestamp,l.Id)
//	if err != nil {
//		fmt.Println("exec failed, ", err)
//		return
//	}
//	row, err := res.RowsAffected()
//	if err != nil {
//		fmt.Println("rows failed, ",err)
//	}
//	fmt.Println("update succ:",row)
//}
//
//func Delete(id uuid.UUID){
//	res, err := Db.Exec("delete from log where id=?", id)
//	if err != nil {
//		fmt.Println("exec failed, ", err)
//		return
//	}
//
//	row,err := res.RowsAffected()
//	if err != nil {
//		fmt.Println("rows failed, ",err)
//	}
//	fmt.Println("delete succ: ",row)
//}

func Close() {
	fmt.Println("Close DB")
	Db.Close()
}
