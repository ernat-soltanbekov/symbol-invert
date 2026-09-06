package reader_test

import (
	"slices"
	"symbol-invert/internal/reader"
	"testing"
)

func TestReadLines(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		filename string
		want     []string
		wantErr  bool
	}{
		{
			name:     "read file with several lines",
			filename: "../../testdata/reader_test.txt",
			want: []string{
				"hello",
				"world",
				"123",
			},
			wantErr: false,
		},
		{
			name:     "read empty file",
			filename: "../../testdata/empty.txt",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "file does not exist",
			filename: "../../testdata/not_exists.txt",
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := reader.ReadLines(tt.filename)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReadLines() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("ReadLines() succeeded unexpectedly")
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("ReadLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
