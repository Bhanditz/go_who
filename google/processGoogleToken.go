package google

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mchirico/go_who/util"
	"log"
)

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

func GetGoogleUserFromToken(input []byte) (map[string]interface{}, error) {

	r := GoogleResponse{}
	err := json.Unmarshal(input, &r)
	if err != nil {
		return nil, err
	}

	jwtToken := r.Id_token

	claims := jwt.MapClaims{}
	//TODO: Figure out how to check for google's key
	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<public key>"), nil
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
		return m, nil
	}

	return m, errors.New("Email Not Verified")

}
