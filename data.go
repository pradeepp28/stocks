package stocks

type respData struct {
	Symbol         string `json:"symbol"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	CloseYesterday string `json:"close_yesterday"`
	Currency       string `json:"currency"`
	MarketCap      string `json:"market_cap"`
	Volume         string `json:"volume"`
	Timezone       string `json:"timezone"`
	TimezoneName   string `json:"timezone_name"`
	GmtOffset      string `json:"gmt_offset"`
	LastTradeTime  string `json:"last_trade_time"`
}

type data struct {
	respData
	StockExchangeShort string `json:"stock_exchange_short"`
}

type stockExchangeResponse struct {
	Message          string `json:"message"`
	SymbolsRequested int    `json:"symbols_requested"`
	SymbolsReturned  int    `json:"symbols_returned"`
	Datas            []data `json:"data"`
}
