# TAK (Teleport Attack Kit)

## What is it?

TAK is a collection of tools demonstrated at Blackhat USA 2026 as part of the talk "Beam Me Up, Luke - A Review of Teleport Attack Scenarios".

This tool collection aims to help researchers explore the Teleport platform, and combined allow some pretty cool attack scenarios.

## What is included?

### certificate-tool

`certificate-tool` is a tool used for generating various types of certificates for use with Teleport.

### log-viewer

`log-viewer` is a tool used exploiting a flaw in Teleport in which a single Node service credentials can access all Teleport audit logs for the entire cluster's servers.

### node-hijack

`node-hijack` is a tool used for exploiting a flaw in Teleport in which a single compromised Node's credentials can be used to modify Auth Server metadata for other Node services.

This leads to the ability for a malicious Node to hijack any node within the Teleport cluster.

### reverse-tunnel

`reverse-tunnel` is a tool used for providing the reverse-tunnel functionality for Teleport.

### tunnel-manager

`tunnel-manager` is a SOCAT like tool which establishes new TCP tunnels using a specified ALPN value.

### IronRDP

A fork of the `IronRDP` client to include Teleport's virtual smartcard for exploiting Teleport certificate attacks.
