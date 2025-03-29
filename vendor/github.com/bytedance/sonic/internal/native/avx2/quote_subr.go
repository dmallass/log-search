// +build !noasm !appengine
// Code generated by asm2asm, DO NOT EDIT.

package avx2

import (
	`github.com/bytedance/sonic/loader`
)

const (
    _entry__quote = 144
)

const (
    _stack__quote = 72
)

const (
    _size__quote = 2880
)

var (
    _pcsp__quote = [][2]uint32{
        {0x1, 0},
        {0x6, 8},
        {0x8, 16},
        {0xa, 24},
        {0xc, 32},
        {0xd, 40},
        {0x11, 48},
        {0xb10, 72},
        {0xb11, 48},
        {0xb13, 40},
        {0xb15, 32},
        {0xb17, 24},
        {0xb19, 16},
        {0xb1a, 8},
        {0xb1e, 0},
        {0xb40, 72},
    }
)

var _cfunc_quote = []loader.CFunc{
    {"_quote_entry", 0,  _entry__quote, 0, nil},
    {"_quote", _entry__quote, _size__quote, _stack__quote, _pcsp__quote},
}
