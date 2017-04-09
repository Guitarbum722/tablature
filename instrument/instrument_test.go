package instrument

import (
	"testing"
)

const targetTestVersion = 1

func TestTestVersion(t *testing.T) {
	if testVersion != targetTestVersion {
		t.Fatalf("Found testVersion = %v, want %v.", testVersion, targetTestVersion)
	}
}

// TestNewInstrument confirms that the default values for each test case are appropriate
func TestNewInstrument(t *testing.T) {
	for _, tt := range newInstrumentCases {
		got := NewInstrument(tt.inputType)
		fb := got.Fretboard()
		for k := range fb {
			if _, ok := fb[k]; !ok {
				t.Fatalf("NewInstrument(%q) = %q which is not a valid default value.", tt.inputType, fb[k])
			}
		}

		gotOrder := got.Order()
		wantOrder := tt.expectedType.Order()
		for i := range gotOrder {
			if gotOrder[i] != wantOrder[i] {
				t.Fatalf("Order() = %q ::: got %q want %q", gotOrder, gotOrder[i], wantOrder[i])
			}
		}
	}
}