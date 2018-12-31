package google

import (
	"encoding/json"
	"github.com/levigross/grequests"
	. "github.com/mchirico/go_who/util"
	"log"
	"os/user"
	"testing"
)

func TestReadGoogle(t *testing.T) {
	usr, _ := user.Current()
	file := usr.HomeDir + "/.go_who_secret_google"

	data, err := ReadFile(file)
	if err != nil {
		t.Fail()
	}
	log.Printf("data: %v\n", data)
	g := GoogleSecret{}
	err = json.Unmarshal([]byte(data), &g)
	if err != nil {
		t.Fail()
	}
	log.Printf("here: %v\n", g.Web.ClientID)
}

func TestGetGoogleTokenRaw(t *testing.T) {

	g := GoogleSecret{}
	g.GetGoogleSecret(".go_who_secret_dummy")
	log.Println(g.Web.RedirectUris)

	ro := grequests.RequestOptions{}
	m := map[string]string{}
	m["client_id"] = g.Web.ClientID
	m["client_secret"] = g.Web.ClientSecret
	m["code"] = "code"
	m["grant_type"] = "authorization_code"
	m["redirect_uri"] = "https://who.aipiggybot.io/auth/google/callback"
	ro.Data = m
	//r, _ := grequests.Post("https://httpbin.org/post", &ro)
	r, _ := grequests.Post(g.Web.TokenUri, &ro)
	log.Println(r)

}
