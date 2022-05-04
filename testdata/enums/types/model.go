package types

import (
	"github.com/swaggo/swag/testdata/enums/consts"
)

type Class int

const (
	None Class = -1
	A    Class = consts.Base + (iota+1-1)*2/2%100 - (1&1 | 1) + (2 ^ 2) // AAA
	B                                                                   /* BBB */
	C
	D = C + 1
	F = Class(5)
	//G is not enum
	G = H + 10
	//H is not enum
	H = 10
	//I is not enum
	I = int(F + 2)
)

const J = 1 << uint16(I)

type Mask int

const (
	Mask1 Mask = 0x02 << iota >> 1 // Mask1
	Mask2                          /* Mask2 */
	Mask3                          // Mask3
	Mask4           