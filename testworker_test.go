package sectorstorage

import (
	"context"

	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-storage/storage"

	"github.com/filecoin-project/sector-storage/mock"
	"github.com/filecoin-project/sector-storage/sealtasks"
	"github.com/filecoin-project/sector-storage/stores"
	"github.com/filecoin-project/sector-storage/storiface"
)

type testWorker struct {
	acceptTasks map[sealtasks.TaskType]struct{}
	lstor       *stores.Local

	mockSeal *mock.SectorMgr
}

func newTestWorker(wcfg WorkerConfig, lstor *stores.Local) *testWorker {
	ssize, err := wcfg.SealProof.SectorSize()
	if err != nil {
		panic(err)
	}

	acceptTasks := map[sealtasks.TaskType]struct{}{}
	for _, taskType := range wcfg.TaskTypes {
		acceptTasks[taskType] = struct{}{}
	}

	return &testWorker{
		acceptTasks: acceptTasks,
		lstor:       lstor,

		mockSeal: mock.NewMockSectorMgr(ssize),
	}
}

func (t *testWorker) SealPreCommit1(ctx context.Context, sector abi.SectorID, ticket abi.SealRandomness, pieces []abi.PieceInfo) (storage.PreCommit1Out, error) {
	return t.mockSeal.SealPreCommit1(ctx, sector, ticket, pieces)
}

func (t *testWorker) NewSector(ctx context.Context, sector abi.SectorID) error {
	panic("implement me")
}

func (t *testWorker) AddPiece(ctx context.Context, sector abi.SectorID, pieceSizes []abi.UnpaddedPieceSize, newPieceSize abi.UnpaddedPieceSize, pieceData storage.Data) (abi.PieceInfo, error) {
	return t.mockSeal.AddPiece(ctx, sector, pieceSizes, newPieceSize, pieceData)
}

func (t *testWorker) SealPreCommit2(ctx context.Context, sector abi.SectorID, pc1o storage.PreCommit1Out) (storage.SectorCids, error) {
	panic("implement me")
}

func (t *testWorker) SealCommit1(ctx context.Context, sector abi.SectorID, ticket abi.SealRandomness, seed abi.InteractiveSealRandomness, pieces []abi.PieceInfo, cids storage.SectorCids) (storage.Commit1Out, error) {
	panic("implement me")
}

func (t *testWorker) SealCommit2(ctx context.Context, sector abi.SectorID, c1o storage.Commit1Out) (storage.Proof, error) {
	panic("implement me")
}

func (t *testWorker) FinalizeSector(ctx context.Context, sector abi.SectorID) error {
	panic("implement me")
}

func (t *testWorker) Fetch(ctx context.Context, id abi.SectorID, fileType stores.SectorFileType, b bool) error {
	return nil
}

func (t *testWorker) TaskTypes(ctx context.Context) (map[sealtasks.TaskType]struct{}, error) {
	return t.acceptTasks, nil
}

func (t *testWorker) Paths(ctx context.Context) ([]stores.StoragePath, error) {
	return t.lstor.Local(ctx)
}

func (t *testWorker) Info(ctx context.Context) (storiface.WorkerInfo, error) {
	res := ResourceTable[sealtasks.TTPreCommit2][abi.RegisteredProof_StackedDRG2KiBSeal]

	return storiface.WorkerInfo{
		Hostname: "testworkerer",
		Resources: storiface.WorkerResources{
			MemPhysical: res.MinMemory * 3,
			MemSwap:     0,
			MemReserved: res.MinMemory,
			CPUs:        32,
			GPUs:        nil,
		},
	}, nil
}

func (t *testWorker) Closing(ctx context.Context) (<-chan struct{}, error) {
	return ctx.Done(), nil
}

func (t *testWorker) Close() error {
	panic("implement me")
}

var _ Worker = &testWorker{}
