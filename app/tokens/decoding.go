package tokens

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type JWTToken struct {
	Header    Header    `json:"header"`
	Payload   Payload   `json:"payload"`
	Signature Signature `json:"signature"`
}

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload map[string]interface{}

type Signature string

func tokenDecode(token string) (JWTToken, error) {
	var parts = strings.Split(token, ".")
	if len(parts) != 3 {
		return JWTToken{}, errors.New("invalid token")
	}
	header, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return JWTToken{}, err
	}
	payload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return JWTToken{}, err
	}
	signature, err := base64.URLEncoding.DecodeString(parts[2])
	if err != nil {
		return JWTToken{}, err
	}
	var t = JWTToken{
		Header:    Header{},
		Payload:   Payload{},
		Signature: Signature(signature),
	}
	err = json.Unmarshal(header, &t.Header)
	if err != nil {
		return JWTToken{}, err
	}
	err = json.Unmarshal(payload, &t.Payload)
	if err != nil {
		return JWTToken{}, err
	}
	return t, nil
}
