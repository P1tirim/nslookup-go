package main

import (
	"fmt"
	"log"

	"github.com/P1tirim/nslookup-go"
)

func main() {
	ips, err := nslookup.LookupIP("google.com", "")
	if err != nil {
		log.Fatal(err)
	}

	for _, ip := range ips {
		fmt.Println("IP: " + ip.String())
	}

	texts, err := nslookup.LookupTEXT("google.com", "")
	if err != nil {
		log.Fatal(err)
	}

	for _, text := range texts {
		fmt.Println("text: " + text)
	}

	cnames, err := nslookup.LookupCNAME("time.apple.com", "")
	if err != nil {
		log.Fatal(err)
	}

	for _, cname := range cnames {
		fmt.Println("CNAME: " + cname)
	}

	// Advande usage
	query := nslookup.NewQueryDNS("google.com", nslookup.TypeA)

	resp, err := query.Lookup("8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}

	for _, answer := range resp.Answers {
		if answer.Type == nslookup.TypeA {
			fmt.Println("IPv4: " + answer.Data.(nslookup.AnswerTypeString).Data)
		}
	}

	query = nslookup.NewQueryDNS("google.com", nslookup.TypeAAAA)

	resp, err = query.Lookup("")
	if err != nil {
		log.Fatal(err)
	}

	for _, answer := range resp.Answers {
		if answer.Type == nslookup.TypeAAAA {
			fmt.Println("IPv6: " + answer.Data.(nslookup.AnswerTypeString).Data)
		}
	}

	query = nslookup.NewQueryDNS("google.com", nslookup.TypeTXT)

	resp, err = query.Lookup("")
	if err != nil {
		log.Fatal(err)
	}

	for _, answer := range resp.Answers {
		if answer.Type == nslookup.TypeTXT {
			t := answer.Data.([]nslookup.AnswerTypeTXT)

			for i := 0; i < len(t); i++ {
				fmt.Println("text: " + t[i].Text)
			}
		}
	}
}
