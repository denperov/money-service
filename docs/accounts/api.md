## API

### Get available accounts
```shell script
curl -XPOST localhost:8080/get_accounts -d '{}'
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
curl -XPOST localhost:8080/get_payments -d '{}'
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
curl -XPOST localhost:8080/send_payment -d '{"transfer":{"from_account":"bob123","to_account":"alice456","amount":"100.01"}}'
```
response:
```json
{}
```

### Errors
#### No errors
HTTP code: 200

The client receives a response for method.

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