// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mdhttp

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"sync/atomic"

	"github.com/Depado/bfchroma"
	"github.com/microcosm-cc/bluemonday"
	bf "gopkg.in/russross/blackfriday.v2"
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

checkdir:

	if mdf.FileInfo, err = f.Stat(); err != nil {
		http.Error(w, fmt.Sprintf("error stating %q: %s", mdf.Path, err), http.StatusInternalServerError)
		return
	}

	if mdf.FileInfo.IsDir() {
		// Check to see if there is an "index.md" file in the directory, and
		// load it if so.
		ipath := strings.TrimSuffix(mdf.Path, "/") + "/index.md"
		if f, err = md.fs.Open(ipath); err == nil {
			mdf.Path = ipath
			goto checkdir
		}

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

	tocCache  atomic.Value
	bodyCache atomic.Value
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

var tocSplitAt = []byte("</nav>")
var tocReplacer = strings.NewReplacer(
	"<nav>", "",
	"<ul>", `<ul class="nav flex-column">`,
	"<li>", `<li class="nav-item">`,
	"<a ", `<a class="nav-link" `,
)

// HTML returns the rendered HTML version of the markdown file, by default
// sanitized ith UGC policies (bluemonday), and common markdown format
// (blackfriday).
func (mdf *MarkdownFile) HTML() (toc, body string) {
	if tocCache := mdf.tocCache.Load(); tocCache != nil {
		toc = tocCache.(string)
	}
	if bodyCache := mdf.bodyCache.Load(); bodyCache != nil {
		body = bodyCache.(string)
	}

	if toc != "" || body != "" {
		return toc, body
	}

	raw := bf.Run(
		mdf.Content,
		bf.WithRenderer(bfchroma.NewRenderer(
			bfchroma.Style("github"),
			bfchroma.Extend(bf.NewHTMLRenderer(bf.HTMLRendererParameters{
				Flags: bf.CommonHTMLFlags | bf.TOC,
			})),
		)),
	)

	// Allow chroma attributes through.
	policy := bluemonday.UGCPolicy()
	policy = policy.AllowElements("span")
	policy = policy.AllowAttrs("style").OnElements("span")

	if i := bytes.Index(raw, tocSplitAt); i > -1 {
		toc = tocReplacer.Replace(string(raw[:i]))
		body = string(policy.SanitizeBytes(raw[i+len(tocSplitAt)+1:]))

		mdf.tocCache.Store(toc)
		mdf.bodyCache.Store(body)

		return toc, body
	}

	body = string(policy.SanitizeBytes(raw))
	mdf.bodyCache.Store(body)
	return "", body
}

// Body returns a template.HTML wrapped HTML representation of the Markdown
// body.
func (mdf *MarkdownFile) Body() template.HTML {
	_, body := mdf.HTML()
	return template.HTML(body)
}

// TOC returns a template.HTML wrapped HTML representation of the Markdown
// table of contents.
func (mdf *MarkdownFile) TOC() template.HTML {
	toc, _ := mdf.HTML()
	return template.HTML(toc)
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
