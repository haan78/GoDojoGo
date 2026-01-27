package data

import (
	"GoDojoGo/lib"
	"database/sql"
)

type feeType struct {
	Single  float64 `db:"single_fee"`
	Worker  float64 `db:"worker_fee"`
	Student float64 `db:"student_fee"`
}

type AddUserActivityType struct {
	UserId   int64           `db:"user_id" json:"user_id"`
	Activity string          `db:"activity" json:"activity"`
	Fee      float64         `db:"fee" json:"fee"`
	Payment  sql.NullFloat64 `db:"payment" json:"payment"`
	Date     sql.NullString  `db:"date" json:"date"`
	Time     sql.NullString  `db:"time" json:"time"`
}

type GetUserActivitiesType struct {
	UserActivityId int64           `db:"useractivity_id" json:"useractivity_id"`
	Activity       string          `db:"activity" json:"activity"`
	Fee            float64         `db:"fee" json:"fee"`
	Payment        sql.NullFloat64 `db:"payment" json:"payment"`
	Date           sql.NullString  `db:"date" json:"date"`
	Time           sql.NullString  `db:"time" json:"time"`
}

type UserActivityPaymentType struct {
	UserActivityId int64   `db:"useractivity_id" json:"useractivity_id"`
	Payment        float64 `db:"payment" json:"payment"`
}

func AddUserActivity(d *AddUserActivityType) (int64, error) {
	db, err := lib.DbConnect()
	if err == nil {
		result, err := db.NamedExec(`INSERT INTO useractivity (user_id, activity, fee, payment, "date", "time") VALUES (:user_id, :activity, :fee, :payment, :date, :time)`, d)
		if err == nil {
			id, err := result.LastInsertId()
			if err == nil {
				return id, nil
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

func DelUserActivity(useractivity_id int64) error {
	db, err := lib.DbConnect()
	if err == nil {
		_, err := db.Exec("DELETE FROM useractivity WHERE useractivity_id = ?", useractivity_id)
		return err
	} else {
		return err
	}
}

func GetUserActivities(userId int64) ([]GetUserActivitiesType, error) {
	db, err := lib.DbConnect()
	if err == nil {
		cmd := `SELECT
					ua.useractivity_id, 
					ua.activity, 
					ua.fee, 
					ua.payment, 
					ua.date, 
					ua.time 
						FROM useractivity ua 
							WHERE ua.user_id = ? 
								ORDER BY ua.date DESC, ua.time ASC`
		result, err := lib.GenericQuery[GetUserActivitiesType](db, cmd, userId)
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func UserActivityPayment(d *UserActivityPaymentType) error {
	db, err := lib.DbConnect()
	if err == nil {
		_, err := db.NamedExec("UPDATE useractivity ua SET ua.payment = :payment WHERE useractivity_id = :useractivity_id", d)
		return err
	} else {
		return err
	}
}
