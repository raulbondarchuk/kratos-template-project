package utils

func StrOrNull(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
