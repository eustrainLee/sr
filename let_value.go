// sender factory: LetValue
package sr

import "context"

type letValueSender[T any, R any] struct {
	s Sender[T]
	f func(context.Context, T, Receiver[R])
}

func LetValue[T any, R any](s Sender[T], f func(context.Context, T, Receiver[R])) Sender[R] {
	return letValueSender[T, R]{f: f, s: s}
}

func (s letValueSender[T, R]) Connect(r Receiver[R]) OperationState {
	return letValueSenderState[T, R]{s: s.s, f: s.f, r: r}
}

func (s letValueSender[T, R]) Tag() SenderTag {
	return s.s.Tag()
}

type letValueSenderState[T any, R any] struct {
	s Sender[T]
	f func(context.Context, T, Receiver[R])
	r Receiver[R]
}

func (state letValueSenderState[T, R]) Start(ctx context.Context) {
	state.s.Connect(letValueReceiver[T, R]{ctx: ctx, f: state.f, r: state.r}).Start(ctx)
}

type letValueReceiver[T any, R any] struct {
	ctx context.Context
	f   func(context.Context, T, Receiver[R])
	r   Receiver[R]
}

func (lvr letValueReceiver[T, R]) SetValue(v T) {
	lvr.f(lvr.ctx, v, lvr.r)
}

func (lvr letValueReceiver[T, R]) SetError(err error) {
	lvr.r.SetError(err)
}

func (lvr letValueReceiver[T, R]) SetStoped() {
	lvr.r.SetStoped()
}
