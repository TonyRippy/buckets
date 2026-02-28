# Test Cases

This direcory contains a set of standard test cases that will be run to ensure
consistent behavior across implementations in different programming languages.

Files:
- `parse_errors.csv`: Invalid strategy specs and required error substrings.
- `*_parse.csv`: Parse behaviors for different bucketing strategies. Includes both success and failure cases.
- `*_index.csv`: `IndexOf` behavior for different bucketing strategies.
- `*_range.csv`: `Range` behavior for different bucketing strategies.
- `alignment_string.csv`: Verifying the string representation of bucket alignment.
