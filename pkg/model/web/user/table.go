package user

type TableConfig interface {
	Table() TableData
}

const ID = "id"

type TableData struct {
	Name       string
	PrimaryKey string
	Alias      string
	Uniques    []string
	Indexes    []string
	Joins      []Join
}
type Join struct {
	Table      interface{}
	TableName  string
	ForeignKey string
	Reference  string
}

func (m *User) Table() TableData {
	return TableData{Name: "users", Alias: "us"}
}
