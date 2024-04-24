@a = global i64 0
@tmp = global i1 true
@.textstr = global [4 x i8] c"%d\0A\00"

define i64 @main() {
entry:
        store i64 3, i64* @a
        %arg0 = alloca i64
        store i64 3, i64* %arg0
        %0 = call i64 @f(i64* %arg0)
        store i64 %0, i64* @a
        %1 = call i1 @putinteger(i64* @a)
        store i1 %1, i1* @tmp
        ret i64 0
}

define i64 @f(i64* %x) {
f:
        %b = alloca i64
        %y = alloca i64
        store i64 3, i64* %b
        %0 = load i64, i64* %b
        store i64 %0, i64* %y
        %1 = load i64, i64* %y
        %2 = mul i64 %1, 4
        %3 = load i64, i64* %x
        %4 = add i64 %3, %2
        br label %f.exit

f.exit:
        ret i64 %4
}

define i1 @putinteger(i64* %paramValue) {
putinteger.entry:
        %0 = load i64, i64* %paramValue
        %1 = getelementptr [4 x i8], [4 x i8]* @.textstr, i64 0, i64 0
        %2 = call i32 @printf(i8* %1, i64 %0)
        ret i1 true
}

declare i32 @printf(i8* %format)