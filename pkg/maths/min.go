package maths
func MinMax(nums ...int)  (int, int){
	min := nums[0]
	max := nums[0]
	for _, value := range nums {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}