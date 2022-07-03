package system

type Observer interface {
	Update(data interface{}) error
}
