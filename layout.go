package nglog

type Layout interface {
	Format(Msg, bool) ([]byte, error)
	Description() string
	Colors(bool)
	DisableTimeStamp()
	EnableTimeStamp()
}
