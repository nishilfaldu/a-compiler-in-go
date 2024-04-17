@y = global i64 0
@tmp = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        %0 = call i64 @proc1(i64* @y)
        store i64 %0, i64* @y
        %1 = call i1 @putinteger(i64* @y)
        store i1 %1, i1* @tmp
        ret i64 0
}

define i64 @proc1(i64* %val) {
proc1:
        %0 = load i64, i64* %val
        %1 = add i64 %0, 1
        %2 = load i64*, i64* %val
        store i64 %1, i64* %val
        %3 = call i64 @proc2(i64* %val)
        %4 = load i64*, i64* %val
        store i64 %3, i64* %val
        %5 = load i64, i64* %val
        ret i64 %5
}

define i64 @proc2(i64* %val) {
proc2:
        %0 = load i64, i64* %val
        %1 = add i64 %0, 1
        %2 = load i64*, i64* %val
        store i64 %1, i64* %val
        %3 = call i64 @proc3(i64* %val)
        %4 = load i64*, i64* %val
        store i64 %3, i64* %val
        %5 = load i64, i64* %val
        ret i64 %5
}

define i64 @proc3(i64* %val) {
proc3:
        %0 = load i64, i64* %val
        %1 = add i64 %0, 1
        %2 = load i64*, i64* %val
        store i64 %1, i64* %val
        %3 = load i64, i64* %val
        ret i64 %3
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}

declare i32 @printf(i8* %format)