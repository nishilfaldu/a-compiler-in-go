@i = global i64 0
@c = global [255 x i8]* getelementptr ([1 x i8], [1 x i8]* @c_0, i64 0)
@myarray = global [0 x %!s(<nil>)] [i64 0]
@c_0 = global [1 x i8] c"a"

define i64 @main() {
entry:
        store i64 100, i64* @i
        %0 = icmp sgt i64* @i, 100
        ret i64 0

if.then:
        store i64 1110, i64* @i
        br label %entry
}
