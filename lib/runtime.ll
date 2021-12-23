; ModuleID = 'llvm-link'
source_filename = "llvm-link"
target datalayout = "e-m:w-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-pc-windows-msvc19.30.30706"

%struct._iobuf = type { i8* }
%struct.__crt_locale_pointers = type { %struct.__crt_locale_data*, %struct.__crt_multibyte_data* }
%struct.__crt_locale_data = type opaque
%struct.__crt_multibyte_data = type opaque
%struct.RefCount = type { i8, i32 }
%struct.PVHead = type { %struct.RefCount, i32, %struct.PVH* }
%struct.PVH = type { i8, i32, i32 }
%struct.Structure = type { %struct.RefCount, i16 }
%struct.Closure = type { %struct.RefCount, i8*, i16 }
%struct.PVNode = type { %struct.PVH, i32*, [32 x %struct.PVH*] }
%struct.PVLeaf_uint16 = type { %struct.PVH, [32 x i16] }

$fprintf = comdat any

$_vfprintf_l = comdat any

$__local_stdio_printf_options = comdat any

$printf = comdat any

$"??_C@_0P@LEGHBOFO@malloc?5failed?6?$AA@" = comdat any

$"??_C@_0P@FJBEAGCL@calloc?5failed?6?$AA@" = comdat any

$"??_C@_0BF@FELBPIPP@too?5many?5references?6?$AA@" = comdat any

$"??_C@_0BE@HGGMDAAD@invalid?5pointer?5?$CFp?6?$AA@" = comdat any

$"??_C@_04PEDNGLFL@?$CFld?6?$AA@" = comdat any

$"??_C@_03PPOCCAPH@?$CFf?6?$AA@" = comdat any

$"??_C@_05LFIOBDML@true?6?$AA@" = comdat any

$"??_C@_06NIOGPBNO@false?6?$AA@" = comdat any

$"??_C@_02HAOIJKIC@?$CFc?$AA@" = comdat any

$"??_C@_04GJDJEMBE@?$CFc?$CFc?$AA@" = comdat any

$"??_C@_06CBJCAPLI@?$CFc?$CFc?$CFc?$AA@" = comdat any

$"??_C@_08MMPFKENM@?$CFc?$CFc?$CFc?$CFc?$AA@" = comdat any

$"??_C@_01EEMJAFIK@?6?$AA@" = comdat any

$"??_C@_0BL@BKEFFINM@too?5many?5pnode?5references?6?$AA@" = comdat any

$"??_C@_0CO@JDLAEGEB@pvector?5index?5out?5of?5bounds?3?5got@" = comdat any

$"??_C@_0BD@DJOMGDD@overflow?5required?6?$AA@" = comdat any

$"??_C@_0BM@CICGOBLB@join?5error?0?5depth?5mismatch?6?$AA@" = comdat any

$"??_C@_0CF@JJIBIHF@balance?5failed?5to?5compress?5a?5vec@" = comdat any

@"??_C@_0P@LEGHBOFO@malloc?5failed?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [15 x i8] c"malloc failed\0A\00", comdat, align 1
@"??_C@_0P@FJBEAGCL@calloc?5failed?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [15 x i8] c"calloc failed\0A\00", comdat, align 1
@"??_C@_0BF@FELBPIPP@too?5many?5references?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [21 x i8] c"too many references\0A\00", comdat, align 1
@"??_C@_0BE@HGGMDAAD@invalid?5pointer?5?$CFp?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [20 x i8] c"invalid pointer %p\0A\00", comdat, align 1
@__local_stdio_printf_options._OptionsStorage = internal global i64 0, align 8
@"??_C@_04PEDNGLFL@?$CFld?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [5 x i8] c"%ld\0A\00", comdat, align 1
@"??_C@_03PPOCCAPH@?$CFf?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [4 x i8] c"%f\0A\00", comdat, align 1
@"??_C@_05LFIOBDML@true?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [6 x i8] c"true\0A\00", comdat, align 1
@"??_C@_06NIOGPBNO@false?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [7 x i8] c"false\0A\00", comdat, align 1
@"??_C@_02HAOIJKIC@?$CFc?$AA@" = linkonce_odr dso_local unnamed_addr constant [3 x i8] c"%c\00", comdat, align 1
@"??_C@_04GJDJEMBE@?$CFc?$CFc?$AA@" = linkonce_odr dso_local unnamed_addr constant [5 x i8] c"%c%c\00", comdat, align 1
@"??_C@_06CBJCAPLI@?$CFc?$CFc?$CFc?$AA@" = linkonce_odr dso_local unnamed_addr constant [7 x i8] c"%c%c%c\00", comdat, align 1
@"??_C@_08MMPFKENM@?$CFc?$CFc?$CFc?$CFc?$AA@" = linkonce_odr dso_local unnamed_addr constant [9 x i8] c"%c%c%c%c\00", comdat, align 1
@"??_C@_01EEMJAFIK@?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [2 x i8] c"\0A\00", comdat, align 1
@"??_C@_0BL@BKEFFINM@too?5many?5pnode?5references?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [27 x i8] c"too many pnode references\0A\00", comdat, align 1
@"??_C@_0CO@JDLAEGEB@pvector?5index?5out?5of?5bounds?3?5got@" = linkonce_odr dso_local unnamed_addr constant [46 x i8] c"pvector index out of bounds: got %d, size %d\0A\00", comdat, align 1
@"??_C@_0BD@DJOMGDD@overflow?5required?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [19 x i8] c"overflow required\0A\00", comdat, align 1
@"??_C@_0BM@CICGOBLB@join?5error?0?5depth?5mismatch?6?$AA@" = linkonce_odr dso_local unnamed_addr constant [28 x i8] c"join error, depth mismatch\0A\00", comdat, align 1
@"??_C@_0CF@JJIBIHF@balance?5failed?5to?5compress?5a?5vec@" = linkonce_odr dso_local unnamed_addr constant [37 x i8] c"balance failed to compress a vector\0A\00", comdat, align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8* @heap_malloc(i32 %0) #0 {
  %2 = alloca i32, align 4
  %3 = alloca i8*, align 8
  store i32 %0, i32* %2, align 4
  %4 = load i32, i32* %2, align 4
  %5 = sext i32 %4 to i64
  %6 = call noalias align 16 i8* @malloc(i64 %5)
  store i8* %6, i8** %3, align 8
  %7 = load i8*, i8** %3, align 8
  %8 = icmp eq i8* %7, null
  br i1 %8, label %9, label %12

9:                                                ; preds = %1
  %10 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %11 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %10, i8* getelementptr inbounds ([15 x i8], [15 x i8]* @"??_C@_0P@LEGHBOFO@malloc?5failed?6?$AA@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

12:                                               ; preds = %1
  %13 = load i8*, i8** %3, align 8
  ret i8* %13
}

declare dso_local noalias i8* @malloc(i64) #1

declare dso_local %struct._iobuf* @__acrt_iob_func(i32) #1

; Function Attrs: noinline nounwind optnone uwtable
define linkonce_odr dso_local i32 @fprintf(%struct._iobuf* %0, i8* %1, ...) #0 comdat {
  %3 = alloca i8*, align 8
  %4 = alloca %struct._iobuf*, align 8
  %5 = alloca i32, align 4
  %6 = alloca i8*, align 8
  store i8* %1, i8** %3, align 8
  store %struct._iobuf* %0, %struct._iobuf** %4, align 8
  %7 = bitcast i8** %6 to i8*
  call void @llvm.va_start(i8* %7)
  %8 = load i8*, i8** %6, align 8
  %9 = load i8*, i8** %3, align 8
  %10 = load %struct._iobuf*, %struct._iobuf** %4, align 8
  %11 = call i32 @_vfprintf_l(%struct._iobuf* %10, i8* %9, %struct.__crt_locale_pointers* null, i8* %8)
  store i32 %11, i32* %5, align 4
  %12 = bitcast i8** %6 to i8*
  call void @llvm.va_end(i8* %12)
  %13 = load i32, i32* %5, align 4
  ret i32 %13
}

; Function Attrs: noreturn
declare dso_local void @exit(i32) #2

; Function Attrs: nofree nosync nounwind willreturn
declare void @llvm.va_start(i8*) #3

; Function Attrs: noinline nounwind optnone uwtable
define linkonce_odr dso_local i32 @_vfprintf_l(%struct._iobuf* %0, i8* %1, %struct.__crt_locale_pointers* %2, i8* %3) #0 comdat {
  %5 = alloca i8*, align 8
  %6 = alloca %struct.__crt_locale_pointers*, align 8
  %7 = alloca i8*, align 8
  %8 = alloca %struct._iobuf*, align 8
  store i8* %3, i8** %5, align 8
  store %struct.__crt_locale_pointers* %2, %struct.__crt_locale_pointers** %6, align 8
  store i8* %1, i8** %7, align 8
  store %struct._iobuf* %0, %struct._iobuf** %8, align 8
  %9 = load i8*, i8** %5, align 8
  %10 = load %struct.__crt_locale_pointers*, %struct.__crt_locale_pointers** %6, align 8
  %11 = load i8*, i8** %7, align 8
  %12 = load %struct._iobuf*, %struct._iobuf** %8, align 8
  %13 = call i64* @__local_stdio_printf_options()
  %14 = load i64, i64* %13, align 8
  %15 = call i32 @__stdio_common_vfprintf(i64 %14, %struct._iobuf* %12, i8* %11, %struct.__crt_locale_pointers* %10, i8* %9)
  ret i32 %15
}

; Function Attrs: nofree nosync nounwind willreturn
declare void @llvm.va_end(i8*) #3

; Function Attrs: noinline nounwind optnone uwtable
define linkonce_odr dso_local i64* @__local_stdio_printf_options() #0 comdat {
  ret i64* @__local_stdio_printf_options._OptionsStorage
}

declare dso_local i32 @__stdio_common_vfprintf(i64, %struct._iobuf*, i8*, %struct.__crt_locale_pointers*, i8*) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8* @heap_calloc(i32 %0, i32 %1) #0 {
  %3 = alloca i32, align 4
  %4 = alloca i32, align 4
  %5 = alloca i8*, align 8
  store i32 %1, i32* %3, align 4
  store i32 %0, i32* %4, align 4
  %6 = load i32, i32* %3, align 4
  %7 = sext i32 %6 to i64
  %8 = load i32, i32* %4, align 4
  %9 = sext i32 %8 to i64
  %10 = call noalias align 16 i8* @calloc(i64 %9, i64 %7)
  store i8* %10, i8** %5, align 8
  %11 = load i8*, i8** %5, align 8
  %12 = icmp eq i8* %11, null
  br i1 %12, label %13, label %16

13:                                               ; preds = %2
  %14 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %15 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %14, i8* getelementptr inbounds ([15 x i8], [15 x i8]* @"??_C@_0P@FJBEAGCL@calloc?5failed?6?$AA@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

16:                                               ; preds = %2
  %17 = load i8*, i8** %5, align 8
  ret i8* %17
}

declare dso_local noalias i8* @calloc(i64, i64) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @increase_refcount(%struct.RefCount* %0) #0 {
  %2 = alloca %struct.RefCount*, align 8
  %3 = alloca i32, align 4
  store %struct.RefCount* %0, %struct.RefCount** %2, align 8
  %4 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %5 = icmp ne %struct.RefCount* %4, null
  br i1 %5, label %6, label %24

6:                                                ; preds = %1
  %7 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %8 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %7, i32 0, i32 1
  %9 = load i32, i32* %8, align 4
  store i32 %9, i32* %3, align 4
  %10 = load i32, i32* %3, align 4
  %11 = icmp ne i32 %10, -1
  br i1 %11, label %12, label %23

12:                                               ; preds = %6
  %13 = load i32, i32* %3, align 4
  %14 = icmp eq i32 %13, -2
  br i1 %14, label %15, label %18

15:                                               ; preds = %12
  %16 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %17 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %16, i8* getelementptr inbounds ([21 x i8], [21 x i8]* @"??_C@_0BF@FELBPIPP@too?5many?5references?6?$AA@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

18:                                               ; preds = %12
  %19 = load i32, i32* %3, align 4
  %20 = add i32 %19, 1
  %21 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %22 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %21, i32 0, i32 1
  store i32 %20, i32* %22, align 4
  br label %23

23:                                               ; preds = %18, %6
  br label %24

24:                                               ; preds = %23, %1
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @free_rc(%struct.RefCount* %0) #0 {
  %2 = alloca %struct.RefCount*, align 8
  %3 = alloca i8, align 1
  store %struct.RefCount* %0, %struct.RefCount** %2, align 8
  %4 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %5 = icmp eq %struct.RefCount* %4, null
  br i1 %5, label %6, label %7

6:                                                ; preds = %1
  br label %43

7:                                                ; preds = %1
  %8 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %9 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %8, i32 0, i32 1
  %10 = load i32, i32* %9, align 4
  %11 = icmp eq i32 %10, -1
  br i1 %11, label %12, label %13

12:                                               ; preds = %7
  br label %43

13:                                               ; preds = %7
  %14 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %15 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %14, i32 0, i32 0
  %16 = load i8, i8* %15, align 4
  store i8 %16, i8* %3, align 1
  %17 = load i8, i8* %3, align 1
  %18 = zext i8 %17 to i32
  %19 = icmp eq i32 %18, 3
  br i1 %19, label %20, label %23

20:                                               ; preds = %13
  %21 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %22 = bitcast %struct.RefCount* %21 to %struct.PVHead*
  call void @pv_free(%struct.PVHead* %22)
  br label %43

23:                                               ; preds = %13
  %24 = load i8, i8* %3, align 1
  %25 = zext i8 %24 to i32
  %26 = icmp eq i32 %25, 2
  br i1 %26, label %27, label %30

27:                                               ; preds = %23
  %28 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %29 = bitcast %struct.RefCount* %28 to %struct.Structure*
  call void @free_structure(%struct.Structure* %29)
  br label %42

30:                                               ; preds = %23
  %31 = load i8, i8* %3, align 1
  %32 = zext i8 %31 to i32
  %33 = icmp eq i32 %32, 1
  br i1 %33, label %34, label %37

34:                                               ; preds = %30
  %35 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %36 = bitcast %struct.RefCount* %35 to %struct.Closure*
  call void @free_closure(%struct.Closure* %36)
  br label %41

37:                                               ; preds = %30
  %38 = load %struct.RefCount*, %struct.RefCount** %2, align 8
  %39 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %40 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %39, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @"??_C@_0BE@HGGMDAAD@invalid?5pointer?5?$CFp?6?$AA@", i64 0, i64 0), %struct.RefCount* %38)
  call void @exit(i32 1) #5
  unreachable

41:                                               ; preds = %34
  br label %42

42:                                               ; preds = %41, %27
  br label %43

43:                                               ; preds = %42, %20, %12, %6
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @print_int(i32 %0) #0 {
  %2 = alloca i32, align 4
  store i32 %0, i32* %2, align 4
  %3 = load i32, i32* %2, align 4
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @"??_C@_04PEDNGLFL@?$CFld?6?$AA@", i64 0, i64 0), i32 %3)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define linkonce_odr dso_local i32 @printf(i8* %0, ...) #0 comdat {
  %2 = alloca i8*, align 8
  %3 = alloca i32, align 4
  %4 = alloca i8*, align 8
  store i8* %0, i8** %2, align 8
  %5 = bitcast i8** %4 to i8*
  call void @llvm.va_start(i8* %5)
  %6 = load i8*, i8** %4, align 8
  %7 = load i8*, i8** %2, align 8
  %8 = call %struct._iobuf* @__acrt_iob_func(i32 1)
  %9 = call i32 @_vfprintf_l(%struct._iobuf* %8, i8* %7, %struct.__crt_locale_pointers* null, i8* %6)
  store i32 %9, i32* %3, align 4
  %10 = bitcast i8** %4 to i8*
  call void @llvm.va_end(i8* %10)
  %11 = load i32, i32* %3, align 4
  ret i32 %11
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @print_real(double %0) #0 {
  %2 = alloca double, align 8
  store double %0, double* %2, align 8
  %3 = load double, double* %2, align 8
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @"??_C@_03PPOCCAPH@?$CFf?6?$AA@", i64 0, i64 0), double %3)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @print_bool(i8 %0) #0 {
  %2 = alloca i8, align 1
  store i8 %0, i8* %2, align 1
  %3 = load i8, i8* %2, align 1
  %4 = icmp ne i8 %3, 0
  br i1 %4, label %5, label %7

5:                                                ; preds = %1
  %6 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([6 x i8], [6 x i8]* @"??_C@_05LFIOBDML@true?6?$AA@", i64 0, i64 0))
  br label %9

7:                                                ; preds = %1
  %8 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([7 x i8], [7 x i8]* @"??_C@_06NIOGPBNO@false?6?$AA@", i64 0, i64 0))
  br label %9

9:                                                ; preds = %7, %5
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @print_string(%struct.PVHead* %0) #0 {
  %2 = alloca %struct.PVHead*, align 8
  %3 = alloca i32, align 4
  %4 = alloca i32, align 4
  %5 = alloca i16, align 2
  %6 = alloca i32, align 4
  %7 = alloca i16, align 2
  %8 = alloca i16, align 2
  %9 = alloca i32, align 4
  %10 = alloca i32, align 4
  %11 = alloca i32, align 4
  %12 = alloca i8, align 1
  %13 = alloca i8, align 1
  %14 = alloca i8, align 1
  %15 = alloca i8, align 1
  %16 = alloca i8, align 1
  %17 = alloca i8, align 1
  %18 = alloca i8, align 1
  %19 = alloca i8, align 1
  %20 = alloca i8, align 1
  store %struct.PVHead* %0, %struct.PVHead** %2, align 8
  %21 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %22 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %21, i32 0, i32 1
  %23 = load i32, i32* %22, align 8
  store i32 %23, i32* %3, align 4
  store i32 0, i32* %4, align 4
  br label %24

24:                                               ; preds = %147, %1
  %25 = load i32, i32* %4, align 4
  %26 = load i32, i32* %3, align 4
  %27 = icmp ult i32 %25, %26
  br i1 %27, label %28, label %150

28:                                               ; preds = %24
  %29 = load i32, i32* %4, align 4
  %30 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %31 = call i16 @pv_uint16_get(%struct.PVHead* %30, i32 %29)
  store i16 %31, i16* %5, align 2
  %32 = load i16, i16* %5, align 2
  %33 = zext i16 %32 to i32
  store i32 %33, i32* %6, align 4
  %34 = load i32, i32* %6, align 4
  %35 = icmp ugt i32 %34, 55295
  br i1 %35, label %36, label %65

36:                                               ; preds = %28
  %37 = load i32, i32* %6, align 4
  %38 = icmp ult i32 %37, 57344
  br i1 %38, label %39, label %65

39:                                               ; preds = %36
  %40 = load i16, i16* %5, align 2
  store i16 %40, i16* %7, align 2
  %41 = load i32, i32* %4, align 4
  %42 = add nsw i32 %41, 1
  store i32 %42, i32* %4, align 4
  %43 = load i32, i32* %4, align 4
  %44 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %45 = call i16 @pv_uint16_get(%struct.PVHead* %44, i32 %43)
  store i16 %45, i16* %8, align 2
  %46 = load i16, i16* %7, align 2
  %47 = zext i16 %46 to i32
  %48 = and i32 %47, 63
  %49 = shl i32 %48, 10
  %50 = load i16, i16* %8, align 2
  %51 = zext i16 %50 to i32
  %52 = and i32 %51, 1023
  %53 = or i32 %49, %52
  store i32 %53, i32* %9, align 4
  %54 = load i16, i16* %7, align 2
  %55 = zext i16 %54 to i32
  %56 = ashr i32 %55, 6
  %57 = and i32 %56, 31
  store i32 %57, i32* %10, align 4
  %58 = load i32, i32* %10, align 4
  %59 = add i32 %58, 1
  store i32 %59, i32* %11, align 4
  %60 = load i32, i32* %11, align 4
  %61 = shl i32 %60, 16
  %62 = load i32, i32* %9, align 4
  %63 = or i32 %61, %62
  %64 = trunc i32 %63 to i16
  store i16 %64, i16* %5, align 2
  br label %65

65:                                               ; preds = %39, %36, %28
  %66 = load i32, i32* %6, align 4
  %67 = icmp ule i32 %66, 127
  br i1 %67, label %68, label %72

68:                                               ; preds = %65
  %69 = load i16, i16* %5, align 2
  %70 = zext i16 %69 to i32
  %71 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @"??_C@_02HAOIJKIC@?$CFc?$AA@", i64 0, i64 0), i32 %70)
  br label %146

72:                                               ; preds = %65
  %73 = load i32, i32* %6, align 4
  %74 = icmp ule i32 %73, 2047
  br i1 %74, label %75, label %90

