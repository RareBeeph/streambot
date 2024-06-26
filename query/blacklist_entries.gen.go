// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"streambot/models"
)

func newBlacklistEntry(db *gorm.DB, opts ...gen.DOOption) blacklistEntry {
	_blacklistEntry := blacklistEntry{}

	_blacklistEntry.blacklistEntryDo.UseDB(db, opts...)
	_blacklistEntry.blacklistEntryDo.UseModel(&models.BlacklistEntry{})

	tableName := _blacklistEntry.blacklistEntryDo.TableName()
	_blacklistEntry.ALL = field.NewAsterisk(tableName)
	_blacklistEntry.ID = field.NewUint(tableName, "id")
	_blacklistEntry.CreatedAt = field.NewTime(tableName, "created_at")
	_blacklistEntry.UpdatedAt = field.NewTime(tableName, "updated_at")
	_blacklistEntry.DeletedAt = field.NewField(tableName, "deleted_at")
	_blacklistEntry.UserID = field.NewString(tableName, "user_id")
	_blacklistEntry.UserLogin = field.NewString(tableName, "user_login")
	_blacklistEntry.ChannelID = field.NewString(tableName, "channel_id")

	_blacklistEntry.fillFieldMap()

	return _blacklistEntry
}

type blacklistEntry struct {
	blacklistEntryDo

	ALL       field.Asterisk
	ID        field.Uint
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	UserID    field.String
	UserLogin field.String
	ChannelID field.String

	fieldMap map[string]field.Expr
}

func (b blacklistEntry) Table(newTableName string) *blacklistEntry {
	b.blacklistEntryDo.UseTable(newTableName)
	return b.updateTableName(newTableName)
}

func (b blacklistEntry) As(alias string) *blacklistEntry {
	b.blacklistEntryDo.DO = *(b.blacklistEntryDo.As(alias).(*gen.DO))
	return b.updateTableName(alias)
}

func (b *blacklistEntry) updateTableName(table string) *blacklistEntry {
	b.ALL = field.NewAsterisk(table)
	b.ID = field.NewUint(table, "id")
	b.CreatedAt = field.NewTime(table, "created_at")
	b.UpdatedAt = field.NewTime(table, "updated_at")
	b.DeletedAt = field.NewField(table, "deleted_at")
	b.UserID = field.NewString(table, "user_id")
	b.UserLogin = field.NewString(table, "user_login")
	b.ChannelID = field.NewString(table, "channel_id")

	b.fillFieldMap()

	return b
}

func (b *blacklistEntry) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := b.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (b *blacklistEntry) fillFieldMap() {
	b.fieldMap = make(map[string]field.Expr, 7)
	b.fieldMap["id"] = b.ID
	b.fieldMap["created_at"] = b.CreatedAt
	b.fieldMap["updated_at"] = b.UpdatedAt
	b.fieldMap["deleted_at"] = b.DeletedAt
	b.fieldMap["user_id"] = b.UserID
	b.fieldMap["user_login"] = b.UserLogin
	b.fieldMap["channel_id"] = b.ChannelID
}

func (b blacklistEntry) clone(db *gorm.DB) blacklistEntry {
	b.blacklistEntryDo.ReplaceConnPool(db.Statement.ConnPool)
	return b
}

func (b blacklistEntry) replaceDB(db *gorm.DB) blacklistEntry {
	b.blacklistEntryDo.ReplaceDB(db)
	return b
}

type blacklistEntryDo struct{ gen.DO }

