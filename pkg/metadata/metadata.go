package metadata

import "fmt"

// Index는 메타데이터 타입의 인덱스 번호를 반환합니다.
// 외부 패키지에서 참조할 수 있는 수정되지 않는 고유 번호입니다.
func (mt MetadataType) Index() int {
	return int(mt)
}

// String은 메타데이터 타입의 영문명을 반환합니다.
func (mt MetadataType) String() string {
	if mt.IsValid() {
		return metadataInfos[mt].englishName
	}
	return fmt.Sprintf("Unknown(%d)", int(mt))
}

// KoreanName은 메타데이터 타입의 한글명을 반환합니다.
func (mt MetadataType) KoreanName() string {
	if mt.IsValid() {
		return metadataInfos[mt].koreanName
	}
	return fmt.Sprintf("알 수 없음(%d)", int(mt))
}

// Description은 메타데이터 타입의 설명을 반환합니다.
func (mt MetadataType) Description() string {
	if mt.IsValid() {
		return metadataInfos[mt].description
	}
	return fmt.Sprintf("설명 없음(%d)", int(mt))
}

// IsValid는 메타데이터 타입이 유효한지 확인합니다.
func (mt MetadataType) IsValid() bool {
	return mt >= 0 && mt < MetadataTypeCount
}

