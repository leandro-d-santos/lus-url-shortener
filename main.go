package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var tmpl *template.Template
var urls = make(map[string]string)

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/home", Homepage)
	mux.HandleFunc("/url", CreateUrlHandler)
	mux.HandleFunc("/", RedirectHandler)

	port := os.Getenv("SERVER_PORT")
	log.Printf("Servidor iniciado na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}

func CreateUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	url := r.FormValue("url")
	shortCode := generateShortCode()
	urls[shortCode] = url
	domainName := os.Getenv("DOMAINNAME")
	response := fmt.Sprintf("http://%s/%s", domainName, shortCode)
	data := struct {
		URL string
	}{
		URL: response,
	}
	tmpl.ExecuteTemplate(w, "result-url.html", data)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]
	originalURL, exists := urls[shortCode]
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[random.Intn(len(charset))]
	}
	return string(code)
}
