package datax

import (
	"context"
	"github.com/nfangxu/mpc-go/internal/utils"
	"github.com/pkg/errors"
	"github.com/smallnest/rpcx/client"
	"sync"
	"time"
)

type Data struct {
	data      sync.Map
	ChunkSize int
}

type ChunkReq struct {
	UUID string
	Size int
}

type ChunkRep struct {
	UUID string
	Size int
	Data []byte
	Done bool
}

func (d *Data) Chunk(ctx context.Context, req *ChunkReq, rep *ChunkRep) error {
	v, ok := d.data.LoadAndDelete(req.UUID)
	if !ok {
		return errors.New("empty")
	}
	data := v.([]byte)
	rep.UUID = req.UUID
	if req.Size == 0 { // 一次性拉取
		rep.Done = true
		rep.Data = data
		rep.Size = len(data)
		return nil
	}

	// 分片
	rep.Data, data = Head(data, req.Size)
	rep.Size = len(rep.Data)
	rep.Done = len(data) == 0

	// 还没传完，剩下的放回去
	if !rep.Done {
		d.data.Store(req.UUID, data)
	}
	return nil
}

type Empty struct{}

func (d *Data) CleanAll(_ context.Context, _ *Empty, _ *Empty) error {
	d.data.Range(func(key, value any) bool {
		d.data.Delete(key)
		return true
	})
	return nil
}

func (d *Data) Push(key string, data []byte) error {
	hashid := utils.MD5(key)
	_, in := d.data.Load(hashid)
	if in {
		return errors.New("duplicate id")
	}
	d.data.Store(hashid, data)
	return nil
}

func (d *Data) Pull(c client.XClient, key string) ([]byte, error) {
	data := make([]byte, 0)

	for {
		req := &ChunkReq{
			UUID: utils.MD5(key),
			Size: d.ChunkSize,
		}
		rep := &ChunkRep{}
		if err := utils.Try(func() error {
			return c.Call(context.Background(), "Chunk", req, rep)
		}, 30, time.Second); err != nil {
			return nil, errors.Wrap(err, "Pull data")
		}
		data = append(data, rep.Data...)
		if rep.Done {
			break
		}
	}

	return data, nil
}
