package stocks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"time"
)

const (
	apiToken      = "xr7nTooqcSBSyzZQZqZyQf4AU7t6Nd1tRmAZ0r67qPrtAH2vpbpK8eNlUAhD"
	stockURL      = "https://www.worldtradingdata.com/api/v1/stock"
	clientTimeout = 30 * time.Second
)

type Server interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type server struct {
	apiToken         string
	stockExchangeURL string
}

func New() Server {
	return server{
		apiToken:         apiToken,
		stockExchangeURL: stockURL,
	}
}

func (s server) Get(w http.ResponseWriter, r *http.Request) {

	symbol, err := parseSymbol(r.URL.Path)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	urlString, err := stockExchangeURL(symbol)
	if err != nil {
		log.Printf("Failed to parse url, url: %s, error: %v", stockURL, err)
		fmt.Fprint(w, "Error! Couldn't read stocks url")
		return
	}
	var seResp stockExchangeResponse
	if err := httpGet(urlString, &seResp); err != nil {
		log.Printf("Failed to get data from stock exchange, error: %v", err)
		fmt.Fprint(w, "Error! Failed to get data from stock exchange")
		return
	}

	if seResp.Message != "" && len(seResp.Datas) == 0 {
		log.Printf("Message from exchange server, %s", seResp.Message)
		fmt.Fprint(w, seResp.Message)
		return
	}

	if len(seResp.Datas) == 0 {
		fmt.Fprint(w, "Error! No stock found for given symbol")
		return
	}

	var stockExchange string
	if se := r.URL.Query().Get("stock_exchange"); se != "" {
		stockExchange = se
	}

	resp := formatResponse(stockExchange, seResp.Datas)

	if stockExchange != "" && len(resp) == 0 {
		fmt.Fprint(w, "Error! No data found for given stock exchange")
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func parseSymbol(urlPath string) (string, error) {
	re := regexp.MustCompile("^/stock/[A-Z]+[.]*[A-Z]+$")
	if !re.MatchString(urlPath) {
		return "", errors.New("Error incorrect url, expected format: /stock/{symbol}")
	}
	return path.Base(urlPath), nil
}

func stockExchangeURL(symbol string) (string, error) {
	u, err := url.Parse(stockURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add("symbol", symbol)
	q.Add("api_token", apiToken)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func httpGet(url string, resp interface{}) error {
	if reflect.TypeOf(resp).Elem().Kind() != reflect.Struct || reflect.TypeOf(resp).Kind() != reflect.Ptr {
		return errors.New("resp should be of type pointer to struct")
	}

	hr, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{Timeout: clientTimeout}
	eResp, err := httpClient.Do(hr)
	if err != nil {
		return err
	}
	defer eResp.Body.Close()
	if err := json.NewDecoder(eResp.Body).Decode(resp); err != nil {
		return err
	}
	return nil
}

func formatResponse(stockExchange string, datas []data) map[string]respData {
	resp := make(map[string]respData)
	for _, d := range datas {
		if stockExchange != "" && d.StockExchangeShort == stockExchange {
			resp[d.StockExchangeShort] = d.respData
		} else if stockExchange == "" {
			resp[d.StockExchangeShort] = d.respData
		}
	}
	return resp
}
