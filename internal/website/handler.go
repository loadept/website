package website

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/loadept/website/web/template"
)

var software = []template.Links{
	{
		Label: "wirq",
		Text:  "Proxy MITM para inspeccionar el trafico HTTP/HTTPS escrito en go.",
		URL:   "https://github.com/loadept/wirq",
	},
	{
		Label: "ipseek",
		Text:  "API gratuita para obtener la información y localización de una IP.",
		URL:   "https://ipseek.space",
	},
	{
		Label: "db-mcp",
		Text:  "Servidor MCP para multiples bases de datos (solo lectura).",
		URL:   "https://loadept.com/s/uEZQ",
	},
	{
		Label: "playground",
		Text:  "Repo con mis expermientos de código, pruebas, prototipos en distintos lenguajes.",
		URL:   "https://github.com/loadept/playground",
	},
}

var links = []template.Links{
	{
		Label: "Email",
		Text:  "hi@loadept.com",
		URL:   "mailto:hi@loadept.com?subject=Hi Loadept!",
	},
	{
		Label: "Github",
		Text:  "Loadept",
		URL:   "https://loadept.com/s/gh",
	},
	{
		Label: "Linkedin",
		Text:  "Jesus Gonzalo Machaca Hancco",
		URL:   "https://loadept.com/s/linkedin",
	},
	{
		Label: "Tiktok",
		Text:  "@loadept",
		URL:   "https://loadept.com/s/tiktok",
	},
	{
		Label: "Blog",
		Text:  "blog.loadept.com",
		URL:   "https://blog.loadept.com",
	},
}

var (
	title       = "loadept"
	description = "Mi sitio web personal, comparto recursos relacionados con el desarrollo de software :) ."
)

func NewHandler(appURL *url.URL, githubURL *url.URL) *Handler {
	return new(Handler{appURL: appURL, githubURL: githubURL})
}

type Handler struct {
	appURL    *url.URL
	githubURL *url.URL
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	props := template.Props{
		Title:       title,
		Description: description,
		Image:       h.appURL.JoinPath("static", "images", "draw.webp").String(),
		Canonical:   h.appURL.String(),
		Software:    software,
		Links:       links,
	}
	template.Index(props).Render(r.Context(), w)
}

var vanityPackages = map[string]bool{
	"asyncker": true,
}

func (h *Handler) ServePkg(w http.ResponseWriter, r *http.Request) {
	pkg := chi.URLParam(r, "pkg")

	if !vanityPackages[pkg] {
		http.NotFound(w, r)
		return
	}

	dom := h.appURL.Hostname()
	ghurl := h.githubURL.String()
	template.Gopkg(dom, ghurl, pkg).Render(r.Context(), w)
}
