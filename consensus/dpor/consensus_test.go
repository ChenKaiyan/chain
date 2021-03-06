// Copyright 2018 The cpchain authors
// This file is part of the cpchain library.
//
// The cpchain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The cpchain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the cpchain library. If not, see <http://www.gnu.org/licenses/>.

package dpor

import (
	"fmt"
	"testing"

	"bitbucket.org/cpchain/chain/types"
	"github.com/stretchr/testify/assert"
)

func TestDpor_VerifyHeader(t *testing.T) {

	tests := []struct {
		name          string
		verifySuccess bool
		wantErr       bool
	}{
		{"verifyHeader success", true, false},
		{"verifyHeader failed", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Dpor{
				dh: &fakeDporHelper{verifySuccess: tt.verifySuccess},
			}

			err := c.VerifyHeader(&FakeReader{}, newHeader(), true, newHeader())
			fmt.Println("err:", err)
			if err := c.VerifyHeader(&FakeReader{}, newHeader(), true, newHeader()); (err != nil) != tt.wantErr {
				t.Errorf("Dpor.VerifyHeaders() got = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDpor_VerifyHeaders(t *testing.T) {
	tests := []struct {
		name          string
		verifySuccess bool
		wantErr       bool
	}{
		{"verifyHeader success", true, true},
		{"verifyHeader failed", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Dpor{
				dh: &fakeDporHelper{verifySuccess: tt.verifySuccess},
			}
			_, results := c.VerifyHeaders(
				&FakeReader{},
				[]*types.Header{newHeader()},
				[]bool{true},
				[]*types.Header{newHeader()})

			got := <-results
			fmt.Println("got:", got)
			if tt.wantErr != (got == nil) {
				t.Errorf("Dpor.VerifyHeaders() got = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestDpor_APIs(t *testing.T) {
	c := &Dpor{
		dh: &fakeDporHelper{},
	}
	got := c.APIs(nil)
	assert.Equal(t, 1, len(got), "only 1 api should be created")
}

// ========================================================
