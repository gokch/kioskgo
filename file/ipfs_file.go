package file

import (
	"context"
	"io"
	"os"

	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/repo/fsrepo"
)

type FsRepo struct {
	Repo *fsrepo.FSRepo
}

func NewFsRepo(rootPath string) (*FsRepo, error) {
	repoPath, err := os.MkdirTemp("", rootPath)
	if err != nil {
		return nil, err
	}

	cfg, err := config.Init(io.Discard, 1024*4)
	if err != nil {
		return nil, err
	}

	cfg.Experimental.FilestoreEnabled = true
	cfg.Experimental.UrlstoreEnabled = true
	cfg.Experimental.Libp2pStreamMounting = true
	cfg.Experimental.P2pHttpProxy = true

	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return nil, err
	}
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	return &FsRepo{Repo: repo.(*fsrepo.FSRepo)}, nil
}

func (f *FsRepo) Put(ctx context.Context) {
	f.Repo.FileManager().Put(ctx, nil)
}
