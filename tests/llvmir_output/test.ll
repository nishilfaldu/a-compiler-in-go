@value = global i64 0
@tmp2 = global float 0.0
@out = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        %0 = call i64 @fib(i64 15)
        store i64 %0, i64* @value
        %1 = call i1 @putinteger(i64* @value)
        store i1 %1, i1* @out
        ret i64 0
}

define i64 @fib(i64 %val) {
fib:
        %tmp1 = alloca i64
        %tmp2 = alloca i64
        %ret = alloca i64
        %0 = icmp eq i64 %val, 0
        store i64 0, i64* %ret
        %1 = icmp eq i64 %val, 1
        store i64 1, i64* %ret
        %2 = sub i64 %val, 1
        %3 = load i64*, i64 %2
        store i64 %2, i64* %3
        %4 = call i64 @fib(i64 %val)
        store i64 %4, i64* %tmp1
        %5 = sub i64 %val, 1
        %6 = load i64*, i64 %5
        store i64 %5, i64* %6
        %7 = call i64 @fib(i64 %val)
        store i64 %7, i64* %tmp2
        %8 = load i64, i64* %tmp1
        %9 = add i64 %8, %tmp2
        store i64 %9, i64* %ret
        ret i64* %ret
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}

declare i32 @printf(i8* %format)