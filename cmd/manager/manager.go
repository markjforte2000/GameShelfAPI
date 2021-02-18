package main

import (
	"github.com/markjforte2000/GameShelfAPI/internal/manager"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"time"
)

func main() {
	m := manager.NewGameLibManager()
	lib := m.GetGameLibrary()
	for _, g := range lib {
		util.PrettyPrintGame(g)
	}
	for {
		changes := m.GetChanges()
		for _, g := range changes {
			util.PrettyPrintGame(g)
		}
		time.Sleep(time.Second)
	}
}
