package dispatcher

type HandlerObj struct {
	Callable func(interface{})
}
