package Channel

type Channel[T interface{}] struct {
	size      int
	receivers []chan T
}

func (channel *Channel[T]) Sender() chan<- T {
	c := make(chan T, channel.size)

	go func() {
		for {
			select {
			case data := <-c:
				for _, p := range channel.receivers {
					p <- data
				}
			}
		}
	}()

	return c
}

func (channel *Channel[T]) Receiver() <-chan T {
	c := make(chan T, channel.size)

	channel.receivers = append(channel.receivers, c)

	return c
}

func New[T interface{}](size int) *Channel[T] {
	return &Channel[T]{
		size: size,
	}
}
