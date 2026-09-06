package tolerant_test

import (
	"symbol-invert/internal/tolerant"
	"testing"
)

func TestRecognizeSymbol(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		artFontSymbols map[rune][8]string
		symbolArt      []string
		want           rune
		want2          bool
	}{
		{
			name: "exact match A",
			artFontSymbols: map[rune][8]string{
				'A': {
					"         ",
					"  _|_|   ",
					"_|    _| ",
					"_|_|_|_| ",
					"_|    _| ",
					"_|    _| ",
					"         ",
					"         ",
				},
				'B': {
					"         ",
					"_|_|_|   ",
					"_|    _| ",
					"_|_|_|   ",
					"_|    _| ",
					"_|_|_|   ",
					"         ",
					"         ",
				},
				'C': {
					"         ",
					"  _|_|_| ",
					"_|       ",
					"_|       ",
					"_|       ",
					"  _|_|_| ",
					"         ",
					"         ",
				},
			},
			symbolArt: []string{
				"         ",
				"  _|_|   ",
				"_|    _| ",
				"_|_|_|_| ",
				"_|    _| ",
				"_|    _| ",
				"         ",
				"         ",
			},
			want:  'A',
			want2: false,
		},
		{
			name: "nearest match A with two damaged lines",
			artFontSymbols: map[rune][8]string{
				'A': {
					"         ",
					"  _|_|   ",
					"_|    _| ",
					"_|_|_|_| ",
					"_|    _| ",
					"_|    _| ",
					"         ",
					"         ",
				},
				'B': {
					"         ",
					"_|_|_|   ",
					"_|    _| ",
					"_|_|_|   ",
					"_|    _| ",
					"_|_|_|   ",
					"         ",
					"         ",
				},
				'C': {
					"         ",
					"  _|_|_| ",
					"_|       ",
					"_|       ",
					"_|       ",
					"  _|_|_| ",
					"         ",
					"         ",
				},
			},
			symbolArt: []string{
				"         ",
				"  WRONG  ",
				"_|    _| ",
				"_|_|_|_| ",
				"_|    _| ",
				"  WRONG  ",
				"         ",
				"         ",
			},
			want:  'A',
			want2: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, err := tolerant.RecognizeSymbol(tt.artFontSymbols, tt.symbolArt)
			if err != nil {
				t.Errorf("RecognizeSymbol() unexpected error: %v", err)
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("RecognizeSymbol() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("RecognizeSymbol() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestRecognizeText(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		artFontSymbols map[rune][8]string
		dividedArt     [][][]string
		want           string
		want2          int
	}{
		{
			name: "ABC with one nearest match",
			artFontSymbols: map[rune][8]string{
				'A': {
					"         ",
					"  _|_|   ",
					"_|    _| ",
					"_|_|_|_| ",
					"_|    _| ",
					"_|    _| ",
					"         ",
					"         ",
				},
				'B': {
					"         ",
					"_|_|_|   ",
					"_|    _| ",
					"_|_|_|   ",
					"_|    _| ",
					"_|_|_|   ",
					"         ",
					"         ",
				},
				'C': {
					"         ",
					"  _|_|_| ",
					"_|       ",
					"_|       ",
					"_|       ",
					"  _|_|_| ",
					"         ",
					"         ",
				},
			},
			dividedArt: [][][]string{
				{
					{
						"         ",
						"  _|_|   ",
						"_|    _| ",
						"_|_|_|_| ",
						"_|    _| ",
						"_|    _| ",
						"         ",
						"         ",
					},
					{
						"         ",
						"_|_|_|   ",
						"_|    _| ",
						"WRONG    ",
						"_|    _| ",
						"_|_|_|   ",
						"         ",
						"         ",
					},
					{
						"         ",
						"  _|_|_| ",
						"_|       ",
						"_|       ",
						"_|       ",
						"  _|_|_| ",
						"         ",
						"         ",
					},
				},
			},
			want:  "ABC",
			want2: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, err := tolerant.RecognizeText(tt.artFontSymbols, tt.dividedArt)
			if err != nil {
				t.Errorf("RecognizeText() unexpected error: %v", err)
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("RecognizeText() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("RecognizeText() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
