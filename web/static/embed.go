package static

import (
	"embed"
	"io/fs"
	"path"
	"strings"
)

//go:embed css/global.min.css fonts images
var staticFS embed.FS

type neuteredFS struct{ fs fs.FS }

func (n *neuteredFS) Open(name string) (fs.File, error) {
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		name = "."
	}

	f, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}

	if stat, err := f.Stat(); err != nil {
		f.Close()
		return nil, err
	} else if stat.IsDir() {
		_, err := fs.Stat(n.fs, path.Join(name, "index.html"))
		if err != nil {
			f.Close()
			return nil, fs.ErrNotExist
		}
	}

	return f, nil
}

var FS fs.FS

func init() {
	FS = &neuteredFS{staticFS}
}
