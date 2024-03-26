// Code generated from JsonQuery.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser

import (
	"fmt"
  	"sync"
	"unicode"
	"github.com/antlr4-go/antlr/v4"
)
// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter


type JsonQueryLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames []string
	// TODO: EOF string
}

var JsonQueryLexerLexerStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  ChannelNames           []string
  ModeNames              []string
  LiteralNames           []string
  SymbolicNames          []string
  RuleNames              []string
  PredictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func jsonquerylexerLexerInit() {
  staticData := &JsonQueryLexerLexerStaticData
  staticData.ChannelNames = []string{
    "DEFAULT_TOKEN_CHANNEL", "HIDDEN",
  }
  staticData.ModeNames = []string{
    "DEFAULT_MODE",
  }
  staticData.LiteralNames = []string{
    "", "'('", "')'", "'pr'", "'.'", "'-'", "'['", "']'", "", "", "", "'null'", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "'\\n'",
  }
  staticData.SymbolicNames = []string{
    "", "", "", "", "", "", "", "", "NOT", "LOGICAL_OPERATOR", "BOOLEAN", 
    "NULL", "IN", "EQ", "NE", "GT", "LT", "GE", "LE", "CO", "SW", "EW", 
    "MW", "ATTRNAME", "VERSION", "STRING", "DOUBLE", "INT", "EXP", "NEWLINE", 
    "COMMA", "SP",
  }
  staticData.RuleNames = []string{
    "T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "NOT", "LOGICAL_OPERATOR", 
    "BOOLEAN", "NULL", "IN", "EQ", "NE", "GT", "LT", "GE", "LE", "CO", "SW", 
    "EW", "MW", "ATTRNAME", "ATTR_NAME_CHAR", "DIGIT", "ALPHA", "VERSION", 
    "STRING", "ESC", "UNICODE", "HEX", "DOUBLE", "INT", "EXP", "NEWLINE", 
    "COMMA", "SP",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 0, 31, 287, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 
	4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 
	10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 
	7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 
	20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 
	2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 
	31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 
	7, 36, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 
	1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 97, 8, 
	7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 104, 8, 8, 1, 9, 1, 9, 1, 9, 1, 
	9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 115, 8, 9, 1, 10, 1, 10, 1, 10, 
	1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 126, 8, 11, 1, 12, 1, 
	12, 1, 12, 1, 12, 1, 12, 1, 12, 3, 12, 134, 8, 12, 1, 13, 1, 13, 1, 13, 
	1, 13, 1, 13, 1, 13, 3, 13, 142, 8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 
	14, 3, 14, 149, 8, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 156, 8, 
	15, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 164, 8, 16, 1, 17, 
	1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 3, 17, 172, 8, 17, 1, 18, 1, 18, 1, 
	18, 1, 18, 3, 18, 178, 8, 18, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 184, 8, 
	19, 1, 20, 1, 20, 1, 20, 1, 20, 3, 20, 190, 8, 20, 1, 21, 1, 21, 1, 21, 
	1, 21, 3, 21, 196, 8, 21, 1, 22, 1, 22, 5, 22, 200, 8, 22, 10, 22, 12, 
	22, 203, 9, 22, 1, 23, 1, 23, 1, 23, 3, 23, 208, 8, 23, 1, 24, 1, 24, 1, 
	25, 1, 25, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 27, 1, 27, 1, 27, 
	5, 27, 223, 8, 27, 10, 27, 12, 27, 226, 9, 27, 1, 27, 1, 27, 1, 28, 1, 
	28, 1, 28, 3, 28, 233, 8, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 
	1, 30, 1, 30, 1, 31, 3, 31, 244, 8, 31, 1, 31, 1, 31, 1, 31, 4, 31, 249, 
	8, 31, 11, 31, 12, 31, 250, 1, 31, 3, 31, 254, 8, 31, 1, 32, 1, 32, 1, 
	32, 5, 32, 259, 8, 32, 10, 32, 12, 32, 262, 9, 32, 3, 32, 264, 8, 32, 1, 
	33, 1, 33, 3, 33, 268, 8, 33, 1, 33, 1, 33, 1, 34, 1, 34, 1, 35, 1, 35, 
	5, 35, 276, 8, 35, 10, 35, 12, 35, 279, 9, 35, 1, 36, 1, 36, 5, 36, 283, 
	8, 36, 10, 36, 12, 36, 286, 9, 36, 0, 0, 37, 1, 1, 3, 2, 5, 3, 7, 4, 9, 
	5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 
	29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39, 20, 41, 21, 43, 22, 45, 23, 
	47, 0, 49, 0, 51, 0, 53, 24, 55, 25, 57, 0, 59, 0, 61, 0, 63, 26, 65, 27, 
	67, 28, 69, 29, 71, 30, 73, 31, 1, 0, 9, 3, 0, 45, 45, 58, 58, 95, 95, 
	2, 0, 65, 90, 97, 122, 2, 0, 34, 34, 92, 92, 8, 0, 34, 34, 47, 47, 92, 
	92, 98, 98, 102, 102, 110, 110, 114, 114, 116, 116, 3, 0, 48, 57, 65, 70, 
	97, 102, 1, 0, 48, 57, 1, 0, 49, 57, 2, 0, 69, 69, 101, 101, 2, 0, 43, 
	43, 45, 45, 314, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 
	0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 
	0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 
	0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 
	0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 
	1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1, 0, 0, 0, 0, 
	45, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 
	0, 65, 1, 0, 0, 0, 0, 67, 1, 0, 0, 0, 0, 69, 1, 0, 0, 0, 0, 71, 1, 0, 0, 
	0, 0, 73, 1, 0, 0, 0, 1, 75, 1, 0, 0, 0, 3, 77, 1, 0, 0, 0, 5, 79, 1, 0, 
	0, 0, 7, 82, 1, 0, 0, 0, 9, 84, 1, 0, 0, 0, 11, 86, 1, 0, 0, 0, 13, 88, 
	1, 0, 0, 0, 15, 96, 1, 0, 0, 0, 17, 103, 1, 0, 0, 0, 19, 114, 1, 0, 0, 
	0, 21, 116, 1, 0, 0, 0, 23, 125, 1, 0, 0, 0, 25, 133, 1, 0, 0, 0, 27, 141, 
	1, 0, 0, 0, 29, 148, 1, 0, 0, 0, 31, 155, 1, 0, 0, 0, 33, 163, 1, 0, 0, 
	0, 35, 171, 1, 0, 0, 0, 37, 177, 1, 0, 0, 0, 39, 183, 1, 0, 0, 0, 41, 189, 
	1, 0, 0, 0, 43, 195, 1, 0, 0, 0, 45, 197, 1, 0, 0, 0, 47, 207, 1, 0, 0, 
	0, 49, 209, 1, 0, 0, 0, 51, 211, 1, 0, 0, 0, 53, 213, 1, 0, 0, 0, 55, 219, 
	1, 0, 0, 0, 57, 229, 1, 0, 0, 0, 59, 234, 1, 0, 0, 0, 61, 240, 1, 0, 0, 
	0, 63, 243, 1, 0, 0, 0, 65, 263, 1, 0, 0, 0, 67, 265, 1, 0, 0, 0, 69, 271, 
	1, 0, 0, 0, 71, 273, 1, 0, 0, 0, 73, 280, 1, 0, 0, 0, 75, 76, 5, 40, 0, 
	0, 76, 2, 1, 0, 0, 0, 77, 78, 5, 41, 0, 0, 78, 4, 1, 0, 0, 0, 79, 80, 5, 
	112, 0, 0, 80, 81, 5, 114, 0, 0, 81, 6, 1, 0, 0, 0, 82, 83, 5, 46, 0, 0, 
	83, 8, 1, 0, 0, 0, 84, 85, 5, 45, 0, 0, 85, 10, 1, 0, 0, 0, 86, 87, 5, 
	91, 0, 0, 87, 12, 1, 0, 0, 0, 88, 89, 5, 93, 0, 0, 89, 14, 1, 0, 0, 0, 
	90, 91, 5, 110, 0, 0, 91, 92, 5, 111, 0, 0, 92, 97, 5, 116, 0, 0, 93, 94, 
	5, 78, 0, 0, 94, 95, 5, 79, 0, 0, 95, 97, 5, 84, 0, 0, 96, 90, 1, 0, 0, 
	0, 96, 93, 1, 0, 0, 0, 97, 16, 1, 0, 0, 0, 98, 99, 5, 97, 0, 0, 99, 100, 
	5, 110, 0, 0, 100, 104, 5, 100, 0, 0, 101, 102, 5, 111, 0, 0, 102, 104, 
	5, 114, 0, 0, 103, 98, 1, 0, 0, 0, 103, 101, 1, 0, 0, 0, 104, 18, 1, 0, 
	0, 0, 105, 106, 5, 116, 0, 0, 106, 107, 5, 114, 0, 0, 107, 108, 5, 117, 
	0, 0, 108, 115, 5, 101, 0, 0, 109, 110, 5, 102, 0, 0, 110, 111, 5, 97, 
	0, 0, 111, 112, 5, 108, 0, 0, 112, 113, 5, 115, 0, 0, 113, 115, 5, 101, 
	0, 0, 114, 105, 1, 0, 0, 0, 114, 109, 1, 0, 0, 0, 115, 20, 1, 0, 0, 0, 
	116, 117, 5, 110, 0, 0, 117, 118, 5, 117, 0, 0, 118, 119, 5, 108, 0, 0, 
	119, 120, 5, 108, 0, 0, 120, 22, 1, 0, 0, 0, 121, 122, 5, 73, 0, 0, 122, 
	126, 5, 78, 0, 0, 123, 124, 5, 105, 0, 0, 124, 126, 5, 110, 0, 0, 125, 
	121, 1, 0, 0, 0, 125, 123, 1, 0, 0, 0, 126, 24, 1, 0, 0, 0, 127, 128, 5, 
	101, 0, 0, 128, 134, 5, 113, 0, 0, 129, 130, 5, 69, 0, 0, 130, 134, 5, 
	81, 0, 0, 131, 132, 5, 61, 0, 0, 132, 134, 5, 61, 0, 0, 133, 127, 1, 0, 
	0, 0, 133, 129, 1, 0, 0, 0, 133, 131, 1, 0, 0, 0, 134, 26, 1, 0, 0, 0, 
	135, 136, 5, 110, 0, 0, 136, 142, 5, 101, 0, 0, 137, 138, 5, 78, 0, 0, 
	138, 142, 5, 69, 0, 0, 139, 140, 5, 33, 0, 0, 140, 142, 5, 61, 0, 0, 141, 
	135, 1, 0, 0, 0, 141, 137, 1, 0, 0, 0, 141, 139, 1, 0, 0, 0, 142, 28, 1, 
	0, 0, 0, 143, 144, 5, 103, 0, 0, 144, 149, 5, 116, 0, 0, 145, 146, 5, 71, 
	0, 0, 146, 149, 5, 84, 0, 0, 147, 149, 5, 62, 0, 0, 148, 143, 1, 0, 0, 
	0, 148, 145, 1, 0, 0, 0, 148, 147, 1, 0, 0, 0, 149, 30, 1, 0, 0, 0, 150, 
	151, 5, 108, 0, 0, 151, 156, 5, 116, 0, 0, 152, 153, 5, 76, 0, 0, 153, 
	156, 5, 84, 0, 0, 154, 156, 5, 60, 0, 0, 155, 150, 1, 0, 0, 0, 155, 152, 
	1, 0, 0, 0, 155, 154, 1, 0, 0, 0, 156, 32, 1, 0, 0, 0, 157, 158, 5, 103, 
	0, 0, 158, 164, 5, 101, 0, 0, 159, 160, 5, 71, 0, 0, 160, 164, 5, 69, 0, 
	0, 161, 162, 5, 62, 0, 0, 162, 164, 5, 61, 0, 0, 163, 157, 1, 0, 0, 0, 
	163, 159, 1, 0, 0, 0, 163, 161, 1, 0, 0, 0, 164, 34, 1, 0, 0, 0, 165, 166, 
	5, 108, 0, 0, 166, 172, 5, 101, 0, 0, 167, 168, 5, 76, 0, 0, 168, 172, 
	5, 69, 0, 0, 169, 170, 5, 60, 0, 0, 170, 172, 5, 61, 0, 0, 171, 165, 1, 
	0, 0, 0, 171, 167, 1, 0, 0, 0, 171, 169, 1, 0, 0, 0, 172, 36, 1, 0, 0, 
	0, 173, 174, 5, 99, 0, 0, 174, 178, 5, 111, 0, 0, 175, 176, 5, 67, 0, 0, 
	176, 178, 5, 79, 0, 0, 177, 173, 1, 0, 0, 0, 177, 175, 1, 0, 0, 0, 178, 
	38, 1, 0, 0, 0, 179, 180, 5, 115, 0, 0, 180, 184, 5, 119, 0, 0, 181, 182, 
	5, 83, 0, 0, 182, 184, 5, 87, 0, 0, 183, 179, 1, 0, 0, 0, 183, 181, 1, 
	0, 0, 0, 184, 40, 1, 0, 0, 0, 185, 186, 5, 101, 0, 0, 186, 190, 5, 119, 
	0, 0, 187, 188, 5, 69, 0, 0, 188, 190, 5, 87, 0, 0, 189, 185, 1, 0, 0, 
	0, 189, 187, 1, 0, 0, 0, 190, 42, 1, 0, 0, 0, 191, 192, 5, 109, 0, 0, 192, 
	196, 5, 119, 0, 0, 193, 194, 5, 77, 0, 0, 194, 196, 5, 87, 0, 0, 195, 191, 
	1, 0, 0, 0, 195, 193, 1, 0, 0, 0, 196, 44, 1, 0, 0, 0, 197, 201, 3, 51, 
	25, 0, 198, 200, 3, 47, 23, 0, 199, 198, 1, 0, 0, 0, 200, 203, 1, 0, 0, 
	0, 201, 199, 1, 0, 0, 0, 201, 202, 1, 0, 0, 0, 202, 46, 1, 0, 0, 0, 203, 
	201, 1, 0, 0, 0, 204, 208, 7, 0, 0, 0, 205, 208, 3, 49, 24, 0, 206, 208, 
	3, 51, 25, 0, 207, 204, 1, 0, 0, 0, 207, 205, 1, 0, 0, 0, 207, 206, 1, 
	0, 0, 0, 208, 48, 1, 0, 0, 0, 209, 210, 2, 48, 57, 0, 210, 50, 1, 0, 0, 
	0, 211, 212, 7, 1, 0, 0, 212, 52, 1, 0, 0, 0, 213, 214, 3, 65, 32, 0, 214, 
	215, 5, 46, 0, 0, 215, 216, 3, 65, 32, 0, 216, 217, 5, 46, 0, 0, 217, 218, 
	3, 65, 32, 0, 218, 54, 1, 0, 0, 0, 219, 224, 5, 34, 0, 0, 220, 223, 3, 
	57, 28, 0, 221, 223, 8, 2, 0, 0, 222, 220, 1, 0, 0, 0, 222, 221, 1, 0, 
	0, 0, 223, 226, 1, 0, 0, 0, 224, 222, 1, 0, 0, 0, 224, 225, 1, 0, 0, 0, 
	225, 227, 1, 0, 0, 0, 226, 224, 1, 0, 0, 0, 227, 228, 5, 34, 0, 0, 228, 
	56, 1, 0, 0, 0, 229, 232, 5, 92, 0, 0, 230, 233, 7, 3, 0, 0, 231, 233, 
	3, 59, 29, 0, 232, 230, 1, 0, 0, 0, 232, 231, 1, 0, 0, 0, 233, 58, 1, 0, 
	0, 0, 234, 235, 5, 117, 0, 0, 235, 236, 3, 61, 30, 0, 236, 237, 3, 61, 
	30, 0, 237, 238, 3, 61, 30, 0, 238, 239, 3, 61, 30, 0, 239, 60, 1, 0, 0, 
	0, 240, 241, 7, 4, 0, 0, 241, 62, 1, 0, 0, 0, 242, 244, 5, 45, 0, 0, 243, 
	242, 1, 0, 0, 0, 243, 244, 1, 0, 0, 0, 244, 245, 1, 0, 0, 0, 245, 246, 
	3, 65, 32, 0, 246, 248, 5, 46, 0, 0, 247, 249, 7, 5, 0, 0, 248, 247, 1, 
	0, 0, 0, 249, 250, 1, 0, 0, 0, 250, 248, 1, 0, 0, 0, 250, 251, 1, 0, 0, 
	0, 251, 253, 1, 0, 0, 0, 252, 254, 3, 67, 33, 0, 253, 252, 1, 0, 0, 0, 
	253, 254, 1, 0, 0, 0, 254, 64, 1, 0, 0, 0, 255, 264, 5, 48, 0, 0, 256, 
	260, 7, 6, 0, 0, 257, 259, 7, 5, 0, 0, 258, 257, 1, 0, 0, 0, 259, 262, 
	1, 0, 0, 0, 260, 258, 1, 0, 0, 0, 260, 261, 1, 0, 0, 0, 261, 264, 1, 0, 
	0, 0, 262, 260, 1, 0, 0, 0, 263, 255, 1, 0, 0, 0, 263, 256, 1, 0, 0, 0, 
	264, 66, 1, 0, 0, 0, 265, 267, 7, 7, 0, 0, 266, 268, 7, 8, 0, 0, 267, 266, 
	1, 0, 0, 0, 267, 268, 1, 0, 0, 0, 268, 269, 1, 0, 0, 0, 269, 270, 3, 65, 
	32, 0, 270, 68, 1, 0, 0, 0, 271, 272, 5, 10, 0, 0, 272, 70, 1, 0, 0, 0, 
	273, 277, 5, 44, 0, 0, 274, 276, 5, 32, 0, 0, 275, 274, 1, 0, 0, 0, 276, 
	279, 1, 0, 0, 0, 277, 275, 1, 0, 0, 0, 277, 278, 1, 0, 0, 0, 278, 72, 1, 
	0, 0, 0, 279, 277, 1, 0, 0, 0, 280, 284, 5, 32, 0, 0, 281, 283, 3, 69, 
	34, 0, 282, 281, 1, 0, 0, 0, 283, 286, 1, 0, 0, 0, 284, 282, 1, 0, 0, 0, 
	284, 285, 1, 0, 0, 0, 285, 74, 1, 0, 0, 0, 286, 284, 1, 0, 0, 0, 28, 0, 
	96, 103, 114, 125, 133, 141, 148, 155, 163, 171, 177, 183, 189, 195, 201, 
	207, 222, 224, 232, 243, 250, 253, 260, 263, 267, 277, 284, 0,
}
  deserializer := antlr.NewATNDeserializer(nil)
  staticData.atn = deserializer.Deserialize(staticData.serializedATN)
  atn := staticData.atn
  staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
  decisionToDFA := staticData.decisionToDFA
  for index, state := range atn.DecisionToState {
    decisionToDFA[index] = antlr.NewDFA(state, index)
  }
}

