package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func checkCardTable(table string) {
    db, err := sql.Open("mysql", ConnectionDB)
    defer db.Close()
    checkErr(err)

    data, err := db.Query("select table_name from information_schema.tables where table_schema = 'bankProject' and table_name = ?;", table)
    checkErr(err)

    if !data.Next(){
      _, err = db.Exec("create table clientsCard(Uni int, cardNumber VARCHAR(30), ExpDate DATE, Currency VARCHAR(15), Balance FLOAT); ")
      checkErr(err)
    }
}

func insertCardData(datamap map[int]cardMap){
  db, err := sql.Open("mysql", ConnectionDB)
  checkErr(err)
  defer db.Close()

  stmt, err := db.Prepare("INSERT INTO clientsCard VALUES(?, ?, ?, ?, ?);")
  checkErr(err)
  for _, data := range datamap{
  _, err = stmt.Exec(int(data.IndNum), string(data.ID), string(data.ExpDate), string(data.Currency), float64(data.Balance))
    checkErr(err)
  }
}



func dbCardWriter(database_info map[int]cardMap){

  db, err := sql.Open("mysql", ConnectionDB)
  checkErr(err)
  defer db.Close()

  checkCardTable("clientsCard")
  insertCardData(database_info)

}
