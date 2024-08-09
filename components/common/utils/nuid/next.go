package nuid

var (
	id = New()
)

func NextId() string {
	return id.Next()
}
