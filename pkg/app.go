package pkg

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/levigross/grequests"
	"github.com/mchirico/go_who/configure"
	"github.com/mchirico/go_who/google"
	"github.com/mchirico/go_who/rand"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var store = sessions.NewCookieStore([]byte(rand.RandomString(73)))

type App struct {
	Router *mux.Router
	DB     *sql.DB
	at     *int
	secStr *configure.SecretStruct
}

func (a *App) Initilize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// TODO: Take our harvest and pump in flag for google value
func (a *App) InitSS(secStr *configure.SecretStruct) {
	a.secStr = secStr
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.getRoot).Methods("GET")
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/auth/google/callback", a.getAuthGoogle).Methods("GET")
	a.Router.HandleFunc("/auth/google/callback", a.getAuthGoogle).Methods("POST")

	a.Router.HandleFunc("/auth2", a.getAuth2).Methods("GET")
	a.Router.HandleFunc("/auth2", a.getAuth2).Methods("POST")

	a.Router.HandleFunc("/upload", a.receiveFile).Methods("POST")

	a.Router.HandleFunc("/status", a.status).Methods("GET")
	// a.Router.HandleFunc("/product", a.createProduct).Methods("POST")

}

func (a *App) Run(addr string, writeTimeout int, readTimeout int) {

	srv := &http.Server{
		Handler: a.Router,
		Addr:    fmt.Sprintf(":%s", addr),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func (a *App) getRoot(w http.ResponseWriter, r *http.Request) {

	log.Printf("get Root")
	products, err := getRoot(a.DB, 0, 5)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) receiveFile(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Printf("File name %s\n", header.Filename)
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	contents := Buf.String()
	fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header
	return
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {

	products, err := getProducts(a.DB, 0, 5)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) status(w http.ResponseWriter, r *http.Request) {

	log.Printf("We are in status\n")
	log.Printf("Method: %v\n", r.Method)
	log.Printf("Header: %v\n", r.Header)

	session, err := store.Get(r, "session-user")

	log.Printf("session: %v\n", session.Values["email"])

	if err != nil {
		log.Printf("Session error: %v\n", err)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	s := fmt.Sprintf("{email:\"%s\"}", session.Values["email"])
	_, err = w.Write([]byte(s))
	if err != nil {
		log.Printf("Can not write response: %v\n", err)
	}

}

// Using this one
func (a *App) getAuthGoogle(w http.ResponseWriter, r *http.Request) {

	log.Printf("We are in getAuthGoogle...\n")
	log.Printf("Method: %v\n", r.Method)
	log.Printf("Header: %v\n", r.Header)

	session, err := store.Get(r, "session-user")
	if err != nil {
		log.Printf("This is good... not previous cookie: %v\n", err)

	} else {
		// Don't come back to this point..
		return
	}

	vals := r.URL.Query()
	log.Printf("vals: %v\n", vals)

	code, ok := vals["code"]
	if ok {

		log.Printf("\ncode: %v\n", code)
		log.Printf("\ngoogle.GoogleSecret{}\n")

		g := google.GoogleSecret{}
		g.GetGoogleSecret(".go_who_secret")

		response, e := g.MakeRequest("https://www.googleapis.com/oauth2/v4/token", code[0])
		log.Println(response, e)
		result, err := google.GetGoogleUserFromToken([]byte(response.String()))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, err = w.Write([]byte(`{user:"invalid"}`))
			if err != nil {
				log.Printf("Can not write response: %v\n", err)
			}
			return
		}

		log.Printf("email...\n%v\n", result["email"])

		// Set some session values.
		session.Values["email"] = result["email"]
		session.Values[42] = 43

		session.Save(r, w)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write([]byte(`{a:"one"}`))
		if err != nil {
			log.Printf("Can not write response: %v\n", err)
		}

	}

}

func (a *App) getAuth(w http.ResponseWriter, r *http.Request) {

	log.Printf("We are in getAuth ...\n")
	log.Printf("Method: %v\n", r.Method)
	log.Printf("Header: %v\n", r.Header)

	vals := r.URL.Query()
	code, ok := vals["code"]
	if ok {
		log.Printf("code: %v\n", code[0])

		c := configure.CodeToPassStruct{}
		c.Secret = a.secStr.Secret
		c.Id = a.secStr.Id
		c.Code = code[0]
		c.GrantType = "authorization_code"

		data, err := c.Marshel()
		if err != nil {
			log.Fatalf("c.Marshel %v\n", err)
		}

		log.Printf("CodeToPassStruct: %v\n", string(data))

		ro := grequests.RequestOptions{}

		headers := map[string]string{}
		headers["Content-Type"] = "application/json"
		headers["User-Agent"] = "AiPiggybot (mchirico@gmail.com)"
		ro.Headers = headers
		ro.JSON = data

		defer func() {
			ResponseCode("code", ro, *a.secStr)
			if r := recover(); r != nil {
				log.Println("Recovered in f", r)
				fmt.Fprint(w, "Recovered")
			}
		}()

		fmt.Fprint(w, "code")
		fmt.Printf("data: %v\n", string(data))

	} else {
		fmt.Printf("\n\n Something Wrong:%v\n\n", vals)
	}
}

func (a *App) getAuth2(w http.ResponseWriter, r *http.Request) {

	log.Printf("We are in getAuth2\n")
	log.Printf("Method: %v\n", r.Method)
	log.Printf("Header: %v\n", r.Header)
	vals := r.URL.Query()
	log.Printf("vals:\n%v", vals)
	getToDoistToken(vals)

}

func getToDoistToken(vals map[string][]string) {
	ro := grequests.RequestOptions{}
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	//headers["Authorization"] = fmt.Sprintf("Bearer %s", string(access_token))
	headers["User-Agent"] = "AiPiggybot (mchirico@gmail.com)"
	ro.Headers = headers
	//url := "https://id.getharvest.com/api/v2/accounts"
}

/*

curl -X POST \
  -H "Content-Type: application/json" \
  -H "User-Agent: MyApp (yourname@example.com)" \
  -d "code=$AUTHORIZATION_CODE" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" \
  -d "grant_type=authorization_code" \
  'https://id.getharvest.com/api/v2/oauth2/token'

*/

func ResponseCode(code string, ro grequests.RequestOptions,
	secStr configure.SecretStruct) *grequests.Response {

	log.Printf("\n\nResponseCode...\n")
	response, err := grequests.Post(secStr.Url, &ro)
	if err != nil {
		log.Printf("err (ResponseCode): \n%v\n\n\n", err)
		return nil
	}

	if strings.Contains(response.String(), "error") {
		type E struct {
			Error     string `json:"error"`
			ErrorDesc string `json:"error_description"`
		}
		e := E{}
		response.JSON(&e)
		if e.Error != "" {
			log.Printf("We have a problem: %v", e.ErrorDesc)
		}
		return nil
	}

	log.Printf("REVIEW OUT:\n%v\n\n", response)
	configure.WriteResponseDataToFile(response, "https://id.getharvest.com/api/v2/accounts")

	return response

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Printf("Can not write response: %v\n", response)
	}
}
