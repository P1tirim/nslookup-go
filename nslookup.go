package nslookup

import (
	"encoding/binary"
	"net"
)

// Return response from DNS server in DNS message protocol format. Default server is 8.8.8.8:53.
func (q *QueryDNS) Lookup(server string) (*Response, error) {
	if len(q.Queries) == 0 {
		return nil, ErrEmptyQueries
	}

	if server == "" {
		server = serverDefault
	}

	conn, err := net.Dial("udp", server)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	payload := make([]byte, 0)
	payload = binary.BigEndian.AppendUint16(payload, q.TransactionID)
	payload = binary.BigEndian.AppendUint16(payload, q.Flags)
	payload = binary.BigEndian.AppendUint16(payload, q.QuestionsCount)
	payload = binary.BigEndian.AppendUint16(payload, q.AnswersCount)
	payload = binary.BigEndian.AppendUint16(payload, q.AuthorityCount)
	payload = binary.BigEndian.AppendUint16(payload, q.AdditionalsCount)
	payload = parseDomain(payload, q.Queries[0].Name)
	payload = binary.BigEndian.AppendUint16(payload, q.Queries[0].Type)
	payload = binary.BigEndian.AppendUint16(payload, q.Queries[0].Class)

	_, err = conn.Write(payload)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return parseResponse(buf[:n], len(payload), q.Queries)
}

// Return array of IPv4 and IPv6 of this domain.
func LookupIP(domain, server string) (ips []net.IP, err error) {
	resp, err := LookupIPv4(domain, server)
	if err != ErrNoAnswer {
		return nil, err
	}

	ips = append(ips, resp...)

	resp, err = LookupIPv6(domain, server)
	if err != ErrNoAnswer {
		return nil, err
	}

	ips = append(ips, resp...)

	return ips, nil
}

// Return array of TEXTs of this domain.
func LookupTEXT(domain, server string) (texts []string, err error) {
	query := NewQueryDNS(domain, TypeTXT)

	resp, err := query.Lookup(server)
	if err != nil {
		return nil, err
	}

	for _, answer := range resp.Answers {
		if v, ok := answer.Data.([]AnswerTypeTXT); ok {
			for i := 0; i < len(v); i++ {
				texts = append(texts, v[0].Text)
			}
		}
	}

	return texts, nil
}

// Return array of CNAMEs of this domain.
func LookupCNAME(domain, server string) (cnames []string, err error) {
	return lookupTypeString(domain, server, TypeCNAME)
}

func LookupIPv4(domain, server string) (ips []net.IP, err error) {
	query := NewQueryDNS(domain, TypeA)

	resp, err := query.Lookup(server)
	if err != nil {
		return nil, err
	}

	for _, answer := range resp.Answers {
		if answer.Type != TypeA {
			continue
		}

		if v, ok := answer.Data.(AnswerTypeString); ok {
			ips = append(ips, net.ParseIP(v.Data))
		}
	}

	return ips, nil
}

func LookupIPv6(domain, server string) (ips []net.IP, err error) {
	query := NewQueryDNS(domain, TypeAAAA)

	resp, err := query.Lookup(server)
	if err != nil {
		return nil, err
	}

	for _, answer := range resp.Answers {
		if answer.Type != TypeAAAA {
			continue
		}

		if v, ok := answer.Data.(AnswerTypeString); ok {
			ips = append(ips, net.ParseIP(v.Data))
		}
	}

	return ips, nil
}

func LookupMX(domain, server string) (mx []AnswerTypeMX, err error) {
	query := NewQueryDNS(domain, TypeMX)

	resp, err := query.Lookup(server)
	if err != nil {
		return nil, err
	}

	for _, answer := range resp.Answers {
		if answer.Type != TypeMX {
			continue
		}

		if v, ok := answer.Data.(AnswerTypeMX); ok {
			mx = append(mx, v)
		}
	}

	return mx, nil
}

func LookupPTR(ip net.IP, server string) (ptr []string, err error) {
	if ip == nil {
		return nil, ErrNotValidIP
	}

	return lookupTypeString(reverseIPAddress(ip), server, TypePTR)
}

func LookupNS(domain, server string) (ns []string, err error) {
	return lookupTypeString(domain, server, TypeNS)
}

func lookupTypeString(domain, server string, dnsType uint16) (arr []string, err error) {
	query := NewQueryDNS(domain, dnsType)

	resp, err := query.Lookup(server)
	if err != nil {
		return nil, err
	}

	for _, answer := range resp.Answers {
		if v, ok := answer.Data.(AnswerTypeString); ok {
			arr = append(arr, v.Data)
		}
	}

	return arr, nil
}
