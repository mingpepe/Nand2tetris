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
run_vm_7: vm.exe
	vm.exe -f projects\07\MemoryAccess\BasicTest\BasicTest.vm
	vm.exe -f projects\07\MemoryAccess\PointerTest\PointerTest.vm
	vm.exe -f projects\07\MemoryAccess\StaticTest\StaticTest.vm
	vm.exe -f projects\07\StackArithmetic\SimpleAdd\SimpleAdd.vm
	vm.exe -f projects\07\StackArithmetic\StackTest\StackTest.vm
run_vm_8: vm.exe
	vm.exe -d projects\08\FunctionCalls\FibonacciElement
	vm.exe -bypass=true -d projects\08\FunctionCalls\NestedCall
	vm.exe -bypass=true -f projects\08\FunctionCalls\SimpleFunction\SimpleFunction.vm
	vm.exe -d projects\08\FunctionCalls\StaticsTest
	vm.exe -bypass=true -f projects\08\ProgramFlow\BasicLoop\BasicLoop.vm
	vm.exe -bypass=true -f projects\08\ProgramFlow\FibonacciSeries\FibonacciSeries.vm
run_vm_8_v: vm.exe
	vm.exe -v=true -d projects\08\FunctionCalls\FibonacciElement
	vm.exe -v=true -bypass=true -d projects\08\FunctionCalls\NestedCall
	vm.exe -v=true -bypass=true -f projects\08\FunctionCalls\SimpleFunction\SimpleFunction.vm
	vm.exe -v=true -d projects\08\FunctionCalls\StaticsTest
	vm.exe -v=true -bypass=true -f projects\08\ProgramFlow\BasicLoop\BasicLoop.vm
	vm.exe -v=true -bypass=true -f projects\08\ProgramFlow\FibonacciSeries\FibonacciSeries.vm
assembler.exe: executable\assembler\main.go assembler\assembler.go
	go build -o assembler.exe executable\assembler\main.go
vm.exe: executable\vm\main.go vm\vm.go
	go build -o vm.exe executable\vm\main.go
myapp: MyApp\DirectRAM\Main.jack MyApp\Helloworld\Main.jack MyApp\Error\Main.jack MyApp\Shell\Main.jack
	tools\JackCompiler.bat MyApp\DirectRAM
	tools\JackCompiler.bat MyApp\Helloworld
	tools\JackCompiler.bat MyApp\Error
	tools\JackCompiler.bat MyApp\Shell
os:
	tools\JackCompiler.bat projects\12\MemoryTest
	tools\JackCompiler.bat projects\12\KeyboardTest
	tools\JackCompiler.bat projects\12\StringTest
	tools\JackCompiler.bat projects\12\ArrayTest
	tools\JackCompiler.bat projects\12\SysTest
	tools\JackCompiler.bat projects\12\MathTest