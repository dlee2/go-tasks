package db

import (
    "log"
)

func CreateUser(username, password, email string) error {
    err := TaskQuery("insert into user(username, password, email) values(?,?,?)", username, password, email)
    return err
}

func ValidUser(username, password string) bool {
    var passwordFromDB string
    userSQL := "select password from user where username=?"
    log.Print("validating user ", username)
    rows := database.Query(userSQL, username)

    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&passwordFromDB)
        if err != nil {
            return false
        }
    }
    // if the password matches, return true
    if password == passwordFromDB {
        return true
    }

    return false
}

func GetUserID(username string) (int, error) {
    var userID int
    userSQL := "select id from user where username =?"
    rows := database.Query(userSQL, username)

    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&userID)
        if err != nil {
            return -1, err
        }
    }
    return userID, nil
}