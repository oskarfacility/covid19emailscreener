.PHONY : buildforall

buildforall :
	rm BCG-EMAIL-CATEGO.exe
	GOOS=windows GOARCH=386 go build -o BCG-EMAIL-CATEGO.exe .
	rm BCG-EMAIL-CATEGO_macos
	go build -o BCG-EMAIL-CATEGO_macos .
