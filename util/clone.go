package util

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func CloneNegativeNumeric(numeric pgtype.Numeric) pgtype.Numeric {
	clone := pgtype.Numeric{
		Int:              new(big.Int).Set(numeric.Int),
		Exp:              numeric.Exp,
		NaN:              numeric.NaN,
		Valid:            numeric.Valid,
		InfinityModifier: numeric.InfinityModifier,
	}
	clone.Int = clone.Int.Neg(clone.Int)
	return clone
}
