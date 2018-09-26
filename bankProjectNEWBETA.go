package main

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  "strings"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "os"
  "io"
)

type clientCredentials struct{
  ClientUni []int `xml:"client>uniNumber"`
  ClientName []string `xml:"client>name"`
  ClientSurname []string `xml:"client>surname"`
  ClientBirth []string `xml:"client>dateOfBirth"`
  ClientAddress []string `xml:"client>address"`
  ClientPhone []string `xml:"client>phoneNumber"`
  ClientPassport []string `xml:"client>passportID"`

  CardID []string  `xml:"client>cards>card>cardNumber"`
  CardExpDate []string `xml:"client>cards>card>expirationDate"`
  CardCurrency []string `xml:"client>cards>card>cardCurrency"`
  CardBalance []float64 `xml:"client>cards>card>cardBalance"`

  AccountCurrency []string `xml:"client>accounts>account>accountCurrency"`
  AccountBalance []float64 `xml:"client>accounts>account>accountBalance"`

}

type clientMap struct{
  IndNum int
  FullName string
  BirthDate string
  Address string
  Phone string
  Passport string
}

type cardMap struct{
  IndNum int
  ID string
  ExpDate string
  Currency string
  Balance float64
}

type accountMap struct{
  IndNum int
  Currency string
  Balance float64
}

type clientToXML struct{
  XMLName xml.Name `xml:"clients"`

  IndNum int `xml:"uniNumber"`
  Name string `xml:"name"`
  Surname string `xml:"surname"`
  BirthDate string `xml:"dateOfBirth"`
  Address string `xml:"address"`
  Phone string `xml:"phoneNumber"`
  Passport string `xml:"passportID"`

  CardID string `xml:"cards>card>cardNumber"`
  CardExpDate string `xml:"cards>card>expirationDate"`
  CardCurrency string `xml:"cards>card>cardCurrency"`
  CardBalance float64 `xml:"cards>card>cardBalance"`

  AccCurrency string `xml:"accounts>account>accountCurrency"`
  AccBalance float64 `xml:"accounts>account>accountBalance"`
}

type ClientFullToXML struct{
  XMLName xml.Name `xml:"clients"`
  Client []clientToXML `xml:"client"`
}

 var ConnectionDB string
 var settings string

func RetrieveClient(dt map[int]clientMap, db *sql.DB, Uni int) map[int]clientMap{
  var innerArr clientMap
  var Uniq int
  var Name, Date, Address, Phone, Passport string
  err := db.QueryRow("select * from clientsBio where Uni = ?", Uni).Scan(&Uniq, &Name, &Date, &Address, &Phone, &Passport)
  checkErr(err)
  innerArr = clientMap{Uniq, Name, Date, Address, Phone, Passport}
  dt[0] = innerArr
  return dt
}

func RetrieveCard(dt map[int]cardMap, db *sql.DB, Uni int) map[int]cardMap{
  var ID, ExpDate, Currency string
  var Balance float64
  var CMap cardMap
  var Uniq int
  i := 0
  rows, err := db.Query("select * from clientsCard where Uni = ?", Uni)
  checkErr(err)
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(&Uniq, &ID, &ExpDate, &Currency, &Balance)
    checkErr(err)
    CMap = cardMap{Uniq, ID, ExpDate, Currency, Balance}
    dt[i] = CMap
    i++
  }
  err = rows.Err()
  checkErr(err)
  return dt
}

func RetrieveAccount(dt map[int]accountMap, db *sql.DB, Uni int)map[int]accountMap{
  var Uniq int
  var Currency string
  var Balance float64
  var aMap accountMap
  i := 0
  rows, err := db.Query("select * from clientsAccount where Uni = ?", Uni)
  checkErr(err)
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(&Uniq, &Currency, &Balance)
    checkErr(err)
    aMap = accountMap{Uniq, Currency, Balance}
    dt[i] = aMap
    i++
  }
  err = rows.Err()
  checkErr(err)
  return dt
}

