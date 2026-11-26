// Lab #7 - functional case for the leaderboard (non-personalized) recommendation API.
//
// The case starts the real Gorse REST handler (RestServer.CreateWebService) behind a
// local HTTP listener, backed by SQLite data/cache stores in a temporary directory,
// seeds a fixed 15-item leaderboard, and calls GET /api/non-personalized/popular
// over HTTP. The endpoint is served by RestServer.SearchDocuments, which also serves
// item-to-item, user-to-user, item/user neighbours and collaborative-filtering lists.
//
// Two groups of checks:
//   - invariants: behaviour that must be identical before and after the Lab #7 change
//     (top-n ordering, default n, read-item exclusion, non-numeric n rejected);
//   - the changed case: n=-1. Its expectation depends on the *contract* claimed by
//     this version of the source, stored in lab7/claimed-contract.txt (shared with
//     lab7/e2e.sh). The workflow additionally runs the case against the other
//     contract (LAB7_CONTRACT override) and requires it to FAIL, which shows that the
//     check can tell the two versions apart.
//
// Every observation is printed as a "LAB7 RESULT" line and, when LAB7_RESULTS is
// set, appended as one JSON object per line to that file (kept as a CI artifact).
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/gorse-io/gorse/config"
	"github.com/gorse-io/gorse/storage/cache"
	"github.com/gorse-io/gorse/storage/data"
	"github.com/stretchr/testify/require"
)

const (
	// lab7ContractLegacy: upstream behaviour - a negative n is not rejected and the
	// whole leaderboard is returned (unbounded response).
	lab7ContractLegacy = "legacy"
	// lab7ContractBounded: a negative n is rejected with 400 Bad Request.
	lab7ContractBounded = "bounded"

	// lab7ClaimedContractFile holds the contract this version of the source claims to
	// implement (path relative to the package directory, where `go test` runs).
	lab7ClaimedContractFile = "../lab7/claimed-contract.txt"

	lab7APIKey      = "lab7_api_key"
	lab7Leaderboard = "popular"
	lab7NumItems    = 15
	lab7ReadUser    = "alice"
	lab7ReadItem    = "item-01"
)

type lab7Observation struct {
	Case     string   `json:"case"`
	Query    string   `json:"query"`
	Status   int      `json:"status"`
	Count    int      `json:"count"`
	Ids      []string `json:"ids,omitempty"`
	Body     string   `json:"body,omitempty"`
	Contract string   `json:"contract,omitempty"`
}

// lab7Contract returns the contract to check the changed case against: the claimed one,
// unless LAB7_CONTRACT overrides it (used by the workflow's counter-check).
func lab7Contract(t *testing.T) string {
	if c := os.Getenv("LAB7_CONTRACT"); c != "" {
		return c
	}
	b, err := os.ReadFile(lab7ClaimedContractFile)
	require.NoError(t, err)
	return strings.TrimSpace(string(b))
}

func lab7ItemId(i int) string { return fmt.Sprintf("item-%02d", i) }

// lab7StartServer deploys the REST handler on a local HTTP listener with a fixed dataset.
func lab7StartServer(t *testing.T) *httptest.Server {
	t.Helper()
	ctx := t.Context()
	s := &RestServer{}
	s.Config = config.GetDefaultConfig()
	s.Config.Server.APIKey = lab7APIKey
	var err error
	s.DataClient, err = data.Open(fmt.Sprintf("sqlite://%s/data.db", t.TempDir()), "")
	require.NoError(t, err)
	s.CacheClient, err = cache.Open(fmt.Sprintf("sqlite://%s/cache.db", t.TempDir()), "")
	require.NoError(t, err)
	require.NoError(t, s.DataClient.Init())
	require.NoError(t, s.CacheClient.Init())
	t.Cleanup(func() {
		_ = s.DataClient.Close()
		_ = s.CacheClient.Close()
	})

	// Fixed leaderboard: item-00 (score 100) ... item-14 (score 86).
	scores := make([]cache.Score, 0, lab7NumItems)
	for i := 0; i < lab7NumItems; i++ {
		scores = append(scores, cache.Score{
			Id:         lab7ItemId(i),
			Score:      float64(100 - i),
			Categories: []string{""},
			Timestamp:  time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
		})
	}
	require.NoError(t, s.CacheClient.AddScores(ctx, cache.NonPersonalized, lab7Leaderboard, scores))
	// alice has already read item-01, so it must be excluded when user-id=alice.
	require.NoError(t, s.DataClient.BatchInsertFeedback(ctx, []data.Feedback{{
		FeedbackKey: data.FeedbackKey{FeedbackType: "read", UserId: lab7ReadUser, ItemId: lab7ReadItem},
		Timestamp:   time.Now().Add(-time.Hour),
	}}, true, true, true))

	s.WebService = new(restful.WebService)
	s.CreateWebService()
	container := restful.NewContainer()
	container.Add(s.WebService)
	srv := httptest.NewServer(container)
	t.Cleanup(srv.Close)
	return srv
}

