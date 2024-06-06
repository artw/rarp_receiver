# Define variables
DOCKER_IMAGE := rpmbuild
CONTAINER_NAME := rpmbuild
ARCHIVE_NAME := rarp_receiver-1.0.0
INPUT_DIR := $(shell pwd)/input
OUTPUT_DIR := $(shell pwd)/output

.PHONY: all build clean

all: build

build:
	# Create output directory if it does not exist
	mkdir -p $(INPUT_DIR) $(OUTPUT_DIR)
	git archive --format=tar.gz --prefix=$(ARCHIVE_NAME) -o $(INPUT_DIR)/$(ARCHIVE_NAME).tar.gz HEAD
	# Build the Docker image
	cd $(TEMP_DIR)
	docker build -t $(DOCKER_IMAGE) .
	# Run the Docker container to build the RPM, with mounts to export the RPM
	docker run --name $(CONTAINER_NAME) --rm -v $(SOURCE_DIR):/root/rpmbuild/SOURCES -v $(OUTPUT_DIR):/root/rpmbuild/RPMS/x86_64 $(DOCKER_IMAGE)

clean:
	# Clean the output directory
	rm -rf $(OUTPUT_DIR) $(INPUT_DIR)

