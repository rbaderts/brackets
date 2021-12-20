package brackets

import (
	"errors"
	"fmt"

	"github.com/gocraft/dbr/v2"
	_ "github.com/gocraft/dbr/v2"
	"github.com/gorilla/sessions"
)

type handler struct{}

/*
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
*/

func ValidateAuthSession(tx dbr.SessionRunner, session *sessions.Session) error {

	fmt.Printf("ValidateAuthSession\n")
	uid, has := session.Values["uid"]
	if !has {
		return errors.New("No session uid")
	}

	fmt.Printf("checking uid in session %v\n", uid)
	//		userId, err := strconv.Atoi(uid.(string))
	/*
		if err != nil {
			return err
		}
	*/
	user, err := LoadUserById(tx, uid.(int64))

	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("Can't load user")
	}
	return nil

}
