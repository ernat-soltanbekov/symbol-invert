package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"symbol-invert/internal/banner"
	"symbol-invert/internal/confidence"
	"symbol-invert/internal/matcher"
	"symbol-invert/internal/printer"
	"symbol-invert/internal/tolerant"
)

func main() {
	// Загружаем все баннеры.
	standardSymbols, err := banner.Banner("standard.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении файла standard: ", err)
		return
	}

	shadowSymbols, err := banner.Banner("shadow.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении файла shadow: ", err)
		return
	}

	thinkertoySymbols, err := banner.Banner("thinkertoy.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении файла thinkertoy: ", err)
		return
	}

	// Проверяем, что пользователь передал хотя бы один аргумент.
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Первый аргумент должен иметь формат --reverse=<fileName>.
	if !strings.HasPrefix(os.Args[1], "--") && len(os.Args) == 2 {
		textArg := os.Args[1] // Здесь будет храниться текст, который необходимо преобразовать в ASCII-арт.

		for _, symbol := range textArg {
			if symbol > 127 {
				fmt.Fprintln(os.Stderr, "Ошибка: программа поддерживает только ASCII символы.")
				os.Exit(1)
			}
		}

		onlyLines := true // Предполагаем, что вход состоит только из пустых строк, пока не обнаружим текст.
		textArg = strings.ReplaceAll(textArg, `\n`, "\n")
		normArgs := strings.Split(textArg, "\n") // Разбиваем строку по последовательности "\n", чтобы обработать каждую строку отдельно.

		for _, arg := range normArgs { // Просматриваем каждую полученную строку.
			if arg != "" {
				onlyLines = false
			}
		}

		if onlyLines { // Если пользователь передал только пустые строки, выводим соответствующее количество пустых строк и завершаем работу.
			for i := 0; i < len(normArgs)-1; i++ {
				fmt.Println()
			}
			return
		}

		printer.PrintLines(os.Stdout, normArgs, "standard.txt") // Передаем подготовленные строки модулю, который выводит ASCII-арт на экран.
		return
	} else if !strings.HasPrefix(os.Args[1], "--reverse=") {
		printUsage()
		return
	}

	// Получаем имя файла из аргумента.
	filename := strings.TrimPrefix(os.Args[1], "--reverse=")

	if filename == "" {
		printUsage()
		return
	}

	confidenceFlag := false
	tolerantFlag := false

	// Проверяем дополнительные флаги и название баннера.
	for _, arg := range os.Args[2:] {
		switch arg {
		case "--confidence":
			confidenceFlag = true
		case "--tolerant":
			tolerantFlag = true
		default:
			printUsage()
			return
		}
	}

	// Имя файла берём из --reverse=<fileName>.
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при открытии файла: ", err)
		return
	}

	defer file.Close()

	var art []string // Переменная, куда будем хранить арт из текстового файла

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		art = append(art, line)
	}

	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении файла: ", err)
		return
	}

	if len(art) == 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: файл пуст")
		return
	}

	// Проверяем, что количество строк ASCII-art кратно 8.
	if len(art)%8 != 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: количество строк ASCII-art должно быть кратно 8")
		return
	}

	// Проверяем корректность каждой 8-строчной строки ASCII-art.
	for start := 0; start < len(art); start += 8 {
		expectedWidth := len(art[start])

		for row := start; row < start+8; row++ {
			// ASCII-art должен содержать только ASCII-символы.
			for _, symbol := range art[row] {
				if symbol > 127 {
					fmt.Fprintln(os.Stderr, "Ошибка: ASCII-art содержит неподдерживаемые символы.")
					return
				}
			}

			// Все 8 строк одного ASCII-art ряда должны иметь одинаковую ширину.
			if len(art[row]) != expectedWidth {
				fmt.Fprintln(os.Stderr, "Ошибка: строки ASCII-art имеют разную ширину.")
				return
			}
		}
	}

	if !confidenceFlag && !tolerantFlag {
		bannerFiles := []string{
			"standard.txt",
			"shadow.txt",
			"thinkertoy.txt",
		}

		for _, bannerFile := range bannerFiles {
			templates, err := matcher.LoadTemplates(bannerFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при загрузке шаблонов:", err)
				return
			}

			recoveredText := ""
			matched := true

			for i := 0; i < len(art); i += 8 {
				var row [8]string
				copy(row[:], art[i:i+8])

				text, err := matcher.MatchRow(row, templates)
				if err != nil {
					matched = false
					break
				}

				if i != 0 {
					recoveredText += "\n"
				}

				recoveredText += text
			}

			if matched {
				fmt.Println(recoveredText)
				return
			}
		}

		fmt.Fprintln(os.Stderr, "Ошибка: не удалось распознать ASCII-art")
		return
	}

	artLines := confidence.DivideIntoLines(art)

	fontCandidates := []map[rune][8]string{
		standardSymbols,
		shadowSymbols,
		thinkertoySymbols,
	}

	var artFontSymbols map[rune][8]string
	var dividedArt [][][]string

	bestScore := -1.0

	for _, symbols := range fontCandidates {
		candidateArt := confidence.FindBorder(
			len(artLines),
			artLines,
		)

		totalConfidence := 0.0
		symbolCount := 0
		valid := true

		for _, lineArt := range candidateArt {
			for _, symbolArt := range lineArt {
				_, symbolConfidence, _, err := confidence.ConfidenceSymbol(
					symbols,
					symbolArt,
				)

				if err != nil {
					valid = false
					break
				}

				totalConfidence += symbolConfidence
				symbolCount++
			}

			if !valid {
				break
			}
		}

		if !valid || symbolCount == 0 {
			continue
		}

		averageConfidence := totalConfidence / float64(symbolCount)

		if averageConfidence > bestScore {
			bestScore = averageConfidence
			artFontSymbols = symbols
			dividedArt = candidateArt
		}
	}

	if artFontSymbols == nil {
		fmt.Fprintln(os.Stderr, "Ошибка: неизвестный или неподдерживаемый шрифт ASCII-арта.")
		return
	}

	// Добавлен tolerant-режим.
	if tolerantFlag {
		text, nearestCount, err := tolerant.RecognizeText(artFontSymbols, dividedArt)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(text)
		fmt.Printf("Recovered by nearest match: %d\n", nearestCount)

		// Если --confidence не передан, завершаем работу.
		// Если переданы оба флага, продолжаем confidence-анализ.
		if !confidenceFlag {
			return
		}
	}

	var MatchesInfo []confidence.ConfidenceStruct
	perfectMatches := 0
	partialMatches := 0
	fails := 0
	recoveredText := ""

	for i := 0; i < len(dividedArt); i++ {
		for j := 0; j < len(dividedArt[i]); j++ {
			matchedSymbol, symbolConfidence, singleQuality, err := confidence.ConfidenceSymbol(artFontSymbols, dividedArt[i][j])
			if err != nil {
				fmt.Println(err)
				return
			}

			MatchesInfo = append(MatchesInfo, confidence.ConfidenceStruct{
				Symbol:     matchedSymbol,
				Confidence: symbolConfidence,
				Quality:    singleQuality,
			})

			// Добавляем распознанный символ в итоговую строку.
			recoveredText += string(matchedSymbol)

			// Проверяем качество совпадения.
			switch singleQuality {
			case "Perfect match":
				perfectMatches++
			case "Partial match":
				partialMatches++
			case "Failed":
				fails++
			}
		}
		if i != len(dividedArt)-1 {
			recoveredText += "\n"
		}
	}

	// Сначала выводим восстановленный текст.
	// В tolerant-режиме текст уже был выведен выше.
	if !tolerantFlag {
		fmt.Println(recoveredText)
	}

	// Если --confidence не передан, выводим только текст.
	if !confidenceFlag {
		return
	}

	overallConfidence, overallQuality := confidence.OverallConfidence(MatchesInfo)

	// Вывод
	fmt.Println("\n--- Confidence Analysis ---")
	fmt.Printf("Overall Confidence: %.2f%%\n\n", overallConfidence)

	fmt.Println("Character Details:")
	for i := 0; i < len(MatchesInfo); i++ {
		fmt.Printf("  Position %d: '%c' - %g%% (%s)\n", i+1, MatchesInfo[i].Symbol, MatchesInfo[i].Confidence, MatchesInfo[i].Quality)
	}

	fmt.Println("\nPattern Quality:")
	fmt.Printf("  Perfect Matches: %d\n", perfectMatches)
	fmt.Printf("  Partial Matches: %d\n", partialMatches)
	fmt.Printf("  Failed Matches: %d\n\n", fails)
	fmt.Printf("Recognition Quality: %s\n", overallQuality)
}

// Функция выводит сообщение Usage из ТЗ.
func printUsage() {
	fmt.Println("Usage: go run . [OPTION]")
	fmt.Println("EX: go run . --reverse=<fileName>")
}
