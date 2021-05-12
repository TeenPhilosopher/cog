package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/replicate/cog/pkg/model"
)

func (c *Client) GetModel(repo *model.Repo, id string) (*model.Model, []*model.Image, error) {
	url := newURL(repo, "v1/repos/%s/%s/models/%s", repo.User, repo.Name, id)
	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("Model not found: %s:%s", repo.String(), id)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("Server returned status %d: %s", resp.StatusCode, body)
	}
	body := struct {
		Model  *model.Model   `json:"model"`
		Images []*model.Image `json:"images"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, nil, err
	}
	return body.Model, body.Images, nil
}
