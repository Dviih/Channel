package Channel

type Channel struct {
	size      int
	receivers []chan interface{}
}

func (channel *Channel) Sender() chan<- interface{} {
	c := make(chan interface{}, channel.size)

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

func (channel *Channel) Receiver() <-chan interface{} {
	c := make(chan interface{}, channel.size)

	channel.receivers = append(channel.receivers, c)

	return c
}

