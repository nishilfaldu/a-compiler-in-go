@putstring.str = global [4 x i8] c"%s\0A\00"
@value = global i64 0
@tmp2 = global float 0.0
@out = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        %arg0 = alloca i64
        store i64 15, i64* %arg0
        %0 = call i64 @fib(i64* %arg0)
        store i64 %0, i64* @value
        %1 = call i1 @putinteger(i64* @value)
        store i1 %1, i1* @out
        ret i64 0
}

declare i32 @printf(i8* %format)

declare i32 @scanf(i8* %format)

define i1 @putstring(i8* %paramValue) {
putstring.entry:
        %0 = getelementptr [4 x i8], [4 x i8]* @putstring.str, i64 0, i64 0
        %1 = call i32 @printf(i8* %0, i8* %paramValue)
        ret i1 true
}

declare i32 @strcmp(i8* %s1, i8* %s2)

define i64 @fib(i64* %val) {
fib:
        %tmp1 = alloca i64
        %tmp2 = alloca i64
        %ret = alloca i64
        %0 = load i64, i64* %val
        %1 = icmp eq i64 %0, 0
        br i1 %1, label %if.then0, label %if.else0

if.then0:
        store i64 0, i64* %ret
        %2 = load i64, i64* %ret
        ret i64 %2

if.else0:
        br label %leave.if0

leave.if0:
        %3 = load i64, i64* %val
        %4 = icmp eq i64 %3, 1
        br i1 %4, label %if.then1, label %if.else1

if.then1:
        store i64 1, i64* %ret
        %5 = load i64, i64* %ret
        ret i64 %5

if.else1:
        br label %leave.if1

leave.if1:
        %6 = load i64, i64* %val
        %7 = sub i64 %6, 1
        store i64 %7, i64* %val
        %8 = call i64 @fib(i64* %val)
        store i64 %8, i64* %tmp1
        %9 = load i64, i64* %val
        %10 = sub i64 %9, 1
        store i64 %10, i64* %val
        %11 = call i64 @fib(i64* %val)
        store i64 %11, i64* %tmp2
        %12 = load i64, i64* %tmp1
        %13 = load i64, i64* %tmp2
        %14 = add i64 %12, %13
        store i64 %14, i64* %ret
        %15 = load i64, i64* %ret
        ret i64 %15
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}