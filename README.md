# Превьювер изображений

[![workflow](https://github.com/dmitriygoldberg/image-previewer/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/dmitriygoldberg/image-previewer/actions)

Сервис предназначен для изготовления preview (создания изображения с новыми размерами на основе имеющегося изображения).

Сервис представляет собой web-сервер (прокси), загружающий изображения, масштабирующий/обрезающий их до нужного формата и возвращающий пользователю.

### Использование привьювера:
http://127.0.0.1/fill/{width}/{height}/{url}

Где: width - ширина, height - высота, url - ссылка на изображение

### Пример:
```http://127.0.0.1/fill/200/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_1024x252.jpg```


## Make команды

- ``make run`` - билд сервиса + поднятие nginx контейнера для работы с сервисом (по умолчанию 127.0.0.1 порт 80)
- ``make stop`` - остановка контейнеров
- ``make build`` - билд сервиса с локального окружения
- ``make clear`` - очистка билдов
- ``make lint`` - запуск линтера кода
- ``make test`` - запуск юнит тестов
