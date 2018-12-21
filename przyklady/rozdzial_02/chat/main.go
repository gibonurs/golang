package main

import (
  "flag"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "sync"
  "text/template"

  "github.com/matryer/goblueprints/chapter1/trace"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/gomniauth/providers/facebook"
  "github.com/stretchr/gomniauth/providers/github"
  "github.com/stretchr/gomniauth/providers/google"
  "github.com/stretchr/objx"
)

// Struktura reprezentująca pojedyczny szablon.
type templateHandler struct {
  once     sync.Once
  filename string
  templ    *template.Template
}

// ServeHTTP obsługuje żądania HTTP.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  })

  data := map[string]interface{}{
    "Host": r.Host,
  }
  if authCookie, err := r.Cookie("auth"); err == nil {
    data["UserData"] = objx.MustFromBase64(authCookie.Value)
  }

  t.templ.Execute(w, data)
}

var host = flag.String("host", ":8080", "Komputer na którym działa aplikacja.")

func main() {

  flag.Parse() // analiza flag wiersza poleceń

  // konfiguracja pakietu gomniauth
  gomniauth.SetSecurityKey("98dfbg7iu2nb4uywevihjw4tuiyub34noilk")
  gomniauth.WithProviders(
    github.New("3d1e6ba69036e0624b61", "7e8938928d802e7582908a5eadaaaf22d64babf1", "http://localhost:8080/auth/callback/github"),
    //google.New("44166123467-o6brs9o43tgaek9q12lef07bk48m3jmf.apps.googleusercontent.com", "rpXpakthfjPVoFGvcf9CVCu7", "http://localhost:8080/auth/callback/google"),
    google.New("708433410979-vrtjtcgie21mej902fq4e2bae0afdt50.apps.googleusercontent.com", "qfq0uZkQjcVogWlwucQDLUDc", "http://localhost:8080/auth/callback/google"),
    facebook.New("537611606322077", "f9f4d77b3d3f4f5775369f5c9f88f65e", "http://localhost:8080/auth/callback/facebook"),
  )

  r := newRoom()
  r.tracer = trace.New(os.Stdout)

  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.HandleFunc("/auth/", loginHandler)
  http.Handle("/room", r)

  // uruchomienie pokoju rozmów
  go r.run()

  // uruchomienie serwera WWW
  log.Println("Uruchamianie serwera WWW pod adresem", *host)
  if err := http.ListenAndServe(*host, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }

}
