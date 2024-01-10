package nslookup

import (
	"encoding/binary"
	"net"
	"net/netip"
)

// TypeA, TypeAAAA, TypeCNAME, TypeNS, TypePTR
type AnswerTypeString struct {
	Data string
}

type AnswerTypeTXT struct {
	TextLength int
	Text       string
}

type AnswerTypeMX struct {
	Preference   int
	MailExchange string
}

func (a *Answer) parseTypeA(answers []byte) error {
	if len(answers) < 4 {
		return ErrInvalidAnswerFromServer
	}

	a.Data = AnswerTypeString{
		Data: net.IPv4(answers[0], answers[1], answers[2], answers[3]).String(),
	}

	return nil
}

func (a *Answer) parseTypeTXT(answers []byte) error {
	if len(answers) < 1 {
		return ErrInvalidAnswerFromServer
	}

	texts := make([]AnswerTypeTXT, 0)
	pointer := 0

	for {
		txtLength := int(answers[0])

		if len(answers) < int(txtLength)+1 {
			return ErrInvalidAnswerFromServer
		}

		texts = append(texts, AnswerTypeTXT{TextLength: txtLength, Text: string(answers[1 : txtLength+1])})
		answers = answers[txtLength+1:]

		pointer += txtLength + 1
		if pointer == int(a.DataLength) {
			break
		}
	}

	a.Data = texts

	return nil
}

// TypeCNAME, TypeNS, TypePTR.
func (a *Answer) parseTypeWithDomain(answers, originalAnswer []byte, domains map[int]string) error {
	if len(answers) < int(a.DataLength) {
		return ErrInvalidAnswerFromServer
	}

	name := parseAnswerDomain(answers, originalAnswer, domains)

	a.Data = AnswerTypeString{
		Data: name[:len(name)-1],
	}

	return nil
}

func (a *Answer) parseTypeAAAA(answers []byte) error {
	if len(answers) < 16 {
		return ErrInvalidAnswerFromServer
	}

	a.Data = AnswerTypeString{
		Data: netip.AddrFrom16([16]byte(answers)).String(),
	}

	return nil
}

func (a *Answer) parseTypeMX(answers, originalAnswer []byte, domains map[int]string) error {
	if len(answers) < int(a.DataLength) {
		return ErrInvalidAnswerFromServer
	}

	answer := AnswerTypeMX{
		Preference:   int(binary.BigEndian.Uint16(answers[:2])),
		MailExchange: parseAnswerDomain(answers[2:], originalAnswer, domains),
	}

	a.Data = answer

	return nil
}

func parseAnswerDomain(arr []byte, originalAnswer []byte, domains map[int]string) string {
	pointer := 0
	name := ""

	for {
		if arr[pointer] == 0xc0 {
			shiftPointer := int(arr[pointer+1])

			if v, ok := domains[shiftPointer]; ok {
				name += v
			} else {
				domain := parseAnswerDomain(originalAnswer[shiftPointer:], originalAnswer, domains)
				name += domain
				domains[shiftPointer] = domain
			}

			break
		}

		length := int(arr[pointer])
		name += string(arr[pointer+1:pointer+length+1]) + "."
		pointer += length + 1

		if arr[pointer] == 0x0 {
			break
		}
	}

	return name
}
