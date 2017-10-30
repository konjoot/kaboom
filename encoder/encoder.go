package encoder

import "io"
import "encoding/json"

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

// Encode JSON from io.Reader to inner protobuf format
func Encode(in io.Reader, rules []FieldType) ([]byte, error) {
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

	out := make([]byte, 0, 0)
	for i := 0; val := range data; i++ {
		var encodedFieldName uint8
		switch rules[i] {
		case VARINT:
			encodedFieldName = i
			encodedFieldName << 3
			encodedFieldName |= VARINT
		case X64BIT:
		case LENGTH_DELIMITED:
		case START_GROUP:
		case END_GROUP:
		case X32BIT:
		default: 
			continue
		}
		// out := append(out, )
	}
}

type FieldType int

const (
	VARINT FieldType = iota
	X64BIT
	LENGTH_DELIMITED
	START_GROUP
	END_GROUP
	X32BIT
)
