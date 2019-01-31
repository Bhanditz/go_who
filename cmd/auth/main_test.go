package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/levigross/grequests"
	"github.com/mchirico/go_who/configure"
	"github.com/mchirico/go_who/pkg"
	"github.com/mchirico/go_who/rand"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"testing"
)

var a pkg.App

func TestMain(m *testing.M) {
	a = pkg.App{}

	a.Initilize()
	code := m.Run()

	os.Exit(code)
}

// Ref: https://play.golang.org/p/UGeNKd-cw34
func TestResponseCode(t *testing.T) {
	type SendData struct {
		Num  float32  `json:"num"`
		Strs []string `json:"strs"`
	}
	ro := grequests.RequestOptions{}
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	ro.Headers = headers
	s := SendData{Num: 6.23, Strs: []string{"one", "two"}}
	ro.JSON = s
	url := "http://httpbin.org/post"
	result, _ := grequests.Post(url, &ro)

	log.Printf("\n\n result:%v\n", result)

	var f interface{}
	b := result.String()
	err := result.JSON(&f)
	if err != nil {
		t.Fail()
	}

	m := f.(map[string]interface{})
	log.Printf("json: %s\n", b)

	log.Printf("m: %v\n type: %T\n\n", m["json"], m["json"])

	for key, value := range m["json"].(map[string]interface{}) {
		fmt.Println("Key:", key, "Value:", value)
	}

	r := m["json"].(map[string]interface{})

	fmt.Println("num", r["num"].(float64))
	fmt.Println("strs[0]", r["strs"].([]interface{})[0])
	fmt.Println("strs[1]", r["strs"].([]interface{})[1])

	if r["strs"].([]interface{})[0] != "one" {
		t.Fail()
	}
	if r["strs"].([]interface{})[1] != "two" {
		t.Fail()
	}
}

func TestGettingSecret(t *testing.T) {

	usr, _ := user.Current()

	file := usr.HomeDir + "/.secretHarvest"

	oSecretStruct := configure.SecretStruct{}

	oSecretStruct.Id = "01223"
	oSecretStruct.Secret = "password"
	oSecretStruct.Url = "http://httpbin.org/post"

	odata, err := json.Marshal(oSecretStruct)

	n, err := writeFile(string(odata),
		file)
	if err != nil {
		log.Printf("error: %v, %v\n", n, err)
		t.Fail()
	}

	data, err := readFile(file)
	if err != nil {
		t.Fail()
	}

	res := configure.SecretStruct{}
	err = json.Unmarshal([]byte(odata), &res)
	if err != nil {
		t.Fail()
	}
	fmt.Println("res.Id: ", res.Id)
	fmt.Println("res.Secret: ", res.Secret)
	fmt.Println("res.Url: ", res.Url)
	fmt.Println("data: ", data)

}

func TestSecret(t *testing.T) {
	type SendData struct {
		Num  float32  `json:"num"`
		Strs []string `json:"strs"`
	}
	ro := grequests.RequestOptions{}
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	ro.Headers = headers
	s := SendData{Num: 6.23, Strs: []string{"one", "two"}}
	ro.JSON = s
	url := "http://httpbin.org/post"
	result, _ := grequests.Post(url, &ro)

	log.Printf("\n\n result:%v\n", result)

	var f interface{}
	b := result.String()
	result.JSON(&f)

	m := f.(map[string]interface{})
	log.Printf("json: %s\n", b)

	log.Printf("m: %v\n type: %T\n\n", m["json"], m["json"])

	for key, value := range m["json"].(map[string]interface{}) {
		fmt.Println("Key:", key, "Value:", value)
	}

	r := m["json"].(map[string]interface{})

	fmt.Println("num", r["num"].(float64))
	fmt.Println("strs[0]", r["strs"].([]interface{})[0])
	fmt.Println("strs[1]", r["strs"].([]interface{})[1])

	if r["strs"].([]interface{})[0] != "one" {
		t.Fail()
	}
	if r["strs"].([]interface{})[1] != "two" {
		t.Fail()
	}

}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return string(data), err
}

func writeFile(data string, file string) (int, error) {
	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		return -1, err
	}

	n, err := f.WriteString(data)

	return n, err
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
		`[{"page":1,"fruits":["pear","orange"]},{"page":2,"fruits":["pear","orange"]}]` {
		t.Errorf("Expected an array. Got %s", body)
	}
}

func NewRecorder() *httptest.ResponseRecorder {
	return &httptest.ResponseRecorder{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
	}
}

// Not working...
func testMetricWithUser(t *testing.T) {
	originalPath := "info:/"
	store := sessions.NewCookieStore(rand.RandKey)
	store.Options.Path = originalPath

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := NewRecorder()
	session, err := store.Get(req, "session-user")
	if err != nil {
		t.Fatal("failed to create session", err)
	}
	session.Values["email"] = "mchirico@gmail.com"
	session.Values[42] = 43
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save session", err)
	}
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body !=
		`{email:"mchirico@gmail.com"}` {
		t.Errorf("Expected an array. Got %s", body)
	}
}

func TestStatus(t *testing.T) {
	originalPath := "info:/"
	store := sessions.NewCookieStore(rand.RandKey)
	store.Options.Path = originalPath

	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := NewRecorder()
	session, err := store.Get(req, "session-user")
	if err != nil {
		t.Fatal("failed to create session", err)
	}
	session.Values["email"] = "mchirico@gmail.com"
	session.Values[42] = 43
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save session", err)
	}
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body !=
		`{email:"mchirico@gmail.com"}` {
		t.Errorf("Expected an email. Got %s", body)
	}

}

func TestInfo(t *testing.T) {
	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	expectedResponse := `{pig.za261.io:"https://bit.ly/2sMoemb",who.aipiggybot.io:"https://bit.ly/2B5hHrw"}`

	if body := response.Body.String(); body !=
		expectedResponse {
		t.Errorf("Expected:\n%s\n Got:\n%s\n", expectedResponse, body)
	}
}

func TestStatusNoEmail(t *testing.T) {
	originalPath := "info:/"
	store := sessions.NewCookieStore(rand.RandKey)
	store.Options.Path = originalPath

	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := NewRecorder()
	session, err := store.Get(req, "session-user")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	session.Values[42] = 43
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save session", err)
	}
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body !=
		`{email:"%!s(<nil>)"}` {
		t.Errorf("Expected an email. Got %s", body)
	}

}

func TestSession(t *testing.T) {
	originalPath := "/"
	store := sessions.NewFilesystemStore("")
	store.Options.Path = originalPath
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	store.Options.Path = "/foo"
	if session.Options.Path != originalPath {
		t.Fatalf("bad session path: got %q, want %q", session.Options.Path, originalPath)
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
