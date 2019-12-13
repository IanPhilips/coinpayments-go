package coinpayments

import (
	"net/http"

	//"fmt"
	"strings"
	"github.com/dghubble/sling"
)

type TransactionService struct {
	sling        *sling.Sling
	ApiPublicKey string
	Params       TransactionBodyParams
}

type Transaction struct {
	Amount         string `json:"amount"`
	Address        string `url:"address"`
	TXNId          string `json:"txn_id"`
	ConfirmsNeeded string `json:"confirms_needed"`
	Timeout        uint32 `json:"timeout"`
	StatusUrl      string `json:"status_url"`
	QRCodeUrl      string `json:"qrcode_url"`
}

type TransactionResponse struct {
	Error  string       `json:"error"`
	Result *Transaction `json:"result"`
}

type TransactionResponseError struct {
	Error  string       `json:"error"`
	Result []*Transaction `json:"result"`
}

type TransactionParams struct {
	Amount     float64 `url:"amount"`
	Currency1  string  `url:"currency1"`
	Currency2  string  `url:"currency2"`
	Address    string  `url:"address"`
	BuyerEmail string  `url:"buyer_email"`
	BuyerName  string  `url:"buyer_name"`
	ItemName   string  `url:"item_name"`
	ItemNumber string  `url:"item_number"`
	Invoice    string  `url:"invoice"`
	Custom     string  `url:"custom"`
	IPNUrl     string  `url:"ipn_url"`
}

type TransactionBodyParams struct {
	APIParams
	TransactionParams
}

func newTransactionService(sling *sling.Sling, apiPublicKey string) *TransactionService {
	transactionService := &TransactionService{
		sling:        sling.Path("api.php"),
		ApiPublicKey: apiPublicKey,
	}
	transactionService.getParams()
	return transactionService
}

func (s *TransactionService) getHMAC() string {
	return getHMAC(getPayload(s.Params))
}

// TODO: if we generate an error in the sling POST request, we get the error:
// json: cannot unmarshal array into Go struct field TransactionResponse.result of type coinpayments.Transaction"
// if you change TransactionResponse.Result to be of type []*Transaction then you can see the error message
// but when we don't get an error, it is a single *Transaction and the slice will throw an error
func (s *TransactionService) NewTransaction(transactionParams *TransactionParams) (TransactionResponse, *http.Response, error) {
	transactionResponse := new(TransactionResponse)
	s.Params.TransactionParams = *transactionParams
	//fmt.Println(getPayload(s.Params))
	//fmt.Println(getHMAC(getPayload(s.Params)))

	resp, err := s.sling.New().Set("HMAC", s.getHMAC()).Post(
		"api.php").BodyForm(s.Params).ReceiveSuccess(transactionResponse)

	// a bit of a hack - we should use a different post method, but catches errors with minimal effort
	if err!= nil && strings.Contains(err.Error(), "json"){
		transactionResponseError := new(TransactionResponseError)
		resp, err := s.sling.New().Set("HMAC", s.getHMAC()).Post(
			"api.php").BodyForm(s.Params).ReceiveSuccess(transactionResponseError)
		transactionResponse.Error = transactionResponseError.Error
		return *transactionResponse, resp, err
	}

	return *transactionResponse, resp, err
}




func (s *TransactionService) getParams() {
	s.Params.Command = "create_transaction"
	s.Params.Key = s.ApiPublicKey
	s.Params.Version = "1"
}
