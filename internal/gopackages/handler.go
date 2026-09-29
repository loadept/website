package gopackages

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var goModuleTmpl = template.Must(template.New("goModuleTmpl").Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="go-import" content="{{ .Domain }}/go/{{ .Package }} git https://github.com/{{ .GithubUser }}/{{ .Package }}" />    
    <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/{{ .Domain }}/go/{{ .Package }}" />      
    <title>{{ .Package }}</title>
  </head>
  <body>
    <p>
      Nothing to see here;
      <a href="https://pkg.go.dev/{{ .Domain }}/go/{{ .Package }}">
        see the package on pkg.go.dev
      </a>.
    </p>
  </body>
</html>
`))

type Handler struct {
	githubUser string
}

func NewHandler(githubUser string) *Handler {
	return new(Handler{githubUser: githubUser})
}

func (p *Handler) ServePkg(w http.ResponseWriter, r *http.Request) {
	domain := r.Host
	pkg := chi.URLParam(r, "pkg")

	goModuleTmpl.Execute(w, map[string]string{
		"Domain":     domain,
		"GithubUser": p.githubUser,
		"Package":    pkg,
	})
}
