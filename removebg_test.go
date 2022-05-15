package removebg

import (
	"os"
	"testing"
)

func getRemoveBG() *RemoveBG {
	os.Getenv("REMOVE_BG_KEY")
	return NewRemoveBG(os.Getenv("REMOVE_BG_KEY"))
}

func TestNewOpt(t *testing.T) {
	t.Log(NewRemoveOption())
}

func TestOptCheck(t *testing.T) {
	opt := NewRemoveOption()
	if opt.check() != nil {
		t.Error(opt)
	}
}

// func TestFromFile(t *testing.T) {
// 	rg := getRemoveBG()
// 	err := rg.RemoveFromFile("", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestFromURL(t *testing.T) {
	rg := getRemoveBG()
	err := rg.RemoveFromURL("https://avatars.githubusercontent.com/u/5282810?v=4", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestFromBase64(t *testing.T) {
	rg := getRemoveBG()
	err := rg.RemoveFromBase64("", nil)
	if err != nil {
		t.Error(err)
	}
}
