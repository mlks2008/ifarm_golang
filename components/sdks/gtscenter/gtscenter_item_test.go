package gtscenter

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_AddItem(t *testing.T) {
	gtscenter := getGtsCenter()

	var id = time.Now().Unix()
	err := gtscenter.AddItem(context.Background(), "1481130751615045", "hero", []string{"10133", "10127"}, []int64{50, 50}, id)

	fmt.Println("Res:", err)
	fmt.Println("")
}

func Test_SubItem(t *testing.T) {
	gtscenter := getGtsCenter()

	var id = time.Now().Unix()
	err := gtscenter.SubItem(context.Background(), "1481130751615045", "hero", []string{"10002"}, []int64{1}, id)

	fmt.Println("Res:", err)
	fmt.Println("")
}

func Test_ReturnItem(t *testing.T) {
	gtscenter := getGtsCenter()

	var id = time.Now().Unix()
	err := gtscenter.FreezeItem(context.Background(), "1481130751615045", "hero", []string{"36341"}, id)
	fmt.Println("Res:", err)
	fmt.Println("")

	err = gtscenter.ReturnItem(context.Background(), "1481130751615045", "hero", id)
	fmt.Println("Res:", err)
	fmt.Println("")
}
