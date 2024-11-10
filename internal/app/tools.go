package app

import (
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type IDArrayParser struct {
	key string
	c   *gin.Context
}

func NewIDArrayParser(key string, c *gin.Context) *IDArrayParser {
	return &IDArrayParser{key, c}
}

// 非法时返回空,不允许元素重复
func (p *IDArrayParser) Parse() []uint {
	ids := p.c.PostFormArray(p.key)
	log.Trace("parse:", ids)
	res := make([]uint, len(ids))
	for i := range ids {
		id, err := strconv.ParseUint(ids[i], 10, 0)
		if err != nil {
			log.Trace(err)
			return []uint{}
		}
		res[i] = uint(id)
	}

	slices.Sort(res)
	res = slices.Compact(res)
	return res
}
