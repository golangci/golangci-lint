//args: -Enestif
//config: linters-settings.nestif.min-complexity=1
package testdata

func _() {
	var b1, b2, b3, b4 bool

	if b1 { // ERROR "`if b1` is deeply nested \\(complexity: 1\\)"
		if b2 { // +1
		}
	}

	if b1 { // ERROR "`if b1` is deeply nested \\(complexity: 3\\)"
		if b2 { // +1
			if b3 { // +2
			}
		}
	}

	if b1 { // ERROR "`if b1` is deeply nested \\(complexity: 5\\)"
		if b2 { // +1
		} else if b3 { // +1
			if b4 { // +2
			}
		} else { // +1
		}
	}

	if b1 { // ERROR "`if b1` is deeply nested \\(complexity: 9\\)"
		if b2 { // +1
			if b3 { // +2
			}
		}

		if b2 { // +1
			if b3 { // +2
				if b4 { // +3
				}
			}
		}
	}

	if b1 == b2 == b3 { // ERROR "`if b1 == b2 == b3` is deeply nested \\(complexity: 1\\)"
		if b4 { // +1
		}
	}
}
