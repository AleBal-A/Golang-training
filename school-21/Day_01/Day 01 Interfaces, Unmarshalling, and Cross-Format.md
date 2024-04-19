# Day 01: Interfaces, Unmarshalling, and Cross-Format Data Handling in Go

### Отчет по Exercise 00: Reading

### Цель задачи:

Разработать код на Go для чтения баз данных в форматах XML и JSON через единый интерфейс. Приложение должно быть способно различать форматы файлов по их расширениям и использовать соответствующую реализацию интерфейса `DBReader` для чтения данных, преобразуя их в единый формат объектов и выводя в противоположный формат.

### Использованные данные:

- Файл **original_database.xml**, содержащий данные в формате XML.

### Результат выполнения программы:

Команда для запуска:

```
./Exercise_00 -f original_database.xml
```

Вывод программы:

```
2024/04/19 17:57:56 JSON data successfully written to: original_database.json
```

### Вывод:

Программа успешно выполнила задачу чтения и преобразования данных между двумя популярными форматами баз данных. Благодаря унифицированному интерфейсу `DBReader`, добавление поддержки новых форматов файлов или изменение логики обработки становится значительно проще, что улучшает масштабируемость и гибкость приложения. Таким образом, разработанный подход не только удовлетворяет текущим требованиям задачи, но и предоставляет основу для дальнейших улучшений и расширений функциональности.

### Отчет по Exercise_01 (сравнения баз данных файлов)

Разработать программу на языке Go для сравнения структур баз данных, хранящихся в форматах XML и JSON. Программа должна определить, какие элементы были добавлены, удалены или изменены между двумя версиями баз данных.

### Исходные данные:

- **original_database.xml**
- **stolen_database.json**

Каждая база данных содержит информацию о различных аспектах кондитерских изделий, включая состав ингредиентов, время приготовления и т.д.

### Результат выполнения программы:

Команда для запуска:

```
./Exercise_01 --old original_database.xml --new stolen_database.json
```

Вывод программы:

```
ADDED cake "Moonshine Muffin"
REMOVED cake "Blueberry Muffin Cake"
CHANGED cooking time for cake "Red Velvet Strawberry Cake" - "45 min" instead of "40 min"
CHANGED unit count for ingredient "Strawberries" for cake: "Red Velvet Strawberry Cake" - "8" instead of "7"
REMOVED unit "pieces" for ingredient "Cinnamon" for cake "Red Velvet Strawberry Cake"
CHANGED unit count for ingredient "Flour" for cake: "Red Velvet Strawberry Cake" - "2" instead of "3"
CHANGED unit for ingredient "Flour" for cake "Red Velvet Strawberry Cake" - "mugs" instead of "cups"
REMOVED ingredient "Vanilla extract" for cake "Red Velvet Strawberry Cake"
ADDED ingredient "Coffee Beans" for cake "Red Velvet Strawberry Cake"

```

### Отчет по Exercise_02 (сравнение больших дампов):

Разработать программу, которая сравнивает содержимое двух файлов с дампами файловых систем и выводит информацию о том, какие файлы были добавлены или удалены между двумя снимками.

### Описание решения:

1. Чтение и сохранение строк из файла `old.txt` в структуру `map`, где ключами служат строки.
2. Последовательное чтение файла `new.txt` и проверка строки на наличие в предварительно заполненной map. Если путь уже существует в map, он удаляется из неё (что означает, что файл не изменился). Если пути нет в map, это означает, что файл был добавлен, и программа выводит соответствующее сообщение.
3. После анализа всех путей в `new.txt`, все оставшиеся в map пути считаются удаленными, так как они были в `old.txt`, но не появились в `new.txt`. Программа выводит сообщения об удаленных файлах.

### Исходные данные:

- **old.txt:**
    
    ```
    /etc/stove/config.xml
    /etc/stove/No_need_file.txt
    
    ```
    
- **new.txt:**
    
    ```
    /etc/stove/config.xml
    /Users/baker/recipes/database.xml
    /Users/baker/New_New_new.txt
    
    ```
    

### Результат выполнения программы:

Команда для запуска:

```
go build 
./Exercise_02 --old old.txt --new new.txt
```

Вывод программы:

```
ADDED /Users/baker/recipes/database.xml
ADDED /Users/baker/New_New_new.txt
REMOVED /etc/stove/No_need_file.txt
```