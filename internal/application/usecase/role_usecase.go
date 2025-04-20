package usecase

import (
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type RoleUseCase struct {
	roleRepo entity.RoleRepository
}

func NewRoleUseCase(repo entity.RoleRepository) *RoleUseCase {
	return &RoleUseCase{roleRepo: repo}
}

type CreateRoleInput struct {
	Name        string
	Permissions []entity.Permission
}

type UpdateRoleInput struct {
	ID          uint
	Name        string
	Permissions []entity.Permission
}

func (uc *RoleUseCase) CreateRole(input *CreateRoleInput) (*entity.Role, error) {
	role := &entity.Role{
		Name:        input.Name,
		Permissions: input.Permissions,
	}

	if err := uc.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (uc *RoleUseCase) UpdateRole(input *UpdateRoleInput) (*entity.Role, error) {
	role, err := uc.roleRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	role.Name = input.Name
	role.Permissions = input.Permissions

	if err := uc.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (uc *RoleUseCase) GetRoleByID(id uint) (*entity.Role, error) {
	return uc.roleRepo.FindByID(id)
}

func (uc *RoleUseCase) GetRoleByName(name string) (*entity.Role, error) {
	return uc.roleRepo.FindByName(name)
}

func (uc *RoleUseCase) ListRoles() ([]entity.Role, error) {
	return uc.roleRepo.List()
}

func (uc *RoleUseCase) DeleteRole(id uint) error {
	return uc.roleRepo.Delete(id)
}

func (uc *RoleUseCase) ValidatePermissions(permissions []entity.Permission) bool {
	validPerms := map[entity.Permission]bool{
		entity.UserCreate:      true,
		entity.UserRead:        true,
		entity.UserUpdate:      true,
		entity.UserDelete:      true,
		entity.RoleCreate:      true,
		entity.RoleRead:        true,
		entity.RoleUpdate:      true,
		entity.RoleDelete:      true,
		entity.AuditLogRead:    true,
		entity.ModuleIntegrate: true,
	}

	for _, p := range permissions {
		if !validPerms[p] {
			return false
		}
	}

	return true
}
