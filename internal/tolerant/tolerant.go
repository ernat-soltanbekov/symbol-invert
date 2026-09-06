package tolerant

import (
	"symbol-invert/internal/confidence"
)

// RecognizeSymbol распознаёт один ASCII-art символ в tolerant-режиме.
// Возвращает найденный символ и true, если символ был восстановлен
// по ближайшему совпадению, а не по точному совпадению.
func RecognizeSymbol(artFontSymbols map[rune][8]string, symbolArt []string) (rune, bool, error) {
	symbol, confidenceValue, _, err := confidence.ConfidenceSymbol(artFontSymbols, symbolArt)
	if err != nil {
		return ' ', false, err
	}

	// Если confidence меньше 100%, значит полного совпадения
	// с шаблоном не было и символ восстановлен по ближайшему совпадению.
	recoveredByNearest := confidenceValue < 100

	return symbol, recoveredByNearest, nil
}

// RecognizeText распознаёт все ASCII-art символы
// и считает количество символов, восстановленных по ближайшему совпадению.
func RecognizeText(artFontSymbols map[rune][8]string, dividedArt [][][]string) (string, int, error) {
	result := ""
	nearestCount := 0

	// Перебираем каждый отдельный ASCII-art символ.
	for index, lineArt := range dividedArt {
		for _, symbolArt := range lineArt {

			symbol, recoveredByNearest, err := RecognizeSymbol(artFontSymbols, symbolArt)
			if err != nil {
				return "", 0, err
			}

			// Добавляем распознанный символ к итоговой строке.
			result += string(symbol)

			// Если символ был восстановлен не точным, а ближайшим совпадением, увеличиваем счётчик.
			if recoveredByNearest {
				nearestCount++
			}
		}
		if index != len(dividedArt)-1 {
			result += "\n"
		}
	}

	return result, nearestCount, nil
}
