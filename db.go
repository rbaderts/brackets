package brackets

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	"log"
	"os"
)

type DBConn interface {
}

var DB *dbr.Connection

func SetupDB() *dbr.Connection {

	pg_user := os.Getenv("POSTGRES_USER")
	pg_pw := os.Getenv("POSTGRES_PASSWORD")
	pg_host := os.Getenv("BRACKETS_DB_HOST")
	pg_db := "brackets"

	fmt.Printf("POSGRES_USER = %s\n", pg_user)
	fmt.Printf("BRACKETS_DB_HOST = %s\n", pg_host)

	connStr := fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable&password=%s",
		pg_user, pg_host, pg_db, pg_pw)

	conn, err := dbr.Open("postgres", connStr, nil)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		log.Fatal(err)
	}
	DB = conn

	return DB

}

func MigrateDB() {

	fmt.Printf("MigrateDB\n")
	//pg_url := os.Getenv("BRACKETS_DB_URI")
	pg_user := os.Getenv("POSTGRES_USER")
	pg_pw := os.Getenv("POSTGRES_PASSWORD")
	pg_host := os.Getenv("BRACKETS_DB_HOST")
	pg_db := "brackets"
	fmt.Printf("POSGRES_USER = %s\n", pg_user)
	fmt.Printf("BRACKETS_DB_HOST = %s\n", pg_host)
	fmt.Printf("POSTGRES_PASSWORD = %s\n", pg_pw)

	pg_url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&password=%s",
		pg_user, pg_pw, pg_host, pg_db, pg_pw)
	fmt.Printf("pg_url = %s\n", pg_url)

	//	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//		"password=%s dbname=%s sslmode=disable",
	//		pg_host, 5432, pg_user, pg_pw, pg_db)

	db, err := dbr.Open("postgres", pg_url, nil)

	if err != nil {
		fmt.Printf("Open err = %v\n", err)
	}

	var m *migrate.Migrate
	for {
		driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			fmt.Printf("got driver err = %v\n", err)
			time.Sleep(1 * time.Second)
			fmt.Printf("Retrying migration...\n")
			continue
		}

		for {
			m, err = migrate.NewWithDatabaseInstance(
				"file://migrations",
				"postgres", driver)

			if err != nil {
				fmt.Printf("got migrate err = %v\n", err)
				time.Sleep(2 * time.Second)
				fmt.Printf("Retrying migration...\n")
				continue
			}
			break
		}
		break
	}

	// file://relative/path
	//	m, err := migrate.NewWithDatabaseInstance("file://migrations",
	//	"postgres", driver)
	//if err != nil {
	//	fmt.Printf("NewWithDatabaseIntsance err = %s\n", err)
	//	}

	err = m.Up()
	if err != nil {
		fmt.Printf("NewWithDatabaseIntsance err = %s\n", err)
	}

	//DBSession = db.NewSession(nil)
	DB = db
	//m, err = migrate.New("file://"+assetDir, pg_url)
	//_ = m

	//if err != nil {
	//		fmt.Printf("migrate err = %s\n", err)
	//	}
}

func NewDBSession() *dbr.Session {
	session := DB.NewSession(nil)
	session.Begin()
	return session
}

/*
func CommitSession(session *dbr.Session) error {
	session.Commit();
}

func RollbackSession(session *dbr.Session) {
	defer session.RollbackUnlessCommitted()
}
*/

/*

func MigrateDB(assetDir string) {


	fmt.Printf("MigrateDB\n")
	//pg_url := os.Getenv("BRACKETS_DB_URI")
	pg_user := os.Getenv("POSTGRES_USER")
	pg_pw := os.Getenv("POSTGRES_PASSWORD")
	pg_host := os.Getenv("BRACKETS_DB_HOST")
	pg_db := "brackets"

	//pg_url := fmt.Sprintf("postgres://%s:'%s'@%s:5432/%s?sslmode=disable", pg_user, pg_pw, pg_host, pg_db)
	pg_url := fmt.Sprintf("postgres://%s:'%s'@%s:5432/%s/sslmode=disable", pg_user, pg_pw, pg_host, pg_db)

	fmt.Printf("pg_url = %s\n", pg_url)
	for {
		var err error
		//db, err = dbr.Open("postgres", db_url, nil)
		conn, err := dbr.Open("postgres", pg_url, nil)
		if err != nil {
			fmt.Printf("err = %v\n", err)
		}

		DB = conn
		DBSession = conn.NewSession(nil)
		if err != nil {
			fmt.Printf("Error = %v\n", err)
			time.Sleep(time.Second * 2)
		} else {
			for {
				if err = DBSession.Ping(); err != nil {
					//continue
					fmt.Printf("DB Ping error, sleeping 2\n")
					time.Sleep(time.Second * 2)
				} else {
					fmt.Printf("DB Ping ok\n")
					break
				}
			}
			break
		}
	}

	fmt.Printf("Calling migrator.New\n")
	m, err := migrate.New("file://"+assetDir, pg_url)

	fmt.Printf("m = %v, err = %v\n", m, err)

	if err != nil {
		fmt.Errorf("err = %v", err)
	}

	fmt.Printf("Calling migrator.Up\n")
	err = m.Up()
	if err != nil {
		fmt.Errorf("Up err = %v", err)
	}

	v, _, _ := m.Version()
	fmt.Printf("database name = %v\n", v)

	SeedAccount(DBSession)
}
*/
func SeedAccount(db *dbr.Session) (*Account, error) {

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	var count int
	err = tx.Select("count(*)").From("accounts").LoadOne(&count)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return CreateAccount(tx, "Default")
	}

	return nil, nil
}
