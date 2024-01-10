# nslookup-go

This library let you do nslookup by UDP connection. You can choose the server, what you want to use or use default 8.8.8.8 server.

## Supported query types

```
TypeA = 1 //IPv4
TypeNS = 2
TypeCNAME = 5
TypePTR = 12
TypeMX = 15
TypeTXT = 16
TyepAAAA= 28 //IPv6
```

## Example usage

``` go 
package main

import (
	"fmt"
	"log"

	"github.com/P1tirim/nslookup-go"
)

func main(){
    // Accept domain as first parametr and ip of DNS server as second.
    // You can not passing DNS server, then will be use standart 8.8.8.8 server
    // Return []net.IP
    ips, err := nslookup.LookupIP("google.com", "")
	if err != nil {
		log.Fatal(err)
	}

	for _, ip := range ips {
		fmt.Println("IP: " + ip.String())
	}
}
```

Another expamples you can find in the folder "examples"

## Advanced usage

If you want to get answer from DNS server by DNS protocol message, then this exapmle for you

``` go
package main

import (
	"fmt"
	"log"

	"github.com/P1tirim/nslookup-go"
)

func main(){
    query := nslookup.NewQueryDNS("google.com", nslookup.TypeA)

    // This function return nslookup.Response
	resp, err := query.Lookup("8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}

    // You need to assert every answer to struct
	for _, answer := range resp.Answers {
		if answer.Type == nslookup.TypeA {
			fmt.Println("IPv4: " + answer.Data.(nslookup.AnswerTypeString).Data)
		}
	}
}
```

The struct, what returns DNS server
``` go
type Response struct {
	TransactionID     uint16
	Flags             uint16
	QuestionsCount    uint16
	AnswersCount      uint16
	AuthorityCounts   uint16
	AdditionalCounts  uint16
	Queries           []Query
	Answers           []Answer
	AuthorityRecords  []Answer
	AdditionalRecords []Answer
}
```