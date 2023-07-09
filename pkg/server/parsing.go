package server

import (
	"fmt"
	"strconv"
	"strings"
)

func parseValue1000(value100Raw string) (int, error) {
	splitValue := strings.Split(value100Raw, ".")

	if len(splitValue) > 2 {
		return 0, fmt.Errorf("Unexpected value100 format, please supply as a number between 0.0 and 100.0%%")
	}

	parseValue100, err := strconv.Atoi(splitValue[0])

	if err != nil {
		return 0, fmt.Errorf("Unexpected value100 format, please supply as a number between 0.0 and 100.0%%")
	}

	if parseValue100 < 0 || parseValue100 > 100 {
		return 0, fmt.Errorf("Value is out of range, please supply as a number between 0.0 and 100.0%%")
	}

	var parseValueDecimal1 int = 0

	if len(splitValue) == 2 {
		parseValueDecimal1, err = strconv.Atoi(splitValue[1][:1])

		if err != nil {
			return 0, fmt.Errorf("Unexpected value100 format, please supply as a number between 0.0 and 100.0%%")
		}
	}

	return (parseValue100 * 10) + parseValueDecimal1, nil
}
