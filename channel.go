package Channel

type Channel struct {
	size      int
	receivers []chan interface{}
}

