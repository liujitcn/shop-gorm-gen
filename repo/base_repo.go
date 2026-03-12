package repo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
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

type QueryOperator int

const (
	EQ QueryOperator = 1

	NEQ QueryOperator = 2
	GT  QueryOperator = 3
	GTE QueryOperator = 4
	LT  QueryOperator = 5
	LTE QueryOperator = 6

	LIKE     QueryOperator = 7
	ILIKE    QueryOperator = 8
	NOT_LIKE QueryOperator = 9

	IN  QueryOperator = 10
	NIN QueryOperator = 11

	IS_NULL     QueryOperator = 12
	IS_NOT_NULL QueryOperator = 13

	BETWEEN QueryOperator = 14
	REGEXP  QueryOperator = 15
	IREGEXP QueryOperator = 16

	CONTAINS     QueryOperator = 17
	STARTS_WITH  QueryOperator = 18
	ENDS_WITH    QueryOperator = 19
	ICONTAINS    QueryOperator = 20
	ISTARTS_WITH QueryOperator = 21
	IENDS_WITH   QueryOperator = 22

	JSON_CONTAINS  QueryOperator = 23
	ARRAY_CONTAINS QueryOperator = 24
	EXISTS         QueryOperator = 25

	SEARCH QueryOperator = 26
	EXACT  QueryOperator = 27
	IEXACT QueryOperator = 28

	ORDER QueryOperator = 100
)

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
			op, ok := parseQueryOperator(tag["type"])
			column := strings.TrimSpace(tag["column"])
			if !ok || column == "" {
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
			queryFieldVal := buildQueryFieldValue(dao.TableName(), column, modelFieldVal.Type())

			switch op {
			case ORDER:
				next, ok := applyOrder(dao, queryFieldVal, condFieldVal)
				if !ok {
					return nil, fmt.Errorf("invalid order value for %s", column)
				}
				dao = next
			default:
				expr, ok := buildConditionExpr(queryFieldVal, op, condFieldVal)
				if !ok {
					return nil, fmt.Errorf("invalid condition op=%s column=%s", queryOperatorName(op), column)
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
		expr, ok := callNoArgMethod(buildQueryFieldValue(dao.TableName(), "sort", sortField.Type()), "Asc")
		if ok {
			return applyOrderExpr(dao, expr)
		}
	}

	updatedAtField := resolveModelField(modelVal, "updated_at")
	if !updatedAtField.IsValid() {
		return dao
	}
	expr, ok := callNoArgMethod(buildQueryFieldValue(dao.TableName(), "updated_at", updatedAtField.Type()), "Desc")
	if !ok {
		return dao
	}
	return applyOrderExpr(dao, expr)
}

func buildQueryFieldValue(table, column string, modelFieldType reflect.Type) reflect.Value {
	if modelFieldType.Kind() == reflect.Ptr {
		modelFieldType = modelFieldType.Elem()
	}
	if modelFieldType == reflect.TypeOf(time.Time{}) {
		return reflect.ValueOf(field.NewTime(table, column))
	}

	switch modelFieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(field.NewString(table, column))
	case reflect.Int:
		return reflect.ValueOf(field.NewInt(table, column))
	case reflect.Int8:
		return reflect.ValueOf(field.NewInt8(table, column))
	case reflect.Int16:
		return reflect.ValueOf(field.NewInt16(table, column))
	case reflect.Int32:
		return reflect.ValueOf(field.NewInt32(table, column))
	case reflect.Int64:
		return reflect.ValueOf(field.NewInt64(table, column))
	case reflect.Uint:
		return reflect.ValueOf(field.NewUint(table, column))
	case reflect.Uint8:
		return reflect.ValueOf(field.NewUint8(table, column))
	case reflect.Uint16:
		return reflect.ValueOf(field.NewUint16(table, column))
	case reflect.Uint32:
		return reflect.ValueOf(field.NewUint32(table, column))
	case reflect.Uint64:
		return reflect.ValueOf(field.NewUint64(table, column))
	case reflect.Float32:
		return reflect.ValueOf(field.NewFloat32(table, column))
	case reflect.Float64:
		return reflect.ValueOf(field.NewFloat64(table, column))
	case reflect.Bool:
		return reflect.ValueOf(field.NewBool(table, column))
	case reflect.Slice:
		if modelFieldType.Elem().Kind() == reflect.Uint8 {
			return reflect.ValueOf(field.NewBytes(table, column))
		}
	}
	return reflect.ValueOf(field.NewField(table, column))
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

func parseQueryOperator(raw string) (QueryOperator, bool) {
	op := strings.ToLower(strings.TrimSpace(raw))
	if op == "" {
		return 0, false
	}
	if num, err := strconv.Atoi(op); err == nil {
		switch QueryOperator(num) {
		case EQ, NEQ, GT, GTE, LT, LTE, LIKE, ILIKE, NOT_LIKE, IN, NIN, IS_NULL, IS_NOT_NULL, BETWEEN,
			REGEXP, IREGEXP, CONTAINS, STARTS_WITH, ENDS_WITH, ICONTAINS, ISTARTS_WITH, IENDS_WITH,
			JSON_CONTAINS, ARRAY_CONTAINS, EXISTS, SEARCH, EXACT, IEXACT, ORDER:
			return QueryOperator(num), true
		default:
			return 0, false
		}
	}

	switch op {
	case "eq", "equal", "exact":
		return EQ, true
	case "iexact":
		return IEXACT, true
	case "neq", "ne", "not_eq":
		return NEQ, true
	case "gt":
		return GT, true
	case "gte", "ge":
		return GTE, true
	case "lt":
		return LT, true
	case "lte", "le":
		return LTE, true
	case "like":
		return LIKE, true
	case "ilike":
		return ILIKE, true
	case "not_like", "notlike":
		return NOT_LIKE, true
	case "in":
		return IN, true
	case "nin", "not_in", "notin":
		return NIN, true
	case "is_null", "isnull":
		return IS_NULL, true
	case "is_not_null", "isnotnull":
		return IS_NOT_NULL, true
	case "between":
		return BETWEEN, true
	case "regexp":
		return REGEXP, true
	case "iregexp":
		return IREGEXP, true
	case "contains":
		return CONTAINS, true
	case "starts_with", "startswith":
		return STARTS_WITH, true
	case "ends_with", "endswith":
		return ENDS_WITH, true
	case "icontains":
		return ICONTAINS, true
	case "istarts_with", "istartswith":
		return ISTARTS_WITH, true
	case "iends_with", "iendswith":
		return IENDS_WITH, true
	case "json_contains", "jsoncontains":
		return JSON_CONTAINS, true
	case "array_contains", "arraycontains":
		return ARRAY_CONTAINS, true
	case "exists":
		return EXISTS, true
	case "search":
		return SEARCH, true
	case "order", "sort":
		return ORDER, true
	default:
		return 0, false
	}
}

func queryOperatorName(op QueryOperator) string {
	switch op {
	case EQ:
		return "eq"
	case NEQ:
		return "neq"
	case GT:
		return "gt"
	case GTE:
		return "gte"
	case LT:
		return "lt"
	case LTE:
		return "lte"
	case LIKE:
		return "like"
	case ILIKE:
		return "ilike"
	case NOT_LIKE:
		return "not_like"
	case IN:
		return "in"
	case NIN:
		return "nin"
	case IS_NULL:
		return "is_null"
	case IS_NOT_NULL:
		return "is_not_null"
	case BETWEEN:
		return "between"
	case REGEXP:
		return "regexp"
	case IREGEXP:
		return "iregexp"
	case CONTAINS:
		return "contains"
	case STARTS_WITH:
		return "starts_with"
	case ENDS_WITH:
		return "ends_with"
	case ICONTAINS:
		return "icontains"
	case ISTARTS_WITH:
		return "istarts_with"
	case IENDS_WITH:
		return "iends_with"
	case JSON_CONTAINS:
		return "json_contains"
	case ARRAY_CONTAINS:
		return "array_contains"
	case EXISTS:
		return "exists"
	case SEARCH:
		return "search"
	case EXACT:
		return "exact"
	case IEXACT:
		return "iexact"
	case ORDER:
		return "order"
	default:
		return "unknown"
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

func buildConditionExpr(modelField reflect.Value, op QueryOperator, condField reflect.Value) (reflect.Value, bool) {
	if condField.Kind() == reflect.Ptr {
		condField = condField.Elem()
	}
	switch op {
	case EQ, EXACT:
		return callMethod(modelField, "Eq", condField)
	case NEQ:
		return callMethod(modelField, "Neq", condField)
	case GT:
		return callMethod(modelField, "Gt", condField)
	case GTE:
		return callMethod(modelField, "Gte", condField)
	case LT:
		return callMethod(modelField, "Lt", condField)
	case LTE:
		return callMethod(modelField, "Lte", condField)
	case LIKE:
		return callMethod(modelField, "Like", condField)
	case ILIKE, ICONTAINS, ISTARTS_WITH, IENDS_WITH, IEXACT:
		s, ok := toLowerString(condField)
		if !ok {
			return reflect.Value{}, false
		}
		lowerField, ok := callNoArgMethod(modelField, "Lower")
		if !ok {
			lowerField = modelField
		}
		pattern := s
		switch op {
		case ICONTAINS:
			pattern = buildContainsPattern(s)
		case ISTARTS_WITH:
			pattern = buildStartsWithPattern(s)
		case IENDS_WITH:
			pattern = buildEndsWithPattern(s)
		case ILIKE:
		case IEXACT:
		}
		return callMethod(lowerField, methodNameByOperator(op), reflect.ValueOf(pattern))
	case NOT_LIKE:
		return callMethod(modelField, "NotLike", condField)
	case IN:
		if condField.Kind() != reflect.Slice && condField.Kind() != reflect.Array {
			return reflect.Value{}, false
		}
		return callVariadicMethod(modelField, "In", condField)
	case NIN:
		if condField.Kind() != reflect.Slice && condField.Kind() != reflect.Array {
			return reflect.Value{}, false
		}
		return callVariadicMethod(modelField, "NotIn", condField)
	case IS_NULL:
		return callNoArgMethod(modelField, "IsNull")
	case IS_NOT_NULL:
		return callNoArgMethod(modelField, "IsNotNull")
	case BETWEEN:
		return buildBetweenExpr(modelField, condField)
	case REGEXP:
		return callMethod(modelField, "Regexp", condField)
	case IREGEXP:
		s, ok := toLowerString(condField)
		if !ok {
			return reflect.Value{}, false
		}
		lowerField, ok := callNoArgMethod(modelField, "Lower")
		if !ok {
			lowerField = modelField
		}
		return callMethod(lowerField, "Regexp", reflect.ValueOf(s))
	case CONTAINS, SEARCH, JSON_CONTAINS, ARRAY_CONTAINS:
		s, ok := toString(condField)
		if !ok {
			return reflect.Value{}, false
		}
		return callMethod(modelField, "Like", reflect.ValueOf(buildContainsPattern(s)))
	case STARTS_WITH:
		s, ok := toString(condField)
		if !ok {
			return reflect.Value{}, false
		}
		return callMethod(modelField, "Like", reflect.ValueOf(buildStartsWithPattern(s)))
	case ENDS_WITH:
		s, ok := toString(condField)
		if !ok {
			return reflect.Value{}, false
		}
		return callMethod(modelField, "Like", reflect.ValueOf(buildEndsWithPattern(s)))
	case EXISTS:
		if condField.Kind() == reflect.Bool {
			if condField.Bool() {
				return callNoArgMethod(modelField, "IsNotNull")
			}
			return callNoArgMethod(modelField, "IsNull")
		}
		return callNoArgMethod(modelField, "IsNotNull")
	default:
		return reflect.Value{}, false
	}
}

func methodNameByOperator(op QueryOperator) string {
	switch op {
	case ILIKE, ICONTAINS, ISTARTS_WITH, IENDS_WITH:
		return "Like"
	case IEXACT:
		return "Eq"
	default:
		return ""
	}
}

func buildBetweenExpr(modelField reflect.Value, condField reflect.Value) (reflect.Value, bool) {
	if condField.Kind() != reflect.Slice && condField.Kind() != reflect.Array {
		return reflect.Value{}, false
	}
	if condField.Len() != 2 {
		return reflect.Value{}, false
	}
	left := condField.Index(0)
	right := condField.Index(1)
	return callTwoArgMethod(modelField, "Between", left, right)
}

func toString(v reflect.Value) (string, bool) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", false
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.String {
		return "", false
	}
	return v.String(), true
}

func toLowerString(v reflect.Value) (string, bool) {
	s, ok := toString(v)
	if !ok {
		return "", false
	}
	return strings.ToLower(s), true
}

func buildContainsPattern(key string) string {
	return fmt.Sprintf("%%%s%%", key)
}

func buildStartsWithPattern(key string) string {
	return fmt.Sprintf("%s%%", key)
}

func buildEndsWithPattern(key string) string {
	return fmt.Sprintf("%%%s", key)
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

func callTwoArgMethod(target reflect.Value, name string, arg1 reflect.Value, arg2 reflect.Value) (reflect.Value, bool) {
	method := target.MethodByName(name)
	if !method.IsValid() || method.Type().NumIn() != 2 {
		return reflect.Value{}, false
	}
	inType1 := method.Type().In(0)
	if !arg1.Type().AssignableTo(inType1) {
		if arg1.Type().ConvertibleTo(inType1) {
			arg1 = arg1.Convert(inType1)
		} else {
			return reflect.Value{}, false
		}
	}
	inType2 := method.Type().In(1)
	if !arg2.Type().AssignableTo(inType2) {
		if arg2.Type().ConvertibleTo(inType2) {
			arg2 = arg2.Convert(inType2)
		} else {
			return reflect.Value{}, false
		}
	}
	out := method.Call([]reflect.Value{arg1, arg2})
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
	cond, ok := expr.Interface().(gen.Condition)
	if !ok {
		return dao
	}
	return dao.Where(cond)
}

func applyOrderExpr(dao gen.Dao, expr reflect.Value) gen.Dao {
	if !expr.IsValid() {
		return dao
	}
	orderExpr, ok := expr.Interface().(field.Expr)
	if !ok {
		return dao
	}
	return dao.Order(orderExpr)
}
