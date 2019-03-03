package main

func MapTo(numbers []int, f func(element int, index int) string) []string {
	var result []string
	
	for i, number := range numbers {
		result = append(result, f(number, i))
	}
	return result;
}

func Convert(arr []int) []string {
	values := [10]string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	var result []string; 
	
	for _, element := range arr{
		switch {
			case (element >= 0) && (element < len(values)) :
				result = append(result, values[element])
			default:
				result = append(result, "unknown")
		}
	}
	return result;
}

func main() {
}
