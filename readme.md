###Ahat Simple Gallary

```
모바일에서 사용하기 적합한 간단한 이미지 갤러리 서버입니다.

Ahat Simple Gallary를 실행시키고 
모바일 웹브라우저 혹은 데스크탑 웹 브라우저로 접속하여
이미지 갤러리 기능을 사용할 수 있습니다.
현재까지 구현된 기능은 아래와 같습니다.
```

```
 - 이미지를 썸네일화 하여 갤러리 목록에 표시
 - 썸네일을 클릭할 경우 원본 이미지를 표시
 - 이미지를 좌 우로 드래그하여 이전/다음 이미지 표시
 - 이미지를 선택 후 이미지 경로 이동 기능
 - 이미지를 선택 후 이미지 삭제 기능
 - 페이지 접속 시 로그인을 하여야 접속 간으
 - 이미지의 이름/크기/날짜별 정렬 기능
```

```
컴파일 방법
 go get github.com/disintegration/imaging
 go get gopkg.in/ini.v1
 go get github.com/gorilla/securecookie
 go build .\ImageCloud.go .\Utils.go .\login.go
```

```
사용방법 1. 윈도우 환경
1) images 파일에 이용한 이미지 파일을 구성한다.
2) ImageCloud.exe 파일을 실행한다.
2-1) 로그파일을 남기고 싶을 경우 CMD 창을 열고 바이너리 위치로 이동하여 ` ImageCloud.exe >> ImageCloud.log ` 형식으로 실행한다.
3) 프로그램이 실행된 후 모바일 브라우저 혹은 브라우저로 접속하여 사용한다.
ex) http://127.0.0.1:9090/
```

