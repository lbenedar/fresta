package json

import "github.com/lbenedar/fresta/internal/models/db"

type Files struct {
	Storages []string `json:"storages"`
	S3       any      `json:"s3"`
}

func (f *Files) ToDB(dest *db.Files) bool {
	if dest == nil {
		return false
	}

	dest.Storages = make(db.FilesStorage, len(f.Storages))
	for i := range f.Storages {
		dest.Storages[i].Storage = f.Storages[i]
	}

	return true
}
