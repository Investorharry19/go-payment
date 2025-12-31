package paystack

import (
	"net/http"
	"time"
)

type PaystackClient struct {
	secretKey string
	baseURL   string
	http      *http.Client
}

func NewPaystackClient(secretKey string) *PaystackClient {
	return &PaystackClient{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
