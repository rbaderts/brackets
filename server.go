package brackets

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gocraft/dbr/v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"strings"

	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/rbaderts/brackets/auth"
	//"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"time"

)

var (
	bracketLiveTemplate *template.Template
	homeTemplate *template.Template
	controlTemplate *template.Template
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

func GetUserId(ctx context.Context) *int64 {
	userId, ok := ctx.Value("uid").(int64)
	fmt.Printf("userId = %d\n", userId)
	if !ok {
		// Log this issue
		return nil
	}
	return &userId
}

func Server() {

	env := &Env{
		DB: DB,
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
		// We might also have a custom log.Logger, our
		// template instance, and a config struct as fields
		// in our Env struct.
	}

	assetBox := packr.NewBox("./web")

	loadTemplates()

	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	//	r.Use(jwtauth.Verifier(tokenAuth))

	FileServer(r, "/static", assetBox)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthenticationRequired)
		r.Get("/home", Handler{env, HomeRenderHandler}.ServeHTTP)
		r.Route("/tournaments/{tournamentID}", func(r chi.Router) {
			r.Get("/bracketlive", Handler{env, BracketLiveRenderHandler}.ServeHTTP)
			r.Get("/control", Handler{env, ControlRenderHandler}.ServeHTTP)
		})
		r.Route("/api", func(r chi.Router) {

			r.Post("/tournaments", Handler{env, CreateTournamentHandler}.ServeHTTP)

			r.Get("/tournaments", Handler{env, GetTournamentsHandler}.ServeHTTP)
			r.Route("/tournaments/{tournamentID}", func(r chi.Router) {
				r.Use(TournamentCtx)
				r.Get("/", Handler{env, GetTournamentHandler}.ServeHTTP)
				r.Put("/generate", Handler{env, GenerateBracketHandler}.ServeHTTP)
				r.Put("/randomize", Handler{env, RandomizePlayersHandler}.ServeHTTP)
				r.Route("/users", func(r chi.Router) {
					r.Get("/", Handler{env, GetPlayerListHandler}.ServeHTTP)
					r.Delete("/", Handler{env, DeletePlayersHandler}.ServeHTTP)
					r.Post("/", Handler{env, PostUserHandler}.ServeHTTP)
					r.Route("/{playerID}", func(r chi.Router) {
						r.Post("/paid", Handler{env, PlayerPaidHandler}.ServeHTTP)
					})
				})
				r.Route("/games/{gameID}/winner", func(r chi.Router) {
					r.Use(TournamentCtx)
					r.Delete("/", Handler{env, DeleteGameResultHandler}.ServeHTTP)
      				r.Route("/{slot}", func(r chi.Router) {
						r.Post("/", Handler{env, PostGameResultHandler}.ServeHTTP)
					})
				})

			})
		})
	})

	r.Get("/callback", AuthCallbackHandler)
	r.Get("/login", LoginHandler)
	r.Get("/logout", LogoutHandler)
	/*
				r.Route("/api", func(r chi.Router) {
					r.Route("/tournaments", func(r chi.Router) {
						r.Post("/", CreateTournamentHandler)
						r.Get("/", ListTournamentHandler)
						r.Route("/{tournamentID}", func(r chi.Router) {
	*/
	fmt.Printf("launching server on 3000\n")

	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Printf("ListenAndServe error = %v\n", err)

	}


}

func HomeRenderHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	session, err := auth.AuthStore.Get(r, "auth-session")
	if err != nil {
		return StatusError{500, err}
	}

	var name string
	var ok bool
	var val interface{}
	if val, ok = session.Values["given_name"]; !ok {
		return StatusError{http.StatusSeeOther, err}
	}

	name = val.(string)

	data := struct {
		UserName string
		TournamentID string
	}{
		name,
		"",
	}

	if err := homeTemplate.Execute(w, data); err != nil {
		return StatusError{500, err}
	}
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

func PlayerPaidHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tournament := r.Context().Value("tournament").(*Tournament)
	_ = tournament

	playerIdStr := chi.URLParam(r, "playerID")
	playerId, err := strconv.Atoi(playerIdStr)

	_ = playerId
	_ = err

	return nil

}

func DeletePlayersHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("DeletePlayersHandler\n")

	userId := r.Context().Value("uid").(int)
	tournament := r.Context().Value("tournament").(*Tournament)

	players, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}
	names := make([]Username, 0)

	if err := json.Unmarshal(players, &names); err != nil {
		return StatusError{500, err}
	}

	justnames := make([]string, 0)
	for _, p := range names {
		justnames = append(justnames, p.Name)
	}

	tournament.RemovePlayers(justnames)

	if err = tournament.Store(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
	}

	playerList := tournament.GetPlayers()

	fmt.Printf("usernames = %s\n", justnames)
	render.JSON(w, r, playerList.Players)

    return nil
}

func GetPlayerListHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	fmt.Printf("GetPlayerListHandler\n")

	tournament := r.Context().Value("tournament").(*Tournament)

	playerList := tournament.GetPlayers()

	render.JSON(w, r, playerList.Players)

    return nil
}

func PostUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	userId := r.Context().Value("uid").(int)

	fmt.Printf("PostUserHandler\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	usernames, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusError{500, err}
	}
	names := make([]Username, 0)

	if err = json.Unmarshal(usernames, &names); err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("usernames = %s\n", names)

	for _, u := range names {

		nm := strings.TrimSpace(u.Name)
		if nm != "" {
			tournament.AddPlayer(nm)
		}
	}

	playerList := tournament.GetPlayers()

	if err := tournament.Store(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
	}
	render.JSON(w, r, playerList.Players)

    return nil
}

func DeleteGameResultHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	userId := r.Context().Value("uid").(int)

	fmt.Printf("DeleteGameResult\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	gameIdStr := chi.URLParam(r, "gameID")
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		return StatusError{500, err}
	}

	if err = tournament.RemoveResult(gameId); err != nil {
		return StatusError{500, err}
	}

	if err = tournament.Store(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
	}

	var t *Tournament
	t, err = LoadTournament(DBSession, int64(tournament.Id))
	if err != nil {
		return StatusError{500, err}
	}

	render.JSON(w, r, t)
	return nil
}

func PostGameResultHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	userId := r.Context().Value("uid").(int)
	fmt.Printf("PostGameResult\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	fmt.Printf("tourn = %v\n", tournament)
	gameIdStr := chi.URLParam(r, "gameID")
	slotStr := chi.URLParam(r, "slot")

	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		return StatusError{500, err}
	}
	_ = gameId

	slot := 0

	if slot, err = strconv.Atoi(slotStr); err != nil {
		return StatusError{500, err}
	}
	_ = slot

//	result := GameResult{gameId, slot, time.Now().Unix()}

    node := tournament.GetNode(gameId)

    if err := tournament.AddResult(node, slot); err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("%s\n", tournament.Root.String())
	if err := tournament.Store(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
	}

	var t *Tournament
	t, err = LoadTournament(DBSession, int64(tournament.Id))
	if err != nil {
		return StatusError{500, err}
	}

	fmt.Printf("%s\n", t.Root.String())
	render.JSON(w, r, t)
	return nil
}


func GetTournamentsHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	userId := r.Context().Value("uid").(int)

	fmt.Printf("userId = %d\n", userId)

	var tRecords []*TournamentRecord
	var err error
	if tRecords, err = ListTournaments(DBSession, int64(userId)); err != nil {
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
	userId := r.Context().Value("uid").(int)

	t := NewTournament2()

	if err := t.Store(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
	}

	fmt.Printf("new Tourn ID = %d\n", t.Id)

	render.JSON(w, r, t)

    return nil
}

func RandomizePlayersHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value("uid").(int)
	fmt.Printf("GenerateBracketHandler\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	var err error
	if err = tournament.Draw(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
	}
	render.JSON(w, r, tournament)

	return nil
}

func GenerateBracketHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	userId := r.Context().Value("uid").(int)
	fmt.Printf("GenerateBracketHandler\n")
	tournament := r.Context().Value("tournament").(*Tournament)

	var err error
	if tournament, err = tournament.BuildBrackets(DBSession, int64(userId)); err != nil {
		return StatusError{500, err}
	}
	render.JSON(w, r, tournament)

	return nil

}

func TournamentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t *Tournament
		var err error

		if tId := chi.URLParam(r, "tournamentID"); tId != "" {
			var id int
			id, err = strconv.Atoi(tId)
			if err != nil {
				render.Render(w, r, ErrNotFound)
				return
			}
			t, err = LoadTournament(DBSession, int64(id))
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
    TournamentID    string
}


func ControlRenderHandler(env *Env, w http.ResponseWriter, r *http.Request) error {

	//t := r.Context().Value("tournament").(*Tournament)
    tId := chi.URLParam(r, "tournamentID")

	session, err := auth.AuthStore.Get(r, "auth-session")
	if err != nil {
		return StatusError{500, err}
	}

	var ok bool
	var val interface{}
	if val, ok = session.Values["given_name"]; !ok {
		return StatusError{http.StatusSeeOther, err}
	}
	name := val.(string)

	data := struct {
		UserName string
		TournamentID string
	}{
		name,
		tId,
	}



	if err := controlTemplate.Execute(w, data); err != nil {
		fmt.Printf("err = %v", err)
		return StatusError{500, err}
	}

     return nil
}

func BracketLiveRenderHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	tId := chi.URLParam(r, "tournamentID")

	session, err := auth.AuthStore.Get(r, "auth-session")
	if err != nil {
		fmt.Printf("err = %d\n", err)
		return StatusError{500, err}
	}

	var ok bool
	var val interface{}
	if val, ok = session.Values["given_name"]; !ok {
		return StatusError{http.StatusSeeOther, err}
	}
	name := val.(string)

	data := struct {
		UserName string
		TournamentID string
	}{
		name,
		tId,
	}

	if err := bracketLiveTemplate.Execute(w, data); err != nil {
		fmt.Printf("err = %v", err)
		return StatusError{500, err}
	}
	return nil
}


func FormatAsDate(t time.Time) string {
	return t.Format(TIME_FORMAT)
}


func loadTemplates() {


	controlTemplate = template.Must(template.New("control").ParseFiles(
		"../../web/templates/control.tmpl",
		"../../web/templates/control_header.tmpl",
		"../../web/templates/control_footer.tmpl",
     	"../../web/templates/nav.tmpl")).Funcs(fmap)

	homeTemplate = template.Must(template.New("home").ParseFiles(
		"../../web/templates/home.tmpl",
		"../../web/templates/home_header.tmpl",
		"../../web/templates/home_footer.tmpl",
		"../../web/templates/nav.tmpl")).Funcs(fmap)

	bracketLiveTemplate = template.Must(template.New("bracket_live").ParseFiles(
		"../../web/templates/bracket_live.tmpl",
		"../../web/templates/bracket_live_header.tmpl",
		"../../web/templates/bracket_live_footer.tmpl",
		"../../web/templates/sidebar.tmpl",
		"../../web/templates/nav.tmpl")).Funcs(fmap)

}

type Username struct {
	Name string `json:"name"`
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
    DB *dbr.Connection
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

