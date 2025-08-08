all:assembler.exe vm.exe analyzer.exe
	
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
	tools\JackCompiler.bat projects\12\OutputTest
	tools\JackCompiler.bat projects\12\ScreenTest
os_test_app:
	python copy_os_for_test.py
	tools\JackCompiler.bat projects\11\Pong
tokenizer_test.exe: executable\tokenizer_test\main.go analyzer\tokenizer.go
	go build -o tokenizer_test.exe executable\tokenizer_test\main.go
test_tokenizer: tokenizer_test.exe
	tokenizer_test.exe -f projects\10\ArrayTest\Main.jack
	tools\TextComparer.bat projects\10\ArrayTest\Main_KMT.xml projects\10\ArrayTest\MainT.xml

	tokenizer_test.exe -d projects\10\ExpressionLessSquare
	tools\TextComparer.bat projects\10\ExpressionLessSquare\Main_KMT.xml projects\10\ExpressionLessSquare\MainT.xml
	tools\TextComparer.bat projects\10\ExpressionLessSquare\Square_KMT.xml projects\10\ExpressionLessSquare\SquareT.xml
	tools\TextComparer.bat projects\10\ExpressionLessSquare\SquareGame_KMT.xml projects\10\ExpressionLessSquare\SquareGameT.xml

	tokenizer_test.exe -d projects\10\Square
	tools\TextComparer.bat projects\10\Square\Main_KMT.xml projects\10\Square\MainT.xml
	tools\TextComparer.bat projects\10\Square\Square_KMT.xml projects\10\Square\SquareT.xml
	tools\TextComparer.bat projects\10\Square\SquareGame_KMT.xml projects\10\Square\SquareGameT.xml

compilation_engine_test.exe: executable\compilation_engine_test\main.go compiler\tokenizer.go compiler\compilation_engine_xml.go
	go build -o compilation_engine_test.exe executable\compilation_engine_test\main.go
test_compilation_engine: compilation_engine_test.exe
	compilation_engine_test.exe -f projects\10\ArrayTest\Main.jack
	tools\TextComparer.bat projects\10\ArrayTest\Main_KM.xml projects\10\ArrayTest\Main.xml

	compilation_engine_test.exe -d projects\10\ExpressionLessSquare
	tools\TextComparer.bat projects\10\ExpressionLessSquare\Main_KM.xml projects\10\ExpressionLessSquare\Main.xml
	tools\TextComparer.bat projects\10\ExpressionLessSquare\Square_KM.xml projects\10\ExpressionLessSquare\Square.xml
	tools\TextComparer.bat projects\10\ExpressionLessSquare\SquareGame_KM.xml projects\10\ExpressionLessSquare\SquareGame.xml

	compilation_engine_test.exe -d projects\10\Square
	tools\TextComparer.bat projects\10\Square\Main_KM.xml projects\10\Square\Main.xml
	tools\TextComparer.bat projects\10\Square\Square_KM.xml projects\10\Square\Square.xml
	tools\TextComparer.bat projects\10\Square\SquareGame_KM.xml projects\10\Square\SquareGame.xml

compiler.exe: executable\compiler_test\main.go compiler\tokenizer.go compiler\compilation_engine_vm.go compiler\symbol_table.go compiler\vm_writer.go
	go build -o compiler.exe executable\compiler_test\main.go
test_compiler: compiler.exe
	compiler.exe -f projects\11\Average\Main.jack
	compiler.exe -f projects\11\ComplexArrays\Main.jack
	compiler.exe -f projects\11\ConvertToBin\Main.jack
	compiler.exe -d projects\11\Pong
	compiler.exe -f projects\11\Seven\Main.jack
	compiler.exe -d projects\11\Square
	