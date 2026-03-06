package repo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseRepo[T, C any] interface {
	Delete(ctx context.Context, ids []int64) error
	UpdateByID(ctx context.Context, entity *T) error
	Create(ctx context.Context, entity *T) error
	Find(ctx context.Context, condition *C) (*T, error)
	FindAll(ctx context.Context, condition *C) ([]*T, error)
	ListPage(ctx context.Context, page, size int64, condition *C) ([]*T, int64, error)
	Count(ctx context.Context, condition *C) (int64, error)
	BatchCreate(ctx context.Context, list []*T) error
}

type baseRepo[T, C any] struct {
	queryDAO  func(ctx context.Context) gen.Dao
	idField   func(ctx context.Context) field.Int64
	id        func(entity *T) int64
	model     *T
	batchSize int
}

func NewBaseRepo[T, C any](
	queryDAO func(ctx context.Context) gen.Dao,
	idField func(ctx context.Context) field.Int64,
	id func(entity *T) int64,
	model *T,
	batchSize int,
) BaseRepo[T, C] {
	return baseRepo[T, C]{
		queryDAO:  queryDAO,
		idField:   idField,
		id:        id,
		model:     model,
		batchSize: batchSize,
	}
}

func (b baseRepo[T, C]) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	res, err := b.queryDAO(ctx).Where(b.idField(ctx).In(ids...)).Delete()
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return nil
	}
	return res.Error
}

func (b baseRepo[T, C]) UpdateByID(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity is nil")
	}
	id := b.id(entity)
	if id == 0 {
		return errors.New("entity id is required")
	}
	res, err := b.queryDAO(ctx).Where(b.idField(ctx).Eq(id)).Updates(entity)
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return nil
	}
	return res.Error
}

func (b baseRepo[T, C]) Create(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity is nil")
	}
	return b.queryDAO(ctx).Create(entity)
}

func (b baseRepo[T, C]) Find(ctx context.Context, condition *C) (*T, error) {
	dao, err := BuildDao(b.queryDAO(ctx), b.model, condition)
	if err != nil {
		return nil, err
	}
	result, err := dao.First()
	if err != nil {
		return nil, err
	}
	item, ok := result.(*T)
	if !ok {
		return nil, fmt.Errorf("unexpected first type %T", result)
	}
	return item, nil
}

func (b baseRepo[T, C]) FindAll(ctx context.Context, condition *C) ([]*T, error) {
	dao, err := BuildDao(b.queryDAO(ctx), b.model, condition)
	if err != nil {
		return nil, err
	}
	var result interface{}
	result, err = dao.Find()
	if err != nil {
		return nil, err
	}
	list, ok := result.([]*T)
	if !ok {
		return nil, fmt.Errorf("unexpected find type %T", result)
	}
	return list, nil
}

func (b baseRepo[T, C]) ListPage(ctx context.Context, page, size int64, condition *C) ([]*T, int64, error) {
	dao, err := BuildDao(b.queryDAO(ctx), b.model, condition)
	if err != nil {
		return nil, 0, err
	}
	offset, limit := PageOffsetLimit(page, size)
	result, err := dao.Offset(int(offset)).Limit(int(limit)).Find()
	if err != nil {
		return nil, 0, err
	}
	list, ok := result.([]*T)
	if !ok {
		return nil, 0, fmt.Errorf("unexpected find type %T", result)
	}
	count, err := dao.Offset(-1).Limit(-1).Count()
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (b baseRepo[T, C]) Count(ctx context.Context, condition *C) (int64, error) {
	dao, err := BuildDao(b.queryDAO(ctx), b.model, condition)
	if err != nil {
		return 0, err
	}
	return dao.Count()
}

func (b baseRepo[T, C]) BatchCreate(ctx context.Context, list []*T) error {
	if len(list) == 0 {
		return nil
	}
	batchSize := b.batchSize
	if batchSize <= 0 {
		batchSize = 100
	}
	return b.queryDAO(ctx).CreateInBatches(list, batchSize)
}

func PageOffsetLimit(page, size int64) (offset, limit int64) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset = (page - 1) * size
	limit = size
	return
}

