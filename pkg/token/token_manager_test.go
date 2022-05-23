/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testingclock "k8s.io/utils/clock/testing"
)

func TestTokenCachingAndExpiration(t *testing.T) {
	type suite struct {
		clock *testingclock.FakeClock
		tg    *fakeTokenGetter
		mgr   *Manager
	}

	type testCase struct {
		name string
		exp  time.Duration
		f    func(t *testing.T, s *suite)
	}

	testCases := []testCase{
		{
			name: "5 minute token not near expiring",
			exp:  time.Minute * 5,
			f: func(t *testing.T, s *suite) {
				s.clock.SetTime(s.clock.Now())

				_, err := s.mgr.GetServiceAccountToken("a", "b", getTokenRequest())

				assert.NoErrorf(t, err, "unexpected error getting token")
				assert.Equal(t, s.tg.count, 1, "expected refresh to not be called, call count was %d", s.tg.count)
			},
		},
		{
			name: "rotate 5 minute token expires in the last minute",
			exp:  time.Minute * 5,
			f: func(t *testing.T, s *suite) {
				s.clock.SetTime(s.clock.Now().Add(4 * time.Minute))

				_, err := s.mgr.GetServiceAccountToken("a", "b", getTokenRequest())

				assert.NoErrorf(t, err, "unexpected error getting token")
				assert.Equal(t, s.tg.count, 2, "expected token to be refreshed, call count was %d", s.tg.count)
			},
		},
		{
			name: "rotate token fails, old token is still valid, doesn't error",
			exp:  time.Hour,
			f: func(t *testing.T, s *suite) {
				s.clock.SetTime(s.clock.Now().Add(4 * time.Minute))
				tg := &fakeTokenGetter{
					err: fmt.Errorf("err"),
				}
				s.mgr.getToken = tg.getToken
				tr, err := s.mgr.GetServiceAccountToken("a", "b", getTokenRequest())

				assert.NoErrorf(t, err, "unexpected error getting token")
				assert.Equal(t, tr.Status.Token, "foo", "unexpected token")
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			clock := testingclock.NewFakeClock(time.Time{}.Add(30 * 24 * time.Hour))
			expSecs := int64(c.exp.Seconds())
			s := &suite{
				clock: clock,
				mgr:   NewManager(nil),
				tg: &fakeTokenGetter{
					tr: &authenticationv1.TokenRequest{
						Spec: authenticationv1.TokenRequestSpec{
							ExpirationSeconds: &expSecs,
						},
						Status: authenticationv1.TokenRequestStatus{
							Token:               "foo",
							ExpirationTimestamp: metav1.Time{Time: clock.Now().Add(c.exp)},
						},
					},
				},
			}
			s.mgr.getToken = s.tg.getToken
			s.mgr.clock = s.clock

			_, err := s.mgr.GetServiceAccountToken("a", "b", getTokenRequest())
			assert.NoErrorf(t, err, "unexpected error getting token")
			assert.Equal(t, s.tg.count, 1, "unexpected client call, call count was %d", s.tg.count)

			_, err = s.mgr.GetServiceAccountToken("a", "b", getTokenRequest())
			assert.NoErrorf(t, err, "unexpected error getting token")
			assert.Equal(t, s.tg.count, 1, "expected token to be served from cache, call count was %d", s.tg.count)

			c.f(t, s)
		})
	}
}

func TestRequiresRefresh(t *testing.T) {
	start := time.Now()

	type testCase struct {
		now, exp      time.Time
		expectRefresh bool
	}

	testCases := []testCase{
		{
			now:           start.Add(1 * time.Minute),
			exp:           start.Add(5 * time.Minute),
			expectRefresh: false,
		},
		{
			now:           start.Add(4 * time.Minute),
			exp:           start.Add(5 * time.Minute),
			expectRefresh: true,
		},
		{
			now:           start.Add(25 * time.Hour),
			exp:           start.Add(60 * time.Hour),
			expectRefresh: true,
		},
		{
			now:           start.Add(10 * time.Minute),
			exp:           start.Add(5 * time.Minute),
			expectRefresh: true,
		},
	}

	for i, c := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			clock := testingclock.NewFakeClock(c.now)
			secs := int64(c.exp.Sub(start).Seconds())
			tr := &authenticationv1.TokenRequest{
				Spec: authenticationv1.TokenRequestSpec{
					ExpirationSeconds: &secs,
				},
				Status: authenticationv1.TokenRequestStatus{
					ExpirationTimestamp: metav1.Time{Time: c.exp},
				},
			}

			mgr := NewManager(nil)
			mgr.clock = clock

			rr := mgr.requiresRefresh(tr)
			assert.Equal(t, rr, c.expectRefresh, "unexpected requiresRefresh result, got: %v, want: %v - %s", rr, c.expectRefresh, c)
		})
	}
}

func TestCleanup(t *testing.T) {
	type testCase struct {
		name              string
		relativeExp       time.Duration
		expectedCacheSize int
	}

	testCases := []testCase{
		{
			name:              "don't cleanup unexpired tokens",
			relativeExp:       -1 * time.Hour,
			expectedCacheSize: 0,
		},
		{
			name:              "cleanup expired tokens",
			relativeExp:       time.Hour,
			expectedCacheSize: 1,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			clock := testingclock.NewFakeClock(time.Time{}.Add(24 * time.Hour))
			mgr := NewManager(nil)
			mgr.clock = clock

			mgr.set("key", &authenticationv1.TokenRequest{
				Status: authenticationv1.TokenRequestStatus{
					ExpirationTimestamp: metav1.Time{Time: mgr.clock.Now().Add(c.relativeExp)},
				},
			})
			mgr.cleanup()

			assert.Equal(t, len(mgr.cache), c.expectedCacheSize, "unexpected number of cache entries after cleanup, got: %d, want: %d", len(mgr.cache), c.expectedCacheSize)
		})
	}
}

type fakeTokenGetter struct {
	count int
	tr    *authenticationv1.TokenRequest
	err   error
}

func (ftg *fakeTokenGetter) getToken(name, namespace string, tr *authenticationv1.TokenRequest) (*authenticationv1.TokenRequest, error) {
	ftg.count++
	return ftg.tr, ftg.err
}

func getTokenRequest() *authenticationv1.TokenRequest {
	return &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			Audiences:         []string{"foo1", "foo2"},
			ExpirationSeconds: getInt64Point(2000),
			BoundObjectRef: &authenticationv1.BoundObjectReference{
				Kind: "pod",
				Name: "foo-pod",
				UID:  "foo-uid",
			},
		},
	}
}

func getInt64Point(v int64) *int64 {
	return &v
}
