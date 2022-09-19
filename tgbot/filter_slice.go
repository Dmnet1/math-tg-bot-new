package tgbot

func Filter(equations []Equation, fn func(equation Equation) bool) []Equation {
	var filtered []Equation
	for _, equation := range equations {
		if fn(equation) {
			filtered = append(filtered, equation)
		}
	}
	return filtered
}
