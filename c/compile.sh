#!/usr/bin/env bash
dir=$(dirname "$0")
$dir/../tools/lemon/lemon "$dir/fortran.y"
gcc "$dir/parser.c" -o "$dir/parser"