# Iris: A Global P2P network for Sharing Threat Intelligence

Iris is a P2P system for collaborative defense proposed by Bc. Martin Řepa developed for Stratosphere lab during his [diploma thesis work](https://www.stratosphereips.org/thesis-projects-list/2022/3/12/global-permissionless-p2p-system-for-sharing-distributed-threat-intelligence).

This repository hosts a reference implementation written in Golang using [LibP2P project](https://github.com/libp2p) along with integration of Iris into [Slips IPS](https://github.com/draliii/StratosphereLinuxIPS) and [Fides Trust Model](https://github.com/lukasforst/fides). 

This project is funded by [NlNet NGI Zero Entrust](https://nlnet.nl/project/Iris-P2P/)


For more details regarding design please see [Design](docs/Design.md). For the architecture/implementation, we refer the reader to [Architecture](docs/architecture.md) or the thesis itself.

### Motivation 

Despite the severity and amount of daily cyberattacks, the best solutions our community has so far are
centralized, threat intelligence shared lists; or centralized, commercially-based defense products.
No system exists yet to automatically connect endpoints globally and share information about new attacks
to improve their security. 

Iris allows collaborative defense in cyberspace with an emphasis on security and privacy concerns.
It is a pure and completely decentralized P2P network that allows peers to (i) share threat intelligence
files, (ii) alert peers about detected attacks, and (iii) ask peers about their opinion on potential
attacks. Iris addresses the problem of confidentiality of local threat intelligence data by
introducing the concept of _Organisations_. Organizations are cryptographically-verified and
trusted groups of peers within the P2P network. They allow Iris to send content only
to pre-trusted groups of peers.

## Dependencies

To run a standalone peer, you need:
* a running redis instance
* golang (1.17) 

### Install Go 1.17
* Go to https://go.dev/dl/ and download/run the installer for your machine
* Run `ls /usr/local/go` in a terminal (for Mac/Linux) or `dir "C:\Go"` in Command Prompt for Windows
* Next, check if `go` is in your path with `echo $PATH` for Mac/Linux or `echo %PATH%` for Windows
* If `/usr/local/go/bin` (Mac/Linux) or `C:\Go\bin` (Windows) is missing, run `echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc; source ~/.bashrc  # or source ~/.zshrc if using zsh` (Mac/Linux) or add `C:\Go\bin` manually to System variables (Windows)
* Verify your installation with `go version`

## Ongoing Improvements: Modifying the file-sharing protocol to send fixed-sized chunks
### Authors
Sabina Sokol, Keerthana Thotakura, Maria Jothish, Tran Ha, Seung-a Baek
### Motivation
Allows for partial downloads of files
### Strategy
1. Peer Requests Metadata → Learns chunk details including file name, total file size, chunk size, and the list of chunk hashes (for verification).
2. Peer Checks Missing Chunks → Compares existing chunks via hashing.
3. Peer Requests Only Missing Chunks → Saves bandwidth.
4. Sender Sends Requested Chunks → Avoids redundant transfers.
5. Peer Reassembles the File if All Chunks Are Available → Completes file transfer efficiently.

## User Guide

### OrgSig Tool

To manage the P2P TI sharing withing an organization, we developed a tool called **orgsig**. Orgsig is a small program written in Golang that can generate organizations or sign existing peers ID using the already generated organisations.

```bash
> make orgsig 
go build cmd/orgsig.go
>  ./orgsig --help
Running v0.0.1 orgsig

Usage of ./orgsig:
  -load-key-path string
    	Path to a file with organisation private key. If not set, new private-key is generated.
  -peer-id string
    	Public ID of a peer to sign. Flag --sign-peer must be set for this option to be valid.
  -save-key-path string
    	If set, value will be used as a path to save organisation private-key.
  -sign-peer
    	Flag to sign peer ID. Flag peer-id can be used to set peerID, otherwise, cli will ask. The signature will be printed to stdout.
```


### Running a Peer

Starting a peer with reference configuration is as simple as running (assuming a Redis instance is running on the localhost):

> make run

![image](https://github.com/user-attachments/assets/bf739119-c699-4125-9a7e-63452ed94161)


### Debugging, Running Multiple Peers

To run multiple peers simultaneously, you can use an already prepared docker-compose file with pre-configured 4 peers.
The network of 4 peers can be started with (note that you must have `docker` and `docker-compose` installed):

```bash
> make network
```

This command starts docker-compose with 4 peers in separate containers and one container with a separate Redis instance. 
Every peer connects to a different Redis channel and waits for messages from Fides. The peers will connect to each other and thus form a small network. The configuration files of every peer can be found in the [dev/](dev) directory. 

To interact with the peers, you must act as Fides Trust Model and send the peers a manual message by publishing some of them through the Redis channels. Example PUBLISH commands can be found in [dev/redisobj.dev](dev/redisobj.dev).

![image](https://github.com/user-attachments/assets/077b660d-041e-43ff-9cf1-9a21fbccac0a)

## Todo/Future Work:
* Signal handling for graceful shutdown
* After a peer connects to the network, search immediately for members of trustworthy organisations. So far, only `connector` does it.
* Implement message (bytes?) rate-limiting per individual peers to mitigate flooding attacks (or adaptive gossip?)
* Use more the Reporting Protocol to report misbehaving peers
* Implement purging of keys after some time (configurable?) in peers' message cache
* responseStorage goroutines should not wait for responses from peers that disconnected during the waiting. Otherwise, when that happens, it's gonna unnecessarily wait until the timeout occurs
* storageResponse goroutines should wait only for responses from peers where requests were successfully sent (err was nil)
* implement purging of file metadata after files expire (viz currently not used field `ElapsedAt`)
* Is the reference basic manager really trimming peers based on their reliability? Need to be checked,

# Related projects
This project is now a submodule of [StratosphereLinuxIPS](https://github.com/stratosphereips/StratosphereLinuxIPS/)

