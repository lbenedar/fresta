package json

import (
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

type System struct {
	ID                    string        `json:"id"`
	Title                 string        `json:"title"`
	Description           string        `json:"description"`
	Authors               []*Author     `json:"authors"`
	URL                   string        `json:"url"`
	License               string        `json:"license"`
	Bugs                  string        `json:"bugs"`
	Changelog             string        `json:"changelog"`
	Media                 []*Media      `json:"media"`
	Version               string        `json:"version"`
	Compatibility         Compatibility `json:"compatibility"`
	Scripts               []string      `json:"scripts"`
	Esmodules             []string      `json:"esmodules"`
	Packs                 []*Pack       `json:"packs"`
	Styles                []*Style      `json:"styles"`
	Languages             []*Language   `json:"languages"`
	PackFolders           []*Folder     `json:"packFolders"`
	Relationships         Relationships `json:"relationships"`
	Socket                bool          `json:"socket"`
	Manifest              string        `json:"manifest"`
	Download              string        `json:"download"`
	Protected             bool          `json:"protected"`
	Exclusive             bool          `json:"exclusive"`
	PersistentStorage     bool          `json:"persistentStorage"`
	DocumentTypes         DocumentTypes `json:"documentTypes"`
	Background            string        `json:"background"`
	Grid                  Grid          `json:"grid"`
	PrimaryTokenAttribute string        `json:"primaryTokenAttribute"`
	Availability          int           `json:"availability"`
	Locked                bool          `json:"locked"`
	Owned                 bool          `json:"owned"`
	Tags                  []string      `json:"tags"`
	HasStorage            bool          `json:"hasStorage"`
	// SystemFlags           SystemFlags         `json:"flags"`
}

func (s *System) ToDB(dest **db.System) bool {
	if dest == nil {
		return false
	}

	system := &db.System{
		ID:                    s.ID,
		Title:                 s.Title,
		Description:           s.Description,
		URL:                   s.URL,
		License:               s.License,
		Bugs:                  s.Bugs,
		Changelog:             s.Changelog,
		Version:               s.Version,
		Socket:                s.Socket,
		Manifest:              s.Manifest,
		Download:              s.Download,
		Protected:             s.Protected,
		Exclusive:             s.Exclusive,
		PersistentStorage:     s.PersistentStorage,
		Background:            s.Background,
		PrimaryTokenAttribute: s.PrimaryTokenAttribute,
		Availability:          s.Availability,
		Locked:                s.Locked,
		Owned:                 s.Owned,
		HasStorage:            s.HasStorage,
	}

	system.Scripts = make([]string, len(s.Scripts))
	copy(system.Scripts, s.Scripts)
	system.Esmodules = make([]string, len(s.Esmodules))
	copy(system.Esmodules, s.Esmodules)
	system.Tags = make([]string, len(s.Tags))
	copy(system.Tags, s.Tags)

	s.Compatibility.ToDB(&system.Compatibility)
	s.Relationships.ToDB(&system.Relationships)
	s.DocumentTypes.ToDB(&system.DocumentTypes)
	s.Grid.ToDB(&system.Grid)

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &system.Authors, s.Authors)
	CopySliceToDBParallel(&wg, &system.Media, s.Media)
	CopySliceToDBParallel(&wg, &system.Styles, s.Styles)
	CopySliceToDBParallel(&wg, &system.Languages, s.Languages)
	CopySliceToDBParallel(&wg, &system.Packs, s.Packs)
	CopySliceToDBParallel(&wg, &system.PackFolders, s.PackFolders)
	wg.Wait()

	*dest = system

	return true
}