75:                                               ; preds = %72
  %76 = load i32, i32* %6, align 4
  %77 = lshr i32 %76, 6
  %78 = and i32 %77, 31
  %79 = add i32 192, %78
  %80 = trunc i32 %79 to i8
  store i8 %80, i8* %12, align 1
  %81 = load i32, i32* %6, align 4
  %82 = and i32 %81, 63
  %83 = add i32 128, %82
  %84 = trunc i32 %83 to i8
  store i8 %84, i8* %13, align 1
  %85 = load i8, i8* %13, align 1
  %86 = zext i8 %85 to i32
  %87 = load i8, i8* %12, align 1
  %88 = zext i8 %87 to i32
  %89 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @"??_C@_04GJDJEMBE@?$CFc?$CFc?$AA@", i64 0, i64 0), i32 %88, i32 %86)
  br label %145

90:                                               ; preds = %72
  %91 = load i32, i32* %6, align 4
  %92 = icmp ule i32 %91, 65535
  br i1 %92, label %93, label %115

93:                                               ; preds = %90
  %94 = load i32, i32* %6, align 4
  %95 = lshr i32 %94, 12
  %96 = and i32 %95, 15
  %97 = add i32 224, %96
  %98 = trunc i32 %97 to i8
  store i8 %98, i8* %14, align 1
  %99 = load i32, i32* %6, align 4
  %100 = lshr i32 %99, 6
  %101 = and i32 %100, 63
  %102 = add i32 128, %101
  %103 = trunc i32 %102 to i8
  store i8 %103, i8* %15, align 1
  %104 = load i32, i32* %6, align 4
  %105 = and i32 %104, 63
  %106 = add i32 128, %105
  %107 = trunc i32 %106 to i8
  store i8 %107, i8* %16, align 1
  %108 = load i8, i8* %16, align 1
  %109 = zext i8 %108 to i32
  %110 = load i8, i8* %15, align 1
  %111 = zext i8 %110 to i32
  %112 = load i8, i8* %14, align 1
  %113 = zext i8 %112 to i32
  %114 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([7 x i8], [7 x i8]* @"??_C@_06CBJCAPLI@?$CFc?$CFc?$CFc?$AA@", i64 0, i64 0), i32 %113, i32 %111, i32 %109)
  br label %144

115:                                              ; preds = %90
  %116 = load i32, i32* %6, align 4
  %117 = lshr i32 %116, 18
  %118 = and i32 %117, 7
  %119 = add i32 240, %118
  %120 = trunc i32 %119 to i8
  store i8 %120, i8* %17, align 1
  %121 = load i32, i32* %6, align 4
  %122 = lshr i32 %121, 12
  %123 = and i32 %122, 63
  %124 = add i32 128, %123
  %125 = trunc i32 %124 to i8
  store i8 %125, i8* %18, align 1
  %126 = load i32, i32* %6, align 4
  %127 = lshr i32 %126, 6
  %128 = and i32 %127, 63
  %129 = add i32 128, %128
  %130 = trunc i32 %129 to i8
  store i8 %130, i8* %19, align 1
  %131 = load i32, i32* %6, align 4
  %132 = and i32 %131, 63
  %133 = add i32 128, %132
  %134 = trunc i32 %133 to i8
  store i8 %134, i8* %20, align 1
  %135 = load i8, i8* %20, align 1
  %136 = zext i8 %135 to i32
  %137 = load i8, i8* %19, align 1
  %138 = zext i8 %137 to i32
  %139 = load i8, i8* %18, align 1
  %140 = zext i8 %139 to i32
  %141 = load i8, i8* %17, align 1
  %142 = zext i8 %141 to i32
  %143 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([9 x i8], [9 x i8]* @"??_C@_08MMPFKENM@?$CFc?$CFc?$CFc?$CFc?$AA@", i64 0, i64 0), i32 %142, i32 %140, i32 %138, i32 %136)
  br label %144

144:                                              ; preds = %115, %93
  br label %145

145:                                              ; preds = %144, %75
  br label %146

146:                                              ; preds = %145, %68
  br label %147

147:                                              ; preds = %146
  %148 = load i32, i32* %4, align 4
  %149 = add nsw i32 %148, 1
  store i32 %149, i32* %4, align 4
  br label %24, !llvm.loop !4

150:                                              ; preds = %24
  %151 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @"??_C@_01EEMJAFIK@?6?$AA@", i64 0, i64 0))
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8 @pv_depth(%struct.PVHead* %0) #0 {
  %2 = alloca i8, align 1
  %3 = alloca %struct.PVHead*, align 8
  %4 = alloca %struct.PVH*, align 8
  store %struct.PVHead* %0, %struct.PVHead** %3, align 8
  %5 = load %struct.PVHead*, %struct.PVHead** %3, align 8
  %6 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %5, i32 0, i32 2
  %7 = load %struct.PVH*, %struct.PVH** %6, align 8
  store %struct.PVH* %7, %struct.PVH** %4, align 8
  %8 = load %struct.PVH*, %struct.PVH** %4, align 8
  %9 = icmp eq %struct.PVH* %8, null
  br i1 %9, label %10, label %11

10:                                               ; preds = %1
  store i8 0, i8* %2, align 1
  br label %15

11:                                               ; preds = %1
  %12 = load %struct.PVH*, %struct.PVH** %4, align 8
  %13 = getelementptr inbounds %struct.PVH, %struct.PVH* %12, i32 0, i32 0
  %14 = load i8, i8* %13, align 4
  store i8 %14, i8* %2, align 1
  br label %15

15:                                               ; preds = %11, %10
  %16 = load i8, i8* %2, align 1
  ret i8 %16
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @pv_length(%struct.PVHead* %0) #0 {
  %2 = alloca %struct.PVHead*, align 8
  store %struct.PVHead* %0, %struct.PVHead** %2, align 8
  %3 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %4 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %3, i32 0, i32 1
  %5 = load i32, i32* %4, align 8
  ret i32 %5
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @pn_incr_ref(%struct.PVH* %0) #0 {
  %2 = alloca %struct.PVH*, align 8
  store %struct.PVH* %0, %struct.PVH** %2, align 8
  %3 = load %struct.PVH*, %struct.PVH** %2, align 8
  %4 = getelementptr inbounds %struct.PVH, %struct.PVH* %3, i32 0, i32 1
  %5 = load i32, i32* %4, align 4
  %6 = icmp ne i32 %5, -1
  br i1 %6, label %7, label %20

7:                                                ; preds = %1
  %8 = load %struct.PVH*, %struct.PVH** %2, align 8
  %9 = getelementptr inbounds %struct.PVH, %struct.PVH* %8, i32 0, i32 1
  %10 = load i32, i32* %9, align 4
  %11 = icmp eq i32 %10, -2
  br i1 %11, label %12, label %15

12:                                               ; preds = %7
  %13 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %14 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %13, i8* getelementptr inbounds ([27 x i8], [27 x i8]* @"??_C@_0BL@BKEFFINM@too?5many?5pnode?5references?6?$AA@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

15:                                               ; preds = %7
  %16 = load %struct.PVH*, %struct.PVH** %2, align 8
  %17 = getelementptr inbounds %struct.PVH, %struct.PVH* %16, i32 0, i32 1
  %18 = load i32, i32* %17, align 4
  %19 = add i32 %18, 1
  store i32 %19, i32* %17, align 4
  br label %20

20:                                               ; preds = %15, %1
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVHead* @pv_new() #0 {
  %1 = alloca %struct.PVHead*, align 8
  %2 = call i8* @heap_calloc(i32 1, i32 24)
  %3 = bitcast i8* %2 to %struct.PVHead*
  store %struct.PVHead* %3, %struct.PVHead** %1, align 8
  %4 = load %struct.PVHead*, %struct.PVHead** %1, align 8
  %5 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %4, i32 0, i32 0
  %6 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %5, i32 0, i32 1
  store i32 1, i32* %6, align 4
  %7 = load %struct.PVHead*, %struct.PVHead** %1, align 8
  %8 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %7, i32 0, i32 0
  %9 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %8, i32 0, i32 0
  store i8 3, i8* %9, align 8
  %10 = load %struct.PVHead*, %struct.PVHead** %1, align 8
  ret %struct.PVHead* %10
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVHead* @pv_construct(%struct.PVH* %0) #0 {
  %2 = alloca %struct.PVH*, align 8
  %3 = alloca %struct.PVHead*, align 8
  store %struct.PVH* %0, %struct.PVH** %2, align 8
  %4 = call %struct.PVHead* @pv_new()
  store %struct.PVHead* %4, %struct.PVHead** %3, align 8
  %5 = load %struct.PVH*, %struct.PVH** %2, align 8
  %6 = load %struct.PVHead*, %struct.PVHead** %3, align 8
  %7 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %6, i32 0, i32 2
  store %struct.PVH* %5, %struct.PVH** %7, align 8
  %8 = load %struct.PVH*, %struct.PVH** %2, align 8
  %9 = getelementptr inbounds %struct.PVH, %struct.PVH* %8, i32 0, i32 2
  %10 = load i32, i32* %9, align 4
  %11 = load %struct.PVHead*, %struct.PVHead** %3, align 8
  %12 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %11, i32 0, i32 1
  store i32 %10, i32* %12, align 8
  %13 = load %struct.PVH*, %struct.PVH** %2, align 8
  call void @pn_incr_ref(%struct.PVH* %13)
  %14 = load %struct.PVHead*, %struct.PVHead** %3, align 8
  ret %struct.PVHead* %14
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVNode* @pn_new(i8 %0) #0 {
  %2 = alloca i8, align 1
  %3 = alloca %struct.PVNode*, align 8
  store i8 %0, i8* %2, align 1
  %4 = call i8* @heap_calloc(i32 1, i32 280)
  %5 = bitcast i8* %4 to %struct.PVNode*
  store %struct.PVNode* %5, %struct.PVNode** %3, align 8
  %6 = load i8, i8* %2, align 1
  %7 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %8 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %7, i32 0, i32 0
  %9 = getelementptr inbounds %struct.PVH, %struct.PVH* %8, i32 0, i32 0
  store i8 %6, i8* %9, align 8
  %10 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  ret %struct.PVNode* %10
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVH* @pl_new(i32 %0) #0 {
  %2 = alloca i32, align 4
  store i32 %0, i32* %2, align 4
  %3 = load i32, i32* %2, align 4
  %4 = call i8* @heap_calloc(i32 1, i32 %3)
  %5 = bitcast i8* %4 to %struct.PVH*
  ret %struct.PVH* %5
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @pn_free(%struct.PVH* %0) #0 {
  %2 = alloca %struct.PVH*, align 8
  %3 = alloca i32, align 4
  %4 = alloca %struct.PVNode*, align 8
  %5 = alloca i32, align 4
  store %struct.PVH* %0, %struct.PVH** %2, align 8
  %6 = load %struct.PVH*, %struct.PVH** %2, align 8
  %7 = icmp eq %struct.PVH* %6, null
  br i1 %7, label %8, label %9

8:                                                ; preds = %1
  br label %69

9:                                                ; preds = %1
  %10 = load %struct.PVH*, %struct.PVH** %2, align 8
  %11 = getelementptr inbounds %struct.PVH, %struct.PVH* %10, i32 0, i32 1
  %12 = load i32, i32* %11, align 4
  store i32 %12, i32* %3, align 4
  %13 = load i32, i32* %3, align 4
  %14 = icmp eq i32 %13, -1
  br i1 %14, label %15, label %16

15:                                               ; preds = %9
  br label %69

16:                                               ; preds = %9
  %17 = load i32, i32* %3, align 4
  %18 = icmp ugt i32 %17, 1
  br i1 %18, label %19, label %24

19:                                               ; preds = %16
  %20 = load i32, i32* %3, align 4
  %21 = sub i32 %20, 1
  %22 = load %struct.PVH*, %struct.PVH** %2, align 8
  %23 = getelementptr inbounds %struct.PVH, %struct.PVH* %22, i32 0, i32 1
  store i32 %21, i32* %23, align 4
  br label %69

24:                                               ; preds = %16
  %25 = load %struct.PVH*, %struct.PVH** %2, align 8
  %26 = getelementptr inbounds %struct.PVH, %struct.PVH* %25, i32 0, i32 0
  %27 = load i8, i8* %26, align 4
  %28 = zext i8 %27 to i32
  %29 = icmp sgt i32 %28, 0
  br i1 %29, label %30, label %66

30:                                               ; preds = %24
  %31 = load %struct.PVH*, %struct.PVH** %2, align 8
  %32 = bitcast %struct.PVH* %31 to %struct.PVNode*
  store %struct.PVNode* %32, %struct.PVNode** %4, align 8
  store i32 0, i32* %5, align 4
  br label %33

33:                                               ; preds = %52, %30
  %34 = load i32, i32* %5, align 4
  %35 = icmp slt i32 %34, 32
  br i1 %35, label %36, label %55

36:                                               ; preds = %33
  %37 = load %struct.PVNode*, %struct.PVNode** %4, align 8
  %38 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %37, i32 0, i32 2
  %39 = load i32, i32* %5, align 4
  %40 = sext i32 %39 to i64
  %41 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %38, i64 0, i64 %40
  %42 = load %struct.PVH*, %struct.PVH** %41, align 8
  %43 = icmp ne %struct.PVH* %42, null
  br i1 %43, label %44, label %51

44:                                               ; preds = %36
  %45 = load %struct.PVNode*, %struct.PVNode** %4, align 8
  %46 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %45, i32 0, i32 2
  %47 = load i32, i32* %5, align 4
  %48 = sext i32 %47 to i64
  %49 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %46, i64 0, i64 %48
  %50 = load %struct.PVH*, %struct.PVH** %49, align 8
  call void @pn_free(%struct.PVH* %50)
  br label %51

51:                                               ; preds = %44, %36
  br label %52

52:                                               ; preds = %51
  %53 = load i32, i32* %5, align 4
  %54 = add nsw i32 %53, 1
  store i32 %54, i32* %5, align 4
  br label %33, !llvm.loop !6

55:                                               ; preds = %33
  %56 = load %struct.PVNode*, %struct.PVNode** %4, align 8
  %57 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %56, i32 0, i32 1
  %58 = load i32*, i32** %57, align 8
  %59 = icmp ne i32* %58, null
  br i1 %59, label %60, label %65

60:                                               ; preds = %55
  %61 = load %struct.PVNode*, %struct.PVNode** %4, align 8
  %62 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %61, i32 0, i32 1
  %63 = load i32*, i32** %62, align 8
  %64 = bitcast i32* %63 to i8*
  call void @free(i8* %64)
  br label %65

65:                                               ; preds = %60, %55
  br label %66

66:                                               ; preds = %65, %24
  %67 = load %struct.PVH*, %struct.PVH** %2, align 8
  %68 = bitcast %struct.PVH* %67 to i8*
  call void @free(i8* %68)
  br label %69

69:                                               ; preds = %66, %19, %15, %8
  ret void
}

declare dso_local void @free(i8*) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @pv_free(%struct.PVHead* %0) #0 {
  %2 = alloca %struct.PVHead*, align 8
  store %struct.PVHead* %0, %struct.PVHead** %2, align 8
  %3 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %4 = icmp eq %struct.PVHead* %3, null
  br i1 %4, label %11, label %5

5:                                                ; preds = %1
  %6 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %7 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %6, i32 0, i32 0
  %8 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %7, i32 0, i32 1
  %9 = load i32, i32* %8, align 4
  %10 = icmp eq i32 %9, -1
  br i1 %10, label %11, label %12

11:                                               ; preds = %5, %1
  br label %39

12:                                               ; preds = %5
  %13 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %14 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %13, i32 0, i32 0
  %15 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %14, i32 0, i32 1
  %16 = load i32, i32* %15, align 4
  %17 = icmp ugt i32 %16, 1
  br i1 %17, label %18, label %27

18:                                               ; preds = %12
  %19 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %20 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %19, i32 0, i32 0
  %21 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %20, i32 0, i32 1
  %22 = load i32, i32* %21, align 4
  %23 = sub i32 %22, 1
  %24 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %25 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %24, i32 0, i32 0
  %26 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %25, i32 0, i32 1
  store i32 %23, i32* %26, align 4
  br label %39

27:                                               ; preds = %12
  %28 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %29 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %28, i32 0, i32 2
  %30 = load %struct.PVH*, %struct.PVH** %29, align 8
  %31 = icmp ne %struct.PVH* %30, null
  br i1 %31, label %32, label %36

32:                                               ; preds = %27
  %33 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %34 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %33, i32 0, i32 2
  %35 = load %struct.PVH*, %struct.PVH** %34, align 8
  call void @pn_free(%struct.PVH* %35)
  br label %36

36:                                               ; preds = %32, %27
  %37 = load %struct.PVHead*, %struct.PVHead** %2, align 8
  %38 = bitcast %struct.PVHead* %37 to i8*
  call void @free(i8* %38)
  br label %39

39:                                               ; preds = %36, %18, %11
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @pn_increment_children_ref(%struct.PVNode* %0) #0 {
  %2 = alloca %struct.PVNode*, align 8
  %3 = alloca %struct.PVH**, align 8
  %4 = alloca i8, align 1
  store %struct.PVNode* %0, %struct.PVNode** %2, align 8
  %5 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %6 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %5, i32 0, i32 2
  %7 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %6, i64 0, i64 0
  store %struct.PVH** %7, %struct.PVH*** %3, align 8
  store i8 0, i8* %4, align 1
  br label %8

8:                                                ; preds = %26, %1
  %9 = load i8, i8* %4, align 1
  %10 = zext i8 %9 to i32
  %11 = icmp slt i32 %10, 32
  br i1 %11, label %12, label %29

12:                                               ; preds = %8
  %13 = load %struct.PVH**, %struct.PVH*** %3, align 8
  %14 = load i8, i8* %4, align 1
  %15 = zext i8 %14 to i64
  %16 = getelementptr inbounds %struct.PVH*, %struct.PVH** %13, i64 %15
  %17 = load %struct.PVH*, %struct.PVH** %16, align 8
  %18 = icmp ne %struct.PVH* %17, null
  br i1 %18, label %19, label %25

19:                                               ; preds = %12
  %20 = load %struct.PVH**, %struct.PVH*** %3, align 8
  %21 = load i8, i8* %4, align 1
  %22 = zext i8 %21 to i64
  %23 = getelementptr inbounds %struct.PVH*, %struct.PVH** %20, i64 %22
  %24 = load %struct.PVH*, %struct.PVH** %23, align 8
  call void @pn_incr_ref(%struct.PVH* %24)
  br label %25

25:                                               ; preds = %19, %12
  br label %26

26:                                               ; preds = %25
  %27 = load i8, i8* %4, align 1
  %28 = add i8 %27, 1
  store i8 %28, i8* %4, align 1
  br label %8, !llvm.loop !7

29:                                               ; preds = %8
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVNode* @pn_copy(%struct.PVNode* %0) #0 {
  %2 = alloca %struct.PVNode*, align 8
  %3 = alloca %struct.PVNode*, align 8
  store %struct.PVNode* %0, %struct.PVNode** %2, align 8
  %4 = call i8* @heap_malloc(i32 280)
  %5 = bitcast i8* %4 to %struct.PVNode*
  store %struct.PVNode* %5, %struct.PVNode** %3, align 8
  %6 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %7 = bitcast %struct.PVNode* %6 to i8*
  %8 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %9 = bitcast %struct.PVNode* %8 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %7, i8* align 8 %9, i64 280, i1 false)
  %10 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %11 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %10, i32 0, i32 1
  %12 = load i32*, i32** %11, align 8
  %13 = icmp ne i32* %12, null
  br i1 %13, label %14, label %27

14:                                               ; preds = %1
  %15 = call i8* @heap_malloc(i32 128)
  %16 = bitcast i8* %15 to i32*
  %17 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %18 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %17, i32 0, i32 1
  store i32* %16, i32** %18, align 8
  %19 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %20 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %19, i32 0, i32 1
  %21 = load i32*, i32** %20, align 8
  %22 = bitcast i32* %21 to i8*
  %23 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %24 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %23, i32 0, i32 1
  %25 = load i32*, i32** %24, align 8
  %26 = bitcast i32* %25 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 4 %22, i8* align 4 %26, i64 128, i1 false)
  br label %27

27:                                               ; preds = %14, %1
  %28 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  call void @pn_increment_children_ref(%struct.PVNode* %28)
  %29 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  ret %struct.PVNode* %29
}

