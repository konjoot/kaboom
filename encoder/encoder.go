package encoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"strings"
)

// // 0000 1000||0000 0001 -> 0x08, 0x01
// // 0001 0010 0000 0111 -> 0x12 0x07
// // 0x12 0x07||0x75 0x73 0x65 0x72 0x5F 0x69 0x64
// // 0001 1010 0100 1000 -> 0x1a, 0x08
// // 0x1a, 0x08||0x73 0x63, 0x6F, 0x70, 0x65, 0x5F, 0x69, 0x64
// return []byte{
// 	0x08, 0x01, // 1|0: 1
// 	0x12, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5F, 0x69, 0x64, // 2|2: user_id
// 	0x1a, 0x08, 0x73, 0x63, 0x6F, 0x70, 0x65, 0x5F, 0x69, 0x64, // 3|2: scope_id
// }, nil

// ParseRules parses string with rules for encoder
// format: "name:type;name:type;..." e.g. "first:string;second:int64"
func ParseRules(in string) ([]*Rule, error) {
	var rules = make([]*Rule, 0, 10)

	reader := bufio.NewReader(strings.NewReader(in))

	var err error
	var i = 0

	for {
		i++

		ruleString, err := reader.ReadString(';')
		if err != nil {
			return rules, err
		}

		ruleParts := strings.SplitN(ruleString, 2)
		if len(ruleParts) != 2 {
			continue
		}

		rules = append(rules, &Rule{Number: i, Name: ruleParts[0], Type: ruleParts[1]})
	}
}

// Encode JSON from io.Reader to inner protobuf format
// not fully implemented
func Encode(in io.Reader, rules []*Rule) ([]byte, error) {
	data := make(map[string]interface{})
	jsonDecoder := json.NewDecoder(in)

	if !jsonDecoder.More() {
		// nothing to encode
		return []byte{}, nil
	}

	if err := jsonDecoder.Decode(data); err != nil {
		// something went wrong during JSON decoding
		return nil, err
	}

	if len(data) < 1 {
		// nothing to encode
		return []byte{}, nil
	}

	out := bytes.NewBuffer(make([]byte, 0, 0))
	var field uint8

	for rule := range rules {

		field = (rule.Number() << 3) | rule.Type()

		err := binary.Write(out, binary.LittleEndian, &field)
		if err != nil {
			return nil, err
		}
		err := binary.Write(out, binary.LittleEndian, &data[rule.Name()])
		if err != nil {
			return nil, err
		}
	}
}

type Rule struct {
	fieldNumber uint8
	fieldName   string
	fieldType   string
}

func (r *Rule) Number() uint8 {
	return 0
}

func (r *Rule) Name() string {
	return ""
}

func (r *Rule) Type() uint8 {
	return 0
}

type FieldType int

const (
	Varing FieldType = iota
	X64Bit
	LengthDelimited
	StartGroup
	EndGroup
	X32Bit
)

type OriginFieldType int

const (
	// Values, stored as Varint data
	Int32 OriginFieldType = iota
	Int64
	Uint32
	Uint64
	Sint32
	Sint64
	Bool
	Enum

	// Values, stored as LengthDelimited data
	String
	Bytes
	EmbeddedMessages
	PackedRepeatedFields
)
