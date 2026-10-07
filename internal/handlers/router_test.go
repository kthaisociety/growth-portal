package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kthaisociety/growth-portal/internal/lumaservice"
)

type fakeMembers struct {
	members []lumaservice.Member
	err     error
}

func (f fakeMembers) CommunityMembers(context.Context) ([]lumaservice.Member, error) {
	return f.members, f.err
}

var builtUI = fstest.MapFS{
	"index.html":      {Data: []byte("<!doctype html><title>Growth Portal</title>")},
	"assets/app-1.js": {Data: []byte("console.log(1)")},
	"favicon.svg":     {Data: []byte("<svg/>")},
	".gitkeep":        {Data: nil},
}

func newTestRouter(m Members, ui fstest.MapFS) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(Deps{Members: m, PingDB: func(context.Context) error { return nil }, UI: ui})
}

func get(r http.Handler, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

func TestDemographicsReturnsAggregatesOnly(t *testing.T) {
	r := newTestRouter(fakeMembers{members: []lumaservice.Member{
		{Answers: []lumaservice.Answer{{QuestionID: "year", Label: "Graduation year", Type: "select", Values: []string{"2027"}}}},
	}}, builtUI)

	w := get(r, "/api/community/demographics")

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(1), body["member_count"])
	assert.NotContains(t, body, "members")
}

func TestDemographicsUpstreamFailureIsABadGateway(t *testing.T) {
	w := get(newTestRouter(fakeMembers{err: errors.New("secret detail")}, builtUI), "/api/community/demographics")

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.NotContains(t, w.Body.String(), "secret detail")
}

func TestFrontendRouting(t *testing.T) {
	r := newTestRouter(fakeMembers{}, builtUI)

	tests := []struct {
		path     string
		status   int
		contains string
	}{
		{"/", http.StatusOK, "Growth Portal"},
		{"/dashboards/interests", http.StatusOK, "Growth Portal"}, // client-side route
		{"/assets/app-1.js", http.StatusOK, "console.log"},
		{"/favicon.svg", http.StatusOK, "<svg/>"},
		{"/api/nope", http.StatusNotFound, "not found"},
		{"/api/health", http.StatusOK, "healthy"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			w := get(r, tt.path)
			assert.Equal(t, tt.status, w.Code)
			assert.Contains(t, w.Body.String(), tt.contains)
		})
	}
}

func TestAssetsAreCachedAndIndexIsNot(t *testing.T) {
	r := newTestRouter(fakeMembers{}, builtUI)

	assert.Contains(t, get(r, "/assets/app-1.js").Header().Get("Cache-Control"), "immutable")
	assert.Equal(t, "no-cache", get(r, "/").Header().Get("Cache-Control"))
}

func TestFrontendNotBuilt(t *testing.T) {
	r := newTestRouter(fakeMembers{}, fstest.MapFS{".gitkeep": {}})

	w := get(r, "/")

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), "frontend not built")
}
