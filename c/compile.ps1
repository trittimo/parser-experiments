..\tools\lemon\lemon.exe fortran.y

Import-Module "C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\Tools\Microsoft.VisualStudio.DevShell.dll";
Enter-VsDevShell d4d6c7a0 -SkipAutomaticLocation

cl "parser.c"