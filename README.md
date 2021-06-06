# backend1
practice work backend #1 golang


Сервер загружает файлы и выдает список из каталога ./files
Вывод осуществляется в формате JSON.
1) Для получения списка файлов используется GET-запос с маршрутом /list (функция listHandlerFunc) 
http://localhost:8080/list
2) Для задания фильтра по расширению файла query-параметр с ключом ext
http://localhost:8080/list?ext=.txt
