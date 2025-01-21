package spell

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/HeapSoil/auler/internal/auler/store"
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/model"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
)

type SpellBiz interface {
	Create(ctx context.Context, username string, r *v1.CreateSpellRequest) (*v1.CreateSpellResponse, error)
	Update(ctx context.Context, username, spellID string, r *v1.UpdateSpellRequest) error
	Delete(ctx context.Context, username, spellID string) error
	DeleteCollection(ctx context.Context, username string, spellIDs []string) error
	Get(ctx context.Context, username, spellID string) (*v1.GetSpellResponse, error)
	List(ctx context.Context, username string, offset, limit int) (*v1.ListSpellResponse, error)
}

type spellBiz struct {
	ds store.IStore
}

var _ SpellBiz = (*spellBiz)(nil)

func New(ds store.IStore) *spellBiz {
	return &spellBiz{ds: ds}
}

// 创建
func (b *spellBiz) Create(ctx context.Context, username string, r *v1.CreateSpellRequest) (*v1.CreateSpellResponse, error) {
	var spellM model.SpellM
	_ = copier.Copy(&spellM, r)
	spellM.Username = username

	if err := b.ds.Spells().Create(ctx, &spellM); err != nil {
		return nil, err
	}

	return &v1.CreateSpellResponse{SpellID: spellM.SpellID}, nil
}

// 删除一条或若干条咒语
func (b *spellBiz) Delete(ctx context.Context, username, spellID string) error {
	if err := b.ds.Spells().Delete(ctx, username, []string{spellID}); err != nil {
		return err
	}

	return nil
}

func (b *spellBiz) DeleteCollection(ctx context.Context, username string, spellIDs []string) error {
	if err := b.ds.Spells().Delete(ctx, username, spellIDs); err != nil {
		return err
	}

	return nil
}

// 获取咒语
func (b *spellBiz) Get(ctx context.Context, username, spellID string) (*v1.GetSpellResponse, error) {
	spell, err := b.ds.Spells().Get(ctx, username, spellID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSpellNotFound
		}

		return nil, err
	}

	var resp v1.GetSpellResponse
	_ = copier.Copy(&resp, spell)

	resp.CreatedAt = spell.CreatedAt.Format("2006-01-02 15:04:05")
	resp.UpdatedAt = spell.UpdatedAt.Format("2006-01-02 15:04:05")

	return &resp, nil
}

// 更新咒语
func (b *spellBiz) Update(ctx context.Context, username, spellID string, r *v1.UpdateSpellRequest) error {
	spellM, err := b.ds.Spells().Get(ctx, username, spellID)
	if err != nil {
		return err
	}

	if r.Title != nil {
		spellM.Title = *r.Title
	}

	if r.Content != nil {
		spellM.Content = *r.Content
	}

	if err := b.ds.Spells().Update(ctx, spellM); err != nil {
		return err
	}

	return nil
}

// 罗列咒语
func (b *spellBiz) List(ctx context.Context, username string, offset, limit int) (*v1.ListSpellResponse, error) {
	count, list, err := b.ds.Spells().List(ctx, username, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list spells from storage", "err", err)
		return nil, err
	}

	spells := make([]*v1.SpellInfo, 0, len(list))
	for _, item := range list {
		spell := item
		spells = append(spells, &v1.SpellInfo{
			Username:  spell.Username,
			SpellID:   spell.SpellID,
			Title:     spell.Title,
			Content:   spell.Content,
			CreatedAt: spell.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: spell.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.ListSpellResponse{TotalCount: count, Spells: spells}, nil
}