; Function Attrs: argmemonly nofree nounwind willreturn
declare void @llvm.memcpy.p0i8.p0i8.i64(i8* noalias nocapture writeonly, i8* noalias nocapture readonly, i64, i1 immarg) #4

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8* @pl_copy(%struct.PVH* %0, i32 %1) #0 {
  %3 = alloca i32, align 4
  %4 = alloca %struct.PVH*, align 8
  %5 = alloca %struct.PVH*, align 8
  store i32 %1, i32* %3, align 4
  store %struct.PVH* %0, %struct.PVH** %4, align 8
  %6 = load i32, i32* %3, align 4
  %7 = call i8* @heap_malloc(i32 %6)
  %8 = bitcast i8* %7 to %struct.PVH*
  store %struct.PVH* %8, %struct.PVH** %5, align 8
  %9 = load %struct.PVH*, %struct.PVH** %5, align 8
  %10 = bitcast %struct.PVH* %9 to i8*
  %11 = load %struct.PVH*, %struct.PVH** %4, align 8
  %12 = bitcast %struct.PVH* %11 to i8*
  %13 = load i32, i32* %3, align 4
  %14 = zext i32 %13 to i64
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 4 %10, i8* align 4 %12, i64 %14, i1 false)
  %15 = load %struct.PVH*, %struct.PVH** %5, align 8
  %16 = bitcast %struct.PVH* %15 to i8*
  ret i8* %16
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @pn_set_child(%struct.PVNode* %0, %struct.PVH* %1, i8 %2) #0 {
  %4 = alloca i8, align 1
  %5 = alloca %struct.PVH*, align 8
  %6 = alloca %struct.PVNode*, align 8
  %7 = alloca i32, align 4
  %8 = alloca i32, align 4
  %9 = alloca %struct.PVH*, align 8
  store i8 %2, i8* %4, align 1
  store %struct.PVH* %1, %struct.PVH** %5, align 8
  store %struct.PVNode* %0, %struct.PVNode** %6, align 8
  store i32 0, i32* %7, align 4
  store i32 0, i32* %8, align 4
  %10 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %11 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %10, i32 0, i32 2
  %12 = load i8, i8* %4, align 1
  %13 = zext i8 %12 to i64
  %14 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %11, i64 0, i64 %13
  %15 = load %struct.PVH*, %struct.PVH** %14, align 8
  store %struct.PVH* %15, %struct.PVH** %9, align 8
  %16 = load %struct.PVH*, %struct.PVH** %9, align 8
  %17 = icmp ne %struct.PVH* %16, null
  br i1 %17, label %18, label %23

18:                                               ; preds = %3
  %19 = load %struct.PVH*, %struct.PVH** %9, align 8
  %20 = getelementptr inbounds %struct.PVH, %struct.PVH* %19, i32 0, i32 2
  %21 = load i32, i32* %20, align 4
  store i32 %21, i32* %7, align 4
  %22 = load %struct.PVH*, %struct.PVH** %9, align 8
  call void @pn_free(%struct.PVH* %22)
  br label %23

23:                                               ; preds = %18, %3
  %24 = load %struct.PVH*, %struct.PVH** %5, align 8
  %25 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %26 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %25, i32 0, i32 2
  %27 = load i8, i8* %4, align 1
  %28 = zext i8 %27 to i64
  %29 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %26, i64 0, i64 %28
  store %struct.PVH* %24, %struct.PVH** %29, align 8
  %30 = load %struct.PVH*, %struct.PVH** %5, align 8
  %31 = icmp ne %struct.PVH* %30, null
  br i1 %31, label %32, label %37

32:                                               ; preds = %23
  %33 = load %struct.PVH*, %struct.PVH** %5, align 8
  %34 = getelementptr inbounds %struct.PVH, %struct.PVH* %33, i32 0, i32 2
  %35 = load i32, i32* %34, align 4
  store i32 %35, i32* %8, align 4
  %36 = load %struct.PVH*, %struct.PVH** %5, align 8
  call void @pn_incr_ref(%struct.PVH* %36)
  br label %37

37:                                               ; preds = %32, %23
  %38 = load i32, i32* %8, align 4
  %39 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %40 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %39, i32 0, i32 0
  %41 = getelementptr inbounds %struct.PVH, %struct.PVH* %40, i32 0, i32 2
  %42 = load i32, i32* %41, align 8
  %43 = add i32 %42, %38
  store i32 %43, i32* %41, align 8
  %44 = load i32, i32* %7, align 4
  %45 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %46 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %45, i32 0, i32 0
  %47 = getelementptr inbounds %struct.PVH, %struct.PVH* %46, i32 0, i32 2
  %48 = load i32, i32* %47, align 8
  %49 = sub i32 %48, %44
  store i32 %49, i32* %47, align 8
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8* @pv_get_leaf(%struct.PVHead* %0, i32* %1) #0 {
  %3 = alloca i32*, align 8
  %4 = alloca %struct.PVHead*, align 8
  %5 = alloca i32, align 4
  %6 = alloca i8, align 1
  %7 = alloca i8*, align 8
  %8 = alloca i32*, align 8
  %9 = alloca i8, align 1
  %10 = alloca i8, align 1
  %11 = alloca i8, align 1
  %12 = alloca i32, align 4
  %13 = alloca i32, align 4
  store i32* %1, i32** %3, align 8
  store %struct.PVHead* %0, %struct.PVHead** %4, align 8
  %14 = load i32*, i32** %3, align 8
  %15 = load i32, i32* %14, align 4
  store i32 %15, i32* %5, align 4
  %16 = load i32, i32* %5, align 4
  %17 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %18 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %17, i32 0, i32 1
  %19 = load i32, i32* %18, align 8
  %20 = icmp uge i32 %16, %19
  br i1 %20, label %21, label %28

21:                                               ; preds = %2
  %22 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %23 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %22, i32 0, i32 1
  %24 = load i32, i32* %23, align 8
  %25 = load i32, i32* %5, align 4
  %26 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %27 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %26, i8* getelementptr inbounds ([46 x i8], [46 x i8]* @"??_C@_0CO@JDLAEGEB@pvector?5index?5out?5of?5bounds?3?5got@", i64 0, i64 0), i32 %25, i32 %24)
  call void @exit(i32 1) #5
  unreachable

28:                                               ; preds = %2
  %29 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %30 = call i8 @pv_depth(%struct.PVHead* %29)
  store i8 %30, i8* %6, align 1
  %31 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %32 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %31, i32 0, i32 2
  %33 = load %struct.PVH*, %struct.PVH** %32, align 8
  %34 = bitcast %struct.PVH* %33 to i8*
  store i8* %34, i8** %7, align 8
  %35 = load i8*, i8** %7, align 8
  %36 = bitcast i8* %35 to %struct.PVNode*
  %37 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %36, i32 0, i32 1
  %38 = load i32*, i32** %37, align 8
  store i32* %38, i32** %8, align 8
  br label %39

39:                                               ; preds = %123, %28
  %40 = load i8, i8* %6, align 1
  %41 = zext i8 %40 to i32
  %42 = icmp ne i32 %41, 0
  br i1 %42, label %43, label %46

43:                                               ; preds = %39
  %44 = load i32*, i32** %8, align 8
  %45 = icmp ne i32* %44, null
  br label %46

46:                                               ; preds = %43, %39
  %47 = phi i1 [ false, %39 ], [ %45, %43 ]
  br i1 %47, label %48, label %138

48:                                               ; preds = %46
  store i8 31, i8* %9, align 1
  store i8 0, i8* %10, align 1
  br label %49

49:                                               ; preds = %62, %48
  %50 = load i8, i8* %9, align 1
  %51 = zext i8 %50 to i32
  %52 = icmp sgt i32 %51, 0
  br i1 %52, label %53, label %60

53:                                               ; preds = %49
  %54 = load i32*, i32** %8, align 8
  %55 = load i8, i8* %9, align 1
  %56 = zext i8 %55 to i64
  %57 = getelementptr inbounds i32, i32* %54, i64 %56
  %58 = load i32, i32* %57, align 4
  %59 = icmp eq i32 %58, 0
  br label %60

60:                                               ; preds = %53, %49
  %61 = phi i1 [ false, %49 ], [ %59, %53 ]
  br i1 %61, label %62, label %65

62:                                               ; preds = %60
  %63 = load i8, i8* %9, align 1
  %64 = add i8 %63, -1
  store i8 %64, i8* %9, align 1
  br label %49, !llvm.loop !8

65:                                               ; preds = %60
  br label %66

66:                                               ; preds = %108, %65
  %67 = load i8, i8* %9, align 1
  %68 = zext i8 %67 to i32
  %69 = load i8, i8* %10, align 1
  %70 = zext i8 %69 to i32
  %71 = icmp sgt i32 %68, %70
  br i1 %71, label %72, label %109

72:                                               ; preds = %66
  %73 = load i8, i8* %10, align 1
  %74 = zext i8 %73 to i32
  %75 = load i8, i8* %9, align 1
  %76 = zext i8 %75 to i32
  %77 = load i8, i8* %10, align 1
  %78 = zext i8 %77 to i32
  %79 = sub nsw i32 %76, %78
  %80 = ashr i32 %79, 1
  %81 = add nsw i32 %74, %80
  %82 = trunc i32 %81 to i8
  store i8 %82, i8* %11, align 1
  %83 = load i32*, i32** %8, align 8
  %84 = load i8, i8* %11, align 1
  %85 = zext i8 %84 to i64
  %86 = getelementptr inbounds i32, i32* %83, i64 %85
  %87 = load i32, i32* %86, align 4
  store i32 %87, i32* %12, align 4
  %88 = load i32, i32* %12, align 4
  %89 = load i32, i32* %5, align 4
  %90 = icmp eq i32 %88, %89
  br i1 %90, label %91, label %96

91:                                               ; preds = %72
  %92 = load i8, i8* %11, align 1
  %93 = zext i8 %92 to i32
  %94 = add nsw i32 %93, 1
  %95 = trunc i32 %94 to i8
  store i8 %95, i8* %9, align 1
  br label %109

96:                                               ; preds = %72
  %97 = load i32, i32* %12, align 4
  %98 = load i32, i32* %5, align 4
  %99 = icmp ult i32 %97, %98
  br i1 %99, label %100, label %105

100:                                              ; preds = %96
  %101 = load i8, i8* %11, align 1
  %102 = zext i8 %101 to i32
  %103 = add nsw i32 %102, 1
  %104 = trunc i32 %103 to i8
  store i8 %104, i8* %10, align 1
  br label %107

105:                                              ; preds = %96
  %106 = load i8, i8* %11, align 1
  store i8 %106, i8* %9, align 1
  br label %107

107:                                              ; preds = %105, %100
  br label %108

108:                                              ; preds = %107
  br label %66, !llvm.loop !9

109:                                              ; preds = %91, %66
  %110 = load i8, i8* %9, align 1
  %111 = zext i8 %110 to i32
  %112 = icmp sgt i32 %111, 0
  br i1 %112, label %113, label %123

113:                                              ; preds = %109
  %114 = load i32*, i32** %8, align 8
  %115 = load i8, i8* %9, align 1
  %116 = zext i8 %115 to i32
  %117 = sub nsw i32 %116, 1
  %118 = sext i32 %117 to i64
  %119 = getelementptr inbounds i32, i32* %114, i64 %118
  %120 = load i32, i32* %119, align 4
  %121 = load i32, i32* %5, align 4
  %122 = sub i32 %121, %120
  store i32 %122, i32* %5, align 4
  br label %123

123:                                              ; preds = %113, %109
  %124 = load i8*, i8** %7, align 8
  %125 = bitcast i8* %124 to %struct.PVNode*
  %126 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %125, i32 0, i32 2
  %127 = load i8, i8* %9, align 1
  %128 = zext i8 %127 to i64
  %129 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %126, i64 0, i64 %128
  %130 = load %struct.PVH*, %struct.PVH** %129, align 8
  %131 = bitcast %struct.PVH* %130 to i8*
  store i8* %131, i8** %7, align 8
  %132 = load i8*, i8** %7, align 8
  %133 = bitcast i8* %132 to %struct.PVNode*
  %134 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %133, i32 0, i32 1
  %135 = load i32*, i32** %134, align 8
  store i32* %135, i32** %8, align 8
  %136 = load i8, i8* %6, align 1
  %137 = add i8 %136, -1
  store i8 %137, i8* %6, align 1
  br label %39, !llvm.loop !10

138:                                              ; preds = %46
  br label %139

139:                                              ; preds = %142, %138
  %140 = load i8, i8* %6, align 1
  %141 = icmp ne i8 %140, 0
  br i1 %141, label %142, label %159

142:                                              ; preds = %139
  %143 = load i32, i32* %5, align 4
  %144 = load i8, i8* %6, align 1
  %145 = zext i8 %144 to i32
  %146 = mul nsw i32 %145, 5
  %147 = lshr i32 %143, %146
  %148 = and i32 %147, 31
  store i32 %148, i32* %13, align 4
  %149 = load i8, i8* %6, align 1
  %150 = add i8 %149, -1
  store i8 %150, i8* %6, align 1
  %151 = load i8*, i8** %7, align 8
  %152 = bitcast i8* %151 to %struct.PVNode*
  %153 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %152, i32 0, i32 2
  %154 = load i32, i32* %13, align 4
  %155 = zext i32 %154 to i64
  %156 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %153, i64 0, i64 %155
  %157 = load %struct.PVH*, %struct.PVH** %156, align 8
  %158 = bitcast %struct.PVH* %157 to i8*
  store i8* %158, i8** %7, align 8
  br label %139, !llvm.loop !11

159:                                              ; preds = %139
  %160 = load i32, i32* %5, align 4
  %161 = load i32*, i32** %3, align 8
  store i32 %160, i32* %161, align 4
  %162 = load i8*, i8** %7, align 8
  ret i8* %162
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i16 @pv_uint16_get(%struct.PVHead* %0, i32 %1) #0 {
  %3 = alloca i32, align 4
  %4 = alloca %struct.PVHead*, align 8
  %5 = alloca i32, align 4
  %6 = alloca %struct.PVLeaf_uint16*, align 8
  store i32 %1, i32* %3, align 4
  store %struct.PVHead* %0, %struct.PVHead** %4, align 8
  %7 = load i32, i32* %3, align 4
  store i32 %7, i32* %5, align 4
  %8 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %9 = call i8* @pv_get_leaf(%struct.PVHead* %8, i32* %5)
  %10 = bitcast i8* %9 to %struct.PVLeaf_uint16*
  store %struct.PVLeaf_uint16* %10, %struct.PVLeaf_uint16** %6, align 8
  %11 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %12 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %11, i32 0, i32 1
  %13 = load i32, i32* %5, align 4
  %14 = and i32 %13, 31
  %15 = zext i32 %14 to i64
  %16 = getelementptr inbounds [32 x i16], [32 x i16]* %12, i64 0, i64 %15
  %17 = load i16, i16* %16, align 2
  ret i16 %17
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8 @pn_right_child_index(%struct.PVNode* %0) #0 {
  %2 = alloca i8, align 1
  %3 = alloca %struct.PVNode*, align 8
  %4 = alloca i8, align 1
  %5 = alloca i8, align 1
  %6 = alloca i32, align 4
  %7 = alloca i8, align 1
  %8 = alloca i8, align 1
  %9 = alloca i32, align 4
  store %struct.PVNode* %0, %struct.PVNode** %3, align 8
  %10 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %11 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %10, i32 0, i32 1
  %12 = load i32*, i32** %11, align 8
  %13 = icmp ne i32* %12, null
  br i1 %13, label %14, label %33

14:                                               ; preds = %1
  store i8 32, i8* %4, align 1
  br label %15

15:                                               ; preds = %30, %14
  %16 = load i8, i8* %4, align 1
  %17 = sext i8 %16 to i32
  %18 = icmp sge i32 %17, 0
  br i1 %18, label %19, label %28

19:                                               ; preds = %15
  %20 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %21 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %20, i32 0, i32 2
  %22 = load i8, i8* %4, align 1
  %23 = add i8 %22, -1
  store i8 %23, i8* %4, align 1
  %24 = sext i8 %23 to i64
  %25 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %21, i64 0, i64 %24
  %26 = load %struct.PVH*, %struct.PVH** %25, align 8
  %27 = icmp eq %struct.PVH* %26, null
  br label %28

28:                                               ; preds = %19, %15
  %29 = phi i1 [ false, %15 ], [ %27, %19 ]
  br i1 %29, label %30, label %31

30:                                               ; preds = %28
  br label %15, !llvm.loop !12

31:                                               ; preds = %28
  %32 = load i8, i8* %4, align 1
  store i8 %32, i8* %2, align 1
  br label %66

33:                                               ; preds = %1
  %34 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %35 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %34, i32 0, i32 0
  %36 = getelementptr inbounds %struct.PVH, %struct.PVH* %35, i32 0, i32 0
  %37 = load i8, i8* %36, align 8
  store i8 %37, i8* %5, align 1
  store i32 1, i32* %6, align 4
  %38 = load i8, i8* %5, align 1
  store i8 %38, i8* %7, align 1
  br label %39

39:                                               ; preds = %46, %33
  %40 = load i8, i8* %7, align 1
  %41 = zext i8 %40 to i32
  %42 = icmp sgt i32 %41, 0
  br i1 %42, label %43, label %49

43:                                               ; preds = %39
  %44 = load i32, i32* %6, align 4
  %45 = shl i32 %44, 5
  store i32 %45, i32* %6, align 4
  br label %46

46:                                               ; preds = %43
  %47 = load i8, i8* %7, align 1
  %48 = add i8 %47, -1
  store i8 %48, i8* %7, align 1
  br label %39, !llvm.loop !13

49:                                               ; preds = %39
  store i8 0, i8* %8, align 1
  %50 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %51 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %50, i32 0, i32 0
  %52 = getelementptr inbounds %struct.PVH, %struct.PVH* %51, i32 0, i32 2
  %53 = load i32, i32* %52, align 8
  store i32 %53, i32* %9, align 4
  br label %54

54:                                               ; preds = %58, %49
  %55 = load i32, i32* %9, align 4
  %56 = load i32, i32* %6, align 4
  %57 = icmp ugt i32 %55, %56
  br i1 %57, label %58, label %64

58:                                               ; preds = %54
  %59 = load i32, i32* %6, align 4
  %60 = load i32, i32* %9, align 4
  %61 = sub i32 %60, %59
  store i32 %61, i32* %9, align 4
  %62 = load i8, i8* %8, align 1
  %63 = add i8 %62, 1
  store i8 %63, i8* %8, align 1
  br label %54, !llvm.loop !14

64:                                               ; preds = %54
  %65 = load i8, i8* %8, align 1
  store i8 %65, i8* %2, align 1
  br label %66

66:                                               ; preds = %64, %31
  %67 = load i8, i8* %2, align 1
  ret i8 %67
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8* @pn_right_child(%struct.PVNode* %0) #0 {
  %2 = alloca %struct.PVNode*, align 8
  store %struct.PVNode* %0, %struct.PVNode** %2, align 8
  %3 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %4 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %3, i32 0, i32 2
  %5 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %6 = call i8 @pn_right_child_index(%struct.PVNode* %5)
  %7 = zext i8 %6 to i64
  %8 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %4, i64 0, i64 %7
  %9 = load %struct.PVH*, %struct.PVH** %8, align 8
  %10 = bitcast %struct.PVH* %9 to i8*
  ret i8* %10
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @pn_update_index_table(%struct.PVNode* %0) #0 {
  %2 = alloca %struct.PVNode*, align 8
  %3 = alloca [32 x i32], align 16
  %4 = alloca i8, align 1
  %5 = alloca i8, align 1
  %6 = alloca i32, align 4
  %7 = alloca i32, align 4
  %8 = alloca i8, align 1
  %9 = alloca i32, align 4
  %10 = alloca i8, align 1
  store %struct.PVNode* %0, %struct.PVNode** %2, align 8
  %11 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %12 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %11, i32 0, i32 0
  %13 = getelementptr inbounds %struct.PVH, %struct.PVH* %12, i32 0, i32 0
  %14 = load i8, i8* %13, align 8
  store i8 %14, i8* %4, align 1
  store i8 0, i8* %5, align 1
  store i32 0, i32* %6, align 4
  store i32 0, i32* %7, align 4
  store i8 0, i8* %8, align 1
  br label %15

15:                                               ; preds = %85, %1
  %16 = load i8, i8* %8, align 1
  %17 = zext i8 %16 to i32
  %18 = icmp slt i32 %17, 32
  br i1 %18, label %19, label %88

19:                                               ; preds = %15
  %20 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %21 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %20, i32 0, i32 2
  %22 = load i8, i8* %8, align 1
  %23 = zext i8 %22 to i64
  %24 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %21, i64 0, i64 %23
  %25 = load %struct.PVH*, %struct.PVH** %24, align 8
  %26 = icmp ne %struct.PVH* %25, null
  br i1 %26, label %27, label %80

