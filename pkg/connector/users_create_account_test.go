package connector

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	createTestOrgID    = "org-1"
	createTestScimPath = "/api/v1/scim/Users"
	createTestMembers  = "/v2/orgs/" + createTestOrgID + "/members"
	scimErrorSchema    = `"schemas":["urn:ietf:params:scim:api:messages:2.0:Error"]`
	scimUniqueness     = `{` + scimErrorSchema + `,"status":"409","scimType":"uniqueness","detail":"User already exists"}`
)

// redirectTransport sends every request to the fake server, keeping the path
// and query, so the hard-coded SCIM and REST hosts both reach it.
type redirectTransport struct{ target *url.URL }

func (t redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.URL.Scheme = t.target.Scheme
	r.URL.Host = t.target.Host
	r.Host = t.target.Host
	return http.DefaultTransport.RoundTrip(r)
}

type fakeMiro struct {
	createStatus int
	createBody   string
	members      []miro.User
	membersFail  int

	mu           sync.Mutex
	createCalls  int
	lookups      []url.Values
	fullScans    int
	cachedLookup bool
	otherCalls   []string
}

func (f *fakeMiro) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodPost && r.URL.Path == createTestScimPath:
		f.createCalls++
		w.Header().Set("X-Ratelimit-Limit", "100")
		w.Header().Set("X-Ratelimit-Remaining", "42")
		w.WriteHeader(f.createStatus)
		_, _ = io.WriteString(w, f.createBody)
	case r.Method == http.MethodGet && r.URL.Path == createTestMembers:
		q := r.URL.Query()
		if q.Has("emails") {
			f.lookups = append(f.lookups, q)
		}
		if r.Header.Get("Cache-Control") != "no-cache" && f.createCalls > 0 {
			f.cachedLookup = true
		}
		w.Header().Set("X-Ratelimit-Limit", "100")
		w.Header().Set("X-Ratelimit-Remaining", "7")
		if f.membersFail != 0 {
			w.WriteHeader(f.membersFail)
			_, _ = io.WriteString(w, `{"status":500,"code":"internalError","message":"boom"}`)
			return
		}
		var data []miro.User
		resp := miro.GetOrganizationMembersResponse{Limit: 50}
		if !q.Has("emails") {
			// Unfiltered listing: one member per page to exercise paging.
			if f.createCalls > 0 {
				f.fullScans++
			}
			idx := 0
			if c := q.Get("cursor"); c != "" {
				idx, _ = strconv.Atoi(c)
			}
			if idx < len(f.members) {
				data = f.members[idx : idx+1]
				if idx+1 < len(f.members) {
					resp.Cursor = strconv.Itoa(idx + 1)
				}
			}
		} else {
			wanted := map[string]bool{}
			for _, e := range strings.Split(q.Get("emails"), ",") {
				wanted[e] = true
			}
			for _, m := range f.members {
				// Model an exact, case-sensitive provider filter.
				if wanted[m.Email] {
					data = append(data, m)
				}
			}
		}
		resp.Data = data
		_ = json.NewEncoder(w).Encode(resp)
	default:
		f.otherCalls = append(f.otherCalls, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}
}

func newCreateAccountTestBuilder(t *testing.T, f *fakeMiro) *userBuilder {
	t.Helper()
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	target, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	hc := &http.Client{Transport: redirectTransport{target: target}}
	return newUserBuilder(miro.New(hc, hc, srv.URL), createTestOrgID)
}

func createAccountInfo(t *testing.T, email string) *v2.AccountInfo {
	t.Helper()
	profile, err := structpb.NewStruct(map[string]any{
		"email":      email,
		"first_name": "Jane",
		"last_name":  "Doe",
	})
	if err != nil {
		t.Fatal(err)
	}
	return &v2.AccountInfo{Profile: profile}
}

var existingMember = miro.User{
	Id:             "3458764500000000001",
	Type:           "organization-member",
	Active:         true,
	License:        "full",
	Role:           "organization_internal_user",
	Email:          "jane.doe@example.com",
	LastActivityAt: "2025-01-02T03:04:05Z",
}

