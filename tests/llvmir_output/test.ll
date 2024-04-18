@x = global i64 0
@i = global i64 0
@max = global i64 0
@tmp = global i64 0
@out = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        store i64 5, i64* @max
        %0 = icmp slt i64* @i, @max
        %1 = call i64 @fib(i64* @i)
        store i64 %1, i64* @x
        %2 = call i1 @putinteger(i64* @x)
        store i1 %2, i1* @out
        %3 = load i64, i64* @i
        %4 = add i64 %3, 1
        store i64 %4, i64* @i
        ret i64 0

for.loop.body:
        %5 = phi i64 [ 0, %entry ]
        br i1 %0, label %for.loop.body, label %leave.for.loop

leave.for.loop:
        ret i64 0
}

define i64 @fib(i64* %val) {
fib:
        %tmp = alloca [2 x i64]
        %loopval = alloca i64
        %ret = alloca i64
        %0 = fneg i64 1
        %1 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 0
        store i64 %0, i64* %1
        %2 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        store i64 1, i64* %2
        %3 = icmp sle i64* %loopval, %val
        %4 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 0
        %5 = load i64, i64* %4
        %6 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        %7 = load i64, i64* %6
        %8 = add i64 %5, %7
        store i64 %8, i64* %ret
        %9 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        %10 = load i64, i64* %9
        %11 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 0
        store i64 %10, i64* %11
        %12 = getelementptr [2 x i64], [2 x i64]* %tmp, i64 0, i64 1
        %13 = load i64, i64* %ret
        store i64 %13, i64* %12
        %14 = load i64, i64* %loopval
        %15 = add i64 %14, 1
        store i64 %15, i64* %loopval
        %16 = load i64, i64* %ret
        ret i64 %16

for.loop.body:
        %17 = phi i64 [ 0, %fib ]
        br i1 %3, label %for.loop.body, label %leave.for.loop

leave.for.loop:
        ret void
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}

declare i32 @printf(i8* %format)