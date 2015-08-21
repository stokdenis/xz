// Copyright 2015 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package newlzma

import "errors"

// Maximum and minimum values for individual parameters.
const (
	MinLC       = 0
	MaxLC       = 8
	MinLP       = 0
	MaxLP       = 4
	MinPB       = 0
	MaxPB       = 4
	MinDictSize = 1 << 12
	MaxDictSize = 1<<32 - 1
)

// MaxProperties is the maximum value for a properties value.
const MaxProperties = (MaxPB+1)*(MaxLP+1)*(MaxLC+1) - 1

// Properties contains the parameters LC, LP and PB. The parameter LC
// defines the number of literal context bits; parameter LP the number
// of literal position bits and PB the number of position bits.
type Properties byte

// NewProperties returns a new properties value. It verifies the validity of
// the arguments.
func NewProperties(lc, lp, pb int) (p Properties, err error) {
	if err = verifyProperties(lc, lp, pb); err != nil {
		return
	}
	return Properties((pb*5+lp)*9 + lc), nil
}

// LC returns the number of literal context bits.
func (p Properties) LC() int {
	return int(p) % 9
}

// LP returns the number of literal position bits.
func (p Properties) LP() int {
	return (int(p) / 9) % 5
}

// PB returns the number of position bits.
func (p Properties) PB() int {
	return (int(p) / 45) % 5
}

// verifyProperties checks the argument for any errors.
func verifyProperties(lc, lp, pb int) error {
	if !(MinLC <= lc && lc <= MaxLC) {
		return errors.New("lc out of range")
	}
	if !(MinLP <= lp && lp <= MaxLP) {
		return errors.New("lp out of range")
	}
	if !(MinPB <= pb && pb <= MaxPB) {
		return errors.New("pb out of range")
	}
	return nil
}
