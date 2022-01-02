// Package netcup implements a DNS record management client compatible
// with the libdns interfaces for netcup.
package netcup

import (
	"context"
	"fmt"
	"sync"

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

	apiSessionID, err := p.login(ctx)
	if err != nil {
		return nil, err
	}
	defer p.logout(ctx, apiSessionID)

	dnsZone, err := p.infoDNSZone(ctx, zone, apiSessionID)
	if err != nil {
		return nil, err
	}

	recordSet, err := p.infoDNSRecords(ctx, zone, apiSessionID)
	if err != nil {
		return nil, err
	}

	return toLibdnsRecords(recordSet.DnsRecords, dnsZone.TTL), nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
// netcup records cannot have individual TTLs, there is one TTL for all records in the zone
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	apiSessionID, err := p.login(ctx)
	if err != nil {
		return nil, err
	}
	defer p.logout(ctx, apiSessionID)

	dnsZone, err := p.infoDNSZone(ctx, zone, apiSessionID)
	if err != nil {
		return nil, err
	}

	recordSetBefore, err := p.infoDNSRecords(ctx, zone, apiSessionID)
	if err != nil {
		return nil, err
	}

	updateRecordSet := dnsRecordSet{
		DnsRecords: toNetcupRecords(records),
	}
	recordSetAfter, err := p.updateDNSRecords(ctx, zone, updateRecordSet, apiSessionID)
	if err != nil {
		return nil, err
	}

	appendedRecords := difference(recordSetAfter.DnsRecords, recordSetBefore.DnsRecords)

	return toLibdnsRecords(appendedRecords, dnsZone.TTL), nil
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
