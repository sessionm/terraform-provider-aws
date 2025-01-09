// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mediapackagev2_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

// The Default AWS Quota for how many MediaPackage V2 Channel Groups you can have is 3
// We'll serialize the tests to prevent hitting that quota
func TestAccMediaPackageV2_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"MediaPackageChannelGroup": {
			acctest.CtBasic:      testAccMediaPackageChannelGroup_basic,
			"description":        testAccMediaPackageChannelGroup_description,
			acctest.CtDisappears: testAccMediaPackageChannelGroup_disappears,
			//"tags":               testAccMediaPackageChannelGroup_tags,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
