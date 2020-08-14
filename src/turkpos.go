package turkpos

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var Modes map[string]string = map[string]string{
	"T": "https://test-dmz.param.com.tr:4443/turkpos.ws/service_turkpos_test.asmx",
	"P": "https://dmzws.param.com.tr/turkpos.ws/service_turkpos_prod.asmx",
}

type API struct {
	Mode string
}

type Request struct {
	XMLName xml.Name `xml:"soap:Envelope,omitempty"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    struct {
		XMLName xml.Name `xml:"soap:Body,omitempty"`
		Payment struct {
			NS string `xml:"xmlns,attr"`
			G  struct {
				ClientCode     interface{} `xml:"CLIENT_CODE,omitempty"`
				ClientUsername interface{} `xml:"CLIENT_USERNAME,omitempty"`
				ClientPassword interface{} `xml:"CLIENT_PASSWORD,omitempty"`
			} `xml:"G,omitempty"`
			GUID            interface{} `xml:"GUID,omitempty"`
			Hash            interface{} `xml:"Islem_Hash,omitempty"`
			Security        interface{} `xml:"Islem_Guvenlik_Tip,omitempty"`
			PosID           interface{} `xml:"SanalPOS_ID,omitempty"`
			OrderID         interface{} `xml:"Siparis_ID,omitempty"`
			Description     interface{} `xml:"Siparis_Aciklama,omitempty"`
			CardOwner       interface{} `xml:"KK_Sahibi,omitempty"`
			CardNumber      interface{} `xml:"KK_No,omitempty"`
			CardMonth       interface{} `xml:"KK_SK_Ay,omitempty"`
			CardYear        interface{} `xml:"KK_SK_Yil,omitempty"`
			CardCvc         interface{} `xml:"KK_CVC,omitempty"`
			GsmNumber       interface{} `xml:"KK_Sahibi_GSM,omitempty"`
			Price           interface{} `xml:"Islem_Tutar,omitempty"`
			Amount          interface{} `xml:"Toplam_Tutar,omitempty"`
			Installment     interface{} `xml:"Taksit,omitempty"`
			IPAddr          interface{} `xml:"IPAdr,omitempty"`
			Referer         interface{} `xml:"Ref_URL,omitempty"`
			CallbackError   interface{} `xml:"Hata_URL,omitempty"`
			CallbackSuccess interface{} `xml:"Basarili_URL,omitempty"`
			Data1           interface{} `xml:"Data1,omitempty"`
			Data2           interface{} `xml:"Data2,omitempty"`
			Data3           interface{} `xml:"Data3,omitempty"`
			Data4           interface{} `xml:"Data4,omitempty"`
			Data5           interface{} `xml:"Data5,omitempty"`
		} `xml:"TP_Islem_Odeme,omitempty"`
	}
}

type Response struct {
	XMLName xml.Name `xml:"TP_Islem_OdemeResponse,omitempty"`
	NS      string   `xml:"xmlns,attr"`
	Payment struct {
		XMLName       xml.Name    `xml:"TP_Islem_OdemeResult,omitempty"`
		TransactionID interface{} `xml:"Islem_ID,omitempty"`
		URL           interface{} `xml:"UCD_URL,omitempty"`
		Code          interface{} `xml:"Sonuc,omitempty"`
		Message       interface{} `xml:"Sonuc_Str,omitempty"`
		Result        interface{} `xml:"Banka_Sonuc_Kod,omitempty"`
	}
}

func (api *API) Payment(request *Request) (response Response) {
	request.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	request.Body.Payment.NS = "https://turkpos.com.tr/"
	response = Response{}
	postdata, _ := xml.Marshal(request)
	res, err := http.Post(Modes[api.Mode]+"?op=TP_Islem_Odeme", "text/xml; charset=utf-8", strings.NewReader(strings.ToLower(xml.Header)+string(postdata)))
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return response
	}
	replacers := strings.NewReplacer(
		`<?xml version="1.0" encoding="utf-8"?>`, ``,
		`<soap:Envelope xmlns:soap="`+request.Soap+`" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">`, ``,
		`</soap:Envelope>`, ``,
		`<soap:Body>`, ``,
		`</soap:Body>`, ``,
	)
	xmldata := strings.ToLower(xml.Header) + replacers.Replace(string(data))
	fmt.Println(xmldata)
	xml.Unmarshal([]byte(xmldata), &response)
	return response
}
