package main

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rbaderts/brackets"

	"github.com/gobuffalo/packr"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	dir := setupMigrationAssets()
	brackets.MigrateDB(dir)
	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
	//	brackets.Initialize()

	/*
	b, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("json marshall error %v\n", err)
	}
	fmt.Printf("winners = %v\n", string(b))

	var s brackets.Tournament
	err = json.Unmarshal(b, &s)

	var c []byte
	c, err = json.Marshal(t)
	fmt.Printf("reconstituted = %v\n", string(c))
	 */
	brackets.Server()

}


func setupMigrationAssets() string {
	tmpDir, err := ioutil.TempDir("./", "migrations_tmp")
	//err := os.Mkdir("./migrations_tmp", 0755)
	if err != nil {
		log.Fatal(err)
	}

	box := packr.NewBox("../../migrations")

	for _, m := range box.List() {
		fmt.Printf("box item %s\n", m)
		bytes, err := box.Find(m)
		if err != nil {
			fmt.Printf("err = %v", err)
			continue
		}
		fmt.Printf("bytes = %v", string(bytes))
		tmpfn := filepath.Join(tmpDir,  m)
		fmt.Printf("tmpFn = %v\n", tmpfn)
		err = ioutil.WriteFile(tmpfn, bytes, 0644)
		if err != nil {
			fmt.Errorf("err = %v", err)
			continue
		}
	}
	return tmpDir
}