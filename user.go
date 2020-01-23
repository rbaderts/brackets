package brackets

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	_	"golang.org/x/crypto/bcrypt"
	"time"
	"log"
)

type UserId int

type User struct {
	Id             int
	Subject        string
	Email          string
	Provider       string
	LastLogin      time.Time
}

func AddProvidedUser(db *dbr.Session, email string, provider string, subject string) (*User, error) {

	var id int64
	err := db.InsertInto("users").
		Pair("subject", subject).
		Pair("email", email).
		Pair("provider", provider).
		Returning("id").Load(&id)

	if err != nil {
		log.Fatalf("Insert User failed: %v", err)
		return nil, err
	}

	var user User
	err = db.Select("*").From("users").Where("id = ?", id).LoadOne(&user)

	if err != nil {
		log.Fatalf("Select User failed: %v", err)
		return nil, err
	}

	return &user, err
}

func AddUser(db *dbr.Session, email string, name string) (*User, error) {

	fmt.Printf("AddUser %s\n", email)
	/*
	sqlStatement := `
INSERT INTO users (age, email, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)o
	*/

	subject := email
	result, err := db.InsertInto("users").
		Pair("subject", subject).
		Pair("email", email).
		Pair("provider", "").
		Returning("id").Exec()

	var id int64
	id, err = result.LastInsertId()

	if err != nil {
		fmt.Printf("err10: %v\n", err)
		log.Fatal(err)
	}

	var user User
	err = db.Select("*").From("users").Where("id = ?", id).LoadOne(&user)

	if err != nil {
		fmt.Printf("err11: %v\n", err)
		log.Fatal(err)
	}

	//_, err := DBConn.Exec("insert into users(id, name, password_digest, password_salt, lastLogin) values($1, $2, $3)",
	// 	 user, hash, nil)
	//	err = DB.Insert(&User{user, string(hash), "", nil})
	return &user, err
}

func LoadUserByEmail(db *dbr.Session, email string) (*User, error) {

	var user User
	err := db.Select("*").From("users").Where("email = ?", email).LoadOne(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoadUserBySubject(db *dbr.Session, subject string) (*User, error) {

	var user User
	err := db.Select("*").From("users").Where("subject = ?", subject).LoadOne(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

/*
func Auth(db *dbr.Session, email string, password string) (*User, error) {

	passwordBytes := []byte(password)

	fmt.Printf("querying for user %s\n", email)
	var user User
	_, err := db.QueryOne(&user, "SELECT * from users where email = ?", email)

	if err != nil {
		fmt.Printf("queried user error = %v\n", err)
		return nil, err
	}
	fmt.Printf("queried user: name = %s\n", user.Email)

	result := bcrypt.CompareHashAndPassword(user.PasswordDigest, passwordBytes)
	if result == nil {
		fmt.Printf("Passwords matched \n")
		return &user, nil
	}
	fmt.Printf("Passwords didn't match \n")
	return nil, result

	//	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	//		DB_USER, DB_PASSWORD, DB_NAME)
	//	db, err := sql.Open("postgres", dbinfo)
	//	if err != nil {
	//		return nil, err
	//	}
	//	defer db.Close()

}
*/
