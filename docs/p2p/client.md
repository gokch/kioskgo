# p2p

This package provides a peer-to-peer (P2P) client that can connect to other peers in the ipfs network, download and upload files, and manage its own connections.

## ClientConfig

The `ClientConfCig` struct contains the configuration for the client.

| Field        | Description                                                                 |
|--------------|-----------------------------------------------------------------------------|
| `RootPath`   | The path to the root directory of the file system.                          |
| `Peers`      | A list of peer addresses to connect to.                                     |
| `SizeWorker` | The number of worker goroutines to use for downloading and uploading files. |
| `ExpireSec`  | The number of seconds to keep a file in the cache before it is evicted.     |


## Client

The `Client` struct represents a P2P client

| Field   | Description         |
|---------|---------------------|
| `dag`   | The DAG mount.      |
| `mq`    | The ants pool.      |
| `host`  | The libp2p host.    |
| `bswap` | The bitswap client. |

## NewClient

The `NewClient` function creates a new client.


func NewClient(ctx context.Context, cfg *ClientConfig) (*Client, error)

**Parameters:**

* `ctx` - The context of the operation.
* `cfg` - The configuration for the client.

**Returns:**

* A pointer to a new `Client` object, or an error if the client could not be created.

## Self

The `Self` function returns the client's self address.

func (c *Client) Self() string {


**Parameters:**

* `c` - The `Client` object.

**Returns:**

The client's self address.

## Connect

The `Connect` function connects the client to a peer.

func (c *Client) Connect(ctx context.Context, targetPeer string) error


**Parameters:**

* `c` - The `Client` object.
* `targetPeer` - The address of the peer to connect to.

**Returns:**

An error if the client could not connect to the peer.

## Disconnect

The `Disconnect` function disconnects the client from a peer.

func (c *Client) Disconnect(ctx context.Context, targetPeer string) error


**Parameters:**

* `c` - The `Client` object.
* `targetPeer` - The address of the peer to disconnect from.

**Returns:**

An error if the client could not disconnect from the peer.

## Close

The `Close` function closes the client.

func (c *Client) Close() error


**Parameters:**

* `c` - The `Client` object.

**Returns:**

An error if the client could not be closed.

## Download

The `Download` function downloads a file from the network.

func (c *Client) Download(ctx context.Context, cid string, path string) error

**Parameters:**

* `c` - The `Client` object.
* `cid` - The CID of the file to download.
* `path` - The path to the file to download to.

**Returns:**

An error if the file could not be downloaded.

## Upload

The `Upload` function uploads a file to the network.


func (c *Client) Upload(ctx context.Context, path string, reader io.Reader) error

**Parameters:**

* `c` - The `Client` object.
* `path` - The path to the file to upload.
* `reader` - The reader for the file to upload.

**Returns:**

An error if the file could not be uploaded.

