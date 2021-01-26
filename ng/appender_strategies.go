package ng

import "compress/flate"

// Archive strategy
type Strategy interface {
	Max() int
	CompressionLevel() int
}
type FileStrategy struct {
	maxToKeep        int // Maximum to keep, can be overridden by maxDays
	maxDays          int // Maximum days to keep
	compressionLevel int // 0 thru 9 for zip only
}

func DefaultFileStrategy() *FileStrategy {
	t := new(FileStrategy)
	t.maxToKeep = 4
	t.compressionLevel = flate.BestCompression
	return t
}
func NewDefaultFileStrategy(maxKeep, compressionLevel int) *FileStrategy {
	t := new(FileStrategy)
	t.maxToKeep = maxKeep
	t.compressionLevel = compressionLevel
	return t
}

func (s *FileStrategy) Max() int {
	return s.maxToKeep
}

func (s *FileStrategy) CompressionLevel() int {
	return s.compressionLevel
}
