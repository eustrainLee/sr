package src

import (
	"context"

	"github.com/eustrainLee/execution/sr"
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

func (s whenAll3Sender[T1, T2, T3]) Tag() sr.SenderTag {
	return sr.SenderTagNone
}

func (s whenAll3Sender[T1, T2, T3]) Connect(r sr.Receiver[lo.Tuple3[T1, T2, T3]]) OperationState {
	return whenAll3OperationState[T1, T2, T3]{s: s, r: r}
}

type whenAll3OperationState[T1, T2, T3 any] struct {
	s whenAll3Sender[T1, T2, T3]
	r sr.Receiver[lo.Tuple3[T1, T2, T3]]
}

func (os whenAll3OperationState[T1, T2, T3]) Start(ctx context.Context) {
	const SenderCount = 3
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	result := lo.Tuple3[T1, T2, T3]{}
	errChan := make(chan error)
	stopedChan := make(chan struct{}, SenderCount)
	v1Chan := make(chan T1, 1)
	go os.s.s1.Connect(sr.ChannelReceiver[T1]{
		ValueChan:  v1Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start(ctx)
	v2Chan := make(chan T2, 1)
	go os.s.s2.Connect(sr.ChannelReceiver[T2]{
		ValueChan:  v2Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start(ctx)
	v3Chan := make(chan T3, 1)
	go os.s.s3.Connect(sr.ChannelReceiver[T3]{
		ValueChan:  v3Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start(ctx)
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
