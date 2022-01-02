// Package netcup implements a DNS record management client compatible
// with the libdns interfaces for netcup.
package netcup

import (
	"context"
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
//
// For each input record, if no ID is given, the first record that matches the host name and type is searched.
// If none is found or the search result doesn't equal the input, a new one is appended.
// For MX records the priority is needed as an additional search parameter.
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

	existingRecordSet, err := p.infoDNSRecords(ctx, zone, apiSessionID)
	if err != nil {
		return nil, err
	}

	netcupRecords := toNetcupRecords(records)
	recordsToAppend := getRecordsToAppend(netcupRecords, existingRecordSet.DnsRecords)
	if len(recordsToAppend) == 0 {
		return []libdns.Record{}, nil
	}
	recordSetToAppend := dnsRecordSet{
		DnsRecords: recordsToAppend,
	}
	if _, err = p.updateDNSRecords(ctx, zone, recordSetToAppend, apiSessionID); err != nil {
		return nil, err
	}

	return toLibdnsRecords(recordsToAppend, dnsZone.TTL), nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
//
// netcup records cannot have individual TTLs, there is one TTL for all records in the zone. So these can not be set.
//
// For each input record, if no ID is given, the first record that matches the host name and type is searched.
// If none is found, the input is appended. If one is found, it is updated accordingly.
// For MX records the priority is needed as an additional search parameter.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
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

	existingRecordSet, err := p.infoDNSRecords(ctx, zone, apiSessionID)
	if err != nil {
		return nil, err
	}

	netcupRecords := toNetcupRecords(records)
	recordsToSet := getRecordsToSet(netcupRecords, existingRecordSet.DnsRecords)
	if len(recordsToSet) == 0 {
		return []libdns.Record{}, nil
	}
	recordSetToSet := dnsRecordSet{
		DnsRecords: recordsToSet,
	}
	if _, err = p.updateDNSRecords(ctx, zone, recordSetToSet, apiSessionID); err != nil {
		return nil, err
	}

	return toLibdnsRecords(recordsToSet, dnsZone.TTL), nil
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
//
// For each input record, if no ID is given, the first record that matches the host name and type is searched and deleted.
// For MX records the priority is needed as an additional search parameter.
// To be safe, the records to delete should include the IDs (for example from GetRecords)
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
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

	netcupRecords := toNetcupRecords(records)
	recordsToDelete := getRecordsToDelete(netcupRecords, recordSet.DnsRecords)
	if len(recordsToDelete) == 0 {
		return []libdns.Record{}, nil
	}
	recordSetToDelete := dnsRecordSet{
		DnsRecords: recordsToDelete,
	}
	if _, err = p.updateDNSRecords(ctx, zone, recordSetToDelete, apiSessionID); err != nil {
		return nil, err
	}

	return toLibdnsRecords(recordsToDelete, dnsZone.TTL), nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
