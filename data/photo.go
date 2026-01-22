package data

import (
	lib "GoDojoGo/lib"
	"context"
	"database/sql"
	"errors"
)

func SavePhotoData(ctx context.Context, userId int64, bArray []byte) error {
	db, err := lib.DbConnect()
	if err == nil {
		cmd := `INSERT INTO userphoto (user_id, image) VALUES (?, ?) ON DUPLICATE KEY UPDATE image = VALUES(image)`
		_, err := db.ExecContext(ctx, cmd, userId, bArray)
		return err
	} else {
		return err
	}
}

func GetUserPhoto(ctx context.Context, userId int64) (bArray []byte, err error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		cmd := `SELECT image FROM userphoto WHERE user_id = ? LIMIT 1`
		err = db.GetContext(ctx, &bArray, cmd, userId)
		if err == nil {
			return bArray, nil
		} else {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.New("photo not found")
			}
			return nil, err
		}
	} else {
		return nil, err
	}
}
