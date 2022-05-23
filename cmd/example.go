package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/aga3000/go-brianfuck"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

const rot13 = `
-,+[                        
    -[                      
        >>++++[>++++++++<-]
                            
        <+<-[               
            >+>+>-[>>>]     
            <[[>+<-]>>+>]
            <<<<<-          
        ]                   
    ]>>>[-]+                
    >--[-[<->+++[-]]]<[     
        ++++++++++++<[      
                            
            >-[>+>>]        
            >[+[<+>-]>+>>]  
            <<<<<-          
        ]                   
        >>[<+>-]            
        >[                  
            -[              
                -<<[-]>>    
            ]<<[<<->>-]>>   
        ]<<[<<+>>-]         
    ]                       
    <[-]                    
    <.[-]                   
    <-,+                    
]                           
`

func main() {
	source := bytes.NewReader([]byte(rot13))
	reader := bytes.NewReader([]byte("HOLA"))
	writer := bufio.NewWriter(os.Stdout)
	runner, err := brainfuck.NewInterpreterRunner(
		reader, writer,
		brainfuck.WithUnknownCharPolicy(brainfuck.IgnoreUnknownCharsPolicy),
	)
	if err != nil {
		log.Fatalf("failed to execute example: %+v", err)
	}
	for charPos := 0; true; charPos++ {
		ch, _, codeReaderErr := source.ReadRune()
		if errors.Is(codeReaderErr, io.EOF) {
			break
		}
		err = runner.Execute(brainfuck.CommandChar(ch))
		if err != nil {
			if errors.Is(err, brainfuck.EOR) {
				break
			}
			log.Fatalf("execution failed at position %d on command %c: %+v", charPos, ch, err)
		}
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("failed to flush writer buffer")
	}
	fmt.Println()
}
