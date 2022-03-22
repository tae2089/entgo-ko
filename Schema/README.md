# Introduction
Schema는  그래프에서 한 Entity로 정의됩니다. <br />
Schema는  field,edge를 갖고 있으며 field는 엔티티의 속성을 나타내며 edge는 relation이라고 생각하시면 됩니다.

Schema를 생성하는 방법은 간단하며 다음과 같습니다
```shell
go run entgo.io/ent/cmd/ent init <생성할 스키마명>
```
