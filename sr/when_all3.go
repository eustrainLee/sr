package sr

import (
	"github.com/samber/lo"
)

type whenAll3Sender[T1, T2, T3 any] struct {
	s1 Sender[T1]
	s2 Sender[T2]
	s3 Sender[T3]
}

func WhenAll3[T1, T2, T3 any](s1 Sender[T1], s2 Sender[T2], s3 Sender[T3]) Sender[lo.Tuple3[T1, T2, T3]] {
	return whenAll3Sender[T1, T2, T3]{s1: s1, s2: s2, s3: s3}
}

func (s whenAll3Sender[T1, T2, T3]) Tag() SenderTag {
	return SenderTagNone
}

func (s whenAll3Sender[T1, T2, T3]) Connect(r Receiver[lo.Tuple3[T1, T2, T3]]) OperationState {
	return whenAll3OperationState[T1, T2, T3]{s: s, r: r}
}

type whenAll3OperationState[T1, T2, T3 any] struct {
	s whenAll3Sender[T1, T2, T3]
	r Receiver[lo.Tuple3[T1, T2, T3]]
}

func (os whenAll3OperationState[T1, T2, T3]) Start() {
	const SenderCount = 3
	result := lo.Tuple3[T1, T2, T3]{}
	errChan := make(chan error)
	stopedChan := make(chan struct{}, SenderCount)
	v1Chan := make(chan T1, 1)
	go os.s.s1.Connect(ChannelReceiver[T1]{
		ValueChan:  v1Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start()
	v2Chan := make(chan T2, 1)
	go os.s.s2.Connect(ChannelReceiver[T2]{
		ValueChan:  v2Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start()
	v3Chan := make(chan T3, 1)
	go os.s.s3.Connect(ChannelReceiver[T3]{
		ValueChan:  v3Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start()
	for i := 0; i < SenderCount; i++ {
		select {
		case result.A = <-v1Chan:
		case result.B = <-v2Chan:
		case result.C = <-v3Chan:
		case err := <-errChan:
			os.r.SetError(err)
			return
		case <-stopedChan:
			os.r.SetStoped()
			return
		}
	}
	os.r.SetValue(result)
}
