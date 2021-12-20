package brackets

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	_ "golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserId int

type User struct {
	AccountId int64  `json:"accountId"`
	Id        int64  `json:"id"`
	Subject   string `json:"subject"`
	Email     string `json:"email"`
	Provider  string `json:"provider"`
	GivenName string `json:"givenName"`
	LastLogin time.Time
}

func getProfileString(key string, profile map[string]interface{}) string {

	v := profile[key]
	if v == nil {
		return ""
	}
	return v.(string)

}

func AddUserFromProfile(tx dbr.SessionRunner, accountId int64, profile map[string]interface{}) (*User, error) {

	var id int

	err := tx.InsertInto("users").
		Pair("account_id", accountId).
		Pair("subject", getProfileString("sub", profile)).
		Pair("email", getProfileString("email", profile)).
		Pair("provider", getProfileString("iss", profile)).
		Pair("given_name", getProfileString("given_name", profile)).
		Pair("picture_url", getProfileString("picture", profile)).
		Returning("id").Load(&id)

	if err != nil {
		log.Fatalf("Insert User failed: %v", err)
		return nil, err
	}

	var user User
	err = tx.Select("*").From("users").Where("id = ?", id).LoadOne(&user)

	if err != nil {
		log.Fatalf("Select User failed: %v", err)
		return nil, err
	}
	return &user, nil

}

/*
func AddProvidedUser(
	db dbr.SessionRunner, accountId int, email string,
	provider string, subject string, givenName string
	) (*User, error) {

	var id int64
	err := db.InsertInto("users").
		Pair("account_id", accountId).
		Pair("subject", subject).
		Pair("email", email).
		Pair("provider", provider).
		Pair("given_name", givenName).
		Pair("picture", pictureUrl).
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

*/

func AddUser(tx dbr.SessionRunner, accountId int64, subject string, email string, name string) (*User, error) {

	fmt.Printf("AddUser %s\n", subject)
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

	var id int64
	err := tx.InsertInto("users").
		Pair("account_id", accountId).
		Pair("subject", subject).
		Pair("email", email).
		Pair("provider", "").
		Pair("given_name", "").
		Returning("id").Load(&id)

	if err != nil {
		fmt.Printf("AddUser InsertIntro error: %v\n", err)
	}

	fmt.Printf("flag1\n")

	var user User
	err = tx.Select("*").From("users").Where("id = ?", id).LoadOne(&user)

	if err != nil {
		fmt.Printf("err11: %v\n", err)
		log.Fatal(err)
	}

	//_, err := DBConn.Exec("insert into users(id, name, password_digest, password_salt, lastLogin) values($1, $2, $3)",
	// 	 user, hash, nil)
	//	err = DB.Insert(&User{user, string(hash), "", nil})
	return &user, err
}

func LoadUserByEmail(tx dbr.SessionRunner, email string) (*User, error) {

	var user User
	err := tx.Select("*").From("users").Where("email = ?", email).LoadOne(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoadUserById(tx dbr.SessionRunner, uid int64) (*User, error) {
	var user User
	err := tx.Select("*").From("users").Where("id = ?", uid).LoadOne(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoadUserBySubject(tx dbr.SessionRunner, subject string) (*User, error) {

	var user User
	err := tx.Select("*").From("users").Where("subject = ?", subject).LoadOne(&user)

	if err != nil {
		fmt.Printf("LoadUserBySubject error: %v\n", err)
		return nil, err
	}
	return &user, nil
}

/*
func Auth(db dbr.SessionRunner, email string, password string) (*User, error) {

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
