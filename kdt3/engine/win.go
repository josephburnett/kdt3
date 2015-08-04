package engine

const (
        Decline = -1
        Neutral = 0
        Incline = 1
)

func eachDirection(K int, fn func ([]int)) {
        direction := make([]int, K)
        var recur func(fn func([]int), depth, dir int)
        recur = func(fn func([]int), depth, dir int) {
                direction[depth] = dir
                if depth == 0 {
                        fn(direction)
                } else {
                        recur(fn, depth-1, Decline)
                        recur(fn, depth-1, Neutral)
                        recur(fn, depth-1, Incline)
                }
        }
        recur(fn, K-1, Decline)
        recur(fn, K-1, Neutral)
        recur(fn, K-1, Incline)
}
