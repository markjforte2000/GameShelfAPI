package database

import (
	"fmt"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"path"
	"testing"
)

func TestAsyncInsertDBManager(t *testing.T) {
	manager := NewAsyncInsertManager(path.Join(TestDir, "async.db"))
	for i := 0; i < 100; i++ {
		go testHelper(t, manager, i)
	}
	_, exists := manager.AccessGameData(&game.GameFile{
		Title:    "TestAsyncInsertDBManager54",
		Year:     "2021",
		Platform: "Unknown",
		FileName: "TestAsyncInsertDBManager.rom54",
	})
	if !exists {
		t.Fatalf("Database did not add file")
	}
}

func testHelper(t *testing.T, manager Manager, num int) {
	g := newTestGame(t)
	g.ID = num
	g.Cover.ID = num
	g.Title += fmt.Sprintf("%v", num)
	g.Filename += fmt.Sprintf("%v", num)
	manager.SaveGameData(g)
}
