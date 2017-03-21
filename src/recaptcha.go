package gocontact

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// RecaptchaResponse stores API response from the reCAPTCHA call
// https://developers.google.com/recaptcha/docs/verify
type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

const recaptchaSite = "https://www.google.com/recaptcha/api/siteverify"

var recaptchaPrivateKey string

// InitRecaptcha should be called in main to load reCAPTCHA credentials
func InitRecaptcha(privateKey string) {
	recaptchaPrivateKey = privateKey
}

// VerifyRecaptcha checks the reCAPTCHA validity given
// remoteip - IP source of original request, NB to check for reverse proxy IPs
// response - g-recaptcha-response value in the submitted form
func VerifyRecaptcha(remoteip, response string) bool {
	r, _ := check(remoteip, response)
	return r.Success
}

func check(remoteip, response string) (r RecaptchaResponse, err error) {
	client := &http.Client{Timeout: 20 * time.Second}
	data := url.Values{}
	data.Set("secret", recaptchaPrivateKey)
	data.Set("remoteip", remoteip)
	data.Set("response", response)
	resp, err := client.PostForm(recaptchaSite, data)
	if err != nil {
		log.Printf("post error: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read error: %s\n", err)
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("invalid JSON: %s\n", err)
	}
	return
}
