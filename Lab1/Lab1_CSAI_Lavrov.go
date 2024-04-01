package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// AStarAlgorithm представляє алгоритм A* для вирішення головоломки.
type AStarAlgorithm struct {
	initArr   [][]int // Початковий стан головоломки
	targetArr [][]int // Цільовий стан головоломки
}

// State представляє стан головоломки для алгоритму A*.
type State struct {
	priorityNum int       // Пріоритет (використовується для пріоритетної черги)
	arr         [][]int   // Поточний стан головоломки
	path        [][][]int // Шлях, який призвів до цього стану
}

// NewAStarAlgorithm ініціалізує новий екземпляр AStarAlgorithm.
func NewAStarAlgorithm(initArr [][]int) (*AStarAlgorithm, error) {
	// Перевірка на валідність початкової головоломки
	if len(initArr) != 3 || len(initArr[0]) != 3 {
		return nil, errors.New("Початковий масив повинен бути матрицею розміром 3x3.")
	}

	// Цільовий стан головоломки
	targetArr := [][]int{
		{1, 2, 3},
		{8, 0, 4},
		{7, 6, 5},
	}

	return &AStarAlgorithm{
		initArr:   initArr,
		targetArr: targetArr,
	}, nil
}

// getZeroPos знаходить позицію нульового елемента в головоломці.
func (a *AStarAlgorithm) getZeroPos(currentArr [][]int) []int {
	for i, row := range currentArr {
		for j, val := range row {
			if val == 0 {
				return []int{i, j}
			}
		}
	}
	panic("Неправильна головоломка: порожній простір (нуль) не знайдено.")
}

// move міняє місцями елементи в головоломці.
func (a *AStarAlgorithm) move(currentArr [][]int, zeroPos, newZeroPos []int) [][]int {
	updatedArr := a.copyArr(currentArr)
	updatedArr[zeroPos[0]][zeroPos[1]] = currentArr[newZeroPos[0]][newZeroPos[1]]
	updatedArr[newZeroPos[0]][newZeroPos[1]] = 0
	return updatedArr
}

// isPosValid перевіряє, чи нова позиція знаходиться в межах матриці.
func (a *AStarAlgorithm) isPosValid(position []int) bool {
	return position[0] >= 0 && position[0] < 3 && position[1] >= 0 && position[1] < 3
}

// getPossibleMoves повертає список можливих наступних матриць.
func (a *AStarAlgorithm) getPossibleMoves(currentArr [][]int) [][][]int {
	zeroPos := a.getZeroPos(currentArr)
	offsets := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // Всі можливі зсуви в горизонтальному та вертикальному напрямках
	var possibleMoves [][][]int

	for _, offset := range offsets {
		newZeroPos := []int{zeroPos[0] + offset[0], zeroPos[1] + offset[1]}
		if a.isPosValid(newZeroPos) {
			possibleMoves = append(possibleMoves, a.move(currentArr, zeroPos, newZeroPos))
		}
	}

	return possibleMoves
}

// calcManhattanDistance обчислює Манхеттенську відстань до цільової матриці.
func (a *AStarAlgorithm) calcManhattanDistance(currentArr [][]int) int {
	distance := 0
	for i, row := range currentArr {
		for j, val := range row {
			if val != 0 {
				distance += abs(i-(val-1)/3) + abs(j-(val-1)%3)
			}
		}
	}
	return distance
}

// algorithm знаходить шлях вирішення або генерує помилку.
func (a *AStarAlgorithm) algorithm() ([][][]int, error) {
	visited := make(map[string]bool)                                                            // відстеження відвіданих станів
	priorityQueue := make(PriorityQueue, 0)                                                     // Пріоритетна черга для вибору наступного кроку
	heap.Init(&priorityQueue)                                                                   // Ініціалізація пріоритетної черги
	heap.Push(&priorityQueue, &State{priorityNum: 0, arr: a.initArr, path: make([][][]int, 0)}) // Додавання початкового стану в чергу

	startTime := time.Now()
	for priorityQueue.Len() > 0 && time.Since(startTime).Milliseconds() < 10000 { // Виконуємо цикл, доки черга не пуста або не вичерпано час
		current := heap.Pop(&priorityQueue).(*State) // Беремо поточний стан з пріоритетної черги

		// Якщо поточний масив дорівнює цільовому, повертаємо шлях вирішення
		if a.deepEquals(current.arr, a.targetArr) {
			res := append(current.path, current.arr)
			return res, nil
		}

		stateHash := a.getStateHash(current) // Отримуємо хеш-рядок поточного стану
		if !visited[stateHash] {             // Перевіряємо, чи ми ще не відвідали цей стан
			visited[stateHash] = true                        // Позначаємо стан як відвіданий
			possibleMoves := a.getPossibleMoves(current.arr) // Отримуємо можливі наступні стани
			for _, arr := range possibleMoves {
				heuristic := a.calcManhattanDistance(arr)
				currentPath := append(current.path, current.arr)                                       // Додаємо поточний стан до шляху
				heap.Push(&priorityQueue, &State{priorityNum: heuristic, arr: arr, path: currentPath}) // Додаємо наступний стан в чергу з його пріоритетом
			}
		}
	}
	return nil, errors.New("Час виконання алгоритму перевищено!")
}

// Run виконує алгоритм A* та виводить результат.
func (a *AStarAlgorithm) Run() {
	start := time.Now()
	arr, err := a.algorithm()

	// Вимір часу роботи
	fmt.Printf("Час вирішення: %d мс.\nШлях:\n", time.Since(start).Milliseconds())

	// Виведення шляху вирішення
	if err != nil {
		fmt.Println("Щось пішло не так!")
		fmt.Println(err.Error())
		return
	}

	numOfStep := 0
	for _, subArr := range arr {
		fmt.Printf("Крок #%d\n", numOfStep)
		for _, i := range subArr {
			fmt.Println(i)
		}
		fmt.Println()
		numOfStep++
	}
}

// copyArr повертає копію заданого масиву.
func (a *AStarAlgorithm) copyArr(arr [][]int) [][]int {
	copyArr := make([][]int, len(arr))
	for i, row := range arr {
		copyArr[i] = make([]int, len(row))
		copy(copyArr[i], row)
	}
	return copyArr
}

// getStateHash повертає хеш-рядок для State.
func (a *AStarAlgorithm) getStateHash(s *State) string {
	return fmt.Sprintf("%v", s.arr)
}

// deepEquals перевіряє, чи дві матриці глибоко рівні.
func (a *AStarAlgorithm) deepEquals(arr1, arr2 [][]int) bool {
	return fmt.Sprintf("%v", arr1) == fmt.Sprintf("%v", arr2)
}

// abs повертає абсолютне значення цілого числа.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// main - головна функція для запуску алгоритму A*.
func main() {
	fmt.Println("Введіть початкову головоломку розміром 3x3 (цифри від 0 до 8, розділені пробілом, рядок за рядком):")
	initArr := make([][]int, 3)
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 3; i++ {
		fmt.Printf("Рядок #%d: ", i+1)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		nums := strings.Split(input, " ")
		for _, numStr := range nums {
			num, _ := strconv.Atoi(numStr)
			initArr[i] = append(initArr[i], num)
		}
	}

	astar, err := NewAStarAlgorithm(initArr)
	if err != nil {
		fmt.Println("Щось пішло не так!")
		fmt.Println(err.Error())
		return
	}

	astar.Run()
}

// PriorityQueue реалізує пріоритетну чергу для структури State.
type PriorityQueue []*State

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priorityNum < pq[j].priorityNum }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

// Push та Pop - функції для heap.Interface
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*State)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
