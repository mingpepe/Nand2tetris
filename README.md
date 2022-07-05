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


## Hack RAM layout

0000 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;SP \
0001 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;LCL \
0002 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;ARG \
0003 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;THIS \
0004 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;THAT \
0005 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Temp start \
.......... &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;\
0012 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Temp end\
0013 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;General start\
.......... &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;\
0015 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;General end\
0016 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Static start\
.......... &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;\
0255 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Static end\
0256 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Stack start\
.......... &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;\
2047 &ensp;|&ensp;&ensp;&ensp;&ensp;| &ensp;Stack end