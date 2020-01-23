package auth

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/coreos/go-oidc"

	//	"github.com/coreos/go-oidc"
//	"github.com/gin-contrib/sessions/cookie"
	"github.com/gorilla/sessions"

	//	oidc "github.com/coreos/go-oidc	"
//	"github.com/go-chi/jwtauth"
//	"github.com/gin-contrib/sessions/cookie"
//	"github.com/gin-gonic/gin"

	//#	"github.com/gin-contrib/sessions/cookie"
//	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)
var (
	SessionKey = "Testkey"
//	AuthStore sessions.CookieStore
	//AuthStore  = sessions.NewCookieStore([]byte(SessionKey))
	AuthStore *sessions.FilesystemStore

)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func init() {
	AuthStore = sessions.NewFilesystemStore("store", []byte(SessionKey))
	gob.Register(map[string]interface{}{})

}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, os.Getenv("BRACKETS_DOMAIN"))
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("BRACKETS_CLIENT_ID"),
		ClientSecret: os.Getenv("BRACKETS_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:3000/callback",
		Endpoint: 	  provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

// routes/callback/callback.go

//func CallbackHandler(w http.ResponseWriter, r *http.Request) {



/*
func AuthenticationRequired(auths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user needs to be signed in to access this service"})
			c.Abort()
			return
		}
		if len(auths) != 0 {
			authType := session.Get("authType")
			if authType == nil || !containsString(auths, authType.(string)) {
				c.JSON(http.StatusForbidden, gin.H{"error": "invalid request, restricted endpoint"})
				c.Abort()
				return
			}
		}
		// add session verification here, like checking if the user and authType
		// combination actually exists if necessary. Try adding caching this (redis)
		// since this middleware might be called a lot
		c.Next()
	}
}
*/

/*
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
 */

func AuthenticationRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("AuthenticationRequired\n")
		//func AuthenticationRequired(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := AuthStore.Get(r, "auth-session")
		if err != nil {
			fmt.Printf("err = %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := session.Values["profile"]; !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)

		} else {
			fmt.Printf("Authenticated\n")

			uid := session.Values["uid"]
			fmt.Printf("uid in session %d\n", uid)
			ctx := context.WithValue(r.Context(), "uid", uid)
			h.ServeHTTP(w, r.WithContext(ctx))

		}

	})
}


/*
func AuthenticationRequired(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("AuthenticationRequired\n")
		//func AuthenticationRequired(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := AuthStore.Get(r, "auth-session")
		if err != nil {
			fmt.Printf("err = %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("f3")
		if _, ok := session.Values["profile"]; !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			//v := session.Values["uid"]
			//c.Set("uid", v)

			fmt.Printf("f4")
			fmt.Printf("Authenticated\n")
			h.ServeHTTP(w, r)
		}

		fmt.Printf("f5")
	}
	return http.HandlerFunc(fn)

}

/*
func AuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("AuthenticationRequired\n")
		//func AuthenticationRequired(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := AuthStore.Get(c.Request, "auth-session")
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := session.Values["profile"]; !ok {
			http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
		} else {

			v := session.Values["uid"]
		    c.Set("uid", v)

			fmt.Printf("Authenticated\n")
			c.Next()
		}
	}
}
 */

func containsString(strings []string, checkFor string) bool {
	for _, s := range strings {
		if (s == checkFor) {
			return true
		}
	}
	return false

}