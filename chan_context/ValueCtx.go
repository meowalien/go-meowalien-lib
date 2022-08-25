package chan_context

import "reflect"

type valueCtx struct {
	WaitContext
	key, val any
}

func stringify(v any) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	}
	return "<not Stringer>"
}

func (c *valueCtx) String() string {
	return contextName(c.WaitContext) + ".WithValue(type " +
		reflect.TypeOf(c.key).String() +
		", val " + stringify(c.val) + ")"
}

func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return value(c.WaitContext, key)
}

func value(c WaitContext, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.WaitContext
		case *waitContext:
			if key == &cancelCtxKey {
				return c
			}
			c = ctx.WaitContext
		//case *timerCtx:
		//	if key == &cancelCtxKey {
		//		return &ctx.cancelCtx
		//	}
		//	c = ctx.WaitContext
		//case *emptyCtx:
		//	return nil
		default:
			return c.Value(key)
		}
	}
}
