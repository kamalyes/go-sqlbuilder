# 业务实战：软删除实现

复制即用的软删除代码。

## 基础软删除

```go
// SoftDeleteUser 软删除用户
func (s *UserService) SoftDeleteUser(ctx context.Context, userID uint) error {
    return repository.SoftDeleteWithDeletedAt[User](ctx, s.repo.GetDB(), userID)
}

// RestoreUser 恢复用户
func (s *UserService) RestoreUser(ctx context.Context, userID uint) error {
    return repository.RestoreDeleted[User](ctx, s.repo.GetDB(), userID)
}

// PermanentlyDeleteUser 彻底删除
func (s *UserService) PermanentlyDeleteUser(ctx context.Context, userID uint) error {
    return repository.PermanentlyDelete[User](ctx, s.repo.GetDB(), userID)
}
```

## 批量操作

```go
// SoftDeleteUsers 批量软删除
func (s *UserService) SoftDeleteUsers(ctx context.Context, userIDs []uint) error {
    for _, id := range userIDs {
        if err := repository.SoftDeleteWithDeletedAt[User](ctx, s.repo.GetDB(), id); err != nil {
            return err
        }
    }
    return nil
}

// RestoreUsers 批量恢复
func (s *UserService) RestoreUsers(ctx context.Context, userIDs []uint) error {
    for _, id := range userIDs {
        if err := repository.RestoreDeleted[User](ctx, s.repo.GetDB(), id); err != nil {
            return err
        }
    }
    return nil
}
```

## 查询控制

```go
// ListActiveUsers 查询未删除用户（默认）
func (s *UserService) ListActiveUsers(ctx context.Context, query *repository.Query) ([]*User, error) {
    return s.repo.List(ctx, query)
}

// ListDeletedUsers 查询已删除用户
func (s *UserService) ListDeletedUsers(ctx context.Context, query *repository.Query) ([]*User, error) {
    return repository.GetDeleted[User](ctx, s.repo.GetDB(), query)
}

// ListAllUsers 查询所有用户（包括已删除）
func (s *UserService) ListAllUsers(ctx context.Context, query *repository.Query) ([]*User, error) {
    return repository.GetWithDeleted[User](ctx, s.repo.GetDB(), query)
}
```

## 回收站服务

```go
// TrashService 回收站服务
type TrashService[T any] struct {
    repo repository.IRepository[T]
}

func NewTrashService[T any](repo repository.IRepository[T]) *TrashService[T] {
    return &TrashService[T]{repo: repo}
}

// ListTrash 获取回收站列表
func (s *TrashService[T]) ListTrash(ctx context.Context) ([]*T, error) {
    return repository.GetDeleted[T](ctx, s.repo.GetDB(), 
        repository.NewQuery().AddOrder("deleted_at", "DESC"))
}

// Restore 恢复
func (s *TrashService[T]) Restore(ctx context.Context, id uint) error {
    return repository.RestoreDeleted[T](ctx, s.repo.GetDB(), id)
}

// Clear 彻底删除
func (s *TrashService[T]) Clear(ctx context.Context, id uint) error {
    return repository.PermanentlyDelete[T](ctx, s.repo.GetDB(), id)
}
```
