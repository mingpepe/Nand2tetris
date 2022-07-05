# Nand2teris

Implementation of https://www.nand2tetris.org/

Data in directory `doc`, `projects`, `tools` all comes from the official site

## Todo

- [x] Project 6 Assembly
- [x] Project 7 VM 1
- [x] Project 8 VM 2
- [x] Project 9 Programming language
- [ ] Project 10 Compiler 1
- [ ] Project 11 Compiler 2
- [ ] Project 12 OS


## Size
* Stack : 1792
* Heap : 1433
* Screen : 8192
* Keyboard : 1
## Hack RAM layout

00000 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;SP \
00001 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;LCL \
00002 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;ARG \
00003 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;THIS \
00004 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;THAT \
00005 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Temp \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
00013 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;General \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
00016 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Static \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
00256 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Stack \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
02048 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Heap \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
16384 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Screen\
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;\
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
..............&ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp; \
24576 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Keyboard