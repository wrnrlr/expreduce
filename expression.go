package cas

import "bytes"
import "fmt"

type Expression struct {
	Parts []Ex
}

func HeadAssertion(ex Ex, head string) (*Expression, bool) {
	expr, isExpr := ex.(*Expression)
	if isExpr {
		sym, isSym := expr.Parts[0].(*Symbol)
		if isSym {
			if sym.Name == head {
				return expr, true
			}
		}
	}
	return &Expression{}, false
}

func (this *Expression) Eval(es *EvalState) Ex {
	// Start by evaluating each argument
	headSym, headIsSym := &Symbol{}, false
	if len(this.Parts) > 0 {
		headSym, headIsSym = this.Parts[0].(*Symbol)
	}
	for i := range this.Parts {
		if headIsSym && i == 1 && IsHoldFirst(headSym) {
			continue
		}
		//if headIsSym && IsAttribute(headSym, "HoldAll", es) {
		if headIsSym && IsHoldAll(headSym) {
			continue
		}
		this.Parts[i] = this.Parts[i].Eval(es)
	}

	// If any of the parts are Sequence, merge them with parts
	// TODO: I should not be attempting to merge the head if it happens to be
	// a Sequence type
	origLen := len(this.Parts)
	offset := 0
	for i := 0; i < origLen; i++ {
		j := i + offset
		e := this.Parts[j]
		seq, isseq := HeadAssertion(e, "Sequence")
		if isseq {
			start := j
			end := j + 1
			if j == 0 {
				this.Parts = append(seq.Parts[1:], this.Parts[end:]...)
			} else if j == len(this.Parts)-1 {
				this.Parts = append(this.Parts[:start], seq.Parts[1:]...)
			} else {
				// All of these deep copies may not be needed.
				this.Parts = append(append(this.DeepCopy().(*Expression).Parts[:start], seq.DeepCopy().(*Expression).Parts[1:]...), this.DeepCopy().(*Expression).Parts[end:]...)
			}
			offset += len(seq.Parts[1:]) - 1
		}
	}

	headAsSym, isHeadSym := this.Parts[0].(*Symbol)
	if isHeadSym {
		headStr := headAsSym.Name
		if headStr == "Power" {
			return this.EvalPower(es)
		}
		if headStr == "Equal" {
			return this.EvalEqual(es)
		}
		if headStr == "SameQ" {
			return this.EvalSameQ(es)
		}
		if headStr == "Plus" {
			return this.EvalPlus(es)
		}
		if headStr == "Times" {
			return this.EvalTimes(es)
		}
		if headStr == "Set" {
			return this.EvalSet(es)
		}
		if headStr == "SetDelayed" {
			return this.EvalSetDelayed(es)
		}
		if headStr == "If" {
			return this.EvalIf(es)
		}
		if headStr == "While" {
			return this.EvalWhile(es)
		}
		if headStr == "MatchQ" {
			return this.EvalMatchQ(es)
		}
		if headStr == "Replace" {
			return this.EvalReplace(es)
		}
		if headStr == "ReplaceRepeated" {
			return this.EvalReplaceRepeated(es)
		}
		if headStr == "BasicSimplify" {
			return this.EvalBasicSimplify(es)
		}
		if headStr == "SetLogging" {
			return this.EvalSetLogging(es)
		}
		if headStr == "Definition" {
			return this.EvalDefinition(es)
		}
		if headStr == "Order" {
			return this.EvalOrder(es)
		}
		if headStr == "Sort" {
			return this.EvalSort(es)
		}
		if headStr == "RandomReal" {
			return this.EvalRandomReal(es)
		}
		if headStr == "SeedRandom" {
			return this.EvalSeedRandom(es)
		}
		if headStr == "UnixTime" {
			return this.EvalUnixTime(es)
		}
		if headStr == "Apply" {
			return this.EvalApply(es)
		}
		if headStr == "Length" {
			return this.EvalLength(es)
		}
		if headStr == "Table" {
			return this.EvalTable(es)
		}
		if headStr == "Clear" {
			return this.EvalClear(es)
		}
		if headStr == "Timing" {
			return this.EvalTiming(es)
		}
		if headStr == "MemberQ" {
			return this.EvalMemberQ(es)
		}

		theRes, isDefined := es.GetDef(headStr, this)
		if isDefined {
			return theRes
		}
	}
	return this
}

