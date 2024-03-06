package jsonutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"irptools/utils/errs"
)

///////////////////////////////////////////////////////////////////

var (
	ErrUnknownField = errors.New("unknown field")
)

const (
	PredTrue  = "$true"
	PredFalse = "$false"
	PredEq    = "$eq"
	PredIn    = "$in"
	PredNot   = "$not"
	PredAnd   = "$and"
	PredOr    = "$or"
)

///////////////////////////////////////////////////////////////////

type Object interface {
	Field(name string) (any, error)
}

type Logic interface {
	Eq(lhs, rhs any) bool
}

type Predicate interface {
	Is(obj Object) bool
}

///////////////////////////////////////////////////////////////////

func NewMappedObject(fields map[string]func() (any, error)) MappedObject {
	return MappedObject{fields: fields}
}

type MappedObject struct {
	fields map[string]func() (any, error)
}

func (this *MappedObject) Field(name string) (any, error) {
	field, ok := this.fields[name]
	if !ok {
		return nil, ErrUnknownField
	}
	return field()
}

///////////////////////////////////////////////////////////////////

func DefaultLogic() Logic {
	return &defaultLogic{}
}

type defaultLogic struct {
}

func (this *defaultLogic) Eq(lhs, rhs any) bool {
	if lhs == nil || rhs == nil {
		return lhs == rhs
	}
	v1 := reflect.ValueOf(lhs)
	v2 := reflect.ValueOf(rhs)
	if v1.Type() != v2.Type() {
		return fmt.Sprintf("%v", lhs) == fmt.Sprintf("%v", rhs)
	}
	return reflect.DeepEqual(lhs, rhs)
}

///////////////////////////////////////////////////////////////////

