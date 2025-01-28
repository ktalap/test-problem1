#!/bin/bash

# Check if dictionary exists
DICT="/usr/share/dict/words"
if [ ! -f "$DICT" ]; then
    echo "Dictionary not found at $DICT"
    exit 1
fi

# 20,000 bytes
target_size=20000

# Generate first file
echo "Generating first file..."
truncate -s 0 large_input1.txt

# Get total word count in dictionary
total_words=$(wc -l < "$DICT")

# Function to get random words with repetition
get_random_words() {
    # Get 100 random words at a time, with repetition
    for i in {1..100}; do
        # Get a random line number
        line_num=$((RANDOM % total_words + 1))
        # Get that line from dictionary
        sed -n "${line_num}p" "$DICT"
    done
}

# Generate first file
while [ $(stat --format=%s large_input1.txt) -lt $target_size ]; do
    get_random_words >> large_input1.txt
    echo "Current size: $(stat --format=%s large_input1.txt) bytes"
done

# Sort first file
echo "Sorting first file..."
sort large_input1.txt -o large_input1.txt

# Generate second file with 70% shared content
echo "Generating second file..."
# Take 70% of the first file
head -n $(( $(wc -l < large_input1.txt) * 7 / 10 )) large_input1.txt > large_input2.txt

# Add new random words to second file
while [ $(stat --format=%s large_input2.txt) -lt $target_size ]; do
    get_random_words >> large_input2.txt
    echo "Current size: $(stat --format=%s large_input2.txt) bytes"
done

# Sort second file
echo "Sorting second file..."
sort large_input2.txt -o large_input2.txt

# Print statistics
echo "File statistics:"
echo "First file:"
echo "Size: $(ls -lh large_input1.txt)"
echo "Lines: $(wc -l < large_input1.txt)"
echo "Unique words: $(sort -u large_input1.txt | wc -l)"
echo
echo "Second file:"
echo "Size: $(ls -lh large_input2.txt)"
echo "Lines: $(wc -l < large_input2.txt)"
echo "Unique words: $(sort -u large_input2.txt | wc -l)"