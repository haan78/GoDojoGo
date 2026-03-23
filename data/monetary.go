package data

import (
	"GoDojoGo/lib"
	"fmt"
)

type MonetaryRecord struct {
	MonetaryId int64         `json:"monetary_id" db:"monetary_id"`
	Fee        float64       `json:"fee" db:"fee"`
	Payment    lib.JSONFloat `json:"payment" db:"payment"`
	Date       lib.JSONDate  `json:"date" db:"date"`
	ActivityId lib.JSONInt   `json:"activity_id" db:"activity_id"`
	UserId     lib.JSONInt   `json:"user_id" db:"user_id"`
}

type UserOfActiviyt struct {
	Name string
}

func MonetaryRecordAdd(mr MonetaryRecord) (int64, error) {
	db, err := lib.DbConnect()
	if err == nil {
		result, err := db.NamedExec(`INSERT INTO monetary (user_id, activity_id, fee, payment, date) VALUES (:user_id, :activity_id, :fee, :payment, :date)`, mr)
		if err == nil {
			liid, err := result.LastInsertId()
			if err == nil {
				return liid, nil
			} else {
				return 0, err
			}
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

func MonetaryRecordDel(MonetaryId int64) error {
	db, err := lib.DbConnect()
	if err == nil {
		_, err := db.Exec(`DELETE FROM monetary WHERE monetary_id = ?`, MonetaryId)
		return err
	} else {
		return err
	}
}

func MonetaryRecordAddByActivity(activityId, userId int64) (MonetaryRecord, error) {
	var mr MonetaryRecord = MonetaryRecord{}
	db, err := lib.DbConnect()
	if err == nil {
		pm, err := lib.GetString(db, `SELECT payment_model FROM user WHERE user_id = ?`, userId)
		if err == nil {
			fn := ""
			var fee float64 = 0
			if pm == "STUDENT" {
				fn = "student_fee"
			} else if pm == "WORKER" {
				fn = "worker_fee"
			} else if pm == "SINGLE" {
				fn = "single_fee"
			} else if pm != "FREE" {
				return mr, fmt.Errorf("unknown payment model (%s)", pm)
			}
			if fn != "" {
				f, err := lib.GetFloat(db, fmt.Sprintf("SELECT %s FROM activity WHERE activity_id = ?", fn), activityId)
				if err == nil {
					fee = f
				} else {
					return mr, err
				}
			}

			mr.ActivityId = lib.JSONIntGet(activityId)
			mr.Date = lib.JSONDateNil()
			mr.Fee = fee
			mr.MonetaryId = 0
			mr.Payment = lib.JSONFloatNil()
			mr.UserId = lib.JSONIntGet(userId)
			mr.Payment.Valid = false

			c, err := lib.GetInt(db, `SELECT COUNT(1) FROM monetary WHERE user_id = ? AND activity_id = ?`)
			if err == nil {
				if c == 0 {
					result, err := db.NamedExec(`INSERT INTO monetary (user_id, activity_id, fee, payment, date) VALUES (:user_id, :activity_id, :fee, :payment, :date)`, mr)
					if err == nil {
						liid, err := result.LastInsertId()
						if err == nil {
							mr.MonetaryId = liid
							return mr, nil
						} else {
							return mr, err
						}
					} else {
						return mr, err
					}
				} else {
					return mr, fmt.Errorf("member already registred on this activity")
				}
			} else {
				return mr, err
			}

		} else {
			return mr, err
		}
	} else {
		return mr, err
	}
}

func MonetaryRecordDelByActivity(activitiyId, userId int64) error {
	db, err := lib.DbConnect()
	if err == nil {
		_, err := db.Exec(`DELETE FROM monetary WHERE user_id = ? AND activity_id = ?`, userId, activitiyId)
		return err
	} else {
		return err
	}
}