func BuildPredicate(rules map[string]any, l Logic) (Predicate, error) {
	preds := make([]Predicate, 0, len(rules))
	for pred, arg := range rules {
		p, _, err := buildPredicate(context{logic: l}, pred, arg)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		preds = append(preds, p)
	}
	return newPredAnd(preds...), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type context struct {
	logic    Logic
	chain    string
	field    string
	hasField bool
}

func (this context) withField() context {
	this.hasField = true
	return this
}

func (this context) withSetField(field string) context {
	if this.field != "" {
		panic("logic error: already with field")
	}
	this.field = field
	return this.withField().joinElem(field)
}

func (this context) joinIndex(idx int) context {
	this.chain += fmt.Sprintf("[%v]", idx)
	return this
}

func (this context) joinElem(elem string) context {
	if elem == "" {
		panic("logic error: elem is empty")
	}
	this.chain += "." + elem
	return this
}

func (this context) String() string {
	return this.chain
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func buildPredicate(ctx context, pred string, rule any) (resPred Predicate, resCtx context, err error) {
	if pred == "" {
		return nil, ctx, errs.Errorf("empty pred %s=%s", ctx, prettyPrintedKeyValue(pred, rule))
	}

	defer func() {
		if err == nil {
			if !resCtx.hasField {
				err = errs.Errorf("without field: %s.%s", resCtx, prettyPrintedKeyValue(pred, rule))
			}
		}
	}()

	if strings.Index(pred, PredTrue) == 0 {
		if !ctx.hasField {
			ctx = ctx.withSetField(PredTrue)
		}
		return newPredConstant(true), ctx, nil
	}

	if strings.Index(pred, PredFalse) == 0 {
		if !ctx.hasField {
			ctx = ctx.withSetField(PredFalse)
		}
		return newPredConstant(false), ctx, nil
	}

	type buildPredFn = func(ctx context, value any) (Predicate, context, error)
	supportedPreds := map[string]buildPredFn{
		PredEq:  buildPredEq,
		PredIn:  buildPredIn,
		PredNot: buildPredNot,
		PredAnd: buildPredAnd,
		PredOr:  buildPredOr,
	}

	build, ok := supportedPreds[pred]
	if ok {
		return build(ctx.joinElem(pred), rule)
	}

	return buildField(ctx, pred, rule)
}

func buildPredEq(ctx context, value any) (Predicate, context, error) {
	if ctx.field == "" {
		// It is not needed by provides error with stack
		return nil, ctx, errs.Error("there is no any specified field to be compared with")
	}
	return newPredEq(ctx.logic, ctx.field, value), ctx, nil
}

func buildPredIn(ctx context, value any) (Predicate, context, error) {
	if ctx.field == "" {
		// It is not needed by provides error with stack
		return nil, ctx, errs.Error("there is no any specified field to be compared with")
	}

	arr, ok := value.([]any)
	if !ok {
		return nil, ctx, errs.Errorf("expected []: unexpected value type %s=%s", ctx, prettyPrintedValue(value))
	}

	return newPredIn(ctx.logic, ctx.field, arr...), ctx, nil
}

func buildPredNot(ctx context, rule any) (Predicate, context, error) {
	m, ok := rule.(map[string]any)
	if !ok {
		return nil, ctx, errs.Errorf("expected {}: unexpected rule type %s=%s", ctx, prettyPrintedValue(rule))
	}
	if len(m) != 1 {
		return nil, ctx, errs.Errorf("expected 1: unexpected rules count=%v %s=%s", len(m), ctx, prettyPrintedValue(rule))
	}

	childPredName, childRule := getKeyValue(m)
	childPred, childCtx, childErr := buildPredicate(ctx, childPredName, childRule)
	if childErr != nil {
		return nil, ctx, errs.Errorf("bad child rule %s=%s: %w", ctx, prettyPrintedKeyValue(childPredName, childRule), childErr)
	}

	return newPredNot(childPred), childCtx, nil
}

func buildPredAnd(ctx context, rule any) (Predicate, context, error) {
	preds, newCtx, err := buildPredsArray(ctx, rule)
	if err != nil {
		return nil, newCtx, errs.Wrap(err)
	}
	return newPredAnd(preds...), newCtx, nil
}

func buildPredOr(ctx context, rule any) (Predicate, context, error) {
	preds, newCtx, err := buildPredsArray(ctx, rule)
	if err != nil {
		return nil, newCtx, errs.Wrap(err)
	}
	return newPredOr(preds...), newCtx, nil
}

func buildPredsArray(ctx context, rule any) ([]Predicate, context, error) {
	rules, ok := rule.([]any)
	if !ok {
		return nil, ctx, errs.Errorf("expected []: unexpected rule type %s=%s", ctx, prettyPrintedValue(rule))
	}

	if len(rules) == 0 {
		return nil, ctx, errs.Errorf("expected not 0: unexpected rules count=0 %s:%s", ctx, prettyPrintedValue(rule))
	}

	newCtx := ctx
	preds := make([]Predicate, 0, len(rules))
	for i, r := range rules {
		ictx := ctx.joinIndex(i)
		m, ok := r.(map[string]any)
		if !ok {
			return nil, ctx, errs.Errorf("expected {}: unexpected rule type %s=%s", ictx, prettyPrintedValue(r))
		}
		if len(m) != 1 {
			return nil, ctx, errs.Errorf("expected 1: unexpected rules count=%v %s=%s", len(m), ictx, prettyPrintedValue(r))
		}

		childPredName, childRule := getKeyValue(m)
		childPred, childCtx, childErr := buildPredicate(ictx, childPredName, childRule)
		if childErr != nil {
			return nil, ctx, errs.Errorf("bad child rule %s=%s: %w", ictx, prettyPrintedKeyValue(childPredName, childRule), childErr)
		}

		if childCtx.hasField {
			newCtx = newCtx.withField()
		}

		preds = append(preds, childPred)
	}

	return preds, newCtx, nil
}

func buildField(ctx context, field string, rule any) (Predicate, context, error) {
	if field == "" {
		panic("logic error: empty field")
	}

	if ctx.field != "" {
		return nil, ctx, errs.Errorf("multiple fields[%s, %s]: %s=%s", ctx.field, field, ctx, prettyPrintedKeyValue(field, rule))
	}

	newCtx := ctx.withSetField(field)
	m, ok := rule.(map[string]any)
	if ok {
		if len(m) != 1 {
			return nil, ctx, fmt.Errorf("expected 1: unexpected rules count=%v %s=%s", len(m), ctx, prettyPrintedKeyValue(field, rule))
		}

		childPred, childRule := getKeyValue(m)
		return buildPredicate(newCtx, childPred, childRule)
	}

	return newPredEq(ctx.logic, field, rule), newCtx, nil
}

///////////////////////////////////////////////////////////////////

func getKeyValue(m map[string]any) (string, any) {
	if len(m) != 1 {
		panic(fmt.Sprintf("logic error: map length is %v", len(m)))
	}
	for k, v := range m {
		return k, v
	}
	return "", nil
}

func prettyPrintedValue(value any) string {
	pretty, e := json.Marshal(value)
	if e != nil {
		return fmt.Sprintf("%v", e)
	}
	return fmt.Sprintf("%v", string(pretty))
}

func prettyPrintedKeyValue(key string, value any) string {
	return prettyPrintedValue(map[string]any{key: value})
}

///////////////////////////////////////////////////////////////////

func newPredConstant(res bool) Predicate {
	return &predConstant{res: res}
}

type predConstant struct {
	res bool
}

func (this *predConstant) Is(obj Object) bool {
	return this.res
}

///////////////////////////////////////////////////////////////////

func newPredAnd(preds ...Predicate) Predicate {
	return &predAnd{preds: preds}
}

type predAnd struct {
	preds []Predicate
}

func (this *predAnd) Is(obj Object) bool {
	for _, p := range this.preds {
		if !p.Is(obj) {
			return false
		}
	}
	return len(this.preds) > 0
}

///////////////////////////////////////////////////////////////////

func newPredOr(preds ...Predicate) Predicate {
	return &predOr{preds: preds}
}

type predOr struct {
	preds []Predicate
}

func (this *predOr) Is(obj Object) bool {
	for _, p := range this.preds {
		if p.Is(obj) {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////////

func newPredNot(pred Predicate) Predicate {
	return &predNot{pred: pred}
}

type predNot struct {
	pred Predicate
}

func (this *predNot) Is(obj Object) bool {
	return !this.pred.Is(obj)
}

///////////////////////////////////////////////////////////////////

func newPredIn(l Logic, field string, values ...any) Predicate {
	return &predIn{
		logic:  l,
		field:  field,
		values: values,
	}
}

type predIn struct {
	logic  Logic
	field  string
	values []any
}

func (this *predIn) Is(obj Object) bool {
	v, err := obj.Field(this.field)
	if err != nil {
		return false
	}

	for _, vv := range this.values {
		if this.logic.Eq(vv, v) {
			return true
		}
	}

	return false
}

///////////////////////////////////////////////////////////////////

func newPredEq(l Logic, field string, value any) Predicate {
	return &predEq{
		logic: l,
		field: field,
		value: value,
	}
}

type predEq struct {
	logic Logic
	field string
	value any
}

func (this *predEq) Is(obj Object) bool {
	v, err := obj.Field(this.field)
	if err != nil {
		return false
	}
	return this.logic.Eq(this.value, v)
}

///////////////////////////////////////////////////////////////////
