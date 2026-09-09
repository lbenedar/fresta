package json

import (
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

type Module struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	Authors           []*Author     `json:"authors"`
	URL               string        `json:"url,omitempty"`
	License           string        `json:"license,omitempty"`
	Readme            string        `json:"readme,omitempty"`
	Bugs              string        `json:"bugs,omitempty"`
	Changelog         string        `json:"changelog,omitempty"`
	Media             []*Media      `json:"media"`
	Version           string        `json:"version"`
	Compatibility     Compatibility `json:"compatibility,omitempty"`
	Scripts           []string      `json:"scripts"`
	Esmodules         []string      `json:"esmodules"`
	Styles            []*Style      `json:"styles"`
	Languages         []*Language   `json:"languages"`
	Packs             []*Pack       `json:"packs"`
	PackFolders       []*Folder     `json:"packFolders"`
	Relationships     Relationships `json:"relationships"`
	Socket            bool          `json:"socket"`
	Manifest          string        `json:"manifest,omitempty"`
	Download          string        `json:"download,omitempty"`
	Protected         bool          `json:"protected"`
	Exclusive         bool          `json:"exclusive"`
	PersistentStorage bool          `json:"persistentStorage"`
	CoreTranslation   bool          `json:"coreTranslation"`
	Library           bool          `json:"library"`
	DocumentTypes     DocumentTypes `json:"documentTypes"`
	Availability      int           `json:"availability"`
	Locked            bool          `json:"locked"`
	Owned             bool          `json:"owned"`
	Tags              []string      `json:"tags"`
	HasStorage        bool          `json:"hasStorage"`
	Active            bool          `json:"active,omitempty"`
	// ModulesFlags           ModulesFlags           `json:"flags,omitempty"`
}

func (m *Module) ToDB(dest **db.Module) bool {
	if dest == nil {
		return false
	}

	module := &db.Module{
		ID:                m.ID,
		Title:             m.Title,
		Description:       m.Description,
		URL:               m.URL,
		License:           m.License,
		Readme:            m.Readme,
		Bugs:              m.Bugs,
		Changelog:         m.Changelog,
		Version:           m.Version,
		Socket:            m.Socket,
		Manifest:          m.Manifest,
		Download:          m.Download,
		Protected:         m.Protected,
		Exclusive:         m.Exclusive,
		PersistentStorage: m.PersistentStorage,
		CoreTranslation:   m.CoreTranslation,
		Library:           m.Library,
		Availability:      m.Availability,
		Locked:            m.Locked,
		Owned:             m.Owned,
		HasStorage:        m.HasStorage,
		Active:            m.Active,
	}

	module.Scripts = make([]string, len(m.Scripts))
	copy(module.Scripts, m.Scripts)
	module.Esmodules = make([]string, len(m.Esmodules))
	copy(module.Esmodules, m.Esmodules)
	module.Tags = make([]string, len(m.Tags))
	copy(module.Tags, m.Tags)

	m.Compatibility.ToDB(&module.Compatibility)
	m.Relationships.ToDB(&module.Relationships)
	m.DocumentTypes.ToDB(&module.DocumentTypes)

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &module.Authors, m.Authors)
	CopySliceToDBParallel(&wg, &module.Media, m.Media)
	CopySliceToDBParallel(&wg, &module.Styles, m.Styles)
	CopySliceToDBParallel(&wg, &module.Languages, m.Languages)
	CopySliceToDBParallel(&wg, &module.Packs, m.Packs)
	CopySliceToDBParallel(&wg, &module.PackFolders, m.PackFolders)
	wg.Wait()

	*dest = module

	return true
}
