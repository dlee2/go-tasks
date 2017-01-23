package main

import (
    "flag"
    "log"
    "net/http"
    "strings"

    "github.com/dlee2/Tasks/config"
    "github.com/dlee2/Tasks/views"
)


func main() {
    values, err := config.ReadConfig("config.json")
    var port *string;

    if err != nil {
        port = flag.String("port", "", "IP address")
        flag.Parse()

        if !strings.HasPrefix(*port, ":") {
            *port = ":" + *port
            log.Println("port is " + *port)
        }

        values.ServerPort = *port
    }

    views.PopulateTemplates()

    http.HandleFunc("/login/", views.LoginFunc)
    http.HandleFunc("/logout/", views.RequiresLogin(views.LogoutFunc))
    http.HandleFunc("/signup/", views.SignUpFunc)

    http.HandleFunc("/add/", views.RequiresLogin(views.AddTaskFunc))
    // http.handleFunc("/complete/", CompleteTask)
    // http.handleFunc("/delete/", DeleteTask)
    // http.handleFunc("/deleted/", ShowDeletedTasks)
    http.HandleFunc("/trash/", views.RequiresLogin(views.TrashTaskFunc))
    // http.handleFunc("/edit/", EditTask)
    // http.handleFunc("/completed/", ShowCompletedTasks)
    // http.handleFunc("/restore/", RestoreTask)
    // http.handleFunc("/add/", AddTask)
    // http.handleFunc("/update/", UpdateTask)
    // http.handleFunc("/search/", SearchTask)
    // http.handleFunc("/login", Login)
    // http.handleFunc("/register", PostRegister)
    // http.handleFunc("/admin", HandleAdmin)
    // http.handleFunc("/add_user", PostAddUser)
    // http.handleFunc("/change", PostChange)
    
    http.HandleFunc("/", views.RequiresLogin(views.ShowAllTasksFunc))
    

    http.Handle("/static/", http.FileServer(http.Dir("public")))
    PORT:= values.ServerPort
    log.Print("Running server on " + PORT)
    log.Fatal(http.ListenAndServe(PORT, nil))
}
