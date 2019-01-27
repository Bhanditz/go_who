package configure

import (
	"encoding/json"
	"fmt"
	"github.com/mchirico/go_who/google"
	. "github.com/mchirico/go_who/util"
	"os/user"
	"testing"
)

/*
This is the file like:
client_secret_620409036887-4fn745fp4f2mrbi50oudcjk75r8vjg51.apps.googleusercontent.com.json

Which is downloaded from Google.
*/
func TestReadGoogleClientID(t *testing.T) {

	tokenFile := `client_secret_620409036887-4fn745fp4f2mrbi50oudcjk75r8vjg51.apps.googleusercontent.com.json`
	usr, _ := user.Current()
	file := usr.HomeDir + "/.go_who/" + tokenFile

	data, err := ReadFile(file)
	if err != nil {
		t.Fail()
	}
	googleToken := google.GoogleToken{}
	err = json.Unmarshal([]byte(data), &googleToken)
	if err != nil {
		panic(err)
	}

	if googleToken.Web.Auth_cert != "https://www.googleapis.com/oauth2/v1/certs" {
		t.Fail()
	}
	fmt.Printf("\n...%v", googleToken.Web.Auth_cert)
}
