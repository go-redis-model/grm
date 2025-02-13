package grm

import "fmt"

// 定义复合错误类型，包含具体错误信息
type PartialError struct {
	Errors map[string]error // Key → 错误原因
}

func (e *PartialError) Error() string {
	return fmt.Sprintf("partial error (%d failures)", len(e.Errors))
}
