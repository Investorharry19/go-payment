package paystack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Investorharry19/go-payment/internal/payment"
)

type initializeRequest struct {
	Email       string `json:"email"`
	Amount      int64  `json:"amount"`
	CallbackURL string `json:"callback_url"`
	Reference   string `json:"reference"`
}

type initializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

func (p *PaystackClient) Authorize(
	ctx context.Context,
	req payment.AuthorizeRequest,
) (payment.AuthorizeResponse, error) {

	payload := initializeRequest{
		Email:       req.Email,
		Amount:      req.Amount,
		CallbackURL: req.CallbackURL,
		Reference:   req.PaymentID, // your internal payment ID
	}

	fmt.Println(payload.Email)

	body, err := json.Marshal(payload)
	if err != nil {
		return payment.AuthorizeResponse{}, fmt.Errorf("marshal paystack request: %w", err)
	}
	fmt.Println(string(body))

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseURL+"/transaction/initialize",
		bytes.NewReader(body),
	)
	if err != nil {
		return payment.AuthorizeResponse{}, fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotency-Key", req.OperationID)

	resp, err := p.http.Do(httpReq)
	if err != nil {
		return payment.AuthorizeResponse{}, fmt.Errorf("paystack request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return payment.AuthorizeResponse{}, fmt.Errorf("paystack returned status %d, failed to decode error: %w", resp.StatusCode, err)
		}
		return payment.AuthorizeResponse{}, fmt.Errorf("paystack error: %s", errResp.Message)
	}

	var psResp initializeResponse
	if err := json.NewDecoder(resp.Body).Decode(&psResp); err != nil {
		return payment.AuthorizeResponse{}, fmt.Errorf("decode paystack response: %w", err)
	}

	if !psResp.Status {
		return payment.AuthorizeResponse{}, fmt.Errorf("paystack error: %s", psResp.Message)
	}

	return payment.AuthorizeResponse{
		Reference:        psResp.Data.Reference,
		AuthorizationURL: psResp.Data.AuthorizationURL,
	}, nil
}

func (p *PaystackClient) Verify(
	ctx context.Context,
	reference string,
) (payment.VerifyResponse, error) {

	url := fmt.Sprintf("%s/transaction/verify/%s", p.baseURL, reference)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return payment.VerifyResponse{}, fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.http.Do(httpReq)
	if err != nil {
		return payment.VerifyResponse{}, fmt.Errorf("paystack request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return payment.VerifyResponse{}, fmt.Errorf("paystack returned status %d, failed to decode error: %w", resp.StatusCode, err)
		}
		return payment.VerifyResponse{}, fmt.Errorf("paystack error: %s", errResp.Message)
	}

	var psResp struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Status    string `json:"status"` // "success", "failed", etc
			Reference string `json:"reference"`
			Amount    int64  `json:"amount"`
			Currency  string `json:"currency"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&psResp); err != nil {
		return payment.VerifyResponse{}, fmt.Errorf("decode paystack response: %w", err)
	}

	if !psResp.Status {
		return payment.VerifyResponse{}, fmt.Errorf("paystack error: %s", psResp.Message)
	}

	// Map Paystack data to internal domain
	return payment.VerifyResponse{
		Reference: psResp.Data.Reference,
		Status:    psResp.Data.Status,
		Amount:    psResp.Data.Amount,
		Currency:  psResp.Data.Currency,
	}, nil
}

func (p *PaystackClient) Refund(
	ctx context.Context,
	req payment.RefundRequest,
) (payment.RefundResponse, error) {

	payload := map[string]interface{}{
		"reference": req.Reference,
		"amount":    req.Amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return payment.RefundResponse{}, fmt.Errorf("marshal paystack refund request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseURL+"/refund",
		bytes.NewReader(body),
	)
	if err != nil {
		return payment.RefundResponse{}, fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.http.Do(httpReq)
	if err != nil {
		return payment.RefundResponse{}, fmt.Errorf("paystack refund request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return payment.RefundResponse{}, fmt.Errorf("paystack returned status %d, failed to decode error: %w", resp.StatusCode, err)
		}
		return payment.RefundResponse{}, fmt.Errorf("paystack error: %s", errResp.Message)
	}

	var psResp struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Reference string `json:"reference"`
			Status    string `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&psResp); err != nil {
		return payment.RefundResponse{}, fmt.Errorf("decode paystack refund response: %w", err)
	}

	if !psResp.Status {
		return payment.RefundResponse{}, fmt.Errorf("paystack refund error: %s", psResp.Message)
	}

	return payment.RefundResponse{
		Reference: psResp.Data.Reference,
		Status:    psResp.Data.Status,
	}, nil
}
