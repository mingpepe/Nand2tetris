all:assembler.exe vm.exe
	
run:
	assembler.exe -f projects\06\add\Add.asm
run_asm:assembler.exe
	assembler.exe -f projects\06\add\Add.asm
	assembler.exe -f projects\06\max\Max.asm
	assembler.exe -f projects\06\max\MaxL.asm
	assembler.exe -f projects\06\pong\Pong.asm
	assembler.exe -f projects\06\pong\PongL.asm
	assembler.exe -f projects\06\rect\RectL.asm
	assembler.exe -f projects\06\rect\RectL.asm
run_vm: vm.exe
	vm.exe -f projects\07\MemoryAccess\BasicTest\BasicTest.vm
	vm.exe -f projects\07\MemoryAccess\PointerTest\PointerTest.vm
	vm.exe -f projects\07\MemoryAccess\StaticTest\StaticTest.vm
	vm.exe -f projects\07\StackArithmetic\SimpleAdd\SimpleAdd.vm
	vm.exe -f projects\07\StackArithmetic\StackTest\StackTest.vm
assembler.exe: executable\assembler\main.go assembler\assembler.go
	go build -o assembler.exe executable\assembler\main.go
vm.exe: executable\vm\main.go vm\vm.go
	go build -o vm.exe executable\vm\main.go