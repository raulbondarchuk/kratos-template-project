package generic

import (
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ========================================
// =============== To Model ===============
// ========================================

// FromDomainGeneric copies fields from the domain entity to the model.
// T — type of the domain entity, M — type of the ORM model.
// The function returns the model and an error if the copy failed.
func ToModelGeneric[T any, M any](domain T) (M, error) {
	var model M
	err := copier.Copy(&model, &domain)
	return model, err
}

// ToModelSliceGeneric converts a slice of domain entities to a slice of models.
// T — type of the domain entity, M — type of the ORM model.
// The function returns the slice of models and an error if the copy failed.
func ToModelSliceGeneric[T any, M any](domains []T) ([]M, error) {
	var result []M
	for _, d := range domains {
		model, err := ToModelGeneric[T, M](d)
		if err != nil {
			return nil, err
		}
		result = append(result, model)
	}
	return result, nil
}

// ========================================
// =============== To Domain ==============
// ========================================

// ToDomainGeneric copies fields from the ORM model to the domain entity.
// M — type of the ORM model, T — type of the domain entity.
// The function returns the domain entity and an error if the copy failed.
func ToDomainGeneric[M any, T any](model M) (T, error) {
	var domain T
	err := copier.Copy(&domain, &model)
	return domain, err
}

// ToDomainSliceGeneric converts a slice of models to a slice of domain entities.
// M — type of the ORM model, T — type of the domain entity.
// The function returns the slice of domain entities and an error if the copy failed.
func ToDomainSliceGeneric[M any, T any](models []M) ([]T, error) {
	var result []T
	for _, m := range models {
		domain, err := ToDomainGeneric[M, T](m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

// ========================================
// =============== To DTO ================
// ========================================

// FromDomainToDTO copies fields from the domain entity to the DTO.
// T — type of the domain entity, D — type of the DTO.
// The function returns the DTO and an error if the copy failed.
func ToDTOGeneric[T any, D any](domain T) (D, error) {
	var dto D
	err := copier.Copy(&dto, &domain)
	return dto, err
}

// ToDTOSliceGeneric converts a slice of domain entities to a slice of DTOs.
// T — type of the domain entity, D — type of the DTO.
// The function returns the slice of DTOs and an error if the copy failed.
func ToDTOSliceGeneric[T any, D any](domains []T) ([]D, error) {
	var result []D
	for _, d := range domains {
		dto, err := ToDTOGeneric[T, D](d)
		if err != nil {
			return nil, err
		}
		result = append(result, dto)
	}
	return result, nil
}

// ========================================
// =============== To Pointer =============
// ========================================

// ToPointerSliceGeneric converts a slice to a slice of pointers.
// T — type of the slice.
// The function returns the slice of pointers.
func ToPointerSliceGeneric[T any](s []T) []*T {
	var result []*T
	for _, t := range s {
		result = append(result, &t)
	}
	return result
}

// ========================================
// =============== To Proto ===============
// ========================================

// ToProtoGeneric copies fields from the domain entity to the proto message.
// D — domain entity, P — proto message.
func ToProtoGeneric[D any, P any](domain D) (P, error) {
	var zero P

	dstT := reflect.TypeOf((*P)(nil)).Elem()
	var dstV reflect.Value
	var dstPtr any
	if dstT.Kind() == reflect.Ptr {
		// P == *SomePB: allocate SomePB and copy into it
		dstV = reflect.New(dstT.Elem()) // *SomePB
		dstPtr = dstV.Interface()
	} else {
		// P == SomePB: create zero and copy into it
		dstV = reflect.New(dstT).Elem()  // SomePB
		dstPtr = dstV.Addr().Interface() // *SomePB
	}

	err := copier.CopyWithOption(dstPtr, &domain, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: &timestamppb.Timestamp{}, // ok for *Timestamp
				Fn: func(src interface{}) (interface{}, error) {
					t := src.(time.Time)
					if t.IsZero() {
						return nil, nil
					}
					return timestamppb.New(t), nil
				},
			},
		},
	})
	if err != nil {
		return zero, err
	}

	if dstT.Kind() == reflect.Ptr {
		return dstV.Interface().(P), nil // Return as *SomePB
	}
	return dstV.Interface().(P), nil // Return as SomePB
}

func ToProtoSliceGeneric[D any, P any](domains []D) ([]P, error) {
	out := make([]P, 0, len(domains))
	for _, d := range domains {
		p, err := ToProtoGeneric[D, P](d)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}
