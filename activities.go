package workshop

import (
	"context"
	"fmt"
	"io"
	"net/http"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

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

func (a *Activities) AnalyzeMetrics(ctx context.Context, metrics string) (string, error) {
	prompt := `You are a system health analyst. Given the following Prometheus node_exporter metrics, respond in this exact format:

Hostname: <value from node_uname_info nodename label>
OS: <value from node_os_info name and version labels>
CPU: <number of unique cpu values in node_cpu_seconds_total> cores
RAM: <node_memory_MemTotal_bytes converted to GB, rounded to 1 decimal>
Disk: <node_filesystem_size_bytes for mountpoint="/" converted to GB, rounded to 0 decimal>GB total, <node_filesystem_avail_bytes for mountpoint="/" converted to GB, rounded to 0 decimal>GB available

<2-3 sentence health summary highlighting notable CPU, memory, or disk values>

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
		return "", fmt.Errorf("analyze metrics: %w", err)
	}
	if len(msg.Content) == 0 {
		return "", fmt.Errorf("empty response from AI")
	}
	return msg.Content[0].Text, nil
}
