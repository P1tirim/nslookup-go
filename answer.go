package main

import (
	"net"
	"net/netip"
)

// TypeA, TypeAAAA, CNAME
type AnswerTypeString struct {
	Data string
}

type AnswerTypeTXT struct {
	TextLength int
	Text       string
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

func (a *Answer) parseTypeCNAME(answers []byte) error {
	if len(answers) < int(a.DataLength) {
		return ErrInvalidAnswerFromServer
	}

	pointer := 0
	cname := ""

	for {
		length := int(answers[pointer])
		cname += string(answers[pointer+1:pointer+length+1]) + "."
		pointer += length + 1

		if pointer == int(a.DataLength) || answers[pointer] == 0xc0 {
			break
		}
	}

	a.Data = AnswerTypeString{
		Data: cname[:len(cname)-1],
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
