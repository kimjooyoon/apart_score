# 아파트 스코어링 시스템 (Apart Score)

## 프로젝트 개요

외부에서 수집된 아파트 정보를 기반으로 DDD(Domain-Driven Design) 기반의 스코어링 시스템입니다.
사용자가 설정한 메타데이터 기준에 따라 아파트 점수를 계산하고 평가합니다.

## 아키텍처

- **언어**: Go
- **패턴**: DDD (Domain-Driven Design)
- **구조**: 모놀리스 (Monolith)
- **패키지 관리**: Go Modules

## 메타데이터 스코어링

아파트의 점수는 다음과 같은 메타데이터 요소들을 기반으로 계산됩니다:

### 주요 메타데이터 후보
- **층수 (Floor Level)**: 중간층에 가까울수록 높은 점수
- **역까지 거리 (Distance to Station)**: 가까울수록 높은 점수
- **엘리베이터 유무 (Elevator Presence)**: 있으면 높은 점수
- **건축년도 (Construction Year)**: 최신 건물일수록 높은 점수
- **건설회사 (Construction Company)**: 신뢰할 수 있는 회사일수록 높은 점수
- **아파트 크기 (Apartment Size)**: 적절한 크기일수록 높은 점수
- **주변 편의시설 (Nearby Amenities)**: 편의시설이 많을수록 높은 점수
- **교통 접근성 (Transportation Access)**: 대중교통 접근성이 좋을수록 높은 점수
- **학군 (School District)**: 좋은 학군일수록 높은 점수
- **범죄율 (Crime Rate)**: 낮을수록 높은 점수
- **녹지율 (Green Space Ratio)**: 높을수록 높은 점수
- **주차장 (Parking)**: 주차 공간이 충분할수록 높은 점수
- **관리비 (Maintenance Fee)**: 적절한 수준일수록 높은 점수
- **난방 방식 (Heating System)**: 효율적인 난방 방식일수록 높은 점수

### 스코어링 예시
```
용인시 A아파트:
- 최대층: 10층, 현재층: 5층 (중간층) → +100점
- 역까지 거리: 5분 → +100점
- 엘리베이터: 있음 → +100점
- 건축년도: 2020년 → +90점
총합: 390점
```

## 프로젝트 구조

```
apart_score/
├── pkg/
│   ├── metadata/
│   │   ├── constants.go    # 메타데이터 상수 정의
│   │   ├── types.go        # 메타데이터 타입 정의
│   │   └── metadata.go     # 메타데이터 메소드
├── internal/
│   ├── domain/
│   │   ├── apartment/      # 아파트 도메인
│   │   ├── scoring/        # 스코어링 도메인
│   │   └── user/           # 사용자 도메인
│   ├── application/        # 애플리케이션 서비스
│   └── infrastructure/     # 인프라스트럭처
├── cmd/
│   └── main.go             # 메인 애플리케이션
└── go.mod
```

## 메타데이터 설계 원칙

1. **수정 불가**: 한 번 정의된 메타데이터는 수정되지 않음 (오직 추가만 가능)
2. **인덱스 기반**: iota를 사용한 순차적 인덱스 번호 부여
3. **다국어 지원**: 영문명과 한글명 모두 지원
4. **외부 참조**: 다른 패키지에서 쉽게 참조 가능

## 시작하기

```bash
# 프로젝트 초기화
go mod init apart_score

# 의존성 설치
go mod tidy

# 빌드 및 실행
go build ./cmd
./apart_score
```

## 개발 가이드라인

- DDD 패턴 준수
- 메타데이터는 pkg 패키지를 통해서만 접근
- 모든 메타데이터 추가는 constants.go에 iota로 순차적 추가
- 테스트 코드 작성 필수
