package ng

// Archive strategy
type Strategy interface {
	Name() string     // Name of file
	Location() string // Location to store old files
	Min() int
	Max() int
	CompressionLevel() int
}
type DefaultStrategy struct {
	name             string // Name of file
	location         string // Location to store old files
	min              int    // Minimum to keep
	max              int    // Maximum to keep
	compressionLevel int    // 0 thru 9 for zip only
}

func NewDefaultStrategy(name, location string, min, max, compressionLevel int) *DefaultStrategy {
	t := new(DefaultStrategy)
	t.name = name
	t.location = location
	t.min = min
	t.max = max
	t.compressionLevel = compressionLevel
	return t
}
func (s *DefaultStrategy) Name() string {
	return s.name
}
func (s *DefaultStrategy) Location() string {
	return s.location
}
func (s *DefaultStrategy) Min() int {
	return s.min
}
func (s *DefaultStrategy) Max() int {
	return s.max
}
func (s *DefaultStrategy) CompressionLevel() int {
	return s.compressionLevel
}
