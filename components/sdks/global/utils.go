package global

import "fmt"

func GetReferenceId(prefix string, recordid int64) string {
	return fmt.Sprintf("%v#%v", prefix, recordid)
}
