package src

import (
	"context"

	"github.com/eustrainLee/execution/sr"
	"github.com/samber/lo"
)

type whenAll5Sender[T1, T2, T3, T4, T5 any] struct {
	s1 Sender[T1]
	s2 Sender[T2]
	s3 Sender[T3]
	s4 Sender[T4]
	s5 Sender[T5]
}

func WhenAll5[T1, T2, T3, T4, T5 any](s1 Sender[T1], s2 Sender[T2], s3 Sender[T3], s4 Sender[T4], s5 Sender[T5]) Sender[lo.Tuple5[T1, T2, T3, T4, T5]] {
	return whenAll5Sender[T1, T2, T3, T4, T5]{s1: s1, s2: s2, s3: s3, s4: s4, s5: s5}
}

func (s whenAll5Sender[T1, T2, T3, T4, T5]) Tag() sr.SenderTag {
	return sr.SenderTagNone
}

func (s whenAll5Sender[T1, T2, T3, T4, T5]) Connect(r sr.Receiver[lo.Tuple5[T1, T2, T3, T4, T5]]) OperationState {
	return whenAll5OperationState[T1, T2, T3, T4, T5]{s: s, r: r}
}

type whenAll5OperationState[T1, T2, T3, T4, T5 any] struct {
	s whenAll5Sender[T1, T2, T3, T4, T5]
	r sr.Receiver[lo.Tuple5[T1, T2, T3, T4, T5]]
}

func (os whenAll5OperationState[T1, T2, T3, T4, T5]) Start(ctx context.Context) {
	const SenderCount = 5
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	result := lo.Tuple5[T1, T2, T3, T4, T5]{}
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
	v4Chan := make(chan T4, 1)
	go os.s.s4.Connect(sr.ChannelReceiver[T4]{
		ValueChan:  v4Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start(ctx)
	v5Chan := make(chan T5, 1)
	go os.s.s5.Connect(sr.ChannelReceiver[T5]{
		ValueChan:  v5Chan,
		ErrorChan:  errChan,
		StopedChan: stopedChan,
	}).Start(ctx)
	for i := 0; i < SenderCount; i++ {
		select {
		case result.A = <-v1Chan:
		case result.B = <-v2Chan:
		case result.C = <-v3Chan:
		case result.D = <-v4Chan:
		case result.E = <-v5Chan:
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