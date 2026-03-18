package data

import (
	lib "GoDojoGo/lib"
	"database/sql"
	"errors"
	"fmt"
)

type ActivityBaseType struct {
	ActivityId int64          `db:"activity_id" json:"activity_id"`
	Name       string         `db:"name" json:"name"`
	Date       string         `db:"date" json:"date"`
	Start      sql.NullString `db:"start" json:"start"`
	End        sql.NullString `db:"end" json:"end"`
	SingleFee  float64        `db:"single_fee" json:"single_fee"`
	WorkerFee  float64        `db:"worker_fee" json:"worker_fee"`
	StudentFee float64        `db:"student_fee" json:"student_fee"`
	Text       sql.NullString `db:"text" json:"text"`
	Repetitive string         `db:"repetitive" json:"repetitive"`
	Active     string         `db:"active" json:"active"`
}

func SetActivity(a *ActivityBaseType) (int64, error) {

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		var cmd string
		if a.ActivityId == 0 {
			cmd = `INSERT INTO activity (name, date, start, end, single_fee, worker_fee, student_fee, text, repetitive, active) 
					VALUES (:name, :date, :start, :end, :single_fee, :worker_fee, :student_fee,  :text, :repetitive, :active)`
		} else {
			cmd = `UPDATE activity SET 
				name = :name, date = :date, start = :start, end = :end, text = :text, repetitive = :repetitive, active = :active
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

func DelActivity(activityId int64) error {
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

func GetActivity() ([]ActivityBaseType, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `SELECT name, date, start, end, single_fee, worker_fee, student_fee, text, repetitive, active FROM activity WHERE activity_date >= CURDATE()`
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

func doRepetitiveActivity() ([]ActivityBaseType, error) {
	selcmd := `SELECT
	a.name, 
	IF(a.repetitive = 'WEEKLY',DATE_ADD(a.date, INTERVAL 1 WEEK), DATE_ADD(a.date, INTERVAL 1 MONTH)) as date,
	a.start, 
	a.end,
	a.single_fee,
	a.worker_fee,
	a.student_fee,
	a.text,
	a.repetitive,
	a.active
		FROM activity a
		LEFT JOIN activity _a ON _a.active = 1 AND _a.repetitive = a.repetitive AND _a.name = a.name AND _a.date > a.date
			WHERE
				_a.activity_id IS NULL
				AND a.active = 1 
				AND a.repetitive IN ('WEEKLY','MONTHLY') 
				AND a.date >= IF(a.repetitive = 'WEEKLY',DATE_ADD(CURDATE(), INTERVAL -1 WEEK), DATE_ADD(CURDATE(), INTERVAL -1 MONTH))
				AND a.date < CURDATE()
				AND CURDATE() = IF(a.repetitive = 'WEEKLY',DATE_ADD(a.date, INTERVAL 1 WEEK), DATE_ADD(a.date, INTERVAL 1 MONTH))`

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		result, err := lib.GenericQuery[ActivityBaseType](db, selcmd)
		if err == nil {
			values := ""
			for i := 0; i < len(result); i++ {

				start := "NULL"
				if result[i].Start.Valid {
					start = "'" + result[i].Start.String + "'"
				}

				end := "NULL"
				if result[i].End.Valid {
					end = "'" + result[i].End.String + "'"
				}

				text := "NULL"
				if result[i].Text.Valid {
					text = "'" + result[i].Text.String + "'"
				}

				if i > 0 {
					values += ","
				}

				values += fmt.Sprintf(`('%s', '%s', %s, %s, %v, %v, %v, '%s', '%s', %v)`,
					result[i].Name, result[i].Date, start, end, result[i].SingleFee, result[i].WorkerFee, result[i].StudentFee, text, result[i].Repetitive, result[i].Active)
			}
			inscmd := `INSERT INTO activity (name, date, start, end, single_fee, worker_fee, student_fee, text, repetitive, active) VALUES ` + values
			_, err := db.Exec(inscmd)
			if err == nil {
				return result, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
