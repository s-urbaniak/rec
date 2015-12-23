package webapp

import (
	"net/http"
	"os"
)

var _ http.FileSystem = (*filesystem)(nil)

// filesystem wraps an underlying http.FileSystems
// and delegates to a custom http.File.
type filesystem struct {
	underlying http.FileSystem
}

func (fs filesystem) Open(name string) (http.File, error) {
	f, err := fs.underlying.Open(name)
	if err != nil {
		return nil, err
	}

	return file{f}, nil
}

var _ http.File = (*file)(nil)

// file embeds an underyling http.File but overrides Readdir
// to prevent directory listings.
type file struct {
	http.File
}

// Readdir returns no directory information
func (f file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

// FileServerNoReaddir returns a http.Handler that behaves like a http.FileServer
// but doesn't provide directory listings.
func FileServerNoReaddir(dir string) http.Handler {
	return http.FileServer(filesystem{http.Dir(dir)})
}
