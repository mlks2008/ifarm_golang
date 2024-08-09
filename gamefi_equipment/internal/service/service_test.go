package service

import (
	"components/common/utils/random"
	"fmt"
	"testing"
)

func Test_range(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(random.InRange(10, 12))
	}
}

func Test_Chooser(t *testing.T) {
	choiceArr := make([]random.Choice[bool, int], 0)
	choiceArr = append(choiceArr, random.NewChoice(true, 1))
	choiceArr = append(choiceArr, random.NewChoice(false, 1))

	chooser, _ := random.NewChooser(choiceArr...)
	for i := 0; i < 5; i++ {
		fmt.Println(chooser.Pick())
	}
}
