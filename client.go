package netcup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiUrl = "https://ccp.netcup.net/run/webservice/servers/endpoint.php?JSON"

func (p *Provider) doRequest(req *http.Request) (*response, error) {
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response response
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("netcup %v failed: %v (reason: %v)", response.Action, response.ShortMessage, response.LongMessage)
	}

	fmt.Printf("netcup %v successful\n", response.Action)

	return &response, nil
}

func (p *Provider) login(ctx context.Context) (string, error) {
	loginRequest := request{
		Action: "login",
		Param: requestParam{
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiPassword:    p.ApiPassword,
		},
	}

	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}

	res, err := p.doRequest(req)
	if err != nil {
		return "", err
	}

	var asd apiSessionData
	if err = json.Unmarshal(res.ResponseData, &asd); err != nil {
		return "", err
	}

	fmt.Printf("Session ID: %v\n", asd.ApiSessionId)

	return asd.ApiSessionId, nil
}
