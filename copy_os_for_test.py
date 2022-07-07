import shutil

base_dir = r'C:\Users\user\Desktop\Nand2tetris\projects\12'
dst_dir = r'C:\Users\user\Desktop\Nand2tetris\projects\11\Pong'
shutil.copyfile(base_dir + r'\ArrayTest\Array.jack', dst_dir + r'\Array.jack')
shutil.copyfile(base_dir + r'\KeyboardTest\Keyboard.jack', dst_dir + r'\Keyboard.jack')
shutil.copyfile(base_dir + r'\MathTest\Math.jack', dst_dir + r'\Math.jack')
shutil.copyfile(base_dir + r'\MemoryTest\Memory.jack', dst_dir + r'\Memory.jack')
shutil.copyfile(base_dir + r'\OutputTest\Output.jack', dst_dir + r'\Output.jack')
shutil.copyfile(base_dir + r'\ScreenTest\Screen.jack', dst_dir + r'\Screen.jack')
shutil.copyfile(base_dir + r'\StringTest\String.jack', dst_dir + r'\String.jack')
shutil.copyfile(base_dir + r'\SysTest\Sys.jack', dst_dir + r'\Sys.jack')

print('Success')