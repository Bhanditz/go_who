package main

import (
	"github.com/mchirico/go_who/pkg"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a pkg.App

func TestMain(m *testing.M) {
	a = pkg.App{}

	a.Initilize()
	code := m.Run()

	os.Exit(code)
}

func TestEmptyProducts(t *testing.T) {

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestRoot(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body !=
		`<a href="https://accounts.google.com/o/oauth2/v2/auth?scope=https://www.googleapis.com/auth/userinfo.email%20https://www.googleapis.com/auth/userinfo.profile%20https://www.googleapis.com/auth/plus.login&access_type=offline&include_granted_scopes=true&state=state_parameter_passthrough_value&redirect_uri=https://who.aipiggybot.io/auth/google/callback&response_type=code&client_id=162501148484-ln73evqghc9slatlgjspevcujsgnvp3s.apps.googleusercontent.com">
		 Google Auth </a> <br><br>
		[{"page":1,"fruits":["pear","orange"]},{"page":2,"fruits":["pear","orange"]}]` {
		t.Errorf("Expected an array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