27:                                               ; preds = %19
  store i32 1, i32* %9, align 4
  %28 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %29 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %28, i32 0, i32 0
  %30 = getelementptr inbounds %struct.PVH, %struct.PVH* %29, i32 0, i32 0
  %31 = load i8, i8* %30, align 8
  store i8 %31, i8* %10, align 1
  br label %32

32:                                               ; preds = %39, %27
  %33 = load i8, i8* %10, align 1
  %34 = zext i8 %33 to i32
  %35 = icmp sgt i32 %34, 0
  br i1 %35, label %36, label %42

36:                                               ; preds = %32
  %37 = load i32, i32* %9, align 4
  %38 = shl i32 %37, 5
  store i32 %38, i32* %9, align 4
  br label %39

39:                                               ; preds = %36
  %40 = load i8, i8* %10, align 1
  %41 = add i8 %40, -1
  store i8 %41, i8* %10, align 1
  br label %32, !llvm.loop !15

42:                                               ; preds = %32
  %43 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %44 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %43, i32 0, i32 2
  %45 = load i8, i8* %8, align 1
  %46 = zext i8 %45 to i64
  %47 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %44, i64 0, i64 %46
  %48 = load %struct.PVH*, %struct.PVH** %47, align 8
  %49 = getelementptr inbounds %struct.PVH, %struct.PVH* %48, i32 0, i32 2
  %50 = load i32, i32* %49, align 4
  store i32 %50, i32* %7, align 4
  %51 = load i32, i32* %7, align 4
  %52 = icmp ugt i32 %51, 0
  br i1 %52, label %53, label %72

53:                                               ; preds = %42
  %54 = load i32, i32* %7, align 4
  %55 = load i32, i32* %9, align 4
  %56 = icmp ult i32 %54, %55
  br i1 %56, label %57, label %72

57:                                               ; preds = %53
  %58 = load i8, i8* %8, align 1
  %59 = zext i8 %58 to i32
  %60 = icmp slt i32 %59, 31
  br i1 %60, label %61, label %72

61:                                               ; preds = %57
  %62 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %63 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %62, i32 0, i32 2
  %64 = load i8, i8* %8, align 1
  %65 = zext i8 %64 to i32
  %66 = add nsw i32 %65, 1
  %67 = sext i32 %66 to i64
  %68 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %63, i64 0, i64 %67
  %69 = load %struct.PVH*, %struct.PVH** %68, align 8
  %70 = icmp ne %struct.PVH* %69, null
  br i1 %70, label %71, label %72

71:                                               ; preds = %61
  store i8 1, i8* %5, align 1
  br label %72

72:                                               ; preds = %71, %61, %57, %53, %42
  %73 = load i32, i32* %7, align 4
  %74 = load i32, i32* %6, align 4
  %75 = add i32 %74, %73
  store i32 %75, i32* %6, align 4
  %76 = load i32, i32* %6, align 4
  %77 = load i8, i8* %8, align 1
  %78 = zext i8 %77 to i64
  %79 = getelementptr inbounds [32 x i32], [32 x i32]* %3, i64 0, i64 %78
  store i32 %76, i32* %79, align 4
  br label %84

80:                                               ; preds = %19
  %81 = load i8, i8* %8, align 1
  %82 = zext i8 %81 to i64
  %83 = getelementptr inbounds [32 x i32], [32 x i32]* %3, i64 0, i64 %82
  store i32 0, i32* %83, align 4
  br label %84

84:                                               ; preds = %80, %72
  br label %85

85:                                               ; preds = %84
  %86 = load i8, i8* %8, align 1
  %87 = add i8 %86, 1
  store i8 %87, i8* %8, align 1
  br label %15, !llvm.loop !16

88:                                               ; preds = %15
  %89 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %90 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %89, i32 0, i32 1
  %91 = load i32*, i32** %90, align 8
  %92 = icmp ne i32* %91, null
  br i1 %92, label %93, label %100

93:                                               ; preds = %88
  %94 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %95 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %94, i32 0, i32 1
  %96 = load i32*, i32** %95, align 8
  %97 = bitcast i32* %96 to i8*
  call void @free(i8* %97)
  %98 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %99 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %98, i32 0, i32 1
  store i32* null, i32** %99, align 8
  br label %100

100:                                              ; preds = %93, %88
  %101 = load i8, i8* %5, align 1
  %102 = icmp ne i8 %101, 0
  br i1 %102, label %103, label %114

103:                                              ; preds = %100
  %104 = call i8* @heap_malloc(i32 128)
  %105 = bitcast i8* %104 to i32*
  %106 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %107 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %106, i32 0, i32 1
  store i32* %105, i32** %107, align 8
  %108 = load %struct.PVNode*, %struct.PVNode** %2, align 8
  %109 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %108, i32 0, i32 1
  %110 = load i32*, i32** %109, align 8
  %111 = bitcast i32* %110 to i8*
  %112 = getelementptr inbounds [32 x i32], [32 x i32]* %3, i64 0, i64 0
  %113 = bitcast i32* %112 to i8*
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 4 %111, i8* align 16 %113, i64 128, i1 false)
  br label %114

114:                                              ; preds = %103, %100
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8 @pn_branch_count(%struct.PVH* %0) #0 {
  %2 = alloca i8, align 1
  %3 = alloca %struct.PVH*, align 8
  %4 = alloca i8, align 1
  store %struct.PVH* %0, %struct.PVH** %3, align 8
  %5 = load %struct.PVH*, %struct.PVH** %3, align 8
  %6 = getelementptr inbounds %struct.PVH, %struct.PVH* %5, i32 0, i32 0
  %7 = load i8, i8* %6, align 4
  %8 = zext i8 %7 to i32
  %9 = icmp sgt i32 %8, 0
  br i1 %9, label %10, label %32

10:                                               ; preds = %1
  store i8 0, i8* %4, align 1
  br label %11

11:                                               ; preds = %27, %10
  %12 = load i8, i8* %4, align 1
  %13 = zext i8 %12 to i32
  %14 = icmp slt i32 %13, 32
  br i1 %14, label %15, label %24

15:                                               ; preds = %11
  %16 = load %struct.PVH*, %struct.PVH** %3, align 8
  %17 = bitcast %struct.PVH* %16 to %struct.PVNode*
  %18 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %17, i32 0, i32 2
  %19 = load i8, i8* %4, align 1
  %20 = zext i8 %19 to i64
  %21 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %18, i64 0, i64 %20
  %22 = load %struct.PVH*, %struct.PVH** %21, align 8
  %23 = icmp ne %struct.PVH* %22, null
  br label %24

24:                                               ; preds = %15, %11
  %25 = phi i1 [ false, %11 ], [ %23, %15 ]
  br i1 %25, label %26, label %30

26:                                               ; preds = %24
  br label %27

27:                                               ; preds = %26
  %28 = load i8, i8* %4, align 1
  %29 = add i8 %28, 1
  store i8 %29, i8* %4, align 1
  br label %11, !llvm.loop !17

30:                                               ; preds = %24
  %31 = load i8, i8* %4, align 1
  store i8 %31, i8* %2, align 1
  br label %37

32:                                               ; preds = %1
  %33 = load %struct.PVH*, %struct.PVH** %3, align 8
  %34 = getelementptr inbounds %struct.PVH, %struct.PVH* %33, i32 0, i32 2
  %35 = load i32, i32* %34, align 4
  %36 = trunc i32 %35 to i8
  store i8 %36, i8* %2, align 1
  br label %37

37:                                               ; preds = %32, %30
  %38 = load i8, i8* %2, align 1
  ret i8 %38
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @pn_branch_sum(%struct.PVNode* %0) #0 {
  %2 = alloca i32, align 4
  %3 = alloca %struct.PVNode*, align 8
  %4 = alloca i32, align 4
  %5 = alloca i8, align 1
  %6 = alloca i8, align 1
  store %struct.PVNode* %0, %struct.PVNode** %3, align 8
  store i32 0, i32* %4, align 4
  %7 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %8 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %7, i32 0, i32 0
  %9 = getelementptr inbounds %struct.PVH, %struct.PVH* %8, i32 0, i32 0
  %10 = load i8, i8* %9, align 8
  store i8 %10, i8* %5, align 1
  %11 = load i8, i8* %5, align 1
  %12 = zext i8 %11 to i32
  %13 = icmp eq i32 %12, 0
  br i1 %13, label %14, label %19

14:                                               ; preds = %1
  %15 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %16 = bitcast %struct.PVNode* %15 to %struct.PVH*
  %17 = getelementptr inbounds %struct.PVH, %struct.PVH* %16, i32 0, i32 2
  %18 = load i32, i32* %17, align 4
  store i32 %18, i32* %2, align 4
  br label %49

19:                                               ; preds = %1
  store i8 0, i8* %6, align 1
  br label %20

20:                                               ; preds = %44, %19
  %21 = load i8, i8* %6, align 1
  %22 = zext i8 %21 to i32
  %23 = icmp slt i32 %22, 32
  br i1 %23, label %24, label %47

24:                                               ; preds = %20
  %25 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %26 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %25, i32 0, i32 2
  %27 = load i8, i8* %6, align 1
  %28 = zext i8 %27 to i64
  %29 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %26, i64 0, i64 %28
  %30 = load %struct.PVH*, %struct.PVH** %29, align 8
  %31 = icmp ne %struct.PVH* %30, null
  br i1 %31, label %33, label %32

32:                                               ; preds = %24
  br label %47

33:                                               ; preds = %24
  %34 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  %35 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %34, i32 0, i32 2
  %36 = load i8, i8* %6, align 1
  %37 = zext i8 %36 to i64
  %38 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %35, i64 0, i64 %37
  %39 = load %struct.PVH*, %struct.PVH** %38, align 8
  %40 = call i8 @pn_branch_count(%struct.PVH* %39)
  %41 = zext i8 %40 to i32
  %42 = load i32, i32* %4, align 4
  %43 = add i32 %42, %41
  store i32 %43, i32* %4, align 4
  br label %44

44:                                               ; preds = %33
  %45 = load i8, i8* %6, align 1
  %46 = add i8 %45, 1
  store i8 %46, i8* %6, align 1
  br label %20, !llvm.loop !18

47:                                               ; preds = %32, %20
  %48 = load i32, i32* %4, align 4
  store i32 %48, i32* %2, align 4
  br label %49

49:                                               ; preds = %47, %14
  %50 = load i32, i32* %2, align 4
  ret i32 %50
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8 @pn_needs_rebalancing(%struct.PVNode* %0, %struct.PVNode* %1) #0 {
  %3 = alloca i8, align 1
  %4 = alloca %struct.PVNode*, align 8
  %5 = alloca %struct.PVNode*, align 8
  %6 = alloca i32, align 4
  %7 = alloca i32, align 4
  %8 = alloca i32, align 4
  store %struct.PVNode* %1, %struct.PVNode** %4, align 8
  store %struct.PVNode* %0, %struct.PVNode** %5, align 8
  %9 = load %struct.PVNode*, %struct.PVNode** %5, align 8
  %10 = icmp eq %struct.PVNode* %9, null
  br i1 %10, label %14, label %11

11:                                               ; preds = %2
  %12 = load %struct.PVNode*, %struct.PVNode** %4, align 8
  %13 = icmp eq %struct.PVNode* %12, null
  br i1 %13, label %14, label %15

14:                                               ; preds = %11, %2
  store i8 0, i8* %3, align 1
  br label %40

15:                                               ; preds = %11
  %16 = load %struct.PVNode*, %struct.PVNode** %5, align 8
  %17 = call i32 @pn_branch_sum(%struct.PVNode* %16)
  %18 = load %struct.PVNode*, %struct.PVNode** %4, align 8
  %19 = call i32 @pn_branch_sum(%struct.PVNode* %18)
  %20 = add i32 %17, %19
  store i32 %20, i32* %6, align 4
  %21 = load %struct.PVNode*, %struct.PVNode** %5, align 8
  %22 = bitcast %struct.PVNode* %21 to %struct.PVH*
  %23 = call i8 @pn_branch_count(%struct.PVH* %22)
  %24 = zext i8 %23 to i32
  %25 = load %struct.PVNode*, %struct.PVNode** %4, align 8
  %26 = bitcast %struct.PVNode* %25 to %struct.PVH*
  %27 = call i8 @pn_branch_count(%struct.PVH* %26)
  %28 = zext i8 %27 to i32
  %29 = add nsw i32 %24, %28
  store i32 %29, i32* %7, align 4
  %30 = load i32, i32* %7, align 4
  %31 = load i32, i32* %6, align 4
  %32 = sub i32 %31, 1
  %33 = lshr i32 %32, 5
  %34 = sub i32 %30, %33
  %35 = sub i32 %34, 1
  store i32 %35, i32* %8, align 4
  %36 = load i32, i32* %8, align 4
  %37 = icmp sgt i32 %36, 1
  br i1 %37, label %38, label %39

38:                                               ; preds = %15
  store i8 1, i8* %3, align 1
  br label %40

39:                                               ; preds = %15
  store i8 0, i8* %3, align 1
  br label %40

40:                                               ; preds = %39, %38, %14
  %41 = load i8, i8* %3, align 1
  ret i8 %41
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVLeaf_uint16* @pl_16_concatenate(%struct.PVLeaf_uint16* %0, %struct.PVLeaf_uint16* %1, %struct.PVLeaf_uint16** %2) #0 {
  %4 = alloca %struct.PVLeaf_uint16*, align 8
  %5 = alloca %struct.PVLeaf_uint16**, align 8
  %6 = alloca %struct.PVLeaf_uint16*, align 8
  %7 = alloca %struct.PVLeaf_uint16*, align 8
  %8 = alloca %struct.PVLeaf_uint16*, align 8
  %9 = alloca i32, align 4
  %10 = alloca %struct.PVLeaf_uint16*, align 8
  store %struct.PVLeaf_uint16** %2, %struct.PVLeaf_uint16*** %5, align 8
  store %struct.PVLeaf_uint16* %1, %struct.PVLeaf_uint16** %6, align 8
  store %struct.PVLeaf_uint16* %0, %struct.PVLeaf_uint16** %7, align 8
  %11 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %12 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %11, i32 0, i32 0
  %13 = getelementptr inbounds %struct.PVH, %struct.PVH* %12, i32 0, i32 2
  %14 = load i32, i32* %13, align 4
  %15 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %16 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %15, i32 0, i32 0
  %17 = getelementptr inbounds %struct.PVH, %struct.PVH* %16, i32 0, i32 2
  %18 = load i32, i32* %17, align 4
  %19 = add i32 %14, %18
  %20 = icmp ule i32 %19, 32
  br i1 %20, label %21, label %76

21:                                               ; preds = %3
  %22 = load %struct.PVLeaf_uint16**, %struct.PVLeaf_uint16*** %5, align 8
  %23 = icmp ne %struct.PVLeaf_uint16** %22, null
  br i1 %23, label %24, label %26

24:                                               ; preds = %21
  %25 = load %struct.PVLeaf_uint16**, %struct.PVLeaf_uint16*** %5, align 8
  store %struct.PVLeaf_uint16* null, %struct.PVLeaf_uint16** %25, align 8
  br label %26

26:                                               ; preds = %24, %21
  %27 = call %struct.PVH* @pl_new(i32 76)
  %28 = bitcast %struct.PVH* %27 to %struct.PVLeaf_uint16*
  store %struct.PVLeaf_uint16* %28, %struct.PVLeaf_uint16** %8, align 8
  %29 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %8, align 8
  %30 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %29, i32 0, i32 1
  %31 = getelementptr inbounds [32 x i16], [32 x i16]* %30, i64 0, i64 0
  %32 = bitcast i16* %31 to i8*
  %33 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %34 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %33, i32 0, i32 1
  %35 = getelementptr inbounds [32 x i16], [32 x i16]* %34, i64 0, i64 0
  %36 = bitcast i16* %35 to i8*
  %37 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %38 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %37, i32 0, i32 0
  %39 = getelementptr inbounds %struct.PVH, %struct.PVH* %38, i32 0, i32 2
  %40 = load i32, i32* %39, align 4
  %41 = zext i32 %40 to i64
  %42 = mul i64 %41, 2
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 4 %32, i8* align 4 %36, i64 %42, i1 false)
  %43 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %8, align 8
  %44 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %43, i32 0, i32 1
  %45 = getelementptr inbounds [32 x i16], [32 x i16]* %44, i64 0, i64 0
  %46 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %47 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %46, i32 0, i32 0
  %48 = getelementptr inbounds %struct.PVH, %struct.PVH* %47, i32 0, i32 2
  %49 = load i32, i32* %48, align 4
  %50 = zext i32 %49 to i64
  %51 = getelementptr inbounds i16, i16* %45, i64 %50
  %52 = bitcast i16* %51 to i8*
  %53 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %54 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %53, i32 0, i32 1
  %55 = getelementptr inbounds [32 x i16], [32 x i16]* %54, i64 0, i64 0
  %56 = bitcast i16* %55 to i8*
  %57 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %58 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %57, i32 0, i32 0
  %59 = getelementptr inbounds %struct.PVH, %struct.PVH* %58, i32 0, i32 2
  %60 = load i32, i32* %59, align 4
  %61 = zext i32 %60 to i64
  %62 = mul i64 %61, 2
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 2 %52, i8* align 4 %56, i64 %62, i1 false)
  %63 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %64 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %63, i32 0, i32 0
  %65 = getelementptr inbounds %struct.PVH, %struct.PVH* %64, i32 0, i32 2
  %66 = load i32, i32* %65, align 4
  %67 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %68 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %67, i32 0, i32 0
  %69 = getelementptr inbounds %struct.PVH, %struct.PVH* %68, i32 0, i32 2
  %70 = load i32, i32* %69, align 4
  %71 = add i32 %66, %70
  %72 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %8, align 8
  %73 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %72, i32 0, i32 0
  %74 = getelementptr inbounds %struct.PVH, %struct.PVH* %73, i32 0, i32 2
  store i32 %71, i32* %74, align 4
  %75 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %8, align 8
  store %struct.PVLeaf_uint16* %75, %struct.PVLeaf_uint16** %4, align 8
  br label %161

76:                                               ; preds = %3
  %77 = load %struct.PVLeaf_uint16**, %struct.PVLeaf_uint16*** %5, align 8
  %78 = icmp ne %struct.PVLeaf_uint16** %77, null
  br i1 %78, label %82, label %79

