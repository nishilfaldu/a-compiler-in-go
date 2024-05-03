@putstring.str = global [4 x i8] c"%s\0A\00"
@tmp = global i64 0
@0 = global [16 x i8] c"enter a string:\00"
@1 = global [16 x i8] c"enter a string:\00"
@.strPrompt = global [39 x i8] c"Enter a string (up to 99 characters): \00"
@.strScanfFmt = global [5 x i8] c"%99s\00"

define i64 @main() {
entry:
        %arg1 = alloca i64
        store i64 2, i64* %arg1
        %0 = call i64 @print_string(i64* %arg1)
        store i64 %0, i64* @tmp
        ret i64 0
}

declare i32 @printf(i8* %format)

declare i32 @scanf(i8* %format, ...)

define i1 @putstring(i8* %paramValue) {
putstring.entry:
        %0 = getelementptr [4 x i8], [4 x i8]* @putstring.str, i64 0, i64 0
        %1 = call i32 @printf(i8* %0, i8* %paramValue)
        ret i1 true
}

declare i32 @strcmp(i8* %s1, i8* %s2)

define i64 @print_string(i64* %level) {
print_string:
        %s = alloca i8*
        %x = alloca i1
        %0 = getelementptr [16 x i8], [16 x i8]* @0, i64 0, i64 0
        %1 = getelementptr [16 x i8], [16 x i8]* @1, i64 0, i64 0
        %2 = call i1 @putstring(i8* %1)
        store i1 %2, i1* %x
        %3 = call i8* @getstring()
        store i8* %3, i8** %s
        %4 = load i64, i64* %level
        %5 = icmp slt i64 %4, 3
        br i1 %5, label %if.then0, label %if.else0

if.then0:
        %6 = load i64, i64* %level
        %7 = add i64 %6, 1
        %8 = load i64, i64* %level
        %9 = add i64 %8, 1
        %arg0 = alloca i64
        store i64 %9, i64* %arg0
        %10 = call i64 @print_string(i64* %arg0)
        store i64 %10, i64* @tmp
        br label %leave.if0

if.else0:
        br label %leave.if0

leave.if0:
        %11 = call i1 @putstring(i8** %s)
        store i1 %11, i1* %x
        %12 = load i1, i1* %x
        %13 = zext i1 %12 to i64
        ret i64 %13
}

define i8* @getstring() {
getstring.entry:
        %str = alloca [100 x i8], align 16
        %0 = getelementptr [39 x i8], [39 x i8]* @.strPrompt, i64 0, i64 0
        %1 = call i32 @printf(i8* %0)
        %2 = getelementptr [5 x i8], [5 x i8]* @.strScanfFmt, i64 0, i64 0
        %3 = bitcast [100 x i8]* %str to i8*
        %4 = call i32 (i8*, ...) @scanf(i8* %2, i8* %3)
        ret i8* %3
}