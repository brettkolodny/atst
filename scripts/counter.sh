#!/bin/bash

# Check if two to four arguments are provided
if [ $# -lt 2 ] || [ $# -gt 4 ]; then
    echo "Usage: $0 <start_number> <end_number> [delay_ms] [exit_code]"
    exit 1
fi

# Assign arguments to variables
x=$1
y=$2
delay_ms=${3:-0}  # Default to 0 if not provided
exit_code=${4:-0}  # Default to 0 if not provided

# Check if inputs are valid numbers
if ! [[ "$x" =~ ^-?[0-9]+$ ]] || ! [[ "$y" =~ ^-?[0-9]+$ ]] || ! [[ "$delay_ms" =~ ^[0-9]+$ ]] || ! [[ "$exit_code" =~ ^[0-9]+$ ]]; then
    echo "Error: start and end must be integers, delay and exit_code must be positive integers"
    exit 1
fi

# Count from x to y
if [ $x -le $y ]; then
    # Count upward
    for ((i=x; i<=y; i++)); do
        echo $i
        if [ $delay_ms -gt 0 ]; then
            sleep $(awk "BEGIN {print $delay_ms/1000}")
        fi
    done
else
    # Count downward
    for ((i=x; i>=y; i--)); do
        echo $i
        if [ $delay_ms -gt 0 ]; then
            sleep $(awk "BEGIN {print $delay_ms/1000}")
        fi
    done
fi

# Exit with error code
exit $exit_code
