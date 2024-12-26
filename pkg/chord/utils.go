package chord

import (
	"math/big"
)

func between(start, x, end *big.Int, inclusive bool) bool {
	startInt := new(big.Int).Set(start)
	xInt := new(big.Int).Set(x)
	endInt := new(big.Int).Set(end)

	if endInt.Cmp(startInt) <= 0 {
		max := new(big.Int).Exp(big.NewInt(2), big.NewInt(m), nil)
		endInt.Add(endInt, max)

		// if x less than start, add 2^m to x as well
		if xInt.Cmp(startInt) < 0 {
			xInt.Add(xInt, max)
		}
	}

	if inclusive {
		return xInt.Cmp(startInt) > 0 && xInt.Cmp(endInt) <= 0
	}

	return xInt.Cmp(startInt) > 0 && xInt.Cmp(endInt) < 0
}