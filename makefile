.PHONY: snapshot
snapshot:
	scaffold --output-dir=":memory:" new --preset="test" --no-prompt --snapshot="stdout" ./ 

.PHONY: test/snapshot
test/snapshot:
	scaffold --output-dir=":memory:" new --preset="test" --no-prompt --snapshot="stdout" ./ | diff -u snapshots/test.snapshot - 

.PHONY: test/snapshot/update
test/snapshot/update:
	scaffold --output-dir=":memory:" new --preset="test" --no-prompt --snapshot="stdout" ./ > snapshots/test.snapshot

.PHONY: test/run
test/run:
	rm -rf /tmp/scaffold-test/
	# render output
	scaffold --output-dir="/tmp/scaffold-test/" new --preset="test" --no-prompt ./

	ls /tmp/scaffold-test/cli-test/

	# run binary, output should be "Hello World!"
	cd /tmp/scaffold-test/cli-test && \
		output=$$(go run main.go hello); \
		if [ "$$output" = "Hello World!" ]; then \
			echo "Output is correct: $$output"; \
		else \
			echo "Output is incorrect: $$output"; \
			exit 1; \
		fi

	
