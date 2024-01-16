<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
"CLI Breaking" for breaking CLI commands.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog (Agoric fork)

## `v0.45.16-alpha.agoric.3` - 2023-12-04

* (vesting) [#342](https://github.com/agoric-labs/cosmos-sdk/pull/342) recipient can return clawback vesting grant to funder

## `v0.45.16-alpha.agoric.2` - 2023-11-08

### Bug Fixes

* (baseapp) [#337](https://github.com/agoric-labs/cosmos-sdk/pull/337) revert #305 which causes test failures in agoric-sdk

## `v0.45.16-alpha.agoric.1` - 2023-09-22

### Improvements

* Agoric/agoric-sdk#8223 Merge [cosmos/cosmos-sdk v0.45.16](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.45.16)
* (vesting) [#303](https://github.com/agoric-labs/cosmos-sdk/pull/303) Improve vestcalc comments and documentation.

### Bug Fixes

* (snapshots) [#304](https://github.com/agoric-labs/cosmos-sdk/pull/304) raise the per snapshot item limit. Fixes [Agoric/agoric-sdk#8325](https://github.com/Agoric/agoric-sdk/issues/8325)
* (baseapp) [#305](https://github.com/agoric-labs/cosmos-sdk/pull/305) Make sure we don't execute blocks beyond the halt height. Port of [cosmos/cosmos-sdk#16639](https://github.com/cosmos/cosmos-sdk/pull/16639)

## `v0.45.11-alpha.agoric.2` - 2023-03-23

### Improvements

* (snapshot) [#13400](https://github.com/cosmos/cosmos-sdk/pull/13400) Fix snapshot checksum issue in golang 1.19.

### API Breaking Changes

* (store) [#11825](https://github.com/cosmos/cosmos-sdk/pull/11825) Make extension snapshotter interface safer to use, renamed the util function `WriteExtensionItem` to `WriteExtensionPayload`.