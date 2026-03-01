package data

import (
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type GetUnsentEmailListType struct {
	UserId   int64          `db:"user_id"`
	Guid     string         `db:"guid"`
	Email    string         `db:"email"`
	Params   sql.NullString `db:"params"`
	Kind     string         `db:"kind"`
	Password sql.NullString `db:"password"`
}

func (use *GetUnsentEmailListType) GetParams() any {
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

func (use *GetUnsentEmailListType) GetTempId() int {
	switch use.Kind {
	case "ACTIVATE":
		return 1
	case "PASSWORD":
		return 6
	default:
		return 0
	}
}

func GetEmailList(db *sqlx.DB, status string, limit int) ([]GetUnsentEmailListType, error) {
	var list []GetUnsentEmailListType
	err := db.Select(&list, "SELECT user_id, guid, email, params, kind, password FROM userguid WHERE expire_time > NOW() AND status = ? ORDER BY expire_time DESC LIMIT ?", status, limit)
	if err == nil {
		return list, nil
	} else {
		return nil, err
	}
}

func SetEmail(tx *sqlx.Tx, guid, status string) error {
	_, err := tx.Exec("UPDATE userguid SET status = ? WHERE guid = ?", status, guid)
	return err
}
