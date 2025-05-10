package test_demo

import (
	"fmt"
	"testing"
)

type Test struct{ Val int }

func TestRef(t *testing.T) {
	list := []Test{{Val: 1}, {Val: 2}}

	// 方式1：值拷贝（失败）
	t1 := list[0]
	t1.Val = 100
	fmt.Println(list) // 输出: [{1} {2}]

	// 方式2：指针操作（成功）
	t2 := &list[1]
	t2.Val = 200
	fmt.Println(list) // 输出: [{1} {200}]
}

func TestRef2(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	x := a[0]
	x = 123
	fmt.Println(a, x)
	y := &a[1]
	*y = 123
	fmt.Println(a, *y)

}
