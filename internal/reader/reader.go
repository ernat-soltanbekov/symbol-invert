package reader

import (
	"bufio"
	"os"
)

// ReadLines построчно читает файл и возвращает все строки в виде среза []string.
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename) // Открываем файл по переданному имени.
	if err != nil {                // Если файл открыть не удалось, возвращаем ошибку.
		return nil, err
	}

	defer file.Close() // Гарантируем закрытие файла после завершения функции.

	var lines []string // Создаём срез, в который будем добавлять строки из файла.

	scanner := bufio.NewScanner(file) // Создаём сканер для последовательного чтения файла построчно.
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() { // Пока сканер успешно читает следующую строку -
		lines = append(lines, scanner.Text()) // - добавляем прочитанную строку в срез.
	}

	if err := scanner.Err(); err != nil { // После чтения проверяем, не возникла ли ошибка сканирования файла.
		return nil, err // Если ошибка возникла, возвращаем её вызывающему коду.
	}

	return lines, nil // Возвращаем все прочитанные строки и nil, если ошибок не было.
}
