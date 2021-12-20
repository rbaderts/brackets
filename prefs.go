package brackets

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocraft/dbr/v2"
)

type PreferenceRecord struct {
	Id           int64
	Subject      string
	TournamentId int64

	Data Preferences
}

/*
   names are dot separated hierarchical paths
*/
type Preference struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (this *Preference) String() string {
	return `"` + this.Name + `" : "` + this.Value + `"}`

}

/*
  Represents a node in the preference name hierarchy
    Leaf nodes, where children is nil, generally represent a pref value
    Ex:
         a.b.c.d = 4                   a
         a.b.e.f = 2    ======      b     c (1)
         a.c = 1                 c     e
                                   d(4)  f(2)

*/

type NameNode struct {
	name       string
	children   map[string]*NameNode
	preference *Preference
}

func newRootNameNode() *NameNode {
	p := new(NameNode)
	p.name = ""
	p.children = make(map[string]*NameNode, 0)
	return p
}

func (this *NameNode) newSubTree(name string) *NameNode {
	p := new(NameNode)
	p.name = name
	p.children = make(map[string]*NameNode, 0)
	this.children[name] = p
	return p
}

func (this *NameNode) newPreferenceNode(name string, preference *Preference) *NameNode {
	p := new(NameNode)
	p.preference = preference
	p.name = name
	p.children = nil
	this.children[name] = p
	return p
}

func (this *NameNode) add(path string, pref Preference) *NameNode {

	parts := strings.Split(path, ".")
	firstName := parts[0]

	if len(parts) == 1 {
		child, has := this.children[firstName]
		if !has {
			child = this.newPreferenceNode(firstName, &pref)
			this.children[firstName] = child
		}
		return child
	} else {
		child, has := this.children[firstName]
		if !has {
			child = this.newSubTree(firstName)
			this.children[firstName] = child
		}
		subpath := strings.Replace(path, firstName+".", "", 1)
		return child.add(subpath, pref)
	}

	return nil

}

/*
 */

type Preferences struct {
	Root *NameNode
	//AllPrefs map[string]*Preference `json:"preferences"`
	AllPrefs []Preference `json:"preferences"`
}

func (this *Preferences) MarshalJSON() ([]byte, error) {

	var bytes []byte
	var err error

	//fmt.Printf("pref count = %d\n", len(this.AllPrefs))
	if bytes, err = json.Marshal(&this.AllPrefs); err != nil {
		return nil, err
	}

	return bytes, nil
}

func (this *Preferences) UnmarshalJSON(bytes []byte) error {

	//	var prefs map[string]*Preference
	var prefs []Preference
	if err := json.Unmarshal(bytes, &prefs); err != nil {
		return err
	}
	this.Root = newRootNameNode()
	this.AllPrefs = make([]Preference, 0)

	for _, pref := range prefs {
		this.AddPreference(pref.Name, pref.Value)
	}

	return nil
}

func NewPreferences() *Preferences {
	p := new(Preferences)
	p.Root = newRootNameNode()
	p.AllPrefs = make([]Preference, 0)
	return p
}

func (this *Preferences) AddPreference(name string, value string) {

	fmt.Printf("AddPreference: %v, %v\n", name, value)

	//	p := Preference{name, value}
	//parts := strings.Split(name, ".")
	//curNode := this.Root

	p := Preference{name, value}
	this.AllPrefs = append(this.AllPrefs, p)
	//	this.AllPrefs[name] = &p
	this.Root.add(name, p)
}

func (this *NameNode) getPreferences() map[string]string {

	res := make(map[string]string, 0)
	if this.children != nil {
		for _, child := range this.children {
			r := child.getPreferences()
			add(&res, r)
		}
		return res
	} else {
		res[this.preference.Name] = this.preference.Value
	}
	return res

}

func (this *NameNode) String() string {

	var b strings.Builder

	b.WriteString(fmt.Sprintf("NameNode:  %s, children: \n", this.name))
	for _, c := range this.children {
		b.WriteString(fmt.Sprintf("     %s\n", c.name))
	}

	return b.String()
}

func (this *Preferences) String() string {
	return "Preferences: \n\n" + this.Root.String()
}

