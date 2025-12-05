package testdata

type UnProtectedEntity struct {
	StringField string
	IntField    int
	IntPtrField *int
}

type Entity struct {
	SubEntityViaProperty *SubEntity
	ProtectedField       string
}

func (e *Entity) SubEntity() *SubEntity {
	return &SubEntity{}
}

func (e *Entity) SubEntityViaInterface() SubEntityInterface {
	return &SubEntity{}
}

type SubEntityInterface interface {
	SetProtectedField(value string)
}

type SubEntity struct {
	ProtectedField string
}

func (s *SubEntity) SetProtectedField(value string) {
	s.ProtectedField = value
}

type Repository interface {
	Read() *Entity
}

type RepositoryImpl struct{}

func (r *RepositoryImpl) Read() *Entity {
	return &Entity{}
}

var repo Repository = &RepositoryImpl{}

func (e *Entity) SetProtectedField(value string) {
	e.ProtectedField = value
}

func SomeFunc1() {
	e := &Entity{}
	e.SetProtectedField("value")
}

func SomeFunc2() {
	e := &Entity{}
	e.ProtectedField = "value" // want "assignment to exported field Entity.ProtectedField is forbidden outside its methods"
}

func SomeFunc3() {
	e := repo.Read()
	e.ProtectedField = "value" // want "assignment to exported field Entity.ProtectedField is forbidden outside its methods"
}

func SomeFunc4() {
	e := repo.Read()
	e.SetProtectedField("value")
}

func SomeFunc5() {
	e := &Entity{}
	sub := e.SubEntity()
	sub.SetProtectedField("value")
}
func SomeFunc6() {
	e := &Entity{}
	sub := e.SubEntity()
	sub.ProtectedField = "value" // want "assignment to exported field SubEntity.ProtectedField is forbidden outside its methods"
}

func SomeFunc7() {
	e := &Entity{}
	sub := e.SubEntityViaInterface()
	sub.SetProtectedField("value")
}

func SomeFunc8() {
	e := &Entity{
		SubEntityViaProperty: &SubEntity{},
	}
	e.SubEntityViaProperty.ProtectedField = "value" // want "assignment to exported field SubEntity.ProtectedField is forbidden outside its methods"
}

func SomeFunc9() {
	e := &Entity{
		SubEntityViaProperty: &SubEntity{},
	}
	e.SubEntityViaProperty.SetProtectedField("value")
	if e.SubEntityViaProperty.ProtectedField != "value" {
	}
}

func SomeFunc10() {
	e := &UnProtectedEntity{}
	e.StringField = "value"
	e.IntField++
	e.IntField--
	e.IntField += 10
	e.IntField -= 10
	e.IntField *= 10
	e.IntField /= 10
	e.IntField = 10
	*(&e.IntField)++
	*(&e.IntField)--
	e.IntPtrField = new(int)
	*e.IntPtrField = 20
}
