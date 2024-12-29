run_test_case () {
    echo "Test case: $TEST_CASE_NAME"

    local ACTUAL_STDERR_OUTPUT
    ACTUAL_STDERR_OUTPUT=$(echo -n "$FILE_CONTENT" | "./csv-validator" 2>&1)
    local ACTUAL_EXIT_CODE=$?

    if [[ "$EXPECTED_EXIT_CODE" != "$ACTUAL_EXIT_CODE" ]]; then
        echo "Failed: Expected exit code to be '$EXPECTED_EXIT_CODE', got '$ACTUAL_EXIT_CODE'"
    elif [[ "$EXPECTED_STDERR_OUTPUT" != "$ACTUAL_STDERR_OUTPUT" ]]; then
        echo "Failed: Expected output to be '$EXPECTED_STDERR_OUTPUT', got '$ACTUAL_STDERR_OUTPUT'"
    else
        echo "Passed"
    fi

    echo ""
}

echo "Running test suite..."
echo ""

TEST_CASE_NAME="Empty file is invalid"
EXPECTED_EXIT_CODE=2
EXPECTED_STDERR_OUTPUT="line 1 column 0: expected at least 1 column"
FILE_CONTENT=""

run_test_case

TEST_CASE_NAME="Inconsistent number of fields is invalid"
EXPECTED_EXIT_CODE=2
EXPECTED_STDERR_OUTPUT="line 2 column 9: expected 3 fields, found 2"
FILE_CONTENT="John,Doe,30
Jane,Doe"

run_test_case

echo "...Done"
