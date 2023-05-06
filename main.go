package main
 
import (
	"fmt"
	"log"
	"context"
	"net/http"
 
	"github.com/gorilla/websocket"

 "github.com/franciscoescher/goopenai"
)
 
// Initialize the ChatGPT client once and reuse it in each function call
var openaiClient *goopenai.Client

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func init() {
    // Replace with your API Key
    apiKey := "blah"

    // Initialize the GPT-3 client
    openaiClient = goopenai.NewClient(apiKey, "organization")

}

func GetResponse2(p string) string {
 r := goopenai.CreateCompletionsRequest{
  Model: "gpt-3.5-turbo",
  Messages: []goopenai.Message{
   {
    Role:    "user",
    Content: p,
   },
  },
  Temperature: 0.7,
 }

 completions, err := openaiClient.CreateCompletions(context.Background(), r)
 if err != nil {
  panic(err)
 }
 return completions.Choices[0].Text
}

// func GetResponse(client *openai.Client, ctx context.Context, quesiton string) string {
// 	req := openai.CompletionRequest{
// 		Model:     openai.GPT3TextDavinci001,
// 		MaxTokens: 300,
// 		Prompt:    quesiton,
// 		Stream:    true,
// 	}

// 	resp, err := client.CreateCompletionStream(ctx, req)
// 	if err != nil {
// 		return "CreateCompletionStream returned error"
// 	}
// 	defer resp.Close()

// 	counter := 0
// 	for {
// 		data, err := resp.Recv()
// 		if err != nil {
// 			if errors.Is(err, io.EOF) {
// 				break
// 			}
// 			return "Stream error"
// 		} else {
// 			counter++
// 			return data.Choices[0].Text

// 		}
// 	}
// 	if counter == 0 {
// 		return "Stream did not return any responses"
// 	}
// 	return ""
// }



 
func reader(conn *websocket.Conn) {
	    // Example usage of the generateTextFromPrompt function


	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))
		prompt := string(p)
		generatedText := prompt

		log.Println("response:"+ generatedText)
		p = []byte(fmt.Sprintf("Echo: %v", string(p)))
		if err := conn.WriteMessage(messageType, []byte(generatedText)); err != nil {
			log.Println(err)
			return
		}
 
	}
}
 
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}
 
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
 
	log.Println("Client Connected")
	p := []byte(fmt.Sprintf("%v", "Hi Client!"))
	err = ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		log.Println(err)
	}
 
	reader(ws)
	defer ws.Close()
}
 
func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}
 
func main() {
	fmt.Println("Hello World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}