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

func (manager *sqliteDBManager) AccessGameDate(gameFile *game.GameFile) *game.Game {
	panic("implement me")
}

func (manager *sqliteDBManager) SaveGameData(g *game.Game) {
	if manager.doesGameExist(g) {
		manager.alterGameDate(g)
	} else {
		manager.insertNewGame(g)
	}
}

func (manager *sqliteDBManager) alterGameDate(g *game.Game) {
	games := manager.queryGameTable(`SELECT * FROM game WHERE id=$1`, g.ID)
	if len(games) == 0 {
		log.Fatalf("Database error - game exists in table but can not be found")
	}
	currentGame := games[0]
	if g.ID != currentGame.ID || g.Filename != currentGame.Filename ||
		!g.ReleaseDate.Equal(currentGame.ReleaseDate) || g.Summary != currentGame.Summary ||
		g.Title != currentGame.Title {
		updateQuery := `UPDATE game SET 
			id=$1,title=$2,releaseDate=$3,summary=$4,filename=$5
			WHERE id=$6`
		manager.executeInsert(updateQuery, g.ID, g.Title, g.ReleaseDate.Unix(),
			g.Summary, g.Filename, currentGame.ID)
	}
	if g.Cover.ID != currentGame.Cover.ID ||
		g.Cover.RemoteURL != currentGame.Cover.RemoteURL {
		updateQuery := `UPDATE artwork SET 
			id=$1,remoteURL=$2
			WHERE id=$3`
		manager.executeInsert(updateQuery, g.Cover.ID, g.Cover.RemoteURL)
	}
	var companiesToRemove []*game.InvolvedCompany
	copy(companiesToRemove, currentGame.InvolvedCompanies)
	for _, company := range currentGame.InvolvedCompanies {
		matchID := false
		for _, other := range g.InvolvedCompanies {
			if company.ID == other.ID && (company.Name != other.Name ||
				company.Developer != other.Developer || company.Publisher != other.Publisher) {
				updateQuery := `UPDATE company SET 
					id=$1,gameID=$2,name=$3,publisher=$4,developer=$5
					WHERE id=$6`
				manager.executeInsert(updateQuery, other.ID, g.ID, other.Name,
					other.Publisher, other.Developer, other.ID)
			}
			if company.ID == other.ID {
				matchID = true
			}
		}
		if !matchID {
			companiesToRemove = append(companiesToRemove, company)
		}
	}
	for _, toRemove := range companiesToRemove {
		updateQuery := `DELETE FROM company WHERE id=$1`
		manager.executeInsert(updateQuery, toRemove.ID)
	}
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
	manager.queryTable(queryString, func(rows *sql.Rows) error {
		g := new(game.Game)
		var timeStamp int64
		err := rows.Scan(&g.ID, &g.Title, &timeStamp, &g.Summary, &g.Filename)
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
			WHERE gameID=$1`, g.ID)
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
	insertStatement := `INSERT INTO game (id, title, releaseDate, summary, filename)
		VALUES ($1, $2, $3, $4, $5)`
	manager.executeInsert(insertStatement, g.ID, g.Title,
		g.ReleaseDate.Unix(), g.Summary, g.Filename)
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
	file, err := os.Create(dbFile)
	if err != nil {
		log.Fatalf("Failed to create database file: %v\n", err)
	}
	file.Close()
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open database file: %v\n", err)
	}
	manager.db = db
	manager.initializeTables()
}

func (manager *sqliteDBManager) initializeTables() {
	manager.createTable(
		`CREATE TABLE game (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		releaseDate INTEGER NOT NULL,
		summary TEXT NOT NULL,
    	filename TEXT NOT NULL UNIQUE 
		);`)
	manager.createTable(
		`CREATE TABLE company (
		id INTEGER NOT NULL ,
		gameID INTEGER NOT NULL,
		name Text NOT NULL,
		publisher INTEGER(1),
		developer INTEGER(1),
		FOREIGN KEY (gameID) REFERENCES game(ID),
		PRIMARY KEY (id, gameID)
		);`)
	manager.createTable(
		`CREATE TABLE genre (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
		);`)
	manager.createTable(
		`CREATE TABLE artwork (
		id INTEGER PRIMARY KEY,
		remoteURL TEXT NOT NULL,
		gameID INTEGER NOT NULL,
		FOREIGN KEY (gameID) REFERENCES game(id)
		);`)
	manager.createTable(
		`CREATE TABLE genreAssociation (
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
