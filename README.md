## [PROJECT] Golang ORM - ENTGO 한국어 버전으로 정리하기
### 프로젝트 관련 설명
- entgo를 한국어로 번역하는 작업
- 간단한 테스트 코드도 작성하기
### 프로젝트 사용된 기술 스택
- entgo
- mock
- suite
### 세팅 및 다운로드
```Shell
go mod init entgo-ko
go get -d entgo.io/ent/cmd/ent
```

### 스키마 생성 방법
```Shell
go run entgo.io/ent/cmd/ent init <Shema 명> # Shema 생성
go run entgo.io/ent/cmd/ent init User # User 생성
```
