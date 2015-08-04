package engine

import (
        "testing"

        "kdt3/model"
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

func TestWinningK1(t *testing.T) {
        K := 1
        player := 1
        root := &model.Cell{
                D: []*model.Cell{
                        &model.Cell{IsClaimed: true, Player: player},
                        &model.Cell{IsClaimed: true, Player: player},
                },
        }
        rules := &model.Rules{InARow: 2}
        if !isWinningVector(K, root, player, rules, []int{0}, []int{Incline}) {
                t.Errorf("Expected K1 win.")
        }
        if !isWinningVector(K, root, player, rules, []int{1}, []int{Decline}) {
                t.Errorf("Expected K1 win.")
        }
}

func TestNotWinningK1(t *testing.T) {
        K := 1
        player := 1
        root := &model.Cell{
                D: []*model.Cell{
                        &model.Cell{IsClaimed: true, Player: player},
                        &model.Cell{IsClaimed: false},
                },
        }
        rules := &model.Rules{InARow: 2}
        if isWinningVector(K, root, player, rules, []int{0}, []int{Incline}) {
                t.Errorf("Expected no K1 win.")
        }
        if isWinningVector(K, root, player, rules, []int{1}, []int{Decline}) {
                t.Errorf("Expected no K1 win.")
        }
}

func TestWinningK2Diagonal(t *testing.T) {
        K := 2
        player := 1
        root := &model.Cell{
                D: []*model.Cell{
                        &model.Cell{
                                D: []*model.Cell{
                                        &model.Cell{IsClaimed: true, Player: player},
                                        &model.Cell{IsClaimed: false},
                                },
                        },
                        &model.Cell{
                                D: []*model.Cell{
                                        &model.Cell{IsClaimed: false},
                                        &model.Cell{IsClaimed: true, Player: player},
                                },
                        },
                },
        }
        rules := &model.Rules{InARow: 2}
        if !isWinningVector(K, root, player, rules, []int{0,0}, []int{Incline, Incline}) {
                t.Errorf("Expected K2 diagonal win")
        }
        if !isWinningVector(K, root, player, rules, []int{1,1}, []int{Decline, Decline}) {
                t.Errorf("Expected K2 diagonal win")
        }
        if isWinningVector(K, root, player, rules, []int{1,0}, []int{Decline, Incline}) {
                t.Errorf("Unexpected K2 diagonal win")
        }
        if isWinningVector(K, root, player, rules, []int{0,1}, []int{Incline, Decline}) {
                t.Errorf("Unexpected K2 diagonal win")
        }
        if isWinningVector(K, root, player, rules, []int{0,0}, []int{Neutral, Decline}) {
                t.Errorf("Unexpected K2 out-of-bounds win")
        }
        if isWinningVector(K, root, player, rules, []int{1,1}, []int{Incline, Neutral}) {
                t.Errorf("Unexpected K2 out-of-bounds win")
        }
}

func TestNotWinningK2(t *testing.T) {
        K := 2
        player := 1
        root := &model.Cell{
                D: []*model.Cell{
                        &model.Cell{
                                D: []*model.Cell{
                                        &model.Cell{IsClaimed: true, Player: player},
                                        &model.Cell{IsClaimed: true, Player: player},
                                },
                        },
                        &model.Cell{
                                D: []*model.Cell{
                                        &model.Cell{IsClaimed: true, Player: player},
                                        &model.Cell{IsClaimed: true, Player: player},
                                },
                        },
                },
        }
        rules := &model.Rules{InARow: 2}
        if isWinningVector(K, root, player, rules, []int{0,0}, []int{Neutral, Neutral}) {
                t.Errorf("Expected no K2 neutral win")
        }
}
