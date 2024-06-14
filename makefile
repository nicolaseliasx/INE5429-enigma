EXECUTABLE = enigma

GOFILE = enigma.go

run: $(GOFILE)
	go build -o $(EXECUTABLE) $(GOFILE)
	./$(EXECUTABLE)

clean:
	rm -f $(EXECUTABLE)
