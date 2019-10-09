package main

var operatorPriority = map[rune]int{
	'*': 4,
	'?': 4,
	'+': 4,
	'.': 3,
	'|': 2,
	'(': 1,
}

func re2post(re string) string {
	opts := make([]rune, 0)
	output := make([]rune, 0)

	pushOperator := func(op rune) {
		priority := operatorPriority[op]
		for len(opts) > 0 {
			top := opts[len(opts)-1]

			if operatorPriority[top] >= priority {
				output = append(output, top)
				opts = opts[:len(opts)-1]
			} else {
				break
			}
		}
		opts = append(opts, op)
	}

	lastIsChar := false
	for _, char := range re {
		switch char {
		case '*', '+', '?':
			pushOperator(char)
			lastIsChar = true
		case '|':
			pushOperator(char)
			lastIsChar = false
		case '(':
			if lastIsChar {
				pushOperator('.')
			}
			opts = append(opts, '(')
			lastIsChar = false
		case ')':
			for len(opts) > 0 {
				op := opts[len(opts)-1]
				opts = opts[:len(opts)-1]
				if op == '(' {
					break
				}
				output = append(output, op)
			}
			lastIsChar = true
		default:
			if lastIsChar {
				pushOperator('.')
			}
			output = append(output, char)
			lastIsChar = true
		}
	}

	for len(opts) > 0 {
		output = append(output, opts[len(opts)-1])
		opts = opts[:len(opts)-1]
	}
	return string(output)
}
