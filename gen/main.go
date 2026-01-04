package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

const (
	DSNFormat = "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
)

func main() {
	godotenv.Load()
	var DBHost = os.Getenv("DB_HOST")
	var DBUser = os.Getenv("DB_USER")
	var DBPassword = os.Getenv("DB_PASSWORD")
	var DBName = os.Getenv("DB_NAME")
	var DBPort = os.Getenv("DB_PORT")
	var DBSsl = os.Getenv("DB_SSL")
	var DBTz = os.Getenv("DB_TZ")
	g := gen.NewGenerator(gen.Config{
		OutPath:          "./query",
		Mode:             gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldCoverable:   true,
		FieldWithTypeTag: true,
	})
	dsn := fmt.Sprintf(
		DSNFormat, DBHost, DBUser, DBPassword, DBName, DBPort, DBSsl, DBTz,
	)
	var gormdb *gorm.DB
	var err error
	if len(os.Args) == 1 {
		gormdb, err = gorm.Open(postgres.Open(dsn))
	} else {
		gormdb, err = gorm.Open(rawsql.New(rawsql.Config{
			FilePath: os.Args[1:],
		}))
	}

	if err != nil {
		panic(err)
	}
	g.UseDB(gormdb) // reuse your gorm db
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
