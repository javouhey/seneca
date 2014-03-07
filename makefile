all: format_main

format_main:
	gofmt -tabs=false -tabwidth=4 -w=true main.go

rebuild:
	./doclean; ./build

testgo:
	./dolist; ./dovet; ./dotest
