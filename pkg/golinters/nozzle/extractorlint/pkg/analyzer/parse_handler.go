package analyzer

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"
)

//nolint:gocyclo
func parseHandler(genDecl *ast.GenDecl) *handlerLint {
	isValidHandler, spec, compLit := verifyHandlerDecl(genDecl)
	if !isValidHandler {
		return nil
	}

	hl := handlerLint{}
	hl.Pos = genDecl.TokPos
	hl.Name = spec.Names[0].Name

	for _, elt := range compLit.Elts {
		kvExpr, ok := elt.(*ast.KeyValueExpr)
		if !ok || kvExpr.Key == nil {
			continue
		}

		key, ok := kvExpr.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch key.Name {
		case "Reg":
			hl.Reg = newHandlerFieldReg(key.NamePos, kvExpr.Value)
		case "ID":
			hl.ID = newHandlerField(key.Name, key.NamePos)
		case "IDFunc":
			hl.IDFunc = newHandlerField(key.Name, key.NamePos)
		case "IDMatchedFunc":
			hl.IDMatchedFunc = newHandlerField(key.Name, key.NamePos)
		case "AfterID":
			hl.AfterID = newHandlerField(key.Name, key.NamePos)
		case "AfterIDFunc":
			hl.AfterIDFunc = newHandlerField(key.Name, key.NamePos)
		case "OnText":
			hl.OnText = newHandlerField(key.Name, key.NamePos)
		case "OnTextFunc":
			hl.OnTextFunc = newHandlerField(key.Name, key.NamePos)
		case "Pop":
			hl.Pop = newHandlerField(key.Name, key.NamePos)
		case "PopFunc":
			hl.PopFunc = newHandlerField(key.Name, key.NamePos)
		case "ProbableHandlers":
			hl.ProbableHandlers = newHandlerField(key.Name, key.NamePos)
		default:
			continue
		}
	}

	return &hl
}

func verifyHandlerDecl(genDecl *ast.GenDecl) (bool, *ast.ValueSpec, *ast.CompositeLit) {
	if len(genDecl.Specs) < 1 {
		return false, nil, nil
	}

	spec, ok := genDecl.Specs[0].(*ast.ValueSpec)
	if !ok || len(spec.Values) < 1 {
		return false, nil, nil
	}

	compLit, ok := spec.Values[0].(*ast.CompositeLit)
	if !ok || compLit.Type == nil {
		return false, nil, nil
	}

	compLitType, ok := compLit.Type.(*ast.SelectorExpr)
	if !ok || compLitType.Sel.Name != "ElementHandler" {
		return false, nil, nil
	}

	if len(spec.Names) < 1 {
		return false, nil, nil
	}

	if len(compLit.Elts) == 0 {
		// this is an empty ElementHandler
		return false, nil, nil
	}

	return true, spec, compLit
}

func newHandlerFieldReg(hfPos token.Pos, keyVal ast.Expr) *handlerFieldReg {
	hfReg := handlerFieldReg{
		handlerField: newHandlerField("Reg", hfPos),
	}

	regCallExpr, ok := keyVal.(*ast.CallExpr)
	if !ok || len(regCallExpr.Args) == 0 {
		hfReg.error(errors.New("could not get call expr"))
		return &hfReg
	}

	regNameArg, ok := regCallExpr.Args[0].(*ast.BasicLit)
	if !ok {
		hfReg.error(errors.New("could not get call arg"))
		return &hfReg
	}

	// get the string used to Register this handler
	hfReg.Name = strings.Trim(regNameArg.Value, `"`)
	hfReg.NamePos = regNameArg.ValuePos + token.Pos(1)

	return &hfReg
}
