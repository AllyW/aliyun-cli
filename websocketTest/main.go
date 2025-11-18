package websockettest

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aliyun/aliyun-cli/v3/cli"
	"github.com/aliyun/aliyun-cli/v3/config"
	"github.com/aliyun/aliyun-cli/v3/i18n"

	openapiClient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	dara "github.com/alibabacloud-go/tea/dara"
	"github.com/alibabacloud-go/tea/tea"
)

func NewWebsocketTestCommand() *cli.Command {
	return &cli.Command{
		Name:   "websocketTest",
		Short:  i18n.T("Websocket Test", "WebSocket测试"),
		Usage:  "aliyun websocketTest",
		Hidden: false,
		Run: func(ctx *cli.Context, args []string) error {
			return runWebsocketTest(ctx, args)
		},
	}
}

func runWebsocketTest(ctx *cli.Context, args []string) error {
	profile, _ := config.LoadProfileWithContext(ctx)
	credential, err := profile.GetCredential(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	conf := &openapiClient.Config{
		Credential: credential,
		RegionId:   tea.String(profile.RegionId),
		Endpoint:   tea.String("openapi-mcp.cn-hangzhou.aliyuncs.com"),
	}

	client, err := openapiClient.NewClient(conf)
	if err != nil {
		return err
	}

	params := &openapiClient.Params{
		Action:      tea.String("ListApiMcpServers"),
		Version:     tea.String("2024-11-30"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		Pathname:    tea.String("/apimcpservers"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	queries := map[string]interface{}{}
	queries["id"] = tea.String("Bt4Td5W1tI31YAsu")

	runtime := &dara.RuntimeOptions{}
	request := &openapiClient.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	response, err := client.CallApi(params, request, runtime)
	if err != nil {
		return err
	}
	// bodyBytes, _ := GetContentFromApiResponse(response)
	fmt.Printf("response: %s\n", response["statusCode"])

	// testAwapWebSocket(ctx, args)
	testAwapWebSocketWithoutHandleRawMessage(ctx, args) // 新增：不重写 HandleRawMessage 的用例
	// testGeneralWebSocket(ctx, args)
	// testSequentialMessageReception(ctx, args)
	return nil
}

func GetContentFromApiResponse(response map[string]any) ([]byte, error) {
	responseBody := response["body"]
	if responseBody == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	switch v := responseBody.(type) {
	case string:
		return []byte(v), nil
	case map[string]any, []any:
		jsonData, _ := json.Marshal(v)
		return jsonData, nil
	case []byte:
		return v, nil
	default:
		return []byte(fmt.Sprintf("%v", v)), nil
	}
}

func printSessionInfo(session *dara.WebSocketSessionInfo) {
	if session == nil {
		fmt.Println("  [Session] nil")
		return
	}
	fmt.Printf("  [Session] ID: %s\n", session.SessionID)
	if session.RequestID != "" {
		fmt.Printf("  [Session] RequestID: %s\n", session.RequestID)
	}
	fmt.Printf("  [Session] ConnectedAt: %s\n", session.ConnectedAt.Format(time.RFC3339))
	fmt.Printf("  [Session] RemoteAddr: %s\n", session.RemoteAddr)
	fmt.Printf("  [Session] LocalAddr: %s\n", session.LocalAddr)
	if len(session.Attributes) > 0 {
		attrJSON, _ := json.Marshal(session.Attributes)
		fmt.Printf("  [Session] Attributes: %s\n", string(attrJSON))
	}
}

type SimpleHandler struct {
	dara.AbstractAwapWebSocketHandler
}

func (h *SimpleHandler) AfterConnectionEstablished(session *dara.WebSocketSessionInfo) error {
	fmt.Println("✓ CLI Connected to WebSocket server")
	printSessionInfo(session)
	return nil
}

func (h *SimpleHandler) HandleAwapMessage(session *dara.WebSocketSessionInfo, message *dara.AwapMessage) error {
	fmt.Println("📨 CLI Received AWAP message:")
	jsonData, _ := json.Marshal(message)
	fmt.Printf("  Message: %s\n", string(jsonData))
	return nil
}

func (h *SimpleHandler) HandleAwapIncomingMessage(session *dara.WebSocketSessionInfo, message *dara.AwapIncomingMessage) error {
	fmt.Println("📬 CLI Received AWAP incoming message:")
	jsonData, _ := json.Marshal(message)
	fmt.Printf("  Incoming Message: %s\n", string(jsonData))
	return nil
}

func (h *SimpleHandler) HandleRawMessage(session *dara.WebSocketSessionInfo, message *dara.WebSocketMessage) error {
	// Parse the AWAP message ourselves and call HandleAwapMessage directly
	// This avoids the issue where AbstractAwapWebSocketHandler.HandleRawMessage
	// can't access the outer SimpleHandler type
	awapMsg, err := dara.ParseAwapMessage(message)
	if err != nil {
		fmt.Printf("[CLI Simple] Failed to parse AWAP message: %v\n", err)
		return err
	}

	// Call HandleAwapMessage directly on h (which is *SimpleHandler)
	// This will call our overridden implementation
	if err := h.HandleAwapMessage(session, awapMsg); err != nil {
		return err
	}

	// Also call HandleAwapIncomingMessage for event types
	if awapMsg.Type == dara.AwapMessageTypeUpstreamTextEvent ||
		awapMsg.Type == dara.AwapMessageTypeUpstreamBinaryEvent ||
		awapMsg.Type == dara.AwapMessageTypeAckRequiredTextEvent ||
		awapMsg.Type == dara.AwapMessageTypeMessageReceiveEvent ||
		awapMsg.Type == dara.AwapMessageTypeDownstreamTextEvent ||
		awapMsg.Type == dara.AwapMessageTypeDownstreamBinaryEvent {
		incoming := &dara.AwapIncomingMessage{
			AwapMessage: *awapMsg,
			RawPayload:  message.Payload,
		}
		return h.HandleAwapIncomingMessage(session, incoming)
	}

	return nil
}

func (h *SimpleHandler) SupportsPartialMessages() bool {
	return false
}

func (h *SimpleHandler) HandleError(session *dara.WebSocketSessionInfo, err error) error {
	fmt.Printf("❌ CLI Error: %v\n", err)
	printSessionInfo(session)
	return nil
}

func (h *SimpleHandler) AfterConnectionClosed(session *dara.WebSocketSessionInfo, code int, reason string) error {
	fmt.Printf("✗ CLI Connection closed (code: %d, reason: %s)\n", code, reason)
	printSessionInfo(session)
	return nil
}

func testAwapWebSocket(ctx *cli.Context, args []string) error {
	fmt.Println("=== WebSocket Example ===")
	profile, _ := config.LoadProfileWithContext(ctx)
	credential, err := profile.GetCredential(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}
	config := &openapiClient.Config{
		Credential: credential,
		Endpoint:   dara.String("dalutest-pre.aliyuncs.com"),
		Protocol:   dara.String("https"),
	}

	apiClient, err := openapiClient.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Setup WebSocket
	params := &openapiClient.Params{
		// Action:      tea.String("ListApiMcpServers"),
		// Version:     tea.String("2024-11-30"),
		// Protocol:    tea.String("HTTPS"),
		// Method:      tea.String("GET"),
		// AuthType:    tea.String("AK"),
		// Style:       tea.String("ROA"),
		// Pathname:    tea.String("/apimcpservers"),
		// ReqBodyType: tea.String("json"),
		// BodyType:    tea.String("json"),
		// Product:  dara.String("DaluTestInner"),
		Action:   tea.String("WebsocketAwapDemoApi"),
		Version:  tea.String("2022-02-02"),
		Protocol: tea.String("wss"),
		Method:   tea.String("GET"),
		Pathname: tea.String("/ws/awap-demo-api"),
		AuthType: tea.String("AK"),
	}

	request := &openapiClient.OpenApiRequest{
		Headers: map[string]*string{
			"Sec-Websocket-Protocol": tea.String("awap"),
		},
	}

	runtime := &dara.RuntimeOptions{
		ReadTimeout:                dara.Int(60000),  // 60 seconds
		ConnectTimeout:             dara.Int(30000),  // 30 seconds (increased for slow networks)
		WebSocketPingInterval:      dara.Int(30000),  // 30秒心跳
		WebSocketHandshakeTimeout:  dara.Int(30000),  // 30秒握手超时（增加以应对网络延迟）
		WebSocketWriteTimeout:      dara.Int(30000),  // 30秒写入超时（增加以应对网络延迟）
		WebSocketEnableReconnect:   dara.Bool(true),  // 启用重连
		WebSocketMaxReconnectTimes: dara.Int(5),      // 最多重连5次
		WebSocketHandler:           &SimpleHandler{}, // 通过 runtime 配置 handler
	}

	fmt.Println("Connecting...")
	// Handler 从 runtime 中获取
	result, err := apiClient.DoWebSocketRequest(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	wsClient := result["wsClient"].(dara.WebSocketClient)
	defer wsClient.Close()

	fmt.Println("\nSending AWAP message...")

	// 方法 1: 使用 SendAwapRequest 发送请求消息（推荐）
	err = apiClient.SendAwapRequest(wsClient, "msg-001", 1, map[string]interface{}{
		"action": "test",
		"data":   "Hello WebSocket!",
	})
	if err != nil {
		log.Printf("Failed to send AWAP request: %v", err)
	}

	// 方法 2: 使用 SendAwapEvent 发送事件消息
	err = apiClient.SendAwapEvent(wsClient, "msg-002", 2, map[string]interface{}{
		"eventType": "userAction",
		"message":   "Hello WebSocket!",
	})
	if err != nil {
		log.Printf("Failed to send AWAP event: %v", err)
	}

	// 方法 3: 手动构建 AWAP 消息
	awapMsg := apiClient.BuildAwapRequest("msg-003", 3, map[string]interface{}{
		"message": "Hello WebSocket!",
	})
	err = apiClient.SendAwapMessage(wsClient, awapMsg)
	if err != nil {
		log.Printf("Failed to send AWAP message: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 方法 4: 发送 AckRequiredTextEvent 类型的消息（需要确认的消息）
	fmt.Println("\n4. Sending AckRequiredTextEvent message...")
	ackRequiredMsg := apiClient.BuildAwapMessage(dara.AwapMessageTypeAckRequiredTextEvent, "msg-004", 4, map[string]interface{}{
		"action":    "ackRequiredTest",
		"data":      "This message requires acknowledgment",
		"timestamp": time.Now().Unix(),
	})
	err = apiClient.SendAwapMessage(wsClient, ackRequiredMsg)
	if err != nil {
		log.Printf("Failed to send AckRequiredTextEvent message: %v", err)
	} else {
		fmt.Println("✓ AckRequiredTextEvent message sent successfully")
	}

	// Wait for response
	time.Sleep(3 * time.Second)

	fmt.Println("\n=== AWAP Example Complete ===")
	return nil
}

type GeneralHandler struct {
	dara.AbstractGeneralWebSocketHandler
}

func (h *GeneralHandler) AfterConnectionEstablished(session *dara.WebSocketSessionInfo) error {
	fmt.Println("✓ CLI General Connected to General WebSocket server")
	printSessionInfo(session)
	return nil
}

func (h *GeneralHandler) HandleGeneralTextMessage(session *dara.WebSocketSessionInfo, message *dara.GeneralMessage) error {
	fmt.Println("📨 CLI General Received General text message:")
	jsonData, _ := json.Marshal(message)
	fmt.Printf("  Message: %s\n", string(jsonData))
	return nil
}

func (h *GeneralHandler) HandleGeneralBinaryMessage(session *dara.WebSocketSessionInfo, data []byte) error {
	fmt.Println("📦 CLI General Received General binary message:")
	fmt.Printf("  Size: %d bytes\n", len(data))
	fmt.Printf("  Content: %s\n", string(data))
	return nil
}

func (h *GeneralHandler) HandleGeneralIncomingMessage(session *dara.WebSocketSessionInfo, message *dara.GeneralIncomingMessage) error {
	if message.IsBinary {
		fmt.Println("📦 CLI General incoming binary message:")
	} else {
		fmt.Println("📨 CLI General incoming message:")
	}
	if message.IsBinary {
		fmt.Printf("  Size: %d bytes\n", len(message.RawPayload))
	} else {
		jsonData, _ := json.Marshal(message.Body)
		fmt.Printf("  Body: %s\n", string(jsonData))
	}
	return nil
}

func (h *GeneralHandler) HandleRawMessage(session *dara.WebSocketSessionInfo, message *dara.WebSocketMessage) error {
	// Parse and handle General messages directly
	// This avoids the issue where AbstractGeneralWebSocketHandler.HandleRawMessage
	// can't access the outer GeneralHandler type
	if message.Type == dara.WebSocketMessageTypeText {
		// Parse as General text message
		generalMsg, err := dara.ParseGeneralMessage(message)
		if err != nil {
			fmt.Printf("[CLI General] Failed to parse General message: %v\n", err)
			return err
		}

		// Check message type from headers
		msgType := generalMsg.Headers["type"]
		fmt.Printf("[CLI General] Received text message, type from header: %s\n", msgType)

		// Create incoming message
		incoming := &dara.GeneralIncomingMessage{
			Headers:    generalMsg.Headers,
			Body:       generalMsg.Body,
			RawPayload: message.Payload,
			IsBinary:   false,
		}

		// Call both handlers directly on h (which is *GeneralHandler)
		if err := h.HandleGeneralTextMessage(session, generalMsg); err != nil {
			return err
		}
		return h.HandleGeneralIncomingMessage(session, incoming)

	} else if message.Type == dara.WebSocketMessageTypeBinary {
		// Handle as binary message
		fmt.Printf("[CLI General] Received binary message\n")

		incoming := &dara.GeneralIncomingMessage{
			Headers:    make(map[string]string),
			Body:       nil,
			RawPayload: message.Payload,
			IsBinary:   true,
		}

		// Call both handlers directly on h (which is *GeneralHandler)
		if err := h.HandleGeneralBinaryMessage(session, message.Payload); err != nil {
			return err
		}
		return h.HandleGeneralIncomingMessage(session, incoming)
	}

	return nil
}

func (h *GeneralHandler) SupportsPartialMessages() bool {
	return false
}

func (h *GeneralHandler) HandleError(session *dara.WebSocketSessionInfo, err error) error {
	fmt.Printf("❌ CLI General Error: %v\n", err)
	printSessionInfo(session)
	return nil
}

func (h *GeneralHandler) AfterConnectionClosed(session *dara.WebSocketSessionInfo, code int, reason string) error {
	fmt.Printf("✗ CLI General Connection closed (code: %d, reason: %s)\n", code, reason)
	printSessionInfo(session)
	return nil
}

// NoHandleRawMessageHandler 是一个不重写 HandleRawMessage 的 handler
// 这个 handler 依赖 AbstractAwapWebSocketHandler.HandleRawMessage
// 来展示 AwapWebSocketHandler 接口的实际使用
type NoHandleRawMessageHandler struct {
	dara.AbstractAwapWebSocketHandler
	messageCount int
	mu           sync.Mutex
}

func (h *NoHandleRawMessageHandler) AfterConnectionEstablished(session *dara.WebSocketSessionInfo) error {
	fmt.Println("✓ CLI NoHandleRawMessage - Connected to WebSocket server")
	fmt.Println("  Note: This handler does NOT override HandleRawMessage")
	fmt.Println("  It uses AbstractAwapWebSocketHandler.HandleRawMessage")
	fmt.Println("  which uses AwapWebSocketHandler interface")
	printSessionInfo(session)
	return nil
}

func (h *NoHandleRawMessageHandler) HandleAwapMessage(session *dara.WebSocketSessionInfo, message *dara.AwapMessage) error {
	h.mu.Lock()
	h.messageCount++
	count := h.messageCount
	h.mu.Unlock()

	fmt.Printf("📨 CLI NoHandleRawMessage - HandleAwapMessage called (#%d):\n", count)
	fmt.Printf("  Type: %s\n", message.Type)
	fmt.Printf("  ID: %s\n", message.ID)
	fmt.Printf("  Seq: %d\n", message.Seq)
	if message.Payload != nil {
		payloadJSON, _ := json.Marshal(message.Payload)
		fmt.Printf("  Payload: %s\n", string(payloadJSON))
	}
	printSessionInfo(session)
	return nil
}

func (h *NoHandleRawMessageHandler) HandleAwapIncomingMessage(session *dara.WebSocketSessionInfo, message *dara.AwapIncomingMessage) error {
	fmt.Printf("📬 CLI NoHandleRawMessage - HandleAwapIncomingMessage called:\n")
	fmt.Printf("  Type: %s\n", message.Type)
	fmt.Printf("  ID: %s\n", message.ID)
	fmt.Printf("  Seq: %d\n", message.Seq)
	fmt.Printf("  RawPayload length: %d bytes\n", len(message.RawPayload))
	printSessionInfo(session)
	return nil
}

func (h *NoHandleRawMessageHandler) HandleError(session *dara.WebSocketSessionInfo, err error) error {
	fmt.Printf("❌ CLI NoHandleRawMessage - Error: %v\n", err)
	printSessionInfo(session)
	return nil
}

func (h *NoHandleRawMessageHandler) AfterConnectionClosed(session *dara.WebSocketSessionInfo, code int, reason string) error {
	h.mu.Lock()
	count := h.messageCount
	h.mu.Unlock()

	fmt.Printf("✗ CLI NoHandleRawMessage - Connection closed (code: %d, reason: %s)\n", code, reason)
	fmt.Printf("  Total messages received: %d\n", count)
	printSessionInfo(session)
	return nil
}

func (h *NoHandleRawMessageHandler) SupportsPartialMessages() bool {
	return false
}

func testAwapWebSocketWithoutHandleRawMessage(ctx *cli.Context, args []string) error {
	fmt.Println("\n=== AWAP WebSocket Test (Without HandleRawMessage Override) ===")
	fmt.Println("This test demonstrates AwapWebSocketHandler interface usage")
	fmt.Println("The handler does NOT override HandleRawMessage,")
	fmt.Println("so it uses AbstractAwapWebSocketHandler.HandleRawMessage")
	fmt.Println("which uses AwapWebSocketHandler interface for type assertion")

	profile, _ := config.LoadProfileWithContext(ctx)
	credential, err := profile.GetCredential(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}
	config := &openapiClient.Config{
		Credential: credential,
		Endpoint:   dara.String("dalutest-pre.aliyuncs.com"),
		Protocol:   dara.String("https"),
	}

	apiClient, err := openapiClient.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	params := &openapiClient.Params{
		Action:   tea.String("WebsocketAwapDemoApi"),
		Version:  tea.String("2022-02-02"),
		Protocol: tea.String("wss"),
		Method:   tea.String("GET"),
		Pathname: tea.String("/ws/awap-demo-api"),
		AuthType: tea.String("AK"),
	}

	request := &openapiClient.OpenApiRequest{
		Headers: map[string]*string{
			"Sec-Websocket-Protocol": tea.String("awap"),
		},
	}

	runtime := &dara.RuntimeOptions{
		ReadTimeout:                dara.Int(60000),              // 60 seconds
		ConnectTimeout:             dara.Int(30000),              // 30 seconds
		WebSocketPingInterval:      dara.Int(30000),              // 30秒心跳
		WebSocketHandshakeTimeout:  dara.Int(30000),              // 30秒握手超时
		WebSocketWriteTimeout:      dara.Int(30000),              // 30秒写入超时
		WebSocketEnableReconnect:   dara.Bool(true),              // 启用重连
		WebSocketMaxReconnectTimes: dara.Int(5),                  // 最多重连5次
		WebSocketHandler:           &NoHandleRawMessageHandler{}, // 使用不重写 HandleRawMessage 的 handler
	}

	fmt.Println("Connecting...")
	result, err := apiClient.DoWebSocketRequest(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	wsClient := result["wsClient"].(dara.WebSocketClient)
	defer wsClient.Close()

	fmt.Println("\nSending AWAP messages...")

	// 方法 1: 使用 SendAwapRequest 发送请求消息
	fmt.Println("1. Sending AWAP request message...")
	err = apiClient.SendAwapRequest(wsClient, "msg-no-handleraw-001", 1, map[string]interface{}{
		"action": "test",
		"data":   "This handler does NOT override HandleRawMessage",
	})
	if err != nil {
		log.Printf("Failed to send AWAP request: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 方法 2: 使用 SendAwapEvent 发送事件消息
	fmt.Println("2. Sending AWAP event message...")
	err = apiClient.SendAwapEvent(wsClient, "msg-no-handleraw-002", 2, map[string]interface{}{
		"eventType": "testEvent",
		"message":   "Testing AwapWebSocketHandler interface usage",
	})
	if err != nil {
		log.Printf("Failed to send AWAP event: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 方法 3: 发送 AckRequiredTextEvent 类型的消息
	fmt.Println("3. Sending AckRequiredTextEvent message...")
	ackRequiredMsg := apiClient.BuildAwapMessage(dara.AwapMessageTypeAckRequiredTextEvent, "msg-no-handleraw-003", 3, map[string]interface{}{
		"action":    "ackRequiredTest",
		"data":      "This message requires acknowledgment",
		"timestamp": time.Now().Unix(),
	})
	err = apiClient.SendAwapMessage(wsClient, ackRequiredMsg)
	if err != nil {
		log.Printf("Failed to send AckRequiredTextEvent message: %v", err)
	}

	// Wait for response
	fmt.Println("\nWaiting for server responses...")
	time.Sleep(3 * time.Second)

	fmt.Println("\n=== NoHandleRawMessage Test Complete ===")
	fmt.Println("This test demonstrates that AwapWebSocketHandler interface")
	fmt.Println("is used when handler does NOT override HandleRawMessage")
	return nil
}

func testGeneralWebSocket(ctx *cli.Context, args []string) error {
	fmt.Println("\n=== General WebSocket Example ===")
	profile, _ := config.LoadProfileWithContext(ctx)
	credential, err := profile.GetCredential(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}
	config := &openapiClient.Config{
		Credential: credential,
		Endpoint:   dara.String("dalutest-pre.aliyuncs.com"),
		Protocol:   dara.String("https"),
	}

	apiClient, err := openapiClient.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	params := &openapiClient.Params{
		Action:   tea.String("WebsocketGeneralDemoApi"),
		Version:  tea.String("2022-02-02"),
		Protocol: tea.String("wss"),
		Method:   tea.String("GET"),
		Pathname: tea.String("/ws/general-demo-api"),
		AuthType: tea.String("AK"),
	}

	request := &openapiClient.OpenApiRequest{
		Headers: map[string]*string{
			"Sec-Websocket-Protocol": tea.String("general"),
		},
	}

	runtime := &dara.RuntimeOptions{
		ReadTimeout:                dara.Int(60000),   // 60 seconds
		ConnectTimeout:             dara.Int(30000),   // 30 seconds
		WebSocketPingInterval:      dara.Int(30000),   // 30秒心跳
		WebSocketHandshakeTimeout:  dara.Int(30000),   // 30秒握手超时
		WebSocketWriteTimeout:      dara.Int(30000),   // 30秒写入超时
		WebSocketEnableReconnect:   dara.Bool(true),   // 启用重连
		WebSocketMaxReconnectTimes: dara.Int(5),       // 最多重连5次
		WebSocketHandler:           &GeneralHandler{}, // 通过 runtime 配置 handler
	}

	fmt.Println("Connecting to General WebSocket...")
	// Handler 从 runtime 中获取
	result, err := apiClient.DoWebSocketRequest(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	wsClient := result["wsClient"].(dara.WebSocketClient)
	defer wsClient.Close()

	fmt.Println("\nSending General messages...")

	// 方法 1: 发送文本消息
	fmt.Println("1. Sending General text message...")
	err = apiClient.SendGeneralTextMessage(wsClient, "Hello General WebSocket!")
	if err != nil {
		log.Printf("Failed to send General text message: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 方法 2: 发送 JSON 消息
	fmt.Println("2. Sending General JSON message...")
	eventData := map[string]interface{}{
		"name": "general-test",
		"object": map[string]interface{}{
			"field1": "general",
			"field2": 2,
			"field3": []string{"test"},
		},
		"list": []map[string]interface{}{
			{"enabled": false, "value": "general"},
		},
		"map": map[string]string{
			"test": "test",
		},
	}
	err = apiClient.SendGeneralJSONMessage(wsClient, eventData)
	if err != nil {
		log.Printf("Failed to send General JSON message: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 方法 3: 发送带自定义头部的消息
	fmt.Println("3. Sending General message with custom headers...")
	generalMsg := apiClient.BuildGeneralJSONMessage(map[string]interface{}{
		"action": "test",
		"data":   "Hello with headers!",
	})
	generalMsg.WithHeader("X-Custom-Header", "custom-value")
	generalMsg.WithHeader("Content-Type", "application/json")
	err = apiClient.SendGeneralMessage(wsClient, generalMsg)
	if err != nil {
		log.Printf("Failed to send General message with headers: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 方法 4: 发送二进制消息
	fmt.Println("4. Sending General binary message...")
	binaryData := []byte("Binary General Data")
	err = apiClient.SendGeneralBinaryMessage(wsClient, binaryData)
	if err != nil {
		log.Printf("Failed to send General binary message: %v", err)
	}

	// Wait for response
	time.Sleep(3 * time.Second)

	fmt.Println("\n=== General Example Complete ===")
	return nil
}

type SequentialHandler struct {
	dara.AbstractAwapWebSocketHandler
	receivedSeq   []int64 // Track received sequence numbers
	mu            sync.Mutex
	expectedCount int
	done          chan struct{}
}

func NewSequentialHandler(expectedCount int) *SequentialHandler {
	return &SequentialHandler{
		receivedSeq:   make([]int64, 0, expectedCount),
		expectedCount: expectedCount,
		done:          make(chan struct{}),
	}
}

func (h *SequentialHandler) AfterConnectionEstablished(session *dara.WebSocketSessionInfo) error {
	fmt.Println("✓ CLI Sequential Test - Connected to WebSocket server")
	printSessionInfo(session)
	return nil
}

func (h *SequentialHandler) HandleAwapMessage(session *dara.WebSocketSessionInfo, message *dara.AwapMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Printf("📨 CLI Sequential Test - HandleAwapMessage called: seq=%d, type=%s, id=%s\n",
		message.Seq, message.Type, message.ID)
	printSessionInfo(session)

	// Print message payload for debugging
	if message.Payload != nil {
		payloadJSON, _ := json.Marshal(message.Payload)
		fmt.Printf("  Payload: %s\n", string(payloadJSON))
	}

	// Track sequence number
	h.receivedSeq = append(h.receivedSeq, message.Seq)
	fmt.Printf("  Progress: %d/%d messages received\n", len(h.receivedSeq), h.expectedCount)

	// Check if we've received all expected messages
	if len(h.receivedSeq) >= h.expectedCount {
		if h.done != nil {
			close(h.done)
			h.done = nil // Prevent double close
		}
	}

	return nil
}

func (h *SequentialHandler) HandleAwapIncomingMessage(session *dara.WebSocketSessionInfo, message *dara.AwapIncomingMessage) error {
	// This is called after HandleAwapMessage, so we don't need to count again
	// Just log for debugging
	fmt.Printf("📬 CLI Sequential Test - HandleAwapIncomingMessage called: seq=%d, type=%s, id=%s\n",
		message.Seq, message.Type, message.ID)
	return nil
}

func (h *SequentialHandler) HandleRawMessage(session *dara.WebSocketSessionInfo, message *dara.WebSocketMessage) error {
	// Parse the AWAP message ourselves and call HandleAwapMessage directly
	// This avoids the issue where AbstractAwapWebSocketHandler.HandleRawMessage
	// can't access the outer SequentialHandler type
	awapMsg, err := dara.ParseAwapMessage(message)
	if err != nil {
		fmt.Printf("[CLI Sequential] Failed to parse AWAP message: %v\n", err)
		return err
	}

	fmt.Printf("[CLI Sequential] Calling HandleAwapMessage directly: seq=%d, type=%s\n", awapMsg.Seq, awapMsg.Type)
	if err := h.HandleAwapMessage(session, awapMsg); err != nil {
		return err
	}

	if awapMsg.Type == dara.AwapMessageTypeUpstreamTextEvent ||
		awapMsg.Type == dara.AwapMessageTypeUpstreamBinaryEvent ||
		awapMsg.Type == dara.AwapMessageTypeAckRequiredTextEvent ||
		awapMsg.Type == dara.AwapMessageTypeMessageReceiveEvent ||
		awapMsg.Type == dara.AwapMessageTypeDownstreamTextEvent ||
		awapMsg.Type == dara.AwapMessageTypeDownstreamBinaryEvent {
		incoming := &dara.AwapIncomingMessage{
			AwapMessage: *awapMsg,
			RawPayload:  message.Payload,
		}
		return h.HandleAwapIncomingMessage(session, incoming)
	}

	return nil
}

func (h *SequentialHandler) SupportsPartialMessages() bool {
	return false
}

func (h *SequentialHandler) HandleError(session *dara.WebSocketSessionInfo, err error) error {
	fmt.Printf("❌ CLI Sequential Test Error: %v\n", err)
	printSessionInfo(session)
	return nil
}

func (h *SequentialHandler) AfterConnectionClosed(session *dara.WebSocketSessionInfo, code int, reason string) error {
	fmt.Printf("✗ CLI Sequential Test - Connection closed (code: %d, reason: %s)\n", code, reason)
	printSessionInfo(session)
	return nil
}

func (h *SequentialHandler) GetReceivedSeq() []int64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]int64, len(h.receivedSeq))
	copy(result, h.receivedSeq)
	return result
}

func (h *SequentialHandler) WaitForCompletion(timeout time.Duration) bool {
	h.mu.Lock()
	done := h.done
	h.mu.Unlock()

	if done == nil {
		return true // Already completed
	}
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func testSequentialMessageReception(ctx *cli.Context, args []string) error {
	fmt.Println("\n=== Sequential Message Reception Test ===")
	fmt.Println("This test verifies that messages are received in the correct order")

	profile, _ := config.LoadProfileWithContext(ctx)
	credential, err := profile.GetCredential(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	config := &openapiClient.Config{
		Credential: credential,
		Endpoint:   dara.String("dalutest-pre.aliyuncs.com"),
		Protocol:   dara.String("https"),
	}

	apiClient, err := openapiClient.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	params := &openapiClient.Params{
		Action:   tea.String("WebsocketAwapDemoApi"),
		Version:  tea.String("2022-02-02"),
		Protocol: tea.String("wss"),
		Method:   tea.String("GET"),
		Pathname: tea.String("/ws/awap-demo-api"),
		AuthType: tea.String("AK"),
	}

	request := &openapiClient.OpenApiRequest{
		Query: map[string]*string{
			"delay":           tea.String("3000"), // Delay between messages (ms)
			"batchSendMsgCnt": tea.String("20"),   // Number of messages to send
		},
		Headers: map[string]*string{
			"Sec-Websocket-Protocol": tea.String("awap"),
		},
	}

	expectedCount := 20
	handler := NewSequentialHandler(expectedCount)

	runtime := &dara.RuntimeOptions{
		ReadTimeout:                dara.Int(120000), // 120 seconds (enough for 20 messages with 3s delay)
		ConnectTimeout:             dara.Int(30000),  // 30 seconds
		WebSocketPingInterval:      dara.Int(30000),  // 30秒心跳
		WebSocketHandshakeTimeout:  dara.Int(30000),  // 30秒握手超时
		WebSocketWriteTimeout:      dara.Int(30000),  // 30秒写入超时
		WebSocketEnableReconnect:   dara.Bool(true),  // 启用重连
		WebSocketMaxReconnectTimes: dara.Int(5),      // 最多重连5次
		WebSocketHandler:           handler,          // 通过 runtime 配置 handler
	}

	fmt.Println("Connecting to WebSocket server...")
	fmt.Printf("Query parameters: delay=%s, batchSendMsgCnt=%s\n",
		dara.StringValue(request.Query["delay"]),
		dara.StringValue(request.Query["batchSendMsgCnt"]))

	result, err := apiClient.DoWebSocketRequest(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	wsClient := result["wsClient"].(dara.WebSocketClient)
	defer wsClient.Close()

	fmt.Println("Connection established, waiting for server to start sending batch messages...")
	time.Sleep(2 * time.Second)

	fmt.Println("\nSending request to trigger batch message sending...")
	err = apiClient.SendAwapRequest(wsClient, "sequential-test-001", 1, map[string]interface{}{
		"action": "batchSend",
		"test":   "sequential",
	})
	if err != nil {
		log.Printf("Failed to send request: %v", err)
	} else {
		fmt.Println("Request sent successfully")
	}

	fmt.Printf("\nWaiting for %d messages to be received (timeout: 100s)...\n", expectedCount)
	fmt.Println("Note: Server will send messages with 3s delay between each, so total time ~60s")

	progressTicker := time.NewTicker(5 * time.Second)
	defer progressTicker.Stop()

	doneChan := make(chan bool, 1)
	go func() {
		timeout := 100 * time.Second // 20 messages * 3s delay + buffer
		doneChan <- handler.WaitForCompletion(timeout)
	}()

	// Show progress while waiting
	for {
		select {
		case success := <-doneChan:
			if success {
				fmt.Println("\n✓ All messages received!")
			} else {
				received := len(handler.GetReceivedSeq())
				fmt.Printf("\n⚠ Timeout waiting for all messages. Received %d/%d messages\n",
					received, expectedCount)
				if received == 0 {
					fmt.Println("\n⚠ No messages received at all. Possible issues:")
					fmt.Println("  1. Server may not support batchSendMsgCnt parameter")
					fmt.Println("  2. Query parameters may not be passed correctly")
					fmt.Println("  3. Server may require a different trigger mechanism")
				}
			}
			goto done
		case <-progressTicker.C:
			received := len(handler.GetReceivedSeq())
			fmt.Printf("  Progress: %d/%d messages received...\n", received, expectedCount)
		}
	}
done:

	receivedSeq := handler.GetReceivedSeq()
	fmt.Printf("\n=== Sequence Verification ===\n")
	fmt.Printf("Expected count: %d\n", expectedCount)
	fmt.Printf("Received count: %d\n", len(receivedSeq))

	if len(receivedSeq) > 0 {
		fmt.Printf("Received sequence numbers: %v\n", receivedSeq)

		isOrdered := true
		for i := 1; i < len(receivedSeq); i++ {
			if receivedSeq[i] < receivedSeq[i-1] {
				isOrdered = false
				fmt.Printf("❌ Out of order detected: seq[%d]=%d < seq[%d]=%d\n",
					i, receivedSeq[i], i-1, receivedSeq[i-1])
				break
			}
		}

		if isOrdered {
			fmt.Println("✅ All messages received in correct order!")
		} else {
			fmt.Println("❌ Messages received out of order!")
		}

		// Check for missing sequence numbers
		missing := []int64{}
		seqMap := make(map[int64]bool)
		for _, seq := range receivedSeq {
			seqMap[seq] = true
		}
		for i := int64(1); i <= int64(expectedCount); i++ {
			if !seqMap[i] {
				missing = append(missing, i)
			}
		}
		if len(missing) > 0 {
			fmt.Printf("⚠ Missing sequence numbers: %v\n", missing)
		} else {
			fmt.Println("✅ No missing sequence numbers")
		}
	} else {
		fmt.Println("❌ No messages received!")
	}

	fmt.Println("\n=== Sequential Message Reception Test Complete ===")
	return nil
}
