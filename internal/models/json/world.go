package json

import (
	"sync"
	"time"

	"github.com/lbenedar/fresta/internal/models/db"
)

type World struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	Authors           []*Author     `json:"authors"`
	Media             []*Media      `json:"media"`
	Version           string        `json:"version"`
	Compatibility     Compatibility `json:"compatibility"`
	Scripts           []string      `json:"scripts"`
	Esmodules         []string      `json:"esmodules"`
	Styles            []*Style      `json:"styles"`
	Languages         []*Language   `json:"languages"`
	Packs             []*Pack       `json:"packs"`
	PackFolders       []*Folder     `json:"packFolders"`
	Relationships     Relationships `json:"relationships"`
	Socket            bool          `json:"socket"`
	Protected         bool          `json:"protected"`
	Exclusive         bool          `json:"exclusive"`
	PersistentStorage bool          `json:"persistentStorage"`
	System            string        `json:"system"`
	Background        string        `json:"background"`
	JoinTheme         string        `json:"joinTheme"`
	CoreVersion       string        `json:"coreVersion"`
	SystemVersion     string        `json:"systemVersion"`
	LastPlayed        string        `json:"lastPlayed"`
	Playtime          int           `json:"playtime"`
	NextSession       time.Time     `json:"nextSession"`
	Demo              Demo          `json:"demo"`
	Availability      int           `json:"availability"`
	Locked            bool          `json:"locked"`
	Owned             bool          `json:"owned"`
	Tags              []string      `json:"tags"`
	HasStorage        bool          `json:"hasStorage"`
	// Flags        any    `json:"flags"`
}

func (w *World) ToDB(dest **db.World) bool {
	if dest == nil {
		return false
	}

	world := &db.World{
		ID:                w.ID,
		Title:             w.Title,
		Description:       w.Description,
		Version:           w.Version,
		Socket:            w.Socket,
		Protected:         w.Protected,
		Exclusive:         w.Exclusive,
		PersistentStorage: w.PersistentStorage,
		System:            w.System,
		Background:        w.Background,
		JoinTheme:         w.JoinTheme,
		CoreVersion:       w.CoreVersion,
		SystemVersion:     w.SystemVersion,
		LastPlayed:        w.LastPlayed,
		Playtime:          w.Playtime,
		NextSession:       w.NextSession,
		Availability:      w.Availability,
		Locked:            w.Locked,
		Owned:             w.Owned,
		HasStorage:        w.HasStorage,
	}

	world.Scripts = make([]string, len(w.Scripts))
	copy(world.Scripts, w.Scripts)
	world.Esmodules = make([]string, len(w.Esmodules))
	copy(world.Esmodules, w.Esmodules)
	world.Tags = make([]string, len(w.Tags))
	copy(world.Tags, w.Tags)

	w.Compatibility.ToDB(&world.Compatibility)
	w.Relationships.ToDB(&world.Relationships)

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &world.Authors, w.Authors)
	CopySliceToDBParallel(&wg, &world.Media, w.Media)
	CopySliceToDBParallel(&wg, &world.Styles, w.Styles)
	CopySliceToDBParallel(&wg, &world.Languages, w.Languages)
	CopySliceToDBParallel(&wg, &world.Packs, w.Packs)
	CopySliceToDBParallel(&wg, &world.PackFolders, w.PackFolders)
	wg.Wait()

	*dest = world

	return true
}

type Demo struct {
	SourceZip any `json:"sourceZip"`
}
