package workpool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyPool(t *testing.T) {
	p := New(nil, 10)
	p.Run()

}

func TestWithWork(t *testing.T) {
	tasks := []*Task{
		NewTask(func() error {
			fmt.Println("Hey yo test 1")
			return nil
		}),
		NewTask(func() error {
			fmt.Println("Hey yo test 2")
			return nil
		}),
		NewTask(func() error {
			fmt.Println("Hey yo test 3")
			return nil
		}),
		NewTask(func() error {
			fmt.Println("Hey yo test 4")
			return nil
		}),
		NewTask(func() error {
			fmt.Println("Hey yo test 5")
			return nil
		}),
		NewTask(func() error {
			fmt.Println("Hey yo test 6")
			return nil
		}),
	}
	p := New(tasks, 10)
	p.Run()

	assert.True(t, len(p.Errors()) == 0)

}

func TestWithError(t *testing.T) {
	tasks := []*Task{
		NewTask(func() error {
			fmt.Println("Hey with error")
			return nil
		}),
		NewTask(func() error {
			fmt.Println("Hey yo with error")
			return nil
		}),
		NewTask(func() error { return fmt.Errorf("error") }),
	}

	p := New(tasks, 10)
	p.Run()

	assert.True(t, len(p.Errors()) == 1)

}
