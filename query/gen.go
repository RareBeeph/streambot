// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q                 = new(Query)
	Message           *message
	RegisteredCommand *registeredCommand
	Subscription      *subscription
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Message = &Q.Message
	RegisteredCommand = &Q.RegisteredCommand
	Subscription = &Q.Subscription
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:                db,
		Message:           newMessage(db, opts...),
		RegisteredCommand: newRegisteredCommand(db, opts...),
		Subscription:      newSubscription(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Message           message
	RegisteredCommand registeredCommand
	Subscription      subscription
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:                db,
		Message:           q.Message.clone(db),
		RegisteredCommand: q.RegisteredCommand.clone(db),
		Subscription:      q.Subscription.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:                db,
		Message:           q.Message.replaceDB(db),
		RegisteredCommand: q.RegisteredCommand.replaceDB(db),
		Subscription:      q.Subscription.replaceDB(db),
	}
}

type queryCtx struct {
	Message           IMessageDo
	RegisteredCommand IRegisteredCommandDo
	Subscription      ISubscriptionDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Message:           q.Message.WithContext(ctx),
		RegisteredCommand: q.RegisteredCommand.WithContext(ctx),
		Subscription:      q.Subscription.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
