package database

import (
	"database/sql"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"os"
	"path"
	"testing"
	"time"
)

const TestDir = "D:/Documents/Program/Web/GameLibrary/GameShelfAPI/.test/"

func TestAlterGame(t *testing.T) {
	original := newTestGame(t)
	copyGame := newTestGame(t)
	if !original.Equal(copyGame) {
		t.Errorf("Original game does not equal copy:\n"+
			"%v\ndoes not equal\n%v", original, copyGame)
	}
	manager := initTestManager(t)
	manager.insertNewGame(original)
	original.Title = "Changed title"
	manager.alterGameDate(original)
	resultingGames := manager.queryGameTable(`SELECT * FROM game WHERE id=0`)
	if len(resultingGames) == 0 {
		t.Error("Query for game did not result in any")
	}
	for _, resultingGame := range resultingGames {
		if resultingGame.Equal(copyGame) {
			t.Error("Original game still exists in table")
		}
	}
	g := resultingGames[0]
	if g.Title != "Changed title" {
		t.Errorf("New title wrong expected: 'Changed title' got: '%v'", g.Title)
	}
}

func TestGameIO(t *testing.T) {
	g := newTestGame(t)
	manager := initTestManager(t)
	manager.insertNewGame(g)
	resultingGames := manager.queryGameTable(`SELECT * FROM game WHERE id=0`)
	if len(resultingGames) == 0 {
		t.Error("Query into game table did not return any valid games")
	}
	resultingGame := resultingGames[0]
	if !g.Equal(resultingGame) {
		t.Errorf("Query for game did not result in same game: \n"+
			"%v\ndoes not equal\n%v",
			util.GameToPrettyString(g), util.GameToPrettyString(resultingGame))
	}
}

func TestInitializeTables(t *testing.T) {
	expectedTables := map[string][]string{
		"game":             {"id", "title", "releaseDate", "summary", "filename"},
		"company":          {"id", "gameID", "name", "publisher", "developer"},
		"genre":            {"id", "name"},
		"artwork":          {"id", "remoteURL", "gameID"},
		"genreAssociation": {"genreID", "gameID"},
	}
	manager := initTestManager(t)
	for table, expectedColumns := range expectedTables {
		testTableHelper(t, manager.db, table, expectedColumns)
	}
	manager.db.Close()
}

func TestCreateTable(t *testing.T) {
	manager := initTestManager(t)
	manager.createTable(`
		CREATE TABLE testTable (
			id INTEGER PRIMARY KEY,
			testField1 TEXT NOT NULL,
			testField2 TEXT
		)
	`)
	testTableHelper(t, manager.db, "testTable", []string{
		"id", "testField1", "testField2",
	})
	manager.db.Close()
}

// Checks to see if table contains the columns given
// returns true if table contains given columns and only given columns
func testTableHelper(t *testing.T, db *sql.DB, table string, expectedColumns []string) bool {
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		t.Fatalf("Failed to get rows from table: %v\n", err)
	}
	existingColumns, err := rows.Columns()
	if err != nil {
		t.Fatalf("Failed to get column information from rows: %v\n", err)
	}
	for _, existingColumn := range existingColumns {
		columnShouldExist := false
		for _, expectedColumn := range expectedColumns {
			if existingColumn == expectedColumn {
				columnShouldExist = true
				break
			}
		}
		if !columnShouldExist {
			t.Fatalf("Table %v exists but it should not", existingColumn)
		}
	}
	for _, expectedColumn := range expectedColumns {
		columnExists := false
		for _, existingColumn := range existingColumns {
			if existingColumn == expectedColumn {
				columnExists = true
				break
			}
		}
		if !columnExists {
			t.Fatalf("Table %v does not exist but it should", expectedColumn)
		}
	}
	return true
}

func initTestManager(t *testing.T) *sqliteDBManager {
	filename := path.Join(TestDir, t.Name()+".db")
	_ = os.Remove(filename)
	manager := new(sqliteDBManager)
	manager.init(filename)
	return manager
}

func newTestGame(t *testing.T) *game.Game {
	return &game.Game{
		ID:    0,
		Title: t.Name() + " Game Title",
		ReleaseDate: time.Date(2021, 2, 14, 12,
			0, 0, 0, time.UTC),
		InvolvedCompanies: []*game.InvolvedCompany{
			{
				Name:      t.Name() + " Company Publisher",
				ID:        1,
				Publisher: true,
				Developer: false,
			},
			{
				Name:      t.Name() + " IO Company Developer",
				ID:        2,
				Publisher: false,
				Developer: true,
			},
		},
		Summary: t.Name() + " Summary",
		Genres: []*game.Genre{
			{
				Name: t.Name() + " Genre 1",
				ID:   0,
			},
			{
				Name: t.Name() + " Genre 2",
				ID:   1,
			},
		},
		Cover: &game.Artwork{
			RemoteURL: "https://www." + t.Name() + ".lan",
			ID:        1,
		},
		Filename: t.Name() + ".rom",
	}
}
