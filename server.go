package brackets

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gocraft/dbr/v2"
	"github.com/gorilla/websocket"
	param "github.com/oceanicdev/chi-param"

	"io/ioutil"
	"log"
	"os"
	"strconv"

	"strings"

	"fmt"
	"html/template"
	"net/http"
	"time"
)

var (
	//	bracketLiveTemplate *template.Template
	//homeTemplate        *template.Template
	//	tournamentsTemplate *template.Template
	//	mainTemplate        *template.Template
	//controlTemplate *template.Template
	//	playersTemplate      *template.Template
	upgrader = websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
)

const (
	TIME_FORMAT = "02/06/2002 3:04PM"
)

func init() {
	fmt.Printf("Server.init\n")
}

var fmap = template.FuncMap{
	"FormatAsDate": FormatAsDate,
	"eq": func(a, b interface{}) bool {
		return a == b
	},
}

//var Assets http.FileSystem = http.Dir("assets")

func GetTournamentId(r *http.Request) (int64, error) {
	if tId := chi.URLParam(r, "tournamentID"); tId != "" {
		id, err := strconv.ParseInt(tId, 10, 64)
		if err != nil {
			return 0, StatusError{500, err}
		}
		return id, nil
	}
	return 0, errors.New("Unable to get tournamentID from request")
}
func GetUserId(ctx context.Context) *int64 {
	userId, ok := ctx.Value(userIdContextKey).(int64)
	fmt.Printf("userId = %d\n", userId)
	if !ok {
		// Log this issue
		return nil
	}
	return &userId
}

