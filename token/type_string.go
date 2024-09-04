// Code generated by "stringer -type=Type"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOL-1]
	_ = x[startValueTokens-2]
	_ = x[IDENT-3]
	_ = x[INT-4]
	_ = x[FLOAT-5]
	_ = x[STRING-6]
	_ = x[LINECOMMENT-7]
	_ = x[BLOCKCOMMENT-8]
	_ = x[endValueTokens-9]
	_ = x[startSingleCharTokens-10]
	_ = x[ASSIGN-11]
	_ = x[PLUS-12]
	_ = x[MINUS-13]
	_ = x[BANG-14]
	_ = x[ASTERISK-15]
	_ = x[SLASH-16]
	_ = x[PERCENT-17]
	_ = x[LT-18]
	_ = x[GT-19]
	_ = x[BITAND-20]
	_ = x[BITOR-21]
	_ = x[BITXOR-22]
	_ = x[BITNOT-23]
	_ = x[COMMA-24]
	_ = x[SEMICOLON-25]
	_ = x[LPAREN-26]
	_ = x[RPAREN-27]
	_ = x[LBRACE-28]
	_ = x[RBRACE-29]
	_ = x[LBRACKET-30]
	_ = x[RBRACKET-31]
	_ = x[COLON-32]
	_ = x[DOT-33]
	_ = x[endSingleCharTokens-34]
	_ = x[startMultiCharTokens-35]
	_ = x[LTEQ-36]
	_ = x[GTEQ-37]
	_ = x[EQ-38]
	_ = x[NOTEQ-39]
	_ = x[INCR-40]
	_ = x[DECR-41]
	_ = x[DOTDOT-42]
	_ = x[OR-43]
	_ = x[AND-44]
	_ = x[LEFTSHIFT-45]
	_ = x[RIGHTSHIFT-46]
	_ = x[LAMBDA-47]
	_ = x[DEFINE-48]
	_ = x[endMultiCharTokens-49]
	_ = x[startIdentityTokens-50]
	_ = x[FUNC-51]
	_ = x[TRUE-52]
	_ = x[FALSE-53]
	_ = x[IF-54]
	_ = x[ELSE-55]
	_ = x[RETURN-56]
	_ = x[FOR-57]
	_ = x[BREAK-58]
	_ = x[CONTINUE-59]
	_ = x[MACRO-60]
	_ = x[QUOTE-61]
	_ = x[UNQUOTE-62]
	_ = x[LEN-63]
	_ = x[FIRST-64]
	_ = x[REST-65]
	_ = x[PRINT-66]
	_ = x[PRINTLN-67]
	_ = x[LOG-68]
	_ = x[ERROR-69]
	_ = x[endIdentityTokens-70]
	_ = x[EOF-71]
}

const _Type_name = "ILLEGALEOLstartValueTokensIDENTINTFLOATSTRINGLINECOMMENTBLOCKCOMMENTendValueTokensstartSingleCharTokensASSIGNPLUSMINUSBANGASTERISKSLASHPERCENTLTGTBITANDBITORBITXORBITNOTCOMMASEMICOLONLPARENRPARENLBRACERBRACELBRACKETRBRACKETCOLONDOTendSingleCharTokensstartMultiCharTokensLTEQGTEQEQNOTEQINCRDECRDOTDOTORANDLEFTSHIFTRIGHTSHIFTLAMBDADEFINEendMultiCharTokensstartIdentityTokensFUNCTRUEFALSEIFELSERETURNFORBREAKCONTINUEMACROQUOTEUNQUOTELENFIRSTRESTPRINTPRINTLNLOGERRORendIdentityTokensEOF"

var _Type_index = [...]uint16{0, 7, 10, 26, 31, 34, 39, 45, 56, 68, 82, 103, 109, 113, 118, 122, 130, 135, 142, 144, 146, 152, 157, 163, 169, 174, 183, 189, 195, 201, 207, 215, 223, 228, 231, 250, 270, 274, 278, 280, 285, 289, 293, 299, 301, 304, 313, 323, 329, 335, 353, 372, 376, 380, 385, 387, 391, 397, 400, 405, 413, 418, 423, 430, 433, 438, 442, 447, 454, 457, 462, 479, 482}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
