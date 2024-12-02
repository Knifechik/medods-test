# medods-test

Часть сервиса аутентификации.\
Единовременно на устройствах может быть запущена одна сессия.\
Не все возможные случаи ошибок для респонса, тесты только для апи слоя.\
Для рефреш операции приходят оба токена для сравнения их общего ID.\
Токены хранятся в куки.\
История сессий не сохраняется.\
Логируются только ошибки.

Для запуска docker compose up --build

Endpoint login localhost:8080/login/(UUID)\
Возвращает два куки с access и refresh токенам.\
Endpoint refresh localhost:8080/refresh\
Принимает два куки access и refresh, возвращает два новых куки.
