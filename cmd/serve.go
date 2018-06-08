package cmd

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"github.com/gobuffalo/pop"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/solarwinds/gitlic-check/augit/handlers"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/saml/samlsp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/url"
	"os"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func samlStrategy() *samlsp.Middleware {
	keyPair, err := tls.LoadX509KeyPair(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"))
	if err != nil {
		log.Fatalf("Failed to create X509 key pair. Error: %v\n", err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		log.Fatalf("Failed to parse X509 cert. Error: %v\n", err)
	}

	mdURL, err := url.Parse(os.Getenv("IDP_METADATA_URL"))
	if err != nil {
		log.Printf("Failed to parse metadata URL. Error: %v\n", err)
	}

	rootURL, err := url.Parse(os.Getenv("ROOT_URL"))
	if err != nil {
		log.Fatalf("Failed to parse root URL. Error: %v\n", err)
	}

	samlOpt := samlsp.Options{
		URL:            *rootURL,
		Key:            keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:    keyPair.Leaf,
		IDPMetadataURL: mdURL,
	}
	if os.Getenv("ENVIRONMENT") != "development" {
		samlOpt.CookieSecure = true
	}
	saml, err := samlsp.New(samlOpt)
	if err != nil {
		log.Fatalf("Failed to create SAML middleware: %v\n", err)
	}
	cookieName := "sw_token_" + os.Getenv("COOKIENAME")
	if cookie, ok := saml.ClientToken.(*samlsp.ClientCookies); ok {
		cookie.Name = cookieName
	}
	return saml
}

func augitHandlers(tx *pop.Connection) *mux.Router {
	sp := samlStrategy()
	r := mux.NewRouter()
	augit := r.PathPrefix("/augit").Subrouter()
	ghudb := models.NewGithubUserDB(tx)

	r.Handle("/", http.HandlerFunc(healthCheck())).Methods("GET")
	augit.Handle("/saml/acs", sp)
	augit.Handle("/user", sp.RequireAccount(http.HandlerFunc(handlers.ShowUser(ghudb)))).Methods("GET")
	augit.Handle("/user", sp.RequireAccount(http.HandlerFunc(handlers.AddUser(ghudb)))).Methods("POST")
	return r
}

func healthCheck() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
}

func runServer() {
	cxn, err := pop.Connect(os.Getenv("ENVIRONMENT"))
	if err != nil {
		log.Fatal(err)
	}
	migrator, err := pop.NewFileMigrator("./migrations", cxn)
	err = migrator.Up()
	if err != nil {
		log.Panic(err)
	}
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), gorillaHandlers.LoggingHandler(os.Stdout, augitHandlers(cxn))))
}
