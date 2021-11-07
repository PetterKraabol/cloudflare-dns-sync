package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	A          string = "A"
	AAAA              = "AAAA"
	CLOUDFLARE        = "https://api.cloudflare.com/client/v4/zones/"
)

func main() {
	zoneId := *flag.String("zone", os.Getenv("CLOUDFLARE_ZONE_ID"), "Cloudflare Zone ID")
	email := *flag.String("email", os.Getenv("CLOUDFLARE_EMAIL"), "Cloudflare email address")
	key := *flag.String("key", os.Getenv("CLOUDFLARE_AUTH_KEY"), "Cloudflare auth key")
	dnsNamesToSync := strings.Split(*flag.String("names", os.Getenv("CLOUDFLARE_SYNC_NAMES"), "Comma-separated DNS names to sync"), ",")

	ipAddresses, err := getExternalIpAddresses()
	if err != nil {
		panic(err)
	}

	// Get dns records
	dnsRecordsResponses, err := getDnsRecords(zoneId, email, key)
	if err != nil {
		panic(err)
	}

	for _, dnsRecordResponse := range dnsRecordsResponses {
		dnsRecord := CreateDnsRecordFrom(dnsRecordResponse)

		// Filter out names not to update
		if !contains(dnsNamesToSync, dnsRecord.Name) {
			continue
		}

		// DNS content is already the external ip address
		if currentContent, ok := ipAddresses[dnsRecord.Type]; !ok || currentContent == dnsRecord.Content {
			continue
		}

		fmt.Println(dnsRecord.Type, dnsRecord.Name, dnsRecord.Content, "->", ipAddresses[dnsRecord.Type])

		dnsRecord.Content = ipAddresses[dnsRecord.Type]

		if err := updateDnsRecord(*dnsRecord, email, key); err != nil {
			panic(err)
		}
	}

}

func updateDnsRecord(record DnsRecord, email string, key string) error {
	data, err := json.Marshal(map[string]string{
		"content": record.Content,
	})

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPatch, CLOUDFLARE+record.ZoneId+"/dns_records/"+record.Id, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	request.Header = http.Header{
		"x-auth-email": []string{email},
		"x-auth-key":   []string{key},
		"Content-Type": []string{"application/json"},
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Could not update DNS record " + string(data) + "\n" +
			"Response status: " + response.Status + "\n" +
			string(bodyBytes))
	}

	return nil
}

func getExternalIpAddresses() (map[string]string, error) {
	ipv4, err := getExternalIpAddress("ipv4")
	if err != nil {
		return nil, err
	}

	ipv6, err := getExternalIpAddress("ipv6")
	if err != nil {
		return nil, err
	}

	return map[string]string{
		A:    ipv4,
		AAAA: ipv6,
	}, nil
}

func getExternalIpAddress(version string) (string, error) {
	response, err := http.Get("https://" + version + ".icanhazip.com/")
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		return strings.TrimSpace(string(bodyBytes)), nil
	}

	return "", err
}

func getDnsRecords(zoneId string, email string, key string) ([]DnsRecordResponseEntry, error) {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, CLOUDFLARE+zoneId+"/dns_records", nil)
	if err != nil {
		return nil, err
	}

	request.Header = http.Header{
		"x-auth-email": []string{email},
		"x-auth-key":   []string{key},
		"Content-Type": []string{"application/json"},
	}

	response, err := client.Do(request)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	var dnsRecordsResponse DnsRecordsResponse
	if err := json.Unmarshal(body, &dnsRecordsResponse); err != nil {
		return nil, err
	}

	return dnsRecordsResponse.Result, nil
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}
