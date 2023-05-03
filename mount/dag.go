package mount

// Mount dag to fileStore
// TODO : fs 의 cid 와 dag 의 cid 가 다를 경우 동기화 처리 필요
// dag 의 block 은 어느 기준으로 Garbage collect?? 블록 전체를 캐싱하고 있으면 안되는데...
// filemanager 도 여기에 넣는게 맞나?
// type Mount struct {
// 	Dag *uih.DagBuilderParams // use MapDataStore
// 	Fs  *file.FileStore       // FileStore
// }

// func NewMount(ctx context.Context, fs *file.FileStore, bs blockstore.Blockstore, rem exchange.Interface) (*Mount, error) {
// 	// make dag service, save dht blocks
// 	// Create a UnixFS graph from our file, parameters described here but can be visualized at https://dag.ipfs.tech/
// 	builder := &uih.DagBuilderParams{
// 		Maxlinks:  uih.DefaultLinksPerBlock, // Default max of 174 links per block
// 		RawLeaves: true,                     // Leave the actual file bytes untouched instead of wrapping them in a dag-pb protobuf wrapper
// 		CidBuilder: cid.V1Builder{ // Use CIDv1 for all links
// 			Codec:    uint64(multicodec.DagPb),
// 			MhType:   uint64(multicodec.Sha3_256), // Use SHA3-256 as the hash function
// 			MhLength: -1,                          // Use the default hash length for the given hash function (in this case 256 bits)
// 		},
// 		Dagserv: merkledag.NewDAGService(blockservice.New(bs, rem)),
// 		NoCopy:  false,
// 	}

// 	mount := &Mount{
// 		Dag: builder,
// 		Fs:  fs,
// 	}

// 	// import blocks in Merkle-DAG from fileStore
// 	err := fs.Iterate("", func(path string, reader *file.Reader) error {
// 		_, err := mount.Upload(ctx, path)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return mount, nil
// }

// func (m *Mount) Download(ctx context.Context, ci cid.Cid, path string) error {
// 	node, err := m.Dag.Dagserv.Get(ctx, ci)
// 	if err != nil {
// 		return err
// 	}

// 	unixFSNode, err := unixfile.NewUnixfsFile(ctx, m.Dag.Dagserv, node)
// 	if err != nil {
// 		return err
// 	}
// 	defer unixFSNode.Close()

// 	err = m.Fs.Overwrite(path, file.NewWriter(unixFSNode))
// 	if err != nil {
// 		return err
// 	}

// 	// put cid
// 	return m.Fs.PutCid(path, ci)
// }

// func (m *Mount) Upload(ctx context.Context, path string) (cid.Cid, error) {
// 	data, err := m.Fs.Get(path)
// 	if err != nil {
// 		return cid.Cid{}, err
// 	}
// 	defer data.Close()

// 	// Split the file up into fixed sized 256KiB chunks
// 	ufsBuilder, err := m.Dag.New(chunker.NewSizeSplitter(data.ReaderFile, chunker.DefaultBlockSize))
// 	if err != nil {
// 		return cid.Cid{}, err
// 	}
// 	nd, err := balanced.Layout(ufsBuilder) // Arrange the graph with a balanced layout
// 	if err != nil {
// 		return cid.Cid{}, err
// 	}

// 	// put cid
// 	ci := nd.Cid()
// 	err = m.Fs.PutCid(path, ci)
// 	if err != nil {
// 		return ci, err
// 	}
// 	return ci, nil
// }