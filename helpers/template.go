package helpers

import (
	"github.com/wangsongyan/wblog/models"
	"time"
)

func ListTag() []*models.Tag {
	tags, _ := models.ListTag()
	return tags
}

func ListArchive() []*models.QrArchive {
	archievs, _ := models.ListPostArchives()
	return archievs
}

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

// 截取字符串
func Substring(source string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(source) {
		end = len(source)
	}
	return source[start:end]
}

// 判断数字是否是奇数
func IsOdd(number int) bool {
	return !IsEven(number)
}

// 判断数字是否是偶数
func IsEven(number int) bool {
	return number%2 == 0
}

func IsActive(arg1, arg2 string) string {
	if arg1 == arg2 {
		return "active"
	}
	return ""
}
