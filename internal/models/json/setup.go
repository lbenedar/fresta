package json

import (
	"encoding/json"
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

type Setup struct {
	CoreUpdate      CoreUpdate                     `json:"coreUpdate"`
	FeaturedContent FeaturedContent                `json:"featuredContent"`
	Files           Files                          `json:"files"`
	IsAdmin         bool                           `json:"isAdmin"`
	IsSetup         bool                           `json:"isSetup"`
	Languages       []*SetupLanguage               `json:"languages"`
	Modules         []*Module                      `json:"modules"`
	News            []*News                        `json:"news"`
	Options         SetupOptions                   `json:"options"`
	PackageWarnings map[string]PackageWarningsData `json:"packageWarnings"`
	Release         Release                        `json:"release"`
	Systems         []*System                      `json:"systems"`
	Worlds          []*World                       `json:"worlds"`
}

func ParseSetup(data []byte) (*Setup, error) {
	var modelSetup []Setup
	err := json.Unmarshal(data, &modelSetup)
	if err != nil {
		return nil, err
	}

	if len(modelSetup) > 1 {
		return nil, ErrorSetupMoreThanOne
	}
	return &modelSetup[0], nil
}

func (s *Setup) ToDB(dest *db.Setup) bool {
	if dest == nil {
		return false
	}

	dest.IsAdmin = s.IsAdmin
	dest.IsSetup = s.IsSetup

	s.CoreUpdate.ToDB(&dest.CoreUpdate)
	s.FeaturedContent.ToDB(&dest.FeaturedContent)
	s.Files.ToDB(&dest.Files)
	s.Options.ToDB(&dest.Options)
	s.Release.ToDB(&dest.Release)

	PackageWarningsToDB(&dest.PackageWarnings, s.PackageWarnings)

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &dest.Languages, s.Languages)
	CopySliceToDBParallel(&wg, &dest.Modules, s.Modules)
	CopySliceToDBParallel(&wg, &dest.News, s.News)
	CopySliceToDBParallel(&wg, &dest.Systems, s.Systems)
	CopySliceToDBParallel(&wg, &dest.Worlds, s.Worlds)
	wg.Wait()

	return true
}

type FeaturedContent struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
	Image   string `json:"image"`
}

func (f *FeaturedContent) ToDB(dest *db.FeaturedContent) bool {
	if dest == nil {
		return false
	}

	dest.Title = f.Title
	dest.Caption = f.Caption
	dest.URL = f.URL
	dest.Image = f.Image

	return true
}

type News struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
	Image   string `json:"image"`
}

func (n *News) ToDB(dest **db.News) bool {
	if dest == nil {
		return false
	}

	news := &db.News{
		Title:   n.Title,
		Caption: n.Caption,
		URL:     n.Caption,
		Image:   n.Image,
	}

	*dest = news

	return true
}
