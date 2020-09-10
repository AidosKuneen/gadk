package gadk

import (
	"bytes"
	"github.com/pquerna/ffjson/ffjson"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"github.com/valyala/fasthttp"
)

//API is for calling APIs.
type API struct {
	client *fasthttp.Client
	endpoint string
}

// NewAPI takes an (optional) endpoint and optional http.Client and returns
// an API struct. If an empty endpoint is supplied, then "http://localhost:14265"
// is used.

func NewAPI(endpoint string, c *fasthttp.Client) *API {
	// if c == nil {
	// 	c = ""
	// }

	if endpoint == "" {
		endpoint = "http://localhost:14265/"
	}
	return &API{client: c, endpoint: endpoint}
}

func handleError(err *ErrorResponse, err1, err2 error) error {
	if err.Error != "" {
		return errors.New(err.Error)
	}
	if err.Exception != "" {
		return errors.New(err.Exception)
	}
	if err1 != nil {
		return err1
	}
	return err2
}

func (api *API) do(cmd interface{}, out interface{}) error {
	b, err := ffjson.Marshal(cmd)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
    req.SetBody(b)
	req.SetRequestURI(api.endpoint)


	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := api.client.Do(req, resp); err != nil {
		panic("handle error")
	}
	body := resp.Body()

	// if resp.StatusCode != fasthttp.StatusOK {
	// 	errResp := &ErrorResponse{}
	// 	err = json.Unmarshal(body, errResp)
	// 	return handleError(errResp, err, fmt.Errorf("http status %d while calling API", resp.StatusCode))
	// }
	if bytes.Contains(body, []byte(`"error"`)) || bytes.Contains(body, []byte(`"exception"`)) {
		errResp := &ErrorResponse{}
		err = ffjson.Unmarshal(body, errResp)
		return handleError(errResp, err, fmt.Errorf("unknown error occured while calling API"))
	}

	if out == nil {
		return nil
	}
	return ffjson.Unmarshal(body, out)
}

//ErrorResponse is for occuring exception while calling API.
type ErrorResponse struct {
	Error     string `json:"error"`
	Exception string `json:"exception"`
}

//GetNodeInfoRequest is for GetNodeInfo API request.
type GetNodeInfoRequest struct {
	Command string `json:"command"`
}

//GetNodeInfoResponse is for GetNode APi response.
type GetNodeInfoResponse struct {
	AppName                          string `json:"appName"`
	AppVersion                       string `json:"appVersion"`
	Duration                         int64  `json:"duration"`
	JREAvailableProcessors           int64  `json:"jreAvailableProcessors"`
	JREFreeMemory                    int64  `json:"jreFreeMemory"`
	JREMaxMemory                     int64  `json:"jreMaxMemory"`
	JRETotalMemory                   int64  `json:"jreTotalMemory"`
	LatestMilestone                  Trytes `json:"latestMilestone"`
	LatestMilestoneIndex             int64  `json:"latestMilestoneIndex"`
	LatestSolidSubmeshMilestone      Trytes `json:"latestSolidSubmeshMilestone"`
	LatestSolidSubmeshMilestoneIndex int64  `json:"latestSolidSubmeshMilestoneIndex"`
	Neighbors                        int64  `json:"neighbors"`
	PacketQueueSize                  int64  `json:"packetQueueSize"`
	Time                             int64  `json:"time"`
	Tips                             int64  `json:"tips"`
	TransactionsToRequest            int64  `json:"transactionsToRequest"`
}

//GetNodeInfo calls GetNodeInfo API.
func (api *API) GetNodeInfo() (*GetNodeInfoResponse, error) {
	resp := &GetNodeInfoResponse{}
	err := api.do(map[string]string{
		"command": "getNodeInfo",
	}, resp)
	return resp, err
}

//Neighbor is a part of response of GetNeighbors API.
type Neighbor struct {
	Address                     Address `json:"address"`
	NumberOfAllTransactions     int64   `json:"numberOfAllTransactions"`
	NumberOfInvalidTransactions int64   `json:"numberOfInvalidTransactions"`
	NumberOfNewTransactions     int64   `json:"numberOfNewTransactions"`
}

//GetNeighborsRequest is for GetNeighbors API request.
type GetNeighborsRequest struct {
	Command string `json:"command"`
}

