package parampos

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"html"
	"log"
	"net/http"
	"strings"
)

var EndPoints map[string]string = map[string]string{
	"TEST": "https://test-dmz.param.com.tr:4443/turkpos.ws/service_turkpos_test.asmx",
	"PROD": "https://dmzws.param.com.tr/turkpos.ws/service_turkpos_prod.asmx",
}

var Currencies map[string]int = map[string]int{
	"TRY": 1000,
	"YTL": 1000,
	"TRL": 1000,
	"TL":  1000,
	"USD": 1001,
	"EUR": 1002,
}

type API struct {
	Mode string
}

type Request struct {
	XMLName xml.Name `xml:"soap:Envelope,omitempty"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    struct {
		XMLName   xml.Name `xml:"soap:Body,omitempty"`
		Hash      *HASH    `xml:"SHA2B64,omitempty"`
		Auth      *Payment `xml:"TP_WMD_UCD,omitempty"`
		Pay       *Payment `xml:"TP_WMD_Pay,omitempty"`
		PreAuth   *Payment `xml:"TP_Islem_Odeme_OnProv_WMD,omitempty"`
		PostAuth  *Payment `xml:"TP_Islem_Odeme_OnProv_Kapa,omitempty"`
		Cancel    *Payment `xml:"TP_Islem_Iptal_Iade_Kismi2,omitempty"`
		PreCancel *Payment `xml:"TP_Islem_Iptal_OnProv,omitempty"`
	}
}

type Response struct {
	XMLName xml.Name
	Body    struct {
		XMLName xml.Name
		Auth    *struct {
			Result Result `xml:"TP_WMD_UCDResult,omitempty"`
		} `xml:"TP_WMD_UCDResponse,omitempty"`
		Pay *struct {
			Result Result `xml:"TP_WMD_PayResult,omitempty"`
		} `xml:"TP_WMD_PayResponse,omitempty"`
		PreAuth *struct {
			Result Result `xml:"TP_Islem_Odeme_OnProv_WMDResult,omitempty"`
		} `xml:"TP_Islem_Odeme_OnProv_WMDResponse,omitempty"`
		PostAuth *struct {
			Result Result `xml:"TP_Islem_Odeme_OnProv_KapaResult,omitempty"`
		} `xml:"TP_Islem_Odeme_OnProv_KapaResponse,omitempty"`
		Cancel *struct {
			Result Result `xml:"TP_Islem_Iptal_Iade_Kismi2Result,omitempty"`
		} `xml:"TP_Islem_Iptal_Iade_Kismi2Response,omitempty"`
		PreCancel *struct {
			Result Result `xml:"TP_Islem_Iptal_OnProvResult,omitempty"`
		} `xml:"TP_Islem_Iptal_OnProvResponse,omitempty"`
		Encrypt *struct {
			Result string `xml:"SHA2B64Result,omitempty"`
		} `xml:"SHA2B64Response,omitempty"`
		Fault *struct {
			XMLName xml.Name
			Code    string `xml:"faultcode,omitempty"`
			String  string `xml:"faultstring,omitempty"`
		}
	}
}

type Result struct {
	UCD             string `xml:"UCD_MD,omitempty"`
	URL             string `xml:"UCD_URL,omitempty"`
	HTML            string `xml:"UCD_HTML,omitempty"`
	BankAuthCode    string `xml:"Bank_AuthCode,omitempty"`
	BankHostRefNum  string `xml:"Bank_HostRefNum,omitempty"`
	BankHostMsg     string `xml:"Bank_HostMsg,omitempty"`
	BankCode        int64  `xml:"Banka_Sonuc_Kod,omitempty"`
	BankTransID     string `xml:"Bank_Trans_ID,omitempty"`
	Code            int64  `xml:"Sonuc,omitempty"`
	Message         string `xml:"Sonuc_Str,omitempty"`
	Description     string `xml:"Sonuc_Ack,omitempty"`
	OrderId         string `xml:"Siparis_ID,omitempty"`
	ReceiptID       string `xml:"Dekont_ID,omitempty"`
	TransactionID   int64  `xml:"Islem_ID,omitempty"`
	TransactionGUID string `xml:"Islem_GUID,omitempty"`
}

type Payment struct {
	NS string `xml:"xmlns,attr"`
	G  struct {
		ClientCode     string `xml:"CLIENT_CODE,omitempty"`
		ClientUsername string `xml:"CLIENT_USERNAME,omitempty"`
		ClientPassword string `xml:"CLIENT_PASSWORD,omitempty"`
	} `xml:"G,omitempty"`
	UCD             string `xml:"UCD_MD,omitempty"`
	GUID            string `xml:"GUID,omitempty"`
	TransactionGUID string `xml:"Islem_GUID,omitempty"`
	Hash            string `xml:"Islem_Hash,omitempty"`
	Type            string `xml:"Islem_Guvenlik_Tip,omitempty"`
	TransactionID   string `xml:"Islem_ID,omitempty"`
	ProvId          string `xml:"Prov_ID,omitempty"`
	ProvAmount      string `xml:"Prov_Tutar,omitempty"`
	OrderId         string `xml:"Siparis_ID,omitempty"`
	Description     string `xml:"Siparis_Aciklama,omitempty"`
	CardHolder      string `xml:"KK_Sahibi,omitempty"`
	CardNumber      string `xml:"KK_No,omitempty"`
	CardMonth       string `xml:"KK_SK_Ay,omitempty"`
	CardYear        string `xml:"KK_SK_Yil,omitempty"`
	CardCode        string `xml:"KK_CVC,omitempty"`
	GsmNumber       string `xml:"KK_Sahibi_GSM,omitempty"`
	Price           string `xml:"Islem_Tutar,omitempty"`
	Total           string `xml:"Toplam_Tutar,omitempty"`
	Amount          string `xml:"Tutar,omitempty"`
	Installment     string `xml:"Taksit,omitempty"`
	IPAddr          string `xml:"IPAdr,omitempty"`
	Referer         string `xml:"Ref_URL,omitempty"`
	CallbackError   string `xml:"Hata_URL,omitempty"`
	CallbackSuccess string `xml:"Basarili_URL,omitempty"`
	Action          string `xml:"Durum,omitempty"`
	Data1           string `xml:"Data1,omitempty"`
	Data2           string `xml:"Data2,omitempty"`
	Data3           string `xml:"Data3,omitempty"`
	Data4           string `xml:"Data4,omitempty"`
	Data5           string `xml:"Data5,omitempty"`
}

type HASH struct {
	NS   string `xml:"xmlns,attr"`
	Data string `xml:"Data,omitempty"`
}

func IPv4(r *http.Request) (ip string) {
	ipv4 := []string{
		r.Header.Get("X-Real-Ip"),
		r.Header.Get("X-Forwarded-For"),
		r.RemoteAddr,
	}
	for _, ipaddress := range ipv4 {
		if ipaddress != "" {
			ip = ipaddress
			break
		}
	}
	return strings.Split(ip, ":")[0]
}

func HEX(data string) (hash string) {
	b, err := hex.DecodeString(data)
	if err != nil {
		log.Println(err)
		return hash
	}
	hash = string(b)
	return hash
}

func SHA1(data string) (hash string) {
	h := sha1.New()
	h.Write([]byte(data))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func B64(data string) (hash string) {
	hash = base64.StdEncoding.EncodeToString([]byte(data))
	return hash
}

func D64(data string) []byte {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return b
}

func Hash(data string) string {
	return B64(HEX(SHA1(data)))
}

func Api(clientid, username, password string) (*API, *Payment) {
	api := new(API)
	payment := new(Payment)
	payment.G.ClientCode = clientid
	payment.G.ClientUsername = username
	payment.G.ClientPassword = password
	return api, payment
}

func (api *API) SetMode(mode string) {
	api.Mode = mode
}

func (payment *Payment) SetGUID(guid string) {
	payment.GUID = guid
}

func (payment *Payment) SetIPAddress(ip string) {
	payment.IPAddr = ip
}

func (payment *Payment) SetPhoneNumber(gsm string) {
	payment.GsmNumber = gsm
}

func (payment *Payment) SetCardHolder(holder string) {
	payment.CardHolder = holder
}

func (payment *Payment) SetCardNumber(number string) {
	payment.CardNumber = number
}

func (payment *Payment) SetCardExpiry(month, year string) {
	payment.CardMonth = month
	payment.CardYear = "20" + year
}

func (payment *Payment) SetCardCode(code string) {
	payment.CardCode = code
}

func (payment *Payment) SetAmount(total string) {
	payment.Price = strings.ReplaceAll(total, ".", ",")
	payment.Total = strings.ReplaceAll(total, ".", ",")
}

func (payment *Payment) SetInstallment(ins string) {
	payment.Installment = ins
}

func (payment *Payment) SetOrderId(oid string) {
	payment.OrderId = oid
}

func (api *API) AuthHash(payment *Payment) (string, error) {
	req := new(Request)
	req.Body.Hash = new(HASH)
	req.Body.Hash.Data = payment.G.ClientCode + payment.GUID + payment.Installment + payment.Price + payment.Total + payment.OrderId
	ctx := context.Background()
	res, err := api.Encrypt(ctx, req)
	return res.Body.Encrypt.Result, err
}

func (api *API) PreAuthHash(payment *Payment) (string, error) {
	req := new(Request)
	req.Body.Hash = new(HASH)
	req.Body.Hash.Data = payment.G.ClientCode + payment.GUID + payment.Price + payment.Total + payment.OrderId + payment.CallbackError + payment.CallbackSuccess
	ctx := context.Background()
	res, err := api.Encrypt(ctx, req)
	return res.Body.Encrypt.Result, err
}

func (api *API) Encrypt(ctx context.Context, req *Request) (res Response, err error) {
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	req.Body.Hash.NS = "https://turkpos.com.tr/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	if res.Body.Encrypt.Result != "" {
		return res, nil
	} else {
		return res, errors.New("empty response")
	}
}

func (api *API) PreAuth(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.Hash, err = api.PreAuthHash(payment)
	if err != nil {
		return res, err
	}
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.PreAuth = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.PreAuth.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.PreAuth.Result.Message)
	}
}

func (api *API) Auth(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.Hash, err = api.AuthHash(payment)
	if err != nil {
		return res, err
	}
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.Auth = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.Auth.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.Auth.Result.Message)
	}
}

func (api *API) PreAuth3D(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.Hash, err = api.PreAuthHash(payment)
	if err != nil {
		return res, err
	}
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.Pay = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.Pay.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.Pay.Result.Message)
	}
}

func (api *API) Auth3D(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.Hash, err = api.AuthHash(payment)
	if err != nil {
		return res, err
	}
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.Pay = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.Pay.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.Pay.Result.Message)
	}
}

func (api *API) PreAuth3Dhtml(ctx context.Context, payment *Payment) (res string, err error) {
	payment.Hash, err = api.PreAuthHash(payment)
	if err != nil {
		return res, err
	}
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.PreAuth = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	data := Response{}
	if err := decoder.Decode(&data); err != nil {
		return res, err
	}
	switch data.Body.PreAuth.Result.Code {
	case 1:
		res = B64(html.UnescapeString(data.Body.PreAuth.Result.HTML))
		return res, nil
	default:
		return res, errors.New(data.Body.PreAuth.Result.Message)
	}
}

func (api *API) Auth3Dhtml(ctx context.Context, payment *Payment) (res string, err error) {
	payment.Hash, err = api.AuthHash(payment)
	if err != nil {
		return res, err
	}
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.Auth = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	data := Response{}
	if err := decoder.Decode(&data); err != nil {
		return res, err
	}
	switch data.Body.Auth.Result.Code {
	case 1:
		res = B64(html.UnescapeString(data.Body.Auth.Result.HTML))
		return res, nil
	default:
		return res, errors.New(data.Body.Auth.Result.Message)
	}
}

func (api *API) PostAuth(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.ProvAmount = strings.ReplaceAll(payment.Total, ",", ".")
	payment.Total = ""
	payment.Price = ""
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.PostAuth = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.PostAuth.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.PostAuth.Result.Message)
	}
}

func (api *API) Refund(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.Amount = strings.ReplaceAll(payment.Total, ",", ".")
	payment.Action = "IADE"
	payment.Total = ""
	payment.Price = ""
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.Cancel = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.Cancel.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.Cancel.Result.Message)
	}
}

func (api *API) Cancel(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.Amount = strings.ReplaceAll(payment.Total, ",", ".")
	payment.Action = "IPTAL"
	payment.Total = ""
	payment.Price = ""
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.Cancel = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.Cancel.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.Cancel.Result.Message)
	}
}

func (api *API) PreCancel(ctx context.Context, payment *Payment) (res Response, err error) {
	payment.NS = "https://turkpos.com.tr/"
	req := new(Request)
	req.Body.PreCancel = payment
	req.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.Body.PreCancel.Result.Code {
	case 1:
		return res, nil
	default:
		return res, errors.New(res.Body.PreCancel.Result.Message)
	}
}
