.PHONY: all build install run clean full

# Define the name of your executable
APP_NAME = WebParseCLI

# Default target to build the project
all: build

# Build the project
build:
	go build -o $(APP_NAME) main.go

# Install the project
install:
	go install ./...

# Run the CLI command
run:
	./$(APP_NAME) fulbo

# Clean up build artifacts
clean:
	rm -f $(APP_NAME)

# Full target
full: build install run clean
