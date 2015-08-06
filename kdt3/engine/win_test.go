package engine

import (
        "reflect"
        "testing"

        "kdt3/model"
)

func TestWinnableBoardK1(t *testing.T) {
        board := &model.Board{
                K: 2,
                Size: 2,
                D: &model.Cell{
                        D: []*model.Cell{
                                &model.Cell{
                                        D: []*model.Cell{
                                                &model.Cell{IsClaimed: true, Player: 1},
                                                &model.Cell{IsClaimed: true, Player: 2},
                                        },
                                },
                                &model.Cell{
                                        D: []*model.Cell{
                                                &model.Cell{},
                                                &model.Cell{IsClaimed: true, Player: 1},
                                        },
                                },
                        },
                },
        }
        rules := &model.Rules{InARow: 2}
        expect := []model.Segment{
                model.Segment{model.Point{0, 0}, model.Point{1, 1},},
                model.Segment{model.Point{1, 1}, model.Point{0, 0},},
        }
        actual := WinnableBoard{board}.GetWins(1, rules)
        if !reflect.DeepEqual(expect, actual) {
                t.Errorf("Expected %v but got %v\n", expect, actual)
        }
}

func TestEachDirection1(t *testing.T) {
        expect := []model.Direction{
                model.Direction{Decline},
                model.Direction{Neutral},
                model.Direction{Incline},
        }
        actual := collect(1)
        if !reflect.DeepEqual(actual, expect) {
                t.Errorf("Expected: %v Actual: %v", expect, actual)
        }
}

func TestEachDirection2(t *testing.T) {
        expect := []model.Direction{
                model.Direction{Decline, Decline},
                model.Direction{Neutral, Decline},
                model.Direction{Incline, Decline},
                model.Direction{Decline, Neutral},
                model.Direction{Neutral, Neutral},
                model.Direction{Incline, Neutral},
                model.Direction{Decline, Incline},
                model.Direction{Neutral, Incline},
                model.Direction{Incline, Incline},
        }
        actual := collect(2)
        if !reflect.DeepEqual(actual, expect) {
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

func collect(K int) []model.Direction {
        actual := make([]model.Direction, 0, 3)
        eachDirection(K, func(direction model.Direction) {
                d := make(model.Direction, len(direction))
                copy(d, direction)
                actual = append(actual, d)
        })
        return actual
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
        if !isWinningVector(K, root, player, rules, model.Point{0}, model.Direction{Incline}) {
                t.Errorf("Expected K1 win.")
        }
        if !isWinningVector(K, root, player, rules, model.Point{1}, model.Direction{Decline}) {
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
        if isWinningVector(K, root, player, rules, model.Point{0}, model.Direction{Incline}) {
                t.Errorf("Expected no K1 win.")
        }
        if isWinningVector(K, root, player, rules, model.Point{1}, model.Direction{Decline}) {
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
        if !isWinningVector(K, root, player, rules, model.Point{0,0}, model.Direction{Incline, Incline}) {
                t.Errorf("Expected K2 diagonal win")
        }
        if !isWinningVector(K, root, player, rules, model.Point{1,1}, model.Direction{Decline, Decline}) {
                t.Errorf("Expected K2 diagonal win")
        }
        if isWinningVector(K, root, player, rules, model.Point{1,0}, model.Direction{Decline, Incline}) {
                t.Errorf("Unexpected K2 diagonal win")
        }
        if isWinningVector(K, root, player, rules, model.Point{0,1}, model.Direction{Incline, Decline}) {
                t.Errorf("Unexpected K2 diagonal win")
        }
        if isWinningVector(K, root, player, rules, model.Point{0,0}, model.Direction{Neutral, Decline}) {
                t.Errorf("Unexpected K2 out-of-bounds win")
        }
        if isWinningVector(K, root, player, rules, model.Point{1,1}, model.Direction{Incline, Neutral}) {
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
        if isWinningVector(K, root, player, rules, model.Point{0,0}, model.Direction{Neutral, Neutral}) {
                t.Errorf("Expected no K2 neutral win")
        }
}

func TestEachPoint(t *testing.T) {
        K := 2
        root := &model.Cell{
                D: []*model.Cell{
                        &model.Cell{
                                D: []*model.Cell{
                                        &model.Cell{},
                                        &model.Cell{},
                                },
                        },
                        &model.Cell{
                                D: []*model.Cell{
                                        &model.Cell{},
                                        &model.Cell{},
                                },
                        },
                },
        }
        expect := []model.Point{
                model.Point{0, 0},
                model.Point{1, 0},
                model.Point{0, 1},
                model.Point{1, 1},
        }
        actual := make([]model.Point, 0)
        eachPoint(K, root, func(p model.Point) {
                point := make(model.Point, len(p))
                copy(point, p)
                actual = append(actual, point)
        })
        if !reflect.DeepEqual(expect, actual) {
                t.Errorf("Expected %v instead of %v\n", expect, actual)
        }
}
