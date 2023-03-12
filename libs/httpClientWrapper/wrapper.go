package httpClientWrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Wrapper[Req, Res any] struct {
	url string
}

func New[Res, Req any](url string) *Wrapper[Res, Req] {
	return &Wrapper[Res, Req]{
		url,
	}
}

func (w *Wrapper[Req, Res]) SendRequest(req Req) (*Res, error) {
	bodyJSON, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling JSON")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, w.url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, errors.Wrap(err, "creating HTTP request")
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "calling HTTP")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var res Res
	err = json.NewDecoder(httpResponse.Body).Decode(&res)
	if err != nil {
		return &res, errors.Wrap(err, "decoding JSON from response")
	}

	return &res, nil
}
