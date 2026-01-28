package data

import (
	lib "GoDojoGo/lib"
	"database/sql"
	"errors"
)

type ActivityBaseType struct {
	ActivityId int64          `db:"activity_id" json:"activity_id"`
	Name       string         `db:"name" json:"name"`
	Date       sql.NullString `db:"date" json:"date"`
	Start      sql.NullString `db:"start" json:"start"`
	End        sql.NullString `db:"end" json:"end"`
	SingleFee  float64        `db:"single_fee" json:"single_fee"`
	WorkerFee  float64        `db:"worker_fee" json:"worker_fee"`
	StudentFee float64        `db:"student_fee" json:"student_fee"`
	Text       sql.NullString `db:"text" json:"text"`
	Repetitive string         `db:"repetitive" json:"repetitive"`
	Active     string         `db:"active" json:"active"`
}

func SetActivty(a *ActivityBaseType) (int64, error) {

	if a.Date.Valid && a.Repetitive == "YES" {
		return 0, errors.New("repetitive activity can't have date")
	}

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		var cmd string
		if a.ActivityId == 0 {
			cmd = `INSERT INTO activity (name, date, start, end, single_fee, worker_fee, student_fee, text, repetitive, active) 
					VALUES (:name, :date, :start, :end, :single_fee, :worker_fee, :student_fee,  :text, :repetitive, :active)`
		} else {
			cmd = `UPDATE activity SET 
				name = :name, date = :date, start = :start, end = :end, text = :text, active = :active
					WHERE activity_id = :activity_id`
		}

		result, err := db.NamedExec(cmd, a)
		if err == nil {
			if a.ActivityId == 0 {
				liid, err := result.LastInsertId()
				if err == nil {
					return liid, nil
				} else {
					return 0, err
				}
			} else {
				return a.ActivityId, nil
			}
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

func DelActivty(activityId int64) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()

		count, err := lib.GetInt(db, "SELECT COUNT(1) FROM useractivity WHERE payment IS NOT NULL AND activity_id = ?", activityId)
		if err != nil {
			if count == 0 {
				cmd := `DELETE FROM activity WHERE activity_id = ? AND deletable = 'YES'`
				_, err := db.Exec(cmd, activityId)
				return err
			} else {
				return errors.New("some user made paymemt this activity")
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func GetActivty() ([]ActivityBaseType, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `SELECT activity_id, name, activity_date, fee, text FROM activity WHERE activity_date IS NULL OR activity_date >= CURDATE()`
		result, err := lib.GenericQuery[ActivityBaseType](db, cmd)
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
