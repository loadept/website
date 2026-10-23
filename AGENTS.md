# AGENTS.md

Personal site loadept.com: Go backend (chi + zombiezen/sqlite) that server-renders HTML with templ, plus a URL shortener (`/s/{code}`) and Go vanity-import pages (`/go/{pkg}`). README is in Spanish.

The frontend is no longer a static-site generator. It was migrated from Astro (`ui/`, deleted) to templ (`web/`). UI markup lives in `templ` components; page data (software list, links) is hardcoded in `internal/website/handler.go`.

## Critical: build artifacts must exist before `go build`

`go build ./...` / `go test ./...` fail on a fresh clone until two generated artifacts are present:

1. `web/template/*_templ.go` — generated from `*.templ`. The templ CLI is a Go tool dependency, so no separate install is needed:
   ```sh
   go tool templ generate
   ```
   The `*_templ.go` files are checked in; re-run this after editing any `*.templ`.
2. `web/static/css/global.min.css` — Tailwind v4 build output, embedded by `web/static/embed.go` (`//go:embed css/global.min.css fonts images`):
   ```sh
   bun install        # root lockfile is bun.lock
   bun run css:min    # -> web/static/css/global.min.css  (watch: bun run css:dev)
   ```

Static assets (`fonts/`, `images/`) and the minified CSS are embedded at compile time via `web/static/embed.go`, a `neuteredFS` wrapper that refuses directory listings.

## Layout

- `server.go` (root): the web server. Loads env config (`config.Load()`; exits fatally on error), opens a SQLite pool (WAL, foreign keys, busy_timeout), and wires chi routes:
  - `/` -> `internal/website` (renders `template.Index`)
  - `/s/{code}` -> `internal/short` redirect
  - `/go/{pkg}` -> `internal/gopackages` (serves a `go-import` meta tag + redirects to pkg.go.dev)
  - `/static/*` -> file server over the embedded `web/static` FS
- `cmd/migrate`: applies the SQLite schema (inline in its `main.go`, no migration files). Run: `go run ./cmd/migrate -db "$DB_PATH" -run`.
- `cmd/shortener`: CLI to insert short URLs: `-name -url [-code]` (base62 code, max 10 chars if custom).
- `internal/`:
  - `config`: **env-only** loader (godotenv + `os.LookupEnv`); no YAML. Vars: `ADDR` (default `127.0.0.1:8080`), `APP_URL` (default `http://localhost:8080`), `DB_PATH` (required), `GITHUB_URL` (required, builds the go-import URLs). Has tests.
  - `middleware`: request logger + client-IP extraction (CF-Connecting-IP / X-Real-IP / X-Forwarded-For). Has tests.
  - `short`: redirect handler + its SQLite repository (in-memory cache + singleflight).
  - `website`: index-page handler (hardcoded software/links data).
  - `gopackages`: vanity-import handler.
- `web/template`: `*.templ` sources + generated `*_templ.go`.
- `web/static`: embedded assets + `embed.go`.
- `.agents/skills/`: agent skills used by AI tooling; versioned in the repo.

## Verify

```sh
gofumpt -l .              # formatter used by CI (format before pushing)
go tool templ generate    # re-run after editing *.templ
go test ./internal/...    # CI runs only ./internal/...
```

CI runs tests via `gotestsum --format testdox ./internal/...`; gotestsum is not installed locally by default, so plain `go test` is enough.

Local `docker build` needs `--network=host` (bridge MTU truncates registry tarballs behind VPN; integrity check failures on native packages).

## CI / deploy

Pushing to `master` triggers a chained pipeline: Test Project -> Build & Push Docker image (ghcr.io) -> SSH deploy to production (`docker compose up -d`). Merging to master = production deploy.

## Known breakage (from the Astro -> templ migration, not yet committed)

- `Dockerfile` still builds the deleted Astro app: it copies `ui/package.json`, runs `bun run build`, and copies `ui/dist`. It must be rewritten to run `go tool templ generate` and `bun run css:min` instead of the Astro stage.
- `.github/dependabot.yml` still tracks `package-ecosystem: bun` at `directory: /ui` (should be `/`).
- `.github/workflows/test.yml` pins Go `1.26.0` while `go.mod` declares `go 1.26.3` (works via toolchain auto-download, but the pins are inconsistent).
- `README.md` still says the site is built with Astro.
- `.gitignore` still carries Astro-era rules (`ui/dist`, `.astro/`). Decide whether `web/static/css/global.min.css` is committed (it currently is untracked) — if it stays untracked, every fresh clone must run `bun run css:min` before `go build`.

## Conventions

- Commit messages: short, lowercase, no conventional-commit prefixes (e.g. `fix go version on dockerfile`).
- Go toolchain pinned to 1.26.x (go.mod: `go 1.26.3`).
- templ CLI is a Go tool dependency (`tool github.com/a-h/templ/cmd/templ` in go.mod) — use `go tool templ generate`, not a global install.
- Bun is only used for Tailwind CSS now: root `package.json` holds `css:dev` / `css:min`, lockfile is `bun.lock`.
- Tailwind v4 (no `tailwind.config.js`); theme tokens are CSS custom properties in `web/static/css/global.css` under `@theme inline`. Dark mode toggles the `.dark` class on `<html>`.
