package nslookup

import "math/rand"

type QueryDNS struct {
	TransactionID    uint16
	Flags            uint16
	QuestionsCount   uint16
	AnswersCount     uint16
	AuthorityCount   uint16
	AdditionalsCount uint16
	Queries          []Query
}

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

type Query struct {
	Name  string
	Type  uint16
	Class uint16
}

type Answer struct {
	Name       uint16
	Type       uint16
	Class      uint16
	TimeToLive uint32
	DataLength uint16
	Data       interface{}
}

func NewQueryDNS(name string, dnsType uint16) *QueryDNS {
	return &QueryDNS{
		TransactionID:    uint16(rand.Uint32()),
		Flags:            0x0100, // Standard query
		QuestionsCount:   1,
		AnswersCount:     0,
		AuthorityCount:   0,
		AdditionalsCount: 0,
		Queries:          []Query{{Name: name, Type: dnsType, Class: 1}},
	}
}
