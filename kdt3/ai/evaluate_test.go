package ai

import (
        "math"
        "testing"

        m "kdt3/model"
)

var game1 *m.Game = &m.Game{
        //  0 |   |
        // ---+---+---
        //    | 1 |
        // ---+---+---
        //    |   |
        Board: &m.Board{
                D: &m.Cell{
                        D: []*m.Cell{
                                &m.Cell{
                                        D: []*m.Cell{
                                                &m.Cell{ Player: 0, IsClaimed: true },
                                                &m.Cell{},
                                                &m.Cell{},
                                        },
                                },
                                &m.Cell{
                                        D: []*m.Cell{
                                                &m.Cell{},
                                                &m.Cell{ Player: 1, IsClaimed: true },
                                                &m.Cell{},
                                        },
                                },
                                &m.Cell{
                                        D: []*m.Cell{
                                                &m.Cell{},
                                                &m.Cell{},
                                                &m.Cell{},
                                        },
                                },
                        },
                },
                K: 2,
                Size: 3,
        },
        Rules: &m.Rules{
                InARow: 3,
        },
        TurnOrder: 1,
}

var game2 *m.Game = &m.Game{
        //  0 | 0 | 1
        // ---+---+---
        //    | 1 |
        // ---+---+---
        //  1 |   |
        Board: &m.Board{
                D: &m.Cell{
                        D: []*m.Cell{
                                &m.Cell{
                                        D: []*m.Cell{
                                                &m.Cell{ Player: 0, IsClaimed: true },
                                                &m.Cell{ Player: 0, IsClaimed: true },
                                                &m.Cell{ Player: 1, IsClaimed: true },
                                        },
                                },
                                &m.Cell{
                                        D: []*m.Cell{
                                                &m.Cell{},
                                                &m.Cell{ Player: 1, IsClaimed: true },
                                                &m.Cell{},
                                        },
                                },
                                &m.Cell{
                                        D: []*m.Cell{
                                                &m.Cell{ Player: 1, IsClaimed: true },
                                                &m.Cell{},
                                                &m.Cell{},
                                        },
                                },
                        },
                },
                K: 2,
                Size: 3,
        },
        Rules: &m.Rules{
                InARow: 3,
        },
        TurnOrder: 1,
}

func TestIsOpenVector2D(t *testing.T) {
        K := game1.Board.K
        D := game1.Board.D
        rules := game1.Rules
        test := func(expect bool, player int, point, direction []int) {
                open := isOpenVector(K, D, player, rules, point, direction)
                if expect && !open {
                        t.Errorf("Expected open vector: %v %v %v", player, point, direction)
                } else if !expect && open {
                        t.Errorf("Expected no open vector: %v %v %v", player, point, direction)
                }
        }
        // Test open vector for 1 from claimed corner
        test(true, 0, []int{0,0}, []int{1,0})
        test(true, 0, []int{0,0}, []int{0,1})
        // Test open vector for 1 from unclaimed corner
        test(true, 0, []int{2,2}, []int{-1,0})
        test(true, 0, []int{2,2}, []int{0,-1})
        // Test open vector for 1 from another unclaimed corner
        test(true, 0, []int{0,2}, []int{1,0})
        test(true, 0, []int{0,2}, []int{0,-1})
        // Test no open vector for 1 through 2's claimed square
        test(false, 0, []int{0,0}, []int{1,1})
        test(false, 0, []int{0,1}, []int{1,0})
        test(false, 0, []int{0,1}, []int{0,1})
        // Test open vector for 2 from unclaimed side
        test(true, 1, []int{0,1}, []int{1,0})
        test(true, 1, []int{1,0}, []int{0,1})
        // Test open vector for 2 from unclaimed corner
        test(true, 1, []int{0,2}, []int{1,-1})
        test(true, 1, []int{2,0}, []int{0,1})
        test(true, 1, []int{2,0}, []int{-1,1})
        // Test no open vector for 2 from unclaimed corners
        test(false, 1, []int{2,2}, []int{-1,-1})
        test(false, 1, []int{0,2}, []int{0,-1})
}

func TestCountOpenVectorsGame1(t *testing.T) {
        board := game1.Board
        rules := game1.Rules
        test := func(expect, player int) {
                open, wins := countOpenVectors(board, player, rules)
                if open != expect {
                        t.Errorf("Unexpected count %v (expected %v)", open, expect)
                }
                if wins != 0 {
                        t.Errorf("Unexpected wins %v (expected 0)", wins)
                }
        }
        test(8, 0)
        test(10, 1)
        test(4, 2)
}

func TestCountOpenVectorsGame2(t *testing.T) {
        board := game2.Board
        rules := game2.Rules
        test := func(expectOpen, expectWins, player int) {
                open, wins := countOpenVectors(board, player, rules)
                if open != expectOpen {
                        t.Errorf("Unexpected count %v (expected %v)", open, expectOpen)
                }
                if wins != expectWins {
                        t.Errorf("Unexpected wins %v (expected %v)", wins, expectWins)
                }
        }
        test(0, 0, 0)
        test(8, 2, 1)
        test(0, 0, 2)
}

func TestEvaluate(t *testing.T) {
        game1Score0 := Evaluate(game1.Board, game1.Rules, 2, 0)
        if game1Score0 != -2 {
                t.Errorf("Unexpected evaluation for game 1 player 0: %v (expected %v)", game1Score0, -2)
        }
        game1Score1 := Evaluate(game1.Board, game1.Rules, 2, 1)
        if game1Score1 != 2 {
                t.Errorf("Unexpected evaluation for game 1 player 1: %v (expected %v)", game1Score1, 2)
        }
        game2Score0 := Evaluate(game2.Board, game2.Rules, 2, 0)
        if game2Score0 != math.MinInt64 {
                t.Errorf("Unexpected evaluation for game 2 player 0: %v (expected %v)", game2Score0, math.MinInt64)
        }
        game2Score1 := Evaluate(game2.Board, game2.Rules, 2, 1)
        if game2Score1 != math.MaxInt64 {
                t.Errorf("Unexpected evaluation for game 2 player 1: %v (expected %v)", game2Score1, math.MaxInt64)
        }
}

func setupBenchmark(K, D int) (board *m.Board, rules *m.Rules) {
        board = m.NewBoard(K, D)
        rules = &m.Rules{ InARow: D }
        return
}
// BenchmarkEvaluateK2D3	   30000	     38052 ns/op
func BenchmarkEvaluateK2D3(b *testing.B) {
        board, rules := setupBenchmark(2,3)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
                Evaluate(board, rules, 2, 0)
        }
}

// BenchmarkEvaluateK3D4	    1000	   1158470 ns/op
func BenchmarkEvaluateK3D4(b *testing.B) {
        board, rules := setupBenchmark(3,4)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
                Evaluate(board, rules, 2, 0)
        }
}

// BenchmarkEvaluateK4D5	      30	  36219983 ns/op
func BenchmarkEvaluateK4D5(b *testing.B) {
        board, rules := setupBenchmark(4,5)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
                Evaluate(board, rules, 2, 0)
        }
}

// BenchmarkEvaluateK5D6	       1	1497060537 ns/op
func BenchmarkEvaluateK5D6(b *testing.B) {
        board, rules := setupBenchmark(5,6)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
                Evaluate(board, rules, 2, 0)
        }
}

// BenchmarkEvaluateK6D7	       1	79817763800 ns/op
func BenchmarkEvaluateK6D7(b *testing.B) {
        board, rules := setupBenchmark(6,7)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
                Evaluate(board, rules, 2, 0)
        }
}
