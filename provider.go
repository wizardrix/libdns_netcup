// Package libdnstemplate implements a DNS record management client compatible
// with the libdns interfaces for <PROVIDER NAME>. TODO: This package is a
// template only. Customize all godocs for actual implementation.
package netcup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/libdns/libdns"
)

// TODO: Providers must not require additional provisioning steps by the callers; it
// should work simply by populating a struct and calling methods on it. If your DNS
// service requires long-lived state or some extra provisioning step, do it implicitly
// when methods are called; sync.Once can help with this, and/or you can use a
// sync.(RW)Mutex in your Provider struct to synchronize implicit provisioning.

// Provider facilitates DNS record manipulation with <TODO: PROVIDER NAME>.
type Provider struct {
	// TODO: put config fields here (with snake_case json
	// struct tags on exported fields), for example:
	CustomerNumber string `json:"customer_number"`
	ApiKey         string `json:"api_key"`
	ApiPassword    string `json:"api_password"`
	mu             sync.Mutex
}

type dnsRecord struct {
	ID               string `json:"id"`
	HostName         string `json:"hostname"`
	RecType          string `json:"type"`
	Destination      string `json:"destination"`
	DeleteRecordFlag bool   `json:"deleterecordflag"`
}

type dnsRecordSet struct {
	DnsRecords []dnsRecord `json:"dnsrecords"`
}

type requestParam struct {
	DomainName     string       `json:"domainname,omitempty"`
	CustomerNumber string       `json:"customernumber"`
	ApiKey         string       `json:"apikey"`
	ApiPassword    string       `json:"apipassword,omitempty"`
	ApiSessionID   string       `json:"apisessionid,omitempty"`
	DnsRecordSet   dnsRecordSet `json:"dnsrecordset,omitempty"`
}

type request struct {
	Action string       `json:"action"`
	Param  requestParam `json:"param"`
}

type response struct {
	Action       string `json:"action"`
	Status       string `json:"status"`
	ShortMessage string `json:"shortmessage"`
	LongMessage  string `json:"longmessage"`
}

type apiSessionData struct {
	ApiSessionId string `json:"apisessionid"`
}

type loginResponse struct {
	response
	ResponseData apiSessionData `json:"responsedata"`
}

type dnsInfoResponse struct {
	response
	ResponseData dnsRecordSet `json:"responsedata"`
}

const apiUrl = "https://ccp.netcup.net/run/webservice/servers/endpoint.php?JSON"

func (p *Provider) doRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
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

	responseBody, err := p.doRequest(req)
	if err != nil {
		return "", err
	}

	var response loginResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}
	if response.Status == "error" {
		return "", fmt.Errorf("login failed: %v (%v)", response.ShortMessage, response.LongMessage)
	}
	fmt.Printf("Login successful, session ID: %v\n", response.ResponseData.ApiSessionId)

	return response.ResponseData.ApiSessionId, nil
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	apiSessionId, err := p.login(ctx)
	if err != nil {
		return nil, err
	}

	infoDNSrecordsRequest := request{
		Action: "infoDnsRecords",
		Param: requestParam{
			DomainName:     zone,
			CustomerNumber: p.CustomerNumber,
			ApiKey:         p.ApiKey,
			ApiSessionID:   apiSessionId,
		},
	}

	requestBody, err := json.Marshal(infoDNSrecordsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	responseBody, err := p.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response dnsInfoResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	if response.Status == "error" {
		return nil, fmt.Errorf("getRecords failed: %v (%v)", response.ShortMessage, response.LongMessage)
	}

	var libDNSRecords []libdns.Record
	for _, record := range response.ResponseData.DnsRecords {
		libDNSRecord := libdns.Record{
			ID:    record.ID,
			Type:  record.RecType,
			Name:  record.HostName,
			Value: record.Destination,
		}
		libDNSRecords = append(libDNSRecords, libDNSRecord)
	}

	return libDNSRecords, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