// lab7Get calls the leaderboard endpoint over HTTP and records the observation.
func lab7Get(t *testing.T, srv *httptest.Server, name, query string) lab7Observation {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		srv.URL+"/api/non-personalized/"+lab7Leaderboard+query, nil)
	require.NoError(t, err)
	req.Header.Set("X-API-Key", lab7APIKey)
	resp, err := srv.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	obs := lab7Observation{Case: name, Query: query, Status: resp.StatusCode}
	if resp.StatusCode == http.StatusOK {
		var items []cache.Score
		require.NoError(t, json.Unmarshal(body, &items), "body: %s", body)
		for _, it := range items {
			obs.Ids = append(obs.Ids, it.Id)
		}
		obs.Count = len(items)
	} else {
		obs.Body = string(body)
	}
	return obs
}

func lab7Record(t *testing.T, obs lab7Observation) {
	t.Helper()
	line, err := json.Marshal(obs)
	require.NoError(t, err)
	t.Logf("LAB7 RESULT %s", line)
	if path := os.Getenv("LAB7_RESULTS"); path != "" {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		require.NoError(t, err)
		defer f.Close()
		_, err = f.Write(append(line, '\n'))
		require.NoError(t, err)
	}
}

func lab7Ids(indices ...int) []string {
	ids := make([]string, 0, len(indices))
	for _, i := range indices {
		ids = append(ids, lab7ItemId(i))
	}
	return ids
}

// TestLab7LeaderboardInvariants: behaviour that the Lab #7 change must NOT alter.
func TestLab7LeaderboardInvariants(t *testing.T) {
	srv := lab7StartServer(t)

	t.Run("top_n", func(t *testing.T) {
		obs := lab7Get(t, srv, "top_n", "?n=3")
		lab7Record(t, obs)
		require.Equal(t, http.StatusOK, obs.Status)
		require.Equal(t, lab7Ids(0, 1, 2), obs.Ids)
	})
	t.Run("default_n", func(t *testing.T) {
		obs := lab7Get(t, srv, "default_n", "")
		lab7Record(t, obs)
		require.Equal(t, http.StatusOK, obs.Status)
		require.Equal(t, lab7Ids(0, 1, 2, 3, 4, 5, 6, 7, 8, 9), obs.Ids) // server.default_n = 10
	})
	t.Run("offset_n", func(t *testing.T) {
		obs := lab7Get(t, srv, "offset_n", "?offset=12&n=5")
		lab7Record(t, obs)
		require.Equal(t, http.StatusOK, obs.Status)
		require.Equal(t, lab7Ids(12, 13, 14), obs.Ids)
	})
	t.Run("exclude_read_items", func(t *testing.T) {
		obs := lab7Get(t, srv, "exclude_read_items", "?n=3&user-id="+lab7ReadUser)
		lab7Record(t, obs)
		require.Equal(t, http.StatusOK, obs.Status)
		require.Equal(t, lab7Ids(0, 2, 3), obs.Ids)
		require.NotContains(t, obs.Ids, lab7ReadItem)
	})
	t.Run("non_numeric_n_rejected", func(t *testing.T) {
		obs := lab7Get(t, srv, "non_numeric_n_rejected", "?n=abc")
		lab7Record(t, obs)
		require.Equal(t, http.StatusBadRequest, obs.Status)
	})
}

// TestLab7NegativeLimit: the case whose behaviour is intentionally changed in Lab #7.
func TestLab7NegativeLimit(t *testing.T) {
	srv := lab7StartServer(t)
	contract := lab7Contract(t)
	obs := lab7Get(t, srv, "negative_n", "?n=-1")
	obs.Contract = contract
	lab7Record(t, obs)

	switch contract {
	case lab7ContractLegacy:
		require.Equal(t, http.StatusOK, obs.Status, "legacy contract: negative n is accepted")
		require.Equal(t, lab7NumItems, obs.Count, "legacy contract: the whole leaderboard is returned")
	case lab7ContractBounded:
		require.Equal(t, http.StatusBadRequest, obs.Status, "bounded contract: negative n is rejected")
		require.Contains(t, obs.Body, "invalid n parameter", "bounded contract: error explains the rejected parameter")
		require.Zero(t, obs.Count, "bounded contract: no items are returned")
	default:
		t.Fatalf("unknown LAB7_CONTRACT %q", contract)
	}
}
