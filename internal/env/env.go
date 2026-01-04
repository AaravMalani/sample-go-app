package env

import (
	"os"

	"github.com/joho/godotenv"
)

var DBHost string
var DBUser string
var DBPassword string
var DBName string
var DBPort string
var DBSsl string
var DBTz string
var JwtSecret []byte

func Setup() {
	godotenv.Load()
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")
	DBSsl = os.Getenv("DB_SSL")
	DBTz = os.Getenv("DB_TZ")
	JwtSecret = ([]byte)(os.Getenv("JWT_SECRET"))
}
