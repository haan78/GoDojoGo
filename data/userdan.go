package data

import lib "GoDojoGo/lib"

type UserDanSetType struct {
	UserId   int64  `db:"user_id" json:"user_id"`
	Dan      string `db:"dan" json:"dan"`
	ExamDate string `db:"exam_date" json:"exam_date"`
	Juri     string `db:"juri" json:"juri"`
	Location string `db:"location" json:"location"`
}

type UserDanDelType struct {
	UserId int64  `db:"user_id" json:"user_id"`
	Dan    string `db:"dan" json:"dan"`
}

func SetUserDan(uds *UserDanSetType) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `INSERT INTO userdan (user_id, dan, exam_date, juri, location) 
					VALUES (:user_id, :dan, :exam_date, :juri, :location)
						ON DUPLICATE KEY UPDATE exam_date = :exam_date, juri = :juri, location = :location`
		_, err := db.NamedExec(cmd, uds)
		return err
	} else {
		return err
	}
}

func DelUserDan(udd *UserDanDelType) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `DELETE FROM userdan WHERE user_id = :user_id AND dan = :dan`
		_, err := db.NamedExec(cmd, udd)
		return err
	} else {
		return err
	}
}

func GetUserDan(userId int64) ([]UserDanSetType, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `SELECT user_id, dan, exam_date, juri, location FROM userdan WHERE user_id = ?`
		result, err := lib.GenericQuery[UserDanSetType](db, cmd, userId)
		if err == nil {
			return result, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
