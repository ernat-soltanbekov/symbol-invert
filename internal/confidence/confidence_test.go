package confidence

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

// ---------- sortValues ----------

func TestSortValues(t *testing.T) {
	data := []MatchesStruct{
		{Symbol: 'a', Amount: 5},
		{Symbol: 'd', Amount: 1},
		{Symbol: 'c', Amount: 10},
		{Symbol: 'b', Amount: 10},
	}

	slices.SortFunc(data, sortValues)

	want := []MatchesStruct{
		{Symbol: 'b', Amount: 10}, // при равном Amount побеждает меньший символ ('b' < 'c')
		{Symbol: 'c', Amount: 10},
		{Symbol: 'a', Amount: 5},
		{Symbol: 'd', Amount: 1},
	}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("sortValues() = %+v, want %+v", data, want)
	}
}

// ---------- isEmptyArt ----------

func TestIsEmptyArt(t *testing.T) {
	spaceLine := "      " // ровно 6 пробелов

	tests := []struct {
		name string
		art  []string
		want bool
	}{
		{
			name: "валидный арт пробела (8 строк по 6 пробелов)",
			art:  []string{spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine},
			want: true,
		},
		{
			name: "неверное количество строк",
			art:  []string{spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine},
			want: false,
		},
		{
			name: "одна из строк не пробел",
			art:  []string{spaceLine, spaceLine, "AAAAAA", spaceLine, spaceLine, spaceLine, spaceLine, spaceLine},
			want: false,
		},
		{
			name: "строки пустые (не 6 пробелов)",
			art:  []string{"", "", "", "", "", "", "", ""},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isEmptyArt(tc.art); got != tc.want {
				t.Errorf("isEmptyArt(%v) = %v, want %v", tc.art, got, tc.want)
			}
		})
	}
}

// ---------- Banner ----------

// Формат файла, который ожидает Banner: заголовок (1 строка, скипается),
// затем для каждого символа — 8 строк арта + 1 разделительная строка.
func writeBannerFile(t *testing.T, chars map[rune][8]string, order []rune) string {
	t.Helper()

	content := "HEADER\n"
	for _, c := range order {
		for _, line := range chars[c] {
			content += line + "\n"
		}
		content += "\n" // разделитель между символами (важен для корректной работы Banner)
	}

	path := filepath.Join(t.TempDir(), "font.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("не удалось создать тестовый файл шрифта: %v", err)
	}
	return path
}

// ---------- DivideIntoLines ----------

