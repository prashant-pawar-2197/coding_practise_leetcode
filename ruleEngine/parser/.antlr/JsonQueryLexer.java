// Generated from d:\project\Syncup\GOlang\lib\git\rule-engine\parser\JsonQuery.g4 by ANTLR 4.8
import org.antlr.v4.runtime.Lexer;
import org.antlr.v4.runtime.CharStream;
import org.antlr.v4.runtime.Token;
import org.antlr.v4.runtime.TokenStream;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.misc.*;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class JsonQueryLexer extends Lexer {
	static { RuntimeMetaData.checkVersion("4.8", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		T__0=1, T__1=2, T__2=3, T__3=4, T__4=5, T__5=6, T__6=7, NOT=8, LOGICAL_OPERATOR=9, 
		BOOLEAN=10, NULL=11, IN=12, EQ=13, NE=14, GT=15, LT=16, GE=17, LE=18, 
		CO=19, SW=20, EW=21, MW=22, ATTRNAME=23, VERSION=24, STRING=25, DOUBLE=26, 
		INT=27, EXP=28, NEWLINE=29, COMMA=30, SP=31;
	public static String[] channelNames = {
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN"
	};

	public static String[] modeNames = {
		"DEFAULT_MODE"
	};

	private static String[] makeRuleNames() {
		return new String[] {
			"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "NOT", "LOGICAL_OPERATOR", 
			"BOOLEAN", "NULL", "IN", "EQ", "NE", "GT", "LT", "GE", "LE", "CO", "SW", 
			"EW", "MW", "ATTRNAME", "ATTR_NAME_CHAR", "DIGIT", "ALPHA", "VERSION", 
			"STRING", "ESC", "UNICODE", "HEX", "DOUBLE", "INT", "EXP", "NEWLINE", 
			"COMMA", "SP"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'('", "')'", "'pr'", "'.'", "'-'", "'['", "']'", null, null, null, 
			"'null'", null, null, null, null, null, null, null, null, null, null, 
			null, null, null, null, null, null, null, "'\n'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, null, null, null, null, null, null, null, "NOT", "LOGICAL_OPERATOR", 
			"BOOLEAN", "NULL", "IN", "EQ", "NE", "GT", "LT", "GE", "LE", "CO", "SW", 
			"EW", "MW", "ATTRNAME", "VERSION", "STRING", "DOUBLE", "INT", "EXP", 
			"NEWLINE", "COMMA", "SP"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}


	public JsonQueryLexer(CharStream input) {
		super(input);
		_interp = new LexerATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@Override
	public String getGrammarFileName() { return "JsonQuery.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public String[] getChannelNames() { return channelNames; }

	@Override
	public String[] getModeNames() { return modeNames; }

	@Override
	public ATN getATN() { return _ATN; }

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\2!\u0121\b\1\4\2\t"+
		"\2\4\3\t\3\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t\t\4\n\t\n\4\13"+
		"\t\13\4\f\t\f\4\r\t\r\4\16\t\16\4\17\t\17\4\20\t\20\4\21\t\21\4\22\t\22"+
		"\4\23\t\23\4\24\t\24\4\25\t\25\4\26\t\26\4\27\t\27\4\30\t\30\4\31\t\31"+
		"\4\32\t\32\4\33\t\33\4\34\t\34\4\35\t\35\4\36\t\36\4\37\t\37\4 \t \4!"+
		"\t!\4\"\t\"\4#\t#\4$\t$\4%\t%\4&\t&\3\2\3\2\3\3\3\3\3\4\3\4\3\4\3\5\3"+
		"\5\3\6\3\6\3\7\3\7\3\b\3\b\3\t\3\t\3\t\3\t\3\t\3\t\5\tc\n\t\3\n\3\n\3"+
		"\n\3\n\3\n\5\nj\n\n\3\13\3\13\3\13\3\13\3\13\3\13\3\13\3\13\3\13\5\13"+
		"u\n\13\3\f\3\f\3\f\3\f\3\f\3\r\3\r\3\r\3\r\5\r\u0080\n\r\3\16\3\16\3\16"+
		"\3\16\3\16\3\16\5\16\u0088\n\16\3\17\3\17\3\17\3\17\3\17\3\17\5\17\u0090"+
		"\n\17\3\20\3\20\3\20\3\20\3\20\5\20\u0097\n\20\3\21\3\21\3\21\3\21\3\21"+
		"\5\21\u009e\n\21\3\22\3\22\3\22\3\22\3\22\3\22\5\22\u00a6\n\22\3\23\3"+
		"\23\3\23\3\23\3\23\3\23\5\23\u00ae\n\23\3\24\3\24\3\24\3\24\5\24\u00b4"+
		"\n\24\3\25\3\25\3\25\3\25\5\25\u00ba\n\25\3\26\3\26\3\26\3\26\5\26\u00c0"+
		"\n\26\3\27\3\27\3\27\3\27\5\27\u00c6\n\27\3\30\3\30\7\30\u00ca\n\30\f"+
		"\30\16\30\u00cd\13\30\3\31\3\31\3\31\5\31\u00d2\n\31\3\32\3\32\3\33\3"+
		"\33\3\34\3\34\3\34\3\34\3\34\3\34\3\35\3\35\3\35\7\35\u00e1\n\35\f\35"+
		"\16\35\u00e4\13\35\3\35\3\35\3\36\3\36\3\36\5\36\u00eb\n\36\3\37\3\37"+
		"\3\37\3\37\3\37\3\37\3 \3 \3!\5!\u00f6\n!\3!\3!\3!\6!\u00fb\n!\r!\16!"+
		"\u00fc\3!\5!\u0100\n!\3\"\3\"\3\"\7\"\u0105\n\"\f\"\16\"\u0108\13\"\5"+
		"\"\u010a\n\"\3#\3#\5#\u010e\n#\3#\3#\3$\3$\3%\3%\7%\u0116\n%\f%\16%\u0119"+
		"\13%\3&\3&\7&\u011d\n&\f&\16&\u0120\13&\2\2\'\3\3\5\4\7\5\t\6\13\7\r\b"+
		"\17\t\21\n\23\13\25\f\27\r\31\16\33\17\35\20\37\21!\22#\23%\24\'\25)\26"+
		"+\27-\30/\31\61\2\63\2\65\2\67\329\33;\2=\2?\2A\34C\35E\36G\37I K!\3\2"+
		"\13\5\2//<<aa\4\2C\\c|\4\2$$^^\n\2$$\61\61^^ddhhppttvv\5\2\62;CHch\3\2"+
		"\62;\3\2\63;\4\2GGgg\4\2--//\2\u013c\2\3\3\2\2\2\2\5\3\2\2\2\2\7\3\2\2"+
		"\2\2\t\3\2\2\2\2\13\3\2\2\2\2\r\3\2\2\2\2\17\3\2\2\2\2\21\3\2\2\2\2\23"+
		"\3\2\2\2\2\25\3\2\2\2\2\27\3\2\2\2\2\31\3\2\2\2\2\33\3\2\2\2\2\35\3\2"+
		"\2\2\2\37\3\2\2\2\2!\3\2\2\2\2#\3\2\2\2\2%\3\2\2\2\2\'\3\2\2\2\2)\3\2"+
		"\2\2\2+\3\2\2\2\2-\3\2\2\2\2/\3\2\2\2\2\67\3\2\2\2\29\3\2\2\2\2A\3\2\2"+
		"\2\2C\3\2\2\2\2E\3\2\2\2\2G\3\2\2\2\2I\3\2\2\2\2K\3\2\2\2\3M\3\2\2\2\5"+
		"O\3\2\2\2\7Q\3\2\2\2\tT\3\2\2\2\13V\3\2\2\2\rX\3\2\2\2\17Z\3\2\2\2\21"+
		"b\3\2\2\2\23i\3\2\2\2\25t\3\2\2\2\27v\3\2\2\2\31\177\3\2\2\2\33\u0087"+
		"\3\2\2\2\35\u008f\3\2\2\2\37\u0096\3\2\2\2!\u009d\3\2\2\2#\u00a5\3\2\2"+
		"\2%\u00ad\3\2\2\2\'\u00b3\3\2\2\2)\u00b9\3\2\2\2+\u00bf\3\2\2\2-\u00c5"+
		"\3\2\2\2/\u00c7\3\2\2\2\61\u00d1\3\2\2\2\63\u00d3\3\2\2\2\65\u00d5\3\2"+
		"\2\2\67\u00d7\3\2\2\29\u00dd\3\2\2\2;\u00e7\3\2\2\2=\u00ec\3\2\2\2?\u00f2"+
		"\3\2\2\2A\u00f5\3\2\2\2C\u0109\3\2\2\2E\u010b\3\2\2\2G\u0111\3\2\2\2I"+
		"\u0113\3\2\2\2K\u011a\3\2\2\2MN\7*\2\2N\4\3\2\2\2OP\7+\2\2P\6\3\2\2\2"+
		"QR\7r\2\2RS\7t\2\2S\b\3\2\2\2TU\7\60\2\2U\n\3\2\2\2VW\7/\2\2W\f\3\2\2"+
		"\2XY\7]\2\2Y\16\3\2\2\2Z[\7_\2\2[\20\3\2\2\2\\]\7p\2\2]^\7q\2\2^c\7v\2"+
		"\2_`\7P\2\2`a\7Q\2\2ac\7V\2\2b\\\3\2\2\2b_\3\2\2\2c\22\3\2\2\2de\7c\2"+
		"\2ef\7p\2\2fj\7f\2\2gh\7q\2\2hj\7t\2\2id\3\2\2\2ig\3\2\2\2j\24\3\2\2\2"+
		"kl\7v\2\2lm\7t\2\2mn\7w\2\2nu\7g\2\2op\7h\2\2pq\7c\2\2qr\7n\2\2rs\7u\2"+
		"\2su\7g\2\2tk\3\2\2\2to\3\2\2\2u\26\3\2\2\2vw\7p\2\2wx\7w\2\2xy\7n\2\2"+
		"yz\7n\2\2z\30\3\2\2\2{|\7K\2\2|\u0080\7P\2\2}~\7k\2\2~\u0080\7p\2\2\177"+
		"{\3\2\2\2\177}\3\2\2\2\u0080\32\3\2\2\2\u0081\u0082\7g\2\2\u0082\u0088"+
		"\7s\2\2\u0083\u0084\7G\2\2\u0084\u0088\7S\2\2\u0085\u0086\7?\2\2\u0086"+
		"\u0088\7?\2\2\u0087\u0081\3\2\2\2\u0087\u0083\3\2\2\2\u0087\u0085\3\2"+
		"\2\2\u0088\34\3\2\2\2\u0089\u008a\7p\2\2\u008a\u0090\7g\2\2\u008b\u008c"+
		"\7P\2\2\u008c\u0090\7G\2\2\u008d\u008e\7#\2\2\u008e\u0090\7?\2\2\u008f"+
		"\u0089\3\2\2\2\u008f\u008b\3\2\2\2\u008f\u008d\3\2\2\2\u0090\36\3\2\2"+
		"\2\u0091\u0092\7i\2\2\u0092\u0097\7v\2\2\u0093\u0094\7I\2\2\u0094\u0097"+
		"\7V\2\2\u0095\u0097\7@\2\2\u0096\u0091\3\2\2\2\u0096\u0093\3\2\2\2\u0096"+
		"\u0095\3\2\2\2\u0097 \3\2\2\2\u0098\u0099\7n\2\2\u0099\u009e\7v\2\2\u009a"+
		"\u009b\7N\2\2\u009b\u009e\7V\2\2\u009c\u009e\7>\2\2\u009d\u0098\3\2\2"+
		"\2\u009d\u009a\3\2\2\2\u009d\u009c\3\2\2\2\u009e\"\3\2\2\2\u009f\u00a0"+
		"\7i\2\2\u00a0\u00a6\7g\2\2\u00a1\u00a2\7I\2\2\u00a2\u00a6\7G\2\2\u00a3"+
		"\u00a4\7@\2\2\u00a4\u00a6\7?\2\2\u00a5\u009f\3\2\2\2\u00a5\u00a1\3\2\2"+
		"\2\u00a5\u00a3\3\2\2\2\u00a6$\3\2\2\2\u00a7\u00a8\7n\2\2\u00a8\u00ae\7"+
		"g\2\2\u00a9\u00aa\7N\2\2\u00aa\u00ae\7G\2\2\u00ab\u00ac\7>\2\2\u00ac\u00ae"+
		"\7?\2\2\u00ad\u00a7\3\2\2\2\u00ad\u00a9\3\2\2\2\u00ad\u00ab\3\2\2\2\u00ae"+
		"&\3\2\2\2\u00af\u00b0\7e\2\2\u00b0\u00b4\7q\2\2\u00b1\u00b2\7E\2\2\u00b2"+
		"\u00b4\7Q\2\2\u00b3\u00af\3\2\2\2\u00b3\u00b1\3\2\2\2\u00b4(\3\2\2\2\u00b5"+
		"\u00b6\7u\2\2\u00b6\u00ba\7y\2\2\u00b7\u00b8\7U\2\2\u00b8\u00ba\7Y\2\2"+
		"\u00b9\u00b5\3\2\2\2\u00b9\u00b7\3\2\2\2\u00ba*\3\2\2\2\u00bb\u00bc\7"+
		"g\2\2\u00bc\u00c0\7y\2\2\u00bd\u00be\7G\2\2\u00be\u00c0\7Y\2\2\u00bf\u00bb"+
		"\3\2\2\2\u00bf\u00bd\3\2\2\2\u00c0,\3\2\2\2\u00c1\u00c2\7o\2\2\u00c2\u00c6"+
		"\7y\2\2\u00c3\u00c4\7O\2\2\u00c4\u00c6\7Y\2\2\u00c5\u00c1\3\2\2\2\u00c5"+
		"\u00c3\3\2\2\2\u00c6.\3\2\2\2\u00c7\u00cb\5\65\33\2\u00c8\u00ca\5\61\31"+
		"\2\u00c9\u00c8\3\2\2\2\u00ca\u00cd\3\2\2\2\u00cb\u00c9\3\2\2\2\u00cb\u00cc"+
		"\3\2\2\2\u00cc\60\3\2\2\2\u00cd\u00cb\3\2\2\2\u00ce\u00d2\t\2\2\2\u00cf"+
		"\u00d2\5\63\32\2\u00d0\u00d2\5\65\33\2\u00d1\u00ce\3\2\2\2\u00d1\u00cf"+
		"\3\2\2\2\u00d1\u00d0\3\2\2\2\u00d2\62\3\2\2\2\u00d3\u00d4\4\62;\2\u00d4"+
		"\64\3\2\2\2\u00d5\u00d6\t\3\2\2\u00d6\66\3\2\2\2\u00d7\u00d8\5C\"\2\u00d8"+
		"\u00d9\7\60\2\2\u00d9\u00da\5C\"\2\u00da\u00db\7\60\2\2\u00db\u00dc\5"+
		"C\"\2\u00dc8\3\2\2\2\u00dd\u00e2\7$\2\2\u00de\u00e1\5;\36\2\u00df\u00e1"+
		"\n\4\2\2\u00e0\u00de\3\2\2\2\u00e0\u00df\3\2\2\2\u00e1\u00e4\3\2\2\2\u00e2"+
		"\u00e0\3\2\2\2\u00e2\u00e3\3\2\2\2\u00e3\u00e5\3\2\2\2\u00e4\u00e2\3\2"+
		"\2\2\u00e5\u00e6\7$\2\2\u00e6:\3\2\2\2\u00e7\u00ea\7^\2\2\u00e8\u00eb"+
		"\t\5\2\2\u00e9\u00eb\5=\37\2\u00ea\u00e8\3\2\2\2\u00ea\u00e9\3\2\2\2\u00eb"+
		"<\3\2\2\2\u00ec\u00ed\7w\2\2\u00ed\u00ee\5? \2\u00ee\u00ef\5? \2\u00ef"+
		"\u00f0\5? \2\u00f0\u00f1\5? \2\u00f1>\3\2\2\2\u00f2\u00f3\t\6\2\2\u00f3"+
		"@\3\2\2\2\u00f4\u00f6\7/\2\2\u00f5\u00f4\3\2\2\2\u00f5\u00f6\3\2\2\2\u00f6"+
		"\u00f7\3\2\2\2\u00f7\u00f8\5C\"\2\u00f8\u00fa\7\60\2\2\u00f9\u00fb\t\7"+
		"\2\2\u00fa\u00f9\3\2\2\2\u00fb\u00fc\3\2\2\2\u00fc\u00fa\3\2\2\2\u00fc"+
		"\u00fd\3\2\2\2\u00fd\u00ff\3\2\2\2\u00fe\u0100\5E#\2\u00ff\u00fe\3\2\2"+
		"\2\u00ff\u0100\3\2\2\2\u0100B\3\2\2\2\u0101\u010a\7\62\2\2\u0102\u0106"+
		"\t\b\2\2\u0103\u0105\t\7\2\2\u0104\u0103\3\2\2\2\u0105\u0108\3\2\2\2\u0106"+
		"\u0104\3\2\2\2\u0106\u0107\3\2\2\2\u0107\u010a\3\2\2\2\u0108\u0106\3\2"+
		"\2\2\u0109\u0101\3\2\2\2\u0109\u0102\3\2\2\2\u010aD\3\2\2\2\u010b\u010d"+
		"\t\t\2\2\u010c\u010e\t\n\2\2\u010d\u010c\3\2\2\2\u010d\u010e\3\2\2\2\u010e"+
		"\u010f\3\2\2\2\u010f\u0110\5C\"\2\u0110F\3\2\2\2\u0111\u0112\7\f\2\2\u0112"+
		"H\3\2\2\2\u0113\u0117\7.\2\2\u0114\u0116\7\"\2\2\u0115\u0114\3\2\2\2\u0116"+
		"\u0119\3\2\2\2\u0117\u0115\3\2\2\2\u0117\u0118\3\2\2\2\u0118J\3\2\2\2"+
		"\u0119\u0117\3\2\2\2\u011a\u011e\7\"\2\2\u011b\u011d\5G$\2\u011c\u011b"+
		"\3\2\2\2\u011d\u0120\3\2\2\2\u011e\u011c\3\2\2\2\u011e\u011f\3\2\2\2\u011f"+
		"L\3\2\2\2\u0120\u011e\3\2\2\2\36\2bit\177\u0087\u008f\u0096\u009d\u00a5"+
		"\u00ad\u00b3\u00b9\u00bf\u00c5\u00cb\u00d1\u00e0\u00e2\u00ea\u00f5\u00fc"+
		"\u00ff\u0106\u0109\u010d\u0117\u011e\2";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}