package netcup

import (
	"time"

	"github.com/libdns/libdns"
)

func toLibdnsRecords(netcupRecords []dnsRecord, ttl int64) []libdns.Record {
	var libdnsRecords []libdns.Record
	for _, record := range netcupRecords {
		libdnsRecord := libdns.Record{
			ID:       record.ID,
			Type:     record.RecType,
			Name:     record.HostName,
			Value:    record.Destination,
			TTL:      time.Duration(ttl * int64(time.Second)),
			Priority: record.Priority,
		}
		libdnsRecords = append(libdnsRecords, libdnsRecord)
	}
	return libdnsRecords
}

func toNetcupRecords(libnsRecords []libdns.Record) []dnsRecord {
	var netcupRecords []dnsRecord
	for _, record := range libnsRecords {
		netcupRecord := dnsRecord{
			ID:          record.ID,
			HostName:    record.Name,
			RecType:     record.Type,
			Destination: record.Value,
			Priority:    record.Priority,
		}
		netcupRecords = append(netcupRecords, netcupRecord)
	}
	return netcupRecords
}

// difference returns the records that are in a but not in b
func difference(a, b []dnsRecord) []dnsRecord {
	bIDmap := make(map[string]struct{}, len(b))
	for _, elm := range b {
		bIDmap[elm.ID] = struct{}{}
	}

	var diff []dnsRecord
	for _, elm := range a {
		if _, found := bIDmap[elm.ID]; !found {
			diff = append(diff, elm)
		}
	}

	return diff
}
