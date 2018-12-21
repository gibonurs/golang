package web

import (
	"html/template"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func init() {
	tmpl, err := template.ParseGlob("templates/*.tmpl.html")
	if err != nil {
		http.Handle("/", errHandler(err.Error(), http.StatusInternalServerError))
		return
	}
	http.Handle("/questions/", templateHandler(tmpl, "question"))
	http.Handle("/", templateHandler(tmpl, "index"))
}

func templateHandler(tmpl *template.Template, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		if ok := ensureUser(ctx, w, r); !ok {
			return
		}
		err := tmpl.ExecuteTemplate(w, name, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// errHandler pobiera http.Handler który zgłasza konkretny błąd.
func errHandler(err string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err, code)
	})
}

// ensureUser upewnia się, że użytkownik jest zalogowany lub przekierwuje go
// na stronę logowania.
// Zwraca true jeśli użytkownik jest zalogowany.
func ensureUser(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	me := user.Current(ctx)
	if me == nil {
		loginURL, err := user.LoginURL(ctx, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
		return false
	}
	return true
}
