package turkpos

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
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
	Body    struct {
		XMLName xml.Name `xml:"soap:Body,omitempty"`
		Payment struct {
			G struct {
				ClientCode     interface{} `xml:"CLIENT_CODE,omitempty"`
				ClientUsername interface{} `xml:"CLIENT_USERNAME,omitempty"`
				ClientPassword interface{} `xml:"CLIENT_PASSWORD,omitempty"`
			} `xml:"G,omitempty"`
			GUID            interface{} `xml:"GUID,omitempty"`
			POSID           interface{} `xml:"SanalPOS_ID,omitempty"`
			HASH            interface{} `xml:"Islem_Hash,omitempty"`
			Security        interface{} `xml:"Islem_Guvenlik_Tip,omitempty"`
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
	} `xml:"soap:Body,omitempty"`
}

type Response struct {
	XMLName xml.Name `xml:"soap:Envelope,omitempty"`
	Body    struct {
		Response struct {
			Result struct {
				Result        string `xml:"Banka_Sonuc_Kod,omitempty"`
				Code          string `xml:"Sonuc,omitempty"`
				Message       string `xml:"Sonuc_Str,omitempty"`
				URL           string `xml:"UCD_URL,omitempty"`
				TransactionID string `xml:"Islem_ID,omitempty"`
			} `xml:"TP_Islem_OdemeResult,omitempty"`
		} `xml:"TP_Islem_OdemeResponse,omitempty"`
	} `xml:"soap:Body,omitempty"`
}

func (api *API) Payment(request *Request) (response *Response) {
	response = new(Response)
	postdata, _ := xml.Marshal(request)
	res, err := http.Post(Modes[api.Mode]+"?op=TP_Islem_Odeme", "text/xml; charset=utf-8", strings.NewReader(strings.ToLower(xml.Header)+string(postdata)))
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	decoder.Decode(&response)
	return response
}
