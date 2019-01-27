package configure

import (
	"fmt"
	"os/user"
	"testing"
)


func TestGetGoogleToken(t *testing.T) {

	tokenFile := `client_secret_620409036887-4fn745fp4f2mrbi50oudcjk75r8vjg51.apps.googleusercontent.com.json`
	usr, _ := user.Current()
	file := usr.HomeDir + "/.go_who/" + tokenFile

	googleToken,err := GetGoogleToken(file)

	if err != nil {
		t.Fail()
	}

	if googleToken.Web.Auth_cert != "https://www.googleapis.com/oauth2/v1/certs" {
		t.Fail()
	}
	fmt.Printf("\n...%v", googleToken.Web.Auth_cert)
}






