import shutil
import os

base_dir = r'projects\12'

files = [
    ('ArrayTest', 'Array.jack'),
    ('KeyboardTest', 'Keyboard.jack'),
    ('MathTest', 'Math.jack'),
    ('MemoryTest', 'Memory.jack'),
    ('OutputTest', 'Output.jack'),
    ('ScreenTest', 'Screen.jack'),
    ('StringTest', 'String.jack'),
    ('SysTest', 'Sys.jack'),

    ('ArrayTest', 'Array.vm'),
    ('KeyboardTest', 'Keyboard.vm'),
    ('MathTest', 'Math.vm'),
    ('MemoryTest', 'Memory.vm'),
    ('OutputTest', 'Output.vm'),
    ('ScreenTest', 'Screen.vm'),
    ('StringTest', 'String.vm'),
    ('SysTest', 'Sys.vm'),
]


dst_dirs = [
    r'projects\11\Average',
    r'projects\11\ComplexArrays',
    r'projects\11\ConvertToBin',
    r'projects\11\Pong',
    r'projects\11\Seven',
    r'projects\11\Square',

    r'projects\12\ArrayTest',
    r'projects\12\KeyboardTest',
    r'projects\12\MathTest',
    r'projects\12\MemoryTest',
    r'projects\12\OutputTest',
    r'projects\12\ScreenTest',
    r'projects\12\StringTest',
    r'projects\12\SysTest',

    r'MyApp\DirectRAM',
    r'MyApp\Error',
    r'MyApp\Helloworld',
    r'MyApp\Shell',
]

for dst_dir in dst_dirs:
    for folder, filename in files:
        src_path = os.path.join(base_dir, folder, filename)
        dst_path = os.path.join(dst_dir, filename)
        if src_path == dst_path:
            continue
        shutil.copyfile(src_path, dst_path)
        print(f'Copied {filename} to {dst_dir}')
    print('')

print('All files copied successfully.')
