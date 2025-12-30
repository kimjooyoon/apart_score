#!/bin/bash

# Go 코드에서 주석을 제거하는 스크립트
# 사용법: ./remove_comments.sh <input_file> [output_file]

if [ $# -lt 1 ]; then
    echo "사용법: $0 <input_file> [output_file]"
    exit 1
fi

input_file="$1"

if [ ! -f "$input_file" ]; then
    echo "파일이 존재하지 않습니다: $input_file"
    exit 1
fi

# 백업 파일 생성
backup_file="${input_file}.bak"
cp "$input_file" "$backup_file"

echo "주석 제거 중: $input_file (백업: $backup_file)"

# awk를 사용하여 Go 주석 제거
# /* */ 블록 주석과 // 주석을 제거 (인라인 주석 포함)
# 단, 문자열 안의 주석은 유지
awk '
BEGIN {
    in_block_comment = 0
    in_string = 0
    string_char = ""
}

{
    line = $0
    result = ""

    for (i = 1; i <= length(line); i++) {
        char = substr(line, i, 1)
        next_char = (i < length(line)) ? substr(line, i+1, 1) : ""

        # 문자열 시작/종료 처리
        if (!in_block_comment) {
            if ((char == "\"" || char == "`") && (i == 1 || substr(line, i-1, 1) != "\\")) {
                if (!in_string) {
                    in_string = 1
                    string_char = char
                } else if (char == string_char) {
                    in_string = 0
                    string_char = ""
                }
            }
        }

        # 문자열 밖에서만 주석 처리
        if (!in_string) {
            # 블록 주석 시작
            if (char == "/" && next_char == "*" && !in_block_comment) {
                in_block_comment = 1
                i++  # * 건너뛰기
                continue
            }

            # 블록 주석 종료
            if (char == "*" && next_char == "/" && in_block_comment) {
                in_block_comment = 0
                i++  # / 건너뛰기
                continue
            }

            # 모든 // 주석 (라인 중간 포함)
            if (char == "/" && next_char == "/") {
                break  # 나머지 라인 무시
            }
        }

        # 블록 주석 중이 아니면 결과에 추가
        if (!in_block_comment) {
            result = result char
        }
    }

    # 결과 라인의 trailing space 제거
    sub(/[ \t]+$/, "", result)

    # 빈 라인이 아니면 출력 (연속 빈 줄 방지)
    if (length(result) > 0) {
        print result
    }
}
' "$input_file" | sed '/^$/N;/^\n$/d' > "$input_file.tmp" && mv "$input_file.tmp" "$input_file"

echo "완료: $input_file (백업: $backup_file)"
