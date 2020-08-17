package turkpos

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var Modes map[string]string = map[string]string{
	"T": "https://test-dmz.param.com.tr:4443/turkpos.ws/service_turkpos_test.asmx",
	"P": "https://dmzws.param.com.tr/turkpos.ws/service_turkpos_prod.asmx",
}

type API struct {
	Mode string
}

type PaymentRequest struct {
	XMLName xml.Name `xml:"soap:Envelope,omitempty"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    struct {
		XMLName xml.Name `xml:"soap:Body,omitempty"`
		Payment struct {
			NS string `xml:"xmlns,attr"`
			G  struct {
				ClientCode     string `xml:"CLIENT_CODE,omitempty"`
				ClientUsername string `xml:"CLIENT_USERNAME,omitempty"`
				ClientPassword string `xml:"CLIENT_PASSWORD,omitempty"`
			} `xml:"G,omitempty"`
			GUID            string `xml:"GUID,omitempty"`
			Hash            string `xml:"Islem_Hash,omitempty"`
			PosID           int    `xml:"SanalPOS_ID,omitempty"`
			OrderID         string `xml:"Siparis_ID,omitempty"`
			Description     string `xml:"Siparis_Aciklama,omitempty"`
			CardOwner       string `xml:"KK_Sahibi,omitempty"`
			CardNumber      string `xml:"KK_No,omitempty"`
			CardMonth       string `xml:"KK_SK_Ay,omitempty"`
			CardYear        string `xml:"KK_SK_Yil,omitempty"`
			CardCvc         string `xml:"KK_CVC,omitempty"`
			GsmNumber       string `xml:"KK_Sahibi_GSM,omitempty"`
			Price           string `xml:"Islem_Tutar,omitempty"`
			Amount          string `xml:"Toplam_Tutar,omitempty"`
			Installment     int    `xml:"Taksit,omitempty"`
			IPAddr          string `xml:"IPAdr,omitempty"`
			Referer         string `xml:"Ref_URL,omitempty"`
			CallbackError   string `xml:"Hata_URL,omitempty"`
			CallbackSuccess string `xml:"Basarili_URL,omitempty"`
			Data1           string `xml:"Data1,omitempty"`
			Data2           string `xml:"Data2,omitempty"`
			Data3           string `xml:"Data3,omitempty"`
			Data4           string `xml:"Data4,omitempty"`
			Data5           string `xml:"Data5,omitempty"`
		} `xml:"TP_Islem_Odeme,omitempty"`
	}
}

type EncryptRequest struct {
	XMLName xml.Name `xml:"soap:Envelope,omitempty"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    struct {
		XMLName xml.Name `xml:"soap:Body,omitempty"`
		Encrypt struct {
			NS   string `xml:"xmlns,attr"`
			Data string `xml:"Data,omitempty"`
		} `xml:"SHA2B64,omitempty"`
	}
}

type Response struct {
	XMLName xml.Name
	Body    struct {
		XMLName xml.Name
		Payment struct {
			Result struct {
				URL           string `xml:"UCD_URL,omitempty"`
				BankCode      int    `xml:"Banka_Sonuc_Kod,omitempty"`
				Code          int    `xml:"Sonuc,omitempty"`
				Message       string `xml:"Sonuc_Str,omitempty"`
				TransactionID int64  `xml:"Islem_ID,omitempty"`
			} `xml:"TP_Islem_OdemeResult,omitempty"`
		} `xml:"TP_Islem_OdemeResponse,omitempty"`

		Encrypt struct {
			Result string `xml:"SHA2B64Result,omitempty"`
		} `xml:"SHA2B64Response,omitempty"`
	}
}

func (api *API) Payment(request *PaymentRequest) (response *Response) {
	hash := &EncryptRequest{}
	hash.Body.Encrypt.Data = request.Body.Payment.G.ClientCode + request.Body.Payment.GUID + strconv.Itoa(request.Body.Payment.PosID) + strconv.Itoa(request.Body.Payment.Installment) + request.Body.Payment.Price + request.Body.Payment.Amount + request.Body.Payment.OrderID + request.Body.Payment.CallbackError + request.Body.Payment.CallbackSuccess
	encrypt := api.Encrypt(hash)
	request.Body.Payment.Hash = encrypt.Body.Encrypt.Result
	request.Body.Payment.NS = "https://turkpos.com.tr/"
	request.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	response = new(Response)
	postdata, _ := xml.Marshal(request)
	res, err := http.Post(Modes[api.Mode]+"?op=TP_Islem_Odeme", "text/xml; charset=utf-8", strings.NewReader(strings.ToLower(xml.Header)+string(postdata)))
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := xml.NewDecoder(res.Body)
	decoder.Decode(&response)
	return response
}

func (api *API) Encrypt(request *EncryptRequest) (response *Response) {
	request.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	request.Body.Encrypt.NS = "https://turkpos.com.tr/"
	response = new(Response)
	postdata, _ := xml.Marshal(request)
	res, err := http.Post(Modes[api.Mode]+"?op=SHA2B64", "text/xml; charset=utf-8", strings.NewReader(strings.ToLower(xml.Header)+string(postdata)))
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := xml.NewDecoder(res.Body)
	decoder.Decode(&response)
	return response
}
