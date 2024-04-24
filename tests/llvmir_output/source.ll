@i = global i64 0
@c = global [255 x i8]* getelementptr ([1 x i8], [1 x i8]* @c_0, i64 0)
@myarray = global [15 x i64] zeroinitializer
@c_0 = global [1 x i8] c"a"

define i64 @main() {
entry:
        store i64 100, i64* @i
        %0 = load i64, i64* @i
        %1 = icmp sgt i64 %0, 100
        br i1 %1, label %if.then, label %if.else

if.then:
        store i64 1110, i64* @i
        br label %leave.if

if.else:
        br label %leave.if

leave.if:
        ret i64 0
}