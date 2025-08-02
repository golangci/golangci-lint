//golangcitest:args -Egocheckerrbeforeuse
//golangcitest:config_path testdata/custom.yml
package testdata

func returns2Values() (int, error) {
	return 0, nil
}

func PositiveWithCustomDistance() {
	i, err := returns2Values()

	print(i)

	if err != nil {
		return
	}
}