func (this *Preferences) findNode(name string, from *NameNode) *NameNode {

	fmt.Printf("findNode: %s from %s\n", name, from.name)
	if len(name) <= 0 {
		return from
	}
	parts := strings.Split(name, ".")
	first := parts[0]

	rest := ""
	match, has := from.children[first]

	if has {
		if len(parts) == 1 {
			return match
		}
		if len(parts) > 1 {
			rest = strings.Trim(name, first+".")
		}
	} else {
		return nil
	}
	return this.findNode(rest, match)

}

/*
   result contains all preferences at or under the specified path
*/
func (this *Preferences) GetPreferences(name string) map[string]string {

	startNode := this.findNode(name, this.Root)

	if startNode != nil {
		fmt.Printf("GetPreferences starting from node: %v\n", startNode)
		return startNode.getPreferences()
	}
	return make(map[string]string, 0)
}

func DefaultPreferences() *Preferences {

	prefs := NewPreferences()

	prefs.AddPreference("brackets.background-color", "#3C3D40")
	prefs.AddPreference("brackets.connector-color", "#FAF30B")
	prefs.AddPreference("brackets.game.border-color", "#97a872")
	prefs.AddPreference("brackets.game.slot1.background-color", "#D5D5D7")
	prefs.AddPreference("brackets.game.slot2.background-color", "#EBEBEF")
	prefs.AddPreference("brackets.game.slot1.font-color", "#272822")
	prefs.AddPreference("brackets.game.slot2.font-color", "#272822")
	prefs.AddPreference("brackets.game.font-color", "#863845")
	prefs.AddPreference("brackets.game.winners.font-color", "#57294b")
	prefs.AddPreference("brackets.game.winners.background-color", "#863845")
	prefs.AddPreference("brackets.game.losers.background-color", "#384863")

	return prefs

}

func LoadPreferences(tx dbr.SessionRunner, subject string) (*Preferences, error) {

	fmt.Printf("LoadPreferences\n")
	result, err := tx.Select("preferences_data").From("preferences").
		Where("subject = ?", subject).Rows()

	//mt.Printf("LoadPreferences 1\n")

	if err != nil || result.Next() == false {
		p := DefaultPreferences()
		StorePreferences(tx, *p, subject)
		fmt.Printf("Select Preferences failed: %v\n", err)
		fmt.Printf("returning default preferences data\n")
		return p, nil
	}
	var data []byte
	if err := result.Scan(&data); err != nil {
		fmt.Printf("Scan preferences data failed: %v\n", err)
		return nil, nil
	}

	result.Close()

	fmt.Printf("data = %v\n", string(data))

	var prefs Preferences

	if err = json.Unmarshal(data, &prefs); err != nil {
		fmt.Printf("Unmarshal preferences data failed: %v", err)
		return nil, nil
	}
	return &prefs, nil
}

func StorePreferences(tx dbr.SessionRunner, prefs Preferences, subject string) error {

	fmt.Printf("Store preferences\n")

	//result, err := db.Select("preferences_data").From("preferences").
	//		Where("user_id = ? and id = 1", userId).Rows()

	q := fmt.Sprintf("select count(*) from preferences where id = 1 and subject = '%v'",
		subject)
	count, err := tx.SelectBySql(q).ReturnInt64()

	update := true
	if count <= 0 {
		update = false
	}

	var data []byte
	if data, err = json.Marshal(&prefs); err != nil {
		fmt.Printf("Store,Marshall Preferences failed: %v", err)
		return err
	}

	if update {
		fmt.Printf("Store preferences - update\n")
		UpdateStmt := tx.Update("preferences").Set("preferences_data", data)
		_, err := UpdateStmt.Where("subject = ? and id = 1", subject).Exec()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
	} else {
		fmt.Printf("Store preferences - insert\n")
		var result sql.Result
		result, err = tx.InsertInto("preferences").
			Columns("id", "subject", "tournament_id", "preferences_data").
			Values(1, subject, 0, data).Exec()
		var c int64
		if result == nil {
			return nil
		}
		c, err = result.RowsAffected()
		if c <= 1 || err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}

	}
	return nil
}

func add(all *map[string]string, newEntries map[string]string) {
	for nm, s := range newEntries {
		(*all)[nm] = s
	}
}
