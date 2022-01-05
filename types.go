package netcup

import (
	"encoding/json"
)

type dnsRecord struct {
	ID           string `json:"id"`
	HostName     string `json:"hostname"`
	RecType      string `json:"type"`
	Priority     int    `json:"priority,string"`
	Destination  string `json:"destination"`
	DeleteRecord bool   `json:"deleterecord"`
}

func (rec *dnsRecord) equals(otherRec dnsRecord) bool {
	return rec.HostName == otherRec.HostName && rec.RecType == otherRec.RecType && rec.Destination == otherRec.Destination && rec.Priority == otherRec.Priority
}

type dnsRecordSet struct {
	DnsRecords []dnsRecord `json:"dnsrecords"`
}

type apiSessionData struct {
	APISessionId string `json:"apisessionid"`
}

// Information about the zone. Name: the zone name, TTL: time to live in seconds
type dnsZone struct {
	Name string `json:"name"`
	TTL  int64  `json:"ttl,string"`
}

type requestParam struct {
	DomainName     string       `json:"domainname,omitempty"`
	CustomerNumber string       `json:"customernumber"`
	APIKey         string       `json:"apikey"`
	APIPassword    string       `json:"apipassword,omitempty"`
	APISessionID   string       `json:"apisessionid,omitempty"`
	DNSRecordSet   dnsRecordSet `json:"dnsrecordset,omitempty"`
}

type request struct {
	Action string       `json:"action"`
	Param  requestParam `json:"param"`
}

type response struct {
	Action       string          `json:"action"`
	Status       string          `json:"status"`
	ShortMessage string          `json:"shortmessage"`
	LongMessage  string          `json:"longmessage"`
	ResponseData json.RawMessage `json:"responsedata"`
}
