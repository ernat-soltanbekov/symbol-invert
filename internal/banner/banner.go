package banner

import (
	"bufio"
	"os"
)

func Banner(filename string) (map[rune][8]string, error) { // Загружает файл шрифта в память и возвращает его содержимое в виде среза строк.
	file, mistake := os.Open(filename) // Открываем файл баннера, который будем читать.
	if mistake != nil {                // Если открыть файл не удалось, сразу возвращаем ошибку вызывающему коду.
		return nil, mistake
	}

	defer file.Close() // После завершения функции обязательно закрываем файл, чтобы освободить системные ресурсы.

	artOfSymbols := make(map[rune][8]string) // Создаем мапу, куда будем записать ASCII-Art каждого символа
	scanner := bufio.NewScanner(file)        // Создаем сканер, который умеет читать файл построчно.

	c := 32              // Показывает ASCII-Art какого символа мы будем записывать в мапу
	i := 0               // Показывает какой номер строки ASCII-Artа мы будем записывать в массив
	var test [8]string   // Объявляем слайс, который будет хранить ASCII-art символа
	scanner.Scan()       // Скипаем первую строку
	for scanner.Scan() { // Последовательно считываем каждую строку файла.
		if i < 8 { // Если мы не завершили записывать ASCII-Art
			test[i] = scanner.Text()
			i++
		} else { // Если завершили, то записываем арт в мапу
			artOfSymbols[rune(c)] = test
			i = 0
			c++
		}
	}

	if i == 8 {
		artOfSymbols[rune(c)] = test
	}

	if err := scanner.Err(); err != nil { // После чтения проверяем, не возникло ли ошибок во время работы сканера.
		return nil, err
	}

	return artOfSymbols, nil // Возвращаем загруженные шаблоны символов.
}

func getSymbolLines(symbolArtList map[rune][8]string, char rune) [8]string { // Находит восемь строк ASCII-арта, соответствующих одному символу.
	return symbolArtList[char] // Возвращаем восемь строк, образующих изображение символа.
}

func renderWord(symbolArtList map[rune][8]string, word string) []string { // Собирает ASCII-арт для целого слова, объединяя изображения всех его символов.
	blocks := make([][8]string, len(word)) // Подготавливаем место для хранения ASCII-блоков каждого символа слова.
	for i, char := range word {            // Последовательно получаем ASCII-представление каждого символа.
		blocks[i] = getSymbolLines(symbolArtList, char)
	}

	result := make([]string, 8) // Будущий результат всегда состоит из восьми строк.

	for stroka := 0; stroka < 8; stroka++ { // Формируем результат построчно, проходя сверху вниз.
		for _, block := range blocks { // На каждой строке последовательно объединяем части всех символов.
			result[stroka] += block[stroka]
		}
	}

	return result // Возвращаем полностью собранный ASCII-арт слова.
}

func Render(filename string, word string) ([]string, error) { // Главная функция модуля: загружает шрифт и строит ASCII-арт для указанного слова.
	symbolArtList, mistake := Banner(filename) // Сначала читаем файл шрифта в память.
	if mistake != nil {                        // Если файл загрузить не удалось, сразу возвращаем ошибку.
		return nil, mistake
	}
	return renderWord(symbolArtList, word), nil // Генерируем ASCII-арт и возвращаем его вызывающему коду.
}
