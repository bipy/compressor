package constant

import (
	"strconv"
	"time"
)

var ID = strconv.FormatInt(time.Now().Unix(), 10)
