package webhook

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *clientImpl) Route1(ctx context.Context, req *Route1Request) (*Route1Response, int, error) {
	resp, err := c.client.R().
		SetBody(req).
		SetResult(Route1Response{}).
		Post(c.baseURL + c.loginRoute)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	code, err := handleErrorResponse(resp)
	if err != nil {
		return nil, code, err
	}

	var response Route1Response
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &response, code, nil
}
