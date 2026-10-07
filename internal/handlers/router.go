// Package handlers is the HTTP layer: the JSON API under /api, and the frontend everywhere else
// (docs/adr/0001). It depends on the interfaces below, not on concrete clients.
package handlers

import (
	"context"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kthaisociety/growth-portal/internal/demographics"
	"github.com/kthaisociety/growth-portal/internal/lumaservice"
)

// Members is what the handlers need from luma-service. *lumaservice.Client implements it.
type Members interface {
	CommunityMembers(ctx context.Context) ([]lumaservice.Member, error)
}

type Deps struct {
	Members Members
	PingDB  func(ctx context.Context) error
	// UI is the built frontend (ui.Dist()).
	UI fs.FS
}

func NewRouter(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api")
	api.GET("/health", health(d.PingDB))
	api.GET("/community/demographics", communityDemographics(d.Members))

	r.NoRoute(frontend(d.UI))
	return r
}

func health(pingDB func(ctx context.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := pingDB(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"service": "growth-portal", "status": "unhealthy", "database": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"service": "growth-portal", "status": "healthy"})
	}
}

// communityDemographics returns aggregate counts only, never individual members (docs/pcr/0005).
func communityDemographics(m Members) gin.HandlerFunc {
	return func(c *gin.Context) {
		members, err := m.CommunityMembers(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "loading members failed"})
			return
		}
		c.JSON(http.StatusOK, demographics.Summarize(members))
	}
}

// frontend serves the single-page app: a file if one exists at the path, otherwise index.html so
// client-side routes work on reload. Unknown /api paths stay JSON 404s.
func frontend(ui fs.FS) gin.HandlerFunc {
	files := http.FileServerFS(ui)
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if p == "/api" || strings.HasPrefix(p, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(path.Clean(p), "/")
		if name != "" && name != "index.html" {
			if info, err := fs.Stat(ui, name); err == nil && !info.IsDir() {
				if strings.HasPrefix(name, "assets/") {
					// Vite puts a content hash in every asset's name.
					c.Header("Cache-Control", "public, max-age=31536000, immutable")
				}
				files.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		index, err := fs.ReadFile(ui, "index.html")
		if err != nil {
			c.String(http.StatusServiceUnavailable, "frontend not built: run `npm run build` in ui/")
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	}
}