func Server() {

	env := &Env{
		DB:   DB,
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
		// We might also have a custom log.Logger, our
		// template instance, and a config struct as fields
		// in our Env struct.
	}

	//	assetBox := packr.NewBox("./web")
	//	htmlBox := packr.NewBox("./dist")

	r := chi.NewRouter()

	corsHandler := cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://localhost:3000", "http://web:3001", "http://localhost:3001", "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:  []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Contenttype"},
		AllowOriginFunc: func(r *http.Request, origin string) bool { return true },

		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            true,
	})
	r.Use(corsHandler)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)

	//	r.Use(jwtauth.Verifier(tokenAuth))

	//	FileServer(r, "/", htmlBox)
	//	FileServer(r, "/static", assetBox)
	FileServer(r, "/static", http.Dir("./web"))
	FileServer(r, "/resources/", http.Dir("./src"))
	//FileServer(r, "/brackets/", http.Dir("./build"))

	//OptionHandler := func(w http.ResponseWriter, _ *http.Request) {
	//	fmt.Printf("OptionHandler\n")
	//return
	////}

	OptionHandler := func(w http.ResponseWriter, _ *http.Request) {
		fmt.Printf("OptionHandler\n")
		return
	}

	r.Group(func(r chi.Router) {
		r.Use(DBConnection)
		//		r.Get("/callback", Handler{env, AuthCallbackHandler}.ServeHTTP)
		//		r.Get("/login", Handler{env, LoginHandler}.ServeHTTP)
		//		r.Get("/devlogin", Handler{env, DevLogin}.ServeHTTP)
		///		r.Get("/logout", Handler{env, LogoutHandler}.ServeHTTP)
		r.Get("/privacy", Handler{env, PrivacyHandler}.ServeHTTP)
	})

	//[]byte("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAosL3+U8ZfT7xx9hqO0qRjxA8xLneYS5aJ8qX87yECxhHzJ/fSJiKOQ7eCFwuoJ6lP+xOL1FHphHl6nhe4MHfpsVEa26oACGB+aDz+uguZdUG8NlXilKMfvCkWABht3d2OnyaWRie6Ngmwc3mFRdq1+I9/F3OjwS2M1PpG+WN5xGRne8fWIMgNfvqF8svo4UpcIKy3sBFZrzEe24JH7s+BJY3BPmIoBJz9cacnUNjhp2jneIvogIy0qHmUK7FMDIQeOL9EUdO/a//WFEpz1mLf0cWAj9zbLffx/tzM3y1rcMB2CIi6I+NE9ng5ixnyJdT3z7ikS75xTq2zHZEaejk5QIDAQAB"),
	//	keycloakPublicKey := []byte("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApcyYwEP63rEHGwjJEI/+SxZDIPlksyaPvHCJ8zm9vJqmTXhhbM8+g9V1ZxwAkWXXiEPYCEYdAQy/02WFW5dE97PU/jNJmFPAP2YpKlD92v6kH/ixk6hDMp5xt4OYNgbUudsr9vlsNhV/TCIHOaRJtBBw4pFocFFZakd1/7lh8KMCHF73MAZcWLA0E4l5yA12XY6by+SV/+tMo5gC+xtlHWsnyGnFUkO2xdP0UqplvaOknK/a9amI8F4nMdQw6ooBtWQo8pXy4Vjkp9mQw5VFPrQk2Nu9D7syrpKPbfDEnhLztt5dL1qzLeagGnW0xey6d7Ag+BG2h711T6l8UagsXQIDAQAB")
	keycloakPublicKey, _ := GetKeycloakRSAKey()
	fmt.Printf("keycloakPublicKey = %s\n", keycloakPublicKey)

	tokenAuth := jwtauth.New("HS256", keycloakPublicKey, nil)

	r.Group(func(r chi.Router) {
		//r.Use(auth.VerifyAuth, DBConnection)
		//	r.Use(DBConnection)

		//r.Use(auth.VerifyAuth, DBConnection)

		//		r.Us/(VerifyUser)
		//r.Get("/", Handler{env, TournamentsRenderHandler}.ServeHTTP)
		//		r.Get("/", Handler{env, HomeRenderHandler}.ServeHTTP)
		//		r.Get("/home", Handler{env, TournamentsRenderHandler}.ServeHTTP)

		//r.Get("/players", Handler{env, PlayersRenderHandler}.ServeHTTP)

		r.Route("/callbacks", func(r chi.Router) {
			/*  This call setups a callback channel for the session associated with the request */
			//r.Get("/", Handler{env, WebserviceHandler}.ServeHTTP)

			r.Route("/{tournamentID:[0-9]+}", func(r chi.Router) {
				//r.Post("/", Handler{env, RegisterCallbackHandler}.ServeHTTP)
				//r.Delete("/", Handler{env, DeleteCallbackHandler}.ServeHTTP)
			})
			/*  This call registers interest in a topic for a channel, a registration ID is returned */

			/*  This call registers interest in a topic for a channel, a registration ID is returned */
		})

		r.Route("/tournaments/{tournamentID:}", func(r chi.Router) {
			//			r.Use(DBConnection, TournamentCtx)
			r.Use(TournamentCtx)
			//			r.Get("/", Handler{env, TournamentRenderHandler}.ServeHTTP)
			//			r.Get("/control", Handler{env, ControlRenderHandler}.ServeHTTP)

			/* This endpoint opens up a websocket to the client and
			   registers the client to receive update events for
			*/
			/*
				r.Route("/updates", func(r chi.Router) {
					r.Get("/", Handler{env, WebserviceHandler}.ServeHTTP)
				})
			*/

			//			r.Get("/watch", Handler{env, WatchRenderHandler}.ServeHTTP)

		})

		r.Route("/api", func(r chi.Router) {
			//			r.Use(transaction)
			r.Use(DBConnection, AuthTokenVerifier(tokenAuth), transaction, UserInterceptor)
			r.Post("/users/loginuser", Handler{env, LoginUserHandler}.ServeHTTP)
			r.Get("/users/{userID}", Handler{env, GetUserHandler}.ServeHTTP)

			r.Get("/players", Handler{env, GetPlayersHandler}.ServeHTTP)
			r.Post("/players", Handler{env, CreatePlayerHandler}.ServeHTTP)
			r.Get("/users/{userID}", Handler{env, GetUserHandler}.ServeHTTP)
			//			r.Get("/players/{playerID}", Handler{env, GetPlayerHandler}.ServeHTTP)
			//			r.Delete("/players/{playerID}", Handler{env, DeletePlayerHandler}.ServeHTTP)

			r.Route("/players/{playerID:[0-9]+}", func(r chi.Router) {
				r.Get("/", Handler{env, GetPlayerHandler}.ServeHTTP)
				r.Delete("/", Handler{env, DeletePlayerHandler}.ServeHTTP)
				r.Put("/", Handler{env, UpdatePlayerHandler}.ServeHTTP)

				r.Post("/image", Handler{env, PostPlayerImageHandler}.ServeHTTP)

			})

			r.Post("/tournaments", Handler{env, CreateTournamentHandler}.ServeHTTP)
			r.Options("/tournaments", OptionHandler)
			r.Get("/tournaments", Handler{env, GetTournamentsHandler}.ServeHTTP)
			r.Delete("/tournaments", Handler{env, DeleteTournamentsHandler}.ServeHTTP)

			r.Route("/tournaments/{tournamentID:[0-9]+}", func(r chi.Router) {
				r.Use(TournamentCtx, transaction)

				r.Get("/", Handler{env, GetTournamentHandler}.ServeHTTP)
				r.Put("/", Handler{env, UpdateTournamentHandler}.ServeHTTP)
				r.Put("/generate", Handler{env, GenerateBracketHandler}.ServeHTTP)
				r.Put("/randomize", Handler{env, RandomizePlayersHandler}.ServeHTTP)
				r.Put("/start", Handler{env, StartTournamentHandler}.ServeHTTP)

				r.Route("/results", func(r chi.Router) {
					//					r.Get("/", Handler{env, GetTournamentResultsHandler}.ServeHTTP)
					r.Post("/", Handler{env, PostTournamentResultHandler}.ServeHTTP)
				})

				r.Route("/games/{gameID:[0-9]+}", func(r chi.Router) {
					r.Delete("/winner", Handler{env, DeleteGameResultHandler}.ServeHTTP)
					r.Post("/winner/{slot}", Handler{env, PostGameResultHandler}.ServeHTTP)
				})
				r.Post("/participants", Handler{env, AddParticipantsHandler}.ServeHTTP)
				r.Get("/participants", Handler{env, GetParticipantsHandler}.ServeHTTP)
				r.Delete("/participants", Handler{env, RemoveParticipantsHandler}.ServeHTTP)
				//				r.Get("/players", Handler{env, GetPlayerListHandler}.ServeHTTP)
				//				r.Delete("/players", Handler{env, DeletePlayersHandler}.ServeHTTP)
				r.Route("/participants/{participantNumber:[0-9]+}", func(r chi.Router) {
					r.Delete("/", Handler{env, RemoveParticipantHandler}.ServeHTTP)
					r.Post("/paid", Handler{env, ParticipantPaidHandler}.ServeHTTP)
					r.Get("/", Handler{env, GetParticipantHandler}.ServeHTTP)
					r.Put("/", Handler{env, UpdateParticipantNameHandler}.ServeHTTP)
					//					r.Put("/", Handler{env, UpdatePlayerHandler}.ServeHTTP)
				})
			})

			r.Get("/preferences", Handler{env, GetPreferencesHandler}.ServeHTTP)
			r.Put("/preferences", Handler{env, UpdatePreferencesHandler}.ServeHTTP)
		})
	})

	/*
		r.Route("/api", func(r chi.Router) {
			r.Route("/tournaments", func(r chi.Router)
				r.Post("/", CreateTournamentHandler)
				r.Get("/", ListTournamentHandler)
				r.Route("/{tournamentID}", func(r chi.Router) {
	*/
	//	fmt.Println("routes: %v\n", docgen.JSONRoutesDoc(r))
	fmt.Printf("launching server on 3000\n")

	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Printf("ListenAndServe error = %v\n", err)

	}

}

