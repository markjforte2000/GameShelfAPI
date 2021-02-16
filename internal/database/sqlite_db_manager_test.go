package database

import (
	"database/sql"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"os"
	"path"
	"testing"
)

type AlterFunc func(t *testing.T, toAlter *game.Game, original *game.Game, manager *sqliteDBManager)

func TestAccessInvalidGameData(t *testing.T) {
	manager := initTestSQLiteManager(t)
	inserted := newTestGame(t)
	manager.SaveGameData(inserted)
	gameFile := &game.GameFile{
		Title:    "Invalid title",
		Year:     "2021",
		Platform: "Unknown",
		FileName: "invalid.rom",
	}
	_, exists := manager.AccessGameData(gameFile)
	if exists {
		t.Errorf("Accessed game should not exist")
	}
}

func TestAccessGameData(t *testing.T) {
	manager := initTestSQLiteManager(t)
	inserted := newTestGame(t)
	manager.SaveGameData(inserted)
	gameFile := &game.GameFile{
		Title:    inserted.Title,
		Year:     "2021",
		Platform: "Unknown",
		FileName: t.Name() + ".rom",
	}
	accessed, exists := manager.AccessGameData(gameFile)
	if !exists {
		t.Error("Accessed game does not exist")
	}
	if !inserted.Equal(accessed) {
		t.Errorf("Accessed game does not equal inserted game")
	}
}

func TestAlterGameData(t *testing.T) {
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.Title = "altered title"
		manager.alterGameData(toAlter)
	})
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.Genres = toAlter.Genres[:1]
		manager.alterGameData(toAlter)
	})
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game, original *game.Game,
		manager *sqliteDBManager) {
		newGenre := &game.Genre{
			Name: "Test Add Genre",
			ID:   5,
		}
		toAlter.Genres = append(toAlter.Genres, newGenre)
		manager.alterGameData(toAlter)
	})
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.Cover.RemoteURL = "https://alteredurl.lan"
		manager.alterGameData(toAlter)
	})
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.InvolvedCompanies = append(toAlter.InvolvedCompanies,
			&game.InvolvedCompany{
				Name:      "added company",
				ID:        3,
				Publisher: true,
				Developer: false,
			})
		manager.alterGameData(toAlter)
	})
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.InvolvedCompanies = toAlter.InvolvedCompanies[:1]
		manager.alterGameData(toAlter)
	})
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.InvolvedCompanies[0].Name = "Altered Name"
		manager.alterGameData(toAlter)
	})
}

func TestAlterBasicGameData(t *testing.T) {
	_ = helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.Title = "altered title"
		manager.alterBasicGameData(toAlter, original)
	})
}

func TestRemoveGenre(t *testing.T) {
	manager := helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.Genres = toAlter.Genres[:1]
		manager.alterGameGenre(toAlter, original)
	})
	// ensure genre is not deleted
	genres := manager.queryGenreTable(`SELECT * FROM genre`)
	if len(genres) != 2 {
		t.Error("Genre was deleted")
	}
}

func TestAddGenre(t *testing.T) {
	manager := helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game, original *game.Game,
		manager *sqliteDBManager) {
		newGenre := &game.Genre{
			Name: "Test Add Genre",
			ID:   5,
		}
		toAlter.Genres = append(toAlter.Genres, newGenre)
		manager.alterGameGenre(toAlter, original)
	})
	// ensure new genre is added and association isn't just created
	genres := manager.queryGenreTable(`SELECT * FROM genre`)
	if len(genres) != 3 {
		t.Error("Genre was not added to database")
	}
}

func TestInsertNewArtwork(t *testing.T) {
	manager := helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.Cover = &game.Artwork{
			RemoteURL: "www.newartwork.lan",
			ID:        2,
		}
		manager.alterGameArtwork(toAlter, original)
	})
	artworks := manager.queryArtworkTable(`SELECT * FROM artwork`)
	if len(artworks) <= 1 {
		t.Error("Failed to add artwork")
	}
}

func TestAlterExistingArtwork(t *testing.T) {
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.Cover.RemoteURL = "https://alteredurl.lan"
		manager.alterGameArtwork(toAlter, original)
	})
}

func TestAddCompany(t *testing.T) {
	manager := helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.InvolvedCompanies = append(toAlter.InvolvedCompanies,
			&game.InvolvedCompany{
				Name:      "added company",
				ID:        3,
				Publisher: true,
				Developer: false,
			})
		manager.alterCompanyData(toAlter, original)
	})
	companies := manager.queryCompanyTable(`SELECT * FROM company`)
	if len(companies) < 3 {
		t.Error("Company was not added")
	}
}

func TestDeleteCompany(t *testing.T) {
	manager := helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.InvolvedCompanies = toAlter.InvolvedCompanies[:1]
		manager.alterCompanyData(toAlter, original)
	})
	companies := manager.queryCompanyTable(`SELECT * FROM company`)
	if len(companies) > 1 {
		t.Error("Company was not deleted")
	}
}

func TestAlterCompany(t *testing.T) {
	helperTestAlterGame(t, func(t *testing.T, toAlter *game.Game,
		original *game.Game, manager *sqliteDBManager) {
		toAlter.InvolvedCompanies[0].Name = "Altered Name"
		manager.alterCompanyData(toAlter, original)
	})
}

func helperTestAlterGame(t *testing.T, alterFunc AlterFunc) *sqliteDBManager {
	original := newTestGame(t)
	copyGame := newTestGame(t)
	if !original.Equal(copyGame) {
		t.Errorf("Original game does not equal copy:\n"+
			"%v\ndoes not equal\n%v", original, copyGame)
	}
	manager := initTestSQLiteManager(t)
	manager.insertNewGame(original)
	alterFunc(t, copyGame, original, manager)
	resultingGames := manager.queryGameTable(`SELECT * FROM game WHERE id=0`)
	if len(resultingGames) == 0 {
		t.Error("Query for game did not result in any")
	}
	for _, resultingGame := range resultingGames {
		if resultingGame.Equal(original) {
			t.Error("Original game still exists in table")
		}
	}
	g := resultingGames[0]
	if g.Equal(original) {
		t.Errorf("Game did not change")
	}
	return manager
}

func TestGameIO(t *testing.T) {
	g := newTestGame(t)
	manager := initTestSQLiteManager(t)
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
		"game":             {"id", "title", "releaseDate", "summary", "filename", "coverID"},
		"company":          {"id", "gameID", "name", "publisher", "developer"},
		"genre":            {"id", "name"},
		"artwork":          {"id", "remoteURL", "gameID"},
		"genreAssociation": {"genreID", "gameID"},
	}
	manager := initTestSQLiteManager(t)
	for table, expectedColumns := range expectedTables {
		testTableHelper(t, manager.db, table, expectedColumns)
	}
	manager.db.Close()
}

func TestCreateTable(t *testing.T) {
	manager := initTestSQLiteManager(t)
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
	if len(existingColumns) != len(expectedColumns) {
		t.Errorf("Column size mismatch: expected %v got %v",
			len(expectedColumns), len(existingColumns))
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

func initTestSQLiteManager(t *testing.T) *sqliteDBManager {
	filename := path.Join(TestDir, t.Name()+".db")
	_ = os.Remove(filename)
	manager := new(sqliteDBManager)
	manager.init(filename)
	return manager
}
