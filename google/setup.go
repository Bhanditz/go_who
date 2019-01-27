package google

import (
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"

	//"github.com/levigross/grequests"
	. "github.com/mchirico/go_who/util"
	"log"
	"os/user"
	"sync"
)

type TokenDetails struct {
	Client_id     string   `json:"client_id"`
	Project_id    string   `json:"project_id"`
	Auth_uri      string   `json:"auth_uri"`
	Token_uri     string   `json:token_uri`
	Auth_cert     string   `json:"auth_provider_x509_cert_url"`
	Client_secret string   `json:"client_secret"`
	Redirect_uris []string `json:"redirect_uris"`
}

type GoogleToken struct {
	Web TokenDetails
}

type GoogleResponse struct {
	Access_token  string `json:"access_token"`
	Expires_in    int    `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
	Scope         string `json:"scope"`
	Token_type    string `json:"token_type"`
	Id_token      string `json:"id_token"`
}

type GoogleWebStruct struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthUri                 string   `json:"auth_uri"`
	TokenUri                string   `json:"token_uri"`
	AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
	JavascriptOrigins       []string `json:"javascript_origins"`
}

type GoogleSecret struct {
	Web  GoogleWebStruct `json:"web"`
	Code string
	sync.Mutex
}

func (g *GoogleSecret) GetGoogleSecret(f_optional ...string) {

	g.Lock()
	defer g.Unlock()

	usr, _ := user.Current()
	file := usr.HomeDir + "/.go_who_secret_google"

	if len(f_optional) == 2 {
		file = fmt.Sprintf("%s/%s", f_optional[0], f_optional[1])
	} else if len(f_optional) == 1 {
		file = usr.HomeDir + "/" + f_optional[0]
	}

	data, err := ReadFile(file)
	if err != nil {
		log.Printf("GetGoogleSecret error reading file: %v\n", file)
	}

	err = json.Unmarshal([]byte(data), &g)
	if err != nil {
		log.Printf("GetGoogleSecret Unmarshal: %v\n", err)
	}

}

/*
"code=4/xAAqNa&\
client_id=1625&\
client_secret=zlx-q&\
redirect_uri=https://who.aipiggybot.io/auth/google/callback&\
grant_type=authorization_code" https://www.googleapis.com/oauth2/v4/token

*/

type Gas struct {
	Code         string `json:"code"`
	ClienID      string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type""`
}

func (g *GoogleSecret) GetToken(code string) {
	g.Lock()
	defer g.Unlock()

	g.Code = code
	log.Printf("code: %v\n", code)
}

func (g *GoogleSecret) MakeRequest(
	url string, code string) (*grequests.Response, error) {

	g.Lock()
	defer g.Unlock()

	g.Code = code

	ro := grequests.RequestOptions{}

	log.Printf("**  Using this code: %v\n", g.Code)

	headers := map[string]string{}
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	ro.Headers = headers
	m := map[string]string{}
	m["client_id"] = g.Web.ClientID
	m["client_secret"] = g.Web.ClientSecret
	m["code"] = code
	m["grant_type"] = "authorization_code"
	m["redirect_uri"] = g.Web.RedirectUris[0]
	ro.Data = m
	r, err := grequests.Post(url, &ro)
	return r, err

}
