package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.yogan.dev/internal/models"
	"snippetbox.yogan.dev/ui"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	// pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// ts, err = ts.ParseGlob("./ui/html/partials/*.html"
		// if err != nil {
		// 	return nil, err
		// }

		// ts, err = ts.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }

		cache[name] = ts
	}

	return cache, nil
}
