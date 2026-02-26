package deff

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const VERSION = "1.0.0+1"

type SettingsType struct {
	PORT               int
	MYSQL_DSN          string
	JWT_SECRET         string
	BREVO_API_KEY      string
	EMAIL_ACTIVATE_URL string
	POOL_CHECK         int
}

var Settings SettingsType

func LoadSettings(path string) error {

	_, err := os.Stat(path)
	if err == nil {
		err := godotenv.Load(path)
		if err != nil {
			return err
		}
	}

	Settings.JWT_SECRET = os.Getenv("JWT_SECRET")
	Settings.MYSQL_DSN = os.Getenv("MYSQL_DSN")
	Settings.BREVO_API_KEY = os.Getenv("BREVO_API_KEY")

	if os.Getenv("PORT") != "" {
		num, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			Settings.PORT = num
		} else {
			return err
		}
	} else {
		Settings.PORT = 1323
	}

	if os.Getenv("POOL_CHECK") != "" {
		num, err := strconv.Atoi(os.Getenv("POOL_CHECK"))
		if err == nil {
			Settings.POOL_CHECK = num
		} else {
			return err
		}
	} else {
		Settings.POOL_CHECK = 3
	}

	return nil
}
