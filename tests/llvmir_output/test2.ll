@putstring.str = global [4 x i8] c"%s\0A\00"
@fibb_result = global i64 0
@out = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        %0 = sub i64 0, 1234
        store i64 %0, i64* @fibb_result
        %arg0 = alloca i64
        store i64 12, i64* %arg0
        %1 = call i64 @fibb(i64* %arg0)
        store i64 %1, i64* @fibb_result
        %2 = call i1 @putinteger(i64* @fibb_result)
        store i1 %2, i1* @out
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

define i64 @fibb(i64* %n) {
fibb:
        %result = alloca i64
        %a = alloca i64
        %b = alloca i64
        %0 = load i64, i64* %n
        %1 = icmp slt i64 %0, 0
        br i1 %1, label %if.then0, label %if.else0

if.then0:
        %2 = sub i64 0, 1
        store i64 %2, i64* %result
        %3 = load i64, i64* %result
        ret i64 %3

if.else0:
        br label %leave.if0

leave.if0:
        %4 = load i64, i64* %n
        %5 = icmp eq i64 %4, 0
        br i1 %5, label %if.then1, label %if.else1

if.then1:
        store i64 0, i64* %result
        %6 = load i64, i64* %result
        ret i64 %6

if.else1:
        br label %leave.if1

leave.if1:
        %7 = load i64, i64* %n
        %8 = icmp eq i64 %7, 1
        br i1 %8, label %if.then2, label %if.else2

if.then2:
        store i64 1, i64* %result
        %9 = load i64, i64* %result
        ret i64 %9

if.else2:
        br label %leave.if2

leave.if2:
        %10 = load i64, i64* %n
        %11 = sub i64 %10, 1
        %12 = load i64, i64* %n
        %13 = sub i64 %12, 1
        %arg0 = alloca i64
        store i64 %13, i64* %arg0
        %14 = call i64 @fibb(i64* %arg0)
        store i64 %14, i64* %a
        %15 = load i64, i64* %n
        %16 = sub i64 %15, 2
        %17 = load i64, i64* %n
        %18 = sub i64 %17, 2
        %arg1 = alloca i64
        store i64 %18, i64* %arg1
        %19 = call i64 @fibb(i64* %arg1)
        store i64 %19, i64* %b
        %20 = load i64, i64* %a
        %21 = load i64, i64* %b
        %22 = add i64 %20, %21
        store i64 %22, i64* %result
        %23 = load i64, i64* %result
        ret i64 %23
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}