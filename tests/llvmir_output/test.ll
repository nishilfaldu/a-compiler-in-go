@putstring.str = global [4 x i8] c"%s\0A\00"
@x = global i1 true
@c = global i8* null
@tmp = global i1 true
@__stdinp = global i8* null

define i64 @main() {
entry:
        store i1 false, i1* @x
        %0 = load i1, i1* @x
        br i1 %0, label %if.then0, label %if.else0

if.then0:
        %1 = call i1 @putstring([1 x i8] c"t")
        store i1 %1, i1* @tmp
        br label %leave.if0

if.else0:
        %2 = call i1 @putstring([1 x i8] c"f")
        store i1 %2, i1* @tmp
        br label %leave.if0

leave.if0:
        %3 = call i8* @getstring()
        store i8* %3, i8** @c
        %4 = icmp eq i8** @c, c"a"
        br i1 %4, label %if.then1, label %if.else1

if.then1:
        %5 = call i1 @putstring([1 x i8] c"t")
        store i1 %5, i1* @tmp
        br label %leave.if1

if.else1:
        %6 = call i1 @putstring([1 x i8] c"f")
        store i1 %6, i1* @tmp
        br label %leave.if1

leave.if1:
        ret i64 0
}

define i1 @putstring(i8* %paramValue) {
putstring.entry:
        %0 = getelementptr [4 x i8], [4 x i8]* @putstring.str, i64 0, i64 0
        %1 = call i32 @printf(i8* %0, i8* %paramValue)
        ret i1 true
}

declare i32 @printf(i8* %format)

declare i64 @getline(i8* %buf, i64* %size, i8* %file)

define i8* @getstring() {
entry:
        %0 = alloca i8*
        %1 = alloca i64
        %2 = alloca i64
        store i8* null, i8** %0
        store i64 0, i64* %1
        %3 = load i8*, i8* @__stdinp
        %4 = call i64 @getline(i8** %0, i64* %1, i8* %3)
        store i64 %4, i64* %2
        %5 = load i8*, i8** %0
        %6 = load i64, i64* %2
        %7 = sub i64 %6, 1
        %8 = getelementptr i8, i8* %5, i64 %7
        store i8 0, i8* %8
        %9 = load i8*, i8** %0
        ret i8* %9
}