/*
func TournamentRenderHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("rendering tournament template (main.html)\n")
	username := r.Context().Value("username").(string)
	/*data := struct {
		UserName     string
		TournamentID string
	}{
		username,
		"",
	}

	//	if err := mainTemplate.Execute(w, data); err != nil {
	//		return StatusError{500, err}
	//	}

	return nil
}
*/

/*
func TournamentsRenderHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	username := r.Context().Value("username").(string)
	data := struct {
		UserName     string
		TournamentID string
	}{
		username,
		"",
	}

	//	if err := tournamentsTemplate.Execute(w, data); err != nil {
	//		return StatusError{500, err}
	//	}
	return nil
}
*/

func HomeRenderHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	/*
		if os.Getenv("AUTH_PROVIDER_DOMAIN") != "localhost:3000" {
			name, err = validateAuth(w, r)
			if err != nil {
				return StatusError{500, err}
			}
		}
	*/

	/*
		session, err := auth.AuthStore.Get(r, "brackets-auth-session")
		if err != nil {
			return StatusError{500, err}
		}

		err = ValidateAuthSession(DBSession, session)
		if err != nil {
			fmt.Printf("session not validated!!!!!\n")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		fmt.Printf("Validated Auth Session\n")

		var name string
		var ok bool
		var val interface{}
		if val, ok = session.Values["given_name"]; !ok {
			return StatusError{http.StatusSeeOther, err}
		}

	*/

	/*
		username := r.Context().Value("username").(string)
		data := struct {
			UserName     string
			TournamentID string
		}{
			username,
			"",
		}

	*/
	//url := "/static/home.html"
	//http.Redirect(w, r, url, http.StatusFound)

	//	if err := homeTemplate.Execute(w, data); err != nil {
	//		return StatusError{500, err}
	//	}
	return nil
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

}

/*
func UpdatePlayerHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tournament := r.Context().Value("tournament").(*Tournament)
	userId := r.Context().Value("uid").(int)

	playerIdStr := chi.URLParam(r, "playerID")
	playerId, err := strconv.Atoi(playerIdStr)
	if err != nil {
		return StatusError{500, err}
	}

	updatedPlayer, err := ioutil.ReadAll(r.Body)

	var p *Player
	if err := json.Unmarshal(updatedPlayer, &p); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return StatusError{500, err}
	}
	/*

	tournament.Players[playerId].Name = p.Name
	if err = tournament.Store(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
	}


*/
/*
	render.JSON(w, r, tournament)
    return nil
}

*/

func PrivacyHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	return nil

}

func GetPreferencesHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	fmt.Printf("GetPreferencesHandler\n")
	//userId := r.Context().Value(userIdContextKey).(int)
	subject := r.Context().Value(subjectContextKey).(string)

	branch, _ := param.QueryString(r, "branch") // returns first value

	//	subset := QueryString(r, "branch")
	preferences, err := LoadPreferences(tx, subject)

	if preferences == nil || err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("branch = %s\n", branch)
	fmt.Printf("preferences = %v\n", preferences)
	result := preferences.GetPreferences(branch)
	fmt.Printf("prefs = %v\n", result)
	render.JSON(w, r, result)
	return nil
}

func UpdatePreferencesHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)

	prefsraw, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return StatusError{500, err}
	}
	var prefArray []Preference
	if err := json.Unmarshal(prefsraw, &prefArray); err != nil {
		return StatusError{500, err}
	}

	prefs, err := LoadPreferences(tx, subject)
	if err != nil {
		return StatusError{500, err}
	}
	prefs.AllPrefs = prefArray

	err = StorePreferences(tx, *prefs, subject)
	if err != nil {
		return StatusError{500, err}
	}
	return nil
}

func UpdateParticipantNameHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tournament := r.Context().Value("tournament").(*Tournament)
	subject := r.Context().Value(subjectContextKey).(string)

	participantNumberStr := chi.URLParam(r, "participantNumber")
	participantNumber, err := strconv.Atoi(participantNumberStr)
	if err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("participantNumber = %v, subject = %v\n", participantNumber, subject)

	tx := getTransaction(r.Context())
	newName, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}

	tournament.Participants[ParticipantNumber(participantNumber)].Name = string(newName)

	if err = tournament.Store(tx, subject); err != nil {
		return StatusError{500, err}
	}

	//render.JSON(w, r, participant)
	return nil
}

func GetParticipantHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tournament := r.Context().Value("tournament").(*Tournament)
	subject := r.Context().Value(subjectContextKey).(string)

	participantNumberStr := chi.URLParam(r, "participantNumber")
	participantNumber, err := strconv.Atoi(participantNumberStr)
	if err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("participantNumber = %v, subject = %v\n", participantNumber, subject)

	participant := tournament.Participants[ParticipantNumber(participantNumber)]

	render.JSON(w, r, participant)
	return nil
}

func ParticipantPaidHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tx := getTransaction(r.Context())
	tournament := r.Context().Value("tournament").(*Tournament)
	subject := r.Context().Value(subjectContextKey).(string)

	participantNumberStr := chi.URLParam(r, "participantNumber")
	participantNumber, err := strconv.Atoi(participantNumberStr)
	if err != nil {
		return StatusError{500, err}
	}

	err = tournament.SetParticipantPaid(tx, subject, participantNumber, true)

	if err != nil {
		return StatusError{500, err}
	}
	list := tournament.GetParticipantList()
	render.JSON(w, r, list.Participants)

	return nil

}

func DeleteTournamentsHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("DeleteTournamentsHandler\n")

	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)
	//tournament := r.Context().Value("tournament").(*Tournament)

	tournaments, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}
	objects := make([]Id, 0)

	if err := json.Unmarshal(tournaments, &objects); err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("objects = %v\n", objects)
	justIds := make([]int64, 0)
	for _, p := range objects {
		fmt.Printf("id = %v\n", p)
		fmt.Printf("id = %d\n", p.Id)
		//		var v int64
		//		v, err = strconv.ParseInt(p.Id, 10, 64)
		justIds = append(justIds, p.Id)
	}

	DeleteTournaments(tx, justIds)

	var tRecords []*TournamentRecord
	if tRecords, err = ListTournaments(tx, subject, false); err != nil {
		return StatusError{500, err}
	}

	render.JSON(w, r, tRecords)

	return nil
}

func GetParticipantsHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tournament := r.Context().Value("tournament").(*Tournament)

	list := tournament.GetParticipantList()

	render.JSON(w, r, list.Participants)

	return nil
}

/*
func CreatePlayerHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	userId := r.Context().Value("uid").(int)
	accountId := r.Context().Value("accountId").(int)

	fmt.Printf("PostPlayerHandler\n")
	tournament := r.Context().Value("tournament").(*Tournament)

}
*/

func AddParticipantsHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	subject := r.Context().Value(subjectContextKey).(string)
	var accountId int64 = 1
	tournament := r.Context().Value("tournament").(*Tournament)

	participants, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("e1\n")
		return StatusError{500, err}
	}
	names := make([]Username, 0)

	fmt.Printf("TouranmentState = %v\n", tournament.State)
	if err := json.Unmarshal(participants, &names); err != nil {
		fmt.Printf("e2\n")
		return StatusError{500, err}

	}
	tx := getTransaction(r.Context())
	for _, p := range names {
		player, err := FindPlayerByName(tx, accountId, p.Name)
		if player == nil {
			player, err = CreatePlayer(tx, accountId, p.Name, "", "")
			if err != nil {
				fmt.Printf("e4\n")
				return StatusError{500, err}

			}
		}

		if tournament.FindParticipantByName(p.Name) != nil {
			fmt.Printf("Warning: %s is already registered, skipping", p.Name)
			return StatusError{500, err}
		}
		err = tournament.AddParticipant(tx, player.Id, p.Name)
		if err != nil {
			fmt.Printf("Unable to add %s to tournament\n", p.Name)
			return StatusError{500, err}
		}
	}

	if err = tournament.Store(tx, subject); err != nil {
		return StatusError{500, err}
	}

	list := tournament.GetParticipantList()

	render.JSON(w, r, list.Participants)

	return nil
}

func RemoveParticipantHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("RemoveParticipantsHandler\n")

	subject := r.Context().Value(subjectContextKey).(string)
	tournament := r.Context().Value("tournament").(*Tournament)

	participantNumberStr := chi.URLParam(r, "participantNumber")
	participantNumber, err := strconv.Atoi(participantNumberStr)
	if err != nil {
		return StatusError{500, err}
	}

	nums := make([]int, 1)
	nums[0] = participantNumber

	tx := getTransaction(r.Context())

	err = tournament.RemoveParticipants(tx, subject, nums)
	if err != nil {
		return StatusError{500, err}
	}

	list := tournament.GetParticipantList()

	render.JSON(w, r, list.Participants)
	return nil

}
func RemoveParticipantsHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("RemoveParticipantsHandler\n")

	subject := r.Context().Value(subjectContextKey).(string)

	tournament := r.Context().Value("tournament").(*Tournament)

	players, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}
	participants := make([]Participant, 0)

	if err := json.Unmarshal(players, &participants); err != nil {
		return StatusError{500, err}
	}

	nums := make([]int, 0)
	for _, p := range participants {
		nums = append(nums, int(p.Number))
	}

	tx := getTransaction(r.Context())

	err = tournament.RemoveParticipants(tx, subject, nums)
	if err != nil {
		return StatusError{500, err}
	}

	list := tournament.GetParticipantList()
	render.JSON(w, r, list.Participants)

	return nil
}

func PostPlayerImageHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return StatusError{500, err}
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return StatusError{500, err}
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return StatusError{500, err}
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")

	playerIdStr := chi.URLParam(r, "playerID")
	playerId, err := strconv.Atoi(playerIdStr)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return StatusError{500, err}
	}

	tx := getTransaction(r.Context())
	var player *Player
	player, err = LoadPlayer(tx, playerId)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return StatusError{500, err}
	}

	err = player.SetImage(tx, tempFile.Name())

	if err != nil {
		fmt.Printf("err = %v\n", err)
		return StatusError{500, err}
	}

	return nil
}

func GetPlayerHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	//userId := r.Context().Value("uid").(int)
	//accountId := r.Context().Value("accountId").(int)

	tx := getTransaction(r.Context())
	playerIdStr := chi.URLParam(r, "playerID")
	playerId, err := strconv.Atoi(playerIdStr)
	if err != nil {
		return StatusError{500, err}
	}

	player, err := LoadPlayer(tx, playerId)

	if err != nil {
		return StatusError{500, err}
	}

	render.JSON(w, r, player)

	return nil
}

func GetUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	userId := r.Context().Value(userIdContextKey).(int64)

	userIdStr := chi.URLParam(r, "userID")
	var targetUid int64 = 0
	if userIdStr == "current" {
		targetUid = userId
	} else {
		tmpId, _ := strconv.Atoi(userIdStr)
		targetUid = int64(tmpId)
	}

	tx := getTransaction(r.Context())
	fmt.Printf("targetUid = %v\n", targetUid)
	user, err := LoadUserById(tx, targetUid)
	fmt.Printf("user = %v\n", user)

	if err != nil {
		return StatusError{500, err}
	}
	render.JSON(w, r, user)

	return nil
}

func DeletePlayerHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	////	accountId := r.Context().Value("accountId").(int)
	var accountId int64 = 1
	userId := r.Context().Value(userIdContextKey).(int)
	_ = userId

	players, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}
	objects := make([]Id, 0)

	if err := json.Unmarshal(players, &objects); err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("objects = %v\n", objects)
	justIds := make([]int64, 0)
	for _, p := range objects {
		fmt.Printf("id = %v\n", p)
		fmt.Printf("id = %d\n", p.Id)
		//		var v int64
		//		v, err = strconv.ParseInt(p.Id, 10, 64)
		justIds = append(justIds, p.Id)
	}

	DeletePlayers(tx, justIds)

	var playerList []Player
	if playerList, err = GetAllPlayersForAccount(tx, accountId); err != nil {
		return StatusError{500, err}
	}

	render.JSON(w, r, playerList)

	return nil
}

func UpdatePlayerHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	playerData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}

	var player Player
	if err = json.Unmarshal(playerData, &player); err != nil {
		return StatusError{500, err}
	}

	err = player.UpdatePlayer(tx)
	if err != nil {
		return StatusError{500, err}
	}

	w.WriteHeader(http.StatusOK)
	return nil

	//return http.StatusCreated

}
func CreatePlayerHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	//userId := r.Context().Value("uid").(int)
	//accountId := r.Context().Value("accountId").(int64)
	var accountId int64 = 1

	playerData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}

	var player Player
	if err = json.Unmarshal(playerData, &player); err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("player = %v\n", player)

	pl, err := CreatePlayer(tx, accountId, player.Name, player.Email, player.Phone)

	//w.WriteHeader(http.StatusOK)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, pl)

	return nil
}

func DeleteGameResultHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)

	fmt.Printf("DeleteGameResult\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	gameIdStr := chi.URLParam(r, "gameID")
	gameId, err := strconv.Atoi(gameIdStr)
	fmt.Printf("game id = %v\n", gameId)
	if err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("DeleteGameResult 1\n")
	if err = tournament.RemoveResult(tx, subject, NodeId(gameId)); err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("DeleteGameResult 2\n")
	if err = tournament.Store(tx, subject); err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("DeleteGameResult 3\n")
	var t *Tournament
	t, err = LoadTournament(tx, int64(tournament.Id))
	if err != nil {
		return StatusError{500, err}
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, t)
	return nil
}
func PostTournamentResultHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	return nil
}

func PostGameResultHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)
	fmt.Printf("PostGameResult\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	fmt.Printf("tourn = %v\n", tournament)
	gameIdStr := chi.URLParam(r, "gameID")
	slotStr := chi.URLParam(r, "slot")

	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		return StatusError{500, err}
	}

	slot := 0

	if slot, err = strconv.Atoi(slotStr); err != nil {
		return StatusError{500, err}
	}
	_ = slot

	//	result := GameResult{gameId, slot, time.Now().Unix()}

	fmt.Printf("GameId = %d, slot = %d\n", gameId, slot)
	node := tournament.Bracket.GetNode(NodeId(gameId))

	if err := tournament.AddResult(tx, node, slot); err != nil {
		fmt.Printf("AddResult returned err: %v\n", err)
		return StatusError{500, err}
	}

	fmt.Printf("%s\n", tournament.Bracket.Root.String())
	if err := tournament.Store(tx, subject); err != nil {
		fmt.Printf("Store returned err: %v\n", err)
		return StatusError{500, err}
	}

	fmt.Printf("%s\n", tournament.Bracket.Root.String())
	render.JSON(w, r, tournament)
	return nil
}

func UpdateTournamentHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	tournament := r.Context().Value("tournament").(*Tournament)
	subject := r.Context().Value(subjectContextKey).(string)

	_ = tournament
	updatedTournament, err := ioutil.ReadAll(r.Body)

	fmt.Printf("%v\n", string(updatedTournament))
	if err != nil {
		fmt.Printf("readTournamentError: %v\n", err)
		return StatusError{500, err}
	}

	var t *Tournament
	if err := json.Unmarshal(updatedTournament, &t); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return StatusError{500, err}
	}
	//tx := NewDBSession()

	if err := t.Store(tx, subject); err != nil {
		fmt.Printf("Store error: %v\n", err)
		return StatusError{500, err}
	}
	w.WriteHeader(http.StatusCreated)
	return nil

}

func GetPlayersHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	//userId := r.Context().Value("uid").(int)
	tx := getTransaction(r.Context())
	//accountId := r.Context().Value("accountId").(int)
	var accountId int64 = 1

	var players []Player
	var err error

	if players, err = GetAllPlayersForAccount(tx, accountId); err != nil {
		return StatusError{500, err}
	}

	render.JSON(w, r, players)

	return nil

}

func GetTournamentsHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	//	token := getJWTToken(r.Context())
	//	fmt.Printf("token = %v\n", token)

	subject := r.Context().Value(subjectContextKey).(string)

	//userId := 0

	tx := getTransaction(r.Context())
	///userId := r.Context().Value("uid").(int)
	//fmt.Printf("userId = %d\n", userId)

	onlyActive := false
	activeStr := QueryString(r, "active")
	if activeStr == "yes" || activeStr == "1" {
		onlyActive = true
	}

	var tRecords []*TournamentRecord
	var err error
	if tRecords, err = ListTournamentsBySubject(tx, subject, onlyActive); err != nil {
		return StatusError{500, err}
	}

	render.JSON(w, r, tRecords)

	return nil

}

func GetTournamentHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tournament := r.Context().Value("tournament").(*Tournament)
	render.JSON(w, r, tournament)
	return nil
}

func CreateTournamentHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)
	userId := r.Context().Value(userIdContextKey).(int64)
	var accountId int64 = 1

	t := NewTournament2(userId, subject, accountId)

	if err := t.Store(tx, subject); err != nil {
		return StatusError{500, err}
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
		//		return
		//		return err
	}

	fmt.Printf("new Tourn ID = %d\n", t.Id)

	render.JSON(w, r, t)

	return nil
}

func StartTournamentHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)
	tournament := r.Context().Value("tournament").(*Tournament)

	var err error
	if err = tournament.Start(tx, subject); err != nil {
		return StatusError{500, err}

	}

	if err = tournament.Store(tx, subject); err != nil {
		return StatusError{500, err}
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func RandomizePlayersHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)
	fmt.Printf("RandomizePlayers\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	if tournament.State == UNDERWAY {
		return StatusError{409, errors.New("Tournament already underway")}
	}

	var err error
	if err = tournament.DrawParticipants(tx, subject); err != nil {
		return StatusError{500, err}
	}

	if err = tournament.Store(tx, subject); err != nil {
		return StatusError{500, err}
	}
	render.JSON(w, r, tournament)

	return nil
}

func GenerateBracketHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	tx := getTransaction(r.Context())
	subject := r.Context().Value(subjectContextKey).(string)
	fmt.Printf("GenerateBracketHandler\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	var err error
	if tournament, err = tournament.BuildBrackets(tx, subject); err != nil {
		return StatusError{500, err}
	}

	//if err = tournament.Store(tx, subject); err != nil {
	//		return StatusError{500, err}
	//	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, tournament)

	return nil

}

func EnsureUser(r *http.Request) int64 {
	subject := r.Context().Value(subjectContextKey).(string)
	fmt.Printf("Ensure User: %s\n", subject)
	tx := getTransaction(r.Context())

	user, err := LoadUserBySubject(tx, subject)
	if err != nil {
		fmt.Printf("LoginUserHandler - Error: %v\n", err)
	}
	if user == nil {
		fmt.Printf("Adding User for %s\n", subject)
		user, err = AddUser(tx, 1, subject, "", "")
		if err != nil {
			fmt.Printf("LoginUser AddUser Error: %v\n", err)
		}
	}

	return int64(user.Id)

}

func LoginUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	subject := r.Context().Value(subjectContextKey).(string)

	fmt.Printf("LoginUserHandler: %s\n", subject)
	tx := getTransaction(r.Context())

	user, err := LoadUserBySubject(tx, subject)
	if err != nil {
		fmt.Printf("LoginUserHandler - Error: %v\n", err)
	}

	if user == nil {
		fmt.Printf("Adding User for %s\n", subject)
		user, err = AddUser(tx, 1, subject, "", "")
		if err != nil {
			fmt.Printf("LoginUser AddUser Error: %v\n", err)
		}
	}

	return nil

}

func TournamentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t *Tournament
		var err error

		sess := getDBSession(r.Context())

		if tId := chi.URLParam(r, "tournamentID"); tId != "" {
			var id int64
			//			id, err = strconv.Atoi(tId)
			id, err = strconv.ParseInt(tId, 10, 64)
			if err != nil {
				render.Render(w, r, ErrNotFound)
				return
			}
			t, err = LoadTournament(sess, id)
			if err != nil {
				render.Render(w, r, ErrNotFound)
				return
			}
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "tournament", t)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type BracketData struct {
	TournamentID string
}

