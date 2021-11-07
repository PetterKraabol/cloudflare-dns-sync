package main

import "time"

type DnsRecordResponseEntry struct {
	Id        string `json:"id"`
	ZoneId    string `json:"zone_id"`
	ZoneName  string `json:"zone_name"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Proxiable bool   `json:"proxiable"`
	Proxied   bool   `json:"proxied"`
	Ttl       int    `json:"ttl"`
	Locked    bool   `json:"locked"`
	Meta      struct {
		AutoAdded           bool   `json:"auto_added"`
		ManagedByApps       bool   `json:"managed_by_apps"`
		ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
		Source              string `json:"source"`
	} `json:"meta"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
}
