package main

import (
	"fmt"
)

// Структура, що представляє бінарний перцептрон
type BinaryPerceptron struct {
	inputSize int        // Розмір вхідних даних
	weights   [3]float64 // Ваги
}

// Конструктор для створення нового бінарного перцептрона
func NewBinaryPerceptron(inputSize int) *BinaryPerceptron {
	return &BinaryPerceptron{
		inputSize: inputSize,
		weights:   [3]float64{0.2, 0.02, 0.002}, // Початкові ваги
	}
}

// Метод для передбачення вихідного значення
func (bp *BinaryPerceptron) predict(inputs [3]float64) (float64, int) {
	// Обчислюємо зважену суму вхідних даних
	weightedSum := bp.weights[0]*inputs[0] + bp.weights[1]*inputs[1] + bp.weights[2]*inputs[2]
	// Передбачене значення (0 або 1)
	prediction := 0
	if weightedSum > 0 {
		prediction = 1
	}
	return weightedSum, prediction
}

// Метод для навчання перцептрона
func (bp *BinaryPerceptron) train(inputs [][3]float64, targets []float64, learningRate float64, epochs int) {
	for epoch := 0; epoch < epochs; epoch++ {
		fmt.Println("Iteration #", epoch+1)
		isError := false
		// Проходимося по усіх вхідних даних
		for i, input := range inputs {
			// Передбачене значення та зважена сума для поточних вхідних даних
			sum, prediction := bp.predict(input)
			// Розрахунок помилки
			error := targets[i] - float64(prediction)
			if error != 0 {
				isError = true
			}
			// Зберігаємо попередні ваги
			prevWeights := bp.weights
			// Оновлюємо ваги згідно з правилом навчання перцептрона
			for j := range prevWeights {
				bp.weights[j] += learningRate * error * input[j]
			}
			// Виведення інформації про крок навчання
			fmt.Printf("#%d w_1=%.2f w_2=%.2f w_3=%.2f x_1=%.2f x_2=%.2f x_3=%.2f a=%.2f Y=%d T=%.2f error=%.2f delta_w_1=%.2f delta_w_2=%.2f delta_w_3=%.2f\n",
				i+1, prevWeights[0], prevWeights[1], prevWeights[2], input[0], input[1], input[2], sum, prediction, targets[i], error,
				learningRate*error*input[0], learningRate*error*input[1], learningRate*error*input[2])
		}
		fmt.Println()
		// Перевірка на закінчення навчання
		if !isError {
			fmt.Println("Training is over")
			break
		}
	}
}

// Функція для обчислення вихідних значень функції
func valuesOfFunc(X [][3]float64) [][]float64 {
	values := make([][]float64, len(X))
	for i, vInput := range X {
		result := 0.0
		if vInput[0] != 0 && vInput[1] == 0 || vInput[2] != 0 {
			result = 1.0
		}
		values[i] = []float64{vInput[0], vInput[1], vInput[2], result}
	}
	return values
}

func main() {
	// Генерування значень для функції
	values := valuesOfFunc([][3]float64{{1, 1, 1}, {0, 1, 1}, {1, 0, 1}, {0, 0, 1}, {1, 1, 0}, {0, 1, 0}, {1, 0, 0}, {0, 0, 0}})
	fmt.Println(values)
	fmt.Println()

	inputs := make([][3]float64, len(values))
	targets := make([]float64, len(values))

	for i, val := range values {
		inputs[i] = [3]float64{val[0], val[1], val[2]}
		targets[i] = val[3]
	}

	// Ініціалізація бінарного перцептрона
	perceptron := NewBinaryPerceptron(3)

	// Навчання бінарного перцептрона
	perceptron.train(inputs, targets, 0.1, 20)

	// Тестування бінарного перцептрона на нових даних
	newData := [][3]float64{{1, 1, 1}, {0, 0, 0}, {1, 1, 0}, {0, 1, 1}}
	for _, data := range newData {
		_, prediction := perceptron.predict(data)
		fmt.Println("Prediction for", data, ":", prediction)
	}
}
