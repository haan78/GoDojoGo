package data

import (
	lib "GoDojoGo/lib"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ActivityBaseType struct {
	ActivityId int64              `db:"activity_id" json:"activity_id"`
	Name       string             `db:"name" json:"name"`
	Start      lib.JSONNullString `db:"start" json:"start"`
	End        lib.JSONNullString `db:"end" json:"end"`
	SingleFee  float64            `db:"single_fee" json:"single_fee"`
	WorkerFee  float64            `db:"worker_fee" json:"worker_fee"`
	StudentFee float64            `db:"student_fee" json:"student_fee"`
	Text       lib.JSONNullString `db:"text" json:"text"`
	Repetitive string             `db:"repetitive" json:"repetitive"`
	Active     string             `db:"active" json:"active"`
	WeekDays   lib.JSONNullString `db:"wdays" json:"wdays"`
	Date       lib.JSONDate       `db:"date" json:"date"`
}

func paresWeekDays(weekdayStr string) ([]int, error) {
	seen := make(map[int]bool)
	var weekdays []int
	for _, ch := range weekdayStr {
		day := int(ch - '0')
		if day < 0 || day > 6 {
			return nil, fmt.Errorf("invalid weekday: %c", ch)
		}
		if !seen[day] {
			seen[day] = true
			weekdays = append(weekdays, day)
		}
	}
	return weekdays, nil
}

func addToCalender(ts *sqlx.Tx, a *ActivityBaseType, liid int64) error {
	switch a.Repetitive {
	case "WEEKLY":
		days, err := lib.NearestWeekdays(a.Date.Time, a.WeekDays.String)
		if err == nil {
			str := ""
			for _, day := range days {
				if str != "" {
					str += ","
				}
				str += fmt.Sprintf(`('%s',%v)`, day, liid)
			}
			_, err := ts.Exec(`INSERT INTO calendar (date, activity_id) VALUES ` + str)
			if err == nil {
				ts.Commit()
				return nil
			} else {
				ts.Rollback()
				return err
			}
		} else {
			ts.Rollback()
			return err
		}
	case "MONTHLY":
		wdays, err := paresWeekDays(a.WeekDays.String)
		if err == nil {
			if len(wdays) == 2 {
				nday, err := lib.InNextMonth(a.Date.Time, 0, wdays[0], wdays[1])
				if err == nil {
					_, err := ts.Exec(`INSERT INTO calendar (date, activity_id) VALUES (?, ?)`, nday, liid)
					if err == nil {
						ts.Commit()
						return nil
					} else {
						ts.Rollback()
						return err
					}
				} else {
					ts.Rollback()
					return err
				}
			} else {
				ts.Rollback()
				return fmt.Errorf("weekdays parameter is wrong")
			}
		} else {
			ts.Rollback()
			return err
		}
	default: //if a.Repetitive == "NO"
		_, err := ts.Exec(`INSERT INTO calendar (date, activity_id) VALUES (?, ?)`, a.Date, liid)
		if err == nil {
			ts.Commit()
			return nil
		} else {
			ts.Rollback()
			return err
		}
	}
}

func ActivityInsert(a *ActivityBaseType) (int64, error) {
	insSql := `INSERT INTO activity (name, start, end, single_fee, worker_fee, student_fee, text, repetitive, :wdays, active) 
					VALUES (:name, :start, :end, :single_fee, :worker_fee, :student_fee,  :text, :repetitive, :wdays, :active)`
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		ts, err := db.Beginx()
		if err == nil {
			result, err := ts.NamedExec(insSql, a)
			if err == nil {
				liid, err := result.LastInsertId()
				if err == nil {
					err := addToCalender(ts, a, liid)
					if err == nil {
						ts.Commit()
						return liid, nil
					} else {
						ts.Rollback()
						return 0, err
					}
				} else {
					ts.Rollback()
					return 0, err
				}
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

func ActivityUpdate(a *ActivityBaseType) error {

	cmd := `UPDATE activity SET 
				name = :name, 
				start = :start, 
				end = :end, 
				single_fee = :single_fee, 
				worker_fee = :worker_fee, 
				student_fee = :student_fee, 
				text = :text, 
				repetitive = :repetitive, 
				wdays = :wdays, 
				active = :active 
					WHERE activity_id = :activity_id`

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		_, err := db.NamedExec(cmd, a)
		if err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func ActivityList(active string) ([]ActivityBaseType, error) {

	getActive := func(a string) string {
		switch a {
		case "":
			return "a.active"
		case "ALL":
			return "a.active"
		case "NO":
			return "'NO'"
		default:
			return "'YES'"
		}
	}
	result := []ActivityBaseType{}
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `SELECT activity_id, name, start, end, single_fee, worker_fee, student_fee, text, repetitive, wdays, active FROM activity a WHERE a.active = %s`
		cmd = fmt.Sprintf(cmd, getActive(active))
		result, err = lib.GenericQuery[ActivityBaseType](db, cmd)
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
