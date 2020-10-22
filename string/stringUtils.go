package stringUtils

import (
	"fmt"
	"math"
)

func HumanReadableSize(value float64) string {

	KB := float64(1024)
	MB := float64(math.Pow(1024,2))
	GB := float64(math.Pow(1024,3))
	TB := float64(math.Pow(1024,4))


	if value > TB {
		return fmt.Sprintf("%.2f TB", value/TB)
	} else if value > GB {
		return fmt.Sprintf("%.2f GB", value/GB)
	} else if value > MB {
		return fmt.Sprintf("%.2f MB", value/MB)
	} else if value > KB {
		return fmt.Sprintf("%.2f KB", value/KB)
	} else {
		return fmt.Sprintf("%.0f B", value)
	}

}

