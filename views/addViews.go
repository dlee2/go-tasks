package views

// import (
//     "log"
//     "net/http"
//     "strconv"

//     "github.com/dlee2/Tasks/sessions"
// )


// func AddComentFunc(w http.ResponseWriter, r *http.Request) {
//     if sessions.IsLoggedIn(r) {
//         if r.Method == "POST" {
//             r.ParseForm()
//             text := r.Form.Get("commentText")
//             id := r.Form.Get("taskID")

//             idInt, err := strconv.Atoi(id)

//             if (err != nil) || (text == "") {
//                 log.Println("unable to convert into integer")
//                 message = "Error adding comment"
//             } else {
//                 // err = db.AddComments(idInt, text)

//                 if err != nil {
//                     log.Println("unable to insert into db")
//                     message = "Comment not added"
//                 } else {
//                     message = "Comment added"
//                 }
//             }

//             http.Redirect(w, r, "/", http.StatusFound)
//         }
//     } else {
//         http.Redirect(w, r, "/login", 302)
//     }
// }