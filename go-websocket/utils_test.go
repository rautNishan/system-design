package gowebsocket

import (
	"reflect"
	"testing"
)

// callUriParser wraps uriParser and converts any panic into a normal
// test failure (with the panic value reported), instead of letting it
// crash the whole test binary. This is what lets us keep running the
// rest of the table even if one case panics.
func callUriParser(t *testing.T, input string) (result *WsUri, err error, panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	result, err = uriParser(input)
	return
}

func TestUriParser(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      *WsUri
		wantErr   bool
		wantPanic bool // set true only for inputs known to currently panic
	}{
		{
			name:  "doc comment example: ws://localhost:3000/ws?hehe",
			input: "ws://localhost:3000/ws?hehe",
			want:  &WsUri{WS: "ws", Host: "localhost", Port: "3000", Path: "ws", Params: []string{"hehe"}},
		},
		{
			name:  "wss scheme with different host/port",
			input: "wss://example.com:443/socket?token=abc",
			want:  &WsUri{WS: "wss", Host: "example.com", Port: "443", Path: "socket", Params: []string{"token=abc"}},
		},
		{
			// No "?" in the path portion at all. thirdPortion[1] (e.g. "ws")
			// gets split on "?" and produces a 1-element slice, so
			// fourthProtion[1:] is an empty-but-non-nil []string{},
			// NOT nil. reflect.DeepEqual treats nil and []string{} as
			// unequal, so this must be []string{} here, not nil.
			name:  "no query params",
			input: "ws://localhost:3000/ws",
			want:  &WsUri{WS: "ws", Host: "localhost", Port: "3000", Path: "ws", Params: []string{}},
		},
		{
			name:    "missing scheme colon entirely",
			input:   "localhost/ws",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err, panicked, panicVal := callUriParser(t, tt.input)

			if tt.wantPanic {
				if !panicked {
					t.Fatalf("expected a panic for input %q (known bug), but got result=%+v err=%v — has the bug been fixed? if so, update this test case's wantPanic/want fields", tt.input, got, err)
				}
				t.Logf("confirmed known bug: input %q panics with: %v", tt.input, panicVal)
				return
			}

			if panicked {
				t.Fatalf("unexpected panic for input %q: %v", tt.input, panicVal)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil (result=%+v)", tt.input, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("input %q\n got:  %+v (Params nil? %v)\n want: %+v (Params nil? %v)",
					tt.input, got, got.Params == nil, tt.want, tt.want.Params == nil)
			}
		})
	}
}
