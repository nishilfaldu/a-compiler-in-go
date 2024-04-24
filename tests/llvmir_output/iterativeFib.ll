@x = global i64 0
@i = global i64 0
@max = global i64 0
@tmp = global i64 0
@out = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        store i64 20, i64* @max
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

define i64 @fib(i64* %val) {
fib:
        %tmp = alloca [2 x i64]
        %loopval = alloca i64
        %ret = alloca i64
        %0 = sub i64 0, 1
        %1 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 0
        store i64 %0, i64* %1
        %2 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        store i64 1, i64* %2
        store i64 0, i64* %loopval
        br label %for.cond

for.cond:
        %3 = load i64, i64* %loopval
        %4 = load i64, i64* %val
        %5 = icmp sle i64 %3, %4
        br i1 %5, label %for.body, label %leave.for.loop

for.body:
        %6 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 0
        %7 = load i64, i64* %6
        %8 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        %9 = load i64, i64* %8
        %10 = add i64 %7, %9
        store i64 %10, i64* %ret
        %11 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        %12 = load i64, i64* %11
        %13 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 0
        store i64 %12, i64* %13
        %14 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        %15 = load i64, i64* %ret
        store i64 %15, i64* %14
        %16 = load i64, i64* %loopval
        %17 = add i64 %16, 1
        store i64 %17, i64* %loopval
        br label %for.cond

leave.for.loop:
        %18 = load i64, i64* %ret
        ret i64 %18
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}

declare i32 @printf(i8* %format)