package data

import (
	"GoDojoGo/Cron/deff"
	"GoDojoGo/Cron/lib"
	"database/sql"
	"encoding/json"
)

type UnsendEmailType struct {
	EmailPoolId int64          `db:"emailpool_id" json:"emailpool_id"`
	Email       string         `db:"email" json:"email"`
	Params      sql.NullString `db:"params" json:"params"`
	Kind        string         `db:"kind" json:"kind"`
}

func (use *UnsendEmailType) GetParams() any {
	if use.Params.Valid {
		var pl map[string]any
		err := json.Unmarshal([]byte(use.Params.String), &pl)
		if err == nil {
			return pl
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func GetUnsentEmails() ([]UnsendEmailType, error) {
	db, err := lib.DbConnect(deff.Settings.MYSQL_DSN)
	if err == nil {
		defer db.Close()
		result, err := lib.GenericQuery[UnsendEmailType](db, "SELECT emailpool_id, email, params, kind FROM emailpool WHERE sent = 'NO' ORDER BY create_at ASC LIMIT 10")
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func SetEmailAsSent(emailpool_id int64) error {
	db, err := lib.DbConnect(deff.Settings.MYSQL_DSN)
	if err == nil {
		defer db.Close()
		_, err := db.Exec("UPDATE emailpool SET sent = 'YES' WHERE emailpool_id = ?", emailpool_id)
		return err
	} else {
		return err
	}
}
