# Iris: A Global P2P network for Sharing Threat Intelligence

Iris is a P2P system for collaborative defense proposed by Bc. Martin Řepa developed for Stratosphere lab during his [diploma thesis work](https://www.stratosphereips.org/thesis-projects-list/2022/3/12/global-permissionless-p2p-system-for-sharing-distributed-threat-intelligence).

This repository forks a reference implementation written in Golang using [LibP2P project](https://github.com/libp2p). 

## Authors - CS 4675/6675 Spring 2025 Final Project
Sabina Sokol, Keerthana Thotakura, Maria Jothish, Tran Ha, Seung-a Baek

## Motivation 

Imagine a large financial institution, BankSecure, detects a sophisticated phishing campaign targeting its customers. The attack exploits vulnerabilities in common online banking authentication methods. BankSecure wants to share this intelligence with other banks and cybersecurity organizations to prevent further damage. However, several concerns arise:

> * Reputation Risk: If it becomes public that BankSecure was targeted and potentially compromised, customers may lose trust in the bank's security, leading to financial losses and reputational damage.
> * Competitive Concerns: Sharing detailed threat intelligence could expose weaknesses in BankSecure’s security posture, which competitors might exploit or regulators might scrutinize.
> * Privacy & Legal Risks: The threat intelligence may include sensitive customer metadata, IP addresses, or internal security logs, which must be shared responsibly to avoid legal or regulatory violations.
> * Centralized Trust Issues: Traditional threat intelligence platforms are often controlled by a single entity, which may misuse or monetize the data, further discouraging open participation.

Because of these concerns, BankSecure hesitates to share its findings, delaying critical information that could protect other institutions from similar attacks.

Despite the severity and amount of daily cyberattacks, the best solutions our community has so far are centralized, threat intelligence shared lists; or centralized, commercially-based defense products. No system exists yet to automatically connect endpoints globally and share information about new attacks to improve their security.

Iris allows collaborative defense in cyberspace with an emphasis on security and privacy concerns. It is a pure and completely decentralized P2P network that allows peers to (i) share threat intelligence files, (ii) alert peers about detected attacks, and (iii) ask peers about their opinion on potential attacks. Iris addresses the problem of confidentiality of local threat intelligence data by introducing the concept of Organisations. Organizations are cryptographically-verified and trusted groups of peers within the P2P network. They allow Iris to send content only to pre-trusted groups of peers.

## Related Work

Early centralized P2P networks, like Gnutella, faced major challenges including link congestion, single points of failure, and high administrative overhead. Relying on a single peer to log and transmit data creates vulnerabilities such as inefficiency, tampering risks, and the inability to detect breaches due to a lack of redundancy. In contrast, decentralized P2P systems distribute data collection across multiple nodes, reducing congestion and improving fault tolerance while enabling cross-peer validation for anomaly detection. **IRIS** builds on this decentralized model for threat intelligence sharing by leveraging distributed hash tables (DHTs) and cryptographic keys to ensure confidentiality and trust. It allows peers to securely share alerts while maintaining control over data privacy and dissemination.

## Novelty/Improvements to IRIS Framework
### Motivation
Allows for partial downloads of files
### Strategy
1. Peer Requests Metadata → Learns chunk details including file name, total file size, chunk size, and the list of chunk hashes (for verification).
2. Peer Checks Missing Chunks → Compares existing chunks via hashing.
3. Peer Requests Only Missing Chunks → Saves bandwidth.
4. Sender Sends Requested Chunks → Avoids redundant transfers.
5. Peer Reassembles the File if All Chunks Are Available → Completes file transfer efficiently.

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

### Steps to Initiate Full and Partial File Transfers with File Chunking
1. Open Terminal and Navigate to IRIS folder
> * make network
> * v

2. Open a terminal for each respective peer in the Docker Desktop
> * Peer 1:
>> * docker exec -it iris-peer1-1 /bin/bash
>> * apt-get update; apt-get install redis-tools -y; redis-cli -h iris-redis-1 -p 6379
> * Peer 2: 
>> * docker exec -it iris-peer2-1 /bin/bash
>> * apt-get update; apt-get install redis-tools -y; redis-cli -h iris-redis-1 -p 6379

3. Run the IRIS commands to share intelligence for partial send
> * Peer 1:
>> * PUBLISH gp2p_tl2nl2 '{"type":"tl2nl_file_share","version":1,"data":{"expired_at":1647162651,"severity":"MAJOR","rights":[],"description":{"size":420},"path":"/root/.bashrc","total_size":420,"chunk_size":100,"chunk_count":5,"available_chunks":[0,1,3]}}'

>> * PUBLISH gp2p_tl2nl1 '{"type":"tl2nl_file_share_download","version":1,"data":{"file_id":"QmS4FkBx1uBDHDLASvDocmfo5FXrXgNv4F8WRDkiNTUFe7","chunks":[0,1,3]}}'

> * Peer 2:
>> * DO THIS FIRST
>> * SUBSCRIBE gp2p_tl2nl1
>> * Monitor logs and messages

4. Run the IRIS commands to share intelligence for full send
> * Peer 1:
>> * PUBLISH gp2p_tl2nl2 '{"type":"tl2nl_file_share","version":1,"data":{"expired_at":1647162651,"severity":"MAJOR","rights":[],"description":{"size":420},"path":"/root/.bashrc","total_size":420,"chunk_size":100,"chunk_count":5,"available_chunks":[0,1,2,3,4]}}'
>> * PUBLISH gp2p_tl2nl1 '{"type":"tl2nl_file_share_download","version":1,"data":{"file_id":"QmS4FkBx1uBDHDLASvDocmfo5FXrXgNv4F8WRDkiNTUFe7","chunks":[0,1,2,3,4]}}'

> * Peer 2:
>> * DO THIS FIRST
>> * SUBSCRIBE gp2p_tl2nl1
>> * Monitor logs and messages


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

