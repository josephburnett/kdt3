package engine

import (
        "testing"
)

func TestEachDirection1(t *testing.T) {
        expect := directions([][]int{
                []int{Decline},
                []int{Neutral},
                []int{Incline},
        })
        actual := collect(1)
        if !actual.equals(expect) {
                t.Errorf("Expected: %v Actual: %v", expect, actual)
        }
}

func TestEachDirection2(t *testing.T) {
        expect := directions([][]int{
                []int{Decline, Decline},
                []int{Neutral, Decline},
                []int{Incline, Decline},
                []int{Decline, Neutral},
                []int{Neutral, Neutral},
                []int{Incline, Neutral},
                []int{Decline, Incline},
                []int{Neutral, Incline},
                []int{Incline, Incline},
        })
        actual := collect(2)
        if !actual.equals(expect) {
                t.Error("Expected: %v Actual: %v", expect, actual)
        }
}

func TestEachDirectionK(t *testing.T) {
        if len(collect(3)) != 27 {
                t.Fail()
        }
        if len(collect(4)) != 81 {
                t.Fail()
        }
        if len(collect(5)) != 243 {
                t.Fail()
        }
}

type directions [][]int

func collect(K int) directions {
        actual := directions(make([][]int, 0, 3))
        eachDirection(K, func(direction []int) {
                d := make([]int, len(direction))
                copy(d, direction)
                actual = append(actual, d)
        })
        return actual
}

func (actual directions) equals(expect directions) bool {
        if len(actual) != len(expect) {
                return false
        }
        for i, _ := range expect {
                a, e := actual[i], expect[i]
                if len(a) != len(e) {
                        return false
                }
                for j, v := range a {
                        if v != e[j] {
                                return false
                        }
                }
        }
        return true
}
