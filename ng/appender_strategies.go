package ng

// Archive strategy
type Strategy interface {
	Max() int
	CompressionLevel() int
}
type DefaultStrategy struct {
	max              int    // Maximum to keep
	compressionLevel int    // 0 thru 9 for zip only
}

func NewDefaultStrategy(maxKeep, compressionLevel int) *DefaultStrategy {
	t := new(DefaultStrategy)
	t.max = maxKeep
	t.compressionLevel = compressionLevel
	return t
}

func (s *DefaultStrategy) Max() int {
	return s.max
}
func (s *DefaultStrategy) CompressionLevel() int {
	return s.compressionLevel
}
