package brackets

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"io"
	"io/ioutil"
	"log"
)

type Image struct {
	format      string          `json:"format"`
	imageData   []byte
}



func CreateImage(db dbr.SessionRunner, format string, reader io.Reader) (int, error) {

	bytes, err := ioutil.ReadAll(reader)

	if err != nil {
		fmt.Printf("err = %v\n", err)
		return 0, err
	}
	var id int
	err = db.InsertInto("images").
		Pair("format", format).
		Pair("imageData", bytes).
		Returning("id").Load(&id)

	if err != nil {
		log.Fatalf("Insert User failed: %v", err)
		return 0, err
	}

    return id, nil

}