type IBlacklistEntryDo interface {
	gen.SubQuery
	Debug() IBlacklistEntryDo
	WithContext(ctx context.Context) IBlacklistEntryDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IBlacklistEntryDo
	WriteDB() IBlacklistEntryDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IBlacklistEntryDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IBlacklistEntryDo
	Not(conds ...gen.Condition) IBlacklistEntryDo
	Or(conds ...gen.Condition) IBlacklistEntryDo
	Select(conds ...field.Expr) IBlacklistEntryDo
	Where(conds ...gen.Condition) IBlacklistEntryDo
	Order(conds ...field.Expr) IBlacklistEntryDo
	Distinct(cols ...field.Expr) IBlacklistEntryDo
	Omit(cols ...field.Expr) IBlacklistEntryDo
	Join(table schema.Tabler, on ...field.Expr) IBlacklistEntryDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IBlacklistEntryDo
	RightJoin(table schema.Tabler, on ...field.Expr) IBlacklistEntryDo
	Group(cols ...field.Expr) IBlacklistEntryDo
	Having(conds ...gen.Condition) IBlacklistEntryDo
	Limit(limit int) IBlacklistEntryDo
	Offset(offset int) IBlacklistEntryDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IBlacklistEntryDo
	Unscoped() IBlacklistEntryDo
	Create(values ...*models.BlacklistEntry) error
	CreateInBatches(values []*models.BlacklistEntry, batchSize int) error
	Save(values ...*models.BlacklistEntry) error
	First() (*models.BlacklistEntry, error)
	Take() (*models.BlacklistEntry, error)
	Last() (*models.BlacklistEntry, error)
	Find() ([]*models.BlacklistEntry, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.BlacklistEntry, err error)
	FindInBatches(result *[]*models.BlacklistEntry, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.BlacklistEntry) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IBlacklistEntryDo
	Assign(attrs ...field.AssignExpr) IBlacklistEntryDo
	Joins(fields ...field.RelationField) IBlacklistEntryDo
	Preload(fields ...field.RelationField) IBlacklistEntryDo
	FirstOrInit() (*models.BlacklistEntry, error)
	FirstOrCreate() (*models.BlacklistEntry, error)
	FindByPage(offset int, limit int) (result []*models.BlacklistEntry, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IBlacklistEntryDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (b blacklistEntryDo) Debug() IBlacklistEntryDo {
	return b.withDO(b.DO.Debug())
}

func (b blacklistEntryDo) WithContext(ctx context.Context) IBlacklistEntryDo {
	return b.withDO(b.DO.WithContext(ctx))
}

func (b blacklistEntryDo) ReadDB() IBlacklistEntryDo {
	return b.Clauses(dbresolver.Read)
}

func (b blacklistEntryDo) WriteDB() IBlacklistEntryDo {
	return b.Clauses(dbresolver.Write)
}

func (b blacklistEntryDo) Session(config *gorm.Session) IBlacklistEntryDo {
	return b.withDO(b.DO.Session(config))
}

func (b blacklistEntryDo) Clauses(conds ...clause.Expression) IBlacklistEntryDo {
	return b.withDO(b.DO.Clauses(conds...))
}

func (b blacklistEntryDo) Returning(value interface{}, columns ...string) IBlacklistEntryDo {
	return b.withDO(b.DO.Returning(value, columns...))
}

func (b blacklistEntryDo) Not(conds ...gen.Condition) IBlacklistEntryDo {
	return b.withDO(b.DO.Not(conds...))
}

func (b blacklistEntryDo) Or(conds ...gen.Condition) IBlacklistEntryDo {
	return b.withDO(b.DO.Or(conds...))
}

func (b blacklistEntryDo) Select(conds ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.Select(conds...))
}

func (b blacklistEntryDo) Where(conds ...gen.Condition) IBlacklistEntryDo {
	return b.withDO(b.DO.Where(conds...))
}

func (b blacklistEntryDo) Order(conds ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.Order(conds...))
}

func (b blacklistEntryDo) Distinct(cols ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.Distinct(cols...))
}

func (b blacklistEntryDo) Omit(cols ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.Omit(cols...))
}

func (b blacklistEntryDo) Join(table schema.Tabler, on ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.Join(table, on...))
}

func (b blacklistEntryDo) LeftJoin(table schema.Tabler, on ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.LeftJoin(table, on...))
}

func (b blacklistEntryDo) RightJoin(table schema.Tabler, on ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.RightJoin(table, on...))
}

