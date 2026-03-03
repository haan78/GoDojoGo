package data

import (
	lib "GoDojoGo/lib"
	"database/sql"
	"fmt"
)

func GenerateUserGuidForForgotPass(email, new_pass string) (string, error) {

	type userInfo struct {
		Name   string `db:"name"`
		UserId int64  `db:"user_id"`
	}

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		users, err := lib.GenericQuery[userInfo](db, "SELECT name, user_id FROM user WHERE email = ?", email)
		if err == nil {
			if len(users) == 1 {
				tx, err := db.Beginx()
				if err == nil {
					guid, err := guidAndEmail(tx, users[0].Name, email, "PASSWORD", new_pass, users[0].UserId)
					if err == nil {
						tx.Commit()
						return guid, nil
					} else {
						tx.Rollback()
						return "", err
					}
				} else {
					return "", err
				}
			} else {
				return "", fmt.Errorf("member not found")
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func SetNewPasswordByGuid(guid string) error {

	type IdAndPassType struct {
		UserId   int64          `db:"user_id"`
		Password sql.NullString `db:"password"`
	}

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		qData, err := lib.GenericQuery[IdAndPassType](db, "SELECT user_id, password FROM userguid WHERE guid = ? AND expire_time > NOW()", guid)
		if err == nil {
			if len(qData) == 1 {
				_, err := db.Exec("UPDATE user SET password = ?, status = 'COMPLETED' WHERE user_id = ?", qData[0].Password, qData[0].UserId)
				return err
			} else {
				return fmt.Errorf("guid not found")
			}
		} else {
			return err
		}

	} else {
		return err
	}
}
