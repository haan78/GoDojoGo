package data

import (
	lib "GoDojoGo/lib"
	"fmt"
)

type CalendarBaseType struct {
	CalendarId int64        `json:"calendar_id" db:"calendar_id"`
	ActivityId int64        `json:"activity_id" db:"activity_id"`
	Date       lib.JSONDate `json:"date" db:"date"`
	Name       string       `json:"name" db:"name"`
}

type CalendarReportType struct {
	CalendarBaseType
	Repetitive   string  `db:"repetitive" json:"repetitive"`
	Active       string  `db:"active" json:"active"`
	Participants int     `json:"participants" db:"participants"`
	TotalIncome  float64 `json:"total_income" db:"total_income"`
	TotalExpense float64 `json:"total_expense" db:"total_expense"`
}

func CalendarDel(sDate string, activityId int) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		count, err := lib.GetInt(db, `SELECT COUNT(*) FROM calendar c INNER JOIN monetary m ON m.calendar_id = c.calendar_id WHERE c.activity_id = ? AND c.date = ?`, activityId, sDate)
		if err == nil {
			if count == 0 {
				ts, err := db.Beginx()
				if err == nil {
					_, err := ts.Exec(`DELETE FROM calendar WHERE activity_id = ? AND date = ?`)
					if err == nil {
						count, err := lib.GetInt(ts, "SELECT COUNT(*) FROM activity a INNER JOIN calendar c ON c.activity_id = a.activity_id WHERE a.activity_id = ?", activityId)
						if err == nil {
							if count == 0 {
								_, err := ts.Exec("DELETE FROM activity WHERE activity_id = ?", activityId)
								if err == nil {
									ts.Commit()
									return nil
								} else {
									ts.Rollback()
									return err
								}
							} else {
								ts.Commit()
								return nil
							}
						} else {
							ts.Rollback()
							return err
						}
					} else {
						ts.Rollback()
						return err
					}
				} else {
					return err
				}
			} else {
				return fmt.Errorf("the event has some members as participant(%v)", count)
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func CalendarList(begin, end, active string) ([]CalendarReportType, error) {

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

	sql := `SELECT 
	c.calendar_id,
	c.activity_id,
	a.name,
	a.repetitive,
	a.active,
	c.date,
	COUNT(distinct m.monetary_id) AS participants,
	SUM(IF(m.payment IS NOT NULL AND m.type = 'INCOME',m.payment,0)) AS total_income,
	SUM(IF(m.payment IS NOT NULL AND m.type = 'EXPENSE',m.payment,0)) AS total_expense
	FROM calendar c 
	INNER JOIN activity a ON a.activity_id = c.activity_id AND a.active = %s
	LEFT JOIN monetary m ON m.calendar_id = c.calendar_id 
		WHERE c.date BETWEEN ? AND ?
			GROUP BY c.calendar_id, c.activity_id, a.name, a.repetitive, a.active, c.date 
			ORDER BY c.date DESC LIMIT 1000`

	sql = fmt.Sprintf(sql, getActive(active))

	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		list, err := lib.GenericQuery[CalendarReportType](db, sql, begin, end)
		if err == nil {
			return list, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func CalendarAddRemoveMemeber(calendar_id, user_id int64) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		count, err := lib.GetInt(db, `SELECT COUNT(*) FROM monetary m WHERE m.calendar_id = ? AND m.user_id = ? AND type = 'INCOME'`, calendar_id, user_id)
		if err == nil {
			if count == 0 {
				fee, err := lib.GetFloat(db, `SELECT 
CASE 
	WHEN u.payment_model = 'STUDENT' THEN a.student_fee
	WHEN u.payment_model = 'WORKER' THEN a.worker_fee 
	ELSE a.single_fee 
END AS fee
FROM calendar c INNER JOIN user u ON u.user_id = ?
INNER JOIN activity a ON a.activity_id = c.activity_id 
WHERE c.calendar_id = ? LIMIT 1`, user_id, calendar_id)
				if err == nil {
					_, err = db.Exec(`INSERT INTO monetary (calendar_id, user_id, type, fee) VALUES (?, ?, 'INCOME', ?)`, calendar_id, user_id, fee)
					return err
				} else {
					return err
				}
			} else {
				_, err := db.Exec(`DELETE FROM monetary WHERE calendar_id = ? AND user_id = ? AND type = 'INCOME'`)
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func CalendarMembers(calendar_id int64) ([]UserDetailType, error) {
	sql := `SELECT 
				u.user_id,				
				u.name,
				u.payment_model,
				u.active,
				m.monetary_id,
				m.fee,
				m.payment,
				m.method
					FROM user u 
					LEFT JOIN monetary m ON m.user_id = u.user_id AND m.calendar_id = ? AND m.type = 'INCOME'
						WHERE m.monetary_id IS NOT NULL OR u.active = 'YES'`

	db, err := lib.DbConnect()
	if err == nil {
		list, err := lib.GenericQuery[UserDetailType](db, sql, calendar_id)
		if err == nil {
			return list, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

}
