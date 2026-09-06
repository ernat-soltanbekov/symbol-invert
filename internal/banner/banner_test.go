package banner

import (
	"testing"
)

func TestBanner(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "загрузка стандартного шрифта",
			filename: "../../standard.txt",
			wantErr:  false,
		},
		{
			name:     "загрузка шрифта shadow",
			filename: "../../shadow.txt",
			wantErr:  false,
		},
		{
			name:     "загрузка шрифта thinkertoy",
			filename: "../../thinkertoy.txt",
			wantErr:  false,
		},
		{
			name:     "ошибка при отсутствии файла шрифта",
			filename: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Banner(tt.filename)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("Banner() failed: %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("Banner() succeeded unexpectedly")
			}

			if len(got) == 0 {
				t.Error("Banner() returned empty result")
			}
		})
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		word     string
		wantErr  bool
	}{
		{
			name:     "рендеринг со стандартным шрифтом",
			filename: "../../standard.txt",
			word:     "Hello",
			wantErr:  false,
		},
		{
			name:     "рендеринг со шрифтом shadow",
			filename: "../../shadow.txt",
			word:     "Hello",
			wantErr:  false,
		},
		{
			name:     "рендеринг со шрифтом thinkertoy",
			filename: "../../thinkertoy.txt",
			word:     "Hello",
			wantErr:  false,
		},
		{
			name:     "ошибка при отсутствии файла шрифта",
			filename: ".",
			word:     "Hello",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.filename, tt.word)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("Render() failed: %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("Render() succeeded unexpectedly")
			}

			if len(got) != 8 {
				t.Errorf("Render() returned %d lines, want 8", len(got))
			}
		})
	}
}
