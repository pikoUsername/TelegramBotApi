package dispatcher

type BaseMiddleware interface {
	func callback()
}
