package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/denperov/money-service/internal/accounts/endpoints"
	"github.com/denperov/money-service/internal/pkg/user_errors"
)

func MakeGetAccountsHandler(e endpoint.Endpoint) *httptransport.Server {
	return httptransport.NewServer(
		e,
		decodeGetAccountsRequest,
		encodeResponse,
		errorEncoderOption,
	)
}

func MakeGetPaymentsHandler(e endpoint.Endpoint) *httptransport.Server {
	return httptransport.NewServer(
		e,
		decodeGetPaymentsRequest,
		encodeResponse,
		errorEncoderOption,
	)
}
func MakeSendPaymentHandler(e endpoint.Endpoint) *httptransport.Server {
	return httptransport.NewServer(
		e,
		decodeSendPaymentRequest,
		encodeResponse,
		errorEncoderOption,
	)
}

func decodeGetAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.GetAccountsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetPaymentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.GetPaymentsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSendPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.SendPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

var errorEncoderOption = httptransport.ServerErrorEncoder(encodeError)

func encodeError(_ context.Context, serverError error, w http.ResponseWriter) {
	var errorText string
	if user_errors.IsUserFriendlyError(serverError) {
		w.WriteHeader(http.StatusBadRequest)
		errorText = serverError.Error()
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		errorText = "server error"
	}

	err := json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: errorText,
	})
	if err != nil {
		log.Print(err)
	}
}
