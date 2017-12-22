// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mdhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// New returns a new markdown renderer. Prefix is url path to strip from the
// URL. fs is the filesystem you want to use, and renderer is an optional
// custom renderer (default uses DefaultRenderer).
func New(prefix string, fs http.FileSystem, renderer Renderer) http.Handler {
	md := &Markdown{prefix: prefix, fs: fs, renderer: renderer}
	md.handler = http.StripPrefix(prefix, http.FileServer(fs))

	return md
}

// Renderer is the render function type which is used by the Markdown
// middleware. This allows you to customize how the markdown is render (e.g.
// you could use your own template wrapping, or custom styles). See the source
// code for DefaultRenderer for how you would implement this.
type Renderer func(w http.ResponseWriter, r *http.Request, path, title string, md []byte)

type Markdown struct {
	prefix   string
	renderer func(w http.ResponseWriter, r *http.Request, path, title string, md []byte)
	fs       http.FileSystem
	handler  http.Handler
}

func (md *Markdown) path(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, md.prefix)
}

func (md *Markdown) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := md.path(r)

	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		md.handler.ServeHTTP(w, r)
		return
	}

	f, err := md.fs.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, fmt.Sprintf("error opening %q: %s", path, err), http.StatusInternalServerError)
		return
	}

	var b []byte
	b, err = ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading %q: %s", path, err), http.StatusInternalServerError)
		return
	}

	title := strings.Title(strings.Replace(strings.TrimSuffix(filepath.Base(path), ".md"), "-", " ", -1))

	if md.renderer == nil {
		DefaultRenderer(w, r, path, title, b)
	}
}

// DefaultRenderer satisfies the Renderer type, and can be used
// with NewMarkdown().
func DefaultRenderer(w http.ResponseWriter, r *http.Request, path, title string, md []byte) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, htmlTemplate, title, bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(md)))
}