79:                                               ; preds = %76
  %80 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %81 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %80, i8* getelementptr inbounds ([19 x i8], [19 x i8]* @"??_C@_0BD@DJOMGDD@overflow?5required?6?$AA@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

82:                                               ; preds = %76
  %83 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %84 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %83, i32 0, i32 0
  %85 = getelementptr inbounds %struct.PVH, %struct.PVH* %84, i32 0, i32 2
  %86 = load i32, i32* %85, align 4
  %87 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %88 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %87, i32 0, i32 0
  %89 = getelementptr inbounds %struct.PVH, %struct.PVH* %88, i32 0, i32 2
  %90 = load i32, i32* %89, align 4
  %91 = add i32 %86, %90
  %92 = sub i32 %91, 32
  store i32 %92, i32* %9, align 4
  %93 = call %struct.PVH* @pl_new(i32 76)
  %94 = bitcast %struct.PVH* %93 to %struct.PVLeaf_uint16*
  %95 = load %struct.PVLeaf_uint16**, %struct.PVLeaf_uint16*** %5, align 8
  store %struct.PVLeaf_uint16* %94, %struct.PVLeaf_uint16** %95, align 8
  %96 = call %struct.PVH* @pl_new(i32 76)
  %97 = bitcast %struct.PVH* %96 to %struct.PVLeaf_uint16*
  store %struct.PVLeaf_uint16* %97, %struct.PVLeaf_uint16** %10, align 8
  %98 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %10, align 8
  %99 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %98, i32 0, i32 1
  %100 = getelementptr inbounds [32 x i16], [32 x i16]* %99, i64 0, i64 0
  %101 = bitcast i16* %100 to i8*
  %102 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %103 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %102, i32 0, i32 1
  %104 = getelementptr inbounds [32 x i16], [32 x i16]* %103, i64 0, i64 0
  %105 = bitcast i16* %104 to i8*
  %106 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %107 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %106, i32 0, i32 0
  %108 = getelementptr inbounds %struct.PVH, %struct.PVH* %107, i32 0, i32 2
  %109 = load i32, i32* %108, align 4
  %110 = zext i32 %109 to i64
  %111 = mul i64 %110, 2
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 4 %101, i8* align 4 %105, i64 %111, i1 false)
  %112 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %10, align 8
  %113 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %112, i32 0, i32 1
  %114 = getelementptr inbounds [32 x i16], [32 x i16]* %113, i64 0, i64 0
  %115 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %116 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %115, i32 0, i32 0
  %117 = getelementptr inbounds %struct.PVH, %struct.PVH* %116, i32 0, i32 2
  %118 = load i32, i32* %117, align 4
  %119 = zext i32 %118 to i64
  %120 = getelementptr inbounds i16, i16* %114, i64 %119
  %121 = bitcast i16* %120 to i8*
  %122 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %123 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %122, i32 0, i32 1
  %124 = getelementptr inbounds [32 x i16], [32 x i16]* %123, i64 0, i64 0
  %125 = bitcast i16* %124 to i8*
  %126 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %127 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %126, i32 0, i32 0
  %128 = getelementptr inbounds %struct.PVH, %struct.PVH* %127, i32 0, i32 2
  %129 = load i32, i32* %128, align 4
  %130 = sub i32 32, %129
  %131 = zext i32 %130 to i64
  %132 = mul i64 %131, 2
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 2 %121, i8* align 4 %125, i64 %132, i1 false)
  %133 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %10, align 8
  %134 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %133, i32 0, i32 0
  %135 = getelementptr inbounds %struct.PVH, %struct.PVH* %134, i32 0, i32 2
  store i32 32, i32* %135, align 4
  %136 = load i32, i32* %9, align 4
  %137 = load %struct.PVLeaf_uint16**, %struct.PVLeaf_uint16*** %5, align 8
  %138 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %137, align 8
  %139 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %138, i32 0, i32 0
  %140 = getelementptr inbounds %struct.PVH, %struct.PVH* %139, i32 0, i32 2
  store i32 %136, i32* %140, align 4
  %141 = load %struct.PVLeaf_uint16**, %struct.PVLeaf_uint16*** %5, align 8
  %142 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %141, align 8
  %143 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %142, i32 0, i32 1
  %144 = getelementptr inbounds [32 x i16], [32 x i16]* %143, i64 0, i64 0
  %145 = bitcast i16* %144 to i8*
  %146 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %6, align 8
  %147 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %146, i32 0, i32 1
  %148 = getelementptr inbounds [32 x i16], [32 x i16]* %147, i64 0, i64 0
  %149 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %7, align 8
  %150 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %149, i32 0, i32 0
  %151 = getelementptr inbounds %struct.PVH, %struct.PVH* %150, i32 0, i32 2
  %152 = load i32, i32* %151, align 4
  %153 = sub i32 32, %152
  %154 = zext i32 %153 to i64
  %155 = getelementptr inbounds i16, i16* %148, i64 %154
  %156 = bitcast i16* %155 to i8*
  %157 = load i32, i32* %9, align 4
  %158 = zext i32 %157 to i64
  %159 = mul i64 %158, 2
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 4 %145, i8* align 2 %156, i64 %159, i1 false)
  %160 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %10, align 8
  store %struct.PVLeaf_uint16* %160, %struct.PVLeaf_uint16** %4, align 8
  br label %161

161:                                              ; preds = %82, %26
  %162 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %4, align 8
  ret %struct.PVLeaf_uint16* %162
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVH* @pn_join_nodes(%struct.PVH* %0, %struct.PVH* %1, %struct.PVH** %2) #0 {
  %4 = alloca %struct.PVH*, align 8
  %5 = alloca %struct.PVH**, align 8
  %6 = alloca %struct.PVH*, align 8
  %7 = alloca %struct.PVH*, align 8
  %8 = alloca %struct.PVNode*, align 8
  %9 = alloca %struct.PVNode*, align 8
  %10 = alloca i32, align 4
  %11 = alloca i32, align 4
  %12 = alloca i8, align 1
  %13 = alloca %struct.PVNode*, align 8
  %14 = alloca i32, align 4
  %15 = alloca %struct.PVNode*, align 8
  %16 = alloca %struct.PVNode*, align 8
  %17 = alloca i8, align 1
  %18 = alloca i8, align 1
  %19 = alloca i8, align 1
  store %struct.PVH** %2, %struct.PVH*** %5, align 8
  store %struct.PVH* %1, %struct.PVH** %6, align 8
  store %struct.PVH* %0, %struct.PVH** %7, align 8
  %20 = load %struct.PVH*, %struct.PVH** %7, align 8
  %21 = bitcast %struct.PVH* %20 to %struct.PVLeaf_uint16*
  %22 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %21, i32 0, i32 0
  %23 = getelementptr inbounds %struct.PVH, %struct.PVH* %22, i32 0, i32 0
  %24 = load i8, i8* %23, align 4
  %25 = zext i8 %24 to i32
  %26 = icmp eq i32 %25, 0
  br i1 %26, label %27, label %44

27:                                               ; preds = %3
  %28 = load %struct.PVH*, %struct.PVH** %6, align 8
  %29 = bitcast %struct.PVH* %28 to %struct.PVLeaf_uint16*
  %30 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %29, i32 0, i32 0
  %31 = getelementptr inbounds %struct.PVH, %struct.PVH* %30, i32 0, i32 0
  %32 = load i8, i8* %31, align 4
  %33 = zext i8 %32 to i32
  %34 = icmp eq i32 %33, 0
  br i1 %34, label %35, label %44

35:                                               ; preds = %27
  %36 = load %struct.PVH**, %struct.PVH*** %5, align 8
  %37 = bitcast %struct.PVH** %36 to %struct.PVLeaf_uint16**
  %38 = load %struct.PVH*, %struct.PVH** %6, align 8
  %39 = bitcast %struct.PVH* %38 to %struct.PVLeaf_uint16*
  %40 = load %struct.PVH*, %struct.PVH** %7, align 8
  %41 = bitcast %struct.PVH* %40 to %struct.PVLeaf_uint16*
  %42 = call %struct.PVLeaf_uint16* @pl_16_concatenate(%struct.PVLeaf_uint16* %41, %struct.PVLeaf_uint16* %39, %struct.PVLeaf_uint16** %37)
  %43 = bitcast %struct.PVLeaf_uint16* %42 to %struct.PVH*
  store %struct.PVH* %43, %struct.PVH** %4, align 8
  br label %209

44:                                               ; preds = %27, %3
  %45 = load %struct.PVH*, %struct.PVH** %7, align 8
  %46 = bitcast %struct.PVH* %45 to %struct.PVNode*
  store %struct.PVNode* %46, %struct.PVNode** %8, align 8
  %47 = load %struct.PVH*, %struct.PVH** %6, align 8
  %48 = bitcast %struct.PVH* %47 to %struct.PVNode*
  store %struct.PVNode* %48, %struct.PVNode** %9, align 8
  %49 = load %struct.PVH*, %struct.PVH** %7, align 8
  %50 = call i8 @pn_branch_count(%struct.PVH* %49)
  %51 = zext i8 %50 to i32
  store i32 %51, i32* %10, align 4
  %52 = load %struct.PVH*, %struct.PVH** %6, align 8
  %53 = call i8 @pn_branch_count(%struct.PVH* %52)
  %54 = zext i8 %53 to i32
  store i32 %54, i32* %11, align 4
  %55 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %56 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %55, i32 0, i32 0
  %57 = getelementptr inbounds %struct.PVH, %struct.PVH* %56, i32 0, i32 0
  %58 = load i8, i8* %57, align 8
  store i8 %58, i8* %12, align 1
  %59 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  %60 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %59, i32 0, i32 0
  %61 = getelementptr inbounds %struct.PVH, %struct.PVH* %60, i32 0, i32 0
  %62 = load i8, i8* %61, align 8
  %63 = zext i8 %62 to i32
  %64 = load i8, i8* %12, align 1
  %65 = zext i8 %64 to i32
  %66 = icmp ne i32 %63, %65
  br i1 %66, label %67, label %70

67:                                               ; preds = %44
  %68 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %69 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %68, i8* getelementptr inbounds ([28 x i8], [28 x i8]* @"??_C@_0BM@CICGOBLB@join?5error?0?5depth?5mismatch?6?$AA@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

70:                                               ; preds = %44
  %71 = load i32, i32* %10, align 4
  %72 = load i32, i32* %11, align 4
  %73 = add i32 %71, %72
  %74 = icmp ule i32 %73, 32
  br i1 %74, label %75, label %124

75:                                               ; preds = %70
  %76 = load %struct.PVH**, %struct.PVH*** %5, align 8
  %77 = icmp ne %struct.PVH** %76, null
  br i1 %77, label %78, label %80

78:                                               ; preds = %75
  %79 = load %struct.PVH**, %struct.PVH*** %5, align 8
  store %struct.PVH* null, %struct.PVH** %79, align 8
  br label %80

80:                                               ; preds = %78, %75
  %81 = load i8, i8* %12, align 1
  %82 = call %struct.PVNode* @pn_new(i8 %81)
  store %struct.PVNode* %82, %struct.PVNode** %13, align 8
  %83 = load %struct.PVNode*, %struct.PVNode** %13, align 8
  %84 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %83, i32 0, i32 2
  %85 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %84, i64 0, i64 0
  %86 = bitcast %struct.PVH** %85 to i8*
  %87 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %88 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %87, i32 0, i32 2
  %89 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %88, i64 0, i64 0
  %90 = bitcast %struct.PVH** %89 to i8*
  %91 = load i32, i32* %10, align 4
  %92 = zext i32 %91 to i64
  %93 = mul i64 %92, 8
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %86, i8* align 8 %90, i64 %93, i1 false)
  %94 = load %struct.PVNode*, %struct.PVNode** %13, align 8
  %95 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %94, i32 0, i32 2
  %96 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %95, i64 0, i64 0
  %97 = load i32, i32* %10, align 4
  %98 = zext i32 %97 to i64
  %99 = getelementptr inbounds %struct.PVH*, %struct.PVH** %96, i64 %98
  %100 = bitcast %struct.PVH** %99 to i8*
  %101 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  %102 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %101, i32 0, i32 2
  %103 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %102, i64 0, i64 0
  %104 = bitcast %struct.PVH** %103 to i8*
  %105 = load i32, i32* %11, align 4
  %106 = zext i32 %105 to i64
  %107 = mul i64 %106, 8
  call void @llvm.memcpy.p0i8.p0i8.i64(i8* align 8 %100, i8* align 8 %104, i64 %107, i1 false)
  %108 = load %struct.PVNode*, %struct.PVNode** %13, align 8
  call void @pn_update_index_table(%struct.PVNode* %108)
  %109 = load %struct.PVNode*, %struct.PVNode** %13, align 8
  call void @pn_increment_children_ref(%struct.PVNode* %109)
  %110 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %111 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %110, i32 0, i32 0
  %112 = getelementptr inbounds %struct.PVH, %struct.PVH* %111, i32 0, i32 2
  %113 = load i32, i32* %112, align 8
  %114 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  %115 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %114, i32 0, i32 0
  %116 = getelementptr inbounds %struct.PVH, %struct.PVH* %115, i32 0, i32 2
  %117 = load i32, i32* %116, align 8
  %118 = add i32 %113, %117
  %119 = load %struct.PVNode*, %struct.PVNode** %13, align 8
  %120 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %119, i32 0, i32 0
  %121 = getelementptr inbounds %struct.PVH, %struct.PVH* %120, i32 0, i32 2
  store i32 %118, i32* %121, align 8
  %122 = load %struct.PVNode*, %struct.PVNode** %13, align 8
  %123 = bitcast %struct.PVNode* %122 to %struct.PVH*
  store %struct.PVH* %123, %struct.PVH** %4, align 8
  br label %209

124:                                              ; preds = %70
  %125 = load %struct.PVH**, %struct.PVH*** %5, align 8
  %126 = icmp ne %struct.PVH** %125, null
  br i1 %126, label %130, label %127

