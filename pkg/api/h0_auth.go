package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func authorization(lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			buf bytes.Buffer
			js  map[string]any
		)

		if _, err := buf.ReadFrom(r.Body); err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}

		js = map[string]any{}

		if err := json.Unmarshal(buf.Bytes(), &js); err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}

		password, ok := js["password"].(string)
		if !ok {
			writeJson(w, fmt.Errorf("incorrect password type %T, expected 'string'", js["password"]), http.StatusBadRequest, lg)
			return
		}

		if password != os.Getenv("TODO_PASSWORD") {
			writeJson(w, fmt.Errorf("incorrect password"), http.StatusUnauthorized, lg)
			return
		}

		jwtToken, err := makeJWT(password)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}
		lg.Println("token created")

		fmt.Println(jwtToken)
		writeJson(w, map[string]string{"token": jwtToken}, http.StatusOK, lg)
	}
}

func auth(next http.HandlerFunc, lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pas := os.Getenv("TODO_PASSWORD")
		if len(pas) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				lg.Printf("cookie error: %v", err)
				http.Error(w, "cookie error "+err.Error(), http.StatusUnauthorized)
				return
			}

			if err = tokenValidation(cookie.Value); err != nil {
				lg.Printf("failed validation: %v", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			lg.Println("successful validation")
		}
		next(w, r)
	}
}

func makeJWT(password string) (string, error) {

	jsheader, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	encodedHeader := base64.URLEncoding.EncodeToString(jsheader)

	pasHash := sha256.Sum256([]byte(password))
	jspayload, err := json.Marshal(map[string]string{"pasHash": base64.URLEncoding.EncodeToString(pasHash[:])})
	if err != nil {
		return "", err
	}
	encodedPayload := base64.URLEncoding.EncodeToString(jspayload)

	headerPayload := encodedHeader + "." + encodedPayload

	sign := signHeaderPayload(headerPayload)

	encodedSign := base64.URLEncoding.EncodeToString(sign)

	return headerPayload + "." + encodedSign, nil
}

func tokenValidation(token string) error {

	baseToken := strings.Split(token, ".")
	if len(baseToken) != 3 {
		return fmt.Errorf("token length not equal 3")
	}

	jsHeader, err := base64.URLEncoding.DecodeString(baseToken[0])
	if err != nil {
		return fmt.Errorf("failed to decode header: %w", err)
	}
	header := map[string]any{}
	if err = json.Unmarshal(jsHeader, &header); err != nil {
		return err
	}
	alg, aOk := header["alg"].(string)
	typ, tOk := header["typ"].(string)
	if !aOk || alg != "HS256" {
		return fmt.Errorf("invalid or missing alg: %s", alg)
	}
	if !tOk || typ != "JWT" {
		return fmt.Errorf("invalid or missing typ: %s", typ)
	}

	jsPayload, err := base64.URLEncoding.DecodeString(baseToken[1])
	if err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}
	payload := map[string]any{}
	if err = json.Unmarshal(jsPayload, &payload); err != nil {
		return err
	}
	receivedPasHash, ok := payload["pasHash"].(string)
	if !ok {
		return fmt.Errorf("invalid type of value in payload")
	}
	ph := sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD")))
	expectedPasHash := base64.URLEncoding.EncodeToString(ph[:])
	if receivedPasHash != expectedPasHash {
		return fmt.Errorf("password changed or invalid")
	}

	expectedSign := signHeaderPayload(baseToken[0] + "." + baseToken[1])
	receivedSign, err := base64.URLEncoding.DecodeString(baseToken[2])
	if err != nil {
		return fmt.Errorf("error decoding received signature: %w", err)
	}
	if !hmac.Equal(expectedSign, receivedSign) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}

func signHeaderPayload(headerPayload string) []byte {

	h := hmac.New(sha256.New, []byte(os.Getenv("SECRET")))
	h.Write([]byte(headerPayload))
	signature := h.Sum(nil)
	return signature
}
