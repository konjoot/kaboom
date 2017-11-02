package encoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"strings"
)

// Innder proto field types
const (
	// Undefined (reserved) field type, such fields will be ignored during encoding
	Undefined uint8 = iota
	Varint
	X64Bit
	LengthDelimited
	StartGroup
	EndGroup
	X32Bit
)

// String representation of rule types
// corresponding string values available
// in rule strings
const (
	// types which encodes as Varint type
	IntString    = "int"
	Int32String  = "int32"
	Int64String  = "int64"
	UintString   = "uint"
	Uint32String = "uint32"
	Uint64String = "uint64"
	SintString   = "sint"
	Sint32String = "sint32"
	Sint64String = "sint64"
	BoolString   = "bool"

	// types which encodes as LengthDelimited types
	StringString = "string"
	BytesString  = "bytes"
)

// Rule is the interface which every rule should provide
type Rule interface {
	Number() uint8
	Name() string
	Type() uint8
}

// RuleSorter is the container for sorting slises of Rule
// implements sort.Interface
type RuleSorter []Rule

func (rs RuleSorter) Len() int           { return len(rs) }
func (rs RuleSorter) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs RuleSorter) Less(i, j int) bool { return rs[i].Number() < rs[j].Number() }

// ParseRules parses string with rules for encoder
// format: "name:type;name:type;..." e.g. "first:string;second:int64"
func ParseRules(in string) ([]Rule, error) {

	var (
		i          uint8
		err        error
		ruleString string
		rules      = make([]Rule, 0, 10)
		reader     = bufio.NewReader(strings.NewReader(in))
	)

	for i = 1; err != io.EOF; i++ {

		ruleString, err = reader.ReadString(';')
		if err != nil && err != io.EOF {
			return rules, err
		}

		ruleString = strings.TrimSuffix(ruleString, ";")

		ruleParts := strings.SplitN(ruleString, ":", 2)

		if len(ruleParts) != 2 {
			continue
		}

		rules = append(rules, &rule{fieldNumber: i, fieldName: ruleParts[0], fieldType: ruleParts[1]})
	}

	return rules, nil
}

// Encode JSON from io.Reader to inner protobuf format
// not fully implemented
func Encode(in io.Reader, rules []Rule) ([]byte, error) {
	data := make(map[string]interface{})
	jsonDecoder := json.NewDecoder(in)

	if !jsonDecoder.More() {
		// nothing to decode
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

	for _, rule := range rules {

		// skip undefined fields
		if rule.Type() == Undefined {
			continue
		}

		field = (rule.Number() << 3) | rule.Type()

		err := binary.Write(out, binary.LittleEndian, &field)
		if err != nil {
			return nil, err
		}
		err = binary.Write(out, binary.LittleEndian, data[rule.Name()])
		if err != nil {
			return nil, err
		}
	}

	return out.Bytes(), nil
}

// unexported struct which implements Rule interface
type rule struct {
	fieldNumber uint8
	fieldName   string
	fieldType   string
}

func (r *rule) Number() uint8 {
	return r.fieldNumber
}

func (r *rule) Name() string {
	return r.fieldName
}

func (r *rule) Type() uint8 {
	switch r.fieldType {
	case IntString, Int32String, Int64String,
		UintString, Uint32String, Uint64String,
		SintString, Sint32String, Sint64String,
		BoolString:
		return Varint
	case StringString,
		BytesString:
		return LengthDelimited
	}

	return Undefined
}
