package errs

import "golang.org/x/exp/constraints"

func CheckRequiredObject(obj any, name string) error {
	return implCheckRequiredObject(obj, name)
}
func ThrowCheckRequiredObject(obj any, name string) {
	ThrowIf(implCheckRequiredObject(obj, name))
}

func CheckRequiredString(val string, name string) error {
	return implCheckRequiredString(val, name)
}
func ThrowCheckRequiredString(val string, name string) {
	ThrowIf(implCheckRequiredString(val, name))
}

func CheckRequiredStringPtr(val *string, name string) error {
	return implCheckRequiredStringPtr(val, name)
}
func ThrowCheckRequiredStringPtr(val *string, name string) {
	ThrowIf(implCheckRequiredStringPtr(val, name))
}

func CheckRequiredStringPtrIfNotNil(val *string, name string) error {
	return implCheckRequiredStringPtrIfNotNil(val, name)
}
func ThrowCheckRequiredStringPtrIfNotNil(val *string, name string) {
	ThrowIf(implCheckRequiredStringPtrIfNotNil(val, name))
}

func CheckZero[T constraints.Integer | constraints.Float](val T, name string) error {
	return implCheckZero(val, name)
}
func ThrowCheckZero[T constraints.Integer | constraints.Float](val T, name string) {
	ThrowIf(implCheckZero(val, name))
}

func CheckZeroPtr[T constraints.Integer | constraints.Float](val *T, name string) error {
	return implCheckZeroPtr(val, name)
}
func ThrowCheckZeroPtr[T constraints.Integer | constraints.Float](val *T, name string) {
	ThrowIf(implCheckZeroPtr(val, name))
}

func CheckZeroPtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) error {
	return implCheckZeroPtrIfNotNil(val, name)
}
func ThrowCheckZeroPtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) {
	ThrowIf(implCheckZeroPtrIfNotNil(val, name))
}

func CheckNotZero[T constraints.Integer | constraints.Float](val T, name string) error {
	return implCheckNotZero(val, name)
}
func ThrowCheckNotZero[T constraints.Integer | constraints.Float](val T, name string) {
	ThrowIf(implCheckNotZero(val, name))
}

func CheckNotZeroPtr[T constraints.Integer | constraints.Float](val *T, name string) error {
	return implCheckNotZeroPtr(val, name)
}
func ThrowCheckNotZeroPtr[T constraints.Integer | constraints.Float](val *T, name string) {
	ThrowIf(implCheckNotZeroPtr(val, name))
}

func CheckNotZeroPtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) error {
	return implCheckNotZeroPtrIfNotNil(val, name)
}
func ThrowCheckNotZeroPtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) {
	ThrowIf(implCheckNotZeroPtrIfNotNil(val, name))
}

func CheckNotNegative[T constraints.Signed | constraints.Float](val T, name string) error {
	return implCheckNotNegative(val, name)
}
func ThrowCheckNotNegative[T constraints.Signed | constraints.Float](val T, name string) {
	ThrowIf(implCheckNotNegative(val, name))
}

func CheckNotNegativePtr[T constraints.Signed | constraints.Float](val *T, name string) error {
	return implCheckNotNegativePtr(val, name)
}
func ThrowCheckNotNegativePtr[T constraints.Signed | constraints.Float](val *T, name string) {
	ThrowIf(implCheckNotNegativePtr(val, name))
}

func CheckNotNegativePtrIfNotNil[T constraints.Signed | constraints.Float](val *T, name string) error {
	return implCheckNotNegativePtrIfNotNil(val, name)
}
func ThrowCheckNotNegativePtrIfNotNil[T constraints.Signed | constraints.Float](val *T, name string) {
	ThrowIf(implCheckNotNegativePtrIfNotNil(val, name))
}

func CheckPositive[T constraints.Integer | constraints.Float](val T, name string) error {
	return implCheckPositive(val, name)
}
func ThrowCheckPositive[T constraints.Integer | constraints.Float](val T, name string) {
	ThrowIf(implCheckPositive(val, name))
}

func CheckPositivePtr[T constraints.Integer | constraints.Float](val *T, name string) error {
	return implCheckPositivePtr(val, name)
}
func ThrowCheckPositivePtr[T constraints.Integer | constraints.Float](val *T, name string) {
	ThrowIf(implCheckPositivePtr(val, name))
}

func CheckPositivePtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) error {
	return implCheckPositivePtrIfNotNil(val, name)
}
func ThrowCheckPositivePtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) {
	ThrowIf(implCheckPositivePtrIfNotNil(val, name))
}

