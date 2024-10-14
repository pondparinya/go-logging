package main

//---------------------------------------------------------

import (
	"encoding/json"
	"fmt"
	"time"

	"go.kbtg.tech/733/go-ecslog"
)

//---------------------------------------------------------

func test_CensorFields(v Sample) {
  fmt.Printf("### test_CensorFields:\n")

  {
    e := ecslog.CensorFields(&v)
    fmt.Printf("%v\n", e)
    s, _ := json.Marshal(v)
    fmt.Printf("%s\n", string(s))
  }

  fmt.Printf("\n")

  {
    l := []Sample{v}
    e := ecslog.CensorFields(&l)
    fmt.Printf("%v\n", e)
    s, _ := json.Marshal(l)
    fmt.Printf("%s\n", string(s))
  }

  fmt.Printf("\n\n")
}

//---------------------------------------------------------

func test_CensorValue(v Sample) {
  fmt.Printf("### test_CensorValue:\n")

  {
    s := ecslog.CensorValue(v)
    fmt.Printf("%s\n", s)
  }

  {
    s := ecslog.CensorValue(v)
    fmt.Printf("%s\n", s)
  }

  fmt.Printf("\n")

  {
    l := []Sample{v}
    s := ecslog.CensorValue(l)
    fmt.Printf("%s\n", s)
  }

  {
    l := []Sample{v}
    s := ecslog.CensorValue(l)
    fmt.Printf("%s\n", s)
  }

  fmt.Printf("\n\n")
}

//---------------------------------------------------------

type Account struct {
  Number string
  Name   string
}

type Address struct {
  Line1    string
  Line2    string
  Province string
  Road     string
  Soi      string
  Tumbol   string
}

type Sample struct {
  Account     Account
  Address     Address
  Birth       *time.Time
  License     string
  Card        string
  Email       string
  CitizenID   string
  Passport    string
  Tax         string
  MobileNO    string
  Username    string
  Password    string
  Nationality string
  Race        string
  Firstname   string
  Surname     string
  // Name
}

func main() {

  t := time.Now()
  v := Sample{
    Account: Account{
      Number: "Number",
      Name:   "Name",
    },

    Address: Address {
      Line1:    "Line1",
      Line2:    "Line2",
      Province: "Province",
      Road:     "Road",
      Soi:      "Soi",
      Tumbol:   "Tumbol",
    },

    Birth:       &t,
    License:     "License",
    Card:        "Card",
    Email:       "Email",
    CitizenID:   "CitizenID",
    Passport:    "Passport",
    Tax:         "Tax",
    MobileNO:    "MobileNO",
    Username:    "Username",
    Password:    "Password",
    Nationality: "Nationality",
    Race:        "Race",
    Firstname:   "Firstname",
    Surname:     "Surname",
  }

  test_CensorFields(v)
  test_CensorValue(v)
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------

