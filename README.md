# symbol-invert

## О проекте

symbol-invert — программа на Go для преобразования текста в ASCII-art и обратного преобразования ASCII-art в обычный текст.

Программа поддерживает:

- вывод обычного текста в формате ASCII-art;
- восстановление текста из ASCII-art с помощью --reverse;
- анализ уверенности распознавания с помощью --confidence;
- tolerant-режим для частично повреждённого ASCII-art с помощью --tolerant.

При обратном преобразовании программа работает с шаблонами:

- standard.txt
- shadow.txt
- thinkertoy.txt

Для reverse-режима подходящий шаблон определяется программой автоматически.

## Требования

Для запуска проекта необходимы:

- Go;
- терминал или командная строка;
- файлы standard.txt, shadow.txt и thinkertoy.txt в корне проекта.

Проверить установленную версию Go можно командой:

    go version

## Установка

Склонируйте репозиторий:

    git clone <repository-url>

Перейдите в папку проекта:

    cd symbol-invert

Если проект уже скачан, достаточно открыть терминал в корневой папке проекта.

Дополнительные внешние библиотеки для работы проекта не требуются.

При необходимости можно проверить и синхронизировать go.mod командой:

    go mod tidy

## Запуск проекта

Основные варианты запуска:

Преобразование обычного текста в ASCII-art:

    go run . "hello"

Обратное преобразование ASCII-art в текст:

    go run . --reverse=testdata/example00.txt

Обратное преобразование с анализом confidence:

    go run . --reverse=testdata/example00.txt --confidence

Обратное преобразование в tolerant-режиме:

    go run . --reverse=testdata/noisy.txt --tolerant

## Сборка проекта

Проект можно собрать в исполняемый файл:

    go build .

После сборки в текущей директории будет создан исполняемый файл проекта.

На macOS и Linux его можно запустить так:

    ./symbol-invert

Например:

    ./symbol-invert --reverse=testdata/example00.txt

## Обычное преобразование текста в ASCII-art

Для преобразования обычного текста в ASCII-art:

    go run . "hello"

По умолчанию используется баннер:

    standard

Программа поддерживает ASCII-символы.

Также поддерживается перенос строки через \n.

Пример:

    go run . "Hello\nWorld"

## Обратное преобразование

Для восстановления обычного текста из ASCII-art используется флаг:

    --reverse=<fileName>

Пример:

    go run . --reverse=testdata/example00.txt

Результат:

    Hello World

Имя файла должно передаваться именно через знак =.

Правильно:

    go run . --reverse=testdata/example00.txt

Неправильно:

    go run . --reverse testdata/example00.txt

Если аргументы переданы в неправильном формате, программа выводит:

    Usage: go run . [OPTION]
    EX: go run . --reverse=<fileName>

## Формат входного файла

Файл для reverse-режима должен содержать ASCII-art.

Основные требования:

- количество строк должно быть кратно 8;
- каждая ASCII-art строка состоит из 8 физических строк;
- все 8 строк одного блока должны иметь одинаковую ширину;
- должны использоваться поддерживаемые ASCII-символы;
- файл не должен быть пустым.

Пример:

    go run . --reverse=file.txt

## Определение баннера

В проекте используются три banner-файла:

    standard.txt
    shadow.txt
    thinkertoy.txt

В reverse-режиме пользователю не нужно указывать название баннера вручную.

Программа загружает доступные шаблоны и пытается определить, какой из них соответствует входному ASCII-art.

## Как работает обратное преобразование

Каждая строка ASCII-art состоит из 8 физических строк.

Программа:

1. Читает ASCII-art из файла.
2. Проверяет корректность входных данных.
3. Проверяет, что количество строк кратно 8.
4. Проверяет одинаковую ширину строк внутри каждого 8-строчного блока.
5. Проверяет, что ASCII-art содержит только поддерживаемые ASCII-символы.
6. Загружает шаблоны banner-файлов.
7. Разделяет ASCII-art на отдельные символы.
8. Сравнивает полученные блоки с шаблонами.
9. Восстанавливает исходный текст.
10. Выводит результат.

Символы ASCII-art могут иметь разную ширину, поэтому программа не делит строку на блоки фиксированного размера.

Пробел между словами также обрабатывается как отдельный символ.

## Confidence

Флаг --confidence включает анализ уверенности распознавания.

Запуск:

    go run . --reverse=testdata/example00.txt --confidence

Сначала программа выводит восстановленный текст, затем подробный анализ.

Пример:

    Hello World

    --- Confidence Analysis ---
    Overall Confidence: 100.00%

    Character Details:
      Position 1: 'H' - 100.00% (Perfect match)
      Position 2: 'e' - 100.00% (Perfect match)
      Position 3: 'l' - 100.00% (Perfect match)

    Pattern Quality:
      Perfect Matches: 11
      Partial Matches: 0
      Failed Matches: 0

    Recognition Quality: Excellent

