package printer

import (
	"fmt"
	"io"
	"symbol-invert/internal/banner"
)

func PrintLines(w io.Writer, normArgs []string, font string) { // Выводит готовый ASCII-арт в любой поток вывода, сохраняя структуру строк исходного текста.
	for _, arg := range normArgs { // Последовательно обрабатываем каждую строку из подготовленного списка.
		if arg == "" { // Если встретилась пустая строка, выводим пустую строку и сразу переходим к следующей.
			fmt.Fprintln(w)
			continue
		}

		result, err := banner.Render(font, arg) // Генерируем ASCII-арт для текущей строки.
		if err != nil {                         // Если при генерации произошла ошибка, выводим сообщение в указанный поток и прекращаем работу функции.
			fmt.Fprintln(w, "Error:", err)
			return
		}

		for _, line := range result { // Построчно отправляем готовый ASCII-арт в переданный поток вывода.
			fmt.Fprintln(w, line)
		}
	}
}
