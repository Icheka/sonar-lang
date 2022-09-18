package inputs

type SourceInput interface {
	Read(args []interface{}) string
}
