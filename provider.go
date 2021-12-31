// Package libdnstemplate implements a DNS record management client compatible
// with the libdns interfaces for <PROVIDER NAME>. TODO: This package is a
// template only. Customize all godocs for actual implementation.
package netcup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

	res, err := p.doRequest(req)
	if err != nil {
		return nil, err
	}

	var recordSet dnsRecordSet
	if err = json.Unmarshal(res.ResponseData, &recordSet); err != nil {
		return nil, err
	}

	var libDNSRecords []libdns.Record
	for _, record := range recordSet.DnsRecords {
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
