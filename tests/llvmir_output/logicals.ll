@putstring.str = global [4 x i8] c"%s\0A\00"
@x = global i1 true
@c = global i8* null
@tmp = global i1 true
@0 = global [2 x i8] c"t\00"
@1 = global [2 x i8] c"t\00"
@2 = global [2 x i8] c"f\00"
@3 = global [2 x i8] c"f\00"
@.strPrompt = global [39 x i8] c"Enter a string (up to 99 characters): \00"
@.strScanfFmt = global [5 x i8] c"%99s\00"
@4 = global [2 x i8] c"a\00"
@5 = global [2 x i8] c"t\00"
@6 = global [2 x i8] c"t\00"
@7 = global [2 x i8] c"f\00"
@8 = global [2 x i8] c"f\00"

define i64 @main() {
entry:
        store i1 true, i1* @x
        %0 = load i1, i1* @x
        br i1 %0, label %if.then0, label %if.else0

if.then0:
        %1 = getelementptr [2 x i8], [2 x i8]* @0, i64 0, i64 0
        %2 = getelementptr [2 x i8], [2 x i8]* @1, i64 0, i64 0
        %3 = call i1 @putstring(i8* %2)
        store i1 %3, i1* @tmp
        br label %leave.if0

if.else0:
        %4 = getelementptr [2 x i8], [2 x i8]* @2, i64 0, i64 0
        %5 = getelementptr [2 x i8], [2 x i8]* @3, i64 0, i64 0
        %6 = call i1 @putstring(i8* %5)
        store i1 %6, i1* @tmp
        br label %leave.if0

leave.if0:
        %7 = call i8* @getstring()
        store i8* %7, i8** @c
        %8 = getelementptr [2 x i8], [2 x i8]* @4, i64 0, i64 0
        %9 = load i8*, i8** @c
        %10 = call i32 @strcmp(i8* %9, i8* %8)
        %11 = icmp eq i32 %10, 0
        br i1 %11, label %if.then1, label %if.else1

if.then1:
        %12 = getelementptr [2 x i8], [2 x i8]* @5, i64 0, i64 0
        %13 = getelementptr [2 x i8], [2 x i8]* @6, i64 0, i64 0
        %14 = call i1 @putstring(i8* %13)
        store i1 %14, i1* @tmp
        br label %leave.if1

if.else1:
        %15 = getelementptr [2 x i8], [2 x i8]* @7, i64 0, i64 0
        %16 = getelementptr [2 x i8], [2 x i8]* @8, i64 0, i64 0
        %17 = call i1 @putstring(i8* %16)
        store i1 %17, i1* @tmp
        br label %leave.if1

leave.if1:
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