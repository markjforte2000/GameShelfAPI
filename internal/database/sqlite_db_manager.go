package database

import (
	"database/sql"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type rowHandler func(rows *sql.Rows) error

type sqliteDBManager struct {
	db *sql.DB
}

type alteredCompanyData struct {
	company *game.InvolvedCompany
	altered bool
	deleted bool
	create  bool
}

type alteredGenreData struct {
	genre             *game.Genre
	createGenre       bool
	createAssociation bool
	deleteAssociation bool
}

func (manager *sqliteDBManager) AccessGameData(gameFile *game.GameFile) (*game.Game, bool) {
	games := manager.queryGameTable(`SELECT * FROM game WHERE filename=$1`, gameFile.FileName)
	if games == nil {
		log.Printf("Entry not found in database for %s", gameFile.FileName)
		return nil, false
	}
	log.Printf("Entry found in database for %s", gameFile.FileName)
	return games[0], true
}

func (manager *sqliteDBManager) SaveGameData(g *game.Game) {
	if manager.doesGameExist(g) {
		manager.alterGameData(g)
	} else {
		manager.insertNewGame(g)
	}
}

func (manager *sqliteDBManager) alterGameData(alteredGame *game.Game) {
	games := manager.queryGameTable(`SELECT * FROM game WHERE id=$1`, alteredGame.ID)
	if len(games) == 0 {
		log.Fatalf("Database error - game exists in table but can not be found")
	}
	original := games[0]
	// check basic company data
	if (alteredGame.ID == original.ID) && (alteredGame.Filename != original.Filename ||
		!alteredGame.ReleaseDate.Equal(original.ReleaseDate) || alteredGame.Summary != original.Summary ||
		alteredGame.Title != original.Title) {
		manager.alterBasicGameData(alteredGame, original)
	}
	// check artwork
	if !alteredGame.Cover.Equal(original.Cover) {
		manager.alterGameArtwork(alteredGame, original)
	}
	manager.alterCompanyData(alteredGame, original)
	manager.alterGameGenre(alteredGame, original)
}

func (manager *sqliteDBManager) alterBasicGameData(alteredGame *game.Game,
	originalGame *game.Game) {
	updateQuery := `UPDATE game SET 
			title=$1,releaseDate=$2,summary=$3,filename=$4
			WHERE id=$5`
	manager.executeInsert(updateQuery, alteredGame.Title, alteredGame.ReleaseDate.Unix(),
		alteredGame.Summary, alteredGame.Filename, originalGame.ID)
}

func (manager *sqliteDBManager) alterGameGenre(alteredGame *game.Game,
	originalGame *game.Game) {
	alteredGenres := manager.getAlteredGenreInfo(alteredGame, originalGame)
	for _, info := range alteredGenres {
		if info.createGenre {
			manager.insertGenre(info.genre)
		}
		if info.createAssociation {
			manager.insertGenreAssociation(info.genre, alteredGame.ID)
		} else if info.deleteAssociation {
			manager.deleteGenreAssociation(alteredGame.ID, info.genre.ID)
		}
	}
}

func (manager *sqliteDBManager) deleteGenreAssociation(gameID int, genreID int) {
	query := `DELETE FROM genreAssociation WHERE gameID=$1 AND genreID=$2`
	manager.executeInsert(query, gameID, genreID)
}

func (manager *sqliteDBManager) getAlteredGenreInfo(alteredGame *game.Game,
	originalGame *game.Game) []*alteredGenreData {
	alteredGenreMap := make(map[int]*alteredGenreData)
	// load original genres
	for _, originalGenre := range originalGame.Genres {
		alteredGenreMap[originalGenre.ID] = &alteredGenreData{
			genre:             originalGenre,
			createGenre:       false,
			createAssociation: false,
			deleteAssociation: true,
		}
	}
	// load altered genres
	for _, alteredGenre := range alteredGame.Genres {
		if data, exists := alteredGenreMap[alteredGenre.ID]; exists {
			data.deleteAssociation = false
		} else {
			alteredGenreMap[alteredGenre.ID] = &alteredGenreData{
				genre:             alteredGenre,
				createGenre:       false,
				createAssociation: true,
				deleteAssociation: false,
			}
			if !manager.doesGenreExist(alteredGenre) {
				alteredGenreMap[alteredGenre.ID].createGenre = true
			}
		}
	}
	// convert to list
	var genreData []*alteredGenreData
	for _, data := range alteredGenreMap {
		genreData = append(genreData, data)
	}
	return genreData
}

func (manager *sqliteDBManager) alterGameArtwork(alteredGame *game.Game,
	originalGame *game.Game) {
	if !manager.doesArtworkExist(alteredGame.Cover.ID) {
		manager.insertArtwork(alteredGame.Cover, alteredGame.ID)
	}
	if originalGame.Cover.ID == alteredGame.Cover.ID && !originalGame.Cover.Equal(alteredGame.Cover) {
		updateQuery := `UPDATE artwork SET remoteURL=$1 WHERE id=$2`
		manager.executeInsert(updateQuery, alteredGame.Cover.RemoteURL, alteredGame.Cover.ID)
	} else if originalGame.Cover.ID != alteredGame.Cover.ID {
		updateQuery := `UPDATE game SET coverID=$1 WHERE id=$2`
		manager.executeInsert(updateQuery, alteredGame.Cover.ID, alteredGame.ID)
	}
}

func (manager *sqliteDBManager) doesArtworkExist(artworkID int) bool {
	artworks := manager.queryArtworkTable(`SELECT * from artwork where id=$1`, artworkID)
	return len(artworks) > 0
}

func (manager *sqliteDBManager) alterCompanyData(alteredGame *game.Game,
	originalGame *game.Game) {
	alteredCompanyInfo := getAlteredCompanyInfo(alteredGame, originalGame)
	for _, info := range alteredCompanyInfo {
		if info.deleted {
			manager.deleteCompany(info.company)
		} else if info.altered {
			manager.alterCompany(info.company)
		} else if info.create {
			manager.insertCompanyData(info.company, alteredGame.ID)
		}
	}
}

func (manager *sqliteDBManager) deleteCompany(company *game.InvolvedCompany) {
	deleteQuery := `DELETE FROM company WHERE id=$1`
	manager.executeInsert(deleteQuery, company.ID)
}

func (manager *sqliteDBManager) alterCompany(company *game.InvolvedCompany) {
	updateQuery := `UPDATE company SET 
					name=$1,publisher=$2,developer=$3
					WHERE id=$4`
	manager.executeInsert(updateQuery, company.Name, company.Publisher, company.Developer, company.ID)
}

func getAlteredCompanyInfo(alteredGame *game.Game, originalGame *game.Game) []*alteredCompanyData {
	alteredCompanyMap := make(map[int]*alteredCompanyData)
	// load original companies
	for _, company := range originalGame.InvolvedCompanies {
		alteredCompanyMap[company.ID] = &alteredCompanyData{
			company: company,
			altered: false,
			deleted: true,
			create:  false,
		}
	}
	// load altered companies
	for _, company := range alteredGame.InvolvedCompanies {
		if data, exists := alteredCompanyMap[company.ID]; exists {
			data.deleted = false
			if data.company.ID == company.ID && !data.company.Equal(company) {
				data.company = company
				data.altered = true
			}
		} else {
			alteredCompanyMap[company.ID] = &alteredCompanyData{
				company: company,
				altered: false,
				deleted: false,
				create:  true,
			}
		}
	}
	// convert to list
	var alteredCompanies []*alteredCompanyData
	for _, companyData := range alteredCompanyMap {
		alteredCompanies = append(alteredCompanies, companyData)
	}
	return alteredCompanies
}

func (manager *sqliteDBManager) insertNewGame(g *game.Game) {
	manager.insertGameData(g)
	for _, company := range g.InvolvedCompanies {
		manager.insertCompanyData(company, g.ID)
	}
	manager.insertArtwork(g.Cover, g.ID)
	for _, genre := range g.Genres {
		if !manager.doesGenreExist(genre) {
			manager.insertGenre(genre)
		}
		manager.insertGenreAssociation(genre, g.ID)
	}
}

func (manager *sqliteDBManager) queryGameTable(queryString string,
	args ...interface{}) []*game.Game {
	var games []*game.Game
	var coverID int
	manager.queryTable(queryString, func(rows *sql.Rows) error {
		g := new(game.Game)
		var timeStamp int64
		err := rows.Scan(&g.ID, &g.Title, &timeStamp, &g.Summary, &g.Filename, &coverID)
		if err != nil {
			return err
		}
		g.ReleaseDate = util.UnixTimestampToDate(timeStamp)
		games = append(games, g)
		return nil
	}, args...)
	for _, g := range games {
		// load artwork
		relatedArtwork := manager.queryArtworkTable(`SELECT * FROM artwork
			WHERE id=$1`, coverID)
		if len(relatedArtwork) > 0 {
			g.Cover = relatedArtwork[0]
		}
		// load companies
		involvedCompanies := manager.queryCompanyTable(`SELECT * FROM company
			WHERE gameID=$1`, g.ID)
		g.InvolvedCompanies = involvedCompanies
		// load related genres
		genreAssociations := manager.queryGenreAssociationTable(`
			SELECT * FROM genreAssociation WHERE gameID=$1`, g.ID)
		g.Genres = []*game.Genre{}
		for _, association := range genreAssociations {
			genres := manager.queryGenreTable(
				`SELECT * FROM genre WHERE id=$1`, association.genreID)
			if len(genres) != 0 {
				g.Genres = append(g.Genres, genres[0])
			}
		}
	}
	return games
}

func (manager *sqliteDBManager) queryArtworkTable(queryString string,
	args ...interface{}) []*game.Artwork {
	var artwork []*game.Artwork
	manager.queryTable(queryString, func(rows *sql.Rows) error {
		art := new(game.Artwork)
		var gameID int
		err := rows.Scan(&art.ID, &art.RemoteURL, &gameID)
		if err != nil {
			return err
		}
		artwork = append(artwork, art)
		return nil
	}, args...)
	return artwork
}

func (manager *sqliteDBManager) queryCompanyTable(queryString string,
	args ...interface{}) []*game.InvolvedCompany {
	var companies []*game.InvolvedCompany
	manager.queryTable(queryString, func(rows *sql.Rows) error {
		company := new(game.InvolvedCompany)
		var gameID int
		err := rows.Scan(&company.ID, &gameID, &company.Name,
			&company.Publisher, &company.Developer)
		if err != nil {
			return err
		}
		companies = append(companies, company)
		return nil
	}, args...)
	return companies
}

func (manager *sqliteDBManager) queryGenreTable(queryString string,
	args ...interface{}) []*game.Genre {
	var genres []*game.Genre
	manager.queryTable(queryString, func(rows *sql.Rows) error {
		genre := new(game.Genre)
		err := rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return err
		}
		genres = append(genres, genre)
		return nil
	}, args...)
	return genres
}