func BuildDao(dao gen.Dao, model, condition any) (gen.Dao, error) {
	if model == nil {
		if condition == nil {
			return dao, nil
		}
		return nil, errors.New("model is nil")
	}

	modelVal := reflect.Indirect(reflect.ValueOf(model))
	if !modelVal.IsValid() || modelVal.Kind() != reflect.Struct {
		if condition == nil {
			return dao, nil
		}
		return nil, errors.New("model must be struct")
	}

	if condition != nil {
		condVal := reflect.Indirect(reflect.ValueOf(condition))
		if !condVal.IsValid() || condVal.Kind() != reflect.Struct {
			dao = applyUpdatedAtDescIfExists(dao, modelVal)
			return dao, nil
		}

		condType := condVal.Type()
		for i := 0; i < condVal.NumField(); i++ {
			fieldMeta := condType.Field(i)
			tag := parseQueryTag(fieldMeta.Tag.Get("query"))
			op := normalizeOperator(tag["type"])
			column := strings.TrimSpace(tag["column"])
			if op == "" || column == "" {
				continue
			}

			condFieldVal := condVal.Field(i)
			if !isFilterValue(condFieldVal) {
				continue
			}

			modelFieldVal := resolveModelField(modelVal, column)
			if !modelFieldVal.IsValid() {
				return nil, fmt.Errorf("column %s not found on model", column)
			}

			switch op {
			case "order":
				next, ok := applyOrder(dao, modelFieldVal, condFieldVal)
				if !ok {
					return nil, fmt.Errorf("invalid order value for %s", column)
				}
				dao = next
			default:
				expr, ok := buildConditionExpr(modelFieldVal, op, condFieldVal)
				if !ok {
					return nil, fmt.Errorf("invalid condition op=%s column=%s", op, column)
				}
				dao = applyWhere(dao, expr)
			}
		}
	}

	dao = applyUpdatedAtDescIfExists(dao, modelVal)
	return dao, nil
}

func applyUpdatedAtDescIfExists(dao gen.Dao, modelVal reflect.Value) gen.Dao {
	sortField := resolveModelField(modelVal, "sort")
	if sortField.IsValid() {
		expr, ok := callNoArgMethod(sortField, "Asc")
		if ok {
			return applyOrderExpr(dao, expr)
		}
	}

	updatedAtField := resolveModelField(modelVal, "updated_at")
	if !updatedAtField.IsValid() {
		return dao
	}
	expr, ok := callNoArgMethod(updatedAtField, "Desc")
	if !ok {
		return dao
	}
	return applyOrderExpr(dao, expr)
}

func buildLikeValue(key string) string {
	return fmt.Sprintf("%%%s%%", key)
}

func parseQueryTag(raw string) map[string]string {
	result := make(map[string]string)
	if raw == "" {
		return result
	}
	for _, part := range strings.Split(raw, ";") {
		kv := strings.SplitN(strings.TrimSpace(part), ":", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.ToLower(strings.TrimSpace(kv[0]))
		v := strings.TrimSpace(kv[1])
		if k != "" && v != "" {
			result[k] = v
		}
	}
	return result
}

func normalizeOperator(op string) string {
	switch strings.ToLower(strings.TrimSpace(op)) {
	case "eq", "equal":
		return "eq"
	case "in":
		return "in"
	case "contains", "like":
		return "contains"
	case "order", "sort":
		return "order"
	default:
		return ""
	}
}

func isFilterValue(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return !v.IsNil()
	case reflect.Slice, reflect.Array:
		return v.Len() > 0
	default:
		return !v.IsZero()
	}
}

func resolveModelField(modelVal reflect.Value, column string) reflect.Value {
	if !modelVal.IsValid() || column == "" {
		return reflect.Value{}
	}
	fieldName := toModelFieldName(column)
	if v := modelVal.FieldByName(fieldName); v.IsValid() {
		return v
	}
	for i := 0; i < modelVal.NumField(); i++ {
		sf := modelVal.Type().Field(i)
		if strings.EqualFold(sf.Name, fieldName) || strings.EqualFold(sf.Name, column) {
			return modelVal.Field(i)
		}
	}
	return reflect.Value{}
}

