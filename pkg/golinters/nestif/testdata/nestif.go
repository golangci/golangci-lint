//golangcitest:args -Enestif
//golangcitest:config_path testdata/nestif.yml
package testdata

func _() {
	var b1, b2, b3, b4 bool

	if b1 { // want "`if b1` has complex nested blocks \\(complexity: 1\\)"
		if b2 { // +1
		}
	}

	if b1 { // want "`if b1` has complex nested blocks \\(complexity: 3\\)"
		if b2 { // +1
			if b3 { // +2
			}
		}
	}

	if b1 { // want "`if b1` has complex nested blocks \\(complexity: 5\\)"
		if b2 { // +1
		} else if b3 { // +1
			if b4 { // +2
			}
		} else { // +1
		}
	}

	if b1 { // want "`if b1` has complex nested blocks \\(complexity: 9\\)"
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

	if b1 == b2 == b3 { // want "`if b1 == b2 == b3` has complex nested blocks \\(complexity: 1\\)"
		if b4 { // +1
		}
	}
}
