.PHONY: test
test:
	scaffold --output-dir=./test new ./

snapshot:
	scaffold --output-dir=":memory:" new --preset="test" --no-prompt --snapshot="stdout" ./ 
