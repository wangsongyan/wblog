package helpers

import (
	"strings"
	"time"

	"wblog/models"
)

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

// 截取字符串
func Substring(source string, start, end int) string {
	rs := []rune(source)
	length := len(rs)
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	return string(rs[start:end])
}

// 判断数字是否是奇数
func IsOdd(number int) bool {
	return !IsEven(number)
}

// 判断数字是否是偶数
func IsEven(number int) bool {
	return number%2 == 0
}

func Add(a1, a2 int) int {
	return a1 + a2
}

func Minus(a1, a2 int) int {
	return a1 - a2
}

func ListTag() (tagstr string) {
	tags, err := models.ListTag()
	if err != nil {
		return
	}
	tagNames := make([]string, 0)
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	tagstr = strings.Join(tagNames, ",")
	return
}
