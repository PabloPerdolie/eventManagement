````mermaid
sequenceDiagram
    participant U as Пользователь
    participant C as Клиент (React)
    participant AG as API Gateway
    participant CS as Core Service (Go)
    participant DB as База данных (PostgreSQL)
    participant NS as Notification Service (Go)
    participant CMS as Communication Service (Go)
    participant ES as Внешняя служба

    U ->> C: Вводит данные мероприятия
    C ->> AG: POST /events (данные мероприятия, токен)
    AG ->> AG: Проверяет токен
    AG -->> C: Ошибка (если токен недействителен)
    AG ->> CS: Перенаправляет запрос (если токен валиден)
    CS ->> DB: Сохраняет мероприятие
    DB -->> CS: Подтверждение сохранения
    CS -->> AG: Успешный ответ
    AG -->> C: Успешный ответ
    C -->> U: Отображает мероприятие

    U ->> C: Назначает задачу участнику
    C ->> AG: POST /tasks (задача, ID участника)
    AG ->> CS: Перенаправляет запрос
    CS ->> DB: Сохраняет задачу
    DB -->> CS: Подтверждение
    CS ->> NS: Данные для уведомления о задаче
    NS ->> ES: Отправляет уведомление о задаче
    ES -->> NS: Подтверждение отправки
    NS -->> CS: Уведомление отправлено
    CS -->> AG: Успешный ответ
    AG -->> C: Успешный ответ
````