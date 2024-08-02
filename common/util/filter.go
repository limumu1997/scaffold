package util

import (
	"bytes"
)

// FilterArrayByRule 根据规则 rule 过滤源切片 src 中的元素，将不符合规则的元素过滤掉，返回过滤后的切片
// Deprecated: this function simply calls bytes.ReplaceAll
func FilterArrayByRule(src, rule []byte) []byte {
	i := 0
	for k, v := range src {
		if bytes.Contains(rule, src[k:k+1]) {
			continue
		}
		src[i] = v
		i++
	}
	return src[:i]
}

// EscapeBytes 字符串转译，src 是输入的字节切片，target 是要转义的目标字节序列（在此例中为 0xA6, 0x01），escape 是转义后的字节（在此例中为 0xA6）
// Deprecated: this function simply calls bytes.ReplaceAll
func EscapeBytes(src []byte, target []byte, escape byte) []byte {
	var escaped []byte
	for i := 0; i < len(src); i++ {
		if i < len(src)-1 && src[i] == target[0] && src[i+1] == target[1] {
			escaped = append(escaped, escape)
			i++ // 跳过下一个字节
		} else {
			escaped = append(escaped, src[i])
		}
	}
	return escaped
}

// CountByteSequence 用于统计一个 byte 数组中出现指定字节序列的次数
// Deprecated: this function simply calls bytes.Count
func CountByteSequence(arr []byte, seq []byte) int {
	count := 0
	for i := 0; i < len(arr)-len(seq)+1; i++ {
		if bytes.Equal(arr[i:i+len(seq)], seq) {
			count++
		}
	}
	return count
}

// FindIndex 找查 seq 在指定数组中的 Last Index
// Deprecated: this function simply calls bytes.Index + len(seq)
func FindIndex(arr []byte, seq []byte) int {
	for i := 0; i < len(arr)-len(seq)+1; i++ {
		if bytes.Equal(arr[i:i+len(seq)], seq) {
			return i + len(seq)
		}
	}
	return 0
}
