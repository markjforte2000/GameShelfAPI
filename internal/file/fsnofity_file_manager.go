package file

import (
	"github.com/fsnotify/fsnotify"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
	"os"
	"path/filepath"
)

// file monitor using fsnotify package
// does not support nested files
type fsnotifyFileManager struct {
	watcher       *fsnotify.Watcher
	rootDirectory string
	handler       NewFileHandler
}

func (manager *fsnotifyFileManager) GetCurrentFiles() []*game.GameFile {
	var files []string
	err := filepath.Walk(manager.rootDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error reading file %v: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileName := filepath.Base(path)
		files = append(files, fileName)
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking root directory: %v\n", err)
	}
	var gameFiles []*game.GameFile
	for _, file := range files {
		gameFile := parseFileName(file)
		if gameFile == nil {
			continue
		}
		logGameFile(gameFile)
		gameFiles = append(gameFiles, gameFile)
	}
	return gameFiles
}

func (manager *fsnotifyFileManager) init(rootDirectory string, handler NewFileHandler) {
	manager.rootDirectory = rootDirectory
	manager.handler = handler
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Unable to create fsnotify file Watcher: %v", err)
	}
	err = watcher.Add(rootDirectory)
	if err != nil {
		log.Fatalf("Failed to add root directory %v to watcher: %v", rootDirectory, err)
	}
	manager.watcher = watcher
	go manager.watchRootDirectory()
}

func (manager *fsnotifyFileManager) watchRootDirectory() {
	for {
		select {
		case event, ok := <-manager.watcher.Events:
			if !ok {
				return
			}
			if event.Op == fsnotify.Write || event.Op == fsnotify.Chmod {
				return
			}
			log.Printf("Detected file change: %v\n", event)
			base := filepath.Base(event.Name)
			gameFile := parseFileName(base)
			if gameFile == nil {
				return
			}
			logGameFile(gameFile)
			manager.handler(gameFile, translateOp(event.Op))
		}
	}
}

func translateOp(op fsnotify.Op) Op {
	switch op {
	case fsnotify.Remove:
		return Delete
	case fsnotify.Create:
		return Create
	case fsnotify.Rename:
		return Rename
	}
	return Null
}
