# medods-test

Часть сервиса аутентификации.
Единовременно на устройствах может быть запущена одна сессия.
Не все возможные случаи ошибок для респонса, тесты только для апи слоя.
В обычных случаях для рефреш операции должен приходить только рефреш токен, но тогда не понимал зачем взаимосвязь токенов, поэтому приходят оба для сравнения.
Токены хранятся в куки.
История сессий не сохраняется.
Логируются только ошибки.

