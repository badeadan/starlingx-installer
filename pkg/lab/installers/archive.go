package installers

import (
	"archive/tar"
	"bytes"
	"time"
)

type TarWriter struct {
	*tar.Writer
}

func (tw *TarWriter) WriteFileBytes(name string, mode int64, buffer *bytes.Buffer) error {
	delta, _ := time.ParseDuration("-30s")
	modtime := time.Now().Add(delta)
	err := tw.WriteHeader(&tar.Header{
		Name:    name,
		Mode:    mode,
		Size:    int64(buffer.Len()),
		ModTime: modtime,
	})
	if err == nil {
		_, err = tw.Write(buffer.Bytes())
	}
	return err
}

