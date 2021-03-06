package main

import (
  "context"
  "flag"
  "gopkg.in/mgo.v2"
  "log"
  "net/http"
)

func main() {
  var (
    addr = flag.String("addr", ":8080", "adres punktu końcowego")
    mongo = flag.String("mongo", "localhost", "adres mongodb")
  )
  log.Println("Nawiązywanie połączenia z bazą MongoDB", *mongo)
  db, err := mgo.Dial(*mongo)
  if err != nil {
    log.Fatalln("Nie udało się nazwiązać połączenia z MongoDB:", err)
  }
  defer db.Close()
  s := &Server{
    db: db,
  }
  mux := http.NewServeMux()
  mux.HandleFunc("/polls/", withCORS(withAPIKey(s.handlePolls)))
  log.Println("Uruchamianie serwera WWW pod adresem", *addr)
  http.ListenAndServe(":8080", mux)
  log.Println("Zatrzymywanie...")
}


type contextKey struct {
  name string
}

var contextKeyAPIKey = &contextKey{"klucz-api"}

func APIKey(ctx context.Context) (string, bool) {
  key, ok := ctx.Value(contextKeyAPIKey).(string)
  return key, ok
}


func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if !isValidAPIKey(key) {
      respondErr(w, r, http.StatusUnauthorized, "Nieprawidłowy klucz API")
      return
    }
    ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
    fn(w, r.WithContext(ctx))
  }
}


func isValidAPIKey(key string) bool {
  return key == "abc123"
}


func withCORS(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Expose-Headers", "Location")
    fn(w, r)
  }
}

// Server reprezentuje serwer API.
type Server struct {
  db *mgo.Session
}
