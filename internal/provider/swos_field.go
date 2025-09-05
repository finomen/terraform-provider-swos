package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type syncedField[M any, B any] interface {
	Sync(backend *B, model *M)
	Read(backend *B, model *M)
	Name() string
	Attribute() schema.Attribute
}

type syncedFieldImpl[T any, B any, M any, V attr.Value] struct {
	backendGet func(backend *B) *T
	modelGet   func(model *M) *V

	fromModel func(mv V) T
	toModel   func(v T) V

	name      string
	attribute schema.Attribute
}

func (s *syncedFieldImpl[T, B, M, V]) Sync(backend *B, model *M) {
	if s.backendGet == nil || s.modelGet == nil {
		return
	}
	mv := s.modelGet(model)
	if (*mv).IsUnknown() {
		*mv = s.toModel(*s.backendGet(backend))
	} else {
		*s.backendGet(backend) = s.fromModel(*mv)
	}
}

func (s *syncedFieldImpl[T, B, M, V]) Read(backend *B, model *M) {
	if s.backendGet == nil || s.modelGet == nil {
		return
	}
	mv := s.modelGet(model)
	*mv = s.toModel(*s.backendGet(backend))
}

func (s *syncedFieldImpl[T, B, M, V]) Name() string {
	return s.name
}

func (s *syncedFieldImpl[T, B, M, V]) Attribute() schema.Attribute {
	return s.attribute
}