func FormatAsDate(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

type Username struct {
	Name string `json:"name"`
}

type Id struct {
	Id int64 `json:"id"`
}

func (this Username) String() string {
	return fmt.Sprintf(`{"name": "%s"}`, this.Name)
}

type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// A (simple) example of our application-wide configuration.
type Env struct {
	DB   *dbr.Connection
	Port string
	Host string
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func uploadFile(env *Env, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return StatusError{500, err}
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return StatusError{500, err}
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return StatusError{500, err}
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	return nil
}

// Query will get a query parameter by key.
func QueryString(r *http.Request, key string) string {
	/*
		if rctx := RouteContext(r.Context()); rctx != nil {
			return rctx.QueryString(key)
		}
	*/
	return ""
}

// QueryInt will get a query parameter by key and convert it to an int or return an error.
func QueryStringInt(r *http.Request, key string) (int, error) {
	val, err := strconv.Atoi(QueryString(r, key))
	if err != nil {
		return 0, err
	}
	return val, nil
}

func UserInterceptor(next http.Handler) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {
		uid := EnsureUser(r)
		r = r.WithContext(context.WithValue(r.Context(), userIdContextKey, uid))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(h)

}

func DBConnection(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		//		ctx := r.Context()
		fmt.Printf("Calling DB.NewSession\n")
		session := DB.NewSession(nil)
		r = r.WithContext(context.WithValue(r.Context(), sessionContextKey, session))

		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(h)
}

func AuthTokenVerifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("AuthToknVerifier\n")

			authHeader := r.Header.Get("Authorization")
			fmt.Printf("authHeader = %v\n", authHeader)
			if authHeader == "" || len(authHeader) <= 0 {
				w.WriteHeader(http.StatusInternalServerError)
				return

			}
			fmt.Printf("f2\n")
			pair := strings.Split(authHeader, " ")
			token := pair[1]
			fmt.Printf("Token = %s\n", token)

			decodedToken, err := jwt.Parse([]byte(token),
				jwt.WithDecrypt(jwa.RSA_OAEP_256, "aCzy27d42eWhrVrp36mpXUcf1LiWHwvQnUXz5E7NXZ4"))

			if err != nil {
				fmt.Printf("errorr decoding - %v\n", err)
			}
			fmt.Printf("subject = %s\n", decodedToken.Subject())

			r = r.WithContext(context.WithValue(r.Context(), subjectContextKey, decodedToken.Subject()))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
}

/*
func AuthTokenVerifier(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("AuthToknVerifier")
		authHeader := r.Header.Get("Authorization")
		ctx := r.Context()
		pair := strings.Split(authHeader, " ")
		token := pair[1]
		fmt.Printf("Token = %s\n", token)
		r = r.WithContext(context.WithValue(r.Context(), jwtauth.TokenCtxKey, token))

		next.ServeHTTP(w, r)
	}

}
*/

/*
func VerifyUser(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		sess := getDBSession(r.Context())
		userId := r.Context().Value(userIdContextKey).(int64)
		u, err := LoadUserById(sess, userId)
		fmt.Printf("user = %v\n", u)
		if err != nil {
			http.Redirect(w, r, auth.Oauth2Config.AuthCodeURL(auth.ExampleAppState), http.StatusFound)
			return
		}
		if u == nil {
			http.Redirect(w, r, auth.Oauth2Config.AuthCodeURL(auth.ExampleAppState), http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), "username",
			u.GivenName)
		fmt.Printf("user naem = %s\n", u.GivenName)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
	return http.HandlerFunc(h)

}
*/

type contextKey int

const (
	txContextKey      contextKey = iota
	sessionContextKey contextKey = iota
	subjectContextKey contextKey = iota
	userIdContextKey  contextKey = iota
)

type CustomHandler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

func (h CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.H(h.Env, w, r); err != nil {

		// handle returned error here.
		w.WriteHeader(503)
		w.Write([]byte("bad"))
	}
}

/*
func TransactionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := NewDBSession()
		ctx := context.WithValue(r.Context(), "tx", tx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/

//func transaction(handler CustomHandler) http.HandlerFunc {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

func transaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sess := getDBSession(r.Context())
		fmt.Printf("sess = %v\n", sess)
		tx, err := sess.BeginTx(r.Context(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Open transaction failed: %s \n", err.Error())
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), txContextKey, tx))
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("%v", r)
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("Transaction is being rolled back: %s \n", err.Error())
				tx.Rollback()
				return
			}
		}()

		next.ServeHTTP(w, r)
		if err != nil {
			err = tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Transaction is being rolled back: %s \n", err.Error())
			return
		}

		err = tx.Commit()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Transaction commit failed: %s \n", err.Error())
		} else {
			fmt.Println("Transaction has been committed")
		}
	})
}

func getDBSession(ctx context.Context) *dbr.Session {
	sessionValue := ctx.Value(sessionContextKey)
	if sessionValue != nil {
		tx := sessionValue.(*dbr.Session)
		return tx
	}

	return nil
}
func getTransaction(ctx context.Context) *dbr.Tx {
	txValue := ctx.Value(txContextKey)
	if txValue != nil {
		tx := txValue.(*dbr.Tx)
		return tx
	}

	return nil
}

func getJWTToken(ctx context.Context) *jwt.Token {
	tokenValue := ctx.Value(jwtauth.TokenCtxKey)
	fmt.Printf("tokenValue = %v\n", tokenValue)
	if tokenValue != nil {
		token := tokenValue.(jwt.Token)
		return &token
	}
	errorValue := ctx.Value(jwtauth.ErrorCtxKey)
	fmt.Printf("errorValue = %v\n", errorValue)
	if errorValue != nil {
		err := errorValue.(error)
		fmt.Printf("error = %v\n", err)
	}

	return nil
}

/*
func GetKeycloakRSAKey() (string, error) {

	authURL := os.Getenv("AUTH_PROVIDER_DOMAIN")
	url := authURL + "/" + ".well-known/openid-configuration"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	result := map[string]interface{}{}
	json.Unmarshal(bodyBytes, &result)

	jwksUri := result["jwks_uri"]

	req, err = http.NewRequest("GET", jwksUri, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	result = map[string]interface{}{}
	json.Unmarshal(bodyBytes, &result)

	keys := result["keys"].(map[string]interface{})
	for k, v := range result["keys"] {

	}
}
*/

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func (this JSONWebKeys) String() string {
	var b strings.Builder
	b.WriteString("Use = ")
	b.WriteString(this.Use)
	b.WriteString("\nKid = ")
	b.WriteString(this.Kid)
	return b.String()
}

func GetKeycloakRSAKey() (string, error) {

	authURL := os.Getenv("AUTH_PROVIDER_DOMAIN")
	url := authURL + "/protocol/openid-connect/certs"
	cert := ""
	resp, err := http.Get(url)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		fmt.Printf("jwks.Keys[%s] = %v\n", k, jwks.Keys[k])
		if jwks.Keys[k].Use == "enc" {
			cert = jwks.Keys[k].X5c[0]
			break
		}
	}
	if cert == "" {
		return "", errors.New("Unable to get public key")
	}
	return cert, nil
}
