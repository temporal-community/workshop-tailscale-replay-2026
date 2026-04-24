package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

// newActivityEnv returns a test activity environment with the given
// Activities registered. Using the env (instead of calling methods
// directly) wires up the activity.* context helpers, e.g. GetLogger.
func newActivityEnv(a *Activities) *testsuite.TestActivityEnvironment {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestActivityEnvironment()
	env.RegisterActivity(a)
	return env
}

func TestFetchMetrics(t *testing.T) {
	const fakeMetrics = "# HELP node_cpu_seconds_total\nnode_cpu_seconds_total{cpu=\"0\",mode=\"idle\"} 12345.0\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, fakeMetrics)
	}))
	t.Cleanup(srv.Close)

	a := &Activities{HTTP: srv.Client(), MetricsURL: srv.URL}
	val, err := newActivityEnv(a).ExecuteActivity(a.FetchMetrics)
	require.NoError(t, err)
	var got string
	require.NoError(t, val.Get(&got))
	require.Equal(t, fakeMetrics, got)
}

func TestAnalyzeMetrics(t *testing.T) {
	want := HealthReport{
		Hostname: "test-host",
		OS:       "Test OS 1.0",
		CPU:      "8 cores",
		RAM:      "16.0 GB",
		Disk:     "500GB total, 250GB available",
		Summary:  "System looks healthy.",
	}
	wantJSON, err := json.Marshal(want)
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/messages", r.URL.Path)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "claude-haiku-4-5", body["model"])
		msgs := body["messages"].([]any)
		require.Len(t, msgs, 1)
		msg := msgs[0].(map[string]any)
		require.Equal(t, "user", msg["role"])
		require.NotNil(t, msg["content"])

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "msg_test",
			"type":        "message",
			"role":        "assistant",
			"model":       "claude-haiku-4-5",
			"stop_reason": "end_turn",
			"content": []map[string]any{
				{"type": "text", "text": string(wantJSON)},
			},
			"usage": map[string]any{"input_tokens": 10, "output_tokens": 5},
		})
	}))
	t.Cleanup(srv.Close)

	a := NewActivities(srv.Client(), "", srv.URL, "claude-haiku-4-5")
	val, err := newActivityEnv(a).ExecuteActivity(a.AnalyzeMetrics, "some metrics")
	require.NoError(t, err)
	var got HealthReport
	require.NoError(t, val.Get(&got))
	require.Equal(t, want, got)
}

func TestAnalyzeMetrics_StripsCodeFences(t *testing.T) {
	want := HealthReport{
		Hostname: "fenced-host",
		Summary:  "ok",
	}
	wantJSON, err := json.Marshal(want)
	require.NoError(t, err)
	fencedResponse := "```json\n" + string(wantJSON) + "\n```"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "msg_test",
			"type":        "message",
			"role":        "assistant",
			"model":       "claude-haiku-4-5",
			"stop_reason": "end_turn",
			"content": []map[string]any{
				{"type": "text", "text": fencedResponse},
			},
			"usage": map[string]any{"input_tokens": 10, "output_tokens": 5},
		})
	}))
	t.Cleanup(srv.Close)

	a := NewActivities(srv.Client(), "", srv.URL, "claude-haiku-4-5")
	val, err := newActivityEnv(a).ExecuteActivity(a.AnalyzeMetrics, "some metrics")
	require.NoError(t, err)
	var got HealthReport
	require.NoError(t, val.Get(&got))
	require.Equal(t, want, got)
}

func TestAnalyzeMetrics_AIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"type": "error",
			"error": map[string]any{
				"type":    "invalid_request_error",
				"message": "no route found for model",
			},
		})
	}))
	t.Cleanup(srv.Close)

	a := NewActivities(srv.Client(), "", srv.URL, "claude-haiku-4-5")
	_, err := newActivityEnv(a).ExecuteActivity(a.AnalyzeMetrics, "some metrics")
	require.Error(t, err)
}
