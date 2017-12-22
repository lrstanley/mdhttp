// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mdhttp

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// New returns a new markdown renderer. Prefix is url path to strip from the
// URL. fs is the filesystem you want to use, and renderer is an optional
// custom renderer (default uses DefaultRenderer). New() differs from
// NewMiddleware() in that New() only allows rendering markdown files, and
// will not serve other static files.
func New(prefix string, fs http.FileSystem, renderer Renderer) http.Handler {
	md := &Markdown{prefix: prefix, fs: fs, renderer: renderer}
	md.handler = http.StripPrefix(prefix, http.FileServer(fs))
	return md
}

// NewMiddleware is like New() however it allows serving non-markdown files.
// Note that because this is a middleware, if the middleware encounters
// a markdown file, it doesn't continue executing the middleware chain, and
// will return early. Make sure to have this be the last middleware in the
// chain to have all of middleware ahead of this get executed properly for
// pages with markdown content on them.
func NewMiddleware(prefix string, fs http.FileSystem, renderer Renderer) func(next http.Handler) http.Handler {
	md := &Markdown{prefix: prefix, fs: fs, renderer: renderer}
	md.handler = http.StripPrefix(prefix, http.FileServer(fs))
	return md.Register
}

// Renderer is the render function type which is used by the Markdown
// middleware. This allows you to customize how the markdown is rendered (e.g.
// you could use your own template wrapping, or custom styles). See the source
// code for DefaultRenderer for how you would implement this.
type Renderer func(w http.ResponseWriter, r *http.Request, mdr *MarkdownFile)

// Markdown is an instantiated handler.
type Markdown struct {
	prefix   string
	renderer func(w http.ResponseWriter, r *http.Request, mdr *MarkdownFile)
	fs       http.FileSystem
	handler  http.Handler
}

// Register is a way of using mdhttp as a middleware that seemlessly renders
// markdown files.
func (md *Markdown) Register(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(strings.ToLower(r.URL.Path), ".md") {
			md.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (md *Markdown) path(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, md.prefix)
}

func (md *Markdown) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mdf := &MarkdownFile{Path: md.path(r)}

	var f http.File
	var err error

	if f, err = md.fs.Open(mdf.Path); err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, fmt.Sprintf("error opening %q: %s", mdf.Path, err), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if mdf.FileInfo, err = f.Stat(); err != nil {
		http.Error(w, fmt.Sprintf("error stating %q: %s", mdf.Path, err), http.StatusInternalServerError)
		return
	}

	if mdf.FileInfo.IsDir() {
		md.handler.ServeHTTP(w, r)
		return
	}

	if !strings.HasSuffix(strings.ToLower(mdf.Path), ".md") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	tp := textproto.NewReader(bufio.NewReader(f))

	hdr, err := tp.ReadMIMEHeader()
	if err != nil {
		if _, ok := err.(textproto.ProtocolError); !ok {
			http.Error(w, fmt.Sprintf("error reading %q: %s", mdf.Path, err), http.StatusInternalServerError)
			return
		}

		// Skip, they likely aren't using custom headers/attributes.
		_, _ = f.Seek(0, 0)
		mdf.Content, err = ioutil.ReadAll(f)
	} else {
		mdf.attr = hdr
		mdf.Content, err = ioutil.ReadAll(tp.R)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error reading %q: %s", mdf.Path, err), http.StatusInternalServerError)
		return
	}

	// Make a name based on the filename
	mdf.Title = strings.Title(titleReplacer.Replace(strings.TrimSuffix(mdf.FileInfo.Name(), ".md")))

	if title := mdf.GetAttr("Title"); title != "" {
		mdf.Title = title
	}

	if md.renderer == nil {
		DefaultRenderer(w, r, mdf)
	}
}

var titleReplacer = strings.NewReplacer(
	"_", " ",
	"-", " ",
	".", " ",
)

// MarkdownFile contains all of the necessary components to render the markdown
// file.
type MarkdownFile struct {
	Title    string
	Path     string
	FileInfo os.FileInfo
	Content  []byte
	attr     textproto.MIMEHeader
}

// GetAttr gets one of the attributes (metadata) which was defined at the top
// of the markdown file. For example, you can pull "Created" or "Title" from
// the following example:
//
//   Created: <some date>
//   Title: Example Page
//
//   ## Your Markdown Here
func (mdf *MarkdownFile) GetAttr(name string) string {
	if mdf.attr == nil {
		return ""
	}

	return mdf.attr.Get(name)
}

// HTML returns the rendered HTML version of the markdown file, by default
// sanitized ith UGC policies (bluemonday), and common markdown format
// (blackfriday).
func (mdf *MarkdownFile) HTML() []byte {
	return bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(mdf.Content))
}

// HTMLTemplate is the same as the HTML method, however it wraps the returned
// bytes in template.HTML(), for use with injecting in an html/template template.
func (mdf *MarkdownFile) HTMLTemplate() template.HTML {
	return template.HTML(mdf.HTML())
}

// DefaultRenderer satisfies the Renderer type, and can be used with
// New() and NewMiddleware().
func DefaultRenderer(w http.ResponseWriter, r *http.Request, mdf *MarkdownFile) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	t, err := template.New("").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, mdf)
	if err != nil {
		log.Printf("error executing template: %s", err)
	}
}
