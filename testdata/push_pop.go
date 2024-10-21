package main

func pushPopIncrementality(j int) int {
    result := j

    for i := 1; i <= 10; i++ {
        result += i
    }

    if result%2 == 0 {
        result++
    }

    return result
}