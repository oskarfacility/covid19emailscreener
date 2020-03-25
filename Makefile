.PHONY : buildforall

buildforall :
	make clean
	GOOS=windows GOARCH=386 go build -o EmailCategorizer.exe .
	go build -o EmailCategorizer_macos .

clean:
	rm -f EmailCategorizer.exe
	rm -f EmailCategorizer_macos

