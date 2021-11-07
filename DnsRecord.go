package main

type DnsRecord struct {
	Id      string `json:"id"`
	ZoneId  string `json:"zone_id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func CreateDnsRecordFrom(entry DnsRecordResponseEntry) *DnsRecord {
	return &DnsRecord{
		Id:      entry.Id,
		ZoneId:  entry.ZoneId,
		Type:    entry.Type,
		Name:    entry.Name,
		Content: entry.Content,
	}
}
