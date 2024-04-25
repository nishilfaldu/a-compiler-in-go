@putstring.str = global [4 x i8] c"%s\0A\00"
@x = global i64 0
@i = global i64 0
@max = global i64 0
@out = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        store i64 5, i64* @max
        store i64 0, i64* @i
        br label %for.cond

for.cond:
        %0 = load i64, i64* @i
        %1 = load i64, i64* @max
        %2 = icmp slt i64 %0, %1
        br i1 %2, label %for.body, label %leave.for.loop

for.body:
        %3 = call i64 @fib(i64* @i)
        store i64 %3, i64* @x
        %4 = call i1 @putinteger(i64* @x)
        store i1 %4, i1* @out
        %5 = load i64, i64* @i
        %6 = add i64 %5, 1
        store i64 %6, i64* @i
        br label %for.cond

leave.for.loop:
        ret i64 0
}

define i1 @putstring(i8* %paramValue) {
putstring.entry:
        %0 = getelementptr [4 x i8], [4 x i8]* @putstring.str, i64 0, i64 0
        %1 = call i32 @printf(i8* %0, i8* %paramValue)
        ret i1 true
}

declare i32 @printf(i8* %format)

declare i32 @strcmp(i8* %s1, i8* %s2)

define i64 @fib(i64* %val) {
fib:
        %0 = load i64, i64* %val
        %1 = icmp eq i64 %0, 0
        br i1 %1, label %if.then0, label %if.else0

if.then0:
        ret i64 0

if.else0:
        br label %leave.if0

leave.if0:
        %2 = load i64, i64* %val
        %3 = icmp eq i64 %2, 1
        br i1 %3, label %if.then1, label %if.else1

if.then1:
        ret i64 1

if.else1:
        br label %leave.if1

leave.if1:
        %4 = call i64 @sub(i64* %val)
        %5 = call i64 @sub(i64* %val)
        %arg0 = alloca i64
        store i64 %5, i64* %arg0
        %6 = call i64 @fib(i64* %arg0)
        %7 = load i64, i64* %val
        %8 = add i64 %7, %6
        ret i64 %8
}

define i64 @sub(i64* %val1) {
sub:
        %0 = load i64, i64* %val1
        %1 = sub i64 %0, 1
        ret i64 %1
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}