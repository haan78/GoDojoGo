package data

import (
	lib "GoDojoGo/lib"
	"fmt"
)

type DebtType struct {
	UserId     int64        `json:"user_id" db:"user_id"`
	Name       string       `json:"name" db:"name"`
	Debt       float64      `json:"debt" db:"debt"`
	Activities int64        `json:"activities" db:"activities"`
	First      lib.JSONDate `json:"first" db:"first"`
	Last       lib.JSONDate `json:"last" db:"last"`
}

type MemberDebtType struct {
	MonetaryId int64              `json:"monetary_id" db:"monetary_id"`
	ActivityId lib.JSONInt        `json:"activity_id" db:"activity_id"`
	Date       lib.JSONDate       `json:"date" db:"date"`
	Text       lib.JSONNullString `json:"text" db:"text"`
	Fee        float64            `json:"fee" db:"fee"`
}

type IncomeType struct {
	MonetaryId int64              `json:"monetary_id" db:"monetary_id"`
	ActivityId lib.JSONInt        `json:"activity_id" db:"activity_id"`
	Date       lib.JSONDate       `json:"date" db:"date"`
	Text       lib.JSONNullString `json:"text" db:"text"`
	Payment    lib.JSONFloat      `json:"payment" db:"payment"`
}

func MemberDebts(user_id int64) ([]MemberDebtType, error) {

	sql := `SELECT 
	m.monetary_id, a.activity_id, COALESCE(a.date, DATE(m.created_at) ) as date, COALESCE(a.text, m.text) as text
	FROM monetary m
	LEFT JOIN activity a ON a.activity_id = m.activity_id 
		WHERE m.user_id  = ? AND m.type = 'INCOME' ORDER BY date ASC`

	db, err := lib.DbConnect()
	if err == nil {
		result, err := lib.GenericQuery[MemberDebtType](db, sql, user_id)
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func Debts(active string) ([]DebtType, error) {
	w1 := ""
	switch active {
	case "YES":
		w1 = "and u.active = 'YES'"
	case "NO":
		w1 = "and u.active = 'NO'"
	}

	sql := fmt.Sprintf(`SELECT
	m.user_id, u.name, SUM(m.fee) as debt, COUNT(distinct m.activity_id) AS activities, MIN(a.date) AS first, MAX(a.date) AS last
	FROM monetary m
	INNER JOIN user u on u.user_id = m.user_id
	LEFT JOIN activity a ON a.activity_id = m.activity_id
		WHERE m.payment is null AND m.type = 'INCOME' %s
			group by m.user_id, u.name order by debt DESC`, w1)

	db, err := lib.DbConnect()
	if err == nil {
		result, err := lib.GenericQuery[DebtType](db, sql)
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func AllMonetaryActions(start, end string) ([]IncomeType, error) {
	sql := `SELECT 
				m.monetary_id,
				a.activity_id,
				m.date,
				COALESCE(a.text, m.text) as text,
				m.payment,
				m.type
				FROM monetary m
					LEFT JOIN activity a ON a.activity_id = m.activity_id 
						WHERE m.payment IS NOT NULL AND (m.date BETWEEN ? AND ?)
							ORDER BY m.date ASC`

	db, err := lib.DbConnect()
	if err == nil {
		result, err := lib.GenericQuery[IncomeType](db, sql, start, end)
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
