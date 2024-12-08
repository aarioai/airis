package debug

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateRmlogScript 生成删除日志的脚本
func generateRmlogScript(targetDir string) error {
	// 这里一般仅使用一次，不用移到常驻常量，浪费内存
	const scriptContent = `#!/bin/bash

# Show usage instructions
show_usage() {
    echo "Usage: $0"
    echo "    -d delete logs before date (YYYY-MM-DD)"
    echo "    -f log file format, default: %Y-%m-%d.log"
    echo "    -h show help"
    echo "Examples:"
    echo "  $0 -d 2024-03-01            # Delete logs before 2024-03-01"
    echo "  $0 -d day|week|month        # Delete logs before one day/week/month ago"
    echo "  $0 -f panic-%Y-%m-%d.log  # Delete logs with format panic-YYYY-MM-DD.log"
    echo "  $0                          # Delete logs before one month ago"
    exit 1
}

# Default values
BEFORE_DATE=""
FILE_NAME_FORMAT="%Y-%m-%d.log"
# Parse command line arguments
while getopts "d:f:h" opt; do
    case $opt in
        d) BEFORE_DATE="$OPTARG" ;;
        f) FILE_NAME_FORMAT="$OPTARG" ;;
        h) show_usage ;;
        ?) show_usage ;;
    esac
done

if [ -z "$BEFORE_DATE" ]; then
    show_usage
    exit 1
fi

case $BEFORE_DATE in
    day) BEFORE_DATE=$(date -d "1 day ago" +%Y-%m-%d) ;;
    week) BEFORE_DATE=$(date -d "1 week ago" +%Y-%m-%d) ;;
    month) BEFORE_DATE=$(date -d "1 month ago" +%Y-%m-%d) ;;
esac

# Validate date format
if ! date -d "$BEFORE_DATE" >/dev/null 2>&1; then
    echo "[error] invalid date format '$BEFORE_DATE'. Please use YYYY-MM-DD format"
    exit 1
fi

# Show operation details
echo "rm $FILE_NAME_FORMAT before $BEFORE_DATE"
echo "continue? [y/N]"
read -r confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo "operation cancelled"
    exit 0
fi

file_extension="${FILE_NAME_FORMAT##*.}"
date_pattern="${FILE_NAME_FORMAT%.*}"

# Convert date pattern to regex pattern
regex_pattern=$(echo "$date_pattern" | sed -e 's/%Y/[0-9]\{4\}/g' \
                                         -e 's/%m/[0-9]\{2\}/g' \
                                         -e 's/%d/[0-9]\{2\}/g')

# Find and delete log files before specified date
for file in *."$file_extension"; do
    [[ -f "$file" ]] || continue

    filename_no_ext="${file%.*}"

    if ! [[ "$filename_no_ext" =~ ^${regex_pattern}$ ]]; then
        continue
    fi

    parsed_date=$(date -d "$(echo "$filename_no_ext" | sed -E \
        -e "s/.*([0-9]{4})-([0-9]{2})-([0-9]{2}).*/\1-\2-\3/")" "+%Y-%m-%d" 2>/dev/null)
    
    if [ $? -ne 0 ]; then
        continue
    fi

    if [[ "$standard_date" < "$BEFORE_DATE" ]]; then
        if rm "$file"; then
            echo "rm $file"
        else
            echo "[error] rm $file"
        fi
    fi
done`

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create target dir failed: %w", err)
	}

	scriptPath := filepath.Join(targetDir, "rmlog.sh")
	if _, err := os.Stat(scriptPath); err == nil {
		return nil
	}
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("write script failed: %w", err)
	}

	return nil
}
