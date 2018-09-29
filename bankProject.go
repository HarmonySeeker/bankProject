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


type clientFullToXML struct{
  XMLName xml.Name `xml:"clients"`
  Client []clientToXML `xml:"client"`
}

type clientToXML struct{
  XMLName xml.Name `xml:"client"`

  IndNum int `xml:"uniNumber"`
  Name string `xml:"name"`
  Surname string `xml:"surname"`
  BirthDate string `xml:"dateOfBirth"`
  Address string `xml:"address"`
  Phone string `xml:"phoneNumber"`
  Passport string `xml:"passportID"`
  Cards Cards  `xml:"cards"`
  Accounts Accounts `xml:"accounts"`
}

type Cards struct{
    Card []Card `xml:"card"`
}
type Card struct{
    CardID string `xml:"cardNumber"`
    CardExpDate string `xml:"expirationDate"`
    CardCurrency string `xml:"cardCurrency"`
    CardBalance float64 `xml:"cardBalance"`
}

type Accounts struct{
    Account []Account `xml:"account"`
}
type Account struct{
    AccCurrency string `xml:"accountCurrency"`
    AccBalance float64 `xml:" accountBalance"`
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
  ConnectionDB = "user:password@tcp(127.0.0.1:3306)/bankProject"
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
      XMLOutput := clientToXML{}
      db, err := sql.Open("mysql", ConnectionDB)
      checkErr(err)
      defer db.Close()

      RetrieveClient(clientInfo, db, UniqueClient)
      RetrieveCard(cardInfo, db, UniqueClient)
      RetrieveAccount(accountInfo, db, UniqueClient)

      Name := strings.Split(string(clientInfo[0].FullName), " ")
      XMLOutput = clientToXML{IndNum: clientInfo[0].IndNum, Name: Name[0], Surname: Name[1], BirthDate: clientInfo[0].BirthDate, Address: clientInfo[0].Address, Phone: clientInfo[0].Phone, Passport: clientInfo[0].Passport}

      for _, data := range cardInfo{
        XMLOutput.Cards.Card = append(XMLOutput.Cards.Card, Card{CardID: data.ID, CardExpDate: data.ExpDate, CardCurrency: data.Currency, CardBalance: data.Balance})
      }

      for _, data := range cardInfo{
        XMLOutput.Accounts.Account = append(XMLOutput.Accounts.Account, Account{AccCurrency: data.Currency, AccBalance: data.Balance})
      }


      filename := "BankCredentialsOUT.xml"
      file, _ := os.Create(filename)

      xmlWriter := io.Writer(file)

      enc := xml.NewEncoder(xmlWriter)
      enc.Indent(" ","  ")
      if err := enc.Encode(XMLOutput); err != nil {
                 checkErr(err)
         }

  } else {fmt.Println("Error: Wrong settings parameter")}
}
