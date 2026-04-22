package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type HealthReport struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	CPU      string `json:"cpu"`
	RAM      string `json:"ram"`
	Disk     string `json:"disk"`
	Summary  string `json:"summary"`
}

type Activities struct {
	HTTP       *http.Client
	MetricsURL string
	ai         anthropic.Client
	aiModel    string
}

func NewActivities(httpClient *http.Client, metricsURL, aiURL, aiModel string) *Activities {
	return &Activities{
		HTTP:       httpClient,
		MetricsURL: metricsURL,
		aiModel:    aiModel,
		ai: anthropic.NewClient(
			option.WithBaseURL(aiURL),
			option.WithAPIKey("x"),
			option.WithHTTPClient(httpClient),
		),
	}
}

func (a *Activities) FetchMetrics(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.MetricsURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := a.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch metrics: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("fetch metrics: unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 128*1024))
	if err != nil {
		return "", fmt.Errorf("read metrics: %w", err)
	}
	return string(body), nil
}

func (a *Activities) AnalyzeMetrics(ctx context.Context, metrics string) (HealthReport, error) {
	prompt := `You are a system health analyst. Given the following Prometheus node_exporter metrics, respond with ONLY a JSON object, no prose, no markdown, no code fences. Use exactly these keys:

{
  "hostname": "<value from node_uname_info nodename label>",
  "os": "<value from node_os_info name and version labels>",
  "cpu": "<N cores, the number of unique cpu values in node_cpu_seconds_total>",
  "ram": "<node_memory_MemTotal_bytes converted to GB, rounded to 1 decimal, like \"36.0 GB\">",
  "disk": "<NGB total, MGB available for mountpoint=/ rounded to 0 decimals>",
  "summary": "<2-3 sentence health summary highlighting notable CPU, memory, or disk values>"
}

Metrics:
` + metrics

	msg, err := a.ai.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(a.aiModel),
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return HealthReport{}, fmt.Errorf("analyze metrics: %w", err)
	}
	if len(msg.Content) == 0 {
		return HealthReport{}, fmt.Errorf("empty response from AI")
	}

	raw := strings.TrimSpace(msg.Content[0].Text)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var report HealthReport
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		return HealthReport{}, fmt.Errorf("parse health report: %w (raw: %q)", err, raw)
	}
	return report, nil
}
