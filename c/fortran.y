%token_type {int}

%left PLUS MINUS.
%left DIVIDE TIMES.

%include {
#include <stdio.h>
#include "fortran.h"
}

%syntax_error {
  printf("Syntax error!\n");
}

program ::= expr(A). { printf("Result=%d\n", A); }

expr(A) ::= expr(B) MINUS  expr(C). { A = B - C; }
expr(A) ::= expr(B) PLUS  expr(C). { A = B + C; }
expr(A) ::= expr(B) TIMES  expr(C). { A = B * C; }
expr(A) ::= expr(B) DIVIDE expr(C). {
    if(C != 0){
        A = B / C;
    }else{
        printf("divide by zero\n");
    }
}

expr(A) ::= INTEGER(B). { A = B; }

