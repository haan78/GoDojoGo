package libbrevo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BREVO_EMAIL_SEND_URL    = "https://api.sendinblue.com/v3/smtp/email"
	BREOV_ACCOUNT_CHECK_URL = "https://api.sendinblue.com/v3/account"
)

type apiError struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
}

// SendinblueTemplateEmail sends a transactional email using a Brevo (Sendinblue) template.
// - email: recipient email
// - templateID: Brevo templateId
// - params: template params (can be nil). Use map[string]any or a struct.
func SendinblueTemplateEmail(apiKey string, email string, templateID int, params any) error {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Build request payload
	payload := map[string]any{
		"templateId": templateID,
		"to": []map[string]string{
			{"email": email},
		},
		"params": params, // can be nil -> JSON null
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, BREVO_EMAIL_SEND_URL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Non-2xx => error (include Brevo payload if possible)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var ae apiError
		if json.Unmarshal(respBody, &ae) == nil && (ae.Code != nil || ae.Message != "") {
			return fmt.Errorf("brevo error: %v / %s (http %d)", ae.Code, ae.Message, resp.StatusCode)
		}
		return fmt.Errorf("brevo http error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	// Mirror PHP behavior: if response JSON includes {code,message}, treat as error
	var ae apiError
	if json.Unmarshal(respBody, &ae) == nil && ae.Code != nil && fmt.Sprint(ae.Code) != "" {
		return fmt.Errorf("brevo error: %v / %s", ae.Code, ae.Message)
	}

	// Mirror PHP behavior: if JSON decode fails, treat as failure
	if !json.Valid(respBody) {
		return fmt.Errorf("email send failed, try again later (%s)", string(respBody))
	}

	return nil
}

type PlanVertical struct {
	PlanCategory string `json:"planCategory"`
	PlanType     string `json:"planType"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
}
type AccountResponse struct {
	PlanVerticals []PlanVertical `json:"planVerticals"`
}

func SendinblueCheck(apiKey string) ([]PlanVertical, error) {
	var result []PlanVertical
	if apiKey == "" {
		return result, fmt.Errorf("api key is empty")
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		BREOV_ACCOUNT_CHECK_URL,
		nil,
	)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("brevo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("brevo unhealthy, status code: %d", resp.StatusCode)
	}

	var account AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return result, fmt.Errorf("failed to decode brevo response: %w", err)
	}

	// Check if any planVertical is active
	for _, plan := range account.PlanVerticals {
		if plan.Status == "active" {
			result = append(result, plan)
		}
	}

	if len(result) > 0 {
		return result, nil
	} else {
		return account.PlanVerticals, fmt.Errorf("brevo account has no active plan")
	}

}
