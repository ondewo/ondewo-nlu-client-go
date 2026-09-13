// Copyright 2020-2026 ONDEWO GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Everything in this file names the ONDEWO NLU API specifically: its services, one of its
// messages, one of its enums. It is the ONLY file a sibling client (s2t, t2s, sip, csi, vtsi,
// survey) has to rewrite - generated_code_test.go and auth_test.go are product agnostic and
// copy over unchanged.
package tests

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	nlu "github.com/ondewo/ondewo-nlu-client-go/v7/api/ondewo/nlu"
	qa "github.com/ondewo/ondewo-nlu-client-go/v7/api/ondewo/qa"
)

// protoFileCount is the number of .proto files below ondewo-nlu-api/ondewo that the compiler
// consumed. Every one of them has to end up in the global descriptor registry when this package
// is linked; a proto that silently stopped being compiled is otherwise invisible until a
// consumer misses a type.
const protoFileCount = 19

// services is every gRPC service this product exposes, keyed by the fully qualified proto name
// the ServiceDesc must declare.
var services = map[string]*grpc.ServiceDesc{
	"ondewo.nlu.Agents":            &nlu.Agents_ServiceDesc,
	"ondewo.nlu.AiServices":        &nlu.AiServices_ServiceDesc,
	"ondewo.nlu.CcaiProjects":      &nlu.CcaiProjects_ServiceDesc,
	"ondewo.nlu.Contexts":          &nlu.Contexts_ServiceDesc,
	"ondewo.nlu.EntityTypes":       &nlu.EntityTypes_ServiceDesc,
	"ondewo.nlu.Intents":           &nlu.Intents_ServiceDesc,
	"ondewo.nlu.LlmEvaluations":    &nlu.LlmEvaluations_ServiceDesc,
	"ondewo.nlu.Operations":        &nlu.Operations_ServiceDesc,
	"ondewo.nlu.ProjectRoles":      &nlu.ProjectRoles_ServiceDesc,
	"ondewo.nlu.ProjectStatistics": &nlu.ProjectStatistics_ServiceDesc,
	"ondewo.nlu.Rags":              &nlu.Rags_ServiceDesc,
	"ondewo.nlu.ServerStatistics":  &nlu.ServerStatistics_ServiceDesc,
	"ondewo.nlu.Sessions":          &nlu.Sessions_ServiceDesc,
	"ondewo.nlu.Users":             &nlu.Users_ServiceDesc,
	"ondewo.nlu.Utilities":         &nlu.Utilities_ServiceDesc,
	"ondewo.nlu.Webhook":           &nlu.Webhook_ServiceDesc,
	"ondewo.qa.QA":                 &qa.QA_ServiceDesc,
}

// clientConstructors is the generated New<Service>Client of every service above. A client SDK
// that compiles but whose constructors are missing is useless, and the two generators that
// produce them (protoc-gen-go, protoc-gen-go-grpc) can disagree - so both halves are listed.
var clientConstructors = map[string]func(grpc.ClientConnInterface) any{
	"ondewo.nlu.Agents":            func(cc grpc.ClientConnInterface) any { return nlu.NewAgentsClient(cc) },
	"ondewo.nlu.AiServices":        func(cc grpc.ClientConnInterface) any { return nlu.NewAiServicesClient(cc) },
	"ondewo.nlu.CcaiProjects":      func(cc grpc.ClientConnInterface) any { return nlu.NewCcaiProjectsClient(cc) },
	"ondewo.nlu.Contexts":          func(cc grpc.ClientConnInterface) any { return nlu.NewContextsClient(cc) },
	"ondewo.nlu.EntityTypes":       func(cc grpc.ClientConnInterface) any { return nlu.NewEntityTypesClient(cc) },
	"ondewo.nlu.Intents":           func(cc grpc.ClientConnInterface) any { return nlu.NewIntentsClient(cc) },
	"ondewo.nlu.LlmEvaluations":    func(cc grpc.ClientConnInterface) any { return nlu.NewLlmEvaluationsClient(cc) },
	"ondewo.nlu.Operations":        func(cc grpc.ClientConnInterface) any { return nlu.NewOperationsClient(cc) },
	"ondewo.nlu.ProjectRoles":      func(cc grpc.ClientConnInterface) any { return nlu.NewProjectRolesClient(cc) },
	"ondewo.nlu.ProjectStatistics": func(cc grpc.ClientConnInterface) any { return nlu.NewProjectStatisticsClient(cc) },
	"ondewo.nlu.Rags":              func(cc grpc.ClientConnInterface) any { return nlu.NewRagsClient(cc) },
	"ondewo.nlu.ServerStatistics":  func(cc grpc.ClientConnInterface) any { return nlu.NewServerStatisticsClient(cc) },
	"ondewo.nlu.Sessions":          func(cc grpc.ClientConnInterface) any { return nlu.NewSessionsClient(cc) },
	"ondewo.nlu.Users":             func(cc grpc.ClientConnInterface) any { return nlu.NewUsersClient(cc) },
	"ondewo.nlu.Utilities":         func(cc grpc.ClientConnInterface) any { return nlu.NewUtilitiesClient(cc) },
	"ondewo.nlu.Webhook":           func(cc grpc.ClientConnInterface) any { return nlu.NewWebhookClient(cc) },
	"ondewo.qa.QA":                 func(cc grpc.ClientConnInterface) any { return qa.NewQAClient(cc) },
}

