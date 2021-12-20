package brackets

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPreferences_AddPreference(*testing.T) {
	fmt.Printf("TestPreferences_AddPreferences\n")

	prefs := NewPreferences()

	prefs.AddPreference("brackets.bkg", "#877892")
	prefs.AddPreference("brackets.game.border", "#97a872")
	prefs.AddPreference("brackets.game.slot1.bkg", "#a72892")
	prefs.AddPreference("brackets.game.slot2.bkg", "#8748c2")
	prefs.AddPreference("brackets.game.slot1.txt", "#272822")
	prefs.AddPreference("brackets.game.slot2.txt", "#a73882")
	prefs.AddPreference("brackets.game.txt", "#171812")
	//	prefs.AddPreference("brackets.game.txt", "#274832")

	fmt.Printf("prefs = %v\n", prefs)

	matches := prefs.GetPreferences("brackets.game.slot1")
	for nm, v := range matches {
		fmt.Printf("pref name: %s, value %v\n", nm, v)
	}
}

func TestPreferences_Codec(*testing.T) {

	prefs := NewPreferences()

	prefs.AddPreference("brackets.bkg", "#877892")
	prefs.AddPreference("brackets.game.border", "#97a872")
	prefs.AddPreference("brackets.game.slot1.bkg", "#a72892")
	prefs.AddPreference("brackets.game.slot2.bkg", "#8748c2")
	prefs.AddPreference("brackets.game.slot1.txt", "#272822")
	prefs.AddPreference("brackets.game.slot2.txt", "#a73882")
	prefs.AddPreference("brackets.game.txt", "#171812")
	//	prefs.AddPreference("brackets.game.txt", "#274832")

	bytes, err := json.Marshal(prefs)
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}

	fmt.Printf("json = %s\n", string(bytes))

	var p Preferences
	err = json.Unmarshal(bytes, &p)

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}

	res := p.GetPreferences("brackets.game.slot1")

	fmt.Printf("res has %d, should have %d", len(res), 2)


	fmt.Printf("preferences after decoding: %v\n", p)
}
