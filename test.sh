#!/bin/bash

# Path to your_program.sh
PROGRAM_PATH="./your_program.sh"

# Directory containing test cases
TESTS_DIR="./tests"

ALL_TESTS_PASSED=true
PASS_COUNT=0
ALL_COUNT=0

# Loop through each test directory in the tests folder
for TEST_DIR in "$TESTS_DIR"/*; do
    FAILED=false


    # Define the path for the test input and expected output
    TEST_INPUT="$TEST_DIR/test.lox"
    EXPECTED_OUTPUT="$TEST_DIR/expected"
    EXPECTED_STDERR="$TEST_DIR/expected_stderr"

    # Run your_program.sh with the test input and capture the output
    ACTUAL_OUTPUT=$(sh "$PROGRAM_PATH" tokenize "$TEST_INPUT" 2>/dev/null)
    ACTUAL_OUTPUT_STDERR=$(sh "$PROGRAM_PATH" tokenize "$TEST_INPUT" 2>&1 >/dev/null)

    echo "Running test $(basename "$TEST_DIR")..."
    # Compare actual output with expected output
    if diff <(echo "$ACTUAL_OUTPUT") "$EXPECTED_OUTPUT" >/dev/null; then
        echo "[STDOUT] Success: Test $(basename "$TEST_DIR") passed."
    else
        echo "[STDOUT] Error: Test $(basename "$TEST_DIR") failed:"
        # Show the diff
        diff --color -y <(echo "$ACTUAL_OUTPUT") "$EXPECTED_OUTPUT"
        FAILED=true
    fi

    # if EXPECTED_STDERR exists, compare it with ACTUAL_OUTPUT_STDERR
    # else, ACTUAL_OUTPUT_STDERR should be empty
    if [ -f "$EXPECTED_STDERR" ]; then
        if diff <(echo "$ACTUAL_OUTPUT_STDERR") "$EXPECTED_STDERR" >/dev/null; then
            echo "[STDERR] Success: Test $(basename "$TEST_DIR") passed."
        else
            echo "[STDERR] Error: Test $(basename "$TEST_DIR") failed:"
            # Show the diff
            diff --color -y <(echo "$ACTUAL_OUTPUT_STDERR") "$EXPECTED_STDERR"
            ALL_TESTS_PASSED=false
        fi
    else
        if [ -n "$ACTUAL_OUTPUT_STDERR" ]; then
            echo "[STDERR] Error: Test $(basename "$TEST_DIR") failed:"
            echo "Expected empty stderr, but got:"
            echo "$ACTUAL_OUTPUT_STDERR"
        FAILED=true
        fi
    fi

    if [ "$FAILED" = true ]; then
        ALL_TESTS_PASSED=false
    else
        PASS_COUNT=$((PASS_COUNT + 1))
    fi
    ALL_COUNT=$((ALL_COUNT + 1))
done

echo "Pass: $PASS_COUNT, Fail: $((ALL_COUNT - PASS_COUNT))"

# Exit with error if any test failed
if [ "$ALL_TESTS_PASSED" = false ]; then
    echo "Some tests failed."
    exit 1
else
    echo "All tests passed."
fi
