package websockettest

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/aliyun/aliyun-cli/v3/cli"
	"github.com/aliyun/aliyun-cli/v3/config"
	"github.com/aliyun/aliyun-cli/v3/i18n"

	openapiClient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	websocketutils "github.com/alibabacloud-go/darabonba-openapi/v2/websocketutils"
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
	// profile, _ := config.LoadProfileWithContext(ctx)
	// credential, err := profile.GetCredential(ctx, nil)
	// if err != nil {
	// 	return fmt.Errorf("failed to get credential: %w", err)
	// }

	// conf := &openapiClient.Config{
	// 	Credential: credential,
	// 	RegionId:   tea.String(profile.RegionId),
	// 	Endpoint:   tea.String("openapi-mcp.cn-hangzhou.aliyuncs.com"),
	// }

	// client, err := openapiClient.NewClient(conf)
	// if err != nil {
	// 	return err
	// }

	// params := &openapiClient.Params{
	// 	Action:      tea.String("ListApiMcpServers"),
	// 	Version:     tea.String("2024-11-30"),
	// 	Protocol:    tea.String("HTTPS"),
	// 	Method:      tea.String("GET"),
	// 	AuthType:    tea.String("AK"),
	// 	Style:       tea.String("ROA"),
	// 	Pathname:    tea.String("/apimcpservers"),
	// 	ReqBodyType: tea.String("json"),
	// 	BodyType:    tea.String("json"),
	// }
	// queries := map[string]interface{}{}
	// queries["id"] = tea.String("Bt4Td5W1tI31YAsu")

	// runtime := &dara.RuntimeOptions{}
	// request := &openapiClient.OpenApiRequest{
	// 	Query: openapiutil.Query(queries),
	// }
	// response, err := client.CallApi(params, request, runtime)
	// if err != nil {
	// 	return err
	// }
	// // bodyBytes, _ := GetContentFromApiResponse(response)
	// fmt.Printf("response: %s\n", response["statusCode"])

	testAwapWebSocketBinary(ctx, args)
	// testAwapWebSocket(ctx, args)
	// testAwapWebSocketWithoutHandleAwapMessage(ctx, args) // 重写 HandleRawMessage 的用例
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
	fmt.Printf("  Type: %s\n", message.Format)
	fmt.Printf("  ID: %s\n", message.ID)
	fmt.Printf("  Headers: %+v\n", message.Headers)
	return nil
}