127:                                              ; preds = %124
  %128 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %129 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %128, i8* getelementptr inbounds ([19 x i8], [19 x i8]* @"??_C@_0BD@DJOMGDD@overflow?5required?6?$AA@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

130:                                              ; preds = %124
  %131 = load i32, i32* %10, align 4
  %132 = load i32, i32* %11, align 4
  %133 = add i32 %131, %132
  %134 = sub i32 %133, 32
  store i32 %134, i32* %14, align 4
  %135 = load i8, i8* %12, align 1
  %136 = call %struct.PVNode* @pn_new(i8 %135)
  store %struct.PVNode* %136, %struct.PVNode** %15, align 8
  %137 = load i8, i8* %12, align 1
  %138 = call %struct.PVNode* @pn_new(i8 %137)
  store %struct.PVNode* %138, %struct.PVNode** %16, align 8
  store i8 0, i8* %17, align 1
  br label %139

139:                                              ; preds = %153, %130
  %140 = load i8, i8* %17, align 1
  %141 = zext i8 %140 to i32
  %142 = load i32, i32* %10, align 4
  %143 = icmp ult i32 %141, %142
  br i1 %143, label %144, label %156

144:                                              ; preds = %139
  %145 = load i8, i8* %17, align 1
  %146 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %147 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %146, i32 0, i32 2
  %148 = load i8, i8* %17, align 1
  %149 = zext i8 %148 to i64
  %150 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %147, i64 0, i64 %149
  %151 = load %struct.PVH*, %struct.PVH** %150, align 8
  %152 = load %struct.PVNode*, %struct.PVNode** %16, align 8
  call void @pn_set_child(%struct.PVNode* %152, %struct.PVH* %151, i8 %145)
  br label %153

153:                                              ; preds = %144
  %154 = load i8, i8* %17, align 1
  %155 = add i8 %154, 1
  store i8 %155, i8* %17, align 1
  br label %139, !llvm.loop !19

156:                                              ; preds = %139
  store i8 0, i8* %18, align 1
  br label %157

157:                                              ; preds = %176, %156
  %158 = load i8, i8* %18, align 1
  %159 = zext i8 %158 to i32
  %160 = load i32, i32* %10, align 4
  %161 = sub i32 32, %160
  %162 = icmp ult i32 %159, %161
  br i1 %162, label %163, label %179

163:                                              ; preds = %157
  %164 = load i8, i8* %18, align 1
  %165 = zext i8 %164 to i32
  %166 = load i32, i32* %10, align 4
  %167 = add i32 %165, %166
  %168 = trunc i32 %167 to i8
  %169 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  %170 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %169, i32 0, i32 2
  %171 = load i8, i8* %18, align 1
  %172 = zext i8 %171 to i64
  %173 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %170, i64 0, i64 %172
  %174 = load %struct.PVH*, %struct.PVH** %173, align 8
  %175 = load %struct.PVNode*, %struct.PVNode** %16, align 8
  call void @pn_set_child(%struct.PVNode* %175, %struct.PVH* %174, i8 %168)
  br label %176

176:                                              ; preds = %163
  %177 = load i8, i8* %18, align 1
  %178 = add i8 %177, 1
  store i8 %178, i8* %18, align 1
  br label %157, !llvm.loop !20

179:                                              ; preds = %157
  store i8 0, i8* %19, align 1
  br label %180

180:                                              ; preds = %198, %179
  %181 = load i8, i8* %19, align 1
  %182 = zext i8 %181 to i32
  %183 = load i32, i32* %14, align 4
  %184 = icmp ult i32 %182, %183
  br i1 %184, label %185, label %201

185:                                              ; preds = %180
  %186 = load i8, i8* %19, align 1
  %187 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  %188 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %187, i32 0, i32 2
  %189 = load i8, i8* %19, align 1
  %190 = zext i8 %189 to i32
  %191 = load i32, i32* %10, align 4
  %192 = sub i32 32, %191
  %193 = add i32 %190, %192
  %194 = zext i32 %193 to i64
  %195 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %188, i64 0, i64 %194
  %196 = load %struct.PVH*, %struct.PVH** %195, align 8
  %197 = load %struct.PVNode*, %struct.PVNode** %15, align 8
  call void @pn_set_child(%struct.PVNode* %197, %struct.PVH* %196, i8 %186)
  br label %198

198:                                              ; preds = %185
  %199 = load i8, i8* %19, align 1
  %200 = add i8 %199, 1
  store i8 %200, i8* %19, align 1
  br label %180, !llvm.loop !21

201:                                              ; preds = %180
  %202 = load %struct.PVNode*, %struct.PVNode** %16, align 8
  call void @pn_update_index_table(%struct.PVNode* %202)
  %203 = load %struct.PVNode*, %struct.PVNode** %15, align 8
  call void @pn_update_index_table(%struct.PVNode* %203)
  %204 = load %struct.PVNode*, %struct.PVNode** %15, align 8
  %205 = bitcast %struct.PVNode* %204 to %struct.PVH*
  %206 = load %struct.PVH**, %struct.PVH*** %5, align 8
  store %struct.PVH* %205, %struct.PVH** %206, align 8
  %207 = load %struct.PVNode*, %struct.PVNode** %16, align 8
  %208 = bitcast %struct.PVNode* %207 to %struct.PVH*
  store %struct.PVH* %208, %struct.PVH** %4, align 8
  br label %209

209:                                              ; preds = %201, %80, %35
  %210 = load %struct.PVH*, %struct.PVH** %4, align 8
  ret %struct.PVH* %210
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVNode* @pn_replace_child(%struct.PVNode* %0, i8 %1, %struct.PVH* %2) #0 {
  %4 = alloca %struct.PVH*, align 8
  %5 = alloca i8, align 1
  %6 = alloca %struct.PVNode*, align 8
  %7 = alloca %struct.PVNode*, align 8
  store %struct.PVH* %2, %struct.PVH** %4, align 8
  store i8 %1, i8* %5, align 1
  store %struct.PVNode* %0, %struct.PVNode** %6, align 8
  %8 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %9 = call %struct.PVNode* @pn_copy(%struct.PVNode* %8)
  store %struct.PVNode* %9, %struct.PVNode** %7, align 8
  %10 = load i8, i8* %5, align 1
  %11 = load %struct.PVH*, %struct.PVH** %4, align 8
  %12 = load %struct.PVNode*, %struct.PVNode** %7, align 8
  call void @pn_set_child(%struct.PVNode* %12, %struct.PVH* %11, i8 %10)
  %13 = load %struct.PVNode*, %struct.PVNode** %7, align 8
  call void @pn_update_index_table(%struct.PVNode* %13)
  %14 = load %struct.PVNode*, %struct.PVNode** %7, align 8
  ret %struct.PVNode* %14
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVNode* @pn_remove_child(%struct.PVNode* %0, i8 %1) #0 {
  %3 = alloca %struct.PVNode*, align 8
  %4 = alloca i8, align 1
  %5 = alloca %struct.PVNode*, align 8
  %6 = alloca %struct.PVNode*, align 8
  %7 = alloca %struct.PVH*, align 8
  %8 = alloca i32, align 4
  %9 = alloca i8, align 1
  store i8 %1, i8* %4, align 1
  store %struct.PVNode* %0, %struct.PVNode** %5, align 8
  %10 = load %struct.PVNode*, %struct.PVNode** %5, align 8
  %11 = call %struct.PVNode* @pn_copy(%struct.PVNode* %10)
  store %struct.PVNode* %11, %struct.PVNode** %6, align 8
  %12 = load %struct.PVNode*, %struct.PVNode** %5, align 8
  %13 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %12, i32 0, i32 2
  %14 = load i8, i8* %4, align 1
  %15 = zext i8 %14 to i64
  %16 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %13, i64 0, i64 %15
  %17 = load %struct.PVH*, %struct.PVH** %16, align 8
  store %struct.PVH* %17, %struct.PVH** %7, align 8
  %18 = load %struct.PVH*, %struct.PVH** %7, align 8
  %19 = icmp ne %struct.PVH* %18, null
  br i1 %19, label %22, label %20

20:                                               ; preds = %2
  %21 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  store %struct.PVNode* %21, %struct.PVNode** %3, align 8
  br label %66

22:                                               ; preds = %2
  %23 = load %struct.PVH*, %struct.PVH** %7, align 8
  %24 = getelementptr inbounds %struct.PVH, %struct.PVH* %23, i32 0, i32 2
  %25 = load i32, i32* %24, align 4
  store i32 %25, i32* %8, align 4
  %26 = load %struct.PVH*, %struct.PVH** %7, align 8
  call void @pn_free(%struct.PVH* %26)
  %27 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %28 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %27, i32 0, i32 2
  %29 = load i8, i8* %4, align 1
  %30 = zext i8 %29 to i64
  %31 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %28, i64 0, i64 %30
  store %struct.PVH* null, %struct.PVH** %31, align 8
  %32 = load i8, i8* %4, align 1
  store i8 %32, i8* %9, align 1
  br label %33

33:                                               ; preds = %51, %22
  %34 = load i8, i8* %9, align 1
  %35 = zext i8 %34 to i32
  %36 = icmp slt i32 %35, 31
  br i1 %36, label %37, label %54

37:                                               ; preds = %33
  %38 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %39 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %38, i32 0, i32 2
  %40 = load i8, i8* %9, align 1
  %41 = zext i8 %40 to i32
  %42 = add nsw i32 %41, 1
  %43 = sext i32 %42 to i64
  %44 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %39, i64 0, i64 %43
  %45 = load %struct.PVH*, %struct.PVH** %44, align 8
  %46 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %47 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %46, i32 0, i32 2
  %48 = load i8, i8* %9, align 1
  %49 = zext i8 %48 to i64
  %50 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %47, i64 0, i64 %49
  store %struct.PVH* %45, %struct.PVH** %50, align 8
  br label %51

51:                                               ; preds = %37
  %52 = load i8, i8* %9, align 1
  %53 = add i8 %52, 1
  store i8 %53, i8* %9, align 1
  br label %33, !llvm.loop !22

54:                                               ; preds = %33
  %55 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %56 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %55, i32 0, i32 2
  %57 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %56, i64 0, i64 31
  store %struct.PVH* null, %struct.PVH** %57, align 8
  %58 = load i32, i32* %8, align 4
  %59 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  %60 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %59, i32 0, i32 0
  %61 = getelementptr inbounds %struct.PVH, %struct.PVH* %60, i32 0, i32 2
  %62 = load i32, i32* %61, align 8
  %63 = sub i32 %62, %58
  store i32 %63, i32* %61, align 8
  %64 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  call void @pn_update_index_table(%struct.PVNode* %64)
  %65 = load %struct.PVNode*, %struct.PVNode** %6, align 8
  store %struct.PVNode* %65, %struct.PVNode** %3, align 8
  br label %66

66:                                               ; preds = %54, %20
  %67 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  ret %struct.PVNode* %67
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8 @pn_fits_into_one_node(%struct.PVH* %0, %struct.PVH* %1) #0 {
  %3 = alloca i8, align 1
  %4 = alloca %struct.PVH*, align 8
  %5 = alloca %struct.PVH*, align 8
  store %struct.PVH* %1, %struct.PVH** %4, align 8
  store %struct.PVH* %0, %struct.PVH** %5, align 8
  %6 = load %struct.PVH*, %struct.PVH** %5, align 8
  %7 = getelementptr inbounds %struct.PVH, %struct.PVH* %6, i32 0, i32 0
  %8 = load i8, i8* %7, align 4
  %9 = zext i8 %8 to i32
  %10 = icmp eq i32 %9, 0
  br i1 %10, label %11, label %22

11:                                               ; preds = %2
  %12 = load %struct.PVH*, %struct.PVH** %5, align 8
  %13 = getelementptr inbounds %struct.PVH, %struct.PVH* %12, i32 0, i32 2
  %14 = load i32, i32* %13, align 4
  %15 = load %struct.PVH*, %struct.PVH** %4, align 8
  %16 = getelementptr inbounds %struct.PVH, %struct.PVH* %15, i32 0, i32 2
  %17 = load i32, i32* %16, align 4
  %18 = add i32 %14, %17
  %19 = icmp ule i32 %18, 32
  br i1 %19, label %20, label %21

20:                                               ; preds = %11
  store i8 1, i8* %3, align 1
  br label %33

21:                                               ; preds = %11
  store i8 0, i8* %3, align 1
  br label %33

22:                                               ; preds = %2
  %23 = load %struct.PVH*, %struct.PVH** %5, align 8
  %24 = call i8 @pn_branch_count(%struct.PVH* %23)
  %25 = zext i8 %24 to i32
  %26 = load %struct.PVH*, %struct.PVH** %4, align 8
  %27 = call i8 @pn_branch_count(%struct.PVH* %26)
  %28 = zext i8 %27 to i32
  %29 = add nsw i32 %25, %28
  %30 = icmp sle i32 %29, 32
  br i1 %30, label %31, label %32

31:                                               ; preds = %22
  store i8 1, i8* %3, align 1
  br label %33

32:                                               ; preds = %22
  store i8 0, i8* %3, align 1
  br label %33

33:                                               ; preds = %32, %31, %21, %20
  %34 = load i8, i8* %3, align 1
  ret i8 %34
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVNode* @pn_make_parent(%struct.PVH* %0) #0 {
  %2 = alloca %struct.PVH*, align 8
  %3 = alloca %struct.PVNode*, align 8
  store %struct.PVH* %0, %struct.PVH** %2, align 8
  %4 = load %struct.PVH*, %struct.PVH** %2, align 8
  %5 = getelementptr inbounds %struct.PVH, %struct.PVH* %4, i32 0, i32 0
  %6 = load i8, i8* %5, align 4
  %7 = zext i8 %6 to i32
  %8 = add nsw i32 %7, 1
  %9 = trunc i32 %8 to i8
  %10 = call %struct.PVNode* @pn_new(i8 %9)
  store %struct.PVNode* %10, %struct.PVNode** %3, align 8
  %11 = load %struct.PVH*, %struct.PVH** %2, align 8
  %12 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  call void @pn_set_child(%struct.PVNode* %12, %struct.PVH* %11, i8 0)
  %13 = load %struct.PVNode*, %struct.PVNode** %3, align 8
  ret %struct.PVNode* %13
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @rassign(%struct.PVH** %0, %struct.PVH* %1) #0 {
  %3 = alloca %struct.PVH*, align 8
  %4 = alloca %struct.PVH**, align 8
  store %struct.PVH* %1, %struct.PVH** %3, align 8
  store %struct.PVH** %0, %struct.PVH*** %4, align 8
  %5 = load %struct.PVH**, %struct.PVH*** %4, align 8
  %6 = load %struct.PVH*, %struct.PVH** %5, align 8
  %7 = icmp ne %struct.PVH* %6, null
  br i1 %7, label %8, label %11

8:                                                ; preds = %2
  %9 = load %struct.PVH**, %struct.PVH*** %4, align 8
  %10 = load %struct.PVH*, %struct.PVH** %9, align 8
  call void @pn_free(%struct.PVH* %10)
  br label %11

11:                                               ; preds = %8, %2
  %12 = load %struct.PVH*, %struct.PVH** %3, align 8
  %13 = icmp ne %struct.PVH* %12, null
  br i1 %13, label %14, label %16

14:                                               ; preds = %11
  %15 = load %struct.PVH*, %struct.PVH** %3, align 8
  call void @pn_incr_ref(%struct.PVH* %15)
  br label %16

16:                                               ; preds = %14, %11
  %17 = load %struct.PVH*, %struct.PVH** %3, align 8
  %18 = load %struct.PVH**, %struct.PVH*** %4, align 8
  store %struct.PVH* %17, %struct.PVH** %18, align 8
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVHead* @pv_concatenate(%struct.PVHead* %0, %struct.PVHead* %1) #0 {
  %3 = alloca %struct.PVHead*, align 8
  %4 = alloca %struct.PVHead*, align 8
  %5 = alloca %struct.PVHead*, align 8
  %6 = alloca %struct.PVH*, align 8
  %7 = alloca %struct.PVH*, align 8
  %8 = alloca i8*, align 8
  %9 = alloca i64, align 8
  %10 = alloca i64, align 8
  %11 = alloca i8, align 1
  %12 = alloca i8, align 1
  %13 = alloca %struct.PVH*, align 8
  %14 = alloca %struct.PVH*, align 8
  %15 = alloca i8, align 1
  %16 = alloca %struct.PVNode*, align 8
  %17 = alloca %struct.PVNode*, align 8
  %18 = alloca %struct.PVNode*, align 8
  %19 = alloca %struct.PVH*, align 8
  %20 = alloca %struct.PVNode*, align 8
  %21 = alloca i8, align 1
  %22 = alloca %struct.PVH*, align 8
  %23 = alloca %struct.PVNode*, align 8
  %24 = alloca i8, align 1
  %25 = alloca %struct.PVNode*, align 8
  %26 = alloca %struct.PVNode*, align 8
  %27 = alloca %struct.PVNode*, align 8
  %28 = alloca %struct.PVNode*, align 8
  %29 = alloca %struct.PVNode*, align 8
  %30 = alloca %struct.PVNode*, align 8
  %31 = alloca %struct.PVNode*, align 8
  %32 = alloca %struct.PVH*, align 8
  %33 = alloca %struct.PVH*, align 8
  %34 = alloca %struct.PVNode*, align 8
  store %struct.PVHead* %1, %struct.PVHead** %4, align 8
  store %struct.PVHead* %0, %struct.PVHead** %5, align 8
  %35 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  %36 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %35, i32 0, i32 2
  %37 = load %struct.PVH*, %struct.PVH** %36, align 8
  store %struct.PVH* %37, %struct.PVH** %6, align 8
  %38 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %39 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %38, i32 0, i32 2
  %40 = load %struct.PVH*, %struct.PVH** %39, align 8
  store %struct.PVH* %40, %struct.PVH** %7, align 8
  %41 = load %struct.PVH*, %struct.PVH** %6, align 8
  %42 = icmp ne %struct.PVH* %41, null
  br i1 %42, label %50, label %43

43:                                               ; preds = %2
  %44 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %45 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %44, i32 0, i32 0
  %46 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %45, i32 0, i32 1
  %47 = load i32, i32* %46, align 4
  %48 = add i32 %47, 1
  store i32 %48, i32* %46, align 4
  %49 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  store %struct.PVHead* %49, %struct.PVHead** %3, align 8
  br label %423

50:                                               ; preds = %2
  %51 = load %struct.PVH*, %struct.PVH** %7, align 8
  %52 = icmp ne %struct.PVH* %51, null
  br i1 %52, label %60, label %53

53:                                               ; preds = %50
  %54 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  %55 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %54, i32 0, i32 0
  %56 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %55, i32 0, i32 1
  %57 = load i32, i32* %56, align 4
  %58 = add i32 %57, 1
  store i32 %58, i32* %56, align 4
  %59 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  store %struct.PVHead* %59, %struct.PVHead** %3, align 8
  br label %423

60:                                               ; preds = %50
  %61 = load %struct.PVH*, %struct.PVH** %6, align 8
  %62 = getelementptr inbounds %struct.PVH, %struct.PVH* %61, i32 0, i32 0
  %63 = load i8, i8* %62, align 4
  %64 = zext i8 %63 to i32
  %65 = add nsw i32 %64, 1
  %66 = zext i32 %65 to i64
  %67 = call i8* @llvm.stacksave()
  store i8* %67, i8** %8, align 8
  %68 = alloca %struct.PVH*, i64 %66, align 16
  store i64 %66, i64* %9, align 8
  %69 = load %struct.PVH*, %struct.PVH** %7, align 8
  %70 = getelementptr inbounds %struct.PVH, %struct.PVH* %69, i32 0, i32 0
  %71 = load i8, i8* %70, align 4
  %72 = zext i8 %71 to i32
  %73 = add nsw i32 %72, 1
  %74 = zext i32 %73 to i64
  %75 = alloca %struct.PVH*, i64 %74, align 16
  store i64 %74, i64* %10, align 8
  store i8 0, i8* %11, align 1
  store i8 0, i8* %12, align 1
  br label %76

76:                                               ; preds = %81, %60
  %77 = load %struct.PVH*, %struct.PVH** %6, align 8
  %78 = getelementptr inbounds %struct.PVH, %struct.PVH* %77, i32 0, i32 0
  %79 = load i8, i8* %78, align 4
  %80 = icmp ne i8 %79, 0
  br i1 %80, label %81, label %91

81:                                               ; preds = %76
  %82 = load %struct.PVH*, %struct.PVH** %6, align 8
  %83 = load i8, i8* %11, align 1
  %84 = add i8 %83, 1
  store i8 %84, i8* %11, align 1
  %85 = zext i8 %83 to i64
  %86 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %85
  store %struct.PVH* %82, %struct.PVH** %86, align 8
  %87 = load %struct.PVH*, %struct.PVH** %6, align 8
  %88 = bitcast %struct.PVH* %87 to %struct.PVNode*
  %89 = call i8* @pn_right_child(%struct.PVNode* %88)
  %90 = bitcast i8* %89 to %struct.PVH*
  store %struct.PVH* %90, %struct.PVH** %6, align 8
  br label %76, !llvm.loop !23

91:                                               ; preds = %76
  %92 = load %struct.PVH*, %struct.PVH** %6, align 8
  %93 = load i8, i8* %11, align 1
  %94 = zext i8 %93 to i64
  %95 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %94
  store %struct.PVH* %92, %struct.PVH** %95, align 8
  br label %96

96:                                               ; preds = %101, %91
  %97 = load %struct.PVH*, %struct.PVH** %7, align 8
  %98 = getelementptr inbounds %struct.PVH, %struct.PVH* %97, i32 0, i32 0
  %99 = load i8, i8* %98, align 4
  %100 = icmp ne i8 %99, 0
  br i1 %100, label %101, label %112

101:                                              ; preds = %96
  %102 = load %struct.PVH*, %struct.PVH** %7, align 8
  %103 = load i8, i8* %12, align 1
  %104 = add i8 %103, 1
  store i8 %104, i8* %12, align 1
  %105 = zext i8 %103 to i64
  %106 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %105
  store %struct.PVH* %102, %struct.PVH** %106, align 8
  %107 = load %struct.PVH*, %struct.PVH** %7, align 8
  %108 = bitcast %struct.PVH* %107 to %struct.PVNode*
  %109 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %108, i32 0, i32 2
  %110 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %109, i64 0, i64 0
  %111 = load %struct.PVH*, %struct.PVH** %110, align 8
  store %struct.PVH* %111, %struct.PVH** %7, align 8
  br label %96, !llvm.loop !24

112:                                              ; preds = %96
  %113 = load %struct.PVH*, %struct.PVH** %7, align 8
  %114 = load i8, i8* %12, align 1
  %115 = zext i8 %114 to i64
  %116 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %115
  store %struct.PVH* %113, %struct.PVH** %116, align 8
  store %struct.PVH* null, %struct.PVH** %13, align 8
  store %struct.PVH* null, %struct.PVH** %14, align 8
  %117 = load i8, i8* %11, align 1
  %118 = zext i8 %117 to i64
  %119 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %118
  %120 = load %struct.PVH*, %struct.PVH** %119, align 8
  call void @rassign(%struct.PVH** %13, %struct.PVH* %120)
  %121 = load i8, i8* %12, align 1
  %122 = zext i8 %121 to i64
  %123 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %122
  %124 = load %struct.PVH*, %struct.PVH** %123, align 8
  call void @rassign(%struct.PVH** %14, %struct.PVH* %124)
  br label %125

125:                                              ; preds = %380, %112
  %126 = load i8, i8* %11, align 1
  %127 = zext i8 %126 to i32
  %128 = icmp sgt i32 %127, 0
  br i1 %128, label %133, label %129

129:                                              ; preds = %125
  %130 = load i8, i8* %12, align 1
  %131 = zext i8 %130 to i32
  %132 = icmp sgt i32 %131, 0
  br label %133

133:                                              ; preds = %129, %125
  %134 = phi i1 [ true, %125 ], [ %132, %129 ]
  br i1 %134, label %135, label %381

135:                                              ; preds = %133
  store i8 0, i8* %15, align 1
  %136 = load %struct.PVH*, %struct.PVH** %13, align 8
  %137 = icmp ne %struct.PVH* %136, null
  br i1 %137, label %138, label %164

138:                                              ; preds = %135
  %139 = load %struct.PVH*, %struct.PVH** %14, align 8
  %140 = icmp ne %struct.PVH* %139, null
  br i1 %140, label %141, label %164

141:                                              ; preds = %138
  %142 = load %struct.PVH*, %struct.PVH** %13, align 8
  %143 = getelementptr inbounds %struct.PVH, %struct.PVH* %142, i32 0, i32 0
  %144 = load i8, i8* %143, align 4
  %145 = zext i8 %144 to i32
  %146 = icmp sgt i32 %145, 0
  br i1 %146, label %147, label %164

147:                                              ; preds = %141
  %148 = load %struct.PVH*, %struct.PVH** %14, align 8
  %149 = bitcast %struct.PVH* %148 to %struct.PVNode*
  %150 = load %struct.PVH*, %struct.PVH** %13, align 8
  %151 = bitcast %struct.PVH* %150 to %struct.PVNode*
  %152 = call i8 @pn_needs_rebalancing(%struct.PVNode* %151, %struct.PVNode* %149)
  %153 = zext i8 %152 to i32
  %154 = icmp ne i32 %153, 0
  br i1 %154, label %155, label %164

155:                                              ; preds = %147
  %156 = load %struct.PVH*, %struct.PVH** %14, align 8
  %157 = bitcast %struct.PVH* %156 to %struct.PVNode*
  %158 = load %struct.PVH*, %struct.PVH** %13, align 8
  %159 = bitcast %struct.PVH* %158 to %struct.PVNode*
  call void @pn_balance_level(%struct.PVNode* %159, %struct.PVNode* %157, %struct.PVNode** %16, %struct.PVNode** %17)
  %160 = load %struct.PVNode*, %struct.PVNode** %16, align 8
  %161 = bitcast %struct.PVNode* %160 to %struct.PVH*
  call void @rassign(%struct.PVH** %13, %struct.PVH* %161)
  %162 = load %struct.PVNode*, %struct.PVNode** %17, align 8
  %163 = bitcast %struct.PVNode* %162 to %struct.PVH*
  call void @rassign(%struct.PVH** %14, %struct.PVH* %163)
  store i8 1, i8* %15, align 1
  br label %164

164:                                              ; preds = %155, %147, %141, %138, %135
  %165 = load %struct.PVH*, %struct.PVH** %13, align 8
  %166 = icmp ne %struct.PVH* %165, null
  br i1 %166, label %167, label %276

167:                                              ; preds = %164
  %168 = load %struct.PVH*, %struct.PVH** %14, align 8
  %169 = icmp ne %struct.PVH* %168, null
  br i1 %169, label %170, label %276

170:                                              ; preds = %167
  %171 = load i8, i8* %11, align 1
  %172 = zext i8 %171 to i32
  %173 = icmp eq i32 %172, 0
  br i1 %173, label %174, label %199

174:                                              ; preds = %170
  %175 = load %struct.PVH*, %struct.PVH** %14, align 8
  %176 = load %struct.PVH*, %struct.PVH** %13, align 8
  %177 = call i8 @pn_fits_into_one_node(%struct.PVH* %176, %struct.PVH* %175)
  %178 = icmp ne i8 %177, 0
  br i1 %178, label %179, label %194

179:                                              ; preds = %174
  %180 = load i8, i8* %12, align 1
  %181 = zext i8 %180 to i32
  %182 = sub nsw i32 %181, 1
  %183 = sext i32 %182 to i64
  %184 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %183
  %185 = load %struct.PVH*, %struct.PVH** %184, align 8
  %186 = bitcast %struct.PVH* %185 to %struct.PVNode*
  store %struct.PVNode* %186, %struct.PVNode** %18, align 8
  %187 = load %struct.PVH*, %struct.PVH** %14, align 8
  %188 = load %struct.PVH*, %struct.PVH** %13, align 8
  %189 = call %struct.PVH* @pn_join_nodes(%struct.PVH* %188, %struct.PVH* %187, %struct.PVH** null)
  store %struct.PVH* %189, %struct.PVH** %19, align 8
  call void @rassign(%struct.PVH** %13, %struct.PVH* null)
  %190 = load %struct.PVH*, %struct.PVH** %19, align 8
  %191 = load %struct.PVNode*, %struct.PVNode** %18, align 8
  %192 = call %struct.PVNode* @pn_replace_child(%struct.PVNode* %191, i8 0, %struct.PVH* %190)
  %193 = bitcast %struct.PVNode* %192 to %struct.PVH*
  call void @rassign(%struct.PVH** %14, %struct.PVH* %193)
  br label %198

194:                                              ; preds = %174
  %195 = load %struct.PVH*, %struct.PVH** %13, align 8
  %196 = call %struct.PVNode* @pn_make_parent(%struct.PVH* %195)
  %197 = bitcast %struct.PVNode* %196 to %struct.PVH*
  call void @rassign(%struct.PVH** %13, %struct.PVH* %197)
  br label %198

198:                                              ; preds = %194, %179
  br label %199

199:                                              ; preds = %198, %170
  %200 = load i8, i8* %12, align 1
  %201 = zext i8 %200 to i32
  %202 = icmp eq i32 %201, 0
  br i1 %202, label %203, label %231

203:                                              ; preds = %199
  %204 = load %struct.PVH*, %struct.PVH** %14, align 8
  %205 = load %struct.PVH*, %struct.PVH** %13, align 8
  %206 = call i8 @pn_fits_into_one_node(%struct.PVH* %205, %struct.PVH* %204)
  %207 = icmp ne i8 %206, 0
  br i1 %207, label %208, label %226

208:                                              ; preds = %203
  %209 = load i8, i8* %11, align 1
  %210 = zext i8 %209 to i32
  %211 = sub nsw i32 %210, 1
  %212 = sext i32 %211 to i64
  %213 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %212
  %214 = load %struct.PVH*, %struct.PVH** %213, align 8
  %215 = bitcast %struct.PVH* %214 to %struct.PVNode*
  store %struct.PVNode* %215, %struct.PVNode** %20, align 8
  %216 = load %struct.PVNode*, %struct.PVNode** %20, align 8
  %217 = call i8 @pn_right_child_index(%struct.PVNode* %216)
  store i8 %217, i8* %21, align 1
  %218 = load %struct.PVH*, %struct.PVH** %14, align 8
  %219 = load %struct.PVH*, %struct.PVH** %13, align 8
  %220 = call %struct.PVH* @pn_join_nodes(%struct.PVH* %219, %struct.PVH* %218, %struct.PVH** null)
  store %struct.PVH* %220, %struct.PVH** %22, align 8
  %221 = load %struct.PVH*, %struct.PVH** %22, align 8
  %222 = load i8, i8* %21, align 1
  %223 = load %struct.PVNode*, %struct.PVNode** %20, align 8
  %224 = call %struct.PVNode* @pn_replace_child(%struct.PVNode* %223, i8 %222, %struct.PVH* %221)
  %225 = bitcast %struct.PVNode* %224 to %struct.PVH*
  call void @rassign(%struct.PVH** %13, %struct.PVH* %225)
  call void @rassign(%struct.PVH** %14, %struct.PVH* null)
  br label %230

226:                                              ; preds = %203
  %227 = load %struct.PVH*, %struct.PVH** %14, align 8
  %228 = call %struct.PVNode* @pn_make_parent(%struct.PVH* %227)
  %229 = bitcast %struct.PVNode* %228 to %struct.PVH*
  call void @rassign(%struct.PVH** %14, %struct.PVH* %229)
  br label %230

230:                                              ; preds = %226, %208
  br label %231

231:                                              ; preds = %230, %199
  %232 = load i8, i8* %12, align 1
  %233 = zext i8 %232 to i32
  %234 = icmp sgt i32 %233, 0
  br i1 %234, label %235, label %253

235:                                              ; preds = %231
  %236 = load %struct.PVH*, %struct.PVH** %14, align 8
  %237 = load i8, i8* %12, align 1
  %238 = zext i8 %237 to i64
  %239 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %238
  %240 = load %struct.PVH*, %struct.PVH** %239, align 8
  %241 = icmp eq %struct.PVH* %236, %240
  br i1 %241, label %246, label %242

242:                                              ; preds = %235
  %243 = load i8, i8* %15, align 1
  %244 = zext i8 %243 to i32
  %245 = icmp ne i32 %244, 0
  br i1 %245, label %246, label %253

246:                                              ; preds = %242, %235
  %247 = load i8, i8* %12, align 1
  %248 = zext i8 %247 to i32
  %249 = sub nsw i32 %248, 1
  %250 = sext i32 %249 to i64
  %251 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %250
  %252 = load %struct.PVH*, %struct.PVH** %251, align 8
  call void @rassign(%struct.PVH** %14, %struct.PVH* %252)
  br label %253

253:                                              ; preds = %246, %242, %231
  %254 = load i8, i8* %11, align 1
  %255 = zext i8 %254 to i32
  %256 = icmp sgt i32 %255, 0
  br i1 %256, label %257, label %275

257:                                              ; preds = %253
  %258 = load %struct.PVH*, %struct.PVH** %13, align 8
  %259 = load i8, i8* %11, align 1
  %260 = zext i8 %259 to i64
  %261 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %260
  %262 = load %struct.PVH*, %struct.PVH** %261, align 8
  %263 = icmp eq %struct.PVH* %258, %262
  br i1 %263, label %268, label %264

264:                                              ; preds = %257
  %265 = load i8, i8* %15, align 1
  %266 = zext i8 %265 to i32
  %267 = icmp ne i32 %266, 0
  br i1 %267, label %268, label %275

268:                                              ; preds = %264, %257
  %269 = load i8, i8* %11, align 1
  %270 = zext i8 %269 to i32
  %271 = sub nsw i32 %270, 1
  %272 = sext i32 %271 to i64
  %273 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %272
  %274 = load %struct.PVH*, %struct.PVH** %273, align 8
  call void @rassign(%struct.PVH** %13, %struct.PVH* %274)
  br label %275

275:                                              ; preds = %268, %264, %253
  br label %366

276:                                              ; preds = %167, %164
  %277 = load %struct.PVH*, %struct.PVH** %13, align 8
  %278 = icmp ne %struct.PVH* %277, null
  br i1 %278, label %279, label %320

279:                                              ; preds = %276
  %280 = load i8, i8* %11, align 1
  %281 = zext i8 %280 to i32
  %282 = icmp sgt i32 %281, 0
  br i1 %282, label %283, label %299

283:                                              ; preds = %279
  %284 = load i8, i8* %11, align 1
  %285 = zext i8 %284 to i32
  %286 = sub nsw i32 %285, 1
  %287 = sext i32 %286 to i64
  %288 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %287
  %289 = load %struct.PVH*, %struct.PVH** %288, align 8
  %290 = bitcast %struct.PVH* %289 to %struct.PVNode*
  store %struct.PVNode* %290, %struct.PVNode** %23, align 8
  %291 = load %struct.PVNode*, %struct.PVNode** %23, align 8
  %292 = call i8 @pn_right_child_index(%struct.PVNode* %291)
  store i8 %292, i8* %24, align 1
  %293 = load %struct.PVH*, %struct.PVH** %13, align 8
  %294 = load i8, i8* %24, align 1
  %295 = load %struct.PVNode*, %struct.PVNode** %23, align 8
  %296 = call %struct.PVNode* @pn_replace_child(%struct.PVNode* %295, i8 %294, %struct.PVH* %293)
  store %struct.PVNode* %296, %struct.PVNode** %25, align 8
  %297 = load %struct.PVNode*, %struct.PVNode** %25, align 8
  %298 = bitcast %struct.PVNode* %297 to %struct.PVH*
  call void @rassign(%struct.PVH** %13, %struct.PVH* %298)
  br label %303

299:                                              ; preds = %279
  %300 = load %struct.PVH*, %struct.PVH** %13, align 8
  %301 = call %struct.PVNode* @pn_make_parent(%struct.PVH* %300)
  %302 = bitcast %struct.PVNode* %301 to %struct.PVH*
  call void @rassign(%struct.PVH** %13, %struct.PVH* %302)
  br label %303

303:                                              ; preds = %299, %283
  %304 = load i8, i8* %12, align 1
  %305 = zext i8 %304 to i32
  %306 = icmp sgt i32 %305, 0
  br i1 %306, label %307, label %319

307:                                              ; preds = %303
  %308 = load i8, i8* %12, align 1
  %309 = zext i8 %308 to i32
  %310 = sub nsw i32 %309, 1
  %311 = sext i32 %310 to i64
  %312 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %311
  %313 = load %struct.PVH*, %struct.PVH** %312, align 8
  %314 = bitcast %struct.PVH* %313 to %struct.PVNode*
  store %struct.PVNode* %314, %struct.PVNode** %26, align 8
  %315 = load %struct.PVNode*, %struct.PVNode** %26, align 8
  %316 = call %struct.PVNode* @pn_remove_child(%struct.PVNode* %315, i8 0)
  store %struct.PVNode* %316, %struct.PVNode** %27, align 8
  %317 = load %struct.PVNode*, %struct.PVNode** %27, align 8
  %318 = bitcast %struct.PVNode* %317 to %struct.PVH*
  call void @rassign(%struct.PVH** %14, %struct.PVH* %318)
  br label %319

319:                                              ; preds = %307, %303
  br label %365

320:                                              ; preds = %276
  %321 = load i8, i8* %12, align 1
  %322 = zext i8 %321 to i32
  %323 = icmp sgt i32 %322, 0
  br i1 %323, label %324, label %364

324:                                              ; preds = %320
  %325 = load i8, i8* %12, align 1
  %326 = zext i8 %325 to i32
  %327 = icmp sgt i32 %326, 0
  br i1 %327, label %328, label %341

328:                                              ; preds = %324
  %329 = load i8, i8* %12, align 1
  %330 = zext i8 %329 to i32
  %331 = sub nsw i32 %330, 1
  %332 = sext i32 %331 to i64
  %333 = getelementptr inbounds %struct.PVH*, %struct.PVH** %75, i64 %332
  %334 = load %struct.PVH*, %struct.PVH** %333, align 8
  %335 = bitcast %struct.PVH* %334 to %struct.PVNode*
  store %struct.PVNode* %335, %struct.PVNode** %28, align 8
  %336 = load %struct.PVH*, %struct.PVH** %14, align 8
  %337 = load %struct.PVNode*, %struct.PVNode** %28, align 8
  %338 = call %struct.PVNode* @pn_replace_child(%struct.PVNode* %337, i8 0, %struct.PVH* %336)
  store %struct.PVNode* %338, %struct.PVNode** %29, align 8
  %339 = load %struct.PVNode*, %struct.PVNode** %29, align 8
  %340 = bitcast %struct.PVNode* %339 to %struct.PVH*
  call void @rassign(%struct.PVH** %14, %struct.PVH* %340)
  br label %345

341:                                              ; preds = %324
  %342 = load %struct.PVH*, %struct.PVH** %14, align 8
  %343 = call %struct.PVNode* @pn_make_parent(%struct.PVH* %342)
  %344 = bitcast %struct.PVNode* %343 to %struct.PVH*
  call void @rassign(%struct.PVH** %14, %struct.PVH* %344)
  br label %345

345:                                              ; preds = %341, %328
  %346 = load i8, i8* %11, align 1
  %347 = zext i8 %346 to i32
  %348 = icmp sgt i32 %347, 0
  br i1 %348, label %349, label %363

349:                                              ; preds = %345
  %350 = load i8, i8* %11, align 1
  %351 = zext i8 %350 to i32
  %352 = sub nsw i32 %351, 1
  %353 = sext i32 %352 to i64
  %354 = getelementptr inbounds %struct.PVH*, %struct.PVH** %68, i64 %353
  %355 = load %struct.PVH*, %struct.PVH** %354, align 8
  %356 = bitcast %struct.PVH* %355 to %struct.PVNode*
  store %struct.PVNode* %356, %struct.PVNode** %30, align 8
  %357 = load %struct.PVNode*, %struct.PVNode** %30, align 8
  %358 = call i8 @pn_right_child_index(%struct.PVNode* %357)
  %359 = load %struct.PVNode*, %struct.PVNode** %30, align 8
  %360 = call %struct.PVNode* @pn_remove_child(%struct.PVNode* %359, i8 %358)
  store %struct.PVNode* %360, %struct.PVNode** %31, align 8
  %361 = load %struct.PVNode*, %struct.PVNode** %31, align 8
  %362 = bitcast %struct.PVNode* %361 to %struct.PVH*
  call void @rassign(%struct.PVH** %13, %struct.PVH* %362)
  br label %363

363:                                              ; preds = %349, %345
  br label %364

364:                                              ; preds = %363, %320
  br label %365

365:                                              ; preds = %364, %319
  br label %366

366:                                              ; preds = %365, %275
  %367 = load i8, i8* %11, align 1
  %368 = zext i8 %367 to i32
  %369 = icmp sgt i32 %368, 0
  br i1 %369, label %370, label %373

370:                                              ; preds = %366
  %371 = load i8, i8* %11, align 1
  %372 = add i8 %371, -1
  store i8 %372, i8* %11, align 1
  br label %373

373:                                              ; preds = %370, %366
  %374 = load i8, i8* %12, align 1
  %375 = zext i8 %374 to i32
  %376 = icmp sgt i32 %375, 0
  br i1 %376, label %377, label %380

377:                                              ; preds = %373
  %378 = load i8, i8* %12, align 1
  %379 = add i8 %378, -1
  store i8 %379, i8* %12, align 1
  br label %380

380:                                              ; preds = %377, %373
  br label %125, !llvm.loop !25

381:                                              ; preds = %133
  %382 = load %struct.PVH*, %struct.PVH** %13, align 8
  %383 = icmp ne %struct.PVH* %382, null
  br i1 %383, label %384, label %411

384:                                              ; preds = %381
  %385 = load %struct.PVH*, %struct.PVH** %14, align 8
  %386 = icmp ne %struct.PVH* %385, null
  br i1 %386, label %387, label %411

387:                                              ; preds = %384
  %388 = load %struct.PVH*, %struct.PVH** %14, align 8
  %389 = load %struct.PVH*, %struct.PVH** %13, align 8
  %390 = call %struct.PVH* @pn_join_nodes(%struct.PVH* %389, %struct.PVH* %388, %struct.PVH** %33)
  store %struct.PVH* %390, %struct.PVH** %32, align 8
  %391 = load %struct.PVH*, %struct.PVH** %13, align 8
  call void @pn_free(%struct.PVH* %391)
  %392 = load %struct.PVH*, %struct.PVH** %14, align 8
  call void @pn_free(%struct.PVH* %392)
  %393 = load %struct.PVH*, %struct.PVH** %33, align 8
  %394 = icmp ne %struct.PVH* %393, null
  br i1 %394, label %395, label %410

395:                                              ; preds = %387
  %396 = load %struct.PVH*, %struct.PVH** %32, align 8
  %397 = getelementptr inbounds %struct.PVH, %struct.PVH* %396, i32 0, i32 0
  %398 = load i8, i8* %397, align 4
  %399 = zext i8 %398 to i32
  %400 = add nsw i32 %399, 1
  %401 = trunc i32 %400 to i8
  %402 = call %struct.PVNode* @pn_new(i8 %401)
  store %struct.PVNode* %402, %struct.PVNode** %34, align 8
  %403 = load %struct.PVH*, %struct.PVH** %32, align 8
  %404 = load %struct.PVNode*, %struct.PVNode** %34, align 8
  call void @pn_set_child(%struct.PVNode* %404, %struct.PVH* %403, i8 0)
  %405 = load %struct.PVH*, %struct.PVH** %33, align 8
  %406 = load %struct.PVNode*, %struct.PVNode** %34, align 8
  call void @pn_set_child(%struct.PVNode* %406, %struct.PVH* %405, i8 1)
  %407 = load %struct.PVNode*, %struct.PVNode** %34, align 8
  call void @pn_update_index_table(%struct.PVNode* %407)
  %408 = load %struct.PVNode*, %struct.PVNode** %34, align 8
  %409 = bitcast %struct.PVNode* %408 to %struct.PVH*
  store %struct.PVH* %409, %struct.PVH** %32, align 8
  br label %410

410:                                              ; preds = %395, %387
  br label %419

411:                                              ; preds = %384, %381
  %412 = load %struct.PVH*, %struct.PVH** %13, align 8
  %413 = icmp ne %struct.PVH* %412, null
  br i1 %413, label %414, label %416

414:                                              ; preds = %411
  %415 = load %struct.PVH*, %struct.PVH** %13, align 8
  store %struct.PVH* %415, %struct.PVH** %32, align 8
  br label %418

416:                                              ; preds = %411
  %417 = load %struct.PVH*, %struct.PVH** %14, align 8
  store %struct.PVH* %417, %struct.PVH** %32, align 8
  br label %418

418:                                              ; preds = %416, %414
  br label %419

419:                                              ; preds = %418, %410
  %420 = load %struct.PVH*, %struct.PVH** %32, align 8
  %421 = call %struct.PVHead* @pv_construct(%struct.PVH* %420)
  store %struct.PVHead* %421, %struct.PVHead** %3, align 8
  %422 = load i8*, i8** %8, align 8
  call void @llvm.stackrestore(i8* %422)
  br label %423

423:                                              ; preds = %419, %53, %43
  %424 = load %struct.PVHead*, %struct.PVHead** %3, align 8
  ret %struct.PVHead* %424
}

