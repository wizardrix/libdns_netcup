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
		return nil, fmt.Errorf("[netcup] %v: %v", response.ShortMessage, response.LongMessage)
	}

	fmt.Printf("[netcup] %v: %v\n", response.ShortMessage, response.LongMessage)

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

func (p *Provider) logout(ctx context.Context, apiSessionID string) error {
	loginRequest := request{
		Action: "logout",
		Param: requestParam{
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiSessionID:   apiSessionID,
		},
	}

	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	res, err := p.doRequest(req)
	if err != nil {
		return err
	}

	var asd apiSessionData
	if err = json.Unmarshal(res.ResponseData, &asd); err != nil {
		return err
	}

	return nil
}

func (p *Provider) getDNSZone(ctx context.Context, zone string, apiSessionID string) (*dnsZone, error) {
	infoDNSZoneRequest := request{
		Action: "infoDnsZone",
		Param: requestParam{
			DomainName:     zone,
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiSessionID:   apiSessionID,
		},
	}

	requestBody, err := json.Marshal(infoDNSZoneRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	res, err := p.doRequest(req)
	if err != nil {
		return nil, err
	}

	var dz dnsZone
	if err = json.Unmarshal(res.ResponseData, &dz); err != nil {
		return nil, err
	}

	fmt.Printf("Zone %v TTL: %v\n", dz.Name, dz.TTL)

	return &dz, nil
}
