package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
GitHub returns a Handler that authenticates via GitHub's Authorization for
Webhooks scheme (https://developer.github.com/webhooks/securing/#validating-payloads-from-github)

Writes a http.StatusUnauthorized if authentication fails
*/
func GitHub(secret string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		requestSignature := req.Header.Get("X_HUB_SIGNATURE")

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusUnauthorized)
		}

		mac := hmac.New(sha1.New, []byte(secret))
		mac.Write([]byte(body))
		calculatedSignature := fmt.Sprintf("sha1=%x", mac.Sum(nil))

		if !SecureCompare(requestSignature, calculatedSignature) {
			fmt.Printf("request: %s, calculated: %s\n", requestSignature, calculatedSignature)
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}
}