; Function Attrs: nofree nosync nounwind willreturn
declare i8* @llvm.stacksave() #3

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @pn_balance_level(%struct.PVNode* %0, %struct.PVNode* %1, %struct.PVNode** %2, %struct.PVNode** %3) #0 {
  %5 = alloca %struct.PVNode**, align 8
  %6 = alloca %struct.PVNode**, align 8
  %7 = alloca %struct.PVNode*, align 8
  %8 = alloca %struct.PVNode*, align 8
  %9 = alloca %struct.PVNode*, align 8
  %10 = alloca %struct.PVNode*, align 8
  %11 = alloca %struct.PVH*, align 8
  %12 = alloca %struct.PVH*, align 8
  %13 = alloca i8, align 1
  %14 = alloca i8, align 1
  %15 = alloca %struct.PVH*, align 8
  %16 = alloca %struct.PVH*, align 8
  store %struct.PVNode** %3, %struct.PVNode*** %5, align 8
  store %struct.PVNode** %2, %struct.PVNode*** %6, align 8
  store %struct.PVNode* %1, %struct.PVNode** %7, align 8
  store %struct.PVNode* %0, %struct.PVNode** %8, align 8
  %17 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %18 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %17, i32 0, i32 0
  %19 = getelementptr inbounds %struct.PVH, %struct.PVH* %18, i32 0, i32 0
  %20 = load i8, i8* %19, align 8
  %21 = call %struct.PVNode* @pn_new(i8 %20)
  store %struct.PVNode* %21, %struct.PVNode** %9, align 8
  %22 = load %struct.PVNode*, %struct.PVNode** %7, align 8
  %23 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %22, i32 0, i32 0
  %24 = getelementptr inbounds %struct.PVH, %struct.PVH* %23, i32 0, i32 0
  %25 = load i8, i8* %24, align 8
  %26 = call %struct.PVNode* @pn_new(i8 %25)
  store %struct.PVNode* %26, %struct.PVNode** %10, align 8
  %27 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %28 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %27, i32 0, i32 2
  %29 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %28, i64 0, i64 0
  %30 = load %struct.PVH*, %struct.PVH** %29, align 8
  store %struct.PVH* %30, %struct.PVH** %11, align 8
  store %struct.PVH* null, %struct.PVH** %12, align 8
  store i8 0, i8* %13, align 1
  store i8 1, i8* %14, align 1
  br label %31

31:                                               ; preds = %136, %4
  %32 = load i8, i8* %14, align 1
  %33 = zext i8 %32 to i32
  %34 = icmp slt i32 %33, 64
  br i1 %34, label %35, label %139

35:                                               ; preds = %31
  %36 = load %struct.PVH*, %struct.PVH** %11, align 8
  %37 = icmp ne %struct.PVH* %36, null
  br i1 %37, label %38, label %93

38:                                               ; preds = %35
  %39 = load i8, i8* %14, align 1
  %40 = zext i8 %39 to i32
  %41 = icmp slt i32 %40, 32
  br i1 %41, label %42, label %49

42:                                               ; preds = %38
  %43 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %44 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %43, i32 0, i32 2
  %45 = load i8, i8* %14, align 1
  %46 = zext i8 %45 to i64
  %47 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %44, i64 0, i64 %46
  %48 = load %struct.PVH*, %struct.PVH** %47, align 8
  store %struct.PVH* %48, %struct.PVH** %12, align 8
  br label %58

49:                                               ; preds = %38
  %50 = load %struct.PVNode*, %struct.PVNode** %7, align 8
  %51 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %50, i32 0, i32 2
  %52 = load i8, i8* %14, align 1
  %53 = zext i8 %52 to i32
  %54 = sub nsw i32 %53, 32
  %55 = sext i32 %54 to i64
  %56 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %51, i64 0, i64 %55
  %57 = load %struct.PVH*, %struct.PVH** %56, align 8
  store %struct.PVH* %57, %struct.PVH** %12, align 8
  br label %58

58:                                               ; preds = %49, %42
  %59 = load %struct.PVH*, %struct.PVH** %12, align 8
  %60 = icmp ne %struct.PVH* %59, null
  br i1 %60, label %61, label %92

61:                                               ; preds = %58
  %62 = load %struct.PVH*, %struct.PVH** %12, align 8
  %63 = load %struct.PVH*, %struct.PVH** %11, align 8
  %64 = call %struct.PVH* @pn_join_nodes(%struct.PVH* %63, %struct.PVH* %62, %struct.PVH** %15)
  store %struct.PVH* %64, %struct.PVH** %16, align 8
  %65 = load i8, i8* %13, align 1
  %66 = zext i8 %65 to i32
  %67 = icmp sgt i32 %66, 0
  br i1 %67, label %68, label %73

68:                                               ; preds = %61
  %69 = load %struct.PVH*, %struct.PVH** %11, align 8
  %70 = icmp ne %struct.PVH* %69, null
  br i1 %70, label %71, label %73

