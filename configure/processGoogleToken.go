package configure

import ("encoding/json"
	. "github.com/mchirico/go_who/util"
	"log"
)

func GetGoogleToken(file string) (GoogleToken, error) {
	data, err := ReadFile(file)
	if err != nil {
		log.Printf("Cannot read Google Token file")
	}
	googleToken := GoogleToken{}
	err = json.Unmarshal([]byte(data), &googleToken)
	if err != nil {
		return googleToken,err
	}
	return googleToken,err
}
