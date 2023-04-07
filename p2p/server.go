package p2p

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ipfs/boxo/routing/http/server"
	"github.com/ipfs/boxo/routing/http/types"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

// waitlist 수신
// waitlist 에 맞게
type Server struct {
	P2P
	ipfsdht *dht.IpfsDHT
	svc     server.ContentRouter
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) FindProviders(ctx context.Context, key cid.Cid) ([]types.ProviderResponse, error) {
	addrsInfo, err := s.ipfsdht.FindProviders(ctx, key)
	if err != nil {
		return nil, err
	}

	res := make([]types.ProviderResponse, 0, len(addrsInfo))
	for _, addr := range addrsInfo {
		bitswapProvider := &types.ReadBitswapProviderRecord{
			Protocol: "p2p",
			Schema:   types.SchemaBitswap,
			ID:       &addr.ID,
		}
		for _, a := range addr.Addrs {
			bitswapProvider.Addrs = append(bitswapProvider.Addrs, types.Multiaddr{Multiaddr: a})
		}
		res = append(res, bitswapProvider)
	}

	return res, nil
}

func (s *Server) ProvideBitswap(ctx context.Context, req *server.BitswapWriteProvideRequest) (time.Duration, error) {
	for _, cid := range req.Keys {
		s.ipfsdht.ProviderStore().AddProvider(ctx, cid.Hash(), peer.AddrInfo{ID: req.ID, Addrs: req.Addrs})
	}

	return req.AdvisoryTTL, nil
}

func (s *Server) Provide(ctx context.Context, req *server.WriteProvideRequest) (types.ProviderResponse, error) {
	return nil, errors.New("not supported")
}

func (s *Server) provide(w http.ResponseWriter, httpReq *http.Request) {
	req := types.WriteProvidersRequest{}
	err := json.NewDecoder(httpReq.Body).Decode(&req)
	_ = httpReq.Body.Close()
	if err != nil {
		writeErr(w, "Provide", http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		return
	}

	resp := types.WriteProvidersResponse{}

	for i, prov := range req.Providers {
		switch v := prov.(type) {
		case *types.WriteBitswapProviderRecord:
			err := v.Verify()
			if err != nil {
				writeErr(w, "Provide", http.StatusForbidden, errors.New("signature verification failed"))
				return
			}

			keys := make([]cid.Cid, len(v.Payload.Keys))
			for i, k := range v.Payload.Keys {
				keys[i] = k.Cid

			}
			addrs := make([]multiaddr.Multiaddr, len(v.Payload.Addrs))
			for i, a := range v.Payload.Addrs {
				addrs[i] = a.Multiaddr
			}
			advisoryTTL, err := s.svc.ProvideBitswap(httpReq.Context(), &server.BitswapWriteProvideRequest{
				Keys:        keys,
				Timestamp:   v.Payload.Timestamp.Time,
				AdvisoryTTL: v.Payload.AdvisoryTTL.Duration,
				ID:          *v.Payload.ID,
				Addrs:       addrs,
			})
			if err != nil {
				writeErr(w, "Provide", http.StatusInternalServerError, fmt.Errorf("delegate error: %w", err))
				return
			}
			resp.ProvideResults = append(resp.ProvideResults,
				&types.WriteBitswapProviderRecordResponse{
					Protocol:    v.Protocol,
					Schema:      v.Schema,
					AdvisoryTTL: &types.Duration{Duration: advisoryTTL},
				},
			)
		case *types.UnknownProviderRecord:
			provResp, err := s.svc.Provide(httpReq.Context(), &server.WriteProvideRequest{
				Protocol: v.Protocol,
				Schema:   v.Schema,
				Bytes:    v.Bytes,
			})
			if err != nil {
				writeErr(w, "Provide", http.StatusInternalServerError, fmt.Errorf("delegate error: %w", err))
				return
			}
			resp.ProvideResults = append(resp.ProvideResults, provResp)
		default:
			writeErr(w, "Provide", http.StatusBadRequest, fmt.Errorf("provider record %d does not contain a protocol", i))
			return
		}
	}
	writeResult(w, "Provide", resp)
}

func (s *Server) findProviders(w http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	cidStr := vars["cid"]
	cid, err := cid.Decode(cidStr)
	if err != nil {
		writeErr(w, "FindProviders", http.StatusBadRequest, fmt.Errorf("unable to parse CID: %w", err))
		return
	}
	providers, err := s.svc.FindProviders(httpReq.Context(), cid)
	if err != nil {
		writeErr(w, "FindProviders", http.StatusInternalServerError, fmt.Errorf("delegate error: %w", err))
		return
	}
	response := types.ReadProvidersResponse{Providers: providers}
	writeResult(w, "FindProviders", response)
}

func writeResult(w http.ResponseWriter, method string, val any) {
	w.Header().Add("Content-Type", "application/json")

	// keep the marshaling separate from the writing, so we can distinguish bugs (which surface as 500)
	// from transient network issues (which surface as transport errors)
	b, err := marshalJSONBytes(val)
	if err != nil {
		writeErr(w, method, http.StatusInternalServerError, fmt.Errorf("marshaling response: %w", err))
		return
	}

	_, err = io.Copy(w, bytes.NewBuffer(b))
	if err != nil {
		return
	}
}

func writeErr(w http.ResponseWriter, method string, statusCode int, cause error) {
	w.WriteHeader(statusCode)
	causeStr := cause.Error()
	if len(causeStr) > 1024 {
		causeStr = causeStr[:1024]
	}
	_, err := w.Write([]byte(causeStr))
	if err != nil {
		return
	}
}

// marshalJSONBytes is needed to avoid changes
// on the original bytes due to HTML escapes.
func marshalJSONBytes(val any) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}

	// remove last \n added by Encode
	return buf.Bytes()[:buf.Len()-1], nil
}
