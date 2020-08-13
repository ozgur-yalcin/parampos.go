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