func (h *SimpleHandler) HandleRawMessage(session *dara.WebSocketSessionInfo, message *dara.WebSocketMessage) error {
	fmt.Println("📨 CLI Received AWAP message from HandleRawMessage!!!! shouldn't be")
	return nil
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

func testAwapWebSocketBinary(ctx *cli.Context, args []string) error {
	fmt.Println("=== WebSocket Binary Example ===")
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
		Action:               tea.String("WebsocketAwapDemoApi"),
		Version:              tea.String("2022-02-02"),
		Protocol:             tea.String("wss"),
		Method:               tea.String("GET"),
		Pathname:             tea.String("/ws/awap-demo-api"),
		AuthType:             tea.String("AK"),
		WebsocketSubProtocol: tea.String("awap"),
	}

	request := &openapiClient.OpenApiRequest{}

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
	result, err := apiClient.CallApi(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	websocketObj := result["websocketClient"].(*websocketutils.WebSocketClient)
	defer websocketObj.Close()

	sessionInfo := websocketObj.GetSessionInfo()
	sessionId := sessionInfo.SessionID
	if sessionId == "" {
		return fmt.Errorf("session ID is empty")
	}

	fmt.Println("\n1. Sending AWAP message...")

	// 方法 1: 发送 binary 信息
	awapMsgBinary := websocketutils.NewAwapMessage(dara.AwapMessageType("UpstreamBinaryEvent"),
		"msg-001",
		[]byte("Hello WebSocket!"),
	)
	awapMsgBinary.WithHeader("session-id", sessionId)
	err = websocketObj.SendAwapBinaryMessage(awapMsgBinary)
	if err != nil {
		log.Printf("Failed to send AWAP message: %v", err)
	} else {
		fmt.Printf("AWAP message sent successfully, type: %s\n", dara.AwapMessageType("UpstreamBinaryEvent"))
	}
	time.Sleep(5 * time.Second)

	apiClientTrigger, err := openapiClient.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// trigger binary from server
	params = &openapiClient.Params{
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
		Action:      tea.String("WebsocketServerExecute"),
		Version:     tea.String("2022-02-02"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		Pathname:    tea.String("/ws_server/execute"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	request = &openapiClient.OpenApiRequest{
		Query: openapiutil.Query(map[string]interface{}{
			"sessionId": sessionId,
			"action":    "sendBinary",
		}),
	}

	runtimeTrigger := &dara.RuntimeOptions{}

	fmt.Println("Triggering binary from server...")
	result, err = apiClientTrigger.CallApi(params, request, runtimeTrigger)
	if err != nil {
		log.Fatalf("Failed to trigger binary from server: %v", err)
	}
	fmt.Printf("Trigger binary from server result: %+v\n", result)

	time.Sleep(10 * time.Second)
	fmt.Println("Waiting for 10 seconds...")

	fmt.Println("\n=== WebSocket Binary Example Complete ===")
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
		Action:               tea.String("WebsocketAwapDemoApi"),
		Version:              tea.String("2022-02-02"),
		Protocol:             tea.String("wss"),
		Method:               tea.String("GET"),
		Pathname:             tea.String("/ws/awap-demo-api"),
		AuthType:             tea.String("AK"),
		WebsocketSubProtocol: tea.String("awap"),
	}

	request := &openapiClient.OpenApiRequest{}

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
	result, err := apiClient.CallApi(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	websocketObj := result["websocketClient"].(*websocketutils.WebSocketClient)
	defer websocketObj.Close()

	fmt.Println("\n1. Sending AWAP message...")

	// 方法 1: 使用 SendAwapRequest 发送请求消息（推荐）
	err = websocketObj.SendRawAwapTextMessageWithId(dara.AwapMessageType("UpstreamTextEvent"),
		"msg-001",
		map[string]interface{}{"action": "test", "data": "Hello WebSocket!"})
	if err != nil {
		log.Printf("Failed to send AWAP request: %v", err)
	} else {
		fmt.Printf("AWAP message sent successfully, type: %s\n", dara.AwapMessageType("UpstreamTextEvent"))
	}

	// 方法 2: 手动构建 AWAP 消息
	awapMsg := websocketutils.NewAwapMessage(dara.AwapMessageType("UpstreamTextEvent"),
		"msg-002",
		map[string]interface{}{"message": "Hello WebSocket!"},
	)
	err = websocketObj.SendAwapTextMessage(awapMsg)
	if err != nil {
		log.Printf("Failed to send AWAP message: %v", err)
	} else {
		fmt.Printf("AWAP message sent successfully, type: %s\n", dara.AwapMessageType("UpstreamTextEvent"))
	}

	// 方法 3: binary 信息
	awapMsgBinary := websocketutils.NewAwapMessage(dara.AwapMessageType("UpstreamBinaryEvent"),
		"msg-003",
		[]byte("Hello WebSocket!"),
	)
	err = websocketObj.SendAwapBinaryMessage(awapMsgBinary)
	if err != nil {
		log.Printf("Failed to send AWAP message: %v", err)
	} else {
		fmt.Printf("AWAP message sent successfully, type: %s\n", dara.AwapMessageType("UpstreamBinaryEvent"))
	}

	time.Sleep(1 * time.Second)

	// 方法 4: 发送 AckRequiredTextEvent 类型的消息（等待响应）
	fmt.Println("\n4. Sending AckRequiredTextEvent message and waiting for ACK...")
	ackResponse, err := websocketObj.SendRawAwapRequestWithAck(
		"msg-004",
		map[string]interface{}{"action": "ackRequiredTest", "data": "This message requires acknowledgment", "timestamp": time.Now().Unix()},
		30*time.Second,
	)
	if err != nil {
		log.Printf("❌ Failed to send AckRequiredTextEvent or timed out waiting for response: %v", err)
	} else {
		fmt.Printf("✅ Received acknowledgment:\n")
		fmt.Printf("  Response Type: %s\n", ackResponse.Type)
		if ackResponse.Headers != nil {
			if ackID, ok := ackResponse.Headers["ack-id"]; ok {
				fmt.Printf("  Ack-ID: %s\n", ackID)
			}
		}
		if ackResponse.Payload != nil {
			payloadJSON, _ := json.Marshal(ackResponse.Payload)
			fmt.Printf("  Payload: %s\n", string(payloadJSON))
		}
	}

	// Wait for other responses
	time.Sleep(2 * time.Second)

	fmt.Println("\n=== AWAP Example Complete ===")
	return nil
}

type GeneralHandler struct {
	dara.AbstractGeneralWebSocketHandler
}

func (h *GeneralHandler) AfterConnectionEstablished(session *dara.WebSocketSessionInfo) error {
	fmt.Println("✓ [CLI General] Connected to General WebSocket server")
	printSessionInfo(session)
	return nil
}

func (h *GeneralHandler) HandleRawMessage(session *dara.WebSocketSessionInfo, message *dara.WebSocketMessage) error {
	// Parse and handle General messages directly
	if message.Type == dara.WebSocketMessageTypeText {
		// Parse as General text message
		generalMsg, err := dara.ParseGeneralMessage(message)
		if err != nil {
			fmt.Printf("[CLI General] Failed to parse General message: %v\n", err)
			return err
		}

		fmt.Printf("[CLI General] Received text message: %+v\n", message)
		fmt.Printf("[CLI General] Received text message: %s\n", generalMsg.Body)

	} else if message.Type == dara.WebSocketMessageTypeBinary {
		// Handle as binary message
		fmt.Printf("[CLI General] Received binary message\n")
		fmt.Printf("[CLI General] Received binary message: %s\n", message.Payload)
	}

	return nil
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

// NoHandleAwapMessageHandler 是一个重写 HandleRawMessage 的 handler
type NoHandleAwapMessageHandler struct {
	dara.AbstractAwapWebSocketHandler
	messageCount int
	mu           sync.Mutex
}

func (h *NoHandleAwapMessageHandler) AfterConnectionEstablished(session *dara.WebSocketSessionInfo) error {
	fmt.Println("✓ CLI NoHandleAwapMessage - Connected to WebSocket server")
	fmt.Println("  Note: This handler does NOT override HandleAwapMessage")
	fmt.Println("  It overrides HandleRawMessage instead")
	printSessionInfo(session)
	return nil
}

func (h *NoHandleAwapMessageHandler) HandleRawMessage(session *dara.WebSocketSessionInfo, message *dara.WebSocketMessage) error {
	h.mu.Lock()
	h.messageCount++
	count := h.messageCount
	h.mu.Unlock()

	fmt.Printf("📨 CLI NoHandleAwapMessage - HandleRawMessage called (#%d):\n", count)
	fmt.Printf("  Type: %+v\n", message.Type)
	fmt.Printf("  Headers: %+v\n", message.Headers)
	fmt.Printf("  Payload: %s\n", string(message.Payload))
	printSessionInfo(session)
	return nil
}

func (h *NoHandleAwapMessageHandler) HandleError(session *dara.WebSocketSessionInfo, err error) error {
	fmt.Printf("❌ CLI NoHandleAwapMessage - Error: %v\n", err)
	printSessionInfo(session)
	return nil
}

func (h *NoHandleAwapMessageHandler) AfterConnectionClosed(session *dara.WebSocketSessionInfo, code int, reason string) error {
	h.mu.Lock()
	count := h.messageCount
	h.mu.Unlock()

	fmt.Printf("✗ CLI NoHandleAwapMessage - Connection closed (code: %d, reason: %s)\n", code, reason)
	fmt.Printf("  Total messages received: %d\n", count)
	printSessionInfo(session)
	return nil
}

func testAwapWebSocketWithoutHandleAwapMessage(ctx *cli.Context, args []string) error {
	fmt.Println("\n=== AWAP WebSocket Test (HandleRawMessage Override) ===")
	fmt.Println("This test demonstrates AwapWebSocketHandler interface usage")
	fmt.Println("The handler does NOT override HandleAwapMessage, but override HandleRawMessage")

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
		Action:               tea.String("WebsocketAwapDemoApi"),
		Version:              tea.String("2022-02-02"),
		Protocol:             tea.String("wss"),
		Method:               tea.String("GET"),
		Pathname:             tea.String("/ws/awap-demo-api"),
		AuthType:             tea.String("AK"),
		WebsocketSubProtocol: tea.String("awap"),
	}

	request := &openapiClient.OpenApiRequest{}

	runtime := &dara.RuntimeOptions{
		ReadTimeout:                dara.Int(60000),               // 60 seconds
		ConnectTimeout:             dara.Int(30000),               // 30 seconds
		WebSocketPingInterval:      dara.Int(30000),               // 30秒心跳
		WebSocketHandshakeTimeout:  dara.Int(30000),               // 30秒握手超时
		WebSocketWriteTimeout:      dara.Int(30000),               // 30秒写入超时
		WebSocketEnableReconnect:   dara.Bool(true),               // 启用重连
		WebSocketMaxReconnectTimes: dara.Int(5),                   // 最多重连5次
		WebSocketHandler:           &NoHandleAwapMessageHandler{}, // 使用不重写 HandleRawMessage 的 handler
	}

	fmt.Println("Connecting...")
	result, err := apiClient.CallApi(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	websocketObj := result["websocketClient"].(*websocketutils.WebSocketClient)
	defer websocketObj.Close()

	fmt.Println("\nSending AWAP messages...")

	// 方法 1: 使用 SendAwapRequest 发送请求消息
	fmt.Println("1. Sending AWAP request message...")
	err = websocketObj.SendRawAwapTextMessageWithId(dara.AwapMessageType("UpstreamTextEvent"),
		"msg-no-handleraw-001",
		map[string]interface{}{"action": "test", "data": "This handler does NOT override HandleRawMessage"},
	)
	if err != nil {
		log.Printf("Failed to send AWAP request: %v", err)
	} else {
		fmt.Printf("AWAP message sent successfully, type: %s\n", dara.AwapMessageType("UpstreamTextEvent"))
	}

	time.Sleep(1 * time.Second)

	// 方法 2: binary 信息
	awapMsgBinary := websocketutils.NewAwapMessage(dara.AwapMessageType("UpstreamBinaryEvent"),
		"msg-003",
		[]byte("Hello WebSocket!"),
	)
	err = websocketObj.SendAwapBinaryMessage(awapMsgBinary)
	if err != nil {
		log.Printf("Failed to send AWAP message: %v", err)
	} else {
		fmt.Printf("AWAP message sent successfully, type: %s\n", dara.AwapMessageType("UpstreamBinaryEvent"))
	}

	time.Sleep(1 * time.Second)

	// Wait for other responses
	fmt.Println("\nWaiting for other server responses...")
	time.Sleep(3 * time.Second)

	fmt.Println("\n=== NoHandleAwapMessageHandler Test Complete ===")
	fmt.Println("This test demonstrates that AwapWebSocketHandler interface")
	fmt.Println("is used when handler overrides HandleRawMessage")
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
		Action:               tea.String("WebsocketGeneralDemoApi"),
		Version:              tea.String("2022-02-02"),
		Protocol:             tea.String("wss"),
		Method:               tea.String("GET"),
		Pathname:             tea.String("/ws/general-demo-api"),
		AuthType:             tea.String("AK"),
		WebsocketSubProtocol: tea.String("general"),
	}

	request := &openapiClient.OpenApiRequest{}

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
	result, err := apiClient.CallApi(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	websocketObj := result["websocketClient"].(*websocketutils.WebSocketClient)
	defer websocketObj.Close()

	fmt.Println("\nSending General messages...")

	// 方法 1: 发送文本消息
	fmt.Println("1. Sending General text message...")
	err = websocketObj.SendGeneralTextMessage("Hello General WebSocket!")
	if err != nil {
		log.Printf("Failed to send General text message: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 方法 2: 发送 JSON 消息
	fmt.Println("2. Sending General JSON message...")
	jsonData, _ := json.Marshal(map[string]interface{}{
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
	})
	err = websocketObj.SendGeneralTextMessage(string(jsonData))
	if err != nil {
		log.Printf("Failed to send General JSON message: %v", err)
	}

	// 方法 3: 发送二进制消息
	fmt.Println("3. Sending General binary message...")
	binaryData := []byte("Binary General Data")
	err = websocketObj.SendGeneralBinaryMessage(binaryData)
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
	seq, err := strconv.ParseInt(message.Headers["seq"], 10, 64)
	if err != nil {
		return err
	}
	fmt.Printf("📨 CLI Sequential Test - HandleAwapMessage called: seq=%d, type=%s, id=%s\n",
		seq, message.Type, message.ID)
	// Print message payload for debugging
	if message.Payload != nil {
		payloadJSON, _ := json.Marshal(message.Payload)
		fmt.Printf("  Payload: %s\n", string(payloadJSON))
	}

	// Track sequence number
	h.receivedSeq = append(h.receivedSeq, seq)
	fmt.Printf("  Progress: %d/%d messages received\n", len(h.receivedSeq), h.expectedCount)

	// Check if we've received all expected messages
	if len(h.receivedSeq) > h.expectedCount {
		if h.done != nil {
			close(h.done)
			h.done = nil // Prevent double close
		}
	}

	return nil
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
		Action:               tea.String("WebsocketAwapDemoApi"),
		Version:              tea.String("2022-02-02"),
		Protocol:             tea.String("wss"),
		Method:               tea.String("GET"),
		Pathname:             tea.String("/ws/awap-demo-api"),
		AuthType:             tea.String("AK"),
		WebsocketSubProtocol: tea.String("awap"),
	}

	request := &openapiClient.OpenApiRequest{
		Query: map[string]*string{
			"delay":           tea.String("3000"), // Delay between messages (ms)
			"batchSendMsgCnt": tea.String("10"),   // Number of messages to send
		},
	}

	expectedCount := 10
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

	result, err := apiClient.CallApi(params, request, runtime)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	websocketObj := result["websocketClient"].(*websocketutils.WebSocketClient)
	defer websocketObj.Close()

	fmt.Println("Connection established, waiting for server to start sending batch messages...")
	time.Sleep(2 * time.Second)

	fmt.Println("\nSending request to trigger batch message sending...")
	err = websocketObj.SendRawAwapTextMessageWithId(dara.AwapMessageType("UpstreamTextEvent"),
		"sequential-test-001",
		map[string]interface{}{"action": "batchSend", "test": "sequential"})
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
