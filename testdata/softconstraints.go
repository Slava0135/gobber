package main

func compareAndIncrement(a, b int) int {
    if a > b {
        c := a + 1

        if (c > b) {
            return 1
        } else {
            return -1
        }
    }

    return 42
}