package json

import "github.com/lbenedar/fresta/internal/models/db"

type Folder struct {
	Name    string    `json:"name"`
	Sorting string    `json:"sorting"`
	Color   string    `json:"color,omitempty"`
	Packs   []string  `json:"packs"`
	Folders []*Folder `json:"folders,omitempty"`
}

func (f *Folder) ToDB(dest **db.Folder) bool {
	if dest == nil {
		return false
	}

	folder := &db.Folder{
		Name:    f.Name,
		Sorting: f.Sorting,
		Color:   f.Color,
	}

	folder.Packs = make([]string, len(f.Packs))
	copy(folder.Packs, f.Packs)
	CopySliceToDB(&folder.Folders, f.Folders)

	*dest = folder

	return true
}

type WorldFolder struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ID          string `json:"_id"`
	Folder      string `json:"folder"`
	Sorting     string `json:"sorting"`
	Sort        int    `json:"sort"`
	Stats       Stats  `json:"_stats,omitempty"`
	Description string `json:"description"`
	Color       string `json:"color"`
	// Flags       any    `json:"flags,omitempty"`
}

func (w *WorldFolder) ToDB(dest **db.WorldFolder) bool {
	if dest == nil {
		return false
	}

	worldFolder := &db.WorldFolder{
		Name:        w.Name,
		Type:        w.Type,
		ID:          w.ID,
		Folder:      w.Folder,
		Sorting:     w.Sorting,
		Sort:        w.Sort,
		Description: w.Description,
		Color:       w.Color,
	}

	w.Stats.ToDB(&worldFolder.Stats)

	*dest = worldFolder

	return true
}
