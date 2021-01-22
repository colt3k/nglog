package ng

import "compress/flate"

// Archive strategy
type Strategy interface {
	Max() int
	CompressionLevel() int
}
type FileStrategy struct {
	max              int    // Maximum to keep
	compressionLevel int    // 0 thru 9 for zip only
}

func DefaultFileStrategy() *FileStrategy {
	t := new(FileStrategy)
	t.max = 4
	t.compressionLevel = flate.BestCompression
	return t
}
func NewDefaultFileStrategy(maxKeep, compressionLevel int) *FileStrategy {
	t := new(FileStrategy)
	t.max = maxKeep
	t.compressionLevel = compressionLevel
	return t
}

func (s *FileStrategy) Max() int {
	return s.max
}
func (s *FileStrategy) CompressionLevel() int {
	return s.compressionLevel
}
