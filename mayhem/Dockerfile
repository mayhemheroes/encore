FROM golang:1.23 AS builder

WORKDIR /src
COPY . /src/encore

# Copy fuzz test into the package directory
COPY mayhem/fuzz_test.go /src/encore/runtime/types/uuid/fuzz_native_test.go

WORKDIR /src/encore/runtime/types/uuid

# Compile native Go fuzz test as libfuzzer binary
RUN go test -c -fuzz=FuzzFromString -o /fuzz.libfuzzer .

FROM ubuntu:22.04
COPY --from=builder /fuzz.libfuzzer /fuzz.libfuzzer
