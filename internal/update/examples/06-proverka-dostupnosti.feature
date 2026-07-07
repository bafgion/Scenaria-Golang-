@smoke
Функционал: Примеры для новичков
Сценарий: Кнопка видна, но недоступна для нажатия
	Допустим открыт "data:text/html,<html><body><button id=submit disabled>Отправить</button></body></html>"
	Тогда вижу "#submit"
	И проверяю что недоступно "#submit"
	И закрываю браузер

Сценарий: Кнопка доступна для нажатия
	Допустим открыт "data:text/html,<html><body><button id=submit>Отправить</button></body></html>"
	Тогда проверяю что доступно "#submit"
	И нажимаю "#submit"
	И закрываю браузер
