package lab

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func DiscoverTemplates(t *template.Template) (*template.Template, error) {
	root := filepath.Clean("./templates")
	prefix := len(root) + 1
	err := filepath.Walk(root, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
			if e1 != nil {
				return e1
			}
			buffer, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			name := path[prefix:]
			t, err = t.New(name).Parse(string(buffer))
			if err != nil {
				return err
			}
		}

		return nil
	})

	return t, err
}
