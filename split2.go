package main

import (
	"fmt"
	"regexp"
	"strings"
)

func processRawMessage(message string, maxSize int) []string {
	lines := strings.Split(message, "\n")
	var chunks []string
	var currentChunk strings.Builder
	codeBlock := false
	for _, line := range lines {
		line = strings.TrimRight(line, " ")
		if line == "```" {
			codeBlock = !codeBlock
		}
		if len(line) > maxSize {
			i := 0
			for i < len(line) {
				remainingLineSpace := maxSize - currentChunk.Len()
				if len(line[i:]) <= remainingLineSpace {
					currentChunk.WriteString(line[i:])
					i += len(line[i:])
				} else {
					currentChunk.WriteString(line[i : i+remainingLineSpace])
					chunks = append(chunks, currentChunk.String())
					currentChunk.Reset()
					i += remainingLineSpace
				}
				if codeBlock && currentChunk.Len() > 0 {
					currentChunk.WriteString("```")
				}
			}
		} else {
			if currentChunk.Len()+len(line)+1 > maxSize {
				if codeBlock {
					currentChunk.WriteString("```")
				}
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
				if codeBlock {
					currentChunk.WriteString("```")
				}
			}
			if currentChunk.Len() > 0 {
				currentChunk.WriteString("\n")
			}
			currentChunk.WriteString(line)
		}
	}
	if currentChunk.Len() > 0 {
		if codeBlock {
			currentChunk.WriteString("```")
		}
		chunks = append(chunks, currentChunk.String())
	}
	return chunks
}

// normalizeBackticks adds leading and trailing newlines to triple backticks (```)
// If the leading or trailing newline is already present, it doesn't duplicate it
func normalizeBackticks(input string) string {
	// Regular expression that matches a triple backtick with any character (or none) on either side
	re := regexp.MustCompile(".{0,1}```.{0,1}")

	// Replace function adds newlines before and after every triple backtick
	// If the newline is already present, it doesn't duplicate it
	return re.ReplaceAllStringFunc(input, func(s string) string {
		if s == "\n```\n" {
			return s
		}
		if strings.Count(s, "`") != 3 {
			return s
		}
		leading := string(s[0])
		trailing := string(s[len(s)-1])
		if leading == "`" {
			leading = ""
		} else if leading != "\n" {
			leading = leading + "\n"
		}
		if trailing == "`" {
			trailing = ""
		} else if trailing != "\n" {
			trailing = "\n" + trailing
		}
		return leading + "```" + trailing
	})
}

// Dummy main function to run some tests
func main() {
	message := `
This is a long test message with multiple lines and a @mention.

<c>
Some example code here. This block is intended to be large enough to ensure that the processRawMessage function will have to split it across multiple messages.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed auctor accumsan neque, non elementum nibh condimentum id. Sed volutpat nulla eu lorem vulputate, eu lacinia turpis condimentum. Donec vehicula ex sed nisl consequat volutpat. Praesent scelerisque cursus libero, a fringilla erat bibendum a. Aliquam at leo tellus. Morbi sagittis vehicula posuere. Donec a libero risus. Quisque aliquam lectus at mauris bibendum, at vehicula nisl elementum. Sed ultrices dui non ante sollicitudin, id fermentum nulla iaculis.
<c>

This is to test the processRawMessage function to ensure it splits the message correctly when necessary.
Let's include a couple more @mentions and even @here.

<c>
Another large code block to test. This one is also intended to be large enough to necessitate splitting across multiple messages.
Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Suspendisse ullamcorper aliquam velit, non viverra ante ullamcorper non. Duis sit amet ante auctor, facilisis orci id, dapibus justo. Sed volutpat erat ut metus elementum, a aliquet dui lobortis. Maecenas porttitor, sapien a eleifend sollicitudin, sem erat pretium mauris, vitae hendrerit est neque vitae nisl. Nullam quis augue eget nulla fringilla scelerisque a vitae est. Nullam volutpat risus ac mauris congue, in condimentum ex pellentesque.
<c>

End of message.
`
	message = strings.ReplaceAll(message, "<c>", "```")
	normal := normalizeBackticks(message)
	fmt.Printf("Original (size: %d):\n%s\n-------\n", len(normal), normal)
	messages := processRawMessage(message, 100)
	fmt.Println("Splits:")
	for _, msg := range messages {
		fmt.Printf("------- (size: %d)\n%s\n", len(msg), msg)
	}
}
