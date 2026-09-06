package confidence

import (
	"cmp"
	"fmt"
	"slices"
)

type MatchesStruct struct {
	Symbol rune
	Amount int
}

type ConfidenceStruct struct {
	Symbol     rune
	Confidence float64
	Quality    string
}

// Gap описывает один найденный "пробельный" промежуток
type Gap struct {
	Start  int // индекс колонки, с которой начинается пробел
	Length int // сколько колонок подряд пустые
}

// Данная функция сортирует массив по убыванию
func sortValues(a, b MatchesStruct) int {
	if a.Amount != b.Amount {
		return b.Amount - a.Amount
	} else {
		return cmp.Compare(a.Symbol, b.Symbol)
	}
}

// Данная функция проверяет, не является ли данный символ-арт артом пробела
func isEmptyArt(symbolArt []string) bool {
	if len(symbolArt) != 8 {
		return false
	}
	for _, line := range symbolArt { // Проверяем каждую строчку
		// fmt.Printf("%d-st line = %d", i+1, len(line))
		if line != "      " { // Если какая то строчка не содержит 6 пробелов подряд
			// fmt.Println("NO")
			return false // То, это не пробел
		}
		// fmt.Println("YES")
	}
	return true // Если все 8 строк состоят из 6 пробелов, то возвращаем true. Это арт пробела
}

// Данная функция делит арт на арт-строки
func DivideIntoLines(art []string) [][]string {
	var artLines [][]string
	for i := 0; i < len(art); i += 8 {
		artLines = append(artLines, art[i:i+8]) // Добавляем по 8 строчек кода в artLines. Это 1 арт-строка
	}
	return artLines

}

// Данная функция ищет границы между буквами и возвращает трехмерный массив, где арт поделен на буквы
func FindBorder(amountOfLines int, artLines [][]string) [][][]string {
	dividedArt := make([][][]string, amountOfLines)
	for i, artLine := range artLines {
		if len(artLine) != 8 {
			continue
		}

		width := len(artLine[0]) // Задаем ширину арта

		if width == 0 {
			continue
		}

		const spaceWidth = 6 // ширина ASCII-Art пробела в твоём шрифте

		isColEmpty := func(col int) bool { // Проверяем, является ли столбец пустым
			for _, line := range artLine {
				if col >= len(line) || line[col] != ' ' {
					return false
				}
			}
			return true
		}

		j := 0
		for j < width { // Перебираем столбцы арта
			if isColEmpty(j) {
				// нашли начало промежутка пустых колонок — считаем, сколько их подряд
				gapStart := j
				for j < width && isColEmpty(j) {
					j++
				}
				gapLength := j - gapStart

				if gapLength >= spaceWidth {
					// промежуток достаточно широкий — это настоящий пробел в слове
					spaceCount := gapLength / spaceWidth

					for n := 0; n < spaceCount; n++ {
						start := gapStart + n*spaceWidth
						end := start + spaceWidth

						var spaceBlock []string
						for _, line := range artLine {
							spaceBlock = append(spaceBlock, line[start:end])
						}
						dividedArt[i] = append(dividedArt[i], spaceBlock)
					}
				}
				// если промежуток короче spaceWidth — это просто технический зазор
				// между соседними буквами, его игнорируем и не добавляем как символ
				continue
			}

			// нашли непустую колонку — здесь начинается буква
			letterStart := j
			for j < width && !isColEmpty(j) {
				j++
			}

			letterEnd := j + 1
			if letterEnd > width {
				letterEnd = width
			}

			var letterBlock []string
			for _, line := range artLine {
				if letterEnd > len(line) {
					letterEnd = len(line)
				}
				letterBlock = append(letterBlock, line[letterStart:letterEnd])
			}
			dividedArt[i] = append(dividedArt[i], letterBlock)
		}
	}

	return dividedArt
}

// Функция, которая вычисляет уровень соответствие (confidence) каждого символа
func ConfidenceSymbol(artFontSymbols map[rune][8]string, symbolArt []string) (rune, float64, string, error) { // Определяет соответствие каждой буквы
	if artFontSymbols == nil {
		return ' ', 0, "Failed", fmt.Errorf("ошибка: шрифт ASCII-арта не определён")
	}

	if len(symbolArt) != 8 {
		return ' ', 0, "Failed", fmt.Errorf("ошибка: блок символа должен содержать 8 строк")
	}

	var matchingLines int
	var symbol rune
	if isEmptyArt(symbolArt) { // Если арт является пробелом
		matchingLines = 8
		symbol = ' ' // То записываем этот пробел в переменнную "symbol"
	} else {
		amountOfMatches := make(map[rune]int) // Мапа, где хранится сколько строк совпали с артом какого-то символа

		for j := ' '; j <= '~'; j++ {
			if _, ok := artFontSymbols[j]; ok {
				amountOfMatches[j] = 0
			}
		}

		for i := 0; i < 8; i++ { // Анализируем каждую строку арта и выявляем, с артом какой буквы оно совпало.
			for j := ' '; j <= '~'; j++ { // Арт отражает символ, чье количество совпадении равен 8
				template, ok := artFontSymbols[j]
				if !ok {
					continue
				}

				if template[i] == symbolArt[i] {
					amountOfMatches[j]++
				}
			}
		}

		if len(amountOfMatches) == 0 {
			return ' ', 0, "Failed", fmt.Errorf("ошибка: шрифт ASCII-арта не содержит поддерживаемых символов")
		}

		sortedMatches := make([]MatchesStruct, 0, len(amountOfMatches)) // Так как, количество совпадении храним в мапу, то результат он выводит рандомно
		for k, v := range amountOfMatches {                             // Поэтому результаты перекидываем в массив структур
			sortedMatches = append(sortedMatches, MatchesStruct{Symbol: k, Amount: v})
		}
		slices.SortFunc(sortedMatches, sortValues) // Сортируем массив по убыванию
		symbol = sortedMatches[0].Symbol
		matchingLines = sortedMatches[0].Amount // Нужный нам символ стоит первым. Нужно просто взять его количество совпадающих строк
	}
	confidence := float64(matchingLines) / 8 * 100

	// Выявляем словесное описание целостности символа
	var quality string
	if confidence == 100 {
		quality = "Perfect match"
	} else if confidence >= 50 {
		quality = "Partial match"
	} else {
		quality = "Failed"
	}
	return symbol, confidence, quality, nil
}

// Функция, которая вычисляет общую целостность арта путем нахождения среднего арифметического целостностей всех символов арта
func OverallConfidence(MatchesInfo []ConfidenceStruct) (float64, string) {
	if len(MatchesInfo) == 0 {
		return 0, "Poor"
	}

	sumOfConfidence := 0.0
	var averageConfidence float64
	for i := 0; i < len(MatchesInfo); i++ {
		sumOfConfidence += MatchesInfo[i].Confidence // Находим сумму целостностей
	}
	averageConfidence = sumOfConfidence / float64(len(MatchesInfo)) // Делим на количество

	// Выявляем словесное описание целостности арта
	var quality string
	if averageConfidence >= 95 {
		quality = "Excellent"
	} else if averageConfidence >= 80 {
		quality = "Good"
	} else if averageConfidence >= 60 {
		quality = "Fair"
	} else {
		quality = "Poor"
	}
	return averageConfidence, quality
}
