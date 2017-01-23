package views

import (
    "net/http"
    "log"

    "github.com/dlee2/Tasks/db"
)

func SignUpFunc(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()

        username := r.Form.Get("username")
        password := r.Form.Get("password")
        email := r.Form.Get("email")

        log.Println(username, password, email)

        err := db.CreateUser(username, password, email)
        if err != nil {
            http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
        } else {
            http.Redirect(w, r, "/login", 302)
        }
    }
}