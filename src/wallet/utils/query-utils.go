package utils

type queryutil struct {}

func QueryUtil() *queryutil {
	return &queryutil{}
}

func (u *queryutil) Build(field string, value string) string {
	return `{"`+field+`": "`+value+`"}`
}