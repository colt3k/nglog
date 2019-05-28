package ng

type TriggerPolicy interface {
	Rotate(string) bool
}

// tests using cron
type CronTriggerPolicy struct {
	schedule    string //	i.e. 0 0 * * * ?
	evalOnStart bool
}

// checks while running or on start
type SizeTriggerPolicy struct {
	maxSize     int64
	evalOnStart bool
}

// tests on interval of hours with a possible max random delay of X
type TimeTriggerPolicy struct {
	interval       int // based on most specific date pattern field i.e. hours (default)
	maxRandomDelay int // 0 no delay otherwise seconds to delay after interval
}

func NewSizeTriggerPolicy(maxsize int64, evalOnStart bool) *SizeTriggerPolicy {
	t := new(SizeTriggerPolicy)
	t.maxSize = maxsize
	t.evalOnStart = evalOnStart
	return t
}

func (s *SizeTriggerPolicy) Rotate(fileName string) bool {
	// TODO check rotation triggers
	return false
}

func NewTimeTriggerPolicy(interval, maxRandomDelay int) *TimeTriggerPolicy {
	t := new(TimeTriggerPolicy)
	t.interval = interval
	t.maxRandomDelay = maxRandomDelay
	return t
}

func (s *TimeTriggerPolicy) Rotate(fileName string) bool {
	// TODO check rotation triggers
	return false
}

func NewCronTriggerPolicy(schedule string, evalOnStart bool) *CronTriggerPolicy {
	t := new(CronTriggerPolicy)
	t.schedule = schedule
	t.evalOnStart = evalOnStart
	return t
}

func (s *CronTriggerPolicy) Rotate(fileName string) bool {
	// TODO check rotation triggers
	return false
}
