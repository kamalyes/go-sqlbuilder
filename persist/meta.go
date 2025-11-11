/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:13:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 21:15:03
 * @FilePath: \go-sqlbuilder\persist\query_param.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package persist

type Paging struct {
	Page   int
	Offset int
	Limit  int
	Total  int64
}

func NewPaging(page int, limit int) *Paging {
	offset := (page - 1) * limit
	return &Paging{Page: page, Offset: offset, Limit: limit}
}
