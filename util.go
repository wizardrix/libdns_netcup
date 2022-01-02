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
// func difference(a, b []dnsRecord) []dnsRecord {
// 	bIDmap := make(map[string]struct{}, len(b))
// 	for _, elm := range b {
// 		bIDmap[elm.ID] = struct{}{}
// 	}

// 	var diff []dnsRecord
// 	for _, elm := range a {
// 		if _, found := bIDmap[elm.ID]; !found {
// 			diff = append(diff, elm)
// 		}
// 	}

// 	return diff
// }

func findRecordByID(id string, records []dnsRecord) *dnsRecord {
	for _, record := range records {
		if record.ID == id {
			return &record
		}
	}

	return nil
}

func findRecordByNameAndType(hostName string, recType string, records []dnsRecord) *dnsRecord {
	for _, record := range records {
		if record.HostName == hostName && record.RecType == recType {
			return &record
		}
	}

	return nil
}

func findRecordByNameAndTypeAndPriority(hostName string, recType string, priority int, records []dnsRecord) *dnsRecord {
	for _, record := range records {
		if record.HostName == hostName && record.RecType == recType && record.Priority == priority {
			return &record
		}
	}

	return nil
}

// func findRecordByExample(exampleRecord dnsRecord, records []dnsRecord) *dnsRecord {
// 	for _, record := range records {
// 		if record.equals(exampleRecord) {
// 			return &record
// 		}
// 	}

// 	return nil
// }

func getRecordsToAppend(updateRecords []dnsRecord, existingRecords []dnsRecord) []dnsRecord {
	var recordsToAppend []dnsRecord
	for _, record := range updateRecords {
		var foundRecord *dnsRecord
		if record.RecType != "MX" {
			foundRecord = findRecordByNameAndType(record.HostName, record.RecType, existingRecords)
		} else {
			foundRecord = findRecordByNameAndTypeAndPriority(record.HostName, record.RecType, record.Priority, existingRecords)
		}
		if foundRecord == nil || foundRecord.Destination != record.Destination {
			recordsToAppend = append(recordsToAppend, record)
		}
	}
	return recordsToAppend
}

func getRecordsToSet(updateRecords []dnsRecord, existingRecords []dnsRecord) []dnsRecord {
	var recordsToUpdate []dnsRecord
	var recordsToAppend []dnsRecord
	for _, record := range updateRecords {
		var foundRecord *dnsRecord
		if record.ID != "" {
			foundRecord = findRecordByID(record.ID, existingRecords)
		} else if record.RecType != "MX" {
			foundRecord = findRecordByNameAndType(record.HostName, record.RecType, existingRecords)
		} else {
			foundRecord = findRecordByNameAndTypeAndPriority(record.HostName, record.RecType, record.Priority, existingRecords)
		}

		if foundRecord != nil && !foundRecord.equals(record) {
			record.ID = foundRecord.ID
			recordsToUpdate = append(recordsToUpdate, record)
		} else if foundRecord == nil {
			recordsToAppend = append(recordsToAppend, record)
		}
	}
	return append(recordsToUpdate, recordsToAppend...)
}
