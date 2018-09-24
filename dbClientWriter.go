package main

import (
  "log"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
  if err != nil{
    log.Fatal(err)
  }
}

func checkClientsTable(table string) {
    db, err := sql.Open("mysql",
    "Wanderer:Art.156.DBW.426@tcp(127.0.0.1:3306)/bankProject")
    defer db.Close()
    checkErr(err)

    data, err := db.Query("select table_name from information_schema.tables where table_schema = 'bankProject' and table_name = ?;", table)
    checkErr(err)

    if !data.Next(){
      _, err = db.Exec("create table clientsBio(Uni int, FullName VARCHAR(30), BirthDate DATE, Address VARCHAR(30), Phone VARCHAR(15), PassportID VARCHAR(15)); ")
      checkErr(err)
    }
}

func insertClientsData(datamap map[int]clientMap){
  db, err := sql.Open("mysql",
  "Wanderer:Art.156.DBW.426@tcp(127.0.0.1:3306)/bankProject")
  checkErr(err)
  defer db.Close()

  stmt, err := db.Prepare("INSERT INTO clientsBio VALUES(?, ?, ?, ?, ?, ?);")
  checkErr(err)
  for _, data := range datamap{
  _, err := stmt.Exec(int(data.IndNum), string(data.FullName), string(data.BirthDate), string(data.Address), string(data.Phone), string(data.Passport))
    checkErr(err)
  }
}



func dbClientWriter(database_info map[int]clientMap){

  db, err := sql.Open("mysql",
  "Wanderer:Art.156.DBW.426@tcp(127.0.0.1:3306)/bankProject")
  checkErr(err)
  defer db.Close()

  checkClientsTable("clientsBio")
  insertClientsData(database_info)

}
