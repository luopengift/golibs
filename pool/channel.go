package pool

type Channel struct {
	channel chan bool
}

func NewChannel(max int) *Channel {
	ch := new(Channel)
	ch.channel = make(chan bool, max)
	return ch
}

func (ch *Channel) Add(delta int) {
	for i := 0; i < delta; i++ {
		ch.channel <- true
	}
}

func (ch *Channel) Done() {
	<-ch.channel
}

func (ch *Channel) Size() int {
	return cap(ch.channel)
}

func (ch *Channel) Len() int {
	return len(ch.channel)
}

func (ch *Channel) Close() error {
	close(ch.channel)
	return nil
}
