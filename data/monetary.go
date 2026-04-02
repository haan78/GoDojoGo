package data

import (
	"GoDojoGo/lib"
)

type MonetaryRecord struct {
	MonetaryId int64         `json:"monetary_id" db:"monetary_id"`
	Fee        float64       `json:"fee" db:"fee"`
	Payment    lib.JSONFloat `json:"payment" db:"payment"`
	Date       lib.JSONDate  `json:"date" db:"date"`
	ActivityId lib.JSONInt   `json:"activity_id" db:"activity_id"`
	UserId     lib.JSONInt   `json:"user_id" db:"user_id"`
}

type ExpenseType struct {
	MonetaryId int64        `json:"monetary_id" db:"monetary_id"`
	Payment    float64      `json:"payment" db:"payment"`
	Date       lib.JSONDate `json:"date" db:"date"`
	ActivityId lib.JSONInt  `json:"activity_id" db:"activity_id"`
	UserId     lib.JSONInt  `json:"user_id" db:"user_id"`
	Text       string       `json:"text" db:"text"`
}

type SellType struct {
	Payment float64      `json:"payment" db:"payment"`
	Text    string       `json:"text" db:"text"`
	Date    lib.JSONDate `json:"date" db:"date"`
}

func Sell(s SellType) (int64, error) {
	db, err := lib.DbConnect()
	if err == nil {
		result, err := db.Exec(`INSERT INTO monetary (fee, payment, date, text) VALUES (?, ?, ?, ?)`, s.Payment, s.Payment, s.Date, s.Text)
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

func MonetaryRecordDelByActivity(activitiyId, userId int64) error {
	db, err := lib.DbConnect()
	if err == nil {
		_, err := db.Exec(`DELETE FROM monetary WHERE user_id = ? AND activity_id = ? AND payment IS NULL AND type = 'INCOME'`, userId, activitiyId)
		return err
	} else {
		return err
	}
}

func MonetaryAddExpense(ex ExpenseType) (int64, error) {
	db, err := lib.DbConnect()
	if err == nil {
		result, err := db.Exec(`INSERT INTO monetary (type, activity_id, user_id, date, text, payment, fee) VALUES ('EXPENSE',? ,?, ?, ?, ?, ?)`, ex.ActivityId, ex.UserId, ex.Date, ex.Text, ex.Payment, ex.Payment)
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
