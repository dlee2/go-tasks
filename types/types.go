package types

type Task struct {
    ID      int
    Title   string
    Content string
    Created string
    Priority string
    Category string
    CompletedMsg string

}

type Tasks []Task

type Context struct {
    Tasks      []Task
    Navigation string
    Search     string
    Message    string
    CSRFToken  string

}

type Status struct {
    StatusCode int    `json:"status_code"`
    Message    string `json:"message"`
}

type Category struct {
    ID int `json:"category_id"`
    Name string `json:"category_name"`
    Created string `json:"created_date"`
}

type Categories []Category