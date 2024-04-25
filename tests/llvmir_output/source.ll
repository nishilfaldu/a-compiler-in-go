@putstring.str = global [4 x i8] c"%s\0A\00"
@i = global i64 0
@c = global i8* getelementptr (i8*, [1 x i8]* @c_0, i64 0)
@myarray = global [15 x i64] zeroinitializer
@0 = global [2 x i8] c"a\00"
@c_0 = global [1 x i8] c"a"

define i64 @main() {
entry:
        store i64 100, i64* @i
        %0 = getelementptr [2 x i8], [2 x i8]* @0, i64 0, i64 0
        %1 = load i64, i64* @i
        %2 = icmp sgt i64 %1, 100
        br i1 %2, label %if.then0, label %if.else0

if.then0:
        store i64 1110, i64* @i
        br label %leave.if0

if.else0:
        br label %leave.if0

leave.if0:
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