package data

import (
	g "GoDojoGo/deff"
	lib "GoDojoGo/lib"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type GetUserType struct {
	UserId int64  `db:"user_id"`
	Name   string `db:"name"`
	EMail  string `db:"email"`
	Role   string `db:"role"`
}

func GetUser(email string, pass string, admin bool) (*GetUserType, error) {
	//time.Sleep(6 * time.Second) // 2 saniye duraklatır
	db, err := lib.DbConnect()
	if err == nil {
		rList, err := lib.GenericQuery[GetUserType](db, "SELECT user_id, name, email, role FROM user WHERE email = ? AND password = MD5(?) LIMIT 1", email, pass)
		if err == nil {
			if len(rList) > 0 {
				if admin {
					if rList[0].Role == "ADMIN" {
						return &rList[0], nil
					} else {
						return nil, fmt.Errorf("user not found (%v:%v)", email, pass)
					}
				} else {
					return &rList[0], nil
				}
			} else {
				return nil, fmt.Errorf("user not found (%v:%v)", email, pass)
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

type AddUserType struct {
}

type UserDetailType struct {
	UserId       int64  `db:"user_id" json:"user_id"`
	Name         string `db:"name" json:"name"`
	BDate        string `db:"bdate" json:"bdate"`
	Gender       string `db:"gender" json:"gender"`
	Gsm          string `db:"gsm" json:"gsm"`
	Email        string `db:"email" json:"email"`
	Active       string `db:"active" json:"active"`
	PatmentModel string `db:"payment_model" json:"payment_model"`
}

type emailParamType struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func CreateUser(ud *UserDetailType) (int64, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()

		tx, err := db.Beginx()
		if err == nil {
			cmd := "INSERT INTO user (name, bdate, gender, gsm, email, active, payment_model) VALUES (:name, :bdate, :gender, :gsm, :email, :active, :payment_model)"
			result, err := tx.NamedExec(cmd, ud)
			if err == nil {
				liid, err := result.LastInsertId()
				if err == nil {
					barr, err := json.Marshal(&emailParamType{Name: ud.Name, Url: g.Settings.EMAIL_ACTIVATE_URL})
					if err == nil {
						cmd := `INSERT INTO emailpool (email, kind, params) VALUES (?, 'ACTIVATE', ?)`
						_, err := tx.Exec(cmd, ud.Email, string(barr))
						if err == nil {
							tx.Commit()
							return liid, nil
						} else {
							tx.Rollback()
							return 0, err
						}
					} else {
						tx.Rollback()
						return 0, err
					}
				} else {
					tx.Rollback()
					return 0, err
				}
			} else {
				return 0, err
			}
		} else {
			return 0, nil
		}
	} else {
		return 0, err
	}
}

func UpdateUser(ud *UserDetailType) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `UPDATE user 
					SET name = :name, bdate = :bdate, gender = :gender, gsm = :gsm, email = :email, active = :active, payment_model = :payment_model
						WHERE user_id = :user_id`
		_, err := db.NamedExec(cmd, ud)
		if err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func CreateOrUpdateUser(ud *UserDetailType) (int64, error) {
	if ud.UserId == 0 {
		return CreateUser(ud)
	} else {
		return ud.UserId, UpdateUser(ud)
	}
}

func SetUserPassword(user_id int64, pass string) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `UPDATE user SET "password" = MD5(?) WHERE user_id = ?`
		_, err := db.Exec(cmd, pass, user_id)
		return err
	} else {
		return err
	}
}

func GenerateUserGuid(user_id int64) (string, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		guid := uuid.New().String()
		cmd := `INSERT INTO userguid (user_id, guid) VALUES (?, ?)`
		_, err := db.Exec(cmd, user_id, guid)
		return guid, err
	} else {
		return "", err
	}
}