func (b blacklistEntryDo) Group(cols ...field.Expr) IBlacklistEntryDo {
	return b.withDO(b.DO.Group(cols...))
}

func (b blacklistEntryDo) Having(conds ...gen.Condition) IBlacklistEntryDo {
	return b.withDO(b.DO.Having(conds...))
}

func (b blacklistEntryDo) Limit(limit int) IBlacklistEntryDo {
	return b.withDO(b.DO.Limit(limit))
}

func (b blacklistEntryDo) Offset(offset int) IBlacklistEntryDo {
	return b.withDO(b.DO.Offset(offset))
}

func (b blacklistEntryDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IBlacklistEntryDo {
	return b.withDO(b.DO.Scopes(funcs...))
}

func (b blacklistEntryDo) Unscoped() IBlacklistEntryDo {
	return b.withDO(b.DO.Unscoped())
}

func (b blacklistEntryDo) Create(values ...*models.BlacklistEntry) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Create(values)
}

func (b blacklistEntryDo) CreateInBatches(values []*models.BlacklistEntry, batchSize int) error {
	return b.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (b blacklistEntryDo) Save(values ...*models.BlacklistEntry) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Save(values)
}

func (b blacklistEntryDo) First() (*models.BlacklistEntry, error) {
	if result, err := b.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.BlacklistEntry), nil
	}
}

func (b blacklistEntryDo) Take() (*models.BlacklistEntry, error) {
	if result, err := b.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.BlacklistEntry), nil
	}
}

func (b blacklistEntryDo) Last() (*models.BlacklistEntry, error) {
	if result, err := b.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.BlacklistEntry), nil
	}
}

func (b blacklistEntryDo) Find() ([]*models.BlacklistEntry, error) {
	result, err := b.DO.Find()
	return result.([]*models.BlacklistEntry), err
}

func (b blacklistEntryDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.BlacklistEntry, err error) {
	buf := make([]*models.BlacklistEntry, 0, batchSize)
	err = b.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (b blacklistEntryDo) FindInBatches(result *[]*models.BlacklistEntry, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return b.DO.FindInBatches(result, batchSize, fc)
}

func (b blacklistEntryDo) Attrs(attrs ...field.AssignExpr) IBlacklistEntryDo {
	return b.withDO(b.DO.Attrs(attrs...))
}

func (b blacklistEntryDo) Assign(attrs ...field.AssignExpr) IBlacklistEntryDo {
	return b.withDO(b.DO.Assign(attrs...))
}

func (b blacklistEntryDo) Joins(fields ...field.RelationField) IBlacklistEntryDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Joins(_f))
	}
	return &b
}

func (b blacklistEntryDo) Preload(fields ...field.RelationField) IBlacklistEntryDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Preload(_f))
	}
	return &b
}

func (b blacklistEntryDo) FirstOrInit() (*models.BlacklistEntry, error) {
	if result, err := b.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.BlacklistEntry), nil
	}
}

func (b blacklistEntryDo) FirstOrCreate() (*models.BlacklistEntry, error) {
	if result, err := b.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.BlacklistEntry), nil
	}
}

func (b blacklistEntryDo) FindByPage(offset int, limit int) (result []*models.BlacklistEntry, count int64, err error) {
	result, err = b.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = b.Offset(-1).Limit(-1).Count()
	return
}

func (b blacklistEntryDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = b.Count()
	if err != nil {
		return
	}

	err = b.Offset(offset).Limit(limit).Scan(result)
	return
}

func (b blacklistEntryDo) Scan(result interface{}) (err error) {
	return b.DO.Scan(result)
}

func (b blacklistEntryDo) Delete(models ...*models.BlacklistEntry) (result gen.ResultInfo, err error) {
	return b.DO.Delete(models)
}

func (b *blacklistEntryDo) withDO(do gen.Dao) *blacklistEntryDo {
	b.DO = *do.(*gen.DO)
	return b
}
