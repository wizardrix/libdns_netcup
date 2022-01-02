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

func (p *Provider) doRequest(ctx context.Context, req request) (*response, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResp.Body)
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

	res, err := p.doRequest(ctx, loginRequest)
	if err != nil {
		return "", err
	}

	var asd apiSessionData
	if err = json.Unmarshal(res.ResponseData, &asd); err != nil {
		return "", err
	}

	return asd.ApiSessionId, nil
}

func (p *Provider) logout(ctx context.Context, apiSessionID string) {
	logoutRequest := request{
		Action: "logout",
		Param: requestParam{
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiSessionID:   apiSessionID,
		},
	}

	p.doRequest(ctx, logoutRequest)
}

func (p *Provider) infoDNSZone(ctx context.Context, zone string, apiSessionID string) (*dnsZone, error) {
	infoDNSZoneRequest := request{
		Action: "infoDnsZone",
		Param: requestParam{
			DomainName:     zone,
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiSessionID:   apiSessionID,
		},
	}

	res, err := p.doRequest(ctx, infoDNSZoneRequest)
	if err != nil {
		return nil, err
	}

	var dz dnsZone
	if err = json.Unmarshal(res.ResponseData, &dz); err != nil {
		return nil, err
	}

	return &dz, nil
}

func (p *Provider) infoDNSRecords(ctx context.Context, zone string, apiSessionID string) (*dnsRecordSet, error) {
	infoDNSrecordsRequest := request{
		Action: "infoDnsRecords",
		Param: requestParam{
			DomainName:     zone,
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiSessionID:   apiSessionID,
		},
	}

	res, err := p.doRequest(ctx, infoDNSrecordsRequest)
	if err != nil {
		return nil, err
	}

	var recordSet dnsRecordSet
	if err = json.Unmarshal(res.ResponseData, &recordSet); err != nil {
		return nil, err
	}

	return &recordSet, err
}

func (p *Provider) updateDNSRecords(ctx context.Context, zone string, updateRecordSet dnsRecordSet, apiSessionID string) (*dnsRecordSet, error) {
	updateDNSrecordsRequest := request{
		Action: "updateDnsRecords",
		Param: requestParam{
			DomainName:     zone,
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiSessionID:   apiSessionID,
			DnsRecordSet:   updateRecordSet,
		},
	}

	res, err := p.doRequest(ctx, updateDNSrecordsRequest)
	if err != nil {
		return nil, err
	}

	var recordSet dnsRecordSet
	if err = json.Unmarshal(res.ResponseData, &recordSet); err != nil {
		return nil, err
	}

	return &recordSet, err
}
