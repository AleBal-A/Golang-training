all: test

test:
	# Запуск тестовых запросов от клиента на сервер
	go run cclient.go -k AA -c 200 -m 50
	go run cclient.go -k AA -c 2 -m 50 | jq -r '.thanks'
	go run cclient.go -k AA -c 2 -m 60
	go run cclient.go -k DE -c 2 -m 42 | jq -r '.thanks'