func (this *Expression) Replace(r *Expression, es *EvalState) Ex {
	oldVars := es.GetDefinedSnapshot()
	es.log.Debugf(es.Pre() + "In Expression.Replace. First trying IsMatchQ(this, r.Parts[1], es).")
	es.log.Debugf(es.Pre()+"Rule r is: %s", r.ToString())

	matchq := IsMatchQ(this, r.Parts[1], es)
	toreturn := r.Parts[2].DeepCopy().Eval(es)
	es.ClearPD()
	es.defined = oldVars
	if matchq {
		es.log.Debugf(es.Pre()+"After MatchQ, rule is: %s", r.ToString())
		es.log.Debugf(es.Pre()+"MatchQ succeeded. Returning r.Parts[2]: %s", r.Parts[2].ToString())
		return toreturn
	}

	thisSym, thisSymOk := this.Parts[0].(*Symbol)
	lhsExpr, lhsExprOk := r.Parts[1].(*Expression)
	if lhsExprOk {
		otherSym, otherSymOk := lhsExpr.Parts[0].(*Symbol)
		if thisSymOk && otherSymOk {
			if thisSym.Name == otherSym.Name {
				if IsOrderless(thisSym) {
					es.log.Debugf(es.Pre() + "r.Parts[1] is Orderless. Now running CommutativeReplace")
					replaced := this.Parts[1:len(this.Parts)]
					CommutativeReplace(&replaced, lhsExpr.Parts[1:len(lhsExpr.Parts)], r.Parts[2], es)
					this.Parts = this.Parts[0:1]
					this.Parts = append(this.Parts, replaced...)
				}
			}
		}
	}

	for i := range this.Parts {
		this.Parts[i] = this.Parts[i].Replace(r, es)
	}
	return this.Eval(es)
}

