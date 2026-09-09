package json

import (
	"encoding/json"
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

type Game struct {
	UserID          string                         `json:"userId"`
	Release         Release                        `json:"release"`
	World           World                          `json:"world"`
	System          System                         `json:"system"`
	Modules         []*Module                      `json:"modules"`
	DemoMode        bool                           `json:"demoMode"`
	IdleLogout      bool                           `json:"idleLogout"`
	Addresses       Addresses                      `json:"addresses"`
	Files           Files                          `json:"files"`
	Options         GameOptions                    `json:"options"`
	ActiveUsers     []string                       `json:"activeUsers"`
	Paused          bool                           `json:"paused"`
	PackageWarnings map[string]PackageWarningsData `json:"packageWarnings"`
	Packs           []*Pack                        `json:"packs"`
	Messages        []*Message                     `json:"messages,omitempty"`
	Combats         []*Combat                      `json:"combats"`
	CardDeck        []*CardDeck                    `json:"cards"`
	Users           []*User                        `json:"users"`
	Macros          []*Macro                       `json:"macros"`
	Folders         []*WorldFolder                 `json:"folders"`
	Items           []*Item                        `json:"items"`
	Settings        []*Setting                     `json:"settings"`
	Journals        []*Journal                     `json:"journal"`
	Tables          []*Table                       `json:"tables"`
	Playlists       []*Playlist                    `json:"playlists"`
	Scenes          []*Scene                       `json:"scenes"`
	Actors          []*Actor                       `json:"actors"`
	CoreUpdate      CoreUpdate                     `json:"coreUpdate"`
	SystemUpdate    SystemUpdate                   `json:"systemUpdate"`
	// Model       Model     `json:"model"`
	// Template        Template                       `json:"template"`
}

func ParseGame(data []byte) (*Game, error) {
	var modelGame []Game
	err := json.Unmarshal(data, &modelGame)
	if err != nil {
		return nil, err
	}

	if len(modelGame) > 1 {
		return nil, ErrorSetupMoreThanOne
	}
	return &modelGame[0], nil
}

func (g *Game) ToDB(dest *db.Game) bool {
	dest.UserID = g.UserID
	dest.DemoMode = g.DemoMode
	dest.IdleLogout = g.IdleLogout
	dest.Paused = g.Paused
	PackageWarningsToDB(&dest.PackageWarnings, g.PackageWarnings)

	g.Addresses.ToDB(&dest.Addresses)
	g.Files.ToDB(&dest.Files)
	g.Release.ToDB(&dest.Release)
	g.World.ToDB(&dest.World)
	g.System.ToDB(&dest.System)
	g.Options.ToDB(&dest.Options)
	g.CoreUpdate.ToDB(&dest.CoreUpdate)
	g.SystemUpdate.ToDB(&dest.SystemUpdate)

	dest.ActiveUsers = make([]string, len(g.ActiveUsers))
	copy(dest.ActiveUsers, g.ActiveUsers)

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &dest.Modules, g.Modules)
	CopySliceToDBParallel(&wg, &dest.Packs, g.Packs)
	CopySliceToDBParallel(&wg, &dest.Messages, g.Messages)
	CopySliceToDBParallel(&wg, &dest.Combats, g.Combats)
	CopySliceToDBParallel(&wg, &dest.CardDeck, g.CardDeck)
	CopySliceToDBParallel(&wg, &dest.Users, g.Users)
	CopySliceToDBParallel(&wg, &dest.Macros, g.Macros)
	CopySliceToDBParallel(&wg, &dest.Folders, g.Folders)
	CopySliceToDBParallel(&wg, &dest.Items, g.Items)
	CopySliceToDBParallel(&wg, &dest.Settings, g.Settings)
	CopySliceToDBParallel(&wg, &dest.Journals, g.Journals)
	CopySliceToDBParallel(&wg, &dest.Tables, g.Tables)
	CopySliceToDBParallel(&wg, &dest.Playlists, g.Playlists)
	CopySliceToDBParallel(&wg, &dest.Actors, g.Actors)
	// go CopySliceToDBParallel(&wg, &dest.Scenes, g.Scenes)
	wg.Wait()

	return true
}