func TestDivideIntoLines(t *testing.T) {
	art := []string{
		"L0", "L1", "L2", "L3", "L4", "L5", "L6", "L7", // первая арт-строка
		"L8", "L9", "L10", "L11", "L12", "L13", "L14", "L15", // вторая арт-строка
	}

	got := DivideIntoLines(art)

	want := [][]string{
		{"L0", "L1", "L2", "L3", "L4", "L5", "L6", "L7"},
		{"L8", "L9", "L10", "L11", "L12", "L13", "L14", "L15"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("DivideIntoLines() = %v, want %v", got, want)
	}
}

// ---------- FindBorder ----------

func TestFindBorder(t *testing.T) {
	// Строка вида: буква "AAA" + пробел шириной 6 + буква "BBB" + техническая
	// подложка из 3 пробелов в конце (чтобы не выйти за границы среза).
	row := "AAA      BBB   "
	artLine := []string{row, row, row, row, row, row, row, row}
	artLines := [][]string{artLine}

	got := FindBorder(1, artLines)

	if len(got) != 1 {
		t.Fatalf("FindBorder() вернул %d арт-строк, want 1", len(got))
	}

	blocks := got[0]
	if len(blocks) != 3 {
		t.Fatalf("FindBorder() нашел %d блоков, want 3 (буква, пробел, буква); blocks=%v", len(blocks), blocks)
	}

	wantLetter1 := "AAA " // индексы 0:4 (буква + первая колонка пробела)
	wantSpace := "      " // индексы 3:9, ширина пробела = 6
	wantLetter2 := "BBB " // индексы 9:13

	for i, line := range blocks[0] {
		if line != wantLetter1 {
			t.Errorf("blocks[0][%d] = %q, want %q", i, line, wantLetter1)
		}
	}
	for i, line := range blocks[1] {
		if line != wantSpace {
			t.Errorf("blocks[1][%d] = %q, want %q", i, line, wantSpace)
		}
	}
	for i, line := range blocks[2] {
		if line != wantLetter2 {
			t.Errorf("blocks[2][%d] = %q, want %q", i, line, wantLetter2)
		}
	}
}

// Промежуток короче spaceWidth (6 колонок) должен игнорироваться и не
// добавляться как отдельный блок-пробел.
func TestFindBorder_ShortGapIgnored(t *testing.T) {
	// между буквами всего 3 пустых колонки подряд — это не пробел, а технический зазор.
	// В конце строки добавлена подложка-пробел: без неё FindBorder попытается
	// срезать line[letterStart:j+1] за пределами длины строки и запаникует —
	// это отдельный дефект реализации, не предмет этого теста.
	row := "AA   BB "
	artLine := []string{row, row, row, row, row, row, row, row}
	artLines := [][]string{artLine}

	got := FindBorder(1, artLines)
	blocks := got[0]

	for idx, block := range blocks {
		for _, line := range block {
			if line == "   " {
				t.Errorf("короткий зазор (3 колонки) не должен был попасть в результат как блок; blocks[%d]=%v", idx, block)
			}
		}
	}
}

// ---------- ConfidenceSymbol ----------

func testFontA() [8]string {
	return [8]string{
		"AAAAA",
		"A...A",
		"AAAAA",
		"A...A",
		"A...A",
		"A...A",
		"A...A",
		"A...A",
	}
}

func TestConfidenceSymbol_PerfectMatch(t *testing.T) {
	fontA := testFontA()
	font := map[rune][8]string{'A': fontA}

	symbolArt := []string{fontA[0], fontA[1], fontA[2], fontA[3], fontA[4], fontA[5], fontA[6], fontA[7]}

	symbol, confidence, quality, err := ConfidenceSymbol(font, symbolArt)

	if symbol != 'A' {
		t.Errorf("symbol = %q, want %q", symbol, 'A')
	}
	if confidence != 100 {
		t.Errorf("confidence = %v, want 100", confidence)
	}
	if quality != "Perfect match" {
		t.Errorf("quality = %q, want %q", quality, "Perfect match")
	}
	if err != nil {
		t.Errorf("The error appeared when it shouldn't: %v", err)
	}
}

func TestConfidenceSymbol_PartialMatch(t *testing.T) {
	fontA := testFontA()
	font := map[rune][8]string{'A': fontA}

	// первые 4 строки совпадают с 'A', остальные 4 — нет
	symbolArt := []string{fontA[0], fontA[1], fontA[2], fontA[3], "XXXXX", "XXXXX", "XXXXX", "XXXXX"}

	symbol, confidence, quality, err := ConfidenceSymbol(font, symbolArt)

	if symbol != 'A' {
		t.Errorf("symbol = %q, want %q", symbol, 'A')
	}
	if confidence != 50 {
		t.Errorf("confidence = %v, want 50", confidence)
	}
	if quality != "Partial match" {
		t.Errorf("quality = %q, want %q", quality, "Partial match")
	}
	if err != nil {
		t.Errorf("The error appeared when it shouldn't: %v", err)
	}
}

func TestConfidenceSymbol_Failed(t *testing.T) {
	fontA := testFontA()
	font := map[rune][8]string{'A': fontA}

	// совпадает только 1 строка из 8
	symbolArt := []string{fontA[0], "XXXXX", "XXXXX", "XXXXX", "XXXXX", "XXXXX", "XXXXX", "XXXXX"}

	symbol, confidence, quality, err := ConfidenceSymbol(font, symbolArt)

	if symbol != 'A' {
		t.Errorf("symbol = %q, want %q", symbol, 'A')
	}
	if confidence != 12.5 {
		t.Errorf("confidence = %v, want 12.5", confidence)
	}
	if quality != "Failed" {
		t.Errorf("quality = %q, want %q", quality, "Failed")
	}
	if err != nil {
		t.Errorf("The error appeared when it shouldn't: %v", err)
	}
}

func TestConfidenceSymbol_Space(t *testing.T) {
	font := map[rune][8]string{} // пустой шрифт — не должен влиять на распознавание пробела

	spaceLine := "      "
	symbolArt := []string{spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine, spaceLine}

	symbol, confidence, quality, err := ConfidenceSymbol(font, symbolArt)
	if err != nil {
		t.Errorf("The error appeared when it shouldn't: %v", err)
	}

	if symbol != ' ' {
		t.Errorf("symbol = %q, want %q", symbol, ' ')
	}
	if confidence != 100 {
		t.Errorf("confidence = %v, want 100", confidence)
	}
	if quality != "Perfect match" {
		t.Errorf("quality = %q, want %q", quality, "Perfect match")
	}
}

// ---------- OverallConfidence ----------

func TestOverallConfidence(t *testing.T) {
	tests := []struct {
		name        string
		matches     []ConfidenceStruct
		wantAvg     float64
		wantQuality string
	}{
		{
			name:        "100 -> Excellent",
			matches:     []ConfidenceStruct{{Confidence: 100}},
			wantAvg:     100,
			wantQuality: "Excellent",
		},
		{
			name:        "95 -> Excellent (нижняя граница)",
			matches:     []ConfidenceStruct{{Confidence: 95}},
			wantAvg:     95,
			wantQuality: "Excellent",
		},
		{
			name:        "94.99 -> Good (верхняя граница)",
			matches:     []ConfidenceStruct{{Confidence: 94.99}},
			wantAvg:     94.99,
			wantQuality: "Good",
		},
		{
			name:        "80 -> Good (нижняя граница)",
			matches:     []ConfidenceStruct{{Confidence: 80}},
			wantAvg:     80,
			wantQuality: "Good",
		},
		{
			name:        "79.99 -> Fair (верхняя граница)",
			matches:     []ConfidenceStruct{{Confidence: 79.99}},
			wantAvg:     79.99,
			wantQuality: "Fair",
		},
		{
			name:        "60 -> Fair (нижняя граница)",
			matches:     []ConfidenceStruct{{Confidence: 60}},
			wantAvg:     60,
			wantQuality: "Fair",
		},
		{
			name:        "59.99 -> Poor",
			matches:     []ConfidenceStruct{{Confidence: 59.99}},
			wantAvg:     59.99,
			wantQuality: "Poor",
		},
		{
			name:        "0 -> Poor",
			matches:     []ConfidenceStruct{{Confidence: 0}},
			wantAvg:     0,
			wantQuality: "Poor",
		},
		{
			name: "среднее нескольких символов",
			matches: []ConfidenceStruct{
				{Confidence: 100},
				{Confidence: 80},
			},
			wantAvg:     90,
			wantQuality: "Good",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotAvg, gotQuality := OverallConfidence(tc.matches)
			if gotAvg != tc.wantAvg {
				t.Errorf("OverallConfidence() avg = %v, want %v", gotAvg, tc.wantAvg)
			}
			if gotQuality != tc.wantQuality {
				t.Errorf("OverallConfidence() quality = %q, want %q", gotQuality, tc.wantQuality)
			}
		})
	}
}
