package data

import (
	lib "GoDojoGo/lib"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
						return nil, fmt.Errorf("user not found (%v)", email)
					}
				} else {
					return &rList[0], nil
				}
			} else {
				return nil, fmt.Errorf("user not found (%v)", email)
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

func CreateUser(ud *UserDetailType) (int64, error) {

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()

		tx, err := db.Beginx()
		if err == nil {
			cmd := "INSERT INTO user (name, bdate, gender, gsm, email, active, payment_model) VALUES (:name, :bdate, :gender, :gsm, :email, :active, :payment_model)"
			result, err := tx.NamedExec(cmd, ud)
			if err == nil {
				user_id, err := result.LastInsertId()
				if err == nil {
					_, err := guidAndEmail(tx, ud.Name, ud.Email, "ACTIVATE", "", user_id)
					return user_id, err
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

func ChangeUserPassword(email, oldp, newp string) error {
	db, err := lib.DbConnect()
	if err == nil {
		result, err := db.Exec(`UPDATE user SET password = MD5(?) WHERE email = ? AND password = MD5(?) LIMIT 1`, newp, email, oldp)
		if err == nil {
			arc, err := result.RowsAffected()
			if err == nil {
				if arc == 1 {
					return nil
				} else {
					return fmt.Errorf("Member not found")
				}
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func guidAndEmail(tx *sqlx.Tx, name, email, kind, pass string, user_id int64) (string, error) {

	type emailParmaType struct {
		Name string `json:"name"`
		Code string `json:"code"`
		Guid string `json:"guid"`
	}

	sql1 := `SELECT COUNT(*) FROM userguid WHERE user_id = ? AND kind = ? AND status <> 'VERIFIED' AND expire_time > NOW()`
	sql2 := `INSERT INTO userguid (user_id, guid, code, kind, email, params, expire_time, status) VALUES (?, ?, ?, ?, ADDTIME(NOW(),'01:00:00'),'WAITING')`
	if pass != "" {
		sql2 = `INSERT INTO userguid (user_id, guid, code, kind, email, params, password, expire_time, status) VALUES (?, ?, ?, ?, MD5(?), ADDTIME(NOW(),'01:00:00'),'WAITING')`
	}

	code, err := lib.GenerateRandomCode()
	if err == nil {
		guid := uuid.New().String() + "-" + fmt.Sprint(user_id)
		jdata, err := json.Marshal(&emailParmaType{Name: name, Code: code, Guid: guid})
		result, err := tx.Query(sql1, user_id, kind)
		if err == nil {
			var c int64
			err := result.Scan(&c)
			if err == nil {
				if c == 0 {
					_, err := tx.Exec(sql2, user_id, guid, code, kind, email, string(jdata))
					return guid, err
				} else {
					return "", fmt.Errorf("there is a waiting proccess")
				}
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}

}

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
						return guid, tx.Commit()
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
				_, err := db.Exec("UPDATE user SET password = ? WHERE user_id = ?", qData[0].Password, qData[0].UserId)
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

func UserAll() ([]UserDetailType, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		list, err := lib.GenericQuery[UserDetailType](db, `SELECT user_id, name, bdate, gender, gsm, email, active, payment_model fROM user`)
		if err == nil {
			return list, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
