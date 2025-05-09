package google

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// API Request

// URL: https://www.google.com/recaptcha/api/siteverify

// METHOD: POST
// POST Parameter 	Description
// secret 	Required. The shared key between your site and reCAPTCHA.
// response 	Required. The user response token provided by the reCAPTCHA client-side integration on your site.
// remoteip 	Optional. The user's IP address.

// API Response

// The response is a JSON object:

// {
//   "success": true|false,
//   "challenge_ts": timestamp,  // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
//   "hostname": string,         // the hostname of the site where the reCAPTCHA was solved
//   "error-codes": [...]        // optional
// }

type (
	ReCaptchaAPI struct {
		secret string `env:"RECAPTCHA_SECRET"`
	}
)

// NewReCaptchaAPI creates a new ReCapchaAPI.
func NewReCaptchaAPI(secret string) *ReCaptchaAPI {
	return &ReCaptchaAPI{
		secret: secret,
	}
}

// Verify verifies the reCAPTCHA response.
func (r *ReCaptchaAPI) Verify(response, remoteip string) (bool, error) {
	// Create a new HTTP request.
	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}

	// Set the form parameters.
	q := req.URL.Query()
	q.Add("secret", r.secret)
	q.Add("response", response)
	if remoteip != "" && !strings.HasPrefix(remoteip, "127.") && !strings.HasPrefix(remoteip, "192.") {
		q.Add("remoteip", remoteip)
	}
	req.URL.RawQuery = q.Encode()

	// Send the HTTP request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response.
	var result struct {
		Success     bool     `json:"success"`
		ChallengeTS string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		ErrorCodes  []string `json:"error-codes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("error decoding response: %w", err)
	}

	// Check if the response was successful.
	if !result.Success {
		return false, fmt.Errorf("reCAPTCHA verification failed: %v", result.ErrorCodes)
	}

	return true, nil
}
