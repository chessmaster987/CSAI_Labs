import numpy as np


class GenAlg:
    def __init__(self, initialValues, target, amountOfChromosomes):
        # Ініціалізуємо початкові значення класу
        self.initialValues = initialValues  # Початкові значення хромосом
        self.target = target  # Цільове значення, яке ми шукаємо
        self.amountOfChromosomes = amountOfChromosomes  # Кількість хромосом у популяції

    def run(self):
        # Запускаємо головний алгоритм
        print("Generating start chromosomes")
        # Генеруємо початкову популяцію хромосом
        chromosomes = self.generate_start_chromosomes()
        for arr in chromosomes:
            print(arr)

        result = self.alg(chromosomes, float('inf'), 0)
        print("Result is", result)
        return result

    def alg(self, chromosomes, prevAvg, step):
        # Головний алгоритм генетичного пошуку
        print("Step", step)
        if step > 10000:
            print("No solution was found")
            return None

        # Оцінка значень функції для кожної хромосоми
        valuesOfFunction = np.array(
            [self.calc_function(arr) for arr in chromosomes])
        if self.target in valuesOfFunction:
            return chromosomes[np.where(valuesOfFunction == self.target)[0][0]]

        # Обчислюємо середнє значення функції
        avgFunc = np.mean(valuesOfFunction)
        print("\nValues of function:", valuesOfFunction)
        print("Average function value is:", avgFunc)
        # Проводимо мутації, якщо середнє значення функції зростає
        if avgFunc > prevAvg:
            print("Avg > previous avg. Some mutation")
            chromosomes = self.mutation(chromosomes)
            valuesOfFunction = np.array(
                [self.calc_function(arr) for arr in chromosomes])
            if self.target in valuesOfFunction:
                return chromosomes[np.where(valuesOfFunction == self.target)[0][0]]

        # Обчислюємо коефіцієнти для вибору батьків
        coef = self.calc_coef(valuesOfFunction)
        print("\nCoefficients:", coef)

        # Генеруємо нових нащадків
        childrens = self.get_children(chromosomes, coef)
        print("\nChildrens:")
        for arr in childrens:
            print(arr)
        # Рекурсивно запускаємо головний алгоритм для нової популяції
        return self.alg(childrens, avgFunc, step + 1)

    def generate_start_chromosomes(self):
        # Генерація випадкових початкових хромосом
        chromosomes = np.random.randint(
            0, self.target + 1, size=(self.amountOfChromosomes, len(self.initialValues)))
        return chromosomes

    def calc_function(self, values):
        # Обчислення значення функції для заданих значень хромосом
        return np.sum(np.array(self.initialValues) * np.array(values))

    def calc_coef(self, values):
        # Обчислення коефіцієнтів для вибору батьків
        r = 1.0 / np.abs(self.target - values)
        rSum = np.sum(r)
        coef = r / rSum
        return coef

    def get_children(self, chromosomes, coefficients):
        # Генерація нових нащадків
        sortedCoef = np.sort(coefficients)[::-1]
        sortedIndices = np.argsort(coefficients)[::-1]

        chromosomes = chromosomes[sortedIndices]
        uniqParent = set()
        while len(uniqParent) != len(sortedCoef):
            parent1Index = np.random.choice(
                range(len(sortedCoef)), p=sortedCoef)
            p_without_parent1 = np.delete(sortedCoef, parent1Index)
            # Нормалізуємо ймовірності
            p_without_parent1 /= np.sum(p_without_parent1)
            parent2Index = np.random.choice(
                range(len(sortedCoef) - 1), p=p_without_parent1)

            if parent1Index == parent2Index:
                continue
            parent1 = tuple(chromosomes[parent1Index])
            parent2 = tuple(chromosomes[parent2Index])
            uniqParent.add((parent1, parent2))

        children = np.array([self.crossover(parents)
                            for parents in uniqParent])
        return children

    def crossover(self, chromosomesParents):
        # Схрещування батьків для створення нащадка
        crossoverPoint = np.random.randint(1, len(chromosomesParents[0]) + 1)

        child = np.zeros(len(chromosomesParents[0]), dtype=int)
        child[:crossoverPoint] = chromosomesParents[0][:crossoverPoint]
        child[crossoverPoint:] = chromosomesParents[1][crossoverPoint:]

        return child

    def mutation(self, chromosomes):
        # Мутація хромосоми
        numOfChanges = np.random.randint(1, len(chromosomes) + 1)
        for _ in range(numOfChanges):
            i = np.random.randint(len(chromosomes))
            j = np.random.randint(len(chromosomes[i]))
            chromosomes[i][j] = np.random.randint(0, self.target + 1)
        return chromosomes


# Приклад використання для мого варіанту:
initialValues = [2, 8, 5, 7, -6]
target = 7
amountOfChromosomes = 20
genAlg = GenAlg(initialValues, target, amountOfChromosomes)
result = genAlg.run()
