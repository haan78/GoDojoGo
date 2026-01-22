package data

import (
	lib "GoDojoGo/lib"
	"database/sql"
	"errors"
)

type ActivityBaseType struct {
	ActivityId   int64           `db:"activity_id" json:"activity_id"`
	Name         string          `db:"name" json:"name"`
	ActivityDate sql.NullString  `db:"activity_date" json:"activity_date"`
	Fee          sql.NullFloat64 `db:"fee" json:"fee"`
	Text         sql.NullString  `db:"text" json:"text"`
}

func CreateActivty(a *ActivityBaseType) (int64, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `INSERT INTO activity (name, activity_date, fee, text) VALUES (:name, :activity_date, :fee, :text)`
		result, err := db.NamedExec(cmd, a)
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

func UpdateActivty(a *ActivityBaseType) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `SELECT activity_date FROM activity WHERE activity_id = ?`
		ad, err := lib.GetValue(db, cmd, a.ActivityId)
		if err == nil {
			if ns, ok := ad.(sql.NullString); ok {
				if ns.Valid == a.ActivityDate.Valid {
					cmd = `UPDATE activity SET name = :name, activity_date = :activity_date, fee = :fee, text = :text 
							WHERE activity_id = :activity_id`
					_, err := db.NamedExec(cmd, a)
					return err
				} else {
					return errors.New("activity date type can't change")
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

func DelActivty(activityId int64) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `DELETE FROM activity WHERE activity_id = ? AND deletable = 'YES'`
		_, err := db.Exec(cmd, activityId)
		return err
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
