//go:build integration

package versions_test

import (
	"encoding/json"
	"fmt"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"
	errors "github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Negative: Get DIDDoc versions", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing utils.DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	expectedDidDereferencing := testCase.ExpectedResult.(utils.DereferencingResult)
	utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"cannot get DIDDoc versions with an existent DID, but not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/versions",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleMainnetDid,
			),
			ResolutionType: string(types.TEXT),
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc versions with not existent DID and not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/versions",
				testconstants.TestHostAddress,
				testconstants.NotExistentMainnetDid,
			),
			ResolutionType: string(types.TEXT),
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc versions with not existent DID",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/versions",
				testconstants.TestHostAddress,
				testconstants.NotExistentTestnetDid,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc versions with invalid DID",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/versions",
				testconstants.TestHostAddress,
				testconstants.InvalidDid,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDid,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.MethodNotSupportedHttpCode,
		},
	),
)
