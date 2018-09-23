package main

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  //"database/sql"
  //"github.com/go-sql-driver/mysql"
  //"strings"
)

type clientCredentials struct{
  ClientUni []string `xml:"client>uniNumber"`
  ClientName []string `xml:"client>name"`
  ClientSurname []string `xml:"client>surname"`
  ClientBirth []string `xml:"client>dateOfBirth"`
  ClientAddress []string `xml:"client>address"`
  ClientPhone []string `xml:"client>phoneNumber"`
  ClientPassport []string `xml:"client>passportID"`

  CardID []string  `xml:"client>cards>card>cardNumber"`
  CardExpDate []string `xml:"client>cards>card>expirationDate"`
  CardCurrency []string `xml:"client>cards>card>cardCurrency"`
  CardBalance []string `xml:"client>cards>card>cardBalance"`

  AccountCurrency []string `xml:"client>accounts>account>accountCurrency"`
  AccountBalance []string `xml:"client>accounts>account>accountBalance"`

}

type clientMap struct{
  IndNum string
  BirthDate string
  Address string
  Phone string
  Passport string
}

type cardMap struct{
  IndNum string
  ID string
  ExpDate string
  Currency string
  Balance string
}

type accountMap struct{
  IndNum string
  Currency string
  Balance string
}

func main(){
  var XMLdata clientCredentials

  fmt.Println("\nStarted:\n")
  fmt.Println("Importing XML data...\n")

  XMLinput, _ := ioutil.ReadFile("bankCredentials.xml")
  xml.Unmarshal(XMLinput, &XMLdata)

  clientInfo := make(map[int]clientMap)
  cardInfo := make(map[int]cardMap)
  accountInfo := make(map[int]accountMap)

  for idx, _ := range XMLdata.ClientName{
    clientUni := XMLdata.ClientUni[idx]

    for idx, _  := range XMLdata.ClientName{
      //clientMapIdx := strings.Join([]string{XMLdata.ClientName[idx], XMLdata.ClientSurname[idx]} , " ")
      clientInfo[idx] = clientMap{clientUni, XMLdata.ClientBirth[idx], XMLdata.ClientAddress[idx], XMLdata.ClientPhone[idx], XMLdata.ClientPassport[idx]}
    }

    for idx, _ := range XMLdata.CardID{
      //cardMapIdx := strings.Join([]string{XMLdata.ClientName[idx], XMLdata.ClientSurname[idx]} , " ")
      cardInfo[idx] = cardMap{clientUni, XMLdata.CardID[idx], XMLdata.CardExpDate[idx], XMLdata.CardCurrency[idx], XMLdata.CardBalance[idx]}
    }

    for idx, _ := range XMLdata.CardID{
      //accountMapIdx := strings.Join([]string{XMLdata.ClientName[idx], XMLdata.ClientSurname[idx]} , " ")
      accountInfo[idx] = accountMap{clientUni, XMLdata.AccountCurrency[idx], XMLdata.AccountBalance[idx]}
    }
  }

  fmt.Println("clientMap:")
  for idx, data := range clientInfo{
    fmt.Println(idx,":", data.IndNum, data.BirthDate, data.Address, data.Phone, data.Passport)
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
  fmt.Println()
}
