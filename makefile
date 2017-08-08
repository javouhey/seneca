all: format_main

format_main:
	go fmt main.go

rebuild:
	./doclean; ./build

testgo:
	./dolist; ./dovet; ./dotest
