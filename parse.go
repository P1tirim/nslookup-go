package nslookup

import (
	"encoding/binary"
	"errors"
	"strings"
)

var (
	ErrNoSuchName              error = errors.New("no such name")
	ErrNoAnswer                      = errors.New("no answer")
	ErrInvalidAnswerFromServer       = errors.New("invalid answer from server")
	ErrUnsupportedDNSType            = errors.New("unsupported dns type")
	ErrEmptyQueries                  = errors.New("empty queries in QueryDNS")
	ErrFormatQuery                   = errors.New("the query was constructed incorrectly")
	ErrInternalServerDNS             = errors.New("internal error in the DNS server")
	ErrRefused                       = errors.New("the query refused by DNS server")
	ErrYXDomain                      = errors.New("name exists when it should not")
	ErrYXRRSet                       = errors.New("RR Set Exists when it should not")
	ErrNXRRSet                       = errors.New("RR Set that should exist does not")
	ErrNotAuth                       = errors.New("not authorized")
)

const (
	// IPv4
	TypeA     = 1
	TypeCNAME = 5
	TypeTXT   = 16
	// IPv6
	TypeAAAA = 28
)

const serverDefault = "8.8.8.8:53"

func parseDomain(payload []byte, domain string) []byte {
	a := strings.Split(domain, ".")

	for _, v := range a {
		payload = append(payload, byte(len(v)))
		payload = append(payload, []byte(v)...)
	}

	payload = append(payload, 0x00)

	return payload
}

func parseResponse(response []byte, requestLength int, queries []Query) (resp *Response, err error) {
	if len(response) < 11 {
		return nil, ErrInvalidAnswerFromServer
	}

	resp = &Response{
		TransactionID:     binary.BigEndian.Uint16(response[0:2]),
		Flags:             binary.BigEndian.Uint16(response[2:4]),
		QuestionsCount:    binary.BigEndian.Uint16(response[4:6]),
		AnswersCount:      binary.BigEndian.Uint16(response[6:8]),
		AuthorityCounts:   binary.BigEndian.Uint16(response[8:10]),
		AdditionalCounts:  binary.BigEndian.Uint16(response[10:12]),
		Queries:           queries,
		Answers:           make([]Answer, 0),
		AuthorityRecords:  make([]Answer, 0),
		AdditionalRecords: make([]Answer, 0),
	}

	switch resp.Flags {
	case 0x8181:
		return nil, ErrFormatQuery
	case 0x8182:
		return nil, ErrInternalServerDNS
	case 0x8183:
		return nil, ErrNoSuchName
	case 0x8184:
		return nil, ErrUnsupportedDNSType
	case 0x8185:
		return nil, ErrRefused
	case 0x8186:
		return nil, ErrYXDomain
	case 0x8187:
		return nil, ErrYXRRSet
	case 0x8188:
		return nil, ErrNXRRSet
	case 0x8189:
		return nil, ErrNotAuth
	}

	response = response[requestLength:]
	if len(response) == 0 || resp.AnswersCount == 0 {
		return nil, ErrNoAnswer
	}

	if response[0] != 0xc0 {
		response = removeStartWrongBytes(response)
	}

	if resp.AnswersCount != 0 {
		response, resp.Answers, err = parseAnswer(response, int(resp.AnswersCount))
		if err != nil {
			return nil, err
		}
	}

	if resp.AuthorityCounts != 0 {
		response, resp.AuthorityRecords, err = parseAnswer(response, int(resp.AuthorityCounts))
		if err != nil {
			return nil, err
		}
	}

	if resp.AdditionalCounts != 0 {
		_, resp.AdditionalRecords, err = parseAnswer(response, int(resp.AdditionalCounts))
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func parseAnswer(response []byte, count int) ([]byte, []Answer, error) {
	answers := make([]Answer, 0)

	for i := 0; i < count; i++ {
		if len(response) < 11 {
			return nil, nil, ErrInvalidAnswerFromServer
		}

		answer := Answer{
			Name:       binary.BigEndian.Uint16(response[0:2]),
			Type:       binary.BigEndian.Uint16(response[2:4]),
			Class:      binary.BigEndian.Uint16(response[4:6]),
			TimeToLive: binary.BigEndian.Uint32(response[6:10]),
			DataLength: binary.BigEndian.Uint16(response[10:12]),
		}

		response = response[12:]

		var err error

		switch answer.Type {
		case TypeA:
			err = answer.parseTypeA(response)
		case TypeTXT:
			err = answer.parseTypeTXT(response)
		case TypeCNAME:
			err = answer.parseTypeCNAME(response)
		case TypeAAAA:
			err = answer.parseTypeAAAA(response)
		default:
			return nil, nil, ErrUnsupportedDNSType
		}

		if err != nil {
			return nil, nil, err
		}

		response = response[answer.DataLength:]
		answers = append(answers, answer)
	}

	return response, answers, nil
}

// Answer should be start form 0xc0 byte. This function
// removes bytes, which change standart array length
func removeStartWrongBytes(data []byte) []byte {
	for i, v := range data {
		if v == 0xc0 {
			return data[i:]
		}
	}

	return data
}
