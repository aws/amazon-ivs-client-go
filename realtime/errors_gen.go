package realtime

import "fmt"

// ErrorType is the category of an error. Many error codes share one type, so
// match on the type to handle a group of errors together.
type ErrorType int32

const (
	ErrorTypeGeneric          ErrorType = 1  // A failure with no more specific category.
	ErrorTypeNotSupported     ErrorType = 2  // The operation is not supported.
	ErrorTypeNoSource         ErrorType = 3  // A required source is missing.
	ErrorTypeInvalidData      ErrorType = 4  // Data was malformed or could not be parsed.
	ErrorTypeInvalidState     ErrorType = 5  // The operation is not allowed in the current state.
	ErrorTypeInvalidParameter ErrorType = 6  // An argument was invalid.
	ErrorTypeTimeout          ErrorType = 7  // The operation timed out.
	ErrorTypeNetwork          ErrorType = 8  // A network connection failed.
	ErrorTypeNetworkIO        ErrorType = 9  // Sending or receiving over the network failed.
	ErrorTypeAuthorization    ErrorType = 10 // The request was not authorized.
	ErrorTypeNotAvailable     ErrorType = 11 // A required resource is not available.
)

func (v ErrorType) String() string {
	switch v {
	case ErrorTypeGeneric:
		return "Generic"
	case ErrorTypeNotSupported:
		return "NotSupported"
	case ErrorTypeNoSource:
		return "NoSource"
	case ErrorTypeInvalidData:
		return "InvalidData"
	case ErrorTypeInvalidState:
		return "InvalidState"
	case ErrorTypeInvalidParameter:
		return "InvalidParameter"
	case ErrorTypeTimeout:
		return "Timeout"
	case ErrorTypeNetwork:
		return "Network"
	case ErrorTypeNetworkIO:
		return "NetworkIO"
	case ErrorTypeAuthorization:
		return "Authorization"
	case ErrorTypeNotAvailable:
		return "NotAvailable"
	}
	return fmt.Sprintf("ErrorType(%d)", int32(v))
}

// ErrorCode is the result and error codes for the IVS Linux Broadcast SDK. OK
// (0) means success. Existing values never change meaning. Undocumented codes
// may be observed at runtime (SDK-internal errors); treat an unknown code as an
// internal error.
type ErrorCode int32

