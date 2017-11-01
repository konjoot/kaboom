# kaboom
tiny gRPC cli tool

## Proof-Of-Concept

- [x] gRPC Base mock server
- [x] split the functionality into three main parts (encoder, processor, decoder)
- [ ] [encoder] read from stdin, write to stdout encoded message, to stderr - logs and errors:
    - [x] use encoding/binary package for protobuf Encoding\Decoding
    - [ ] rules parser for the encoder
- [ ] [processor] read from stdin, call gRPC, write answer to stdout, to stderr - logs and errors
- [ ] [decoder] read from stdin, write to stdout decoded string, to stderr - logs and errors
- [ ] clean code
- [ ] [encoder] encoding for int64 values
- [ ] [encoder] encoding for string values
- [ ] [encoder] encoding for int32 values
- [ ] [processor] gRPC request processing