// expectedMethods pins RPCs by name. The descriptor cross-check in generated_code_test.go proves
// the two generators agree with each other; it cannot notice an RPC that was renamed upstream,
// because both halves would be renamed together. These are spelled out so that a rename is a
// failing test rather than a silently broken consumer.
var expectedMethods = map[string][]string{
	"ondewo.nlu.Agents":   {"CreateAgent", "GetAgent", "UpdateAgent", "DeleteAgent", "ListAgents"},
	"ondewo.nlu.Contexts": {"ListContexts", "GetContext", "CreateContext", "UpdateContext", "DeleteContext", "DeleteAllContexts"},
	// StreamingDetectIntent is bidirectional, so it lives in ServiceDesc.Streams, not .Methods -
	// the lookup has to consider both.
	"ondewo.nlu.Sessions": {"DetectIntent", "StreamingDetectIntent", "ListSessions", "GetSession"},
	"ondewo.nlu.Users":    {"CreateUser", "GetUser", "DeleteUser", "ListUsers"},
	"ondewo.qa.QA":        {"GetAnswer", "GetServerState", "ListProjectIds"},
}

// TestMessageRoundTripsThroughTheWire is the core assertion about generated message code: a value
// built in go, serialized and parsed back is the same value. It covers a scalar, a repeated
// key/value map of a nested message type and two well-known Timestamps, so a generator that
// mis-numbers a field or loses a nested type fails here.
func TestMessageRoundTripsThroughTheWire(t *testing.T) {
	t.Parallel()

	original := &nlu.Context{
		Name:          "projects/4c9a/agent/sessions/1f2e/contexts/greeting",
		LifespanCount: 7,
		Parameters: map[string]*nlu.Context_Parameter{
			"city": {
				Name:          "city",
				DisplayName:   "City",
				Value:         "Vienna",
				ValueOriginal: "vienna",
				CreatedBy:     "6b1d0c3e-0f3a-4a53-9a4c-1b0a5a6f7c8d",
			},
		},
		CreatedAt:  timestamppb.New(referenceTime),
		ModifiedAt: timestamppb.New(referenceTime),
		CreatedBy:  "6b1d0c3e-0f3a-4a53-9a4c-1b0a5a6f7c8d",
		ModifiedBy: "6b1d0c3e-0f3a-4a53-9a4c-1b0a5a6f7c8d",
	}

	wire, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("proto.Marshal(%T) failed: %v", original, err)
	}
	if len(wire) == 0 {
		t.Fatal("proto.Marshal produced 0 bytes for a fully populated message")
	}

	parsed := &nlu.Context{}
	if err := proto.Unmarshal(wire, parsed); err != nil {
		t.Fatalf("proto.Unmarshal failed: %v", err)
	}

	if !proto.Equal(original, parsed) {
		t.Fatalf("round trip changed the message:\n original = %v\n  parsed = %v", original, parsed)
	}
	if got, want := parsed.GetParameters()["city"].GetValue(), "Vienna"; got != want {
		t.Errorf("nested map value after round trip = %q, want %q", got, want)
	}
	if got, want := parsed.GetCreatedAt().AsTime().UTC(), referenceTime.UTC(); !got.Equal(want) {
		t.Errorf("timestamp after round trip = %v, want %v", got, want)
	}
}

// TestProto3ExplicitPresenceSurvivesTheWire guards the field kind that generators get wrong: a
// proto3 `optional` scalar has to keep the difference between "set to the zero value" and "not
// set". The angular target of the same compiler lost exactly this distinction, which made a
// false/0/"" unsendable; protoc-gen-go models it as a pointer, and this asserts it stays that way.
func TestProto3ExplicitPresenceSurvivesTheWire(t *testing.T) {
	t.Parallel()

	t.Run("zero value set explicitly is transmitted", func(t *testing.T) {
		t.Parallel()

		wire, err := proto.Marshal(&nlu.Context{LifespanTime: proto.Float32(0)})
		if err != nil {
			t.Fatalf("proto.Marshal failed: %v", err)
		}
		if len(wire) == 0 {
			t.Fatal("an explicitly set zero value was dropped from the wire - proto3 presence is lost")
		}

		parsed := &nlu.Context{}
		if err := proto.Unmarshal(wire, parsed); err != nil {
			t.Fatalf("proto.Unmarshal failed: %v", err)
		}
		if parsed.LifespanTime == nil {
			t.Fatal("LifespanTime is nil after the round trip, want a pointer to 0")
		}
		if got := *parsed.LifespanTime; got != 0 {
			t.Errorf("LifespanTime = %v, want 0", got)
		}
	})

	t.Run("unset stays unset", func(t *testing.T) {
		t.Parallel()

		wire, err := proto.Marshal(&nlu.Context{Name: "unset"})
		if err != nil {
			t.Fatalf("proto.Marshal failed: %v", err)
		}

		parsed := &nlu.Context{}
		if err := proto.Unmarshal(wire, parsed); err != nil {
			t.Fatalf("proto.Unmarshal failed: %v", err)
		}
		if parsed.LifespanTime != nil {
			t.Errorf("LifespanTime = %v after a round trip that never set it, want nil", *parsed.LifespanTime)
		}
	})
}

