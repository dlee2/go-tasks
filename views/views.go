// UploadedFileHandler is used to handle the uploaded file related request
package views

import (
    "crypto/md5"
    "fmt"
    "html/template"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/dlee2/Tasks/db"
    "github.com/dlee2/Tasks/sessions"
)


var homeTemplate *template.Template
var deletedTemplate *template.Template
var completedTemplate *template.Template
var loginTemplate *template.Template
var editTemplate *template.Template
var searchTemplate *template.Template
var templates *template.Template

var message string // message will store the message to be shown as notification
var err error


// PopulateTemplates is used to parse all templates present in
// the templates folder
func PopulateTemplates() {
    var allFiles []string
    templatesDir := "./public/static/templates/"
    files, err := ioutil.ReadDir(templatesDir)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    for _, file := range files {
        filename := file.Name()
        if strings.HasSuffix(filename, ".html") {
            allFiles = append(allFiles, templatesDir + filename)
        }
    }

    templates, err = template.ParseFiles(allFiles...)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    homeTemplate = templates.Lookup("home.html")
    deletedTemplate = templates.Lookup("deleted.html")
    editTemplate = templates.Lookup("edit.html")
    searchTemplate = templates.Lookup("search.html")
    completedTemplate = templates.Lookup("completed.html")
    loginTemplate = templates.Lookup("login.html")

}

func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        username := sessions.GetUserName(r)
        context, err := db.GetTasks(username, "pending", "")
        log.Println(context)
        if err != nil {
            http.Redirect(w, r, "/", http.StatusInternalServerError)
        } else {
            context.CSRFToken = "abcd"
            message = ""
            expiration := time.Now().Add(365 * 24 * time.Hour)
            cookie := http.Cookie{Name: "csrftoken", Value: "abcd", Expires: expiration}
            http.SetCookie(w, &cookie)
            log.Println("Rendering template")
            homeTemplate.Execute(w, context)

        } 
    }
}

func AddTaskFunc(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            log.Println(err)
        }

        taskPriority, err := strconv.Atoi(r.FormValue("priority"))
        if err != nil {
            log.Print("Unable to convert priority to integer")
        }

        priorityList := []int{1,2,3}

        found := false
        for _, priority := range priorityList {
            if taskPriority == priority {
                found = true
            }
        }
        // check if the task priority radio button was selected
        if !found {
            taskPriority = 1
        }

        hide, err := strconv.Atoi(r.FormValue("hide"))
        if err != nil {
            log.Print("Unable to convert hide to integer")
        }
        title := template.HTMLEscapeString(r.Form.Get("title"))
        content := template.HTMLEscapeString(r.Form.Get("content"))
        formToken := template.HTMLEscapeString(r.Form.Get("CSRFToken"))

        cookie, _ := r.Cookie("csrftoken")
        if (formToken == cookie.Value){

            username := sessions.GetUserName(r)

            if handler != nil {
                r.ParseMultipartForm(32 << 20) //defined as max file size
                defer file.Close()

                randomFileName := md5.New()
                io.WriteString(randomFileName, strconv.FormatInt(time.Now().Unix(), 10))
                io.WriteString(randomFileName, handler.Filename)
                token := fmt.Sprintf("%x", randomFileName.Sum(nil))
                f, err := os.OpenFile("./files/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
                if err != nil {
                    log.Println(err)
                    return
                }
                defer f.Close()
                io.Copy(f, file)


                filelink := "<br> <a href=/files/" + handler.Filename + ">" + handler.Filename + "</a>"
                content = content + filelink

                fileTruth := db.AddFile(handler.Filename, token, username)
                if fileTruth != nil {
                    message = "Error adding filename in db"
                    log.Println("error adding task to db")
                }
            }

            truth := db.AddTask(title, content, "", username, taskPriority, hide)
            if truth != nil {
                message = "Error adding task"
                log.Println("error adding task to the db")
            } else {
                message = "Task added"
                log.Println("added task to db")
            }
            http.Redirect(w, r, "/", http.StatusFound)
        } else {
            log.Fatal("CSRF mismatch")
        }
    } else {
        message = "Method not allowed"
        http.Redirect(w, r, "/", http.StatusFound)
    }
}

func TrashTaskFunc(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        id, err := strconv.Atoi(r.URL.Path[len("/trash/"):])
        if err != nil {
            log.Println("TrashFunc", err)
        } else {
            username := sessions.GetUserName(r)
            err := db.TrashTask(username, id)
            if err != nil {
                log.Println("Error trashing task ", id)
            } else {
                http.Redirect(w, r, "/", http.StatusFound)
            }
        }
    }
}

func UploadedFileHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        log.Println("into the handler")
        token := r.URL.Path[len("/files/"):]

        log.Println("serving the file ./files/" + token)
        http.ServeFile(w, r, "./files/" + token)
    }
}