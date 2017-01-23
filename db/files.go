package db


func AddFile(fileName, token, username string) error {
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	err = TaskQuery("insert into files values(?,?,?,datetime())", fileName, token, userID)

	return err
}