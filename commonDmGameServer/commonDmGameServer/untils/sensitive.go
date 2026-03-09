package untils

import (
	"dmGameServer/zlog"
	"github.com/importcjj/sensitive"
)

var sensitiveFilter *SensitiveFilter

type SensitiveFilter struct {
	filter *sensitive.Filter
}

func GetSensitiveFilter() *SensitiveFilter {
	if sensitiveFilter == nil {
		sensitiveFilter, _ = NewSensitiveFilter("./mgc.txt")

	}
	return sensitiveFilter
}

// NewSensitiveFilter 创建一个新的敏感词过滤器实例
func NewSensitiveFilter(dictPath string) (*SensitiveFilter, error) {
	filter := sensitive.New()
	err := filter.LoadWordDict(dictPath)
	if err != nil {
		zlog.Logger.Error().Msgf("加载敏感词字典失败: %v", err)
		return nil, err
	}
	return &SensitiveFilter{filter: filter}, nil
}

// ContainsSensitiveWord 检查文本中是否包含敏感词
func (sf *SensitiveFilter) ContainsSensitiveWord(text string) bool {
	return sf.filter.Filter(text) != text
}

// ReplaceSensitiveWords 替换文本中的敏感词为指定符号
func (sf *SensitiveFilter) ReplaceSensitiveWords(text string, replaceWith rune) string {
	return sf.filter.Replace(text, replaceWith)
}