const (
	ErrorCodeOK                                                         ErrorCode = 0     // Success.
	ErrorCodeSocketCreationFailed                                       ErrorCode = 201   // Network connection failed during initialization due to system resource constraints. Verify network availability and retry the stream.
	ErrorCodeSocketConfigurationFailed                                  ErrorCode = 202   // Network connection setup failed due to configuration issues. Check network connectivity and retry the stream.
	ErrorCodeSocketConnectFailed                                        ErrorCode = 203   // Unable to establish network connection to the streaming server. Verify network connectivity and that the stream URL is accessible.
	ErrorCodeSocketSendFailed                                           ErrorCode = 204   // Failed to send data over the network connection. Check network stability and restart the stream.
	ErrorCodeSocketRecvFailed                                           ErrorCode = 205   // Failed to receive data from the streaming server. Check network stability and restart the stream.
	ErrorCodeSocketRecvFailedEof                                        ErrorCode = 206   // The streaming server closed the connection unexpectedly. Verify no other users are streaming to this endpoint and restart the stream.
	ErrorCodeSocketRemoteHangup                                         ErrorCode = 207   // Network connection closed unexpectedly during stream playback. Check network stability and restart the stream.
	ErrorCodeSocketRemoteError                                          ErrorCode = 208   // Network error occurred during stream playback. Verify network connectivity and restart the stream.
	ErrorCodeSocketGetRttFailedNotSupportedOnThisSystem                 ErrorCode = 212   // Network latency measurement is unavailable on this device. Streaming will continue without latency metrics.
	ErrorCodeSocketSendLowWaterFailed                                   ErrorCode = 213   // Network send buffer configuration failed. Check network conditions and retry the stream.
	ErrorCodeCouldNotResolveAddress                                     ErrorCode = 214   // Unable to resolve the streaming server hostname. Check network connectivity and verify the stream URL is correct.
	ErrorCodeSocketGetRttFailedNotConnected                             ErrorCode = 215   // Network latency measurement was requested while no network connection was established. Latency metrics are unavailable until the connection is restored.
	ErrorCodeGetAverageOutgoingRateFailed                               ErrorCode = 300   // Network metrics calculation failed due to invalid configuration. This is an internal error that should not affect streaming, but calculated transmission statistics may be impacted.
	ErrorCodeGetEstimatedCongestionLevelFailed                          ErrorCode = 301   // Network congestion measurement failed due to invalid configuration. This is an internal error that should not affect streaming, but calculated transmission statistics may be impacted.
	ErrorCodeTlsSendFailedBadWriteRetry                                 ErrorCode = 400   // Secure connection data transmission failed due to protocol violation. Restart the stream to re-establish the connection.
	ErrorCodeTlsSendFailedAttemptedToSendWhileHandshaking               ErrorCode = 401   // Data transmission attempted before secure connection was established. Restart the stream.
	ErrorCodeTlsSendFailedAttemptedToSendWhileDisconnected              ErrorCode = 402   // Data transmission attempted after connection was closed. Restart the stream to re-establish the connection.
	ErrorCodeTlsRecvFailedAttemptedToRecvWhileHandshaking               ErrorCode = 403   // Data reception attempted before secure connection was established. Restart the stream.
	ErrorCodeTlsRecvFailedAttemptedToRecvWhileDisconnected              ErrorCode = 404   // Data reception attempted after connection was closed. Restart the stream to re-establish the connection.
	ErrorCodeTlsRecvFailedAttemptedToRecvAfterReceivingShutdownFromPeer ErrorCode = 405   // The streaming server closed the connection. Verify no other users are streaming to this endpoint and restart the stream.
	ErrorCodeTlsRecvFailedAttemptedToRecvWhileShuttingDown              ErrorCode = 406   // Data reception attempted while connection is closing. This is expected during stream shutdown.
	ErrorCodeTlsHandshakeFailedCouldNotCreateSslInstance                ErrorCode = 407   // Secure connection initialization failed. This is usually an allocation failure, the device may not be able to allocate more memory to the process.
	ErrorCodeTlsSendFailedOpenSslError                                  ErrorCode = 408   // Failed to send data over secure connection. Check network stability and restart the stream. This could also be caused by a firewall restricting access to the provided domain.
	ErrorCodeTlsRecvFailedOpenSslError                                  ErrorCode = 409   // Failed to receive data over secure connection. Check network stability and restart the stream.
	ErrorCodeTlsHandshakeFailedOpenSslError                             ErrorCode = 410   // Secure connection handshake failed. Verify network connectivity and restart the stream. This could also be caused by a firewall restricting access to the provided domain.
	ErrorCodeTlsFailedToGetX509Param                                    ErrorCode = 411   // Certificate verification setup failed during secure connection. Restart the stream to retry the connection. This could also be caused by a firewall restricting access to the provided domain.
	ErrorCodeTlsSetSniFailed                                            ErrorCode = 412   // Server hostname configuration failed during secure connection setup. Verify the stream URL and restart the stream. This could also be caused by a firewall restricting access to the provided domain.
	ErrorCodeTlsSetExpectedHostFailed                                   ErrorCode = 413   // Certificate hostname verification setup failed. Verify the stream URL is correct and restart the stream. This could also be caused by a firewall restricting access to the provided domain.
	ErrorCodeTlsHandshakeTimeout                                        ErrorCode = 414   // The secure connection handshake timed out. The server accepted the connection but did not complete the TLS negotiation. Verify that the endpoint expects a secure connection and check network connectivity.
	ErrorCodeEventLoopFailedToCreateTimer                               ErrorCode = 500   // Internal timing mechanism failed. This may prevent efficient DNS resolution.
	ErrorCodeWebSocketHandshakeInvalidUpgradeResponse                   ErrorCode = 600   // Server did not accept WebSocket connection request.
	ErrorCodeWebSocketHandshakeInvalidAcceptKey                         ErrorCode = 601   // WebSocket handshake validation failed.
	ErrorCodeWebSocketHandshakeProtocolNegotiationFailed                ErrorCode = 602   // Server selected a WebSocket protocol not supported by the client. Verify the events URL is correct and restart the session.
	ErrorCodeWebSocketHandshakeTimeout                                  ErrorCode = 603   // WebSocket connection took too long to establish. Check network connectivity and retry the connection.
	ErrorCodeWebSocketHandshakeInvalidStatusCode                        ErrorCode = 604   // Server rejected WebSocket connection with unexpected status code. Verify the events URL supports WebSocket connections.
	ErrorCodeWebSocketHandshakeInvalidHeaders                           ErrorCode = 605   // Server response headers do not meet WebSocket protocol requirements. Verify the events URL is correct and supports WebSocket.
	ErrorCodeWebSocketFrameInvalidOpcode                                ErrorCode = 620   // Received malformed data from the streaming server. Restart the session to re-establish the connection.
	ErrorCodeWebSocketFrameInvalidLength                                ErrorCode = 621   // Received data with invalid length from the streaming server. Restart the connection to re-establish the connection.
	ErrorCodeWebSocketFrameInvalidMask                                  ErrorCode = 622   // Server violated WebSocket protocol by sending masked data. Restart the session to re-establish the connection.
	ErrorCodeWebSocketFrameControlFrameTooLarge                         ErrorCode = 623   // Server sent control message exceeding protocol limit. Restart the session to re-establish the connection.
	ErrorCodeWebSocketFrameInvalidUtf8                                  ErrorCode = 624   // Received text message with invalid character encoding from the server. Restart the session to re-establish the connection.
	ErrorCodeWebSocketFrameReservedBitsSet                              ErrorCode = 625   // Server violated WebSocket protocol by setting reserved bits. Restart the session to re-establish the connection.
	ErrorCodeWebSocketFrameIncompleteData                               ErrorCode = 626   // Received partial data from the streaming server. This is typically handled automatically by waiting for more data and should not impact the session.
	ErrorCodeWebSocketFrameControlFrameFragmented                       ErrorCode = 627   // Server violated WebSocket protocol by fragmenting a control message. Restart the session to re-establish the connection.
	ErrorCodeWebSocketProtocolMessageTooLarge                           ErrorCode = 641   // Received message exceeds maximum allowed size. Restart the session to re-establish the connection.
	ErrorCodeWebSocketProtocolConnectionNotEstablished                  ErrorCode = 642   // Attempted to send data before WebSocket connection was ready. This is an internal SDK and the session should be retried.
	ErrorCodeWebSocketConnectionTimeout                                 ErrorCode = 660   // WebSocket connection took too long to establish. Check network connectivity and restart the session.
	ErrorCodeWebSocketConnectionLost                                    ErrorCode = 661   // The streaming server closed the WebSocket connection. Check network connectivity and restart the session.
	ErrorCodeWebSocketConnectionCloseTimeout                            ErrorCode = 662   // Server did not acknowledge connection close request in time. The connection will be terminated forcefully.
	ErrorCodeWebSocketConnectionPongTimeout                             ErrorCode = 663   // The client sent a WebSocket PING on the idle WebSocket but no PONG was received within the timeout period. It is likely the socket has disconnected but hasn't yet reached the TCP idle timeout. The WebSocket will automatically attempt to reconnect.
	ErrorCodeStageInvalidToken                                          ErrorCode = 1000  // Stage token is invalid or malformed. Generate a new token and create a new stage session.
	ErrorCodeStageExpiredToken                                          ErrorCode = 1001  // Stage token has expired. Generate a new token and create a new stage session.
	ErrorCodeNoSession                                                  ErrorCode = 1002  // Cannot perform operation because no active stage session exists. Join a stage before attempting this operation.
	ErrorCodeAlreadyJoiningStage                                        ErrorCode = 1005  // Cannot join stage because a join operation is already in progress or completed. Wait for the current operation to finish.
	ErrorCodeStageNotReady                                              ErrorCode = 1011  // The Stage session is not ready due to a failure is the setup process. Check for errors during initialization, and recreate the local Stage instance.
	ErrorCodeSelfSubscribeNotSupport                                    ErrorCode = 1012  // Cannot subscribe to your own participant stream. Subscribe to other participants instead. This should not be possible using the StageStrategy.
	ErrorCodeNoParticipantToSubscribeTo                                 ErrorCode = 1013  // Participant not found in the stage. This is likely because the remote participant unpublished at the same time this SDK attempted to subscribe.
	ErrorCodePublishFailed                                              ErrorCode = 1015  // Failed to publish audio/video to the stage after multiple retry attempts. Check network connectivity and retry.
	ErrorCodeSubscribeFailed                                            ErrorCode = 1016  // Failed to subscribe to participant stream after multiple retry attempts. Check network connectivity and retry.
	ErrorCodeOfferCreationFailPublish                                   ErrorCode = 1020  // Failed to create WebRTC offer for publishing. Retry the publish operation.
	ErrorCodeOfferCreationFailSubscribe                                 ErrorCode = 1021  // Failed to create WebRTC offer for subscribing. Retry the subscribe operation.
	ErrorCodePubSubAborted                                              ErrorCode = 1023  // Publish or subscribe operation was cancelled. This is expected when leaving or unsubscribing in the middle of an operation.
	ErrorCodeStageAtCapacity                                            ErrorCode = 1024  // Stage has reached maximum participant capacity. Wait for a participant to leave or create a new stage.
	ErrorCodeInternalServerError                                        ErrorCode = 1025  // An unexpected server error occurred. Retry the operation or contact support if the issue persists.
	ErrorCodeStageTokenRejected                                         ErrorCode = 1026  // Stage token lacks required permissions or has been revoked. Generate a new token with appropriate permissions.
	ErrorCodeStageAuthenticationError                                   ErrorCode = 1027  // Stage authentication failed during join operation. The stage may have been deleted or is no longer available.
	ErrorCodeInsufficientInputResolution                                ErrorCode = 1050  // Input video resolution does not meet stage requirements. Increase camera resolution or adjust stage configuration.
	ErrorCodeTokenExchangeInvalidToken                                  ErrorCode = 1051  // The token provided to the exchange token API was invalid and couldn't be parsed. Create a new token and try again.
	ErrorCodeTokenExchangeIncompatibleToken                             ErrorCode = 1052  // The token provided to the exchange token API was not compatible with the existing Stage. This is likely because properties that must remain constant (jti, topic, etc) had changed.
	ErrorCodeConnectionTokenInvalid                                     ErrorCode = 1053  // The connection token passed to the IVS RealTimeConnection is invalid or malformed. Generate a new connection token and create a new connection.
	ErrorCodeConnectionTokenExpired                                     ErrorCode = 1054  // The connection token passed to the IVS RealTimeConnection has expired. Generate a new connection token and create a new connection.
	ErrorCodeMultihostGroupStatePipelineMissingMultihostEventBus        ErrorCode = 1100  // Internal stage event system is not initialized. Restart the stage session.
	ErrorCodeMultihostGroupStatePipelineMissingEventGenerator           ErrorCode = 1101  // Internal event generation system is not initialized. Restart the stage session.
	ErrorCodeResourceDeleted                                            ErrorCode = 1206  // Stage or participant resource no longer exists on the server. The stage may have been deleted or the session expired.
	ErrorCodeRemoteDisconnectReasonUnknown                              ErrorCode = 1208  // Disconnected from stage by the server for an unspecified reason. Rejoin the stage with a new token if needed.
	ErrorCodeRemoteDisconnectStageDeleted                               ErrorCode = 1209  // Disconnected because the stage was deleted.
	ErrorCodeRemoteDisconnectTokenReused                                ErrorCode = 1210  // Disconnected because another participant joined with the same token. Use unique tokens for each participant.
	ErrorCodeRemoteDisconnectParticipantDisconnected                    ErrorCode = 1211  // Disconnected from stage by server request. Rejoin the stage with a new token if needed.
	ErrorCodeQuotaExceeded                                              ErrorCode = 1212  // A service-side quota or limit was exceeded. This is a general-purpose error covering multiple quotas; read the error message for the specific limit that was hit and how to resolve it.
	ErrorCodeStagePresenceInconsistentDelta                             ErrorCode = 1215  // The presence roster became unavailable because the SDK received a presence delta it could not reconcile with the roster it was tracking (for example a participant leaving or being updated when the SDK never saw them join, or joining when they were already present). Because the delta stream can no longer be trusted, the SDK stops applying it and marks the roster unavailable. This is a non-fatal error and does not affect media: publishers still reach the client through the participant/subscribe callbacks. The app is also notified via the consolidated roster callback (Android StageRenderer.onRosterChanged / iOS stage:rosterDidChange:) with an UnavailableError availability; any locally held roster should be cleared and roster UI hidden. The roster recovers if a later valid snapshot arrives.
	ErrorCodeStagePresenceResyncExhausted                               ErrorCode = 1216  // The presence roster became unavailable because the SDK repeatedly could not obtain a usable recovery snapshot (its recovery buffer kept overflowing while awaiting one) and gave up after several attempts. This is a non-fatal error and does not affect media: publishers still reach the client through the participant/subscribe callbacks. The app is also notified via the consolidated roster callback (Android StageRenderer.onRosterChanged / iOS stage:rosterDidChange:) with an UnavailableError availability; any locally held roster should be cleared and roster UI hidden. The roster recovers if a later valid snapshot arrives.
	ErrorCodeStagePresenceSnapshotRequestTimeout                        ErrorCode = 1218  // The presence roster became unavailable because a recovery snapshot request timed out with no server response, and the SDK gave up. This is a non-fatal error and does not affect media: publishers still reach the client through the participant/subscribe callbacks. The app is also notified via the consolidated roster callback (Android StageRenderer.onRosterChanged / iOS stage:rosterDidChange:) with an UnavailableError availability; any locally held roster should be cleared and roster UI hidden. The roster recovers if a later valid snapshot arrives.
	ErrorCodeStagePresenceInconsistentSnapshot                          ErrorCode = 1219  // The presence roster became unavailable because the SDK received a presence snapshot it could not treat as authoritative — it listed the same participant more than once, which should never happen. Because the snapshot cannot be trusted, the SDK discards it and marks the roster unavailable. This is a non-fatal error and does not affect media: publishers still reach the client through the participant/subscribe callbacks. The app is also notified via the consolidated roster callback (Android StageRenderer.onRosterChanged / iOS stage:rosterDidChange:) with an UnavailableError availability; any locally held roster should be cleared and roster UI hidden. The roster recovers if a later valid snapshot arrives.
	ErrorCodeTokenExchangeInvalidResponse                               ErrorCode = 1230  // Invalid response from the server for a token exchange request. Reconnect to the Stage using the new token, or create a new token and connect again.
	ErrorCodeTokenExchangeFailedResponse                                ErrorCode = 1231  // The token exchange request was rejected. The error comment will contain more details. A new token must be created and used to join a new SDK Stage instance.
	ErrorCodeTokenExchangeTimeout                                       ErrorCode = 1232  // The token exchange request timed out. Because a token exchange puts the SDK into an in between state of 2 tokens while waiting for confirmation from the backend, a timeout is enforced to prevent undefined behaviors.
	ErrorCodeSignalingRetryBudgetExhausted                              ErrorCode = 1240  // A request to the IVS service was rate limited repeatedly until timeout exhaustion.
	ErrorCodeEventsRetriesExhausted                                     ErrorCode = 1300  // Failed to connect to stage after multiple retry attempts. Check network connectivity and retry joining the stage.
	ErrorCodeRealTimeConnectionReleased                                 ErrorCode = 1311  // Cannot perform this operation because the RealTimeConnection provided to the Stage has been released. Keep the RealTimeConnection alive for as long as any Stage that was created with it may be used.
	ErrorCodePeerConnectionNetworkError                                 ErrorCode = 1400  // Media connection lost due to network error. Check network connectivity and rejoin the stage.
	ErrorCodePeerConnectionSetRemoteDescriptionFail                     ErrorCode = 1401  // Failed to establish media connection with participant. Retry the publish or subscribe operation.
	ErrorCodePeerConnectionParseRemoteDescriptionError                  ErrorCode = 1404  // Received invalid connection data from participant.
	ErrorCodePeerConnectionUnexpectedRemoteClose                        ErrorCode = 1407  // Participant disconnected unexpectedly. The participant may have left the stage or lost network connectivity.
	ErrorCodePeerConnectionSetParametersError                           ErrorCode = 1410  // Failed to configure video settings for the connection. Check system resources, this may have been a memory allocation failure.
	ErrorCodePeerConnectionRequestKeyframeNoSender                      ErrorCode = 1411  // Cannot request keyframe because video sender is not ready. Ensure video is being published.
	ErrorCodePeerConnectionAddTransceiversFailed                        ErrorCode = 1412  // Failed to configure media transceivers for the connection. This is a libWebRTC failure and may indicate the system is low on resources.
	ErrorCodePeerConnectionVideoSenderNotReady                          ErrorCode = 1413  // Video sender is not ready for configuration. Ensure video is being published before adjusting settings.
	ErrorCodePeerConnectionVideoRtpParametersNotReady                   ErrorCode = 1414  // Video encoding parameters are not ready. Wait for the connection to fully establish before adjusting video settings.
	ErrorCodePeerConnectionSetVideoRtpParametersFailed                  ErrorCode = 1415  // Failed to update video encoding settings. This is a libWebRTC failure and may indicate the system is low on resources, or the provided configuration is not supported.
	ErrorCodePeerConnectionMethodNotSupportedError                      ErrorCode = 1420  // Connection operation failed due to invalid configuration. The error message will contain more detailed information.
	ErrorCodePeerConnectionInitFail                                     ErrorCode = 1421  // Failed to initialize media connection. This is a libWebRTC failure and may indicate the system is low on resources. Restat the session to try again.
	ErrorCodePeerConnectionNoMatchingMediaType                          ErrorCode = 1423  // Operation failed because the requested media type is not available. Ensure audio or video is being published or subscribed before requesting RTC stats.
	ErrorCodePeerConnectionTimedOutConnecting                           ErrorCode = 1424  // Media connection is taking too long to establish. This operation will be automatically retried, but this may indicate a network connectivity issue.
	ErrorCodePeerConnectionOnFailure                                    ErrorCode = 1425  // WebRTC connection operation failed. This operation will be automatically retried, but this may indicate a network connectivity issue.
	ErrorCodePeerConnectionIceConnectionFailed                          ErrorCode = 1426  // Media connection failed to establish network path. This operation will be automatically retried, but this may indicate a network connectivity issue.
	ErrorCodePeerConnectionCallbackSenderNotReady                       ErrorCode = 1427  // Cannot retrieve video sender configuration because sender is not initialized. Ensure video is being published.
	ErrorCodePeerConnectionCallbackReceiverNotReady                     ErrorCode = 1428  // Cannot retrieve video receiver configuration because receiver is not initialized. Ensure video is being subscribed.
	ErrorCodeNoValidPeerConnection                                      ErrorCode = 1430  // Cannot perform operation because media connection is not available. Ensure publish or subscribe is active.
	ErrorCodePeerConnectionCreateDataChannelFailed                      ErrorCode = 1431  // The stage connection could not create the channel that carries messages, so messages cannot be sent or received on it.
	ErrorCodePeerConnectionDataChannelClosed                            ErrorCode = 1432  // A message was handed to the transport, but the channel closed before reporting what became of it, so it may or may not have been delivered.
	ErrorCodePeerConnectionDataChannelSendEncodeFailed                  ErrorCode = 1433  // A message could not be sent because its contents could not be encoded. A required field was empty; the error message names the field at fault.
	ErrorCodePeerConnectionDataChannelNotOpen                           ErrorCode = 1434  // A message was not sent because the channel that carries it was not open, usually because the local participant was not publishing, or stopped before the message went out. The error message carries the id of the message affected.
	ErrorCodePeerConnectionDataChannelSendRejected                      ErrorCode = 1435  // A message was refused before being sent, because it was larger than the maximum, the send rate was above the maximum, or the transport would not accept it. Read the error message for which limit was hit, and send smaller or fewer messages.
	ErrorCodePeerConnectionDataChannelReceiveDecodeFailed               ErrorCode = 1436  // A message arrived but could not be read, so it was dropped.
	ErrorCodePeerConnectionCueSendCancelled                             ErrorCode = 1438  // A cue was handed to the transport, but the application released the stream it was sent on before its outcome arrived, so it may or may not have been delivered.
	ErrorCodeNoPublishAudioPacketsSent                                  ErrorCode = 1450  // No audio packets are being sent through the publish connection. Check that audio input is active and not muted.
	ErrorCodeNoPublishVideoPacketsSent                                  ErrorCode = 1451  // No video packets are being sent through the publish connection. Check that video samples are being submitted.
	ErrorCodeNoSubscribeAudioPacketsReceived                            ErrorCode = 1452  // No audio packets are being received from the subscribed participant. The participant may have stopped publishing audio or network connectivity is lost.
	ErrorCodeNoSubscribeVideoPacketsReceived                            ErrorCode = 1453  // No video packets are being received from the subscribed participant. The participant may have stopped publishing video without muting or network connectivity is lost.
	ErrorCodeSetLayerRequestFailed                                      ErrorCode = 1510  // Failed to set video quality layer for participant stream. This operation will be automatically retried, but this may indicate a network connectivity issue.
	ErrorCodeSetLayerRequestThrottled                                   ErrorCode = 1511  // Video quality layer selection is being rate limited. This operation will be automatically retried after a delay.
	ErrorCodeGetLayerRequestFailed                                      ErrorCode = 1512  // Failed to retrieve available video quality layers for participant stream. This operation will be automatically retried, but this may indicate a network connectivity issue.
	ErrorCodeGetLayerRequestThrottled                                   ErrorCode = 1513  // Video quality layer query is being rate limited. This operation will be automatically retried after a delay.
	ErrorCodeStageSetAudioCallbackWhileJoined                           ErrorCode = 10000 // The Stage audio callback was set while the Stage is joined. Set it before joining or after leaving.
	ErrorCodeStageStreamNotAudio                                        ErrorCode = 10001 // A PCM callback was set on a stream that is not an audio stream. Set PCM callbacks only on audio streams.
	ErrorCodeStageStreamEnded                                           ErrorCode = 10002 // A PCM callback was set on a stream that has ended. Register on the stream delivered by the next streams-added event.
	ErrorCodeSubtitleStreamNotPublishing                                ErrorCode = 10003 // A subtitle cue was sent on a stream the local participant is not publishing. Return the stream from the streams-to-publish strategy callback first.
	ErrorCodeAudioSourceBufferSizeMismatch                              ErrorCode = 10004 // The pushed audio buffer size does not match the frame count, channel count, and sample format of the audio source.
	ErrorCodeChannelTranscriptionNoParticipant                          ErrorCode = 10006 // A transcription was written while no participant has subscribed audio, so it cannot be attributed.
	ErrorCodeChannelTranscriptionInvalidTimeRange                       ErrorCode = 10007 // A transcription was written with an end time that is not after its start time.
	ErrorCodeChannelTranscriptionStale                                  ErrorCode = 10008 // A transcription was written with a start time that lags the live audio position by more than the maximum transcription age.
	ErrorCodeChannelTranscriptionCaptureTimeUnavailable                 ErrorCode = 10009 // No absolute capture time was received on the subscribed audio stream, so transcriptions cannot be delivered. The publisher must send the abs-capture-time RTP header extension.
	ErrorCodeChannelTranscriberNotAttached                              ErrorCode = 10010 // A transcription was written before the transcriber's strategy was passed to a Stage.
	ErrorCodeInvalidArgument                                            ErrorCode = 10800 // A required argument was null or out of range. Check the arguments passed to the function that returned this code.
)

