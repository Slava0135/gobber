package main

func compareElement(array []int, index int, value int) int {
    if index < 0 || index >= len(array) {
        return -1 // Индекс вне границ
    }
    element := array[index]
    if element > value {
        return 1 // Элемент больше
    } else if element < value {
        return -1 // Элемент меньше
    }
    return 0 // Элемент равен
}


type Person struct {
    Name string
    Age  int
}

func compareAge(people []*Person, index int, value int) int {
    if index < 0 || index >= len(people) {
        return -1 // Индекс вне границ
    }
    age := people[index].Age // Достаем возраст по индексу

    if age > value {
        return 1 // Возраст больше
    } else if age < value {
        return -1 // Возраст меньше
    }
    return 0 // Возраст равен
}