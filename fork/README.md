*(Submitting as a PR since this repo has its issue tracker disabled)*

On behalf of Coinbase and others preparing for the Ethereum merge, we request the team behind
EthereumPoW clarify important but currently underspecified aspects of their fork chain and its
activation. At the time of this writing, this EthereumPoW github repo leaves
[chainid unchanged](https://github.com/ethereumpow/go-ethereum/blob/704b52ec1e3eba0605d3d5ade0b9ae8562d879de/params/config.go#L60).
While a code comment suggests an intent to change chainid, we are unaware of any
concrete communications around how the change will be implemented and ultimately activated. As the
team must know, using the same chainid as the PoS Ethereum chain post-merge will allow transactions
from one chain to be replayed on the other. These replay attacks place funds at serious risk, and
many users will not have the time or expertise to protect themselves.

We have seen allusions from some EthereumPoW communication channels
([Discord](https://discord.gg/invite/nthXNEv), [Twitter](https://twitter.com/EthereumPoW), and
[Medium](https://medium.com/@ETHW)) that there is an intent to change the EthereumPoW chainid no
later than the merge point. We strongly support this approach as it would prohibit any possibility
of transaction replaying. But these same communication channels (and the code in this repo)
indicate that no such fork implementation would be available let alone activate prior to it. Unless
the intent is an effective EthereumPoW chain halt until the fork implementation is formally
released, these statements appear contradictory.

In order to clear up this uncertainty and allow us to take appropriate actions to protect user
funds, we request that the EthereumPoW team promptly and authoritatively communicate if and when an
EthereumPoW chainid change is to take place, and how the change will be deployed and activated.
