package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s - %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

var (
	flagServiceAccountName  string
	flagServiceAccountToken string
)

func main() {
	flag.StringVar(&flagServiceAccountName, "account", "local", "name of service account to advertise")
	flag.StringVar(&flagServiceAccountToken, "token", "", "service account bearer token")
	flag.Parse()

	if flagServiceAccountToken == "" {
		log.Fatal("service account token required")
	}

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/computeMetadata/v1/project/numeric-project-id",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "0123456789")
		})

	r.HandleFunc("/computeMetadata/v1/project/project-id",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "1234567890")
		})

	r.HandleFunc("/computeMetadata/v1/instance/service-accounts/{name}/email",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s@google.internal", flagServiceAccountName)
		})

	r.HandleFunc("/computeMetadata/v1/instance/service-accounts/{name}/scopes",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "email")
		})

	r.HandleFunc("/computeMetadata/v1/instance/service-accounts/{name}",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, `{"scopes":"email","email":"%s","aliases":["default"]}`, mux.Vars(r)["name"])
		})

	r.HandleFunc("/computeMetadata/v1/instance/service-accounts/{name}/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, `{"access_token":"%s","expires_in":3000,"token_type":"Bearer"}`, flagServiceAccountToken)
		})

	r.HandleFunc("/computeMetadata/v1/instance/service-accounts",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "default/\n%s@google.internal", flagServiceAccountName)
		})

	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("metadata-flavor", "Google")
			fmt.Fprintf(w, "computeMetadata/")
		})

	log.Fatal(http.Serve(ln, logRequest(r)))
}
