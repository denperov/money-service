## API

### Get available accounts
```shell script
curl -X GET localhost:8080/accounts
```
response:
```json
{
  "accounts": [
    {
      "id": "bob123",
      "currency": "USD",
      "balance": "100.00"
    },
    {
      "id": "alice456",
      "currency": "USD",
      "balance": "0.01"
    }
  ]
}
```

### Get all payments
```shell script
curl -X GET localhost:8080/payments
```
response:
```json
{
  "payments": [
    {
      "direction": "outgoing",
      "account": "alice456",
      "to_account": "bob123",
      "amount": "100.00"
    },
    {
      "direction": "incoming",
      "account": "bob123",
      "from_account": "alice456",
      "amount": "100.00"
    }
  ]
}
```

### Send payment
```shell script
curl -X POST localhost:8080/transfers -d '{"transfer":{"from_account":"bob123","to_account":"alice456","amount":"100.01"}}'
```
response:
```json
{}
```

### Errors
#### No errors
HTTP code: 200

The client receives a response from method.

#### User errors
HTTP code: 400

Error description response:
```json
{
  "error": "user friendly error description"
}
```

#### Server errors
HTTP code: 500

Error description response, always the same:
```json
{
  "error": "server error"
}
```
