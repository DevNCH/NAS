package database

import "fmt"

func Migrate() error {

	db := GetDB()

	_, err := db.Exec(CreateUsersTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(CreateFilesTable)
	if err != nil {
		return err
	}

	fmt.Println("Migrações executadas!")

	return nil
}