func (this *Expression) ToString() string {
	headAsSym, isHeadSym := this.Parts[0].(*Symbol)
	if isHeadSym {
		headStr := headAsSym.Name
		if headStr == "Times" {
			return this.ToStringTimes()
		} else if headStr == "Plus" {
			return this.ToStringPlus()
		} else if headStr == "Power" {
			return this.ToStringPower()
		} else if headStr == "Equal" {
			return this.ToStringEqual()
		} else if headStr == "SameQ" {
			return this.ToStringSameQ()
		} else if headStr == "Replace" {
			return this.ToStringReplace()
		} else if headStr == "ReplaceRepeated" {
			return this.ToStringReplaceRepeated()
		} else if headStr == "Pattern" {
			return this.ToStringPattern()
		} else if headStr == "Blank" {
			return this.ToStringBlank()
		} else if headStr == "BlankSequence" {
			return this.ToStringBlankSequence()
		} else if headStr == "BlankNullSequence" {
			return this.ToStringBlankNullSequence()
		} else if headStr == "Rule" {
			return this.ToStringRule()
		} else if headStr == "Set" {
			return this.ToStringSet()
		} else if headStr == "SetDelayed" {
			return this.ToStringSetDelayed()
		} else if headStr == "List" {
			return this.ToStringList()
		}
	}

	// Default printing format
	var buffer bytes.Buffer
	buffer.WriteString(this.Parts[0].ToString())
	buffer.WriteString("[")
	for i, e := range this.Parts {
		if i == 0 {
			continue
		}
		buffer.WriteString(e.ToString())
		if i != len(this.Parts)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}

func IsAttribute(sm *Symbol, attr string, es *EvalState) bool {
	if sm.Name == "MemberQ" {
		return attr == "Protected"
	} else if sm.Name == "Attributes" {
		return attr == "Protected" || attr == "HoldAll" || attr == "Listable"
	} else if sm.Name == "List" {
		return attr == "Protected" || attr == "Locked"
	} else if sm.Name == "Pattern" {
		return attr == "Protected" || attr == "HoldFirst"
	} else if sm.Name == "Blank" {
		return attr == "Protected"
	} else if sm.Name == "Rule" {
		return attr == "Protected" || attr == "SequenceHold"
	} else if sm.Name == "Times" || sm.Name == "Plus" {
		return attr == "Flat" || attr == "Listable" || attr == "NumericFunction" || attr == "OneIdentity" || attr == "Orderless" || attr == "Protected"
	} else if sm.Name == "Power" {
		return attr == "Listable" || attr == "NumericFunction" || attr == "OneIdentity" || attr == "Protected"
	} else if sm.Name == "ReplaceRepeated" {
		return attr == "Protected"
	} else if sm.Name == "Equal" {
		return attr == "Protected"
	} // This is probably slow because it requires tons of Defined copies
	fmt.Printf("IsAttribute(%v, %v)\n", sm.Name, attr)
	res := (&Expression{[]Ex{
		&Symbol{"MemberQ"},
		&Expression{[]Ex{
			&Symbol{"Attributes"},
			sm,
		}},
		&Symbol{attr},
	}}).Eval(es)
	resSym, resIsSym := res.(*Symbol)
	if resIsSym {
		return resSym.Name == "True"
	}
	return false
}

// TODO: convert to a map
func IsOrderless(sym *Symbol) bool {
	if sym.Name == "Times" {
		return true
	} else if sym.Name == "Plus" {
		return true
	}
	return false
}

// TODO: convert to a map
func IsHoldFirst(sym *Symbol) bool {
	if sym.Name == "Set" {
		return true
	} else if sym.Name == "Pattern" {
		return true
	}
	return false
}

// TODO: convert to a map
func IsHoldAll(sym *Symbol) bool {
	if sym.Name == "SetDelayed" {
		return true
	}
	if sym.Name == "Table" {
		return true
	}
	if sym.Name == "Clear" {
		return true
	}
	if sym.Name == "Timing" {
		return true
	}
	if sym.Name == "Hold" {
		return true
	}
	return false
}

func (this *Expression) IsEqual(otherEx Ex, es *EvalState) string {
	thisEvaled := this.Eval(es)
	otherEx = otherEx.Eval(es)
	this, ok := thisEvaled.(*Expression)
	if !ok {
		return thisEvaled.IsEqual(otherEx, es)
	}

	other, ok := otherEx.(*Expression)
	if !ok {
		return "EQUAL_UNK"
	}

	thisSym, thisSymOk := this.Parts[0].(*Symbol)
	otherSym, otherSymOk := other.Parts[0].(*Symbol)
	if thisSymOk && otherSymOk {
		if thisSym.Name == otherSym.Name {
			if IsOrderless(thisSym) {
				return CommutativeIsEqual(this.Parts[1:len(this.Parts)], other.Parts[1:len(other.Parts)], es)
			}
		}
	}

	return FunctionIsEqual(this.Parts, other.Parts, es)
}

func (this *Expression) IsSameQ(otherEx Ex, es *EvalState) bool {
	other, ok := otherEx.(*Expression)
	if !ok {
		return false
	}

	thisSym, thisSymOk := this.Parts[0].(*Symbol)
	otherSym, otherSymOk := other.Parts[0].(*Symbol)
	if thisSymOk && otherSymOk {
		if thisSym.Name == otherSym.Name {
			if IsOrderless(thisSym) {
				return this.IsEqual(otherEx, es) == "EQUAL_TRUE"
			}
		}
	}

	return FunctionIsSameQ(this.Parts, other.Parts, es)
}

func (this *Expression) IsMatchQ(otherEx Ex, es *EvalState) bool {
	headAsSym, isHeadSym := this.Parts[0].(*Symbol)
	if isHeadSym {
		headStr := headAsSym.Name
		if IsBlankTypeOnly(otherEx) {
			if IsBlankTypeCapturing(otherEx, this, headStr, es) {
				return true
			}
			return false
		}
	}
	//thisEx := this.Eval(es)
	//otherEx = otherEx.Eval(es)
	//this, ok := thisEx.(*Expression)
	//if !ok {
	//return IsMatchQ(thisEx, otherEx, es)
	//}
	other, otherOk := otherEx.(*Expression)
	if !otherOk {
		return false
	}

	thisSym, thisSymOk := this.Parts[0].(*Symbol)
	otherSym, otherSymOk := other.Parts[0].(*Symbol)
	if thisSymOk && otherSymOk {
		if thisSym.Name == otherSym.Name {
			if IsOrderless(thisSym) {
				return CommutativeIsMatchQ(this.Parts[1:len(this.Parts)], other.Parts[1:len(other.Parts)], es)
			}
		}
	}

	return NonCommutativeIsMatchQ(this.Parts, other.Parts, es)
}

func (this *Expression) DeepCopy() Ex {
	var thiscopy = &Expression{}
	for i := range this.Parts {
		thiscopy.Parts = append(thiscopy.Parts, this.Parts[i].DeepCopy())
	}
	return thiscopy
}

// Implement the sort.Interface
func (this *Expression) Len() int {
	return len(this.Parts) - 1
}

func (this *Expression) Less(i, j int) bool {
	return ExOrder(this.Parts[i+1], this.Parts[j+1]) == 1
}

func (this *Expression) Swap(i, j int) {
	this.Parts[j+1], this.Parts[i+1] = this.Parts[i+1], this.Parts[j+1]
}

// Apply
func (this *Expression) EvalApply(es *EvalState) Ex {
	if len(this.Parts) != 3 {
		return this
	}

	sym, isSym := this.Parts[1].(*Symbol)
	list, isList := HeadAssertion(this.Parts[2].DeepCopy(), "List")
	if isSym && isList {
		toReturn := &Expression{[]Ex{sym}}
		toReturn.Parts = append(toReturn.Parts, list.Parts[1:]...)
		return toReturn.Eval(es)
	}
	return this
}
