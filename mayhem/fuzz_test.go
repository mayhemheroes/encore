//go:build !gofuzz

package uuid

import "testing"

func FuzzFromString(f *testing.F) {
	f.Add([]byte("550e8400-e29b-41d4-a716-446655440000"))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = FromString(string(data))
	})
}
