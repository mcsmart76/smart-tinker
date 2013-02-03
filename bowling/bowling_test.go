package bowling

import "fmt"

func ExampleGutter() {
	fmt.Println(Score([]int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
	// Output: 0
}

func ExampleOnes() {
	fmt.Println(Score([]int{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	// Output: 20
}

func ExampleSpare1() {
	fmt.Println(Score([]int{
		1, 9, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	// Output: 29
}

func ExampleSpare2() {
	fmt.Println(Score([]int{
		2, 8, 5, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
	// Output: 21
}

func ExampleStrike1() {
	fmt.Println(Score([]int{
		10, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	// Output: 30
}

func ExampleStrike2() {
	fmt.Println(Score([]int{
		10, 3, 4, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
	// Output: 24
}

func ExampleTurkey() {
	fmt.Println(Score([]int{
		10, 10, 10, 10, 10, 10,
                0, 0, 0, 0, 0, 0, 0, 0}))
	// Output: 150
}

func ExampleTenthFrame1() {
	fmt.Println(Score([]int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 3, 4}))
	// Output: 7
}

func ExampleTenthFrame2() {
	fmt.Println(Score([]int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 3, 7, 5}))
	// Output: 15
}

func ExampleTenthFrame3() {
	fmt.Println(Score([]int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 3, 7, 10}))
	// Output: 20
}

func ExampleTenthFrame4() {
	fmt.Println(Score([]int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 10, 3, 7}))
	// Output: 20
}

func ExampleTenthFrame5() {
	fmt.Println(Score([]int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 10, 10, 7}))
	// Output: 27
}

func ExamplePerfectGame() {
	fmt.Println(Score([]int{
		10, 10, 10, 10, 10, 10,
                10, 10, 10, 10, 10, 10}))
	// Output: 300
}

func ExamplePerfectGameWithExtraRoll() {
	fmt.Println(Score([]int{
		10, 10, 10, 10, 10, 10,
                10, 10, 10, 10, 10, 10, 0}))
	// Output: -1
}

func ExampleInvalidRoll1() {
	fmt.Println(Score([]int{-1}))
	// Output:
	// -1
}

func ExampleInvalidRoll2() {
	fmt.Println(Score([]int{11}))
	// Output:
	// -1
}

func ExampleInvalidFrame() {
	fmt.Println(Score([]int{9, 9}))
	// Output:
	// -1
}

func ExampleInvalidGame() {
	fmt.Println(Score([]int{0, 0, 0, 0}))
	// Output:
	// -1
}