func CheckValid(obj Validatable, name string) error {
	return implCheckValid(obj, name)
}
func ThrowCheckValid(obj Validatable, name string) {
	ThrowIf(implCheckValid(obj, name))
}

func CheckValidIfNotNil(obj Validatable, name string) error {
	return implCheckValidIfNotNil(obj, name)
}
func ThrowCheckValidIfNotNil(obj Validatable, name string) {
	ThrowIf(implCheckValidIfNotNil(obj, name))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func implCheckRequiredObject(obj any, name string) error {
	return checkRequiredObject(obj, name)
}

func implCheckRequiredString(val string, name string) error {
	return checkRequiredString(val, name)
}

func implCheckRequiredStringPtr(val *string, name string) error {
	err := checkRequiredObject(val, name)
	if err != nil {
		return err
	}
	return checkRequiredString(*val, name)
}

func implCheckRequiredStringPtrIfNotNil(val *string, name string) error {
	if val != nil {
		return checkRequiredString(*val, name)
	}
	return nil
}

func implCheckZero[T constraints.Integer | constraints.Float](val T, name string) error {
	return checkZero(val, name)
}

func implCheckZeroPtr[T constraints.Integer | constraints.Float](val *T, name string) error {
	err := checkRequiredObject(val, name)
	if err != nil {
		return err
	}
	return checkZero(*val, name)
}

func implCheckZeroPtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) error {
	if val != nil {
		return checkZero(*val, name)
	}
	return nil
}

func implCheckNotZero[T constraints.Integer | constraints.Float](val T, name string) error {
	return checkNotZero(val, name)
}

func implCheckNotZeroPtr[T constraints.Integer | constraints.Float](val *T, name string) error {
	err := checkRequiredObject(val, name)
	if err != nil {
		return err
	}
	return checkNotZero(*val, name)
}

func implCheckNotZeroPtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) error {
	if val != nil {
		return checkNotZero(*val, name)
	}
	return nil
}

func implCheckNotNegative[T constraints.Signed | constraints.Float](val T, name string) error {
	return checkNotNegative(val, name)
}

func implCheckNotNegativePtr[T constraints.Signed | constraints.Float](val *T, name string) error {
	err := checkRequiredObject(val, name)
	if err != nil {
		return err
	}
	return checkNotNegative(*val, name)
}

func implCheckNotNegativePtrIfNotNil[T constraints.Signed | constraints.Float](val *T, name string) error {
	if val != nil {
		return checkNotNegative(*val, name)
	}
	return nil
}

func implCheckPositive[T constraints.Integer | constraints.Float](val T, name string) error {
	return checkPositive(val, name)
}

func implCheckPositivePtr[T constraints.Integer | constraints.Float](val *T, name string) error {
	err := checkRequiredObject(val, name)
	if err != nil {
		return err
	}
	return checkPositive(*val, name)
}

func implCheckPositivePtrIfNotNil[T constraints.Integer | constraints.Float](val *T, name string) error {
	if val != nil {
		return checkPositive(*val, name)
	}
	return nil
}

func implCheckValid(obj Validatable, name string) error {
	return checkValid(obj, name)
}

func implCheckValidIfNotNil(obj Validatable, name string) error {
	if !isNil(obj) {
		return checkValid(obj, name)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const checkImplSkip = 3

func checkRequiredObject(obj any, name string) error {
	if isNil(obj) {
		return makeRequiredErrorWithSkip(checkImplSkip, name)
	}
	return nil
}

func checkRequiredString(val string, name string) error {
	if len(val) <= 0 {
		return makeRequiredErrorWithSkip(checkImplSkip, name)
	}
	return nil
}

func checkZero[T constraints.Integer | constraints.Float](val T, name string) error {
	if val != 0 {
		return makeNotZeroErrorWithSkip(checkImplSkip, name)
	}
	return nil
}

func checkNotZero[T constraints.Integer | constraints.Float](val T, name string) error {
	if val == 0 {
		return makeZeroErrorWithSkip(checkImplSkip, name)
	}
	return nil
}

func checkNotNegative[T constraints.Signed | constraints.Float](val T, name string) error {
	if val < 0 {
		return makeNegativeErrorWithSkip(checkImplSkip, name)
	}
	return nil
}

func checkPositive[T constraints.Integer | constraints.Float](val T, name string) error {
	if val <= 0 {
		return makeNotPositiveErrorWithSkip(checkImplSkip, name)
	}
	return nil
}

func checkValid(obj Validatable, name string) error {
	if isNil(obj) {
		return makeRequiredErrorWithSkip(checkImplSkip, name)
	}
	err := obj.Validate()
	if err != nil {
		return makeNotValidErrorWithSkip(checkImplSkip, err, name)
	}
	return nil
}
