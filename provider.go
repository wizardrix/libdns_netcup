// Package netcup implements a DNS record management client compatible
// with the libdns interfaces for netcup.
package netcup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/libdns/libdns"
)

// TODO: Providers must not require additional provisioning steps by the callers; it
// should work simply by populating a struct and calling methods on it. If your DNS
// service requires long-lived state or some extra provisioning step, do it implicitly
// when methods are called; sync.Once can help with this, and/or you can use a
// sync.(RW)Mutex in your Provider struct to synchronize implicit provisioning.

// Provider facilitates DNS record manipulation with netcup.
type Provider struct {
	CustomerNumber string `json:"customer_number"`
	ApiKey         string `json:"api_key"`
	ApiPassword    string `json:"api_password"`
	mutex          sync.Mutex
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	apiSessionId, err := p.login(ctx)
	if err != nil {
		return nil, err
	}
	defer p.logout(ctx, apiSessionId)

	dnsZone, err := p.getDNSZone(ctx, zone, apiSessionId)
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
			ID:       record.ID,
			Type:     record.RecType,
			Name:     record.HostName,
			Value:    record.Destination,
			TTL:      time.Duration(dnsZone.TTL * int64(time.Second)),
			Priority: record.Priority,
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
