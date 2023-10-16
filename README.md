## Реализация целостности БД с помощью распределенной декомпозиции хэша

### Проблема, которую решает сервис
Если злоумышленник получил доступ к одной из двух нод и собственноручно изменит в ней данные, то сервис при проведении операции относительно определенного id проверит на валидность текущих полей по распределенному хэшу 2 нодах. Фактически сервис не позволит осуществить операцию при нарушенной целостности одного из полей к запрашиваемому id.

### Структура системы
2 ноды, в которых имеются изолированные БД.

Первая нода с БД с полями:
- id кошелька 
- баланс
- время последней операции
- hash (начальный хэш)

Вторая БД с полями:
- id кошелька
- hash (конечный хэш)

### Работа сервиса по обеспечению целостности
1.Сервис имеет четыре grpc эндпоинта:
- создание кошелька (создается только с нулевым балансом)
- удаление кошелька
- изменение баланса кошельков
- чтение баланса кошелька
2. При вызове эндпоинта на изменеие балансов, удаление кошелька, чтение баланса происходит такой алгоритм проверки целостности относительно определенного кошелька:
   1. В текущей ноде происходит чтение БД полей id, баланса, времени последней операции, к ним применяется расчет хэша
   2. Полученное значение считывается и сравнивается с конкатенируемым значением hex БД первой ноды и второй ноды по id, если не совпало, то реверт

### Полученный бэнчмарк для 10.000 трансферов (i7-10700 2.9GHz, 16GB RAM)
с логами в консоль 317 tps
без логов в консоль 334 tps

### Возможные проблемы
#### 1. Восстановление данных при нарушении целостности.
#### Возможное решение:
Ведение в отдельной БД лога каждого изменения с фиксацией конечного хэша, что позволит по id кошелька по внешнему хэшу (конечный хэш в отдельной ноде) 
восстановить последнюю запись с корректной валидацией (с некорректными данными в логах это будет сделать невозможно, что исключит риск подделки логов).

БД лога транзакций:
- id кошелька
- тип операции
- время операции
- предыдущий баланс
- принятный новый баланс
- hash (начльный)

Данная база может быть размещена на второй ноде. Соотвественно сервис дполнительно пишет в эту базу изменения при положительной валидации.

#### 2. Отказоустойчивость, при отключении одной из нод.
#### Возможное решение:
Использовать кластеризацию для кадого контура базы, т.е. создать два независимых кластера, которые обслуживают свою БД.

#### 3. При выполнении транзакции в две ноды в одну из них коммит не прошел, а в другую прошел, вследствии чего возможен рассинхрон.
#### Возможное решение:
В сервисе при не подтверждении коммита одной из операции записи в базу в рамках одной операции внедрен механизм отката предыдущего состояния транзакции, в которой коммит был осуществлен.

