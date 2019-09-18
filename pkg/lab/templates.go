package lab

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file"
	"hash/adler32"
	"net"
	"path/filepath"
	"strings"
	"text/template"
)

func DiscoverTemplates(box *packr.Box, prefix string, t *template.Template) (*template.Template, error) {
	err := box.Walk(func(path string, f file.File) error {
		var err error
		var extension = filepath.Ext(path)
		if strings.HasSuffix(path, ".tmpl") {
			path = filepath.Join(prefix, path)
			var path = path[0 : len(path)-len(extension)]
			t, err = t.New(path).Parse(f.String())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return t, err
}

func NewTxtTemplate() *template.Template {
	t := template.New("")
	t = t.Funcs(sprig.TxtFuncMap())
	t = t.Funcs(template.FuncMap{
		"include": func(name string, data interface{}) (string, error) {
			buf := &bytes.Buffer{}
			err := t.ExecuteTemplate(buf, name, data)
			return buf.String(), err
		},
		"NetCidrMask": func(cidr string) (string, error) {
			_, n, err := net.ParseCIDR(cidr)
			mask := ""
			if err == nil {
				mask = fmt.Sprintf("%d.%d.%d.%d",
					n.Mask[0], n.Mask[1], n.Mask[2], n.Mask[3])
			}
			return mask, err
		},
		"NameHash": func(input string) (string, error) {
			c := adler32.Checksum([]byte(input))
			h := (c & 0xffff) + ((c >> 16) & 0xffff)
			return fmt.Sprintf("%x", h), nil
		},
	})
	return t
}