func toModelFieldName(column string) string {
	column = strings.TrimSpace(column)
	if column == "" {
		return ""
	}
	parts := strings.FieldsFunc(column, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		runes := []rune(strings.ToLower(parts[i]))
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}

func buildConditionExpr(modelField reflect.Value, op string, condField reflect.Value) (reflect.Value, bool) {
	if condField.Kind() == reflect.Ptr {
		condField = condField.Elem()
	}
	switch op {
	case "eq":
		return callMethod(modelField, "Eq", condField)
	case "contains":
		if condField.Kind() != reflect.String {
			return reflect.Value{}, false
		}
		return callMethod(modelField, "Like", reflect.ValueOf(buildLikeValue(condField.String())))
	case "in":
		if condField.Kind() != reflect.Slice && condField.Kind() != reflect.Array {
			return reflect.Value{}, false
		}
		return callVariadicMethod(modelField, "In", condField)
	default:
		return reflect.Value{}, false
	}
}

func applyOrder(dao gen.Dao, modelField reflect.Value, condField reflect.Value) (gen.Dao, bool) {
	if condField.Kind() == reflect.Ptr {
		condField = condField.Elem()
	}
	direction := "asc"
	switch condField.Kind() {
	case reflect.String:
		direction = strings.ToLower(strings.TrimSpace(condField.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if condField.Int() < 0 {
			direction = "desc"
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if condField.Uint() == 0 {
			return dao, false
		}
	default:
		return dao, false
	}

	var expr reflect.Value
	var ok bool
	if direction == "desc" {
		expr, ok = callNoArgMethod(modelField, "Desc")
	} else {
		expr, ok = callNoArgMethod(modelField, "Asc")
	}
	if !ok {
		return dao, false
	}
	return applyOrderExpr(dao, expr), true
}

func callMethod(target reflect.Value, name string, arg reflect.Value) (reflect.Value, bool) {
	method := target.MethodByName(name)
	if !method.IsValid() || method.Type().NumIn() != 1 {
		return reflect.Value{}, false
	}
	inType := method.Type().In(0)
	if !arg.Type().AssignableTo(inType) {
		if arg.Type().ConvertibleTo(inType) {
			arg = arg.Convert(inType)
		} else {
			return reflect.Value{}, false
		}
	}
	out := method.Call([]reflect.Value{arg})
	if len(out) == 0 {
		return reflect.Value{}, false
	}
	return out[0], true
}

func callNoArgMethod(target reflect.Value, name string) (reflect.Value, bool) {
	method := target.MethodByName(name)
	if !method.IsValid() || method.Type().NumIn() != 0 {
		return reflect.Value{}, false
	}
	out := method.Call(nil)
	if len(out) == 0 {
		return reflect.Value{}, false
	}
	return out[0], true
}

func callVariadicMethod(target reflect.Value, name string, args reflect.Value) (reflect.Value, bool) {
	method := target.MethodByName(name)
	if !method.IsValid() || !method.Type().IsVariadic() || method.Type().NumIn() != 1 {
		return reflect.Value{}, false
	}
	inElemType := method.Type().In(0).Elem()
	callArgs := make([]reflect.Value, 0, args.Len())
	for i := 0; i < args.Len(); i++ {
		arg := args.Index(i)
		if !arg.Type().AssignableTo(inElemType) {
			if arg.Type().ConvertibleTo(inElemType) {
				arg = arg.Convert(inElemType)
			} else {
				return reflect.Value{}, false
			}
		}
		callArgs = append(callArgs, arg)
	}
	out := method.Call(callArgs)
	if len(out) == 0 {
		return reflect.Value{}, false
	}
	return out[0], true
}

func applyWhere(dao gen.Dao, expr reflect.Value) gen.Dao {
	if !expr.IsValid() {
		return dao
	}
	where := reflect.ValueOf(dao).MethodByName("Where")
	if !where.IsValid() || where.Type().NumIn() != 1 {
		return dao
	}
	inType := where.Type().In(0)
	if !expr.Type().AssignableTo(inType) {
		if expr.Type().ConvertibleTo(inType) {
			expr = expr.Convert(inType)
		} else {
			return dao
		}
	}
	out := where.Call([]reflect.Value{expr})
	if len(out) == 1 {
		if next, ok := out[0].Interface().(gen.Dao); ok {
			return next
		}
	}
	return dao
}

func applyOrderExpr(dao gen.Dao, expr reflect.Value) gen.Dao {
	if !expr.IsValid() {
		return dao
	}
	order := reflect.ValueOf(dao).MethodByName("Order")
	if !order.IsValid() || !order.Type().IsVariadic() || order.Type().NumIn() != 1 {
		return dao
	}
	inElemType := order.Type().In(0).Elem()
	if !expr.Type().AssignableTo(inElemType) {
		if expr.Type().ConvertibleTo(inElemType) {
			expr = expr.Convert(inElemType)
		} else {
			return dao
		}
	}
	out := order.Call([]reflect.Value{expr})
	if len(out) == 1 {
		if next, ok := out[0].Interface().(gen.Dao); ok {
			return next
		}
	}
	return dao
}