//GetNeighborsResponse is for GetNeighbors API resonse.
type GetNeighborsResponse struct {
	Duration  int64
	Neighbors []Neighbor
}

//GetNeighbors calls GetNeighbors API.
func (api *API) GetNeighbors() (*GetNeighborsResponse, error) {
	resp := &GetNeighborsResponse{}
	err := api.do(map[string]string{
		"command": "getNeighbors",
	}, resp)
	return resp, err
}

//AddNeighborsRequest is for AddNeighbors API request.
type AddNeighborsRequest struct {
	Command string `json:"command"`

	// URIS is an array of strings in the form of "udp://identifier:port"
	// where identifier can be either an IP address or a domain name.
	URIS []string `json:"uris"`
}

//AddNeighborsResponse is for AddNeighbors API resonse.
type AddNeighborsResponse struct {
	Duration       int64 `json:"duration"`
	AddedNeighbors int64 `json:"addedNeighbors"`
}

//AddNeighbors calls AddNeighbors API.
func (api *API) AddNeighbors(uris []string) (*AddNeighborsResponse, error) {
	resp := &AddNeighborsResponse{}
	err := api.do(&struct {
		Command string   `json:"command"`
		URIS    []string `json:"uris"`
	}{
		"addNeighbors",
		uris,
	}, resp)
	return resp, err
}

//RemoveNeighborsRequest is for RemoveNeighbors API request.
type RemoveNeighborsRequest struct {
	Command string `json:"command"`

	// URIS is an array of strings in the form of "udp://identifier:port"
	// where identifier can be either an IP address or a domain name.
	URIS []string `json:"uris"`
}

//RemoveNeighborsResponse is for RemoveNeighbors API resonse.
type RemoveNeighborsResponse struct {
	Duration         int64 `json:"duration"`
	RemovedNeighbors int64 `json:"removedNeighbors"`
}

//RemoveNeighbors calls RemoveNeighbors API.
func (api *API) RemoveNeighbors(uris []string) (*RemoveNeighborsResponse, error) {
	resp := &RemoveNeighborsResponse{}
	err := api.do(&struct {
		Command string   `json:"command"`
		URIS    []string `json:"uris"`
	}{
		"removeNeighbors",
		uris,
	}, resp)
	return resp, err
}

//GetTipsRequest is for GetTipsRequest API request.
type GetTipsRequest struct {
	Command string `json:"command"`
}

//GetTipsResponse is for GetTips API resonse.
type GetTipsResponse struct {
	Duration int64    `json:"duration"`
	Hashes   []Trytes `json:"hashes"`
}

//GetTips calls GetTips API.
func (api *API) GetTips() (*GetTipsResponse, error) {
	resp := &GetTipsResponse{}
	err := api.do(map[string]string{
		"command": "getTips",
	}, resp)
	return resp, err
}

//GetPeerAddressesRequest is for GetPeerAddressesRequest API request.
type GetPeerAddressesRequest struct {
	Command string `json:"command"`
}

//GetPeersAddressesResponse is for GetPeerAddresses API response.
type GetPeerAddressesResponse struct {
	Duration int64    `json:"duration"`
	PeerList   []string `json:"peerlist"`
}

//GetPeerAddresses calls GetPeersAddress API.
func (api *API) GetPeerAddresses(*GetPeerAddressesRequest) (*GetPeerAddressesResponse, error) {
	resp := &GetPeerAddressesResponse{}
	err := api.do(&struct {
		Command string `json:"command"`
	}{
		"getPeerAddresses",
	}, resp)
	return resp, err
}

