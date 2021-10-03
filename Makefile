all:
	go build -o assembler.exe executable\assembler\main.go
run:
	assembler.exe -f projects\06\add\Add.asm