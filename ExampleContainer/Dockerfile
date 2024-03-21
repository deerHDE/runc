# Use an official TensorFlow runtime as a parent image
FROM tensorflow/tensorflow:latest

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

WORKDIR /
RUN mkdir -p /tmp
RUN chmod 1777 -R /tmp
USER root

# Install any additional dependencies if required
# RUN pip install -r requirements.txt

# Command to run the Python script
WORKDIR /app
CMD ["python", "train_model.py"]
