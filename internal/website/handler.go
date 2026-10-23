package website

import (
	"net/http"
	"strings"

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
		Label: "db-mcp",
		Text:  "Servidor MCP para multiples bases de datos (solo lectura).",
		URL:   "https://loadept.com/s/uEZQ",
	},
	{
		Label: "cursor-live",
		Text:  "Mira cómo funcionan los websockets con este experimento.",
		URL:   "https://live.loadept.com",
	},
	{
		Label: "ipseek",
		Text:  "API gratuita para obtener la información y localización de una IP.",
		URL:   "https://ipseek.space",
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

func NewHandler(appURL string) *Handler {
	return new(Handler{appURL: appURL})
}

type Handler struct {
	appURL    string
	githubURL string
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	appURL := strings.TrimPrefix(h.appURL, "/")

	props := template.Props{
		Title:       title,
		Description: description,
		Image:       appURL + "/static/images/draw.webp",
		Canonical:   appURL,
		Software:    software,
		Links:       links,
	}
	template.Index(props).Render(r.Context(), w)
}

func (p *Handler) ServePkg(w http.ResponseWriter, r *http.Request) {
	pkg := chi.URLParam(r, "pkg")

	template.Gopkg(r.Host, pkg, p.githubURL).Render(r.Context(), w)
}
