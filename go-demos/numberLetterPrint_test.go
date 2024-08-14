package go_demos

import (
	"fmt"
	"sync"
	"testing"
)

func Test_numberLetterPrint1(t *testing.T) {
	var number = make(chan struct{})
	var letter = make(chan struct{})
	var wait sync.WaitGroup

	wait.Add(1)
	go func() {
		defer wait.Done()
		var i = 1
		for {
			select {
			case _, o := <-number:
				if !o {
					return
				}
				fmt.Print(i)
				i++
				fmt.Print(i)
				i++
				letter <- struct{}{}
			}
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		var i = 'A'
		for {
			select {
			case _, o := <-letter:
				if !o {
					return
				}
				if i >= 'Z' {
					close(number)
					close(letter)
					continue
				}
				fmt.Print(string(i))
				i++
				fmt.Print(string(i))
				i++
				number <- struct{}{}
			}
		}
	}()

	number <- struct{}{}
	wait.Wait()
}

func Test_numberLetterPrint2(t *testing.T) {
	var number = make(chan struct{})
	var letter = make(chan struct{})
	var done = make(chan struct{})

	go func() {
		var i = 1
		for {
			select {
			case <-number:
				fmt.Print(i)
				i++
				fmt.Print(i)
				i++
				letter <- struct{}{}
			}
		}
	}()

	go func() {
		var i = 'A'
		for {
			select {
			case <-letter:
				if i >= 'Z' {
					done <- struct{}{}
					continue
				}
				fmt.Print(string(i))
				i++
				fmt.Print(string(i))
				i++
				number <- struct{}{}
			}
		}
	}()

	number <- struct{}{}
	select {
	case <-done:
	}
}
