AOUT ?= gluebin
SRCFILES := $(shell find $(SRC) -name "*.go")
BUILD ?= ./build

$(BUILD)/$(AOUT): $(SRCFILES)
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(BUILD)/$(AOUT) cmd/gluebin/main.go

all: $(BUILD)/$(AOUT)

clean:
	$(RM) -r $(BUILD)
