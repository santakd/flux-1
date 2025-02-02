// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: conditional.gen.go.tmpl

package values

import (
	"github.com/apache/arrow/go/v7/arrow/memory"
	fluxarray "github.com/influxdata/flux/array"
	"github.com/influxdata/flux/codes"
	"github.com/influxdata/flux/internal/errors"
	"github.com/influxdata/flux/semantic"
)

// VectorConditional delegates the branch selection and necessary array building
// to specialized (by input type) functions in the `array` package, then wraps
// the output as a Value.
//
// For cases where `t` is constant (either by happening to be all the same value,
// or by being a `VectorRepeatValue`) the caller should skip this function and
// return either the consequent or alternate
// directly.
// With constant `t` is expected to be handled higher up, this function will error
// if it receives a `VectorRepeatValue` for this parameter.
//
// Cases where either or both of the consequent `c` and alternate `a` are
// `VectorRepeatValue`s are supported.
func VectorConditional(t, c, a Vector, mem memory.Allocator) (Value, error) {
	var tvr, cvr, avr *Value
	if vr, ok := t.(*VectorRepeatValue); ok {
		tvr = &vr.val
	}
	if vr, ok := c.(*VectorRepeatValue); ok {
		cvr = &vr.val
	}
	if vr, ok := a.(*VectorRepeatValue); ok {
		avr = &vr.val
	}

	if tvr != nil {
		panic("t is constant, should be handled higher up, such as in conditionalVectorEvaluator.Eval")
	}

	switch a.ElementType().Nature() {

	case semantic.Int:
		var (
			x   *fluxarray.Int
			err error
		)
		if cvr != nil && avr != nil {
			x, err = fluxarray.IntConditionalCConstAConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Int(),

				(*avr).Int(),

				mem,
			)
		} else if cvr != nil {
			x, err = fluxarray.IntConditionalCConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Int(),

				a.Arr().(*fluxarray.Int),
				mem,
			)
		} else if avr != nil {
			x, err = fluxarray.IntConditionalAConst(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Int),

				(*avr).Int(),

				mem,
			)
		} else {
			x, err = fluxarray.IntConditional(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Int),
				a.Arr().(*fluxarray.Int),
				mem,
			)
		}
		if err != nil {
			return nil, err
		}
		return NewVectorValue(x, semantic.BasicInt), nil

	case semantic.UInt:
		var (
			x   *fluxarray.Uint
			err error
		)
		if cvr != nil && avr != nil {
			x, err = fluxarray.UintConditionalCConstAConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).UInt(),

				(*avr).UInt(),

				mem,
			)
		} else if cvr != nil {
			x, err = fluxarray.UintConditionalCConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).UInt(),

				a.Arr().(*fluxarray.Uint),
				mem,
			)
		} else if avr != nil {
			x, err = fluxarray.UintConditionalAConst(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Uint),

				(*avr).UInt(),

				mem,
			)
		} else {
			x, err = fluxarray.UintConditional(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Uint),
				a.Arr().(*fluxarray.Uint),
				mem,
			)
		}
		if err != nil {
			return nil, err
		}
		return NewVectorValue(x, semantic.BasicUint), nil

	case semantic.Float:
		var (
			x   *fluxarray.Float
			err error
		)
		if cvr != nil && avr != nil {
			x, err = fluxarray.FloatConditionalCConstAConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Float(),

				(*avr).Float(),

				mem,
			)
		} else if cvr != nil {
			x, err = fluxarray.FloatConditionalCConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Float(),

				a.Arr().(*fluxarray.Float),
				mem,
			)
		} else if avr != nil {
			x, err = fluxarray.FloatConditionalAConst(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Float),

				(*avr).Float(),

				mem,
			)
		} else {
			x, err = fluxarray.FloatConditional(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Float),
				a.Arr().(*fluxarray.Float),
				mem,
			)
		}
		if err != nil {
			return nil, err
		}
		return NewVectorValue(x, semantic.BasicFloat), nil

	case semantic.Bool:
		var (
			x   *fluxarray.Boolean
			err error
		)
		if cvr != nil && avr != nil {
			x, err = fluxarray.BooleanConditionalCConstAConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Bool(),

				(*avr).Bool(),

				mem,
			)
		} else if cvr != nil {
			x, err = fluxarray.BooleanConditionalCConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Bool(),

				a.Arr().(*fluxarray.Boolean),
				mem,
			)
		} else if avr != nil {
			x, err = fluxarray.BooleanConditionalAConst(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Boolean),

				(*avr).Bool(),

				mem,
			)
		} else {
			x, err = fluxarray.BooleanConditional(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Boolean),
				a.Arr().(*fluxarray.Boolean),
				mem,
			)
		}
		if err != nil {
			return nil, err
		}
		return NewVectorValue(x, semantic.BasicBool), nil

	case semantic.String:
		var (
			x   *fluxarray.String
			err error
		)
		if cvr != nil && avr != nil {
			x, err = fluxarray.StringConditionalCConstAConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Str(),

				(*avr).Str(),

				mem,
			)
		} else if cvr != nil {
			x, err = fluxarray.StringConditionalCConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Str(),

				a.Arr().(*fluxarray.String),
				mem,
			)
		} else if avr != nil {
			x, err = fluxarray.StringConditionalAConst(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.String),

				(*avr).Str(),

				mem,
			)
		} else {
			x, err = fluxarray.StringConditional(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.String),
				a.Arr().(*fluxarray.String),
				mem,
			)
		}
		if err != nil {
			return nil, err
		}
		return NewVectorValue(x, semantic.BasicString), nil

	case semantic.Time:
		var (
			x   *fluxarray.Int
			err error
		)
		if cvr != nil && avr != nil {
			x, err = fluxarray.IntConditionalCConstAConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Time().Time().UnixNano(),

				(*avr).Time().Time().UnixNano(),

				mem,
			)
		} else if cvr != nil {
			x, err = fluxarray.IntConditionalCConst(
				t.Arr().(*fluxarray.Boolean),

				(*cvr).Time().Time().UnixNano(),

				a.Arr().(*fluxarray.Int),
				mem,
			)
		} else if avr != nil {
			x, err = fluxarray.IntConditionalAConst(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Int),

				(*avr).Time().Time().UnixNano(),

				mem,
			)
		} else {
			x, err = fluxarray.IntConditional(
				t.Arr().(*fluxarray.Boolean),
				c.Arr().(*fluxarray.Int),
				a.Arr().(*fluxarray.Int),
				mem,
			)
		}
		if err != nil {
			return nil, err
		}
		return NewVectorValue(x, semantic.BasicTime), nil

	default:
		return nil, errors.Newf(codes.Invalid, "unsupported type for vector: %v", a.ElementType())
	}
}
