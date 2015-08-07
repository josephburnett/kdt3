package model

import (
        "reflect"
        "testing"
)

func TestNewBoardK1D3(t *testing.T) {
        expect := &Board{
                K: 1,
                Size: 3,
                D: &Cell{
                        D: []*Cell{
                                &Cell{},
                                &Cell{},
                                &Cell{},
                        },
                },
        }
        actual := NewBoard(1, 3)
        if !reflect.DeepEqual(expect, actual) {
                t.Errorf("Unexpected board: %v wanted: %v\n", actual, expect)
        }
}

func TestNewBoardK2D2(t *testing.T) {
        expect := &Board{
                K: 2,
                Size: 2,
                D: &Cell{
                        D: []*Cell{
                                &Cell{
                                        D: []*Cell{
                                                &Cell{},
                                                &Cell{},
                                        },
                                },
                                &Cell{
                                        D: []*Cell{
                                                &Cell{},
                                                &Cell{},
                                        },
                                },
                        },
                },
        }
        actual := NewBoard(2, 2)
        if !reflect.DeepEqual(expect, actual) {
                t.Errorf("Unexpected board: %v wanted: %v\n", actual, expect)
        }
}
