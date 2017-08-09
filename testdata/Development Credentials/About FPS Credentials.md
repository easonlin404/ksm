# About FairPlay Streaming Credentials

* Every playback app that uses FairPlay Streaming (FPS) must find the media’s key server and establish communication with that server. When messages can be exchanged between the FPS aware client and the key server, the app must send the server an FPS-created SPC message.  This message contains a hash of the Application Certificate identifying your private key which your KSM should use to verify the identity of your client.

* While the Application Certificate does have an expiration date, FPS does not enforce the expiration date at the system level.  You should not enforce the expiration date in your FPS aware client or KSM.