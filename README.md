[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/ozgur-yalcin/parampos.go/blob/master/LICENSE.md)
[![documentation](https://pkg.go.dev/badge/github.com/ozgur-yalcin/parampos.go)](https://pkg.go.dev/github.com/ozgur-yalcin/parampos.go/src)

# Parampos.go
ParamPos (TurkPos) API with golang

# Installation
```bash
go get github.com/ozgur-yalcin/parampos.go
```

# Satış
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	parampos "github.com/ozgur-yalcin/parampos.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                                 // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	clientid = "10738"                                // Müşteri numarası
	username = "Test"                                 // Kullanıcı adı
	password = "Test"                                 // Şifre
	storekey = "0c13d406-873b-403b-9c09-a5766840d98c" // İşyeri anahtarı (GUID)
)

func main() {
	api, req := parampos.Api(clientid, username, password)
	api.SetMode(envmode)

	req.SetStoreKey(storekey)
	req.SetIPAddress("1.2.3.4")           // Müşteri ip adresi (zorunlu)
	req.SetCardHolder("AD SOYAD")         // Kart sahibi (zorunlu)
	req.SetCardNumber("4546711234567894") // Kart numarası (zorunlu)
	req.SetCardExpiry("12", "26")         // Son kullanma tarihi - AA,YY (zorunlu)
	req.SetCardCode("000")                // Kart arkasındaki 3 haneli numara (zorunlu)
	req.SetPhoneNumber("5554443322")      // Müşteri cep telefonu (zorunlu)
	req.SetAmount("1.00")                 // Satış tutarı (zorunlu)
	req.SetInstallment("1")               // Taksit sayısı (zorunlu)

	// Satış
	ctx := context.Background()
	if res, err := api.Auth(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İade
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	parampos "github.com/ozgur-yalcin/parampos.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                                 // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	clientid = "10738"                                // Müşteri numarası
	username = "Test"                                 // Kullanıcı adı
	password = "Test"                                 // Şifre
	storekey = "0c13d406-873b-403b-9c09-a5766840d98c" // İşyeri anahtarı (GUID)
)

func main() {
	api, req := parampos.Api(clientid, username, password)
	api.SetMode(envmode)

	req.SetStoreKey(storekey)
	req.SetAmount("1.00") // İade tutarı (zorunlu)
	req.SetOrderId("")    // Sipariş numarası (zorunlu)

	// İade
	ctx := context.Background()
	if res, err := api.Refund(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İptal
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	parampos "github.com/ozgur-yalcin/parampos.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                                 // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	clientid = "10738"                                // Müşteri numarası
	username = "Test"                                 // Kullanıcı adı
	password = "Test"                                 // Şifre
	storekey = "0c13d406-873b-403b-9c09-a5766840d98c" // İşyeri anahtarı (GUID)
)

func main() {
	api, req := parampos.Api(clientid, username, password)
	api.SetMode(envmode)

	req.SetStoreKey(storekey)
	req.SetAmount("1.00") // İptal tutarı (zorunlu)
	req.SetOrderId("")    // Sipariş numarası (zorunlu)

	// İptal
	ctx := context.Background()
	if res, err := api.Cancel(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```
