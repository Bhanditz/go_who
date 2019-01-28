package google

import (
	"fmt"
	"github.com/levigross/grequests"
	"log"
	"os/user"
	"testing"
)

func TestGetGoogleToken(t *testing.T) {

	tokenFile := `client_secret_620409036887-4fn745fp4f2mrbi50oudcjk75r8vjg51.apps.googleusercontent.com.json`
	usr, _ := user.Current()
	file := usr.HomeDir + "/.go_who/" + tokenFile

	googleToken, err := GetGoogleToken(file)

	if err != nil {
		t.Fail()
	}

	if googleToken.Web.Auth_cert != "https://www.googleapis.com/oauth2/v1/certs" {
		t.Fail()
	}
	fmt.Printf("\n...%v", googleToken.Web.Auth_cert)
}

/*
curl   -d  "code=4/3wAisfTdeQ&\
client_id=620409036887-4fn745fp4f2mrbi50oudcjk75r8vjg51.apps.googleusercontent.com&\
client_secret=CDnx&\
redirect_uri=https://who.aipiggybot.io/auth/google/callback&\
grant_type=authorization_code" https://httpbin.org/post

But to https://httpbin.org/post instead of
    https://www.googleapis.com/oauth2/v4/token

Below is a simple test:
*/

func TestPost(t *testing.T) {

	ro := grequests.RequestOptions{}

	headers := map[string]string{}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	ro.Headers = headers
	m := map[string]string{}
	m["client_id"] = "620409036887-4fn745fp4f2mrbi50oudcjk75r8vjg51.apps.googleusercontent.com"
	m["client_secret"] = "CDnx"
	m["code"] = "4/3wAisfTdeQ"
	m["grant_type"] = "authorization_code"
	m["redirect_uri"] = "https://who.aipiggybot.io/auth/google/callback"
	ro.Data = m
	r, _ := grequests.Post("https://httpbin.org/post", &ro)
	log.Println(r)

}

func TestMakeRequest(t *testing.T) {

	g := GoogleSecret{}
	g.GetGoogleSecret(".go_who_secret_dummy")

	r, e := g.MakeRequest("https://httpbin.org/post", "code...")
	log.Println(r, e)

}

func TestGetGoogleUserFromToken(t *testing.T) {
	input := []byte(`{
  "access_token": "ACCESS_TOKEN",
  "expires_in": 3600,
  "refresh_token": "RefreshToken",
  "scope": "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
  "token_type": "Bearer",
  "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImIxNWEyYjhmN2E2YjNmNmJjMDhiYzFjNTZhODg0MTBlMTQ2ZDAxZmQiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI2MjA0MDkwMzY4ODctNGZuNzQ1ZnA0ZjJtcmJpNTBvdWRjams3NXI4dmpnNTEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2MjA0MDkwMzY4ODctNGZuNzQ1ZnA0ZjJtcmJpNTBvdWRjams3NXI4dmpnNTEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDUwNjcyNjg3MzM5MDczMTM3MjgiLCJlbWFpbCI6Im1jaGlyaWNvQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdF9oYXNoIjoiTi0xaEJEVXJqOWFLT0VMWGVjY0E4USIsIm5hbWUiOiJNaWtlIENoaXJpY28iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDQuZ29vZ2xldXNlcmNvbnRlbnQuY29tLy1tZl9sVHZYb3M5ay9BQUFBQUFBQUFBSS9BQUFBQUFBQVBscy9Eal9rZEFBeGtzWS9zOTYtYy9waG90by5qcGciLCJnaXZlbl9uYW1lIjoiTWlrZSIsImZhbWlseV9uYW1lIjoiQ2hpcmljbyIsImxvY2FsZSI6ImVuIiwiaWF0IjoxNTQ4NTQ1NzA5LCJleHAiOjE1NDg1NDkzMDl9.TsV3bflCbQGj28aYfIPGIPCbt0Qro5Qya7mrprTiRFx4GwOEHe4W5WSZeNl8h3BaW9BNvS5MeEvho-PN9mtCR9Yi24rK3Qc81BdKXOMZMbUg72oLiAqUkWauSabBa4-9wvfZ56-LOF4aDdZDGXXcGmicmNCa7H0X3m2tDYZDY7Fn9iSsfd0DnVrRVb7sdfdlwAmzYqz2Zi7ct4i-44s7RY7U4XntVD2ZYrsaLuxTNFVLX4jsP6_LQ5RE6-yNM4adZT1l21keijz7-fJNPR5fF6kqsz14JuSaHDHaowvOVdXLl67yBss_vOA9U9vVqYzeNlHS5zVvFj_8KlItckwwQw"}`)

	r, err := GetGoogleUserFromToken(input)

	if err.Error() != "Token is expired" {
		t.Fail()
	}

	if r["email_verified"].(bool) != true {
		t.Fail()
	}
	if r["email"].(string) != "mchirico@gmail.com" {
		t.Fail()
	}
}

func TestGetToken(t *testing.T) {
	g := GoogleSecret{}
	g.GetGoogleSecret(".go_who_secret_dummy")
	log.Printf("%v\n", g.Web.ClientID)

}