// JsonQueryLexerInit initializes any static state used to implement JsonQueryLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewJsonQueryLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func JsonQueryLexerInit() {
  staticData := &JsonQueryLexerLexerStaticData
  staticData.once.Do(jsonquerylexerLexerInit)
}

// NewJsonQueryLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewJsonQueryLexer(input antlr.CharStream) *JsonQueryLexer {
  JsonQueryLexerInit()
	l := new(JsonQueryLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
  staticData := &JsonQueryLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "JsonQuery.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// JsonQueryLexer tokens.
const (
	JsonQueryLexerT__0 = 1
	JsonQueryLexerT__1 = 2
	JsonQueryLexerT__2 = 3
	JsonQueryLexerT__3 = 4
	JsonQueryLexerT__4 = 5
	JsonQueryLexerT__5 = 6
	JsonQueryLexerT__6 = 7
	JsonQueryLexerNOT = 8
	JsonQueryLexerLOGICAL_OPERATOR = 9
	JsonQueryLexerBOOLEAN = 10
	JsonQueryLexerNULL = 11
	JsonQueryLexerIN = 12
	JsonQueryLexerEQ = 13
	JsonQueryLexerNE = 14
	JsonQueryLexerGT = 15
	JsonQueryLexerLT = 16
	JsonQueryLexerGE = 17
	JsonQueryLexerLE = 18
	JsonQueryLexerCO = 19
	JsonQueryLexerSW = 20
	JsonQueryLexerEW = 21
	JsonQueryLexerMW = 22
	JsonQueryLexerATTRNAME = 23
	JsonQueryLexerVERSION = 24
	JsonQueryLexerSTRING = 25
	JsonQueryLexerDOUBLE = 26
	JsonQueryLexerINT = 27
	JsonQueryLexerEXP = 28
	JsonQueryLexerNEWLINE = 29
	JsonQueryLexerCOMMA = 30
	JsonQueryLexerSP = 31
)

