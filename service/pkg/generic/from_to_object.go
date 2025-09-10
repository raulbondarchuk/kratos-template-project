package generic

import "github.com/jinzhu/copier"

// FromDomainGeneric copies fields from the domain entity to the model.
// T — type of the domain entity, M — type of the ORM model.
// The function returns the model and an error if the copy failed.
func ToModelGeneric[T any, M any](domain T) (M, error) {
	var model M
	err := copier.Copy(&model, &domain)
	return model, err
}

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
