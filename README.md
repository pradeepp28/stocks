# stocks

stocks server uses https://www.worldtradingdata.com api to get the stocks details.

```
curl -X GET http://127.0.0.1:8080/stock/{symbol} -H 'authtoken: <token>'
```

**symbol**: 
stock exchange symbol (ex: AAPL)

**authToken**:
stocks server doesn't have authentication mechanism, 
instead it sends the given authToken to https://www.worldtradingdata.com as a url encoded value.
sign up (https://www.worldtradingdata.com) and get your personal API token for free 