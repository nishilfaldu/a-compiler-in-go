@i = global i64 0
@c = global [255 x i8]* getelementptr ([1 x i8], [1 x i8]* @c_0, i64 0)
@c_0 = global [1 x i8] c"a"

define i64 @main() {
entry:
        store i64 100, i64* @i
        ret i64 0

}