func main(){
  //settings blog; check README for more info
  ConnectionDB = "Wanderer:Art.156.DBW.426@tcp(127.0.0.1:3306)/bankProject"
  settings = "Retrieve"
  UniqueClient := 77602
  //end of settings blog

  clientInfo := make(map[int]clientMap)
  cardInfo := make(map[int]cardMap)
  accountInfo := make(map[int]accountMap)

  if settings == "Import"{
    var XMLdata clientCredentials

    fmt.Println("\nStarted:\n")
    fmt.Println("Importing XML data...\n")

    XMLinput, _ := ioutil.ReadFile("bankCredentials.xml")
    xml.Unmarshal(XMLinput, &XMLdata)

    for idx, _ := range XMLdata.ClientName{
      clientUni := XMLdata.ClientUni[idx]

      for idx, _  := range XMLdata.ClientName{
        clientFullName := strings.Join([]string{XMLdata.ClientName[idx], XMLdata.ClientSurname[idx]} , " ")
        clientInfo[idx] = clientMap{clientUni, clientFullName, XMLdata.ClientBirth[idx], XMLdata.ClientAddress[idx], XMLdata.ClientPhone[idx], XMLdata.ClientPassport[idx]}
      }

      for idx, _ := range XMLdata.CardID{
        cardInfo[idx] = cardMap{clientUni, XMLdata.CardID[idx], XMLdata.CardExpDate[idx], XMLdata.CardCurrency[idx], XMLdata.CardBalance[idx]}
      }

      for idx, _ := range XMLdata.CardID{
        accountInfo[idx] = accountMap{clientUni, XMLdata.AccountCurrency[idx], XMLdata.AccountBalance[idx]}
      }
    }
    /*
    fmt.Println("clientMap:")
    for idx, data := range clientInfo{
      fmt.Println(idx,":", data.IndNum, data.FullName, data.BirthDate, data.Address, data.Phone, data.Passport)
    }
    fmt.Println()

    fmt.Println("cardMap:")
    for idx, data := range cardInfo{
      fmt.Println(idx,":", data.IndNum, data.ID, data.ExpDate, data.Currency, data.Balance)
    }
    fmt.Println()

    fmt.Println("accountMap:")
    for idx, data := range accountInfo {
      fmt.Println(idx,":", data.IndNum, data.Currency, data.Balance)
    }
    fmt.Println()*/

    //Commented code above can help in debugging

    dbClientWriter(clientInfo)
    dbCardWriter(cardInfo)
    dbAccountWriter(accountInfo)

    fmt.Println("Imported!")
  } else if settings == "Retrieve"{
      XMLOutput := &ClientFullToXML{}
      db, err := sql.Open("mysql", ConnectionDB)
      checkErr(err)
      defer db.Close()

      RetrieveClient(clientInfo, db, UniqueClient)
      RetrieveCard(cardInfo, db, UniqueClient)
      RetrieveAccount(accountInfo, db, UniqueClient)

      for _, data  := range clientInfo{
      Name := strings.Split(string(data.FullName), " ")
      XMLOutput.Client = append(XMLOutput.Client, clientToXML{IndNum: int(data.IndNum), Name: Name[0], Surname: Name[1], BirthDate: string(data.BirthDate), Address: string(data.Address), Phone: string(data.Phone), Passport: string(data.Passport)})
      }

      for _, data := range cardInfo{
      //cardCollection[idx] = {data.ID, data.ExpDate, data.Currency, data.Balance}
      XMLOutput.Client = append(XMLOutput.Client, clientToXML{CardID: string(data.ID), CardExpDate: string(data.ExpDate), CardCurrency: string(data.Currency), CardBalance: float64(data.Balance)})
      }

      for _, data := range accountInfo{
      XMLOutput.Client = append(XMLOutput.Client, clientToXML{AccCurrency: string(data.Currency), AccBalance: float64(data.Balance)})
      }

      filename := "newBankCredentials.xml"
      file, _ := os.Create(filename)

      xmlWriter := io.Writer(file)

      enc := xml.NewEncoder(xmlWriter)
      enc.Indent(" ","  ")
      if err := enc.Encode(XMLOutput); err != nil {
                 checkErr(err)
         }

  } else {fmt.Println("Error: Wrong settings parameter")}
}
