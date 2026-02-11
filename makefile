SHELL := /bin/bash

.PHONY: snapshot
snapshot:
	scaffold new --output-dir=":memory:" --preset="test" --no-prompt --snapshot="stdout" ./

DATE_RE := [0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}

.PHONY: test/snapshot
test/snapshot:
	diff -u \
		<(sed 's/$(DATE_RE)/DATE/g' snapshots/test.snapshot) \
		<(scaffold new --output-dir=":memory:" --preset="test" --no-prompt --snapshot="stdout" ./ | sed 's/$(DATE_RE)/DATE/g')

.PHONY: test/snapshot/update
test/snapshot/update:
	scaffold new --output-dir=":memory:" --preset="test" --no-prompt --snapshot="stdout" ./ > snapshots/test.snapshot

.PHONY: test/run
test/run:
	rm -rf /tmp/scaffold-test/
	# render output
	scaffold --log-level="debug" new --output-dir="/tmp/scaffold-test/" --preset="test" --no-prompt ./

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