func (v ErrorCode) String() string {
	switch v {
	case ErrorCodeOK:
		return "OK"
	case ErrorCodeSocketCreationFailed:
		return "SocketCreationFailed"
	case ErrorCodeSocketConfigurationFailed:
		return "SocketConfigurationFailed"
	case ErrorCodeSocketConnectFailed:
		return "SocketConnectFailed"
	case ErrorCodeSocketSendFailed:
		return "SocketSendFailed"
	case ErrorCodeSocketRecvFailed:
		return "SocketRecvFailed"
	case ErrorCodeSocketRecvFailedEof:
		return "SocketRecvFailedEof"
	case ErrorCodeSocketRemoteHangup:
		return "SocketRemoteHangup"
	case ErrorCodeSocketRemoteError:
		return "SocketRemoteError"
	case ErrorCodeSocketGetRttFailedNotSupportedOnThisSystem:
		return "SocketGetRttFailedNotSupportedOnThisSystem"
	case ErrorCodeSocketSendLowWaterFailed:
		return "SocketSendLowWaterFailed"
	case ErrorCodeCouldNotResolveAddress:
		return "CouldNotResolveAddress"
	case ErrorCodeSocketGetRttFailedNotConnected:
		return "SocketGetRttFailedNotConnected"
	case ErrorCodeGetAverageOutgoingRateFailed:
		return "GetAverageOutgoingRateFailed"
	case ErrorCodeGetEstimatedCongestionLevelFailed:
		return "GetEstimatedCongestionLevelFailed"
	case ErrorCodeTlsSendFailedBadWriteRetry:
		return "TlsSendFailedBadWriteRetry"
	case ErrorCodeTlsSendFailedAttemptedToSendWhileHandshaking:
		return "TlsSendFailedAttemptedToSendWhileHandshaking"
	case ErrorCodeTlsSendFailedAttemptedToSendWhileDisconnected:
		return "TlsSendFailedAttemptedToSendWhileDisconnected"
	case ErrorCodeTlsRecvFailedAttemptedToRecvWhileHandshaking:
		return "TlsRecvFailedAttemptedToRecvWhileHandshaking"
	case ErrorCodeTlsRecvFailedAttemptedToRecvWhileDisconnected:
		return "TlsRecvFailedAttemptedToRecvWhileDisconnected"
	case ErrorCodeTlsRecvFailedAttemptedToRecvAfterReceivingShutdownFromPeer:
		return "TlsRecvFailedAttemptedToRecvAfterReceivingShutdownFromPeer"
	case ErrorCodeTlsRecvFailedAttemptedToRecvWhileShuttingDown:
		return "TlsRecvFailedAttemptedToRecvWhileShuttingDown"
	case ErrorCodeTlsHandshakeFailedCouldNotCreateSslInstance:
		return "TlsHandshakeFailedCouldNotCreateSslInstance"
	case ErrorCodeTlsSendFailedOpenSslError:
		return "TlsSendFailedOpenSslError"
	case ErrorCodeTlsRecvFailedOpenSslError:
		return "TlsRecvFailedOpenSslError"
	case ErrorCodeTlsHandshakeFailedOpenSslError:
		return "TlsHandshakeFailedOpenSslError"
	case ErrorCodeTlsFailedToGetX509Param:
		return "TlsFailedToGetX509Param"
	case ErrorCodeTlsSetSniFailed:
		return "TlsSetSniFailed"
	case ErrorCodeTlsSetExpectedHostFailed:
		return "TlsSetExpectedHostFailed"
	case ErrorCodeTlsHandshakeTimeout:
		return "TlsHandshakeTimeout"
	case ErrorCodeEventLoopFailedToCreateTimer:
		return "EventLoopFailedToCreateTimer"
	case ErrorCodeWebSocketHandshakeInvalidUpgradeResponse:
		return "WebSocketHandshakeInvalidUpgradeResponse"
	case ErrorCodeWebSocketHandshakeInvalidAcceptKey:
		return "WebSocketHandshakeInvalidAcceptKey"
	case ErrorCodeWebSocketHandshakeProtocolNegotiationFailed:
		return "WebSocketHandshakeProtocolNegotiationFailed"
	case ErrorCodeWebSocketHandshakeTimeout:
		return "WebSocketHandshakeTimeout"
	case ErrorCodeWebSocketHandshakeInvalidStatusCode:
		return "WebSocketHandshakeInvalidStatusCode"
	case ErrorCodeWebSocketHandshakeInvalidHeaders:
		return "WebSocketHandshakeInvalidHeaders"
	case ErrorCodeWebSocketFrameInvalidOpcode:
		return "WebSocketFrameInvalidOpcode"
	case ErrorCodeWebSocketFrameInvalidLength:
		return "WebSocketFrameInvalidLength"
	case ErrorCodeWebSocketFrameInvalidMask:
		return "WebSocketFrameInvalidMask"
	case ErrorCodeWebSocketFrameControlFrameTooLarge:
		return "WebSocketFrameControlFrameTooLarge"
	case ErrorCodeWebSocketFrameInvalidUtf8:
		return "WebSocketFrameInvalidUtf8"
	case ErrorCodeWebSocketFrameReservedBitsSet:
		return "WebSocketFrameReservedBitsSet"
	case ErrorCodeWebSocketFrameIncompleteData:
		return "WebSocketFrameIncompleteData"
	case ErrorCodeWebSocketFrameControlFrameFragmented:
		return "WebSocketFrameControlFrameFragmented"
	case ErrorCodeWebSocketProtocolMessageTooLarge:
		return "WebSocketProtocolMessageTooLarge"
	case ErrorCodeWebSocketProtocolConnectionNotEstablished:
		return "WebSocketProtocolConnectionNotEstablished"
	case ErrorCodeWebSocketConnectionTimeout:
		return "WebSocketConnectionTimeout"
	case ErrorCodeWebSocketConnectionLost:
		return "WebSocketConnectionLost"
	case ErrorCodeWebSocketConnectionCloseTimeout:
		return "WebSocketConnectionCloseTimeout"
	case ErrorCodeWebSocketConnectionPongTimeout:
		return "WebSocketConnectionPongTimeout"
	case ErrorCodeStageInvalidToken:
		return "StageInvalidToken"
	case ErrorCodeStageExpiredToken:
		return "StageExpiredToken"
	case ErrorCodeNoSession:
		return "NoSession"
	case ErrorCodeAlreadyJoiningStage:
		return "AlreadyJoiningStage"
	case ErrorCodeStageNotReady:
		return "StageNotReady"
	case ErrorCodeSelfSubscribeNotSupport:
		return "SelfSubscribeNotSupport"
	case ErrorCodeNoParticipantToSubscribeTo:
		return "NoParticipantToSubscribeTo"
	case ErrorCodePublishFailed:
		return "PublishFailed"
	case ErrorCodeSubscribeFailed:
		return "SubscribeFailed"
	case ErrorCodeOfferCreationFailPublish:
		return "OfferCreationFailPublish"
	case ErrorCodeOfferCreationFailSubscribe:
		return "OfferCreationFailSubscribe"
	case ErrorCodePubSubAborted:
		return "PubSubAborted"
	case ErrorCodeStageAtCapacity:
		return "StageAtCapacity"
	case ErrorCodeInternalServerError:
		return "InternalServerError"
	case ErrorCodeStageTokenRejected:
		return "StageTokenRejected"
	case ErrorCodeStageAuthenticationError:
		return "StageAuthenticationError"
	case ErrorCodeInsufficientInputResolution:
		return "InsufficientInputResolution"
	case ErrorCodeTokenExchangeInvalidToken:
		return "TokenExchangeInvalidToken"
	case ErrorCodeTokenExchangeIncompatibleToken:
		return "TokenExchangeIncompatibleToken"
	case ErrorCodeConnectionTokenInvalid:
		return "ConnectionTokenInvalid"
	case ErrorCodeConnectionTokenExpired:
		return "ConnectionTokenExpired"
	case ErrorCodeMultihostGroupStatePipelineMissingMultihostEventBus:
		return "MultihostGroupStatePipelineMissingMultihostEventBus"
	case ErrorCodeMultihostGroupStatePipelineMissingEventGenerator:
		return "MultihostGroupStatePipelineMissingEventGenerator"
	case ErrorCodeResourceDeleted:
		return "ResourceDeleted"
	case ErrorCodeRemoteDisconnectReasonUnknown:
		return "RemoteDisconnectReasonUnknown"
	case ErrorCodeRemoteDisconnectStageDeleted:
		return "RemoteDisconnectStageDeleted"
	case ErrorCodeRemoteDisconnectTokenReused:
		return "RemoteDisconnectTokenReused"
	case ErrorCodeRemoteDisconnectParticipantDisconnected:
		return "RemoteDisconnectParticipantDisconnected"
	case ErrorCodeQuotaExceeded:
		return "QuotaExceeded"
	case ErrorCodeStagePresenceInconsistentDelta:
		return "StagePresenceInconsistentDelta"
	case ErrorCodeStagePresenceResyncExhausted:
		return "StagePresenceResyncExhausted"
	case ErrorCodeStagePresenceSnapshotRequestTimeout:
		return "StagePresenceSnapshotRequestTimeout"
	case ErrorCodeStagePresenceInconsistentSnapshot:
		return "StagePresenceInconsistentSnapshot"
	case ErrorCodeTokenExchangeInvalidResponse:
		return "TokenExchangeInvalidResponse"
	case ErrorCodeTokenExchangeFailedResponse:
		return "TokenExchangeFailedResponse"
	case ErrorCodeTokenExchangeTimeout:
		return "TokenExchangeTimeout"
	case ErrorCodeSignalingRetryBudgetExhausted:
		return "SignalingRetryBudgetExhausted"
	case ErrorCodeEventsRetriesExhausted:
		return "EventsRetriesExhausted"
	case ErrorCodeRealTimeConnectionReleased:
		return "RealTimeConnectionReleased"
	case ErrorCodePeerConnectionNetworkError:
		return "PeerConnectionNetworkError"
	case ErrorCodePeerConnectionSetRemoteDescriptionFail:
		return "PeerConnectionSetRemoteDescriptionFail"
	case ErrorCodePeerConnectionParseRemoteDescriptionError:
		return "PeerConnectionParseRemoteDescriptionError"
	case ErrorCodePeerConnectionUnexpectedRemoteClose:
		return "PeerConnectionUnexpectedRemoteClose"
	case ErrorCodePeerConnectionSetParametersError:
		return "PeerConnectionSetParametersError"
	case ErrorCodePeerConnectionRequestKeyframeNoSender:
		return "PeerConnectionRequestKeyframeNoSender"
	case ErrorCodePeerConnectionAddTransceiversFailed:
		return "PeerConnectionAddTransceiversFailed"
	case ErrorCodePeerConnectionVideoSenderNotReady:
		return "PeerConnectionVideoSenderNotReady"
	case ErrorCodePeerConnectionVideoRtpParametersNotReady:
		return "PeerConnectionVideoRtpParametersNotReady"
	case ErrorCodePeerConnectionSetVideoRtpParametersFailed:
		return "PeerConnectionSetVideoRtpParametersFailed"
	case ErrorCodePeerConnectionMethodNotSupportedError:
		return "PeerConnectionMethodNotSupportedError"
	case ErrorCodePeerConnectionInitFail:
		return "PeerConnectionInitFail"
	case ErrorCodePeerConnectionNoMatchingMediaType:
		return "PeerConnectionNoMatchingMediaType"
	case ErrorCodePeerConnectionTimedOutConnecting:
		return "PeerConnectionTimedOutConnecting"
	case ErrorCodePeerConnectionOnFailure:
		return "PeerConnectionOnFailure"
	case ErrorCodePeerConnectionIceConnectionFailed:
		return "PeerConnectionIceConnectionFailed"
	case ErrorCodePeerConnectionCallbackSenderNotReady:
		return "PeerConnectionCallbackSenderNotReady"
	case ErrorCodePeerConnectionCallbackReceiverNotReady:
		return "PeerConnectionCallbackReceiverNotReady"
	case ErrorCodeNoValidPeerConnection:
		return "NoValidPeerConnection"
	case ErrorCodePeerConnectionCreateDataChannelFailed:
		return "PeerConnectionCreateDataChannelFailed"
	case ErrorCodePeerConnectionDataChannelClosed:
		return "PeerConnectionDataChannelClosed"
	case ErrorCodePeerConnectionDataChannelSendEncodeFailed:
		return "PeerConnectionDataChannelSendEncodeFailed"
	case ErrorCodePeerConnectionDataChannelNotOpen:
		return "PeerConnectionDataChannelNotOpen"
	case ErrorCodePeerConnectionDataChannelSendRejected:
		return "PeerConnectionDataChannelSendRejected"
	case ErrorCodePeerConnectionDataChannelReceiveDecodeFailed:
		return "PeerConnectionDataChannelReceiveDecodeFailed"
	case ErrorCodePeerConnectionCueSendCancelled:
		return "PeerConnectionCueSendCancelled"
	case ErrorCodeNoPublishAudioPacketsSent:
		return "NoPublishAudioPacketsSent"
	case ErrorCodeNoPublishVideoPacketsSent:
		return "NoPublishVideoPacketsSent"
	case ErrorCodeNoSubscribeAudioPacketsReceived:
		return "NoSubscribeAudioPacketsReceived"
	case ErrorCodeNoSubscribeVideoPacketsReceived:
		return "NoSubscribeVideoPacketsReceived"
	case ErrorCodeSetLayerRequestFailed:
		return "SetLayerRequestFailed"
	case ErrorCodeSetLayerRequestThrottled:
		return "SetLayerRequestThrottled"
	case ErrorCodeGetLayerRequestFailed:
		return "GetLayerRequestFailed"
	case ErrorCodeGetLayerRequestThrottled:
		return "GetLayerRequestThrottled"
	case ErrorCodeStageSetAudioCallbackWhileJoined:
		return "StageSetAudioCallbackWhileJoined"
	case ErrorCodeStageStreamNotAudio:
		return "StageStreamNotAudio"
	case ErrorCodeStageStreamEnded:
		return "StageStreamEnded"
	case ErrorCodeSubtitleStreamNotPublishing:
		return "SubtitleStreamNotPublishing"
	case ErrorCodeAudioSourceBufferSizeMismatch:
		return "AudioSourceBufferSizeMismatch"
	case ErrorCodeChannelTranscriptionNoParticipant:
		return "ChannelTranscriptionNoParticipant"
	case ErrorCodeChannelTranscriptionInvalidTimeRange:
		return "ChannelTranscriptionInvalidTimeRange"
	case ErrorCodeChannelTranscriptionStale:
		return "ChannelTranscriptionStale"
	case ErrorCodeChannelTranscriptionCaptureTimeUnavailable:
		return "ChannelTranscriptionCaptureTimeUnavailable"
	case ErrorCodeChannelTranscriberNotAttached:
		return "ChannelTranscriberNotAttached"
	case ErrorCodeInvalidArgument:
		return "InvalidArgument"
	}
	return fmt.Sprintf("ErrorCode(%d)", int32(v))
}

// Error is an error delivered by a failed operation or asynchronous Stage
// event. A function's documented errors are not exhaustive. For any error,
// fatality says whether the Stage can still be used, and the code and message
// identify the failure when reporting it.
type Error struct {
	// Public or SDK-internal error code. Unknown codes represent internal errors.
	Code ErrorCode
	// Category of the error.
	Type ErrorType
	// Platform-specific code such as an HTTP status or errno; zero when
	// unavailable.
	PlatformCode int32
	// Whether the error terminated the Stage or otherwise requires termination
	// handling.
	IsFatal bool
	// Component that produced the error.
	Source string
	// Human-readable error detail.
	Message string
}
