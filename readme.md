# Currency rates service

Service provides an async interface, where the user first makes a request to update the currency rate, and then, after some time, requests the rate. 

To run service locally run from the root:
```
make start
```

Swagger docs: http://localhost:8080/api/v1/swagger/index.html#/

Base currency is `USD`. The following currencies are supported: `USD`, `EUR`, `MXN`. Info is provided by [VAT comply](https://www.vatcomply.com/documentation#rates-base). Currency rates tracks foreign exchange references rates published by the European Central Bank.
The data refreshes around 16:00 CET every working day. 

1. Send a PUT request to initiate an exchange rate update and receive a UUID. Usually it takes about a minute for the update to complete. Optionally, specify the base currency with `?base={base}`.
```bash
curl -X PUT http://localhost:8080/api/v1/rates/eur/update

{
  "Base": "USD",
  "Currency": "EUR",
  "UUID": "2c63ba98-0908-4598-8c7c-870b9d83f3e9"
}
```
2. Retrieve the updated exchange rate by calling `GET /rates?uuid={uuid_from_previous_step`}.
```bash
curl -X GET http://localhost:8080/api/v1/rates?uuid=2c63ba98-0908-4598-8c7c-870b9d83f3e9

{
  "updated_at": "2024-02-29 22:38:18.788931 +0000 UTC",
  "base": "USD",
  "currency": "EUR",
  "value": 0.9237021803855896
}
```
3. Get the latest currency rate from the database `GET /rates/{code}?base={base}`. Default base is `USD`.
```bash
curl -X GET http://localhost:8080/api/v1/rates?uuid=2c63ba98-0908-4598-8c7c-870b9d83f3e9

{
  "updated_at": "2024-02-28 20:00:28.525431 +0000 UTC",
  "base": "EUR",
  "currency": "MXN",
  "value": 18.480100631713867
}
```
### How update works
![](scheme.png)