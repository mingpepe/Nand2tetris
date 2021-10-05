all:
	go build -o assembler.exe executable\assembler\main.go
run:
	assembler.exe -f projects\06\add\Add.asm
run_all:assembler.exe
	assembler.exe -f projects\06\add\Add.asm
	assembler.exe -f projects\06\max\Max.asm
	assembler.exe -f projects\06\max\MaxL.asm
	assembler.exe -f projects\06\pong\Pong.asm
	assembler.exe -f projects\06\pong\PongL.asm
	assembler.exe -f projects\06\rect\RectL.asm
	assembler.exe -f projects\06\rect\RectL.asm
assembler.exe: executable\assembler\main.go assembler\assembler.go
	go build -o assembler.exe executable\assembler\main.go