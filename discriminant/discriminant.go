package discriminant

import "math"

func Discriminant(a, b, c float64) (D float64, notice string, err error) {
	D = b*b - 4*a*c

	if D > 0 {
		notice = "Уравнение имеет два корня"
	}

	if D == 0 {
		notice = "Уравнение имеет один корень"
	}

	if D < 0 {
		notice = "Уравнение не имеет корней"
	}
	return D, notice, nil
}

func X1X2(queryA, queryB, D float64) (x1, x2 float64, err error) {
	x1 = (-queryB + math.Sqrt(D)) / (2 * queryA)
	x2 = (-queryB - math.Sqrt(D)) / (2 * queryA)
	return x1, x2, nil
}

func X(queryA, queryB float64) (x float64, err error) {
	x = -queryB / 2 * queryA

	return x, nil
}

// Если коэффициент b = 0
func SpecialCase1(queryA, queryC float64) (x1, x2 float64, err error) {
	x1 = -(-queryC / queryA)
	x2 = -queryC / queryA

	return x1, -x2, err
}

// Если коэффициент c = 0
func SpecialCase2(queryA, queryB float64) (x1, x2 float64, err error) {
	x1 = 0
	x2 = -queryB / queryA

	return x1, x2, err
}