### Расчёт confidence

Каждый символ ASCII-art состоит из 8 строк.

Confidence рассчитывается по формуле:

    confidence = количество совпавших строк / 8 × 100

Например, если совпали 6 строк из 8:

    6 / 8 × 100 = 75%

Классификация символов:

- 100% — Perfect match;
- от 50% до менее 100% — Partial match;
- меньше 50% — Failed.

Например:

    8 / 8 = 100% → Perfect match
    6 / 8 = 75% → Partial match
    4 / 8 = 50% → Partial match
    3 / 8 = 37.5% → Failed

## Общая уверенность

Overall Confidence рассчитывается как среднее значение confidence всех распознанных символов.

Оценка качества:

- 95% и выше — Excellent;
- 80% и выше — Good;
- 60% и выше — Fair;
- меньше 60% — Poor.

Пример:

    Overall Confidence: 97.73%
    Recognition Quality: Excellent

## Tolerant mode

Флаг --tolerant используется для восстановления текста из частично повреждённого ASCII-art.

Запуск:

    go run . --reverse=testdata/noisy.txt --tolerant

Программа сравнивает повреждённые символы с доступными шаблонами и выбирает наиболее подходящий вариант.

Этот режим предназначен прежде всего для случаев, когда у символа повреждена одна или две строки.

Пример:

    Hello World
    Recovered by nearest match: 2

Recovered by nearest match показывает количество символов, которые были восстановлены не по полному совпадению с шаблоном.

## Проверка входных данных

Перед распознаванием программа выполняет несколько проверок.

### Пустой файл

Если файл пустой:

    Ошибка: файл пуст

### Количество строк

Количество физических строк ASCII-art должно быть кратно 8.

При неправильном количестве строк программа сообщает об ошибке:

    Ошибка: количество строк ASCII-art должно быть кратно 8

### Ширина строк

Все 8 строк одного ASCII-art ряда должны иметь одинаковую ширину.

При нарушении:

    Ошибка: строки ASCII-art имеют разную ширину.

### Неподдерживаемые символы

ASCII-art должен содержать поддерживаемые ASCII-символы.

При обнаружении неподдерживаемого символа:

    Ошибка: ASCII-art содержит неподдерживаемые символы.

## Использование

Формат reverse-команды:

    go run . --reverse=<fileName>

Дополнительные флаги:

    --confidence
    --tolerant

Примеры:

    go run . --reverse=file.txt
    go run . --reverse=file.txt --confidence
    go run . --reverse=file.txt --tolerant

Если формат аргументов неправильный, программа выводит:

    Usage: go run . [OPTION]
    EX: go run . --reverse=<fileName>

## Примеры reverse

    go run . --reverse=testdata/example00.txt

Результат:

    Hello World

    go run . --reverse=testdata/example01.txt

Результат:

    123

    go run . --reverse=testdata/example02.txt

Результат:

    #=\[

    go run . --reverse=testdata/example07.txt

Результат:

    ABCDEFGHIJKLMNOPQRSTUVWXYZ

## Структура проекта

    SYMBOL-INVERT/

    ├── internal/
    │   ├── banner/
    │   │   ├── banner.go
    │   │   └── banner_test.go
    │   │
    │   ├── confidence/
    │   │   ├── confidence.go
    │   │   └── confidence_test.go
    │   │
    │   ├── matcher/
    │   │   ├── matcher.go
    │   │   └── matcher_test.go
    │   │
    │   ├── printer/
    │   │   ├── printer.go
    │   │   └── printer_test.go
    │   │
    │   ├── reader/
    │   │   ├── reader.go
    │   │   └── reader_test.go
    │   │
    │   └── tolerant/
    │       ├── tolerant.go
    │       └── tolerant_test.go
    │
    ├── testdata/
    │   ├── empty.txt
    │   ├── example00.txt
    │   ├── example01.txt
    │   ├── example02.txt
    │   ├── example03.txt
    │   ├── example04.txt
    │   ├── example05.txt
    │   ├── example06.txt
    │   ├── example07.txt
    │   ├── noisy.txt
    │   └── reader_test.txt
    │
    ├── go.mod
    ├── main.go
    ├── README.md
    ├── shadow.txt
    ├── standard.txt
    └── thinkertoy.txt

## Тестирование

Для запуска всех тестов:

    go test ./...

Тесты проверяют:

- чтение файлов;
- загрузку banner-файлов;
- распознавание символов;
- обратное преобразование;
- разделение ASCII-art;
- точные совпадения;
- частичные совпадения;
- расчёт confidence;
- tolerant-распознавание;
- подсчёт символов, восстановленных по ближайшему совпадению.

## Форматирование кода

Для форматирования кода:

    go fmt ./...

После форматирования можно запустить тесты:

    go test ./...

## Используемые технологии

Проект написан на Go.

Используются стандартные возможности языка Go без внешних библиотек.
