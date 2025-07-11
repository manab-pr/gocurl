package helper

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PrintColoredJSON(v interface{}, indent int) {
	ind := strings.Repeat("  ", indent)

	switch val := v.(type) {
	case map[string]interface{}:
		fmt.Println("{")
		length := len(val)
		i := 0
		for k, v := range val {
			fmt.Printf("%s  %s: ", ind, color.YellowString(`"%s"`, k))
			PrintColoredJSON(v, indent+1)
			i++
			if i < length {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Printf("%s}", ind)

	case []interface{}:
		fmt.Println("[")
		for i, v := range val {
			fmt.Printf("%s  ", ind)
			PrintColoredJSON(v, indent+1)
			if i < len(val)-1 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Printf("%s]", ind)

	case string:
		fmt.Print(color.GreenString(`"%s"`, val))
	case float64:
		fmt.Print(color.BlueString("%v", val))
	case bool:
		fmt.Print(color.RedString("%v", val))
	case nil:
		fmt.Print(color.RedString("null"))
	default:
		fmt.Print(val)
	}
}
