@putstring.str = global [4 x i8] c"%s\0A\00"
@jake = global i64 0
@ryan = global [3 x i64] zeroinitializer
@zach = global i64 0
@tmp = global i64 0

define i64 @main() {
entry:
        %0 = call i64 @if_proc()
        store i64 %0, i64* @tmp
        %1 = call i64 @for_proc()
        store i64 %1, i64* @tmp
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

define i64 @if_proc() {
if_proc:
        %declaration = alloca i64
        br i1 true, label %if.then0, label %if.else0

if.then0:
        %0 = load i64, i64* @jake
        %1 = add i64 %0, 1
        store i64 %1, i64* @jake
        br label %leave.if0

if.else0:
        %2 = getelementptr [3 x i64], [3 x i64]* @ryan, i64 0, i64 2
        %3 = load i64, i64* %2
        %4 = load i64, i64* @zach
        %5 = add i64 %4, %3
        store i64 %5, i64* @zach
        br label %leave.if0

leave.if0:
        ret i64 0
}

define i64 @for_proc() {
for_proc:
        %i = alloca i64
        store i64 0, i64* %i
        br label %for.cond

for.cond:
        %0 = load i64, i64* %i
        %1 = load i64, i64* @zach
        %2 = icmp slt i64 %0, %1
        br i1 %2, label %for.body, label %leave.for.loop

for.body:
        %3 = load i64, i64* @zach
        %4 = load i64, i64* %i
        %5 = add i64 %3, %4
        %6 = getelementptr [3 x i64], [3 x i64]* @ryan, i64 0, i64 1
        store i64 %5, i64* %6
        br label %for.cond

leave.for.loop:
        ret i64 0
}