// TestEnumZeroValueIsTheUnspecifiedMember checks the member every proto3 enum must have at 0 and
// the name maps generated beside it. A zero value that is a real choice rather than
// "unspecified" is unrequestable in several of the other clients of this API.
func TestEnumZeroValueIsTheUnspecifiedMember(t *testing.T) {
	t.Parallel()

	var zero nlu.AgentView

	if zero != nlu.AgentView_AGENT_VIEW_UNSPECIFIED {
		t.Errorf("zero value of AgentView = %v, want AGENT_VIEW_UNSPECIFIED", zero)
	}
	if got, want := zero.String(), "AGENT_VIEW_UNSPECIFIED"; got != want {
		t.Errorf("AgentView(0).String() = %q, want %q", got, want)
	}
	if got, want := nlu.AgentView_name[0], "AGENT_VIEW_UNSPECIFIED"; got != want {
		t.Errorf("AgentView_name[0] = %q, want %q", got, want)
	}
	if got, want := nlu.AgentView_value["AGENT_VIEW_FULL"], int32(nlu.AgentView_AGENT_VIEW_FULL); got != want {
		t.Errorf("AgentView_value[AGENT_VIEW_FULL] = %d, want %d", got, want)
	}
	if got, want := int32(nlu.AgentView_AGENT_VIEW_MINIMUM), int32(3); got != want {
		t.Errorf("AGENT_VIEW_MINIMUM = %d, want %d", got, want)
	}
}

// TestUnmarshalRejectsTruncatedInput asserts the generated message reports a parse error instead
// of accepting a malformed payload: field 1 (`name`) is announced as 5 bytes long but only 1
// follows.
func TestUnmarshalRejectsTruncatedInput(t *testing.T) {
	t.Parallel()

	if err := proto.Unmarshal([]byte{0x0a, 0x05, 'a'}, &nlu.Context{}); err == nil {
		t.Fatal("proto.Unmarshal accepted a truncated payload, want an error")
	}
}

// contextsServer is a fake ONDEWO server: it answers GetContext and inherits the "unimplemented"
// behaviour of the generated base type for every other RPC of the service.
type contextsServer struct {
	nlu.UnimplementedContextsServer
}

func (contextsServer) GetContext(_ context.Context, req *nlu.GetContextRequest) (*nlu.Context, error) {
	return &nlu.Context{Name: req.GetName(), LifespanCount: 42}, nil
}

// TestUnaryRPCRoundTripsOverAnInProcessServer drives the generated client stub, the generated
// server stub and the generated ServiceDesc against each other over a real gRPC connection - the
// request is marshalled, routed by the method name baked into the stub, and the response is
// parsed back. Nothing here is mocked except the transport, which is in memory.
func TestUnaryRPCRoundTripsOverAnInProcessServer(t *testing.T) {
	t.Parallel()

	conn := dialInProcess(t, nil, func(srv *grpc.Server) {
		nlu.RegisterContextsServer(srv, contextsServer{})
	})
	client := nlu.NewContextsClient(conn)

	const name = "projects/4c9a/agent/sessions/1f2e/contexts/greeting"
	response, err := client.GetContext(t.Context(), &nlu.GetContextRequest{Name: name})
	if err != nil {
		t.Fatalf("GetContext failed: %v", err)
	}

	if got := response.GetName(); got != name {
		t.Errorf("response name = %q, want %q", got, name)
	}
	if got, want := response.GetLifespanCount(), int32(42); got != want {
		t.Errorf("response lifespan count = %d, want %d", got, want)
	}
}

// TestUnimplementedMethodIsReportedAsUnimplemented pins the other half of the generated server
// contract: an RPC the server does not implement must come back as codes.Unimplemented, not as a
// routing failure or a panic. It also proves the method is routed at all - a method missing from
// the ServiceDesc would surface as a different code.
func TestUnimplementedMethodIsReportedAsUnimplemented(t *testing.T) {
	t.Parallel()

	conn := dialInProcess(t, nil, func(srv *grpc.Server) {
		nlu.RegisterContextsServer(srv, contextsServer{})
	})
	client := nlu.NewContextsClient(conn)

	_, err := client.DeleteAllContexts(t.Context(), &nlu.DeleteAllContextsRequest{})
	if got := status.Code(err); got != codes.Unimplemented {
		t.Fatalf("DeleteAllContexts returned code %v (err = %v), want %v", got, err, codes.Unimplemented)
	}
}
