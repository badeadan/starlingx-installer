package lab

import (
	"path/filepath"
	"strings"
	"text/template"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file"
)


func DiscoverTemplates(box *packr.Box, prefix string, t *template.Template) (*template.Template, error) {
	err := box.Walk(func(path string, f file.File) error {
		var err error
		var extension = filepath.Ext(path)
		if strings.HasSuffix(path, ".tmpl") {
			path = filepath.Join(prefix, path)
			var path = path[0:len(path)-len(extension)]
			t, err = t.New(path).Parse(f.String())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return t, err
}
