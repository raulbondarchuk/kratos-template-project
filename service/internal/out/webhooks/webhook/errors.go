package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

func handleErrorResponse(resp *resty.Response) (int, error) {
	if resp.StatusCode() != http.StatusOK {
		if err := printIfNotJSON(resp); err != nil {
			return resp.StatusCode(), err
		}

		var m map[string]any
		if err := json.Unmarshal(resp.Body(), &m); err != nil {
			return resp.StatusCode(), errors.New("(getlineas) failed to parse error response: " + err.Error())
		}

		if v, ok := m["error"]; ok {
			if s, ok := v.(string); ok && s != "" {
				return resp.StatusCode(), errors.New(s)
			}
			if mm, ok := v.(map[string]any); ok {
				if ms, ok := mm["message"].(string); ok && ms != "" {
					return resp.StatusCode(), errors.New(ms)
				}
			}
		}
		if s, ok := m["message"].(string); ok && s != "" {
			return resp.StatusCode(), errors.New(s)
		}
		return resp.StatusCode(), errors.New("unknown error")
	}
	return http.StatusOK, nil
}

func printIfNotJSON(resp *resty.Response) error {
	contentType := resp.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		fmt.Printf("Non-JSON response body:\n%s\n", resp.String())
		return fmt.Errorf("response is not JSON, got Content-Type: %s | Check logs for more details", contentType)
	}
	return nil
}
