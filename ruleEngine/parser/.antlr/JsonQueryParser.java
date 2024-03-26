// Generated from d:\project\Syncup\GOlang\lib\git\rule-engine\parser\JsonQuery.g4 by ANTLR 4.8
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class JsonQueryParser extends Parser {
	static { RuntimeMetaData.checkVersion("4.8", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		T__0=1, T__1=2, T__2=3, T__3=4, T__4=5, T__5=6, T__6=7, NOT=8, LOGICAL_OPERATOR=9, 
		BOOLEAN=10, NULL=11, IN=12, EQ=13, NE=14, GT=15, LT=16, GE=17, LE=18, 
		CO=19, SW=20, EW=21, MW=22, ATTRNAME=23, VERSION=24, STRING=25, DOUBLE=26, 
		INT=27, EXP=28, NEWLINE=29, COMMA=30, SP=31;
	public static final int
		RULE_query = 0, RULE_attrPath = 1, RULE_subAttr = 2, RULE_value = 3, RULE_listStrings = 4, 
		RULE_subListOfStrings = 5, RULE_listDoubles = 6, RULE_subListOfDoubles = 7, 
		RULE_listInts = 8, RULE_subListOfInts = 9;
	private static String[] makeRuleNames() {
		return new String[] {
			"query", "attrPath", "subAttr", "value", "listStrings", "subListOfStrings", 
			"listDoubles", "subListOfDoubles", "listInts", "subListOfInts"
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

	@Override
	public String getGrammarFileName() { return "JsonQuery.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public JsonQueryParser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	public static class QueryContext extends ParserRuleContext {
		public QueryContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_query; }
	 
		public QueryContext() { }
		public void copyFrom(QueryContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class CompareExpContext extends QueryContext {
		public Token op;
		public AttrPathContext attrPath() {
			return getRuleContext(AttrPathContext.class,0);
		}
		public List<TerminalNode> SP() { return getTokens(JsonQueryParser.SP); }
		public TerminalNode SP(int i) {
			return getToken(JsonQueryParser.SP, i);
		}
		public ValueContext value() {
			return getRuleContext(ValueContext.class,0);
		}
		public TerminalNode EQ() { return getToken(JsonQueryParser.EQ, 0); }
		public TerminalNode NE() { return getToken(JsonQueryParser.NE, 0); }
		public TerminalNode GT() { return getToken(JsonQueryParser.GT, 0); }
		public TerminalNode LT() { return getToken(JsonQueryParser.LT, 0); }
		public TerminalNode GE() { return getToken(JsonQueryParser.GE, 0); }
		public TerminalNode LE() { return getToken(JsonQueryParser.LE, 0); }
		public TerminalNode CO() { return getToken(JsonQueryParser.CO, 0); }
		public TerminalNode SW() { return getToken(JsonQueryParser.SW, 0); }
		public TerminalNode EW() { return getToken(JsonQueryParser.EW, 0); }
		public TerminalNode IN() { return getToken(JsonQueryParser.IN, 0); }
		public TerminalNode MW() { return getToken(JsonQueryParser.MW, 0); }
		public CompareExpContext(QueryContext ctx) { copyFrom(ctx); }
	}
	public static class ParenExpContext extends QueryContext {
		public QueryContext query() {
			return getRuleContext(QueryContext.class,0);
		}
		public TerminalNode NOT() { return getToken(JsonQueryParser.NOT, 0); }
		public TerminalNode SP() { return getToken(JsonQueryParser.SP, 0); }
		public ParenExpContext(QueryContext ctx) { copyFrom(ctx); }
	}
	public static class PresentExpContext extends QueryContext {
		public AttrPathContext attrPath() {
			return getRuleContext(AttrPathContext.class,0);
		}
		public TerminalNode SP() { return getToken(JsonQueryParser.SP, 0); }
		public PresentExpContext(QueryContext ctx) { copyFrom(ctx); }
	}
	public static class LogicalExpContext extends QueryContext {
		public List<QueryContext> query() {
			return getRuleContexts(QueryContext.class);
		}
		public QueryContext query(int i) {
			return getRuleContext(QueryContext.class,i);
		}
		public List<TerminalNode> SP() { return getTokens(JsonQueryParser.SP); }
		public TerminalNode SP(int i) {
			return getToken(JsonQueryParser.SP, i);
		}
		public TerminalNode LOGICAL_OPERATOR() { return getToken(JsonQueryParser.LOGICAL_OPERATOR, 0); }
		public LogicalExpContext(QueryContext ctx) { copyFrom(ctx); }
	}

	public final QueryContext query() throws RecognitionException {
		return query(0);
	}

	private QueryContext query(int _p) throws RecognitionException {
		ParserRuleContext _parentctx = _ctx;
		int _parentState = getState();
		QueryContext _localctx = new QueryContext(_ctx, _parentState);
		QueryContext _prevctx = _localctx;
		int _startState = 0;
		enterRecursionRule(_localctx, 0, RULE_query, _p);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(41);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,2,_ctx) ) {
			case 1:
				{
				_localctx = new ParenExpContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;

				setState(22);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==NOT) {
					{
					setState(21);
					match(NOT);
					}
				}

				setState(25);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==SP) {
					{
					setState(24);
					match(SP);
					}
				}

				setState(27);
				match(T__0);
				setState(28);
				query(0);
				setState(29);
				match(T__1);
				}
				break;
			case 2:
				{
				_localctx = new PresentExpContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(31);
				attrPath();
				setState(32);
				match(SP);
				setState(33);
				match(T__2);
				}
				break;
			case 3:
				{
				_localctx = new CompareExpContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(35);
				attrPath();
				setState(36);
				match(SP);
				setState(37);
				((CompareExpContext)_localctx).op = _input.LT(1);
				_la = _input.LA(1);
				if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << IN) | (1L << EQ) | (1L << NE) | (1L << GT) | (1L << LT) | (1L << GE) | (1L << LE) | (1L << CO) | (1L << SW) | (1L << EW) | (1L << MW))) != 0)) ) {
					((CompareExpContext)_localctx).op = (Token)_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(38);
				match(SP);
				setState(39);
				value();
				}
				break;
			}
			_ctx.stop = _input.LT(-1);
			setState(50);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,3,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					if ( _parseListeners!=null ) triggerExitRuleEvent();
					_prevctx = _localctx;
					{
					{
					_localctx = new LogicalExpContext(new QueryContext(_parentctx, _parentState));
					pushNewRecursionContext(_localctx, _startState, RULE_query);
					setState(43);
					if (!(precpred(_ctx, 3))) throw new FailedPredicateException(this, "precpred(_ctx, 3)");
					setState(44);
					match(SP);
					setState(45);
					match(LOGICAL_OPERATOR);
					setState(46);
					match(SP);
					setState(47);
					query(4);
					}
					} 
				}
				setState(52);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,3,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			unrollRecursionContexts(_parentctx);
		}
		return _localctx;
	}

	public static class AttrPathContext extends ParserRuleContext {
		public TerminalNode ATTRNAME() { return getToken(JsonQueryParser.ATTRNAME, 0); }
		public SubAttrContext subAttr() {
			return getRuleContext(SubAttrContext.class,0);
		}
		public AttrPathContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_attrPath; }
	}

	public final AttrPathContext attrPath() throws RecognitionException {
		AttrPathContext _localctx = new AttrPathContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_attrPath);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(53);
			match(ATTRNAME);
			setState(55);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==T__3) {
				{
				setState(54);
				subAttr();
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SubAttrContext extends ParserRuleContext {
		public AttrPathContext attrPath() {
			return getRuleContext(AttrPathContext.class,0);
		}
		public SubAttrContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_subAttr; }
	}

	public final SubAttrContext subAttr() throws RecognitionException {
		SubAttrContext _localctx = new SubAttrContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_subAttr);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(57);
			match(T__3);
			setState(58);
			attrPath();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ValueContext extends ParserRuleContext {
		public ValueContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_value; }
	 
		public ValueContext() { }
		public void copyFrom(ValueContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class ListOfDoublesContext extends ValueContext {
		public ListDoublesContext listDoubles() {
			return getRuleContext(ListDoublesContext.class,0);
		}
		public ListOfDoublesContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class ListOfStringsContext extends ValueContext {
		public ListStringsContext listStrings() {
			return getRuleContext(ListStringsContext.class,0);
		}
		public ListOfStringsContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class BooleanContext extends ValueContext {
		public TerminalNode BOOLEAN() { return getToken(JsonQueryParser.BOOLEAN, 0); }
		public BooleanContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class NullContext extends ValueContext {
		public TerminalNode NULL() { return getToken(JsonQueryParser.NULL, 0); }
		public NullContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class StringContext extends ValueContext {
		public TerminalNode STRING() { return getToken(JsonQueryParser.STRING, 0); }
		public StringContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class DoubleContext extends ValueContext {
		public TerminalNode DOUBLE() { return getToken(JsonQueryParser.DOUBLE, 0); }
		public DoubleContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class VersionContext extends ValueContext {
		public TerminalNode VERSION() { return getToken(JsonQueryParser.VERSION, 0); }
		public VersionContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class LongContext extends ValueContext {
		public TerminalNode INT() { return getToken(JsonQueryParser.INT, 0); }
		public TerminalNode EXP() { return getToken(JsonQueryParser.EXP, 0); }
		public LongContext(ValueContext ctx) { copyFrom(ctx); }
	}
	public static class ListOfIntsContext extends ValueContext {
		public ListIntsContext listInts() {
			return getRuleContext(ListIntsContext.class,0);
		}
		public ListOfIntsContext(ValueContext ctx) { copyFrom(ctx); }
	}

	public final ValueContext value() throws RecognitionException {
		ValueContext _localctx = new ValueContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_value);
		int _la;
		try {
			setState(75);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,7,_ctx) ) {
			case 1:
				_localctx = new BooleanContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(60);
				match(BOOLEAN);
				}
				break;
			case 2:
				_localctx = new NullContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(61);
				match(NULL);
				}
				break;
			case 3:
				_localctx = new VersionContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(62);
				match(VERSION);
				}
				break;
			case 4:
				_localctx = new StringContext(_localctx);
				enterOuterAlt(_localctx, 4);
				{
				setState(63);
				match(STRING);
				}
				break;
			case 5:
				_localctx = new DoubleContext(_localctx);
				enterOuterAlt(_localctx, 5);
				{
				setState(64);
				match(DOUBLE);
				}
				break;
			case 6:
				_localctx = new LongContext(_localctx);
				enterOuterAlt(_localctx, 6);
				{
				setState(66);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==T__4) {
					{
					setState(65);
					match(T__4);
					}
				}

				setState(68);
				match(INT);
				setState(70);
				_errHandler.sync(this);
				switch ( getInterpreter().adaptivePredict(_input,6,_ctx) ) {
				case 1:
					{
					setState(69);
					match(EXP);
					}
					break;
				}
				}
				break;
			case 7:
				_localctx = new ListOfIntsContext(_localctx);
				enterOuterAlt(_localctx, 7);
				{
				setState(72);
				listInts();
				}
				break;
			case 8:
				_localctx = new ListOfDoublesContext(_localctx);
				enterOuterAlt(_localctx, 8);
				{
				setState(73);
				listDoubles();
				}
				break;
			case 9:
				_localctx = new ListOfStringsContext(_localctx);
				enterOuterAlt(_localctx, 9);
				{
				setState(74);
				listStrings();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ListStringsContext extends ParserRuleContext {
		public SubListOfStringsContext subListOfStrings() {
			return getRuleContext(SubListOfStringsContext.class,0);
		}
		public ListStringsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_listStrings; }
	}

	public final ListStringsContext listStrings() throws RecognitionException {
		ListStringsContext _localctx = new ListStringsContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_listStrings);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(77);
			match(T__5);
			setState(78);
			subListOfStrings();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SubListOfStringsContext extends ParserRuleContext {
		public TerminalNode STRING() { return getToken(JsonQueryParser.STRING, 0); }
		public TerminalNode COMMA() { return getToken(JsonQueryParser.COMMA, 0); }
		public SubListOfStringsContext subListOfStrings() {
			return getRuleContext(SubListOfStringsContext.class,0);
		}
		public SubListOfStringsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_subListOfStrings; }
	}

	public final SubListOfStringsContext subListOfStrings() throws RecognitionException {
		SubListOfStringsContext _localctx = new SubListOfStringsContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_subListOfStrings);
		try {
			setState(85);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,8,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(80);
				match(STRING);
				setState(81);
				match(COMMA);
				setState(82);
				subListOfStrings();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(83);
				match(STRING);
				setState(84);
				match(T__6);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ListDoublesContext extends ParserRuleContext {
		public SubListOfDoublesContext subListOfDoubles() {
			return getRuleContext(SubListOfDoublesContext.class,0);
		}
		public ListDoublesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_listDoubles; }
	}

	public final ListDoublesContext listDoubles() throws RecognitionException {
		ListDoublesContext _localctx = new ListDoublesContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_listDoubles);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(87);
			match(T__5);
			setState(88);
			subListOfDoubles();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SubListOfDoublesContext extends ParserRuleContext {
		public TerminalNode DOUBLE() { return getToken(JsonQueryParser.DOUBLE, 0); }
		public TerminalNode COMMA() { return getToken(JsonQueryParser.COMMA, 0); }
		public SubListOfDoublesContext subListOfDoubles() {
			return getRuleContext(SubListOfDoublesContext.class,0);
		}
		public SubListOfDoublesContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_subListOfDoubles; }
	}

	public final SubListOfDoublesContext subListOfDoubles() throws RecognitionException {
		SubListOfDoublesContext _localctx = new SubListOfDoublesContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_subListOfDoubles);
		try {
			setState(95);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,9,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(90);
				match(DOUBLE);
				setState(91);
				match(COMMA);
				setState(92);
				subListOfDoubles();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(93);
				match(DOUBLE);
				setState(94);
				match(T__6);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ListIntsContext extends ParserRuleContext {
		public SubListOfIntsContext subListOfInts() {
			return getRuleContext(SubListOfIntsContext.class,0);
		}
		public ListIntsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_listInts; }
	}

	public final ListIntsContext listInts() throws RecognitionException {
		ListIntsContext _localctx = new ListIntsContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_listInts);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(97);
			match(T__5);
			setState(98);
			subListOfInts();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SubListOfIntsContext extends ParserRuleContext {
		public TerminalNode INT() { return getToken(JsonQueryParser.INT, 0); }
		public TerminalNode COMMA() { return getToken(JsonQueryParser.COMMA, 0); }
		public SubListOfIntsContext subListOfInts() {
			return getRuleContext(SubListOfIntsContext.class,0);
		}
		public SubListOfIntsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_subListOfInts; }
	}

	public final SubListOfIntsContext subListOfInts() throws RecognitionException {
		SubListOfIntsContext _localctx = new SubListOfIntsContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_subListOfInts);
		try {
			setState(105);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,10,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(100);
				match(INT);
				setState(101);
				match(COMMA);
				setState(102);
				subListOfInts();
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(103);
				match(INT);
				setState(104);
				match(T__6);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public boolean sempred(RuleContext _localctx, int ruleIndex, int predIndex) {
		switch (ruleIndex) {
		case 0:
			return query_sempred((QueryContext)_localctx, predIndex);
		}
		return true;
	}
	private boolean query_sempred(QueryContext _localctx, int predIndex) {
		switch (predIndex) {
		case 0:
			return precpred(_ctx, 3);
		}
		return true;
	}

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\3!n\4\2\t\2\4\3\t\3"+
		"\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t\t\4\n\t\n\4\13\t\13\3\2"+
		"\3\2\5\2\31\n\2\3\2\5\2\34\n\2\3\2\3\2\3\2\3\2\3\2\3\2\3\2\3\2\3\2\3\2"+
		"\3\2\3\2\3\2\3\2\5\2,\n\2\3\2\3\2\3\2\3\2\3\2\7\2\63\n\2\f\2\16\2\66\13"+
		"\2\3\3\3\3\5\3:\n\3\3\4\3\4\3\4\3\5\3\5\3\5\3\5\3\5\3\5\5\5E\n\5\3\5\3"+
		"\5\5\5I\n\5\3\5\3\5\3\5\5\5N\n\5\3\6\3\6\3\6\3\7\3\7\3\7\3\7\3\7\5\7X"+
		"\n\7\3\b\3\b\3\b\3\t\3\t\3\t\3\t\3\t\5\tb\n\t\3\n\3\n\3\n\3\13\3\13\3"+
		"\13\3\13\3\13\5\13l\n\13\3\13\2\3\2\f\2\4\6\b\n\f\16\20\22\24\2\3\3\2"+
		"\16\30\2v\2+\3\2\2\2\4\67\3\2\2\2\6;\3\2\2\2\bM\3\2\2\2\nO\3\2\2\2\fW"+
		"\3\2\2\2\16Y\3\2\2\2\20a\3\2\2\2\22c\3\2\2\2\24k\3\2\2\2\26\30\b\2\1\2"+
		"\27\31\7\n\2\2\30\27\3\2\2\2\30\31\3\2\2\2\31\33\3\2\2\2\32\34\7!\2\2"+
		"\33\32\3\2\2\2\33\34\3\2\2\2\34\35\3\2\2\2\35\36\7\3\2\2\36\37\5\2\2\2"+
		"\37 \7\4\2\2 ,\3\2\2\2!\"\5\4\3\2\"#\7!\2\2#$\7\5\2\2$,\3\2\2\2%&\5\4"+
		"\3\2&\'\7!\2\2\'(\t\2\2\2()\7!\2\2)*\5\b\5\2*,\3\2\2\2+\26\3\2\2\2+!\3"+
		"\2\2\2+%\3\2\2\2,\64\3\2\2\2-.\f\5\2\2./\7!\2\2/\60\7\13\2\2\60\61\7!"+
		"\2\2\61\63\5\2\2\6\62-\3\2\2\2\63\66\3\2\2\2\64\62\3\2\2\2\64\65\3\2\2"+
		"\2\65\3\3\2\2\2\66\64\3\2\2\2\679\7\31\2\28:\5\6\4\298\3\2\2\29:\3\2\2"+
		"\2:\5\3\2\2\2;<\7\6\2\2<=\5\4\3\2=\7\3\2\2\2>N\7\f\2\2?N\7\r\2\2@N\7\32"+
		"\2\2AN\7\33\2\2BN\7\34\2\2CE\7\7\2\2DC\3\2\2\2DE\3\2\2\2EF\3\2\2\2FH\7"+
		"\35\2\2GI\7\36\2\2HG\3\2\2\2HI\3\2\2\2IN\3\2\2\2JN\5\22\n\2KN\5\16\b\2"+
		"LN\5\n\6\2M>\3\2\2\2M?\3\2\2\2M@\3\2\2\2MA\3\2\2\2MB\3\2\2\2MD\3\2\2\2"+
		"MJ\3\2\2\2MK\3\2\2\2ML\3\2\2\2N\t\3\2\2\2OP\7\b\2\2PQ\5\f\7\2Q\13\3\2"+
		"\2\2RS\7\33\2\2ST\7 \2\2TX\5\f\7\2UV\7\33\2\2VX\7\t\2\2WR\3\2\2\2WU\3"+
		"\2\2\2X\r\3\2\2\2YZ\7\b\2\2Z[\5\20\t\2[\17\3\2\2\2\\]\7\34\2\2]^\7 \2"+
		"\2^b\5\20\t\2_`\7\34\2\2`b\7\t\2\2a\\\3\2\2\2a_\3\2\2\2b\21\3\2\2\2cd"+
		"\7\b\2\2de\5\24\13\2e\23\3\2\2\2fg\7\35\2\2gh\7 \2\2hl\5\24\13\2ij\7\35"+
		"\2\2jl\7\t\2\2kf\3\2\2\2ki\3\2\2\2l\25\3\2\2\2\r\30\33+\649DHMWak";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}