grammar west;

program: program_header program_body '.';

program_header: 'program' identifier 'is';

program_body:
	(declaration ';')* 'begin' (statement ';')* 'end program';

declaration:
	('global'?) (procedure_declaration | variable_declaration);

procedure_declaration: procedure_header procedure_body;

procedure_header:
	'procedure' identifier ':' type_mark '(' (parameter_list)? ')';

parameter_list: parameter (',' parameter)* | parameter;

parameter: variable_declaration;

procedure_body:
	(declaration ';')* 'begin' (statement ';')* 'end procedure';

variable_declaration:
	'variable' identifier ':' type_mark ('[' bound ']')?;

type_mark: 'integer' | 'float' | 'string' | 'bool';

bound: number;

statement:
	assignment_statement
	| if_statement
	| loop_statement
	| return_statement;

procedure_call: identifier '(' (argument_list)? ')';

// expression here is "Value" in ast.go
assignment_statement: destination ':=' expression;

destination: identifier ('[' expression ']')?;

if_statement:
	'if' '(' expression ')' 'then' (statement ';')* (
		'else' (statement ';')*
	)? 'end if';

loop_statement:
	'for' '(' assignment_statement ';' expression ')' (
		statement ';'
	)* 'end for';

return_statement: 'return' expression;

identifier:[a-zA-Z][a-zA-Z0-9_]*;

expression: ('not'?) arithOp ( '&' | '|' expression)*;

arithOp: relation ( '+' | '-' arithOp)*;

relation: term ( '<' | '>=' | '<=' | '>' | '==' | '!=' term)*;

term: factor ( '*' | '/' factor)*;

factor:
	'(' expression ')'
	| procedure_call
	| ('-')? (name | number | string | 'true' | 'false');

name: identifier ('[' expression ']')?;

//todo
// argument_list: expression (',' expression)* | expression;
argument_list: expression (',' expression)*;

number:[0-9][0-9_]*[.[0-9_]*];

string: '"'[^"]* '"' ;


