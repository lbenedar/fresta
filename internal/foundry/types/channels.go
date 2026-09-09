package types

import "sync"

type ExchangeChannels struct {
	Msgs          map[int](chan []byte)
	ProgressMsg   map[string](chan struct{})
	ProgressMutex sync.Mutex
}

func CloseMapChannel[K comparable, V any](channel map[K]chan V) {
	for _, v := range channel {
		select {
		case _, ok := <-v:
			if ok {
				close(v)
			}
		default:
			close(v)
		}
	}
	for k := range channel {
		delete(channel, k)
	}
}

func (channels *ExchangeChannels) Close() {
	CloseMapChannel(channels.Msgs)
	CloseMapChannel(channels.ProgressMsg)
}

type ReadChannels struct {
	done      chan struct{}
	reconnect chan struct{}
	msg       chan *WsMessage
	err       chan error
}

func (channels ReadChannels) Err() chan error {
	return channels.err
}

func (channels ReadChannels) Msg() chan *WsMessage {
	return channels.msg
}

func (channels ReadChannels) Done() chan struct{} {
	return channels.done
}

func (channels ReadChannels) Reconnect() chan struct{} {
	return channels.reconnect
}

func CloseChannel[V any](channel chan V) {
	select {
	case _, ok := <-channel:
		if ok {
			close(channel)
		}
	default:
		close(channel)
	}
}

func (channels *ReadChannels) Close() {
	CloseChannel(channels.done)
	CloseChannel(channels.err)
	CloseChannel(channels.msg)
	CloseChannel(channels.reconnect)
}

func InitWsChannels() *ReadChannels {
	return &ReadChannels{
		done:      make(chan struct{}),
		err:       make(chan error, 10),
		msg:       make(chan *WsMessage, 10),
		reconnect: make(chan struct{}),
	}
}