func (manager *sqliteDBManager) queryGenreAssociationTable(queryString string,
	args ...interface{}) []*genreAssociation {
	var genres []*genreAssociation
	manager.queryTable(queryString, func(rows *sql.Rows) error {
		genre := new(genreAssociation)
		err := rows.Scan(&genre.genreID, &genre.gameID)
		if err != nil {
			return err
		}
		genres = append(genres, genre)
		return nil
	}, args...)
	return genres
}

func (manager *sqliteDBManager) queryTable(queryString string,
	handler rowHandler, args ...interface{}) {
	rows, err := manager.db.Query(queryString, args...)
	if err != nil {
		log.Fatalf("Failed to execute query: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = handler(rows)
		if err != nil {
			log.Fatalf("Failed to process row: %v\n", err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("Failed to process rows: %v\n", err)
	}
}

func (manager *sqliteDBManager) doesGameExist(g *game.Game) bool {
	query := `SELECT id FROM game where id=$1`
	return manager.doesRowExist(query, g.ID)
}

func (manager *sqliteDBManager) doesGenreExist(genre *game.Genre) bool {
	query := `SELECT id FROM genre where id=$1`
	return manager.doesRowExist(query, genre.ID)
}

func (manager *sqliteDBManager) doesRowExist(query string, args ...interface{}) bool {
	row := manager.db.QueryRow(query, args...)
	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatalf("Failed to check for row: %v\n", err)
	}
	return true
}

func (manager *sqliteDBManager) insertGameData(g *game.Game) {
	insertStatement := `INSERT INTO game (id, title, releaseDate, summary, filename, coverID)
		VALUES ($1, $2, $3, $4, $5, $6)`
	manager.executeInsert(insertStatement, g.ID, g.Title,
		g.ReleaseDate.Unix(), g.Summary, g.Filename, g.Cover.ID)
}

func (manager *sqliteDBManager) insertCompanyData(company *game.InvolvedCompany, relatedGameID int) {
	insertStatement := `INSERT INTO company (id, gameID, name, publisher, developer)
		VALUES ($1, $2, $3, $4, $5)`
	manager.executeInsert(insertStatement,
		company.ID, relatedGameID, company.Name, company.Publisher, company.Developer)
}

func (manager *sqliteDBManager) insertArtwork(artwork *game.Artwork, gameID int) {
	insertStatement := `INSERT INTO artwork (id, remoteURL, gameID)
		VALUES ($1, $2, $3)`
	manager.executeInsert(insertStatement, artwork.ID, artwork.RemoteURL, gameID)
}

func (manager *sqliteDBManager) insertGenre(genre *game.Genre) {
	insertStatement := `INSERT INTO genre (id, name) VALUES ($1, $2)`
	manager.executeInsert(insertStatement, genre.ID, genre.Name)
}

func (manager *sqliteDBManager) insertGenreAssociation(genre *game.Genre, gameID int) {
	insertStatement := `INSERT INTO genreAssociation (genreID, gameID) VALUES ($1, $2)`
	manager.executeInsert(insertStatement, genre.ID, gameID)
}

func (manager *sqliteDBManager) executeInsert(statement string, args ...interface{}) {
	_, err := manager.db.Exec(statement, args...)
	if err != nil {
		log.Fatalf("Failed to excecute into db: '%v': %v\n", statement, err)
	}
}

func (manager *sqliteDBManager) init(dbFile string) {
	log.Printf("Initializing Database at %v\t", dbFile)
	// check if file exists
	if _, err := os.Stat(dbFile); err == os.ErrNotExist {
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatalf("Failed to create database file: %v\n", err)
		}
		file.Close()
	} else if err != nil {
		log.Fatalf("Failed to get status of database file: %v\n")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database file: %v\n", err)
	}
	manager.db = db
	manager.initializeTables()
}

func (manager *sqliteDBManager) initializeTables() {
	manager.createTable(
		`CREATE TABLE IF NOT EXISTS game (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		releaseDate INTEGER NOT NULL,
		summary TEXT NOT NULL,
    	filename TEXT NOT NULL UNIQUE,
    	coverID INTEGER NOT NULL,
    	FOREIGN KEY (coverID) REFERENCES artwork(ID)
		);`)
	manager.createTable(
		`CREATE TABLE IF NOT EXISTS company (
		id INTEGER NOT NULL ,
		gameID INTEGER NOT NULL,
		name Text NOT NULL,
		publisher INTEGER(1),
		developer INTEGER(1),
		FOREIGN KEY (gameID) REFERENCES game(ID),
		PRIMARY KEY (id, gameID)
		);`)
	manager.createTable(
		`CREATE TABLE IF NOT EXISTS genre (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
		);`)
	manager.createTable(
		`CREATE TABLE IF NOT EXISTS artwork (
		id INTEGER PRIMARY KEY,
		remoteURL TEXT NOT NULL,
		gameID INTEGER NOT NULL,
		FOREIGN KEY (gameID) REFERENCES game(id)
		);`)
	manager.createTable(
		`CREATE TABLE IF NOT EXISTS genreAssociation (
		genreID INTEGER NOT NULL,
		gameID INTEGER NOT NULL,
		FOREIGN KEY (genreID) REFERENCES genre(id),
		FOREIGN KEY (gameID) REFERENCES game(id),
		PRIMARY KEY (genreID, gameID)
		)`)
}

func (manager *sqliteDBManager) createTable(createString string) {
	statement, err := manager.db.Prepare(createString)
	if err != nil {
		log.Fatalf("Failed to create sql statement: %v\n", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}
}
