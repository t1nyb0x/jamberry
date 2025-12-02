package tracktaste

import (
	"context"
	"fmt"
)

// HealthResult はヘルスチェックの結果を表します
type HealthResult struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	BuildTime string            `json:"build_time"`
	GitCommit string            `json:"git_commit"`
	Uptime    string            `json:"uptime"`
	Services  map[string]string `json:"services"`
}

// FetchHealth はTrackTaste APIのヘルスチェックを行い、バージョン情報を取得します
func (c *Client) FetchHealth(ctx context.Context) (*HealthResult, error) {
	endpoint := fmt.Sprintf("%s/healthz", c.baseURL)

	resp, err := doRequest[HealthResult](ctx, c, endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
