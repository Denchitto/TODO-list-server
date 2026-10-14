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
			buf   bytes.Buffer
			jsPas map[string]any
		)

		if _, err := buf.ReadFrom(r.Body); err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}

		jsPas = map[string]any{}

		if err := json.Unmarshal(buf.Bytes(), &jsPas); err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}

		password, ok := jsPas["password"].(string)
		if !ok {
			writeJson(w, fmt.Errorf("incorrect password type %T, expected 'string'", jsPas["password"]), http.StatusBadRequest, lg)
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
	// 1. Заголовок
	jsheader, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	encodedHeader := base64.URLEncoding.EncodeToString(jsheader)
	// 2. Полезная нагрузка
	pasHash := sha256.Sum256([]byte(password))
	jspayload, err := json.Marshal(map[string]string{"pasHash": base64.URLEncoding.EncodeToString(pasHash[:])})
	if err != nil {
		return "", err
	}
	encodedPayload := base64.URLEncoding.EncodeToString(jspayload)

	headerPayload := encodedHeader + "." + encodedPayload
	// 3. Подпись
	sign := signHeaderPayload(headerPayload)

	encodedSign := base64.URLEncoding.EncodeToString(sign)
	// 4. Готовый токен
	return headerPayload + "." + encodedSign, nil
}

func tokenValidation(token string) error {
	// 1. Разбиваем токен и проверяем его...
	baseToken := strings.Split(token, ".")
	if len(baseToken) != 3 {
		return fmt.Errorf("token length not equal 3")
	}
	// 2. ...заголовок...
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
	// 3. ...полезную нагрузку...
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
	// 4. ...подпись.
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
