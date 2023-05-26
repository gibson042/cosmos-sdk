package group

import (
	"fmt"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	groupv1 "cosmossdk.io/api/cosmos/group/v1"
	ormv1alpha1 "cosmossdk.io/api/cosmos/orm/v1alpha1"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// ORMSchema is the schema for the group module
var ORMSchema = &ormv1alpha1.ModuleSchemaDescriptor{
	SchemaFile: []*ormv1alpha1.ModuleSchemaDescriptor_FileEntry{
		{
			Id:            1,
			ProtoFileName: groupv1.File_cosmos_group_v1_state_proto.Path(),
		},
	},
}

func ProposalToPulsar(proposal Proposal) *groupv1.Proposal {
	var messages []*anypb.Any
	for _, msg := range proposal.Messages {
		messages = append(messages, codectypes.GogoAnyToAnyV2(msg))
	}

	return &groupv1.Proposal{
		Id:                 proposal.Id,
		GroupPolicyAddress: proposal.GroupPolicyAddress,
		Proposers:          proposal.Proposers,
		Metadata:           proposal.Metadata,
		SubmitTime:         timestamppb.New(proposal.SubmitTime),
		GroupVersion:       proposal.GroupVersion,
		GroupPolicyVersion: proposal.GroupPolicyVersion,
		Status:             groupv1.ProposalStatus(proposal.Status),
		VotingPeriodEnd:    timestamppb.New(proposal.VotingPeriodEnd),
		ExecutorResult:     groupv1.ProposalExecutorResult(proposal.ExecutorResult),
		Messages:           messages,
		Title:              proposal.Title,
		Summary:            proposal.Summary,
	}
}

func ProposalFromPulsar(cdc codec.Codec, proposal *groupv1.Proposal) Proposal {
	var messages []*codectypes.Any
	for _, msg := range proposal.Messages {
		messages = append(messages, codectypes.AnyV2ToGogoAny(msg))
	}

	return Proposal{
		Id:                 proposal.Id,
		GroupPolicyAddress: proposal.GroupPolicyAddress,
		Proposers:          proposal.Proposers,
		Metadata:           proposal.Metadata,
		SubmitTime:         proposal.SubmitTime.AsTime(),
		GroupVersion:       proposal.GroupVersion,
		GroupPolicyVersion: proposal.GroupPolicyVersion,
		Status:             ProposalStatus(proposal.Status),
		VotingPeriodEnd:    proposal.VotingPeriodEnd.AsTime(),
		ExecutorResult:     ProposalExecutorResult(proposal.ExecutorResult),
		Messages:           messages,
		Title:              proposal.Title,
		Summary:            proposal.Summary,
	}
}

func GroupInfoFromPulsar(groupInfo *groupv1.GroupInfo) GroupInfo { //nolint:revive // naming is ok
	return GroupInfo{
		Id:          groupInfo.Id,
		Admin:       groupInfo.Admin,
		Version:     groupInfo.Version,
		TotalWeight: groupInfo.TotalWeight,
		Metadata:    groupInfo.Metadata,
		CreatedAt:   groupInfo.CreatedAt.AsTime(),
	}
}

func GroupInfoToPulsar(groupInfo GroupInfo) *groupv1.GroupInfo { //nolint:revive // naming is ok
	return &groupv1.GroupInfo{
		Id:          groupInfo.Id,
		Admin:       groupInfo.Admin,
		Version:     groupInfo.Version,
		TotalWeight: groupInfo.TotalWeight,
		Metadata:    groupInfo.Metadata,
		CreatedAt:   timestamppb.New(groupInfo.CreatedAt),
	}
}

func GroupPolicyInfoFromPulsar(cdc codec.Codec, groupPolicyInfo *groupv1.GroupPolicyInfo) GroupPolicyInfo { //nolint:revive // naming is ok
	result := GroupPolicyInfo{
		Address:             groupPolicyInfo.Address,
		GroupId:             groupPolicyInfo.GroupId,
		Admin:               groupPolicyInfo.Admin,
		Metadata:            groupPolicyInfo.Metadata,
		Version:             groupPolicyInfo.Version,
		GroupPolicySequence: groupPolicyInfo.GroupPolicySequence,
	}

	if groupPolicyInfo.DecisionPolicy != nil {
		decisionPolicy := DecisionPolicy(nil)
		err := cdc.UnpackAny(codectypes.AnyV2ToGogoAny(groupPolicyInfo.DecisionPolicy), &decisionPolicy)
		if err != nil {
			panic(fmt.Sprintf("failed to transform decision policy: %s", err))
		}

		result.DecisionPolicy, err = codectypes.NewAnyWithValue(decisionPolicy)
		if err != nil {
			panic(fmt.Sprintf("failed to transform decision policy: %s", err))
		}
	}

	return result
}

func GroupPolicyInfoToPulsar(groupPolicyInfo GroupPolicyInfo) *groupv1.GroupPolicyInfo { //nolint:revive // naming is ok
	result := &groupv1.GroupPolicyInfo{
		Address:             groupPolicyInfo.Address,
		GroupId:             groupPolicyInfo.GroupId,
		Admin:               groupPolicyInfo.Admin,
		Metadata:            groupPolicyInfo.Metadata,
		Version:             groupPolicyInfo.Version,
		GroupPolicySequence: groupPolicyInfo.GroupPolicySequence,
	}

	if groupPolicyInfo.DecisionPolicy != nil {
		result.DecisionPolicy = codectypes.GogoAnyToAnyV2(groupPolicyInfo.DecisionPolicy)
	}

	return result
}

func GroupMemberFromPulsar(groupMember *groupv1.GroupMember) GroupMember { //nolint:revive // naming is ok
	return GroupMember{
		GroupId:       groupMember.GroupId,
		MemberAddress: groupMember.Member.Address,
		Member: &Member{
			Address:  groupMember.Member.Address,
			Weight:   groupMember.Member.Weight,
			Metadata: groupMember.Member.Metadata,
			AddedAt:  groupMember.Member.AddedAt.AsTime(),
		},
	}
}

func GroupMemberToPulsar(groupMember GroupMember) *groupv1.GroupMember { //nolint:revive // naming is ok
	return &groupv1.GroupMember{
		GroupId:       groupMember.GroupId,
		MemberAddress: groupMember.Member.Address,
		Member: &groupv1.Member{
			Address:  groupMember.Member.Address,
			Weight:   groupMember.Member.Weight,
			Metadata: groupMember.Member.Metadata,
			AddedAt:  timestamppb.New(groupMember.Member.AddedAt),
		},
	}
}

func VoteFromPulsar(vote *groupv1.Vote) Vote {
	return Vote{
		ProposalId: vote.ProposalId,
		Voter:      vote.Voter,
		Option:     VoteOption(vote.Option),
		Metadata:   vote.Metadata,
		SubmitTime: vote.SubmitTime.AsTime(),
	}
}

func VoteToPulsar(vote Vote) *groupv1.Vote {
	return &groupv1.Vote{
		ProposalId: vote.ProposalId,
		Voter:      vote.Voter,
		Option:     groupv1.VoteOption(vote.Option),
		Metadata:   vote.Metadata,
		SubmitTime: timestamppb.New(vote.SubmitTime),
	}
}
