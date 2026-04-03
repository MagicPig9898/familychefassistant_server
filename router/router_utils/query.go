package router_utils

import "strconv"

func ParseInt64Query(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}
