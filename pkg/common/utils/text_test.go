package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextNormalizer(t *testing.T) {
	t.Run(
		"normalizer", func(t *testing.T) {
			inputData := []map[string]string{
				{
					"input":  "Nguyễn Vàng Linh",
					"result": "nguyen vang linh",
				},
				{
					"input":  "Trần Văn Được",
					"result": "tran van duoc",
				},
				{
					"input":  "012345",
					"result": "012345",
				},
			}

			// expect 2 completed achievements only

			for _, item := range inputData {
				assert.Equal(
					t,
					item["result"],
					TextNormalizer(item["input"]),
				)
			}
		},
	)
}
