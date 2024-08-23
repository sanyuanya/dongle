package tools

import (
	"fmt"
	"strconv"
	"time"
)

func ValidateTimestamp(timestampStr string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析时间戳: %v", err)
	}

	// 检查时间戳是否在合理范围内（例如：1970年到现在）
	if timestamp < 0 || timestamp > time.Now().Unix()*1000 {
		return time.Time{}, fmt.Errorf("时间戳超出合理范围")
	}

	// 将时间戳转换为 time.Time 类型
	return time.Unix(0, timestamp*int64(time.Millisecond)), nil
}
