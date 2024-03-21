import numpy as np
from tensorflow import keras
from tensorflow.keras import layers

# Generate synthetic dataset
num_samples = 1000
image_height, image_width, num_channels = 32, 32, 3
num_classes = 10

x_train = np.random.random((num_samples, image_height, image_width, num_channels))
y_train = np.random.randint(num_classes, size=(num_samples,))

# Define CNN model
model = keras.Sequential([
    layers.Conv2D(32, (3, 3), activation='relu', input_shape=(image_height, image_width, num_channels)),
    layers.MaxPooling2D((2, 2)),
    layers.Conv2D(64, (3, 3), activation='relu'),
    layers.MaxPooling2D((2, 2)),
    layers.Conv2D(64, (3, 3), activation='relu'),
    layers.Flatten(),
    layers.Dense(64, activation='relu'),
    layers.Dense(num_classes)
])

# Compile the model
model.compile(optimizer='adam',
              loss=keras.losses.SparseCategoricalCrossentropy(from_logits=True),
              metrics=['accuracy'])

# Train the model in batches with epochs
batch_size = 32
num_epochs = 20

model.fit(x_train, y_train, batch_size=batch_size, epochs=num_epochs)
