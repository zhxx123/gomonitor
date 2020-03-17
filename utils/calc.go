package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// 字符串乘法
func FloatMutiplyStr(num1, num2 string) string {
	v2, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		fmt.Printf("FloatMutiplyStr error num1: %s\n", num1)
		return ""
	}
	v3, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		fmt.Printf("FloatMutiplyStr error num1: %s\n", num2)
		return ""
	}
	res := fmt.Sprintf("%.3f", v2*v3)
	return res
}

// 字符串乘法
func FloatMutiplyStrPoint(num1, num2 string) string {
	v2, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		fmt.Printf("FloatMutiplyStr error num1: %s\n", num1)
		return ""
	}
	v3, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		fmt.Printf("FloatMutiplyStr error num1: %s\n", num2)
		return ""
	}
	res := fmt.Sprintf("%.2f", v2*v3)
	return res
}

// PayStringAdd 字符串加
func PayStringAdd(a, b string, pn int) string {
	if pn != 2 && pn != 8 {
		pn = 2
	}
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	// a 表示当前账户余额
	// b 表示 增加的余额
	anum, erra := strconv.ParseFloat(a, 64)
	if erra != nil {
		return "0"
	}
	bnum, errb := strconv.ParseFloat(b, 64)
	if errb != nil {
		return a
	}
	return strconv.FormatFloat(anum+bnum, 'f', pn, 64)
}

// 字符串减
func PayStringSub(a, b string, pn int) (string, error) {
	if pn != 2 && pn != 8 {
		pn = 2
	}
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	// a 表示 当前账户余额
	// b 表示 减少的余额
	anum, erra := strconv.ParseFloat(a, 64)
	if erra != nil {
		return "0", erra
	}
	bnum, errb := strconv.ParseFloat(b, 64)
	if errb != nil {
		return "0", errb
	}
	if anum < bnum {
		return "0", errors.New("not enough account")
	}
	return strconv.FormatFloat(anum-bnum, 'f', pn, 64), nil
}

// string 转换为 浮点数
func ParseStrToFloat(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

// string 转换 int
func ParseStrToInt(b string, defInt int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return defInt
	}
	return id
}

// string 转换 int64
func ParseStrToInt64(value string) (int64, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// 浮点数 转为 string
func ParseFloatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', 8, 64)
}

// 浮点数 转换为 string
func ParseFloatToStrWithPoint(f float64, pn int) string {
	if pn < 0 || pn > 8 {
		pn = 0
	}
	return strconv.FormatFloat(f, 'f', pn, 64)
}

// map string 转结构体对象
func MapStrToStruct(m map[string]string, i interface{}) error {
	bin, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, i)
	if err != nil {
		return err
	}
	return nil
}

//map interface{} 转结构体对象
func MapToStruct(m map[string]interface{}, i interface{}) error {
	bin, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, i)
	if err != nil {
		return err
	}
	return nil
}

// 字符串转结构体
func StrToStruct(bin string, data interface{}) error {
	if err := json.Unmarshal([]byte(bin), data); err != nil {
		return err
	}
	return nil
}

// 结构体 转成 json字符串
func StructToStr(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

// float64 匿名函数比较
type Accuracy func() float64

func (this Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < this()
}

func (this Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this()
}

func (this Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && (math.Abs(a-b) > this() || math.Abs(a) < this())
}

func (this Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this()
}

func (this Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this()
}

//****************** etherscan

func BigIntDivToString(value string, pn int) string {
	lv := len(value)
	newValue := "0"
	zero := ""
	if pn <= 0 {
		return value
	}
	if lv > pn {
		newValue = fmt.Sprintf("%s.%s", value[:lv-pn], value[lv-pn:])
		return newValue
	}
	for i := 0; i < pn-lv; i++ {
		zero = zero + "0"
	}
	return fmt.Sprintf("0.%s%s", zero, value)
}
