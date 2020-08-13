[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/OzqurYalcin/turkpos/blob/master/LICENSE.md)

# TurkPos
TurkPos (ParamPos) API with golang

# Installation
```bash
go get github.com/OzqurYalcin/turkpos
```

# Sanalpos satış işlemi
```go
package main

import (
	turkpos "github.com/OzqurYalcin/turkpos/src"
)

func main() {
	api := &turkpos.API{"T"} // "T","P"
	request := &turkpos.Request{}
	request.Body.Payment.G.ClientCode = "10738"    // Müşteri No
	request.Body.Payment.G.ClientUsername = "Test" // Kullanıcı adı
	request.Body.Payment.G.ClientPassword = "Test" // Şifre
	// Ödeme

	response := api.Payment(request)
	if response.Body.Response.Result.TransactionID != "" {

	}
}
```
