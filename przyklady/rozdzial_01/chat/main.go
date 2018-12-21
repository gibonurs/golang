package main

import (
  "flag"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "sync"
  "text/template"

  "../trace"
)

// Struktura reprezentująca pojedyczny szablon.
type templateHandler struct {
  once     sync.Once
  filename string
  templ    *template.Template
}

// Metoda ServeHTTP obsługuje żądania HTTP.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  })
  t.templ.Execute(w, r)
}

func main() {
  var addr = flag.String("addr", ":8080", "Adres aplikacji internetowej.")
  flag.Parse() // analiza flag wiersza poleceń

  r := newRoom()
  r.tracer = trace.New(os.Stdout)

  http.Handle("/", &templateHandler{filename: "chat.html"})
  http.Handle("/room", r)

  // uruchomienie pokoju rozmów
  go r.run()

  // uruchomienie serwera WWW
  log.Println("Uruchamianie serwera WWW pod adresem", *addr)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }

}
