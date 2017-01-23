package db

import (
    "database/sql"
    "log"
    "strconv"
    "strings"

    _ "github.com/mattn/go-sqlite3" //we want to use sqlite natively
    "github.com/dlee2/Tasks/types"
)

var database Database
var err error
var taskStatus map[string]int

type Database struct {
    db *sql.DB
}

//Begin a transaction
func (db Database) begin() (tx *sql.Tx) {
    tx, err := db.db.Begin()
    if err != nil {
        log.Println(err)
        return nil
    }
    return tx
}

func (db Database) prepare(q string) (stmt *sql.Stmt) {
    stmt, err := db.db.Prepare(q)
    if err != nil {
        log.Println(err)
        return nil
    }
    return stmt
}

func (db Database) query(q string, args ...interface{}) (rows *sql.Rows) {
    rows, err := db.db.Query(q, args...)
    if err != nil {
        log.Println(err)
        return nil
    }
    return rows
}

func init() {
    database.db, err = sql.Open("sqlite3", "./tasks.db")
    taskStatus = map[string]int{"COMPLETE": 1, "PENDING": 2, "DELETED": 3}
    if err != nil {
        log.Fatal(err)
    }
}

func Close() {
    database.db.Close()
}


// Retrieve the tasks
func GetTasks(username, status, category string) (types.Context, error) {
    log.Println("getting tasks for ", username)
    var tasks []types.Task
    var task types.Task
    var context types.Context
    var getTaskSQL string
    var rows *sql.Rows

    basicSQL := `select t.id, title, content, created_date, priority
                from task t, status s, user u where u.username=? and u.id=t.user_id `
    if category == "" {
        switch status {
            case "pending":
                getTaskSQL = basicSQL + " and s.status = 'PENDING' and t.hide !=1"
            case "deleted":
                getTaskSQL = basicSQL + " and s.status = 'DELETED' and t.hide !=1"
            case "completed":
                getTaskSQL = basicSQL + " and s.status = 'COMPLETED' and t.hide !=1"
        }
        getTaskSQL += " order by t.created_date asc"
        rows, err = database.db.Query(getTaskSQL, username)
    } else {
        status = category
        getTaskSQL = basicSQL + " and name = ?  and  s.status='PENDING'  order by priority desc, created_date asc, finish_date asc"
        rows, err = database.db.Query(getTaskSQL, username, category)

    }
    if err != nil {
        log.Println("Something went wrong with the query")
    }

    defer rows.Close()
    for rows.Next() {
        task = types.Task{}
        err = rows.Scan(&task.ID, &task.Title, &task.Content, &task.Created, &task.Priority)

        taskCompleted := 0
        totalTasks := 0

        if strings.HasPrefix(task.Content, "- [") {
            for _, value := range strings.Split(task.Content, "\n") {
                if strings.HasPrefix(value, "- [x]") {
                    taskCompleted += 1
                }
                totalTasks += 1
            }
            task.CompletedMsg = strconv.Itoa(taskCompleted) + " complete out of " + strconv.Itoa(totalTasks)
        }

        tasks = append(tasks, task)
    }

    context = types.Context{Tasks: tasks, Navigation: status}
    return context, nil
}


func AddTask(title, content, category, username string, priority, hide int) error {
    log.Println("AddTask: started function")
    var err error

    userID, err := GetUserID(username)
    if err != nil && (title != "" || content != "") {
        return err
    }
    if category == "" {
        err = TaskQuery(`insert into task(title, content, task_status_id, user_id, created_date, last_modified_at, priority, hide)
              values(?, ?, ?, ?, datetime(), datetime(), ?, ?)`, title, content, taskStatus["PENDING"], userID, priority, hide)
    } else {
        categoryID := GetCategoryIDByName(username, category)
        err = TaskQuery(`insert into task(title, content, task_status_id, cat_id, user_id, created_date, last_modified_at, priority, hide)
              values(?, ?, ?, ?, ?, datetime(), datetime(), ?, ?)`, title, content, taskStatus["PENDING"], categoryID, userID, priority, hide)
    }

    return err
}

func GetCategoryIDByName(username string, category string) int {
    var categoryID int
    getTaskSQL := "select c.id from category c , user u where u.id = c.user_id and name=? and u.username=?"

    rows := database.query(getTaskSQL, category, username)
    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&categoryID)
        if err != nil {
            log.Println(err)
        }
    }
    return categoryID
}

func TrashTask(username string, id int) error {
    err := TaskQuery("update task set task_status_id=?,last_modified_at=datetime() where user_id=(select id from user where username=?) and id=?", taskStatus["DELETED"], username, id)
    return err
}


func TaskQuery(sql string, args ...interface{}) error {
    log.Println("inside of task query")
    SQL := database.prepare(sql)
    tx := database.begin()
    _, err = tx.Stmt(SQL).Exec(args...)
    if err != nil {
        log.Println("TaskQuery: ", err)
        tx.Rollback()
    } else {
        err = tx.Commit()
        if err != nil {
            log.Println(err)
            return err
        }
        log.Println("Commit successful")
    }
    return err
}