//FindTransactionsRequest is for FindTransactions API request.
type FindTransactionsRequest struct {
	Command   string    `json:"command"`
	Bundles   []Trytes  `json:"bundles,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
	Tags      []Trytes  `json:"tags,omitempty"`
	Approvees []Trytes  `json:"approvees,omitempty"`
}

//FindTransactionsResponse is for FindTransaction API resonse.
type FindTransactionsResponse struct {
	Duration int64    `json:"duration"`
	Hashes   []Trytes `json:"hashes"`
}

//FindTransactions calls FindTransactions API.
func (api *API) FindTransactions(ft *FindTransactionsRequest) (*FindTransactionsResponse, error) {
	resp := &FindTransactionsResponse{}
	err := api.do(&struct {
		Command string `json:"command"`
		*FindTransactionsRequest
	}{
		"findTransactions",
		ft,
	}, resp)
	return resp, err
}

//GetTrytesRequest is for GetTrytes API request.
type GetTrytesRequest struct {
	Command string   `json:"command"`
	Hashes  []Trytes `json:"hashes"`
}

//GetTrytesResponse is for GetTrytes API resonse.
type GetTrytesResponse struct {
	Duration int64         `json:"duration"`
	Trytes   []Transaction `json:"trytes"`
}

//GetTrytes calls GetTrytes API.
func (api *API) GetTrytes(hashes []Trytes) (*GetTrytesResponse, error) {
	resp := &GetTrytesResponse{}
	err := api.do(&struct {
		Command string   `json:"command"`
		Hashes  []Trytes `json:"hashes"`
	}{
		"getTrytes",
		hashes,
	}, resp)
	return resp, err
}

//GetInclusionStatesRequest is for GetInclusionStates API request.
type GetInclusionStatesRequest struct {
	Command      string   `json:"command"`
	Transactions []Trytes `json:"transactions"`
	Tips         []Trytes `json:"tips"`
}

//GetInclusionStatesResponse is for GetInclusionStates API resonse.
type GetInclusionStatesResponse struct {
	Duration int64  `json:"duration"`
	States   []bool `json:"states"`
}

//GetInclusionStates calls GetInclusionStates API.
func (api *API) GetInclusionStates(tx []Trytes, tips []Trytes) (*GetInclusionStatesResponse, error) {
	resp := &GetInclusionStatesResponse{}
	err := api.do(&struct {
		Command      string   `json:"command"`
		Transactions []Trytes `json:"transactions"`
		Tips         []Trytes `json:"tips"`
	}{
		"getInclusionStates",
		tx,
		tips,
	}, resp)
	return resp, err
}

//Balance is total balance of an Address.
type Balance struct {
	Address Address
	Value   int64
}

//Balances is slice of Balance.
type Balances []Balance

//Total returns the total balance.
func (bs Balances) Total() int64 {
	var total int64
	for _, b := range bs {
		total += b.Value
	}
	return total
}

//GetBalancesRequest is for GetBalances API request.
type GetBalancesRequest struct {
	Command   string    `json:"command"`
	Addresses []Address `json:"addresses"`
	Threshold int64     `json:"threshold"`
}

//GetBalancesResponse is for GetBalances API resonse.
type GetBalancesResponse struct {
	Duration       int64
	Balances       []int64
	Milestone      Trytes
	MilestoneIndex int64
}

//Balances call GetBalances API and returns address-balance pair struct.
func (api *API) Balances(adr []Address) (Balances, error) {
	r, err := api.GetBalances(adr, 100)
	if err != nil {
		return nil, err
	}
	bs := make(Balances, 0, len(adr))
	for i, bal := range r.Balances {
		if bal <= 0 {
			continue
		}
		b := Balance{
			Address: adr[i],
			Value:   bal,
		}
		bs = append(bs, b)
	}
	return bs, nil
}

//GetBalances calls GetBalances API.
func (api *API) GetBalances(adr []Address, threshold int64) (*GetBalancesResponse, error) {
	if threshold <= 0 {
		threshold = 100
	}
	type getBalancesResponse struct {
		Duration       int64    `json:"duration"`
		Balances       []string `json:"balances"`
		Milestone      Trytes   `json:"milestone"`
		MilestoneIndex int64    `json:"milestoneIndex"`
	}

	resp := &getBalancesResponse{}
	err := api.do(&struct {
		Command   string    `json:"command"`
		Addresses []Address `json:"addresses"`
		Threshold int64     `json:"threshold"`
	}{
		"getBalances",
		adr,
		threshold,
	}, resp)
	r := &GetBalancesResponse{
		Duration:       resp.Duration,
		Balances:       make([]int64, len(resp.Balances)),
		Milestone:      resp.Milestone,
		MilestoneIndex: resp.MilestoneIndex,
	}
	for i, ba := range resp.Balances {
		r.Balances[i], err = strconv.ParseInt(ba, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return r, err
}

//GetTransactionsToApproveRequest is for GetTransactionsToApprove API request.
type GetTransactionsToApproveRequest struct {
	Command string `json:"command"`
	Depth   int64  `json:"depth"`
}

//GetTransactionsToApproveResponse is for GetTransactionsToApprove API resonse.
type GetTransactionsToApproveResponse struct {
	Duration          int64  `json:"duration"`
	TrunkTransaction  Trytes `json:"trunkTransaction"`
	BranchTransaction Trytes `json:"branchTransaction"`
}

//GetTransactionsToApprove calls GetTransactionsToApprove API.
func (api *API) GetTransactionsToApprove(depth int64) (*GetTransactionsToApproveResponse, error) {
	resp := &GetTransactionsToApproveResponse{}
	err := api.do(&struct {
		Command string `json:"command"`
		Depth   int64  `json:"depth"`
	}{
		"getTransactionsToApprove",
		depth,
	}, resp)
	return resp, err
}

//AttachToMeshRequest is for AttachToMesh API request.
type AttachToMeshRequest struct {
	Command            string        `json:"command"`
	TrunkTransaction   Trytes        `json:"trunkTransaction"`
	BranchTransaction  Trytes        `json:"branchTransaction"`
	MinWeightMagnitude int64         `json:"minWeightMagnitude"`
	Trytes             []Transaction `json:"trytes"`
}

//AttachToMeshResponse is for AttachToMesh API resonse.
type AttachToMeshResponse struct {
	Duration int64         `json:"duration"`
	Trytes   []Transaction `json:"trytes"`
}

//AttachToMesh calls AttachToMesh API.
func (api *API) AttachToMesh(att *AttachToMeshRequest) (*AttachToMeshResponse, error) {
	resp := &AttachToMeshResponse{}
	err := api.do(&struct {
		Command string `json:"command"`
		*AttachToMeshRequest
	}{
		"attachToMesh",
		att,
	}, resp)
	return resp, err
}

//InterruptAttachingToMeshRequest is for InterruptAttachingToMesh API request.
type InterruptAttachingToMeshRequest struct {
	Command string `json:"command"`
}

//InterruptAttachingToMesh calls InterruptAttachingToMesh API.
func (api *API) InterruptAttachingToMesh() error {
	err := api.do(map[string]string{
		"command": "interruptAttachingToMesh",
	}, nil)
	return err
}

//BroadcastTransactionsRequest is for BroadcastTransactions API request.
type BroadcastTransactionsRequest struct {
	Command string        `json:"command"`
	Trytes  []Transaction `json:"trytes"`
}

//BroadcastTransactions calls BroadcastTransactions API.
func (api *API) BroadcastTransactions(trytes []Transaction) error {
	err := api.do(&struct {
		Command string        `json:"command"`
		Trytes  []Transaction `json:"trytes"`
	}{
		"broadcastTransactions",
		trytes,
	}, nil)
	return err
}

//StoreTransactionsRequest is for StoreTransactions API request.
type StoreTransactionsRequest struct {
	Command string        `json:"command"`
	Trytes  []Transaction `json:"trytes"`
}

//StoreTransactions calls StoreTransactions API.
func (api *API) StoreTransactions(trytes []Transaction) error {
	err := api.do(&struct {
		Command string        `json:"command"`
		Trytes  []Transaction `json:"trytes"`
	}{
		"storeTransactions",
		trytes,
	}, nil)
	return err
}

//GetLatestInclusion takes the most recent solid milestone as returned by getNodeInfo
//and uses it to get the inclusion states of a list of transaction hashes
func (api *API) GetLatestInclusion(hash []Trytes) ([]bool, error) {
	var gt *GetTrytesResponse
	var ni *GetNodeInfoResponse
	var err1 error
	var err2 error
	wd := sync.WaitGroup{}
	wd.Add(2)
	go func() {
		gt, err1 = api.GetTrytes(hash)
		wd.Done()
	}()
	go func() {
		ni, err2 = api.GetNodeInfo()
		wd.Done()
	}()
	wd.Wait()
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	if len(gt.Trytes) == 0 {
		return nil, errors.New("transaction is not found while GetTrytes")
	}

	resp, err := api.GetInclusionStates(hash, []Trytes{ni.LatestMilestone})
	if err != nil {
		return nil, err
	}
	if len(resp.States) == 0 {
		return nil, errors.New("transaction is not found while GetInclusionStates")
	}
	return resp.States, nil
}