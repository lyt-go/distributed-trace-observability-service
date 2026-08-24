package tracepublish

import (
	"errors"

	"tracing/internal/tracejournal"
)

type Publisher interface{ Publish(string) error }
type Audit interface{ Success(string) }

// Publish 提交轨迹日志事务，成功后向外部发布一次并记一条成功审计。
//
// 提交是发布与审计的前置条件：提交失败则不发布、不记成功，并保留原始提交错误
// （回滚错误只作为次要错误追加，不得掩盖原始提交错误）。正常提交后发布一次、
// 留下一条成功审计。
func Publish(tx *tracejournal.Transaction, publisher Publisher, audit Audit, traceID string) (err error) {
	// 先提交轨迹日志：未提交成功前不得发布或记成功。
	if err = tx.Commit(); err != nil {
		// 提交失败，事务尚未持久化，回滚以清理；回滚错误作为次要错误追加，
		// 原始提交错误必须保留。
		if rb := tx.Rollback(); rb != nil {
			return errors.Join(err, rb)
		}
		return err
	}
	// 正常提交后发布一次。
	if err = publisher.Publish(traceID); err != nil {
		return err
	}
	// 发布成功后记一条成功审计。
	audit.Success(traceID)
	return nil
}
