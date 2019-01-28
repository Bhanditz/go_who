package google

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mchirico/go_who/util"
	"log"
	"net/http"
	"strings"
)

const (
	clientCertURL = "https://www.googleapis.com/oauth2/v1/certs"
)

func convertKey(key string) interface{} {
	certPEM := key
	certPEM = strings.Replace(certPEM, "\\n", "\n", -1)
	certPEM = strings.Replace(certPEM, "\"", "", -1)
	block, _ := pem.Decode([]byte(certPEM))
	cert, _ := x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

	return rsaPublicKey
}

func fetchPublicKeys() (map[string]*json.RawMessage, error) {
	resp, err := http.Get(clientCertURL)
	if err != nil {
		return nil, err
	}

	var objmap map[string]*json.RawMessage
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&objmap)

	return objmap, err
}

func GetGoogleToken(file string) (GoogleToken, error) {
	data, err := util.ReadFile(file)
	if err != nil {
		log.Printf("Cannot read Google Token file")
	}
	googleToken := GoogleToken{}
	err = json.Unmarshal([]byte(data), &googleToken)
	if err != nil {
		return googleToken, err
	}
	return googleToken, err
}

// TODO: This needs to be cleaned up
func GetGoogleUserFromToken(input []byte) (map[string]interface{}, error) {

	keys, err := fetchPublicKeys()
	if err != nil {
		return nil, nil
	}

	r := GoogleResponse{}
	err = json.Unmarshal(input, &r)
	if err != nil {
		return nil, err
	}

	jwtToken := r.Id_token
	claims := jwt.MapClaims{}

	_, err2 := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"]
		rsaPublicKey := convertKey(string(*keys[kid.(string)]))
		return rsaPublicKey, nil
	})

	m := map[string]interface{}{}
	for key, val := range claims {
		//fmt.Printf("Key: %v, value: %v\n", key, val)
		m[key] = val
	}

	verified, ok := m["email_verified"]
	if ok {
		if verified.(bool) != true {
			return m, errors.New("Email Not Verified")
		}
		return m, err2
	}

	return m, errors.New("Email Not Verified")

}
