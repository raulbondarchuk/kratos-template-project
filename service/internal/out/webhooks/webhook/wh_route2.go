package webhook

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *clientImpl) Route2(ctx context.Context, req *Route2Request) (*Route2Response, int, error) {
	resp, err := c.client.R().
		SetQueryParams(map[string]string{"username": req.Username}).
		SetResult(Route2Response{}).
		Get(c.baseURL + c.userRoute)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	code, err := handleErrorResponse(resp)
	if err != nil {
		return nil, code, err
	}

	var response Route2Response
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &response, code, nil
}
