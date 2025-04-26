# Lambda API
 
#### Official API for the [Lambda utility mod](https://github.com/lambda-client/lambda)

## Routes

Statuses:
- ✅ Supported
- 🟡 Deprecated; To be removed soon
- ❌ Unsupported; Not working anymore
- ⏳ Not supported yet

| Route         | Method | Status |                                                                                                                   Information |
|:--------------|:------:|--------|------------------------------------------------------------------------------------------------------------------------------:|
| /login        |  POST  | ✅      |                                           Allows a player to authenticate using a Minecraft username and Mojang session hash. |
| /link/discord |  POST  | ✅      |                                      Links a Discord account to an existing authentication token using a Discord OAuth token. |
|               |        |        |                                                                                                                               |
| /party/create |  POST  | ⏳      | Creates a new party and associates it with the player. <br/> Automatically deletes or leaves existing parties for the player. |
| /party/join   |  PUT   | ⏳      |                                                               Allows a player to join an existing party using a party secret. |~~
| /party/leave  |  PUT   | ⏳      |                                 Removes a player from their current party.<br/>Deletes the party if the player is the leader. |
| /party/delete | DELETE | ⏳      |                                                                        Deletes an existing party if the player is the leader. |
| /party        |  GET   | ⏳      |                                                                           Retrieves details about the player's current party. |
| /party/listen |  GET   | 🟡     |                                                             Streams party updates to the player via Server-Sent Events (SSE). |
|               |        |        |                                                                                                                               |
| /cape         |  GET   | ✅      |                                                                                             Fetches the ID of a player's cape |
| /cape         |  PUT   | ✅      |                                                                                      Updates the player's cape with a new ID. |
| /capes        |  GET   | ✅      |                                                                                     Retrieves cape ids from a list of players |

## Authentication

We understand that not all players trust us to handle account tokens, especially Microsoft ones.
<br>
This is why we've decided to use a unique authentication method that ensure security and anonymity

### How does it work ?
This system authenticates players by verifying their session with Mojang's servers, ensuring they own the Minecraft account they claim to. It is only effective in online-mode servers—offline servers cannot be trusted for this verification.

1. Server sends RSA public key → Client generates AES-128 secret key and SHA-1 session hash (secret + server key).
2. Encrypted handshake with server's public key and client's secret to prevents MITM attacks.
3. Client posts hash to Mojang (https://sessionserver.mojang.com/session/minecraft/join), valid for 5 minutes to prevent replay attacks

[![](https://mermaid.ink/img/pako:eNp9U01z2jAQ_Ss7OuRSPowdsPGhDLFJIG3aDHQ60xgOqr2AEluikpyWEv57hWXaQjrVSVrte2_3rbQjqciQhGSZi-_pmkoN76dzDmYNk_ucblFCJDjHVKsFNJtv4SqZoXw24RnyTMGIp3K70UxwmOK3EpW26OO6L7_mLIV3uIU3UCMn8cImXVWMURLlDLmGG-QoqUYVnnI0YTiaNTtuYAhSifrAdp4xGw-bHXOv1KGUMVXrWiOqNOKjRl2vggvbwCup_0i8Jo8r8tGubixGS36OnChAq4sZrA0cnmnOMhj8I1FZ_SdjWJ20t1mjg9bLF1QvcJ0MJcKdeKR8BbTUa9MZS2k1BVWVoqDcwGDxN_KDeIGHJGYqtfOsL6__0I6PJk1xI6TxSIta5LTOaq4fZ5-grawl7YIZRkmXuv0oGD9hrmRtYFy5dZsM-VZwhIhy-IySLU2njJ4N4mZk6I1Vt4YPs0Fp2uK0wAvb3iQ7zV4KCV0wVZTm9SxIg6wky0ioZYkNUqAs6OFIdgfUnBi_CpyT0GwzKp_mZM73BrOh_EGI4giTolytSbikuTKncpOZlxkzupK0-B2V5gmhjETJNQld3-tXLCTckR8k7Actx-l6ntu_7Pm-7102yJaEnttrOb7n-q7vdB0_6Owb5Gcl67QC3w-6jtN3um4v6Lm9BsGMaSHv7B-tvur-F5vrGG8?type=png)](https://mermaid.live/edit#pako:eNp9U01z2jAQ_Ss7OuRSPowdsPGhDLFJIG3aDHQ60xgOqr2AEluikpyWEv57hWXaQjrVSVrte2_3rbQjqciQhGSZi-_pmkoN76dzDmYNk_ucblFCJDjHVKsFNJtv4SqZoXw24RnyTMGIp3K70UxwmOK3EpW26OO6L7_mLIV3uIU3UCMn8cImXVWMURLlDLmGG-QoqUYVnnI0YTiaNTtuYAhSifrAdp4xGw-bHXOv1KGUMVXrWiOqNOKjRl2vggvbwCup_0i8Jo8r8tGubixGS36OnChAq4sZrA0cnmnOMhj8I1FZ_SdjWJ20t1mjg9bLF1QvcJ0MJcKdeKR8BbTUa9MZS2k1BVWVoqDcwGDxN_KDeIGHJGYqtfOsL6__0I6PJk1xI6TxSIta5LTOaq4fZ5-grawl7YIZRkmXuv0oGD9hrmRtYFy5dZsM-VZwhIhy-IySLU2njJ4N4mZk6I1Vt4YPs0Fp2uK0wAvb3iQ7zV4KCV0wVZTm9SxIg6wky0ioZYkNUqAs6OFIdgfUnBi_CpyT0GwzKp_mZM73BrOh_EGI4giTolytSbikuTKncpOZlxkzupK0-B2V5gmhjETJNQld3-tXLCTckR8k7Actx-l6ntu_7Pm-7102yJaEnttrOb7n-q7vdB0_6Owb5Gcl67QC3w-6jtN3um4v6Lm9BsGMaSHv7B-tvur-F5vrGG8)