package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func checkAccountTable(table string) {
    db, err := sql.Open("mysql", ConnectionDB)
    defer db.Close()
    checkErr(err)

    data, err := db.Query("select table_name from information_schema.tables where table_schema = 'bankProject' and table_name = ?;", table)
    checkErr(err)

    if !data.Next(){
      _, err = db.Exec("create table clientsAccount(Uni int, Currency VARCHAR(20), Balance FLOAT); ")
      checkErr(err)
    }
}

func insertAccountData(datamap map[int]accountMap){
  db, err := sql.Open("mysql", ConnectionDB)
  checkErr(err)
  defer db.Close()

  stmt, err := db.Prepare("INSERT INTO clientsAccount VALUES(?, ?, ?);")
  checkErr(err)
  for _, data := range datamap{
  _, err = stmt.Exec(int(data.IndNum), string(data.Currency), float64(data.Balance))
    checkErr(err)
  }
}



func dbAccountWriter(database_info map[int]accountMap){

  db, err := sql.Open("mysql", ConnectionDB)
  checkErr(err)
  defer db.Close()

  checkAccountTable("clientsAccount")
  insertAccountData(database_info)

}
