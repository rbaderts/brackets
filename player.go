package brackets

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"image"
	"io"
	"log"
	"mime"
	"os"
)

type Player struct {
	Id      int64  `json:"id"         db:"id"`
	Name    string `json:"name"       db:"player_name"`
	Paid    int    `json:"paid"       db:"paid"`
	Email   string `json:"email"      db:"email"`
	Phone   string `json:"phone"      db:"phone"`
	ImageId int    `json:"imageId"    db:"image_id"`
}

type PlayerResult struct {
	Date     int64  `json:"date"`
	Opponent string `json:"opponentId"`
	Win      bool   `json:"win"`
}

func (this *Player) UpdatePlayer(db dbr.SessionRunner) error {

	UpdateStmt := db.Update("players").
		Set("player_name", this.Name).
		Set("email", this.Email).
		Set("phone", this.Phone)
	_, err := UpdateStmt.Where("id = ?", this.Id).Exec()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return nil
}

func CreatePlayer(db dbr.SessionRunner,
	accountId int64,
	name string,
	email string,
	phone string) (*Player, error) {

	var id int64
	err := db.InsertInto("players").
		Pair("account_id", accountId).
		Pair("player_name", name).
		Pair("email", email).
		Pair("phone", phone).
		Returning("id").Load(&id)

	if err != nil {
		log.Fatalf("Create Player failed: %v", err)
		return nil, err
	}

	var player Player
	err = db.Select("*").From("players").Where("id = ?", id).LoadOne(&player)

	if err != nil {
		log.Fatalf("Select User failed: %v", err)
		return nil, err
	}

	return &player, err
}

func DeletePlayers(db dbr.SessionRunner, ids []int64) error {
	for _, id := range ids {
		fmt.Printf("Deleteting player with id %v\n", id)
		err := DeletePlayer(db, int(id))
		if err != nil {
			fmt.Printf("Error deleting player with id: %v\n", id)
		}

	}
	return nil
}
func DeletePlayer(db dbr.SessionRunner, playerId int) error {
	_, err := db.DeleteFrom("players").Where(dbr.Eq("id", playerId)).Exec()

	if err != nil {
		return err
	}
	return nil
}

func LoadPlayer(db dbr.SessionRunner, playerId int) (*Player, error) {
	var player Player
	err := db.Select("*").From("players").Where("id = ?", playerId).LoadOne(&player)

	if err != nil {
		return nil, err
	}
	return &player, nil
}

func FindPlayerByName(db dbr.SessionRunner, accountId int64, name string) (*Player, error) {
	var player Player
	err := db.Select("*").From("players").
		Where("account_id = ? and player_name = ?", accountId, name).LoadOne(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func GetAllPlayersForAccount(db dbr.SessionRunner, accountId int64) ([]Player, error) {
	var players []Player
	_, err := db.Select("*").From("players").
		Where("account_id = ?", accountId).Load(&players)

	if err != nil {
		return nil, err
	}
	if players == nil {
		players = make([]Player, 0)
	}
	return players, nil
}

func (this *Player) SetImage(db dbr.SessionRunner, imagePath string) error {

	imageFile, err := os.Open(imagePath)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	mimeType := guessImageMimeTypes(imageFile)

	var imageId int
	imageId, err = CreateImage(db, mimeType, imageFile)

	if err != nil {
		fmt.Printf("err = %v\n", err)
		return err
	}

	_, err = db.Update("players").Set("image_id", imageId).
		Where("id = ?", this.Id).Exec()
	if err != nil {
		log.Fatalf("Update Tournament failed: %v", err)
		return err
	}

	return nil
}

func AddPlayerResult(db dbr.SessionRunner, playerId int, opponent int, win bool) {

	_, err := db.InsertInto("player_results").
		Pair("player_id", playerId).
		Pair("opponent_id", opponent).
		Pair("win", win).Exec()

	if err != nil {
		fmt.Printf("AddResult error: %v\n", err)
	}
}

func (this *Player) String() string {
	return fmt.Sprintf("%d:%s", this.Id, this.Name)
}

// Guess image format from gif/jpeg/png/webp
func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// Guess image mime types from gif/jpeg/png/webp
func guessImageMimeTypes(r io.Reader) string {
	format, _ := guessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}
