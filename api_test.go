package gadk

import "testing"

const server = "http://78.46.250.88:15555"

func TestAPIGetNodeInfo(t *testing.T) {
	var err error
	var resp *GetNodeInfoResponse

	for i := 0; i < 5; i++ {
		api := NewAPI(server, nil)
		resp, err = api.GetNodeInfo()
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("GetNodeInfo() expected err to be nil but got %v", err)
	}
	if resp.AppName == "" {
		t.Errorf("GetNodeInfo() returned invalid response: %#v", resp)
	}
}

/*
func TestAPIGetNeighbors(t *testing.T) {
	api := NewAPI(server, nil)

	_, err := api.GetNeighbors()
	if err != nil {
		t.Errorf("GetNeighbors() expected err to be nil but got %v", err)
	}
}

func TestAPIAddNeighbors(t *testing.T) {
	api := NewAPI(server, nil)

	resp, err := api.AddNeighbors([]string{"udp://127.0.0.1:14265/"})
	if err != nil {
		t.Errorf("AddNeighbors([]) expected err to be nil but got %v", err)
	} else if resp.AddedNeighbors != 1 {
		t.Errorf("AddNeighbors([]) expected to add %d got %d", 0, resp.AddedNeighbors)
	}
}

func TestAPIRemoveNeighbors(t *testing.T) {
	api := NewAPI(server, nil)

	resp, err := api.RemoveNeighbors([]string{"udp://127.0.0.1:14265/"})
	if err != nil {
		t.Errorf("RemoveNeighbors([]) expected err to be nil but got %v", err)
	} else if resp.RemovedNeighbors != 1 {
		t.Errorf("RemoveNeighbors([]) expected to remove %d got %d", 0, resp.RemovedNeighbors)
	}
}
func TestAPIGetTips(t *testing.T) {
	api := NewAPI(server, nil)

	resp, err := api.GetTips()
	if err != nil {
		t.Fatalf("GetTips() expected err to be nil but got %v", err)
	}

	if len(resp.Hashes) < 1 {
		t.Errorf("GetTips() returned less than one tip")
	}
	t.Log(len(resp.Hashes))
}
*/
func TestAPIFindTransactions(t *testing.T) {
	var err error
	var resp *FindTransactionsResponse

	ftr := &FindTransactionsRequest{Bundles: []Trytes{"DEXRPLKGBROUQMKCLMRPG9HFKCACDZ9AB9HOJQWERTYWERJNOYLW9PKLOGDUPC9DLGSUH9UHSKJOASJRU"}}
	for i := 0; i < 5; i++ {
		api := NewAPI(server, nil)
		resp, err = api.FindTransactions(ftr)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Errorf("FindTransactions([]) expected err to be nil but got %v", err)
	}
	t.Logf("FindTransactions() = %#v", resp)
}

func TestAPIGetTrytes(t *testing.T) {
	var err error
	var resp *GetTrytesResponse

	for i := 0; i < 5; i++ {
		api := NewAPI(server, nil)
		resp, err = api.GetTrytes([]Trytes{})
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Errorf("GetTrytes([]) expected err to be nil but got %v", err)
	}
	t.Logf("GetTrytes() = %#v", resp)
}

func TestAPIGetInclusionStates(t *testing.T) {
	var err error
	var resp *GetInclusionStatesResponse

	for i := 0; i < 5; i++ {
		api := NewAPI(server, nil)
		resp, err = api.GetInclusionStates([]Trytes{}, []Trytes{})
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Errorf("GetInclusionStates([]) expected err to be nil but got %v", err)
	}
	t.Logf("GetInclusionStates() = %#v", resp)
}

func TestAPIGetBalances(t *testing.T) {
	var err error
	var resp *GetBalancesResponse
	for i := 0; i < 5; i++ {
		api := NewAPI(server, nil)
		resp, err = api.GetBalances([]Address{}, 100)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Errorf("GetBalances([]) expected err to be nil but got %v", err)
	}
	t.Logf("GetBalances() = %#v", resp)
}

func TestAPIGetTransactionsToApprove(t *testing.T) {
	var err error
	var resp *GetTransactionsToApproveResponse

	for i := 0; i < 5; i++ {
		api := NewAPI(server, nil)
		resp, err = api.GetTransactionsToApprove(Depth)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Errorf("GetTransactionsToApprove() expected err to be nil but got %v", err)
	} else if resp.BranchTransaction == "" || resp.TrunkTransaction == "" {
		t.Errorf("GetTransactionsToApprove() return empty branch and/or trunk transactions\n%#v", resp)
	}
}

// TODO Fix test
func TestGetLatestInclusion(t *testing.T) {
	var err error
	var resp []bool

	for i := 0; i < 5; i++ {
		api := NewAPI(server, nil)
		resp, err = api.GetLatestInclusion([]Trytes{"AZBDYWZOARNPMYSJIQCCIXXAA9MTTYHXSMKRRXJKCLOUANBMGBXHL9JB9JFKKXZIHCCQWNHONWCS99999"})
		if err == nil && len(resp) > 0 {
			break
		}
	}
	if err != nil {
		t.Errorf("GetLatestInclustion() expected err to be nil but got %v", err)
	}
	if len(resp) == 0 || !resp[0] {
		t.Error("GetLatestInclustion() is invalid len(resp):", len(resp))
	}
}

/*
func TestAPIInterruptAttachingToMesh(t *testing.T) {
	api := NewAPI(server, nil)

	err := api.InterruptAttachingToMesh()
	if err != nil {
		t.Errorf("InterruptAttachingToMesh() expected err to be nil but got %v", err)
	}
}

// XXX: The following tests are failing because I'd rather not just
//      constantly attach/broadcast/store the same transaction
func TestAPIAttachToMesh(t *testing.T) {
	api := NewAPI(server, nil)

	anr := &AttachToMeshRequest{}
	resp, err := api.AttachToMesh(anr)
	if err != nil {
		t.Errorf("AttachToMesh([]) expected err to be nil but got %v", err)
	}
	t.Logf("AttachToMesh() = %#v", resp)
}

func TestAPIBroadcastTransactions(t *testing.T) {
	api := NewAPI(server, nil)

	err := api.BroadcastTransactions([]Transaction{})
	if err != nil {
		t.Errorf("BroadcastTransactions() expected err to be nil but got %v", err)
	}
}

func TestAPIStoreTransactions(t *testing.T) {
	api := NewAPI(server, nil)

	err := api.StoreTransactions([]Trytes{})
	if err != nil {
		t.Errorf("StoreTransactions() expected err to be nil but got %v", err)
	}
}
*/

func TestAPIGetPeerAddresses(t *testing.T) {
	var err error
	var resp *GetPeerAddressesResponse

	api := NewAPI(server, nil)
	resp, err = api.GetPeerAddresses()

	if err != nil {
		t.Fatalf("GetPeerAddresses() expected err to be nil but got %v", err)
	}
	if len(resp.Peers) == 0 {
		t.Error("GetPeerAddresses() returned invalid response with empty peer list")
	}
}
