import numpy as np

class BinaryPerceptron:
    def __init__(self, input_size):
        self.input_size = input_size
        self.weights = np.array([0.2, 0.02, 0.002])

    def predict(self, inputs):
        weighted_sum = np.dot(inputs, self.weights)
        prediction = 1 if weighted_sum > 0 else 0
        return round(weighted_sum, 2), prediction

    def train(self, inputs, targets, learning_rate, epochs):
        epoch_i = 0
        for epoch in range(epochs):
            epoch_i += 1
            print("Iteration #", epoch_i)
            i = 0
            is_error = 0
            for input_data, target in zip(inputs, targets):
                sum_, prediction = self.predict(input_data)
                error = target - prediction
                if error != 0:
                    is_error = 1
                prev_weights = self.weights.copy()
                self.weights += learning_rate * error * input_data
                i += 1
                print(f'#{i} w_1={round(prev_weights[0], 2)} w_2={round(prev_weights[1], 2)}'
                      f' w_3={round(prev_weights[2], 2)} x_1={input_data[0]} x_2={input_data[1]} x_3={input_data[2]}'
                      f' a={sum_} Y={prediction} T={target} error={error}'
                      f' delta_w_1={round(learning_rate * error * input_data[0], 2)}'
                      f' delta_w_2={round(learning_rate * error * input_data[1], 2)}'
                      f' delta_w_3={round(learning_rate * error * input_data[2], 2)}')
            print()
            if is_error == 0:
                print("Training is over")
                break

def values_of_func(X):
    values = []
    for v_input in X:
        values.append((v_input[0], v_input[1], v_input[2], v_input[0] and not v_input[1] or v_input[2]))
    return values

values = values_of_func([[1, 1, 1], [0, 1, 1], [1, 0, 1], [0, 0, 1], [1, 1, 0], [0, 1, 0], [1, 0, 0], [0, 0, 0]])
print(values)
print()

inputs = []
targets = []

for val in values:
    inputs.append(val[:3])
    targets.append(val[3])

perceptron = BinaryPerceptron(input_size=3)

inputs = np.array(inputs)
targets = np.array(targets)

perceptron.train(inputs, targets, learning_rate=0.1, epochs=20)

# Testing data
new_data = np.array([[1, 1, 1], [0, 0, 0], [1, 1, 0], [0, 1, 1]])
for data in new_data:
    _, prediction = perceptron.predict(data)
    print("Prediction for", data, ":", prediction)
