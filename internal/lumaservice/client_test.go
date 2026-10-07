package lumaservice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityMembersSendsTheSecret(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/community/members", r.URL.Path)
		assert.Equal(t, "s3cret", r.Header.Get("X-Service-Secret"))
		_, _ = w.Write([]byte(`{"tier_id":"t","members":[{"id":"c1","email":"a@example.com","created_at":"2026-09-01T10:00:00Z",
			"answers":[{"question_id":"q1","label":"Graduation year","type":"select","values":["2027"]}]}]}`))
	}))
	defer srv.Close()

	members, err := New(srv.URL, "s3cret").CommunityMembers(context.Background())

	require.NoError(t, err)
	require.Len(t, members, 1)
	assert.Equal(t, []string{"2027"}, members[0].Answers[0].Values)
}

func TestCommunityMembersNon200IsAnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	_, err := New(srv.URL, "wrong").CommunityMembers(context.Background())

	assert.ErrorContains(t, err, "status 401")
}
