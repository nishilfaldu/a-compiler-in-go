@a = global i64 0
@tmp = global i1 true
@mainReturn = global i64 0

define i64 @main() {
entry:
        %0 = load i64, i64 3
        store i64 %0, i64* @a
        %1 = call i64 @f(i64 3)
        %2 = load i64, i64 %1
        store i64 %2, i64* @a
        %3 = load i1, i32 %0
        store i1 %3, i1* @tmp
        ret i64* @mainReturn
}

define i64 @f(i64 %x) {
f:
        %b = alloca i64
        %y = alloca i64
        %0 = load i64, i64 3
        store i64 %0, i64* %b
        %1 = load i64, i64* %b
        store i64 %1, i64* %y
        %2 = mul i64* %y, 4
        %3 = add i64 %x, %2
        ret i64 %3
}

define i1 @putinteger(i64 %paramValue) {
putinteger.entry:
        %0 = call i32 @printf([6 x i8] c"%lld\0A\00", i64* @a)
        ret i1 true
}

declare i32 @printf(i8* %format)