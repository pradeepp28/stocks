package stocks

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/magiconair/properties/assert"
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

	testData := []string{"AAPL", "HSBA.L"}

	for _, d := range testData {

		if seURL, err := stockExchangeURL(d); err != nil {
			t.Error("Failed to create stock exchange url")
		} else {
			u, err := url.Parse(seURL)
			if err != nil {
				t.Error(err)
			} else {
				assert.Equal(t, u.Query().Get("symbol"), d, "Query doesn't have expected symbol")
				assert.Equal(t, u.Query().Get("api_token"), apiToken, "Query doesn't have expected api_token")
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
	resp1 := formatResponse("", ds)
	assert.Equal(t, reflect.DeepEqual(expectedResp, resp1), true, "Response doesn't matcg expected response")

	resp2 := formatResponse("NASDAQ", ds)
	assert.Equal(t, reflect.DeepEqual(expectedResp, resp2), true, "Response doesn't matcg expected response")
}
