package main

// 本文件演示 Go 语言的基本数据类型

import (
	"fmt"
	"unsafe"
)

func main() {
	// ========== 整型 ==========
	// Go 整型分为有符号和无符号两大类

	// 1. 有符号整型（可以表示正负数）
	var signedInt8 int8 = 127                   // 8位，范围: -128 ~ 127
	var signedInt16 int16 = 32767               // 16位，范围: -32768 ~ 32767
	var signedInt32 int32 = 2147483647          // 32位，范围: -2147483648 ~ 2147483647
	var signedInt64 int64 = 9223372036854775807 // 64位，范围非常大
	var signedInt int = 1000000                 // 32位或64位，取决于操作系统

	fmt.Printf("有符号整型:\n")
	fmt.Printf("  int8: %d, 大小: %d 字节\n", signedInt8, unsafe.Sizeof(signedInt8))
	fmt.Printf("  int16: %d, 大小: %d 字节\n", signedInt16, unsafe.Sizeof(signedInt16))
	fmt.Printf("  int32: %d, 大小: %d 字节\n", signedInt32, unsafe.Sizeof(signedInt32))
	fmt.Printf("  int64: %d, 大小: %d 字节\n", signedInt64, unsafe.Sizeof(signedInt64))
	fmt.Printf("  int: %d, 大小: %d 字节\n", signedInt, unsafe.Sizeof(signedInt))

	// 2. 无符号整型（只能表示非负数）
	var unsignedUint8 uint8 = 255                    // 8位，范围: 0 ~ 255
	var unsignedUint16 uint16 = 65535                // 16位，范围: 0 ~ 65535
	var unsignedUint32 uint32 = 4294967295           // 32位，范围: 0 ~ 4294967295
	var unsignedUint64 uint64 = 18446744073709551615 // 64位，非常大
	var unsignedUint uint = 1000000                  // 32位或64位

	fmt.Printf("无符号整型:\n")
	fmt.Printf("  uint8: %d, 大小: %d 字节\n", unsignedUint8, unsafe.Sizeof(unsignedUint8))
	fmt.Printf("  uint16: %d, 大小: %d 字节\n", unsignedUint16, unsafe.Sizeof(unsignedUint16))
	fmt.Printf("  uint32: %d, 大小: %d 字节\n", unsignedUint32, unsafe.Sizeof(unsignedUint32))
	fmt.Printf("  uint64: %d, 大小: %d 字节\n", unsignedUint64, unsafe.Sizeof(unsignedUint64))
	fmt.Printf("  uint: %d, 大小: %d 字节\n", unsignedUint, unsafe.Sizeof(unsignedUint))

	// 3. 特殊整型
	var byteVar byte = 255        // byte 是 uint8 的别名
	var runeVar rune = 2147483647 // rune 是 int32 的别名，用于表示 Unicode 码点
	var uintptrVar uintptr = 1000 // 用于存储指针的无符号整数

	fmt.Printf("特殊整型:\n")
	fmt.Printf("  byte: %d, 大小: %d 字节\n", byteVar, unsafe.Sizeof(byteVar))
	fmt.Printf("  rune: %d, 大小: %d 字节\n", runeVar, unsafe.Sizeof(runeVar))
	fmt.Printf("  uintptr: %d, 大小: %d 字节\n", uintptrVar, unsafe.Sizeof(uintptrVar))

	// ========== 浮点型 ==========
	// 用于表示小数

	var float32Var float32 = 3.1415926         // 32位浮点数，精度约6-7位
	var float64Var float64 = 3.141592653589793 // 64位浮点数，精度约15位

	fmt.Printf("\n浮点型:\n")
	fmt.Printf("  float32: %f, 大小: %d 字节\n", float32Var, unsafe.Sizeof(float32Var))
	fmt.Printf("  float64: %f, 大小: %d 字节\n", float64Var, unsafe.Sizeof(float64Var))

	// 浮点数精度问题
	a := 0.1
	b := 0.2
	c := a + b
	fmt.Printf("  精度问题: %.20f\n", c) // 输出: 0.30000000000000004441

	// ========== 复数类型 ==========
	// 用于表示复数（数学中的复数 a + bi）

	var complex64Var complex64 = 3 + 4i   // 实部和虚部都是 float32
	var complex128Var complex128 = 3 + 4i // 实部和虚部都是 float64

	fmt.Printf("\n复数类型:\n")
	fmt.Printf("  complex64: %v, 大小: %d 字节\n", complex64Var, unsafe.Sizeof(complex64Var))
	fmt.Printf("  complex128: %v, 大小: %d 字节\n", complex128Var, unsafe.Sizeof(complex128Var))

	// 获取实部和虚部
	fmt.Printf("  complex128 实部: %f, 虚部: %f\n", real(complex128Var), imag(complex128Var))

	// ========== 布尔型 ==========
	// 只有两个值：true 和 false

	var boolVar1 bool = true
	var boolVar2 bool = false

	fmt.Printf("\n布尔型:\n")
	fmt.Printf("  boolVar1: %v, 大小: %d 字节\n", boolVar1, unsafe.Sizeof(boolVar1))
	fmt.Printf("  boolVar2: %v, 大小: %d 字节\n", boolVar2, unsafe.Sizeof(boolVar2))

	// 布尔运算
	fmt.Printf("  true && false = %v (逻辑与)\n", true && false)
	fmt.Printf("  true || false = %v (逻辑或)\n", true || false)
	fmt.Printf("  !true = %v (逻辑非)\n", !true)

	// ========== 字符串类型 ==========
	// 字符串是不可变的字节序列

	var str1 string = "Hello, Go!" // 双引号声明
	var str2 string = `多行字符串
可以使用换行` // 反引号声明，保留原始格式

	// 字符串长度（字节数）
	fmt.Printf("\n字符串:\n")
	fmt.Printf("  str1: %s, 长度: %d 字节\n", str1, len(str1))
	fmt.Printf("  str2: %s\n", str2)

	// 字符串是 UTF-8 编码
	chineseStr := "你好"
	fmt.Printf("  中文字符串: %s, 字节数: %d, 字符数: %d\n",
		chineseStr, len(chineseStr), len([]rune(chineseStr)))

	// 字符串操作
	concatenated := str1 + " " + "World" // 拼接
	fmt.Printf("  拼接: %s\n", concatenated)

	// 字符串切片
	substr := str1[0:5] // 0 <= index < 5
	fmt.Printf("  切片 [0:5]: %s\n", substr)

	// ========== 零值 ==========
	// 每种类型的默认值（未初始化时的值）

	var (
		zeroInt    int
		zeroFloat  float64
		zeroBool   bool
		zeroString string
		zeroPtr    *int
		zeroSlice  []int
		zeroMap    map[string]int
	)

	fmt.Printf("\n零值:\n")
	fmt.Printf("  int: %d\n", zeroInt)
	fmt.Printf("  float64: %f\n", zeroFloat)
	fmt.Printf("  bool: %v\n", zeroBool)
	fmt.Printf("  string: '%s'\n", zeroString)
	fmt.Printf("  pointer: %v\n", zeroPtr)
	fmt.Printf("  slice: %v\n", zeroSlice)
	fmt.Printf("  map: %v\n", zeroMap)

	// ========== 类型别名 ==========
	// 为已有类型定义别名

	type (
		MyInt    int    // MyInt 是 int 的别名
		MyString string // MyString 是 string 的别名
	)

	var num MyInt = 100
	var name MyString = "别名"
	fmt.Printf("\n类型别名:\n")
	fmt.Printf("  MyInt: %d, 实际类型: %T\n", num, num)
	fmt.Printf("  MyString: %s, 实际类型: %T\n", name, name)
}

// ========== 总结 ==========
// 1. 整型：int8/int16/int32/int64, uint8/uint16/uint32/uint64, byte, rune
// 2. 浮点型：float32, float64
// 3. 复数：complex64, complex128
// 4. 布尔型：true, false
// 5. 字符串：不可变的 UTF-8 字节序列
// 6. 零值：每种类型的默认值
