package test

// 我们在操作数据库的时候，你如dao层中当遇到一个sql.ErrNoRows的时候，是否应该warp这个error,抛给上一层。为什么？应该怎么做，请写出你的代码？
//需要，实际生产中，调用者需要了解错误场景，因此为错误增加上下文信息后再返回非常有必要
import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func QueryInfo() error {
	return errors.Wrap(sql.ErrNoRows, "QueryInfo failed")
}

func Call() error {
	return errors.WithMessage(QueryInfo(), "Call failed")
}

func TestSqlErr(t *testing.T) {
	err := Call()
	if errors.Cause(err) == sql.ErrNoRows {
		fmt.Printf("Data not found, %v\n", err)
		fmt.Printf("%+v\n", err)
		return
	}
	if err != nil {
		// unknown error
	}
}