// syncedResource returns what List emits for id from the same fake.
func syncedResource(t *testing.T, b *userBuilder, id string) *v2.Resource {
	t.Helper()
	token := ""
	for {
		resources, res, err := b.List(context.Background(), nil, rs.SyncOpAttrs{PageToken: pagination.Token{Token: token}})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		for _, r := range resources {
			if r.GetId().GetResource() == id {
				return r
			}
		}
		if res.NextPageToken == "" {
			break
		}
		token = res.NextPageToken
	}
	t.Fatalf("List() did not emit %s", id)
	return nil
}

func TestCreateAccountConflictReturnsExistingAccount(t *testing.T) {
	other := existingMember
	other.Id = "3458764500000000002"
	other.Email = "someone.else@example.com"

	f := &fakeMiro{
		createStatus: http.StatusConflict,
		createBody:   scimUniqueness,
		members:      []miro.User{other, existingMember},
	}
	b := newCreateAccountTestBuilder(t, f)

	// Requested with different casing than Miro stores.
	resp, plaintext, annos, err := b.CreateAccount(context.Background(), createAccountInfo(t, "Jane.Doe@Example.com"), nil)
	if err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}
	rateLimits := 0
	for _, a := range annos {
		if a.MessageIs(&v2.RateLimitDescription{}) {
			rateLimits++
		}
	}
	rl := &v2.RateLimitDescription{}
	if ok, aerr := annos.Pick(rl); aerr != nil || !ok || rateLimits != 1 || rl.GetRemaining() != 7 {
		t.Fatalf("want one rate-limit annotation from the latest lookup (remaining 7), got %d (remaining %d)", rateLimits, rl.GetRemaining())
	}
	if len(plaintext) != 0 {
		t.Fatalf("CreateAccount() returned %d plaintext credentials for an account it did not create", len(plaintext))
	}
	result, ok := resp.(*v2.CreateAccountResponse_AlreadyExistsResult)
	if !ok {
		t.Fatalf("CreateAccount() response = %T, want *v2.CreateAccountResponse_AlreadyExistsResult", resp)
	}
	if f.fullScans != 0 {
		t.Fatalf("full scans = %d, want 0 when the email filter matches", f.fullScans)
	}
	if f.cachedLookup {
		t.Fatal("member lookup after the conflict must bypass the response cache")
	}
	if want := syncedResource(t, b, existingMember.Id); !proto.Equal(result.GetResource(), want) {
		t.Fatalf("AlreadyExistsResult resource differs from synced resource:\n got: %v\nwant: %v", result.GetResource(), want)
	}
	if f.createCalls != 1 {
		t.Fatalf("create calls = %d, want 1", f.createCalls)
	}
	if len(f.lookups) == 0 {
		t.Fatal("expected an organization member lookup by email")
	}
	if len(f.otherCalls) != 0 {
		t.Fatalf("unexpected calls (existing account must not be modified): %v", f.otherCalls)
	}
}