71:                                               ; preds = %68
  %72 = load %struct.PVH*, %struct.PVH** %11, align 8
  call void @pn_free(%struct.PVH* %72)
  br label %73

73:                                               ; preds = %71, %68, %61
  %74 = load %struct.PVH*, %struct.PVH** %15, align 8
  store %struct.PVH* %74, %struct.PVH** %11, align 8
  %75 = load i8, i8* %13, align 1
  %76 = zext i8 %75 to i32
  %77 = icmp slt i32 %76, 32
  br i1 %77, label %78, label %82

78:                                               ; preds = %73
  %79 = load i8, i8* %13, align 1
  %80 = load %struct.PVH*, %struct.PVH** %16, align 8
  %81 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  call void @pn_set_child(%struct.PVNode* %81, %struct.PVH* %80, i8 %79)
  br label %89

82:                                               ; preds = %73
  %83 = load i8, i8* %13, align 1
  %84 = zext i8 %83 to i32
  %85 = sub nsw i32 %84, 32
  %86 = trunc i32 %85 to i8
  %87 = load %struct.PVH*, %struct.PVH** %16, align 8
  %88 = load %struct.PVNode*, %struct.PVNode** %10, align 8
  call void @pn_set_child(%struct.PVNode* %88, %struct.PVH* %87, i8 %86)
  br label %89

89:                                               ; preds = %82, %78
  %90 = load i8, i8* %13, align 1
  %91 = add i8 %90, 1
  store i8 %91, i8* %13, align 1
  br label %92

92:                                               ; preds = %89, %58
  br label %135

93:                                               ; preds = %35
  %94 = load i8, i8* %14, align 1
  %95 = zext i8 %94 to i32
  %96 = icmp slt i32 %95, 32
  br i1 %96, label %97, label %104

97:                                               ; preds = %93
  %98 = load %struct.PVNode*, %struct.PVNode** %8, align 8
  %99 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %98, i32 0, i32 2
  %100 = load i8, i8* %14, align 1
  %101 = zext i8 %100 to i64
  %102 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %99, i64 0, i64 %101
  %103 = load %struct.PVH*, %struct.PVH** %102, align 8
  store %struct.PVH* %103, %struct.PVH** %12, align 8
  br label %113

104:                                              ; preds = %93
  %105 = load %struct.PVNode*, %struct.PVNode** %7, align 8
  %106 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %105, i32 0, i32 2
  %107 = load i8, i8* %14, align 1
  %108 = zext i8 %107 to i32
  %109 = sub nsw i32 %108, 32
  %110 = sext i32 %109 to i64
  %111 = getelementptr inbounds [32 x %struct.PVH*], [32 x %struct.PVH*]* %106, i64 0, i64 %110
  %112 = load %struct.PVH*, %struct.PVH** %111, align 8
  store %struct.PVH* %112, %struct.PVH** %12, align 8
  br label %113

113:                                              ; preds = %104, %97
  %114 = load i8, i8* %13, align 1
  %115 = zext i8 %114 to i32
  %116 = icmp slt i32 %115, 32
  br i1 %116, label %117, label %121

117:                                              ; preds = %113
  %118 = load i8, i8* %13, align 1
  %119 = load %struct.PVH*, %struct.PVH** %12, align 8
  %120 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  call void @pn_set_child(%struct.PVNode* %120, %struct.PVH* %119, i8 %118)
  br label %128

121:                                              ; preds = %113
  %122 = load i8, i8* %13, align 1
  %123 = zext i8 %122 to i32
  %124 = sub nsw i32 %123, 32
  %125 = trunc i32 %124 to i8
  %126 = load %struct.PVH*, %struct.PVH** %12, align 8
  %127 = load %struct.PVNode*, %struct.PVNode** %10, align 8
  call void @pn_set_child(%struct.PVNode* %127, %struct.PVH* %126, i8 %125)
  br label %128

128:                                              ; preds = %121, %117
  %129 = load %struct.PVH*, %struct.PVH** %12, align 8
  %130 = icmp ne %struct.PVH* %129, null
  br i1 %130, label %131, label %134

131:                                              ; preds = %128
  %132 = load i8, i8* %13, align 1
  %133 = add i8 %132, 1
  store i8 %133, i8* %13, align 1
  br label %134

134:                                              ; preds = %131, %128
  br label %135

135:                                              ; preds = %134, %92
  br label %136

136:                                              ; preds = %135
  %137 = load i8, i8* %14, align 1
  %138 = add i8 %137, 1
  store i8 %138, i8* %14, align 1
  br label %31, !llvm.loop !26

139:                                              ; preds = %31
  %140 = load %struct.PVH*, %struct.PVH** %11, align 8
  %141 = icmp ne %struct.PVH* %140, null
  br i1 %141, label %142, label %145

142:                                              ; preds = %139
  %143 = call %struct._iobuf* @__acrt_iob_func(i32 2)
  %144 = call i32 (%struct._iobuf*, i8*, ...) @fprintf(%struct._iobuf* %143, i8* getelementptr inbounds ([37 x i8], [37 x i8]* @"??_C@_0CF@JJIBIHF@balance?5failed?5to?5compress?5a?5vec@", i64 0, i64 0))
  call void @exit(i32 1) #5
  unreachable

145:                                              ; preds = %139
  %146 = load %struct.PVNode*, %struct.PVNode** %10, align 8
  %147 = getelementptr inbounds %struct.PVNode, %struct.PVNode* %146, i32 0, i32 0
  %148 = getelementptr inbounds %struct.PVH, %struct.PVH* %147, i32 0, i32 2
  %149 = load i32, i32* %148, align 8
  %150 = icmp eq i32 %149, 0
  br i1 %150, label %151, label %154

151:                                              ; preds = %145
  %152 = load %struct.PVNode*, %struct.PVNode** %10, align 8
  %153 = bitcast %struct.PVNode* %152 to %struct.PVH*
  call void @pn_free(%struct.PVH* %153)
  store %struct.PVNode* null, %struct.PVNode** %10, align 8
  br label %156

154:                                              ; preds = %145
  %155 = load %struct.PVNode*, %struct.PVNode** %10, align 8
  call void @pn_update_index_table(%struct.PVNode* %155)
  br label %156

156:                                              ; preds = %154, %151
  %157 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  call void @pn_update_index_table(%struct.PVNode* %157)
  %158 = load %struct.PVNode*, %struct.PVNode** %9, align 8
  %159 = load %struct.PVNode**, %struct.PVNode*** %6, align 8
  store %struct.PVNode* %158, %struct.PVNode** %159, align 8
  %160 = load %struct.PVNode*, %struct.PVNode** %10, align 8
  %161 = load %struct.PVNode**, %struct.PVNode*** %5, align 8
  store %struct.PVNode* %160, %struct.PVNode** %161, align 8
  ret void
}

; Function Attrs: nofree nosync nounwind willreturn
declare void @llvm.stackrestore(i8*) #3

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i8 @pv_uint16_equals(%struct.PVHead* %0, %struct.PVHead* %1) #0 {
  %3 = alloca i8, align 1
  %4 = alloca %struct.PVHead*, align 8
  %5 = alloca %struct.PVHead*, align 8
  %6 = alloca i32, align 4
  store %struct.PVHead* %1, %struct.PVHead** %4, align 8
  store %struct.PVHead* %0, %struct.PVHead** %5, align 8
  %7 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  %8 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %7, i32 0, i32 1
  %9 = load i32, i32* %8, align 8
  %10 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %11 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %10, i32 0, i32 1
  %12 = load i32, i32* %11, align 8
  %13 = icmp ne i32 %9, %12
  br i1 %13, label %14, label %15

14:                                               ; preds = %2
  store i8 0, i8* %3, align 1
  br label %48

15:                                               ; preds = %2
  %16 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  %17 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %18 = icmp eq %struct.PVHead* %16, %17
  br i1 %18, label %24, label %19

19:                                               ; preds = %15
  %20 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  %21 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %20, i32 0, i32 1
  %22 = load i32, i32* %21, align 8
  %23 = icmp eq i32 %22, 0
  br i1 %23, label %24, label %25

24:                                               ; preds = %19, %15
  store i8 1, i8* %3, align 1
  br label %48

25:                                               ; preds = %19
  store i32 0, i32* %6, align 4
  br label %26

26:                                               ; preds = %44, %25
  %27 = load i32, i32* %6, align 4
  %28 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  %29 = getelementptr inbounds %struct.PVHead, %struct.PVHead* %28, i32 0, i32 1
  %30 = load i32, i32* %29, align 8
  %31 = icmp ult i32 %27, %30
  br i1 %31, label %32, label %47

32:                                               ; preds = %26
  %33 = load i32, i32* %6, align 4
  %34 = load %struct.PVHead*, %struct.PVHead** %5, align 8
  %35 = call i16 @pv_uint16_get(%struct.PVHead* %34, i32 %33)
  %36 = zext i16 %35 to i32
  %37 = load i32, i32* %6, align 4
  %38 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %39 = call i16 @pv_uint16_get(%struct.PVHead* %38, i32 %37)
  %40 = zext i16 %39 to i32
  %41 = icmp ne i32 %36, %40
  br i1 %41, label %42, label %43

42:                                               ; preds = %32
  store i8 0, i8* %3, align 1
  br label %48

43:                                               ; preds = %32
  br label %44

44:                                               ; preds = %43
  %45 = load i32, i32* %6, align 4
  %46 = add i32 %45, 1
  store i32 %46, i32* %6, align 4
  br label %26, !llvm.loop !27

47:                                               ; preds = %26
  store i8 1, i8* %3, align 1
  br label %48

48:                                               ; preds = %47, %42, %24, %14
  %49 = load i8, i8* %3, align 1
  ret i8 %49
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local %struct.PVHead* @pv_uint16_append(%struct.PVHead* %0, i16 %1) #0 {
  %3 = alloca i16, align 2
  %4 = alloca %struct.PVHead*, align 8
  %5 = alloca %struct.PVLeaf_uint16*, align 8
  %6 = alloca %struct.PVHead*, align 8
  %7 = alloca %struct.PVHead*, align 8
  store i16 %1, i16* %3, align 2
  store %struct.PVHead* %0, %struct.PVHead** %4, align 8
  %8 = call %struct.PVH* @pl_new(i32 76)
  %9 = bitcast %struct.PVH* %8 to %struct.PVLeaf_uint16*
  store %struct.PVLeaf_uint16* %9, %struct.PVLeaf_uint16** %5, align 8
  %10 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %5, align 8
  %11 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %10, i32 0, i32 0
  %12 = getelementptr inbounds %struct.PVH, %struct.PVH* %11, i32 0, i32 2
  store i32 1, i32* %12, align 4
  %13 = load i16, i16* %3, align 2
  %14 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %5, align 8
  %15 = getelementptr inbounds %struct.PVLeaf_uint16, %struct.PVLeaf_uint16* %14, i32 0, i32 1
  %16 = getelementptr inbounds [32 x i16], [32 x i16]* %15, i64 0, i64 0
  store i16 %13, i16* %16, align 4
  %17 = load %struct.PVLeaf_uint16*, %struct.PVLeaf_uint16** %5, align 8
  %18 = bitcast %struct.PVLeaf_uint16* %17 to %struct.PVH*
  %19 = call %struct.PVHead* @pv_construct(%struct.PVH* %18)
  store %struct.PVHead* %19, %struct.PVHead** %6, align 8
  %20 = load %struct.PVHead*, %struct.PVHead** %6, align 8
  %21 = load %struct.PVHead*, %struct.PVHead** %4, align 8
  %22 = call %struct.PVHead* @pv_concatenate(%struct.PVHead* %21, %struct.PVHead* %20)
  store %struct.PVHead* %22, %struct.PVHead** %7, align 8
  %23 = load %struct.PVHead*, %struct.PVHead** %6, align 8
  call void @pv_free(%struct.PVHead* %23)
  %24 = load %struct.PVHead*, %struct.PVHead** %7, align 8
  ret %struct.PVHead* %24
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @free_structure(%struct.Structure* %0) #0 {
  %2 = alloca %struct.Structure*, align 8
  %3 = alloca i32, align 4
  %4 = alloca i16, align 2
  %5 = alloca i8*, align 8
  %6 = alloca i8**, align 8
  store %struct.Structure* %0, %struct.Structure** %2, align 8
  %7 = load %struct.Structure*, %struct.Structure** %2, align 8
  %8 = icmp eq %struct.Structure* %7, null
  br i1 %8, label %9, label %10

9:                                                ; preds = %1
  br label %53

10:                                               ; preds = %1
  %11 = load %struct.Structure*, %struct.Structure** %2, align 8
  %12 = getelementptr inbounds %struct.Structure, %struct.Structure* %11, i32 0, i32 0
  %13 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %12, i32 0, i32 1
  %14 = load i32, i32* %13, align 4
  store i32 %14, i32* %3, align 4
  %15 = load i32, i32* %3, align 4
  %16 = icmp eq i32 %15, 1
  br i1 %16, label %17, label %43

17:                                               ; preds = %10
  %18 = load %struct.Structure*, %struct.Structure** %2, align 8
  %19 = getelementptr inbounds %struct.Structure, %struct.Structure* %18, i32 0, i32 1
  %20 = load i16, i16* %19, align 4
  store i16 %20, i16* %4, align 2
  %21 = load %struct.Structure*, %struct.Structure** %2, align 8
  %22 = getelementptr inbounds %struct.Structure, %struct.Structure* %21, i64 1
  %23 = bitcast %struct.Structure* %22 to i8*
  store i8* %23, i8** %5, align 8
  %24 = load i8*, i8** %5, align 8
  %25 = getelementptr inbounds i8, i8* %24, i64 4
  store i8* %25, i8** %5, align 8
  %26 = load i8*, i8** %5, align 8
  %27 = bitcast i8* %26 to i8**
  store i8** %27, i8*** %6, align 8
  br label %28

28:                                               ; preds = %32, %17
  %29 = load i16, i16* %4, align 2
  %30 = zext i16 %29 to i32
  %31 = icmp sgt i32 %30, 0
  br i1 %31, label %32, label %40

32:                                               ; preds = %28
  %33 = load i8**, i8*** %6, align 8
  %34 = load i8*, i8** %33, align 8
  %35 = bitcast i8* %34 to %struct.RefCount*
  call void @free_rc(%struct.RefCount* %35)
  %36 = load i8**, i8*** %6, align 8
  %37 = getelementptr inbounds i8*, i8** %36, i32 1
  store i8** %37, i8*** %6, align 8
  %38 = load i16, i16* %4, align 2
  %39 = add i16 %38, -1
  store i16 %39, i16* %4, align 2
  br label %28, !llvm.loop !28

40:                                               ; preds = %28
  %41 = load %struct.Structure*, %struct.Structure** %2, align 8
  %42 = bitcast %struct.Structure* %41 to i8*
  call void @free(i8* %42)
  br label %53

43:                                               ; preds = %10
  %44 = load i32, i32* %3, align 4
  %45 = icmp ugt i32 %44, 1
  br i1 %45, label %46, label %52

46:                                               ; preds = %43
  %47 = load i32, i32* %3, align 4
  %48 = sub i32 %47, 1
  %49 = load %struct.Structure*, %struct.Structure** %2, align 8
  %50 = getelementptr inbounds %struct.Structure, %struct.Structure* %49, i32 0, i32 0
  %51 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %50, i32 0, i32 1
  store i32 %48, i32* %51, align 4
  br label %52

52:                                               ; preds = %46, %43
  br label %53

53:                                               ; preds = %52, %40, %9
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @free_closure(%struct.Closure* %0) #0 {
  %2 = alloca %struct.Closure*, align 8
  %3 = alloca i32, align 4
  %4 = alloca i16, align 2
  %5 = alloca i8*, align 8
  %6 = alloca i8**, align 8
  store %struct.Closure* %0, %struct.Closure** %2, align 8
  %7 = load %struct.Closure*, %struct.Closure** %2, align 8
  %8 = icmp eq %struct.Closure* %7, null
  br i1 %8, label %9, label %10

9:                                                ; preds = %1
  br label %51

10:                                               ; preds = %1
  %11 = load %struct.Closure*, %struct.Closure** %2, align 8
  %12 = getelementptr inbounds %struct.Closure, %struct.Closure* %11, i32 0, i32 0
  %13 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %12, i32 0, i32 1
  %14 = load i32, i32* %13, align 4
  store i32 %14, i32* %3, align 4
  %15 = load i32, i32* %3, align 4
  %16 = icmp eq i32 %15, 1
  br i1 %16, label %17, label %41

17:                                               ; preds = %10
  %18 = load %struct.Closure*, %struct.Closure** %2, align 8
  %19 = getelementptr inbounds %struct.Closure, %struct.Closure* %18, i32 0, i32 2
  %20 = load i16, i16* %19, align 8
  store i16 %20, i16* %4, align 2
  %21 = load %struct.Closure*, %struct.Closure** %2, align 8
  %22 = getelementptr inbounds %struct.Closure, %struct.Closure* %21, i64 1
  %23 = bitcast %struct.Closure* %22 to i8*
  store i8* %23, i8** %5, align 8
  %24 = load i8*, i8** %5, align 8
  %25 = bitcast i8* %24 to i8**
  store i8** %25, i8*** %6, align 8
  br label %26

26:                                               ; preds = %30, %17
  %27 = load i16, i16* %4, align 2
  %28 = zext i16 %27 to i32
  %29 = icmp sgt i32 %28, 0
  br i1 %29, label %30, label %38

30:                                               ; preds = %26
  %31 = load i8**, i8*** %6, align 8
  %32 = load i8*, i8** %31, align 8
  %33 = bitcast i8* %32 to %struct.RefCount*
  call void @free_rc(%struct.RefCount* %33)
  %34 = load i8**, i8*** %6, align 8
  %35 = getelementptr inbounds i8*, i8** %34, i32 1
  store i8** %35, i8*** %6, align 8
  %36 = load i16, i16* %4, align 2
  %37 = add i16 %36, -1
  store i16 %37, i16* %4, align 2
  br label %26, !llvm.loop !29

38:                                               ; preds = %26
  %39 = load %struct.Closure*, %struct.Closure** %2, align 8
  %40 = bitcast %struct.Closure* %39 to i8*
  call void @free(i8* %40)
  br label %51

41:                                               ; preds = %10
  %42 = load i32, i32* %3, align 4
  %43 = icmp ugt i32 %42, 1
  br i1 %43, label %44, label %50

44:                                               ; preds = %41
  %45 = load i32, i32* %3, align 4
  %46 = sub i32 %45, 1
  %47 = load %struct.Closure*, %struct.Closure** %2, align 8
  %48 = getelementptr inbounds %struct.Closure, %struct.Closure* %47, i32 0, i32 0
  %49 = getelementptr inbounds %struct.RefCount, %struct.RefCount* %48, i32 0, i32 1
  store i32 %46, i32* %49, align 4
  br label %50

50:                                               ; preds = %44, %41
  br label %51

51:                                               ; preds = %50, %38, %9
  ret void
}

attributes #0 = { noinline nounwind optnone uwtable "frame-pointer"="none" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { "frame-pointer"="none" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { noreturn "frame-pointer"="none" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { nofree nosync nounwind willreturn }
attributes #4 = { argmemonly nofree nounwind willreturn }
attributes #5 = { noreturn }

!llvm.ident = !{!0, !0, !0, !0}
!llvm.module.flags = !{!1, !2, !3}

!0 = !{!"clang version 14.0.0 (https://github.com/llvm/llvm-project.git 52a4a4a53c3ebffe474802dc87cd61a38e1783b5)"}
!1 = !{i32 1, !"wchar_size", i32 2}
!2 = !{i32 7, !"PIC Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 1}
!4 = distinct !{!4, !5}
!5 = !{!"llvm.loop.mustprogress"}
!6 = distinct !{!6, !5}
!7 = distinct !{!7, !5}
!8 = distinct !{!8, !5}
!9 = distinct !{!9, !5}
!10 = distinct !{!10, !5}
!11 = distinct !{!11, !5}
!12 = distinct !{!12, !5}
!13 = distinct !{!13, !5}
!14 = distinct !{!14, !5}
!15 = distinct !{!15, !5}
!16 = distinct !{!16, !5}
!17 = distinct !{!17, !5}
!18 = distinct !{!18, !5}
!19 = distinct !{!19, !5}
!20 = distinct !{!20, !5}
!21 = distinct !{!21, !5}
!22 = distinct !{!22, !5}
!23 = distinct !{!23, !5}
!24 = distinct !{!24, !5}
!25 = distinct !{!25, !5}
!26 = distinct !{!26, !5}
!27 = distinct !{!27, !5}
!28 = distinct !{!28, !5}
!29 = distinct !{!29, !5}
