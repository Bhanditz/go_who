package google

import (
	"encoding/json"
	"fmt"
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

var id_token = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImIxNWEyYjhmN2E2YjNmNmJjMDhiYzFjNTZhODg0MTBlMTQ2ZDAxZmQiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI2MjA0MDkwMzY4ODctNGZuNzQ1ZnA0ZjJtcmJpNTBvdWRjams3NXI4dmpnNTEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2MjA0MDkwMzY4ODctNGZuNzQ1ZnA0ZjJtcmJpNTBvdWRjams3NXI4dmpnNTEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDUwNjcyNjg3MzM5MDczMTM3MjgiLCJlbWFpbCI6Im1jaGlyaWNvQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdF9oYXNoIjoiTi0xaEJEVXJqOWFLT0VMWGVjY0E4USIsIm5hbWUiOiJNaWtlIENoaXJpY28iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDQuZ29vZ2xldXNlcmNvbnRlbnQuY29tLy1tZl9sVHZYb3M5ay9BQUFBQUFBQUFBSS9BQUFBQUFBQVBscy9Eal9rZEFBeGtzWS9zOTYtYy9waG90by5qcGciLCJnaXZlbl9uYW1lIjoiTWlrZSIsImZhbWlseV9uYW1lIjoiQ2hpcmljbyIsImxvY2FsZSI6ImVuIiwiaWF0IjoxNTQ4NTQ1NzA5LCJleHAiOjE1NDg1NDkzMDl9.TsV3bflCbQGj28aYfIPGIPCbt0Qro5Qya7mrprTiRFx4GwOEHe4W5WSZeNl8h3BaW9BNvS5MeEvho-PN9mtCR9Yi24rK3Qc81BdKXOMZMbUg72oLiAqUkWauSabBa4-9wvfZ56-LOF4aDdZDGXXcGmicmNCa7H0X3m2tDYZDY7Fn9iSsfd0DnVrRVb7sdfdlwAmzYqz2Zi7ct4i-44s7RY7U4XntVD2ZYrsaLuxTNFVLX4jsP6_LQ5RE6-yNM4adZT1l21keijz7-fJNPR5fF6kqsz14JuSaHDHaowvOVdXLl67yBss_vOA9U9vVqYzeNlHS5zVvFj_8KlItckwwQw"

//https://play.golang.org/p/2rzuPbU-6Rf
func TestResponseStruct(t *testing.T) {
	input := []byte(`{
  "access_token": "ACCESS_TOKEN",
  "expires_in": 3600,
  "refresh_token": "RefreshToken",
  "scope": "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
  "token_type": "Bearer",
  "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImIxNWEyYjhmN2E2YjNmNmJjMDhiYzFjNTZhODg0MTBlMTQ2ZDAxZmQiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI2MjA0MDkwMzY4ODctNGZuNzQ1ZnA0ZjJtcmJpNTBvdWRjams3NXI4dmpnNTEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2MjA0MDkwMzY4ODctNGZuNzQ1ZnA0ZjJtcmJpNTBvdWRjams3NXI4dmpnNTEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDUwNjcyNjg3MzM5MDczMTM3MjgiLCJlbWFpbCI6Im1jaGlyaWNvQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdF9oYXNoIjoiTi0xaEJEVXJqOWFLT0VMWGVjY0E4USIsIm5hbWUiOiJNaWtlIENoaXJpY28iLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDQuZ29vZ2xldXNlcmNvbnRlbnQuY29tLy1tZl9sVHZYb3M5ay9BQUFBQUFBQUFBSS9BQUFBQUFBQVBscy9Eal9rZEFBeGtzWS9zOTYtYy9waG90by5qcGciLCJnaXZlbl9uYW1lIjoiTWlrZSIsImZhbWlseV9uYW1lIjoiQ2hpcmljbyIsImxvY2FsZSI6ImVuIiwiaWF0IjoxNTQ4NTQ1NzA5LCJleHAiOjE1NDg1NDkzMDl9.TsV3bflCbQGj28aYfIPGIPCbt0Qro5Qya7mrprTiRFx4GwOEHe4W5WSZeNl8h3BaW9BNvS5MeEvho-PN9mtCR9Yi24rK3Qc81BdKXOMZMbUg72oLiAqUkWauSabBa4-9wvfZ56-LOF4aDdZDGXXcGmicmNCa7H0X3m2tDYZDY7Fn9iSsfd0DnVrRVb7sdfdlwAmzYqz2Zi7ct4i-44s7RY7U4XntVD2ZYrsaLuxTNFVLX4jsP6_LQ5RE6-yNM4adZT1l21keijz7-fJNPR5fF6kqsz14JuSaHDHaowvOVdXLl67yBss_vOA9U9vVqYzeNlHS5zVvFj_8KlItckwwQw"}`)

	r := GoogleResponse{}
	err := json.Unmarshal(input, &r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", r)
	if r.Access_token != "ACCESS_TOKEN" &&
		r.Expires_in != 3600 &&
		r.Refresh_token != "RefreshToken" &&
		r.Token_type != "Bearer" &&
		r.Id_token != id_token {
		t.Fail()
	}

}