func TestCreateAccountConflictFallsBackToCaseInsensitiveScan(t *testing.T) {
	stored := existingMember
	stored.Email = "Jane.Doe@Example.com"
	other := existingMember
	other.Id = "3458764500000000002"
	other.Email = "someone.else@example.com"

	tests := []struct {
		name string
		body string
	}{
		{name: "uniqueness scimType", body: scimUniqueness},
		{name: "numeric status without scimType", body: `{` + scimErrorSchema + `,"status":409,"detail":"User already exists"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakeMiro{createStatus: http.StatusConflict, createBody: tt.body, members: []miro.User{other, stored}}
			b := newCreateAccountTestBuilder(t, f)

			// The fake's emails filter is case-sensitive, so only a scan finds it.
			resp, plaintext, _, err := b.CreateAccount(context.Background(), createAccountInfo(t, "jane.doe@example.com"), nil)
			if err != nil {
				t.Fatalf("CreateAccount() error = %v", err)
			}
			if len(plaintext) != 0 {
				t.Fatalf("unexpected plaintext credentials: %v", plaintext)
			}
			result, ok := resp.(*v2.CreateAccountResponse_AlreadyExistsResult)
			if !ok {
				t.Fatalf("CreateAccount() response = %T, want *v2.CreateAccountResponse_AlreadyExistsResult", resp)
			}
			if f.fullScans != 2 {
				t.Fatalf("full scan pages = %d, want 2", f.fullScans)
			}
			if f.cachedLookup {
				t.Fatal("member lookups after the conflict must bypass the response cache")
			}
			if want := syncedResource(t, b, stored.Id); !proto.Equal(result.GetResource(), want) {
				t.Fatalf("AlreadyExistsResult resource differs from synced resource:\n got: %v\nwant: %v", result.GetResource(), want)
			}
		})
	}
}

func TestCreateAccountConflictScanIsBounded(t *testing.T) {
	members := make([]miro.User, maxConflictScanPages+2)
	for i := range members {
		members[i] = existingMember
		members[i].Id = "member-" + strconv.Itoa(i)
		members[i].Email = "member" + strconv.Itoa(i) + "@example.com"
	}
	// The fake serves one member per page, so this one sits past the cap.
	members[len(members)-1].Email = "Jane.Doe@Example.com"

	f := &fakeMiro{createStatus: http.StatusConflict, createBody: scimUniqueness, members: members}
	b := newCreateAccountTestBuilder(t, f)

	resp, _, _, err := b.CreateAccount(context.Background(), createAccountInfo(t, "jane.doe@example.com"), nil)
	if err == nil || resp != nil {
		t.Fatalf("CreateAccount() = %v, %v; want bounded no-match error", resp, err)
	}
	if got := status.Code(err); got != codes.AlreadyExists {
		t.Fatalf("status code = %v, want AlreadyExists (err: %v)", got, err)
	}
	if !strings.Contains(err.Error(), "members scanned") {
		t.Fatalf("error %q does not mention the bounded scan", err)
	}
	if f.fullScans != maxConflictScanPages {
		t.Fatalf("scan pages = %d, want %d", f.fullScans, maxConflictScanPages)
	}
}

func TestCreateAccountConflictWithoutUniqueMatch(t *testing.T) {
	dup := existingMember
	dup.Id = "3458764500000000003"

	tests := []struct {
		name    string
		members []miro.User
		want    string
	}{
		{name: "no match", members: nil, want: "no organization member matches"},
		{name: "ambiguous", members: []miro.User{existingMember, dup}, want: "multiple organization members match"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakeMiro{createStatus: http.StatusConflict, createBody: scimUniqueness, members: tt.members}
			b := newCreateAccountTestBuilder(t, f)

			resp, plaintext, _, err := b.CreateAccount(context.Background(), createAccountInfo(t, "jane.doe@example.com"), nil)
			if err == nil {
				t.Fatalf("CreateAccount() error = nil, response = %v", resp)
			}
			if got := status.Code(err); got != codes.AlreadyExists {
				t.Fatalf("status code = %v, want AlreadyExists (err: %v)", got, err)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err, tt.want)
			}
			if resp != nil || len(plaintext) != 0 {
				t.Fatalf("expected no response or credentials, got %v / %v", resp, plaintext)
			}
			if len(f.lookups) == 0 {
				t.Fatal("expected an organization member lookup by email")
			}
			if len(f.otherCalls) != 0 {
				t.Fatalf("unexpected calls: %v", f.otherCalls)
			}
		})
	}
}

func TestCreateAccountConflictLookupFailure(t *testing.T) {
	f := &fakeMiro{createStatus: http.StatusConflict, createBody: scimUniqueness, membersFail: http.StatusInternalServerError}
	b := newCreateAccountTestBuilder(t, f)

	resp, _, _, err := b.CreateAccount(context.Background(), createAccountInfo(t, "jane.doe@example.com"), nil)
	if err == nil || resp != nil {
		t.Fatalf("CreateAccount() = %v, %v; want lookup error", resp, err)
	}
	if !strings.Contains(err.Error(), "look up existing") {
		t.Fatalf("error %q does not describe the failed lookup", err)
	}
}

func TestCreateAccountOtherErrorsUnchanged(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   codes.Code
	}{
		{name: "bad request", status: http.StatusBadRequest, body: `{` + scimErrorSchema + `,"status":"400","scimType":"invalidValue","detail":"bad"}`, want: codes.InvalidArgument},
		{name: "forbidden", status: http.StatusForbidden, body: `{"status":"403","detail":"nope"}`, want: codes.PermissionDenied},
		{name: "rate limited", status: http.StatusTooManyRequests, body: `{"status":"429","detail":"slow down"}`, want: codes.Unavailable},
		{name: "server error", status: http.StatusInternalServerError, body: `{"status":"500","detail":"boom"}`, want: codes.Unavailable},
		{name: "non-uniqueness conflict", status: http.StatusConflict, body: `{` + scimErrorSchema + `,"status":"409","scimType":"mutability","detail":"immutable"}`, want: codes.AlreadyExists},
		{name: "conflict without a SCIM error body", status: http.StatusConflict, body: `not json`, want: codes.AlreadyExists},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakeMiro{createStatus: tt.status, createBody: tt.body, members: []miro.User{existingMember}}
			b := newCreateAccountTestBuilder(t, f)

			resp, plaintext, annos, err := b.CreateAccount(context.Background(), createAccountInfo(t, "jane.doe@example.com"), nil)
			if err == nil {
				t.Fatalf("CreateAccount() error = nil, response = %v", resp)
			}
			rl := &v2.RateLimitDescription{}
			if ok, aerr := annos.Pick(rl); aerr != nil || !ok {
				t.Fatalf("rate-limit annotation missing from error response (ok=%v, err=%v)", ok, aerr)
			}
			if tt.status == http.StatusTooManyRequests && rl.GetStatus() != v2.RateLimitDescription_STATUS_OVERLIMIT {
				t.Fatalf("rate-limit status = %v, want OVERLIMIT", rl.GetStatus())
			}
			if got := status.Code(err); got != tt.want {
				t.Fatalf("status code = %v, want %v (err: %v)", got, tt.want, err)
			}
			if !strings.Contains(err.Error(), "failed to create miro user") {
				t.Fatalf("error %q lost its original context", err)
			}
			if resp != nil || len(plaintext) != 0 {
				t.Fatalf("expected no response or credentials, got %v / %v", resp, plaintext)
			}
			if len(f.lookups) != 0 {
				t.Fatalf("unexpected member lookups: %v", f.lookups)
			}
		})
	}
}

func TestCreateAccountSuccessUnchanged(t *testing.T) {
	f := &fakeMiro{
		createStatus: http.StatusCreated,
		createBody:   `{"id":"3458764500000000009","type":"user","active":true,"email":"jane.doe@example.com"}`,
	}
	b := newCreateAccountTestBuilder(t, f)

	resp, plaintext, _, err := b.CreateAccount(context.Background(), createAccountInfo(t, "jane.doe@example.com"), nil)
	if err != nil {
		t.Fatalf("CreateAccount() error = %v", err)
	}
	if len(plaintext) != 0 {
		t.Fatalf("unexpected plaintext credentials: %v", plaintext)
	}
	result, ok := resp.(*v2.CreateAccountResponse_SuccessResult)
	if !ok {
		t.Fatalf("CreateAccount() response = %T, want *v2.CreateAccountResponse_SuccessResult", resp)
	}
	if got := result.GetResource().GetId().GetResource(); got != "3458764500000000009" {
		t.Fatalf("resource id = %q", got)
	}
	if len(f.lookups) != 0 {
		t.Fatalf("unexpected member lookups: %v", f.lookups)
	}
}
