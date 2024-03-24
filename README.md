# "Фильмотека"

## Запуск 

Для запуска можно использовать run.sh. Файл необходимо запустить из директории film_library. При использование этого скрипта будет подняты 3 docker контейнера. Контейнер "film" будет содержать в себе непосредственно приложение. Два других контейнера содержат в себе базу данных PostgreSQL и adminer к базе данных.

## Запросы к API

Для того, чтобы делать запросы к API. Необходимо в таблицу user_rights добавить токен и права пользователя для этого токена. Существует два типа прав “admin” и “user”. Пользователь с правами “user” может делать только GET запросы. Таблица содержит в себе следующие колонки user_rights (id, user_right, token). Где в token нужно положить токен пользователя, а в user_right – роль этого пользователя. 

### POST запросы 

"/actor/create" – создает актеров, на вход нужно подать JSON формата: 
{"name_actor": "", "gender": "", "birthday": ""} 

"/film/create" – создает фильм, на вход нужно подать JSON формата: 
{"name_film": "", "description": "", "rating": "", "release_date": "" }

### PUT запросы 

"/actor/change" – изменить существующего актера, на вход нужно подать JSON формата: 
{"name_actor": "", "gender": "", "birthday": "", "replace_name_actor": "", "replace_gender": "", "replace_birthday": ""} 

"/film/add/actors" – добавить существующих актеров в БД в фильм, на вход нужно подать JSON формата:
{" name_film ": "", "release_date": "", actors: [{"name_actor": "", "gender": "", "birthday": ""} 
,] }

"/film/change" - изменить существующий фильм, на вход нужно подать JSON формата:
{"name_film": "", "description": "", "rating": "", "release_date": "", "replace_description": "", "replace_rating": "", "replace_release_date": ""}

<b>Внимание, изменить название фильма нельзя!!!</b>

### DELETE запросы 

"/actor/delete" – удалить существующего актера, на вход нужно подать JSON формата: 
{"name_actor": "", "gender": "", "birthday": "", "replace_name_actor": ""}

"/film/delete" – удалить существующий фильм, на вход нужно подать JSON формата:{" name_film ": "", "release_date": ""}

"/film/remove/actors" – удалить актеров из фильма, на вход нужно подать JSON формата:
{" name_film ": "", "release_date": "", actors: [{"name_actor": "", "gender": "", "birthday": ""} ,] }

### GET запросы

"/films" – показывает все фильмы в БД, опциональные аргументы "sort_by_coloms" отвечает за сортировку по колонке (name, released, rating) и "direction" направление сортировки (ASC и DESC). Возвращаем JSON формата {"name_films": ["",]}

"/films/name/film" - показывает все фильмы в БД поиск фильма происходит фрагменту названия фильма. Атрибут “name_film”. Возвращаем JSON формата {"name_films": ["",]}

"/films/name/actor" - показывает все фильмы в БД поиск фильма происходит фрагменту имени актера. Атрибут “name_actor”. Возвращаем JSON формата {"name_films": ["",]}

"/actors" – показываем всех актеров и фильмы, в которых они снимались. Возвращаем JSON, где ключ это имя актера, а значение список фильмов. 
