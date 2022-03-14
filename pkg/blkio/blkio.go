package blkio

import (
	"context"
	"sync"

	"github.com/rfyiamcool/go-fusehook/pkg/hookfs"
	"golang.org/x/time/rate"
)

// BlkioContext implements hookfs.HookContext
type BlkioContext struct {
	path string
}

// Hook implements hookfs.Hook
type Hook struct {
	sync.Mutex

	ReadLimiter  *rate.Limiter
	WriteLimiter *rate.Limiter
}

func NewHook(reader, writer *rate.Limiter) *Hook {
	hk := &Hook{}
	if reader != nil {
		hk.ReadLimiter = reader
	}
	if writer != nil {
		hk.WriteLimiter = writer
	}
	return hk
}

func (h *Hook) takeReadLimiter(n int) {
	if h.ReadLimiter != nil {
		h.ReadLimiter.WaitN(context.Background(), n)
	}
}

func (h *Hook) takeWriteLimiter(n int) {
	if h.ReadLimiter != nil {
		h.WriteLimiter.WaitN(context.Background(), n)
	}
}

// Init implements hookfs.HookWithInit
func (h *Hook) Init() error {
	return nil
}

// PreOpen implements hookfs.HookOnOpen
func (h *Hook) PreOpen(path string, flags uint32) (bool, hookfs.HookContext, error) {
	ctx := BlkioContext{path: path}
	return false, ctx, nil
}

// PostOpen implements hookfs.HookOnOpen
func (h *Hook) PostOpen(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	return false, nil
}

// PreRead implements hookfs.HookOnRead
func (h *Hook) PreRead(path string, length int64, offset int64) ([]byte, bool, hookfs.HookContext, error) {
	h.takeReadLimiter(int(length))
	ctx := BlkioContext{path: path}
	return nil, false, ctx, nil
}

// PostRead implements hookfs.HookOnRead
func (h *Hook) PostRead(realRetCode int32, realBuf []byte, ctx hookfs.HookContext) ([]byte, bool, error) {
	return nil, false, nil
}

// PreWrite implements hookfs.HookOnWrite
func (h *Hook) PreWrite(path string, buf []byte, offset int64) (bool, hookfs.HookContext, error) {
	h.takeWriteLimiter(len(buf))
	ctx := BlkioContext{path: path}
	return false, ctx, nil
}

// PostWrite implements hookfs.HookOnWrite
func (h *Hook) PostWrite(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	return false, nil
}

// PreMkdir implements hookfs.HookOnMkdir
func (h *Hook) PreMkdir(path string, mode uint32) (bool, hookfs.HookContext, error) {
	ctx := BlkioContext{path: path}
	return false, ctx, nil
}

// PostMkdir implements hookfs.HookOnMkdir
func (h *Hook) PostMkdir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	return false, nil
}

// PreRmdir implements hookfs.HookOnRmdir
func (h *Hook) PreRmdir(path string) (bool, hookfs.HookContext, error) {
	ctx := BlkioContext{path: path}
	return false, ctx, nil
}

// PostRmdir implements hookfs.HookOnRmdir
func (h *Hook) PostRmdir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	return false, nil
}

// PreOpenDir implements hookfs.HookOnOpenDir
func (h *Hook) PreOpenDir(path string) (bool, hookfs.HookContext, error) {
	ctx := BlkioContext{path: path}
	return false, ctx, nil
}

// PostOpenDir implements hookfs.HookOnOpenDir
func (h *Hook) PostOpenDir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	return false, nil
}

// PreFsync implements hookfs.HookOnFsync
func (h *Hook) PreFsync(path string, flags uint32) (bool, hookfs.HookContext, error) {
	ctx := BlkioContext{path: path}
	return false, ctx, nil
}

// PostFsync implements hookfs.HookOnFsync
func (h *Hook) PostFsync(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	return false, nil
}
