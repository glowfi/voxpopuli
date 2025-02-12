package voxsphere

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Voxspheres() ([]Voxsphere, error)
	VoxsphereByID(context.Context, uuid.UUID) (Voxsphere, error)
	AddVoxsphere(context.Context, Voxsphere) error
	UpdateVoxsphere(context.Context, Voxsphere) error
	DeleteVoxsphere(context.Context, uuid.UUID) error
}

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) Voxspheres(ctx context.Context) ([]Voxsphere, error) { panic("implement me") }

func (r *Repo) VoxsphereByID(ctx context.Context, ID uuid.UUID) (Voxsphere, error) {
	panic("implement me")
}

func (r *Repo) AddVoxsphere(ctx context.Context, voxsphere Voxsphere) error { panic("implement me") }

func (r *Repo) UpdateVoxsphere(ctx context.Context, voxsphere Voxsphere) error { panic("implement me") }

func (r *Repo) DeleteVoxsphere(ctx context.Context, ID uuid.UUID) error { panic("implement me") }
