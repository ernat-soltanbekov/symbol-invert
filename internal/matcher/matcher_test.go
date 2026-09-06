package matcher // Тест принадлежит тому же пакету matcher.

import "testing" // Стандартный пакет для тестов.

// fakeTemplates строит две простые, придуманные "буквы" разной ширины —
// специально, чтобы проверить именно логику поиска границ, не завися от
// настоящего файла баннера.
func fakeTemplates() map[rune]Template { // Вспомогательная функция, не тест сама по себе.
	return map[rune]Template{ // Возвращаем готовую карту из двух шаблонов.
		'A': { // Первая "буква" — шириной 3 колонки.
			Char:  'A',
			Width: 3,
			Lines: [8]string{"AAA", "A.A", "A.A", "AAA", "A.A", "A.A", "A.A", "..."},
		},
		'I': { // Вторая "буква" — шириной всего 1 колонка, специально уже.
			Char:  'I',
			Width: 1,
			Lines: [8]string{"I", "I", "I", "I", "I", "I", "I", "I"},
		},
	}
}

// buildRow склеивает построчно шаблоны нескольких символов в одну "восьмёрку",
// как будто это уже готовый ASCII-арт целого слова.
func buildRow(templates map[rune]Template, word string) [8]string { // Принимает шаблоны и слово, которое нужно "нарисовать".
	var row [8]string        // Заготовка под итоговые 8 строк.
	for i := 0; i < 8; i++ { // Идём по каждой из 8 строк отдельно.
		for _, c := range word { // И по каждому символу слова.
			row[i] += templates[c].Lines[i] // Дописываем соответствующую строку его шаблона.
		}
	}
	return row // Возвращаем готовую "восьмёрку".
}

func TestMatchRowSimpleWord(t *testing.T) { // Проверяем, что MatchRow верно восстанавливает слово из "AI".
	templates := fakeTemplates()     // Берём наши придуманные шаблоны.
	row := buildRow(templates, "AI") // Строим блок, как будто это готовый рисунок слова "AI".

	got, err := MatchRow(row, templates) // Вызываем проверяемую функцию.
	if err != nil {                      // Ошибки быть не должно — блок собран из тех же шаблонов.
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if got != "AI" { // Сверяем результат с исходным словом.
		t.Errorf("ожидалось %q, получено %q", "AI", got)
	}
}

func TestMatchRowDifferentWidths(t *testing.T) { // Отдельно проверяем именно разную ширину символов — суть задания.
	templates := fakeTemplates()
	row := buildRow(templates, "IIA") // Узкий, узкий, широкий символ подряд — проверка на смещение границ.

	got, err := MatchRow(row, templates)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got != "IIA" {
		t.Errorf("ожидалось %q, получено %q", "IIA", got)
	}
}

func TestMatchRowUnrecognizedBlock(t *testing.T) { // Проверяем поведение, когда символ вообще не распознаётся.
	templates := fakeTemplates()

	badRow := [8]string{ // Специально собираем "мусорный" блок, не похожий ни на один шаблон.
		"???", "???", "???", "???", "???", "???", "???", "???",
	}

	_, err := MatchRow(badRow, templates) // Вызываем функцию с заведомо нераспознаваемым блоком.
	if err == nil {                       // Мы ОЖИДАЕМ ошибку — если её нет, это баг.
		t.Fatal("ожидалась ошибка при нераспознанном блоке, но её не было")
	}
}
