package deff

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type SettingsType struct {
	PORT          int
	MYSQL_DSN     string
	JWT_SECRET    string
	BREVO_API_KEY string
}

var Settings SettingsType

func LoadSettings() error {
	err := godotenv.Load()
	if err == nil {
		if Settings.JWT_SECRET = os.Getenv("JWT_SECRET"); Settings.JWT_SECRET == "" {
			return errors.New("JWT_SECRET not found in .env file")
		}

		if Settings.MYSQL_DSN = os.Getenv("MYSQL_DSN"); Settings.MYSQL_DSN == "" {
			return errors.New("MYSQL_DSN not found in .env file")
		}

		num, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return err
		} else {
			Settings.PORT = num
		}

		Settings.BREVO_API_KEY = os.Getenv("BREVO_API_KEY")

		return nil

	} else {
		return err
	}
}
