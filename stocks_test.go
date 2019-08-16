package stocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSymbol(t *testing.T) {
	testData := []struct {
		urlPath      string
		expectSymbol string
		expectErr    string
	}{
		{urlPath: "/stock/AAPL", expectSymbol: "AAPL"},
		{urlPath: "/stock/HSBA.L", expectSymbol: "HSBA.L"},
		{urlPath: "/stock/AAPL/NASDAQ", expectErr: "Error incorrect url, expected format: /stock/{symbol}"},
		{urlPath: "/stock/", expectErr: "Error incorrect url, expected format: /stock/{symbol}"},
	}
	for _, d := range testData {
		s, err := parseSymbol(d.urlPath)
		if d.expectSymbol != "" {
			assert.Equal(t, s, d.expectSymbol, fmt.Sprintf("Failed to parse url: %s, err: %v", d.urlPath, err))
		}
		if d.expectErr != "" {
			assert.Equal(t, err.Error(), d.expectErr, "Error not matching expected error")
		}
	}

}

func TestStockExchangeURL(t *testing.T) {

	testData := []struct {
		symbol string
		token  string
	}{
		{symbol: "AAPL", token: "token1"},
		{symbol: "HSBA.L", token: "token2"},
		{symbol: "", token: ""},
	}

	for _, d := range testData {

		if seURL, err := stockExchangeURL(d.symbol, d.token); err != nil {
			t.Error("Failed to create stock exchange url")
		} else {
			u, err := url.Parse(seURL)
			if err != nil {
				t.Error(err)
			} else {
				assert.Equal(t, u.Query().Get("symbol"), d.symbol, "Query doesn't have expected symbol")
				assert.Equal(t, u.Query().Get("api_token"), d.token, "Query doesn't have expected api_token")
			}
		}
	}
}

func TestFormatResponse(t *testing.T) {
	rd := respData{
		Symbol:   "AAPL",
		Name:     "Apple Inc.",
		Timezone: "EDT",
	}
	ds := []data{
		{
			respData:           rd,
			StockExchangeShort: "NASDAQ",
		},
	}

	expectedResp := map[string]respData{
		"NASDAQ": rd,
	}

	testData := []struct {
		se           string
		expectedResp map[string]respData
	}{
		{se: "", expectedResp: expectedResp},
		{se: "NASDAQ", expectedResp: expectedResp},
		{se: "AAA", expectedResp: make(map[string]respData)},
	}

	for _, d := range testData {
		resp := formatResponse(d.se, ds)
		assert.Equal(t, reflect.DeepEqual(d.expectedResp, resp), true, "Response doesn't match expected response")
	}
}

func TestHttpGet(t *testing.T) {
	type foo struct {
		Data string `json:"data"`
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": "foo-data"}`))
	}))

	testCases := []struct {
		name   string
		output foo
		error  bool
	}{
		{
			name:   "Should return expected output",
			output: foo{Data: "foo-data"},
			error:  false,
		},
	}

	for _, tc := range testCases {
		var f foo
		err := httpGet(ts.URL, &f)
		if tc.error {
			assert.Error(t, err, "Expect httpGet() to return error")
		} else {
			assert.NoError(t, err, "Expect httpGet() to not return error")
			assert.Equal(t, tc.output, f, "Expect httpGet() to load expected data")
		}
	}

}
