#!/usr/bin/env bash

api_url=${SERVICE_API_URL:-'http://127.0.0.1:8080'}
echo "${api_url}"

api_post() { \
	printf '\nrequest:  %s %s\n' "$1" "$(echo "$2" | jq -c -C '.')" >&2; \
	status=$(curl -w "%{http_code}" -o "./pipe_for_curl_out" -s "${api_url}$1" -H 'Content-Type: application/json' -d "$2"); \
	response=$(cat "./pipe_for_curl_out"); \
	rm ./pipe_for_curl_out; \
 	printf 'response: %s %s\n' "${status}" "$(echo "${response}" | jq -c -C '.')" >&2; \
}

response() { \
	echo "${response}" | jq -r "${1}"
}

assert_status() { \
	if [ "${status}" -eq "$1" ]
	then echo -e '\033[0;32mOK\033[0m'
	else echo -e '\033[0;31mUnexpected status: '"${status}"'\033[0m'
	fi
}

assert_ok() {
	assert_status 200
}
assert_user_error() {
	assert_status 400
}
assert_server_error() {
	assert_status 500
}
#------------------------------------------------------------------------------

printf "wait for the service"
sleep 10

# I wand to be able to send a payment from one account to another (same currency)
echo ------------------------------------------------------------------------------
echo Transfer monye USD
api_post "/send_payment" '{"transfer":{"from_account":"usd01","to_account":"usd02","amount":"12.34"}}'
assert_ok

echo ------------------------------------------------------------------------------
echo Transfer monye RUB
api_post "/send_payment" '{"transfer":{"from_account":"rub01","to_account":"rub02","amount":"12.34"}}'
assert_ok

echo ------------------------------------------------------------------------------
echo Transfer all monye from the account
api_post "/send_payment" '{"transfer":{"from_account":"usd03","to_account":"usd04","amount":"100.00"}}'
assert_ok

echo ------------------------------------------------------------------------------
echo Transfer zero monye
api_post "/send_payment" '{"transfer":{"from_account":"usd05","to_account":"usd06","amount":"0.00"}}'
assert_user_error

# I want to be able to see all payments
echo ------------------------------------------------------------------------------
echo Get payments
api_post "/get_payments" '{}'
assert_ok

# I want to be able to see available accounts
echo ------------------------------------------------------------------------------
echo Get accounts
api_post "/get_accounts" '{}'
assert_ok

# Only payments within the same currency are supported (no exchanges)
echo ------------------------------------------------------------------------------
echo Transfer money between accounts with different currencies
api_post "/send_payment" '{"transfer":{"from_account":"usd07","to_account":"rub03","amount":"12.34"}}'
assert_user_error

# Balance can't go below zero
echo ------------------------------------------------------------------------------
echo Transfer more money than there is on the account
api_post "/send_payment" '{"transfer":{"from_account":"usd08","to_account":"usd09","amount":"100.01"}}'
assert_user_error

# Nonexistent source account
echo ------------------------------------------------------------------------------
echo Transfer monye
api_post "/send_payment" '{"transfer":{"from_account":"xxx","to_account":"usd02","amount":"12.34"}}'
assert_user_error

# Nonexistent destination account
echo ------------------------------------------------------------------------------
echo Transfer monye
api_post "/send_payment" '{"transfer":{"from_account":"usd01","to_account":"xxx","amount":"12.34"}}'
assert_user_error
