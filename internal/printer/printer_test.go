package printer

import (
	"bytes"
	"testing"

	"symbol-invert/internal/banner"
)

func TestPrintLines(t *testing.T) {
	tests := []struct {
		name     string
		normArgs []string
		font     string
	}{
		{
			name:     "пустой ввод",
			normArgs: []string{""},
			font:     "../../standard.txt",
		},
		{
			name:     "слово hello",
			normArgs: []string{"hello"},
			font:     "../../standard.txt",
		},
		{
			name:     "смешанный регистр",
			normArgs: []string{"HeLlo HuMaN"},
			font:     "../../standard.txt",
		},
		{
			name:     "цифры и пробел",
			normArgs: []string{"1Hello 2There"},
			font:     "../../standard.txt",
		},
		{
			name:     "две строки",
			normArgs: []string{"Hello", "There"},
			font:     "../../standard.txt",
		},
		{
			name:     "пустая строка между строками",
			normArgs: []string{"Hello", "", "There"},
			font:     "../../standard.txt",
		},
		{
			name:     "специальные символы",
			normArgs: []string{"{Hello & There #}"},
			font:     "../../standard.txt",
		},
		{
			name:     "алфавит",
			normArgs: []string{"ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
			font:     "../../standard.txt",
		},
		{
			name:     "стандартный шрифт",
			normArgs: []string{"Hello"},
			font:     "../../standard.txt",
		},
		{
			name:     "шрифт shadow",
			normArgs: []string{"Hello"},
			font:     "../../shadow.txt",
		},
		{
			name:     "шрифт thinkertoy",
			normArgs: []string{"Hello"},
			font:     "../../thinkertoy.txt",
		},
		{
			name:     "несуществующий шрифт",
			normArgs: []string{"Hello"},
			font:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer

			PrintLines(&got, tt.normArgs, tt.font)

			var want bytes.Buffer

			for _, arg := range tt.normArgs {
				if arg == "" {
					want.WriteByte('\n')
					continue
				}

				result, err := banner.Render(tt.font, arg)
				if err != nil {
					want.WriteString("Error: ")
					want.WriteString(err.Error())
					want.WriteByte('\n')
					break
				}

				for _, line := range result {
					want.WriteString(line)
					want.WriteByte('\n')
				}
			}

			if got.String() != want.String() {
				t.Errorf(
					"PrintLines() = %q, want %q",
					got.String(),
					want.String(),
				)
			}
		})
	}
}
