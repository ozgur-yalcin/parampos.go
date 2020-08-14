[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/OzqurYalcin/turkpos/blob/master/LICENSE.md)

# TurkPos
TurkPos (ParamPos) API with golang

# Installation
```bash
go get github.com/OzqurYalcin/turkpos
go get github.com/google/uuid
```

# Sanalpos satış işlemi
```go
package main

import (
	"fmt"
	"strings"

	turkpos "github.com/OzqurYalcin/turkpos/src"
	uuid "github.com/google/uuid"
)

func main() {
	api := &turkpos.API{"T"} // "T": test, "P": production
	request := &turkpos.Request{}
	request.Body.Payment.G.ClientCode = "10738"    // Müşteri No
	request.Body.Payment.G.ClientUsername = "Test" // Kullanıcı adı
	request.Body.Payment.G.ClientPassword = "Test" // Şifre
	// Ödeme
	commission := 0.0094
	amount := 100.00
	request.Body.Payment.GUID = "0c13d406-873b-403b-9c09-a5766840d98c"
	request.Body.Payment.Security = "NS"
	request.Body.Payment.OrderID = uuid.New().String()
	request.Body.Payment.PosID = 1029 // yurtdışı 1023
	request.Body.Payment.Description = "Açıklama"
	request.Body.Payment.CardOwner = "Kart sahibi"
	request.Body.Payment.CardNumber = "4546711234567894"
	request.Body.Payment.CardMonth = "12"
	request.Body.Payment.CardYear = "2026"
	request.Body.Payment.CardCvc = "000"
	request.Body.Payment.GsmNumber = "5554443322"
	request.Body.Payment.Price = strings.ReplaceAll(fmt.Sprintf("%.2f", amount-(amount*commission)), ".", ",")
	request.Body.Payment.Amount = strings.ReplaceAll(fmt.Sprintf("%.2f", amount), ".", ",")
	request.Body.Payment.Installment = 1
	request.Body.Payment.IPAddr = "85.34.78.112"
	request.Body.Payment.Referer = "https://www.example.com/payment"
	request.Body.Payment.CallbackError = "https://www.example.com/payment"
	request.Body.Payment.CallbackSuccess = "https://www.example.com/payment"
	request.Body.Payment.Hash = ""

	response := api.Payment(request)
	fmt.Println(response.Body.Payment.Result.Message)
}
```
