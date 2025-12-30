# 아파트 스코어링 프로젝트 Makefile

.PHONY: help clean build test clean-comments clean-comments-single restore-backups clean-all

# 기본 타겟
help:
	@echo "사용 가능한 명령어들:"
	@echo "  build              - 프로젝트 빌드"
	@echo "  test               - 테스트 실행"
	@echo "  clean-comments     - 모든 Go 파일에서 주석 제거"
	@echo "  clean-comments-single FILE=<file> - 특정 파일에서 주석 제거"
	@echo "  restore-backups    - 백업 파일에서 원본 복원"
	@echo "  clean              - 임시 파일 정리"
	@echo "  clean-all          - 모든 백업과 임시 파일 정리"

# 프로젝트 빌드
build:
	go build -o apart_score ./cmd

# 테스트 실행
test:
	go test ./...

# 모든 Go 파일에서 주석 제거
clean-comments:
	@echo "모든 Go 파일에서 주석 제거 중..."
	@find . -name "*.go" -not -path "./vendor/*" -not -name "*_test.go" | while read file; do \
		echo "처리 중: $$file"; \
		./remove_comments.sh "$$file"; \
	done
	@echo "완료! 백업 파일들은 *.bak 확장자로 저장되었습니다."

# 특정 파일에서 주석 제거
clean-comments-single:
	@if [ -z "$(FILE)" ]; then \
		echo "사용법: make clean-comments-single FILE=<파일경로>"; \
		exit 1; \
	fi
	@if [ ! -f "$(FILE)" ]; then \
		echo "파일이 존재하지 않습니다: $(FILE)"; \
		exit 1; \
	fi
	@echo "$(FILE)에서 주석 제거 중..."
	./remove_comments.sh "$(FILE)"

# 백업 파일에서 원본 복원
restore-backups:
	@echo "백업 파일에서 원본 복원 중..."
	@find . -name "*.bak" | while read backup; do \
		original="$${backup%.bak}"; \
		if [ -f "$$backup" ]; then \
			echo "복원 중: $$original"; \
			cp "$$backup" "$$original"; \
			rm "$$backup"; \
		fi; \
	done
	@echo "복원 완료!"

# 주석 제거 후 빌드 테스트
test-build:
	@echo "주석 제거 후 빌드 테스트..."
	@make clean-comments
	@make build
	@echo "빌드 성공! 주석 제거가 올바르게 작동했습니다."
	@make restore-backups

# 모든 백업과 임시 파일 정리
clean-all:
	@echo "모든 백업과 임시 파일 정리 중..."
	@find . -name "*.bak" -delete
	@find . -name "*.tmp" -delete
	@find . -name "apart_score" -delete
	@echo "정리 완료!"

# 일반 정리
clean:
	
# Clean target: Remove comments and apply goimports (includes gofmt + import cleanup)
clean:
	@echo "🧹 Starting code cleanup process..."
	@echo "📝 Step 1: Removing comments from Go files..."

	# Find all .go files, exclude test files and backup files
	@find . -name "*.go" \
		-not -name "*_test.go" \
		-not -name "*.bak" \
		-not -path "./domain/district_test/*" \
		-not -path "./application/services/test/*" \
		-not -path "./infrastructure/repository/mongodb/test/*" \
		-not -path "./domain/*/test/*" | while read -r file; do \
		echo "Processing: $$file"; \
		./remove_comments.sh "$$file"; \
	done

	@echo "🎨 Step 2: Applying goimports (includes gofmt + import cleanup)..."
	@goimports -w .

	@echo "✨ Step 3: Checking for any remaining issues..."
	@go vet ./...

	@echo "🗑️  Step 4: Removing backup files..."
	@find . -name "*.bak" -type f -delete

	@echo "✅ Code cleanup completed successfully!"
	@echo "   - Comments removed from all Go files"
	@echo "   - Code formatted with goimports"
	@echo "   - Imports cleaned up automatically"
	@echo "   - Static analysis passed"
	@echo "   - Backup files cleaned up"

	@find . -name "*.tmp" -delete
	@find . -name "apart_score" -delete
